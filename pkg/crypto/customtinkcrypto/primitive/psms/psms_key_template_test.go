/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms

import (
	"crypto/rand"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/stretchr/testify/require"
)

func TestPSMSKeyTemplateSuccess(t *testing.T) {
	kt, err := BLS12381G1KeyTemplate(kms.WithAttrs([]string{"2"}))
	require.NoError(t, err)
	if kt.TypeUrl == PsmsSignerKeyTypeURL {
		kh, err := NewMockHandle(kt)
		require.NoError(t, err)

		pubKH, err := kh.Public()
		require.NoError(t, err)

		// now test the PSMS primitives with these keyset handles
		signer, err := NewSignerMock(kh)
		require.NoError(t, err)

		messages := [][]byte{[]byte("epoch"), []byte("msg def"), []byte("msg ghi")}

		sig, err := signer.Sign(messages)
		require.NoError(t, err)

		verifier, err := NewVerifierMock(pubKH)
		require.NoError(t, err)

		err = verifier.Verify(messages, sig)
		require.NoError(t, err)

		revealedIndexes := []int{0, 1, 2}
		nonce := make([]byte, 10)

		_, err = rand.Read(nonce)
		require.NoError(t, err)

		proof, err := verifier.DeriveProof(messages, sig, nonce, revealedIndexes)
		require.NoError(t, err)

		revealedMsgs := [][]byte{messages[0], messages[1], messages[2]}

		err = verifier.VerifyProof(revealedMsgs, proof, nonce)
		require.NoError(t, err)
	}
}
