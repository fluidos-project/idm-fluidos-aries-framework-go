/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package subtlepsms

import (
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
)

// BLS12381G1Signer is the PSMS signer for BLS12-381 curve for keys on a G1 group.
// Currently this is the only available PSMS+ signer in aries-framework-go (see `pkg/doc/psms/psms12381g2pub/psms.go`).
// Other PSMS+ signers can be added later if needed.
type BLS12381G1Signer struct {
	privateKeyBytes []byte
	psmsPrimitive   *psms12381g1pub.PSMSG1Pub
}

// NewBLS12381G2Signer creates a new instance of BLS12381G2Signer with the provided privateKey.
func NewBLS12381G1Signer(privateKey []byte) *BLS12381G1Signer {
	return &BLS12381G1Signer{
		privateKeyBytes: privateKey,
		psmsPrimitive:   psms12381g1pub.New(),
	}
}

// Sign will sign create signature of each message and aggregate it into a single signature using the signer's
// private key.
// returns:
// 		signature in []byte
//		error in case of errors
func (s *BLS12381G1Signer) Sign(messages [][]byte) ([]byte, error) {
	return s.psmsPrimitive.Sign(messages, s.privateKeyBytes)
}
