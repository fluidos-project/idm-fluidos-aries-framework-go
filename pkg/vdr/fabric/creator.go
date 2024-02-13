/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/common/model"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
	"strings"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
)

const (
	schemaResV1                = "https://w3id.org/did-resolution/v1"
	ed25519VerificationKey2018 = "Ed25519VerificationKey2018"
	bls12381G1Key2022		   = "Bls12381G1Key2022"
	jsonWebKey2020             = "JsonWebKey2020"
	x25519KeyAgreementKey2019  = "X25519KeyAgreementKey2019"
)

// modifiedBy key/signature used to update the DID Document.
type modifiedBy struct {
	Key string `json:"key,omitempty"`
	Sig string `json:"sig,omitempty"`
}

type docDelta struct {
	Change     string        `json:"change,omitempty"`
	ModifiedBy *[]modifiedBy `json:"by,omitempty"`
	ModifiedAt time.Time     `json:"when,omitempty"`
}

// Create create new DID Document.
// TODO https://github.com/hyperledger/aries-framework-go/issues/2466
func (v *VDR) Create(didDoc *did.Doc, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
	docOpts := &vdrapi.DIDMethodOpts{Values: make(map[string]interface{})}
	// Apply options
	for _, opt := range opts {
		opt(docOpts)
	}

	store := false

	storeOpt := docOpts.Values["store"]
	if storeOpt != nil {
		var ok bool

		store, ok = storeOpt.(bool)
		if !ok {
			return nil, fmt.Errorf("store opt not boolean")
		}
	}

	if !store {
		docResolution, err := build(didDoc, docOpts)
		if err != nil {
			return nil, fmt.Errorf("create peer fabric DID : %w", err)
		}

		didDoc = docResolution.DIDDocument
	}

	if err := v.storeDID(didDoc, nil); err != nil {
		return nil, err
	}

	return &did.DocResolution{Context: []string{schemaResV1}, DIDDocument: didDoc}, nil
}

//nolint: funlen,gocyclo
func build(didDoc *did.Doc, docOpts *vdrapi.DIDMethodOpts) (*did.DocResolution, error) {
	if len(didDoc.VerificationMethod) == 0 && len(didDoc.KeyAgreement) == 0 {
		return nil, fmt.Errorf("verification method and key agreement are empty, at least one should be set")
	}

	mainVM, keyAgreementVM, err := buildDIDVMs(didDoc)

	if err != nil {
		return nil, err
	}

	// Service model to be included only if service type is provided through opts
	var service []did.Service

	for i := range didDoc.Service {
		if didDoc.Service[i].ID == "" {
			didDoc.Service[i].ID = uuid.New().String()
		}

		if didDoc.Service[i].Type == "" && docOpts.Values[DefaultServiceType] != nil {
			v, ok := docOpts.Values[DefaultServiceType].(string)
			if !ok {
				return nil, fmt.Errorf("defaultServiceType not string")
			}

			didDoc.Service[i].Type = v
		}

		uri, _ := didDoc.Service[i].ServiceEndpoint.URI() // nolint:errcheck

		// nolint:nestif
		if uri == "" && docOpts.Values[DefaultServiceEndpoint] != nil {
			switch didDoc.Service[i].Type {
			case vdrapi.DIDCommServiceType, "IndyAgent":
				v, ok := docOpts.Values[DefaultServiceEndpoint].(string)
				if !ok {
					return nil, fmt.Errorf("defaultServiceEndpoint not string")
				}

				didDoc.Service[i].ServiceEndpoint = model.NewDIDCommV1Endpoint(v)
			case vdrapi.DIDCommV2ServiceType:
				epArrayEntry := stringArray(docOpts.Values[DefaultServiceEndpoint])

				sp := model.Endpoint{}

				if len(epArrayEntry) == 0 {
					sp = model.NewDIDCommV2Endpoint([]model.DIDCommV2Endpoint{{}})
				} else {
					for _, ep := range epArrayEntry {
						err = sp.UnmarshalJSON([]byte(ep))
						if err != nil {
							if strings.EqualFold(err.Error(), "endpoint data is not supported") {
								// if unmarshall failed, then use as string.
								sp = model.NewDIDCommV2Endpoint([]model.DIDCommV2Endpoint{
									{URI: ep, Accept: []string{transport.MediaTypeDIDCommV2Profile}},
								})
							}

							continue
						}

						break
					}
				}

				didDoc.Service[i].ServiceEndpoint = sp
			}
		}
		applyDIDCommKeys(i, didDoc)
		applyDIDCommV2Keys(i, didDoc)

		service = append(service, didDoc.Service[i])
	}

	// Created/Updated time
	t := time.Now()

	var assertion []did.Verification
	if(len(mainVM) > 1){
		assertion = []did.Verification{{
			VerificationMethod: mainVM[1],
			Relationship:       did.AssertionMethod,
		}}
	}

	authentication := []did.Verification{{
		VerificationMethod: mainVM[0],
		Relationship:       did.Authentication,
	}}

	var keyAgreement []did.Verification

	verificationMethods := mainVM

	if keyAgreementVM != nil {
		verificationMethods = append(verificationMethods, keyAgreementVM...)

		for _, ka := range keyAgreementVM {
			keyAgreement = append(keyAgreement, did.Verification{
				VerificationMethod: ka,
				Relationship:       did.KeyAgreement,
			})
		}
	}
	var realID = didDoc.ID

	didDoc, err = NewDoc(
		verificationMethods,
		did.WithService(service),
		did.WithCreatedTime(t),
		did.WithUpdatedTime(t),
		did.WithAuthentication(authentication),
		did.WithAssertion(assertion),
		did.WithKeyAgreement(keyAgreement),
	)
	if err != nil {
		return nil, err
	}

	didDoc.ID = realID
	return &did.DocResolution{DIDDocument: didDoc}, nil
}
func buildDIDVMs(didDoc *did.Doc) ([]did.VerificationMethod, []did.VerificationMethod, error) {
	var mainVM, keyAgreementVM []did.VerificationMethod

	// add all VMs, not only the first one.
	for _, vm := range didDoc.VerificationMethod {
		switch vm.Type {
		case ed25519VerificationKey2018:
			mainVM = append(mainVM, *did.NewVerificationMethodFromBytes(vm.ID, ed25519VerificationKey2018,
				vm.Controller, vm.Value))
		case bls12381G1Key2022:
			mainVM = append(mainVM, *did.NewVerificationMethodFromBytes(vm.ID, bls12381G1Key2022,
				vm.Controller, vm.Value))
		case jsonWebKey2020:
			publicKey1, err := did.NewVerificationMethodFromJWK(vm.ID, jsonWebKey2020, "#id",
				vm.JSONWebKey())
			if err != nil {
				return nil, nil, err
			}

			mainVM = append(mainVM, *publicKey1)
		default:
			return nil, nil, fmt.Errorf("not supported VerificationMethod public key type: %s",
				didDoc.VerificationMethod[0].Type)
		}
	}

	for _, ka := range didDoc.KeyAgreement {
		switch ka.VerificationMethod.Type {
		case x25519KeyAgreementKey2019:
			keyAgreementVM = append(keyAgreementVM, *did.NewVerificationMethodFromBytes(
				ka.VerificationMethod.ID, x25519KeyAgreementKey2019, "",
				ka.VerificationMethod.Value))

		case jsonWebKey2020:
			kaVM, err := did.NewVerificationMethodFromJWK(ka.VerificationMethod.ID, jsonWebKey2020, "",
				ka.VerificationMethod.JSONWebKey())
			if err != nil {
				return nil, nil, err
			}

			keyAgreementVM = append(keyAgreementVM, *kaVM)
		default:
			return nil, nil, fmt.Errorf("not supported KeyAgreement public key type: %s", didDoc.VerificationMethod[0].Type)
		}
	}

	return mainVM, keyAgreementVM, nil
}

func genesisDeltaBytes(doc *did.Doc, by *[]modifiedBy) ([]byte, error) {
	var deltas []docDelta

	// For now, assume the doc is a genesis document
	jsonDoc, err := doc.JSONBytes()
	if err != nil {
		return nil, fmt.Errorf("JSON marshalling of document failed: %w", err)
	}

	docDelta := &docDelta{
		Change:     base64.URLEncoding.EncodeToString(jsonDoc),
		ModifiedBy: by,
		ModifiedAt: time.Now(),
	}

	deltas = append(deltas, *docDelta)

	val, err := json.Marshal(deltas)
	if err != nil {
		return nil, fmt.Errorf("JSON marshalling of document deltas failed: %w", err)
	}

	return val, nil
}

// storeDID saves DID Document along with user key/signature using fabric DLT.
func (v *VDR) storeDID(doc *did.Doc, by *[]modifiedBy) error { //nolint: unparam
	if doc == nil || doc.ID == "" {
		return errors.New("DID and document are mandatory")
	}

	val, err := genesisDeltaBytes(doc, by)
	if err != nil {
		return err
	}

	v.connectGateway()
	_, err = v.contract.SubmitTransaction(SC_METHOD_CREATE_DID, string(doc.ID), string(val))
	//s := string(result)
	if err != nil {
		return fmt.Errorf("Failed to submit transaction (storeDID): %s\n", err)
	}
	return nil

}

// stringArray.
func stringArray(entry interface{}) []string {
	if entry == nil {
		return nil
	}

	entries, ok := entry.([]interface{})
	if !ok {
		if entryStr, ok := entry.(string); ok {
			return []string{entryStr}
		}

		return nil
	}

	var result []string

	for _, e := range entries {
		if e != nil {
			result = append(result, stringEntry(e))
		}
	}

	return result
}

// stringEntry.
func stringEntry(entry interface{}) string {
	if entry == nil {
		return ""
	}

	return entry.(string)
}

func applyDIDCommKeys(i int, didDoc *did.Doc) {
	if didDoc.Service[i].Type == vdrapi.DIDCommServiceType {
		didKey, _ := fingerprint.CreateDIDKey(didDoc.VerificationMethod[0].Value)
		didDoc.Service[i].RecipientKeys = []string{didKey}
		didDoc.Service[i].Priority = 0
	}
}

func applyDIDCommV2Keys(i int, didDoc *did.Doc) {
	if didDoc.Service[i].Type == vdrapi.DIDCommV2ServiceType {
		didDoc.Service[i].RecipientKeys = []string{}
		didDoc.Service[i].Priority = 0

		for _, ka := range didDoc.KeyAgreement {
			kaID := ka.VerificationMethod.ID

			didDoc.Service[i].RecipientKeys = append(didDoc.Service[i].RecipientKeys, kaID)
		}
	}
}
