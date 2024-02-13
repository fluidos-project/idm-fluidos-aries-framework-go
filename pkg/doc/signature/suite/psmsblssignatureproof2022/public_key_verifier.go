/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignatureproof2022

import "github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"

const g1PubKeyType = "Bls12381G1Key2022"

// NewG1PublicKeyVerifier creates a signature verifier that verifies a PsmsBlsSignatureProof2022 signature
// taking Bls12381G1Key2022 public key bytes as input.
func NewG1PublicKeyVerifier(nonce []byte) *verifier.PublicKeyVerifier {
	return verifier.NewPublicKeyVerifier(verifier.NewPSMSG1SignatureProofVerifier(nonce),
		verifier.WithExactPublicKeyType(g1PubKeyType))
}
