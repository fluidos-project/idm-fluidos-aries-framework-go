/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package subtlepsms

import (
	"crypto/rand"
	psms "github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPSMSG2_Sign(t *testing.T) {
	nattr := 1
	pubKey, privKey, err := generateKeyPairRandom(nattr)
	require.NoError(t, err)

	privKeyBytes, err := privKey.Marshal()
	require.NoError(t, err)

	pubKeyBytes, err := pubKey.Marshal()
	require.NoError(t, err)

	blsSigner := NewBLS12381G1Signer(privKeyBytes)
	blsVerifier := NewBLS12381G1Verifier(pubKeyBytes)

	messagesBytes := [][]byte{[]byte("epoch1230871270832321"), []byte("message2")}

	signatureBytes, err := blsSigner.Sign(messagesBytes)
	require.NoError(t, err)
	require.NotEmpty(t, signatureBytes)
	require.NoError(t, blsVerifier.Verify(messagesBytes, signatureBytes))

	// at least one message must be passed
	signatureBytes, err = blsSigner.Sign([][]byte{[]byte("epoch1230871270831212321")})
	require.Error(t, err)
	require.EqualError(t, err, "messages are not defined")
	require.Nil(t, signatureBytes)
}

func TestPSMSG2_DeriveProof(t *testing.T) {
	nattr := 3
	pubKey, privKey, err := generateKeyPairRandom(nattr)
	require.NoError(t, err)

	privKeyBytes, err := privKey.Marshal()
	require.NoError(t, err)

	pubKeyBytes, err := pubKey.Marshal()
	require.NoError(t, err)

	messagesBytes := [][]byte{
		[]byte("epoch12308712708312"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}

	blsSigner := NewBLS12381G1Signer(privKeyBytes)
	blsVerifier := NewBLS12381G1Verifier(pubKeyBytes)

	signatureBytes, err := blsSigner.Sign(messagesBytes)
	require.NoError(t, err)

	require.NoError(t, blsVerifier.Verify(messagesBytes, signatureBytes))

	nonce := []byte("nonce")
	revealedIndexes := []int{0, 2}
	proofBytes, err := blsVerifier.DeriveProof(messagesBytes, signatureBytes, nonce, revealedIndexes)
	require.NoError(t, err)
	require.NotEmpty(t, proofBytes)

	revealedMessages := make([][]byte, len(revealedIndexes))
	for i, ind := range revealedIndexes {
		revealedMessages[i] = messagesBytes[ind]
	}

	require.NoError(t, blsVerifier.VerifyProof(revealedMessages, proofBytes, nonce))
}

func generateKeyPairRandom(nattr int) (*psms.PublicKey, *psms.PrivateKey, error) {
	seed := make([]byte, 32)

	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	return psms.New().GenerateKeyPair(seed, nattr)
}
