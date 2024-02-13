/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignature2022

import "github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"

const g1PubKeyType = "Bls12381G1Key2022"

// NewG1PublicKeyVerifier creates a signature verifier that verifies a BbsBlsSignature2020 signature
// taking Bls12381G2Key2020 public key bytes as input.
func NewG1PublicKeyVerifier() *verifier.PublicKeyVerifier {
	return verifier.NewPublicKeyVerifier(verifier.NewPSMSG1SignatureVerifier(),
		verifier.WithExactPublicKeyType(g1PubKeyType))
}
