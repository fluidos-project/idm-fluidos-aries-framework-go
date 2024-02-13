/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package customtinkcrypto_test

import (
	"crypto/rand"
	"strconv"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/kms/customlocalkms"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/jwkkid"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock/noop"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

type kmsProvider struct {
	store             kms.Store
	secretLockService secretlock.Service
}

func (k *kmsProvider) StorageProvider() kms.Store {
	return k.store
}

func (k *kmsProvider) SecretLock() secretlock.Service {
	return k.secretLockService
}

func TestSignVerifyKeyTypes(t *testing.T) {
	testCases := []struct {
		name    string
		keyType kms.KeyType
	}{
		{
			"P-256",
			kms.ECDSAP256TypeIEEEP1363,
		},
		{
			"P-384",
			kms.ECDSAP384TypeIEEEP1363,
		},
		{
			"P-521",
			kms.ECDSAP521TypeIEEEP1363,
		},
	}

	data := []byte("abcdefg 1234567 1234567 1234567 1234567 1234567 AaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAaAa")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kmsStore, err := kms.NewAriesProviderWrapper(mockstorage.NewMockStoreProvider())
			require.NoError(t, err)

			kmsStorage, err := customlocalkms.New("local-lock://test/master/key/", &kmsProvider{
				store:             kmsStore,
				secretLockService: &noop.NoLock{},
			})
			require.NoError(t, err)

			cr, err := customtinkcrypto.New()
			require.NoError(t, err)

			kid, pkb, err := kmsStorage.CreateAndExportPubKeyBytes(tc.keyType)
			require.NoError(t, err)

			kh, err := kmsStorage.Get(kid)
			require.NoError(t, err)

			pkJWK, err := jwkkid.BuildJWK(pkb, tc.keyType)
			require.NoError(t, err)

			jkBytes, err := pkJWK.PublicKeyBytes()
			require.NoError(t, err)
			require.Equal(t, pkb, jkBytes)

			kh2, err := kmsStorage.PubKeyBytesToHandle(jkBytes, tc.keyType)
			require.NoError(t, err)

			sig, err := cr.Sign(data, kh)
			require.NoError(t, err)

			err = cr.Verify(sig, data, kh2)
			require.NoError(t, err)
		})
	}
}

type mockProvider struct {
	storeProvider storage.Provider
	secretLock    secretlock.Service
}

func (m *mockProvider) StorageProvider() storage.Provider {
	return m.storeProvider
}

func (m *mockProvider) SecretLock() secretlock.Service {
	return m.secretLock
}

func TestMultiSignVerifyKeyTypes(t *testing.T) {
	const testMessage = "test message"
	msg := [][]byte{
		[]byte("epoch123123123123"), []byte(testMessage + "1"), []byte(testMessage + "2"),
		[]byte(testMessage + "3"), []byte(testMessage + "4"), []byte(testMessage + "5"),
	}
	testCases := []struct {
		name    string
		keyType kms.KeyType
		keyOpts kms.KeyOpts
	}{
		{
			"PSMS",
			kms.BLS12381G1Type,
			kms.WithAttrs([]string{strconv.Itoa((len(msg) - 1))}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kmsStore, err := kms.NewAriesProviderWrapper(mockstorage.NewMockStoreProvider())
			require.NoError(t, err)

			kmsStorage, err := customlocalkms.New("local-lock://test/master/key/", &kmsProvider{
				store:             kmsStore,
				secretLockService: &noop.NoLock{},
			})
			require.NoError(t, err)

			cr, err := customtinkcrypto.New()
			require.NoError(t, err)

			kid, pkb, err := kmsStorage.CreateAndExportPubKeyBytes(tc.keyType, tc.keyOpts)
			require.NoError(t, err)

			kh, err := kmsStorage.Get(kid)
			require.NoError(t, err)

			/*
				pkJWK, err := jwkkid.BuildJWK(pkb, tc.keyType)
				require.NoError(t, err)

				jkBytes, err := pkJWK.PublicKeyBytes()
				require.NoError(t, err)
				require.Equal(t, pkb, jkBytes) */

			kh2, err := kmsStorage.PubKeyBytesToHandle(pkb, tc.keyType, tc.keyOpts)
			require.NoError(t, err)

			sig, err := cr.SignMulti(msg, kh)
			require.NoError(t, err)

			err = cr.VerifyMulti(msg, sig, kh2)
			require.NoError(t, err)

			revealedIndexes := []int{0, 3}
			nonce := make([]byte, 32)

			_, err = rand.Read(nonce)
			require.NoError(t, err)

			_, err = cr.DeriveProof([][]byte{}, sig, nonce, revealedIndexes, kh2)
			require.EqualError(t, err, "PSMS derive proof msg: invalid size: 2 revealed indexes is larger than 0 messages")

			proof, err := cr.DeriveProof(msg, sig, nonce, revealedIndexes, kh2)
			require.NoError(t, err)

			err = cr.VerifyProof([][]byte{msg[0], msg[1]}, proof, nonce, kh2)
			require.EqualError(t, err, "verify proof msg: invalid PSMS BLS12-381 zk proof")

			err = cr.VerifyProof([][]byte{msg[0], msg[3]}, proof, nonce[1:], kh2)
			require.EqualError(t, err, "verify proof msg: invalid PSMS BLS12-381 zk proof")

			err = cr.VerifyProof([][]byte{msg[0], msg[3]}, proof, nonce, kh2)
			require.NoError(t, err)

		})
	}
}
