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
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

const (
	// VersionIDOpt version id opt this option is not mandatory.
	VersionIDOpt = "versionID"
	// VersionTimeOpt version time opt this option is not mandatory.
	VersionTimeOpt = "versionTime"
	didLDJson      = "application/did+ld+json"
)

// resolveDID makes DID resolution via FABRIC.
func (v *VDR) resolveDID(didID string) ([]byte, error) {
	v.connectGateway()
	result, err := v.contract.EvaluateTransaction(SC_METHOD_RESOLVE_DID, didID)
	//s := string(result)
	if err != nil {
		return nil, fmt.Errorf("Failed to evaluate transaction: %s\n", err)
	}
	return result, nil
}

// Read implements didresolver.DidMethod.Read interface (https://w3c-ccg.github.io/did-resolution/#resolving-input)
func (v *VDR) Read(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) { //nolint: funlen,gocyclo
	// get the document from the store
	doc, err := v.Get(didID)
	if err != nil {
		return nil, fmt.Errorf("fetching data from store failed: %w", err)
	}

	if doc == nil {
		return nil, vdrapi.ErrNotFound
	}

	return &did.DocResolution{Context: []string{schemaResV1}, DIDDocument: doc}, nil
}

// Get returns Peer DID Document.
func (v *VDR) Get(id string) (*did.Doc, error) {
	if id == "" {
		return nil, errors.New("ID is mandatory")
	}

	deltas, err := v.getDeltas(id)
	if err != nil {
		return nil, fmt.Errorf("delta data fetch from store for did [%s] failed: %w", id, err)
	}

	return assembleDocFromDeltas(deltas)
}
func (v *VDR) getDeltas(id string) ([]docDelta, error) {
	val, err := v.resolveDID(id)
	if errors.Is(err, storage.ErrDataNotFound) {
		return nil, vdrapi.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("fetching data from store failed: %w", err)
	}

	var deltas []docDelta

	err = json.Unmarshal(val, &deltas)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshalling of document deltas failed: %w", err)
	}

	return deltas, nil
}
func assembleDocFromDeltas(deltas []docDelta) (*did.Doc, error) {
	// For now, assume storage contains only one delta(genesis document)
	delta := deltas[0]

	doc, err := base64.URLEncoding.DecodeString(delta.Change)
	if err != nil {
		return nil, fmt.Errorf("decoding of document delta failed: %w", err)
	}

	document, err := did.ParseDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("document ParseDocument() failed: %w", err)
	}

	return document, nil
}
