/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package subtlepsms

import (
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
)

// BLS12381G1Verifier is the PSMS signature/proof verifier for keys on BLS12-381 curve with a point in the G1 group.
// Currently this is the only available PSMS+ verifier in aries-framework-go (see `pkg/doc/psms/psms12381g1pub/psms.go`).
type BLS12381G1Verifier struct {
	signerPubKeyBytes []byte
	psmsPrimitive     *psms12381g1pub.PSMSG1Pub
}

// NewBLS12381G1Verifier creates a new instance of BLS12381G1Verifier with the provided signerPublicKey.
func NewBLS12381G1Verifier(signerPublicKey []byte) *BLS12381G1Verifier {
	return &BLS12381G1Verifier{
		signerPubKeyBytes: signerPublicKey,
		psmsPrimitive:     psms12381g1pub.New(),
	}
}

// Verify will verify an aggregated signature of one or more messages against the signer's public key.
// returns:
// 		error in case of errors or nil if signature verification was successful
func (v *BLS12381G1Verifier) Verify(messages [][]byte, signature []byte) error {
	return v.psmsPrimitive.Verify(messages, signature, v.signerPubKeyBytes)
}

// VerifyProof will verify a PSMS+ signature proof (generated e.g. by DeriveProof()) with the signer's public key.
// returns:
// 		error in case of errors or nil if signature proof verification was successful
func (v *BLS12381G1Verifier) VerifyProof(messages [][]byte, proof, nonce []byte) error {
	return v.psmsPrimitive.VerifyProof(messages, proof, nonce, v.signerPubKeyBytes)
}

// DeriveProof will create a PSMS+ signature proof for a list of revealed messages using PSMS signature
// (can be built using a Signer's Sign() call) and the signer's public key.
// returns:
// 		signature proof in []byte
//		error in case of errors
func (v *BLS12381G1Verifier) DeriveProof(messages [][]byte, signature, nonce []byte,
	revealedIndexes []int) ([]byte, error) {
	return v.psmsPrimitive.DeriveProof(messages, signature, nonce, v.signerPubKeyBytes, revealedIndexes)
}
