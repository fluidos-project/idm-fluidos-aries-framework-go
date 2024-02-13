/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package verifiable

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignatureproof2022"

	jsonutil "github.com/hyperledger/aries-framework-go/pkg/doc/util/json"
)

// GeneratePSMSSelectiveDisclosure generate PSMS+ selective disclosure from one PSMS+ signature.
func (vc *Credential) GeneratePSMSSelectiveDisclosure(revealDoc map[string]interface{},
	nonce []byte, opts ...CredentialOpt) (*Credential, error) {
	if len(nonce) == 0 { //TODO UMU (verifiable/issue) Use current date as "BAD" replay defense
		nonce = []byte(time.Now().String())
	}
	if len(vc.Proofs) == 0 {
		return nil, errors.New("expected at least one proof present")
	}

	vcOpts := getCredentialOpts(opts)
	jsonldProcessorOpts := mapJSONLDProcessorOpts(&vcOpts.jsonldCredentialOpts)

	if vcOpts.publicKeyFetcher == nil {
		return nil, errors.New("public key fetcher is not defined")
	}

	suite := psmsblssignatureproof2022.New()

	vcDoc, err := jsonutil.ToMap(vc)
	if err != nil {
		return nil, err
	}

	keyResolver := &keyResolverAdapter{vcOpts.publicKeyFetcher}

	vcWithSelectiveDisclosureDoc, err := suite.SelectiveDisclosure(vcDoc, revealDoc, nonce,
		keyResolver, jsonldProcessorOpts...)
	if err != nil {
		return nil, fmt.Errorf("create VC selective disclosure: %w", err)
	}

	vcWithSelectiveDisclosureBytes, err := json.Marshal(vcWithSelectiveDisclosureDoc)
	if err != nil {
		return nil, err
	}

	opts = append(opts, WithDisabledProofCheck())
	return ParseCredential(vcWithSelectiveDisclosureBytes, opts...)
}
