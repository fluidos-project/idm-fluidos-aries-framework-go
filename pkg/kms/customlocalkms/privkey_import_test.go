/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package customlocalkms

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	psms2 "github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto/primitive/psms"
	psms "github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"

	"github.com/golang/mock/gomock"
	tinkpb "github.com/google/tink/go/proto/tink_go_proto"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/kms"
	mockkms "github.com/hyperledger/aries-framework-go/pkg/mock/kms"
	mocksecretlock "github.com/hyperledger/aries-framework-go/pkg/mock/secretlock"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock/noop"
)

func TestImportECDSAKeyWithInvalidKey(t *testing.T) {
	k := createKMS(t)
	errPrefix := "import private EC key failed: "

	_, _, err := k.importECDSAKey(nil, kms.ECDSAP256TypeDER)
	require.EqualError(t, err, errPrefix+"private key is nil")

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	_, _, err = k.importECDSAKey(privKey, kms.AES128GCM)
	require.EqualError(t, err, errPrefix+"invalid ECDSA key type", "importECDSAKey should fail with "+
		"unsupported key type")
}

func TestImportEd25519KeyWitnInvalidKey(t *testing.T) {
	k := createKMS(t)
	errPrefix := "import private ED25519 key failed: "

	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	_, _, err = k.importEd25519Key(privKey, kms.ECDSAP256TypeDER)
	require.EqualError(t, err, errPrefix+"invalid key type")

	_, _, err = k.importEd25519Key(nil, kms.ED25519Type)
	require.EqualError(t, err, errPrefix+"private key is nil")
}

func TestImportKeySetInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	errPrefix := "import private EC key failed: "

	flagTests := []struct {
		tcName        string
		kmsProvider   kms.Provider
		ks            *tinkpb.Keyset
		expectedError string
	}{
		{
			tcName: "call importKeySet with nil keyset",
			kmsProvider: &mockProvider{
				storage: newInMemoryKMSStore(),
				secretLock: &mocksecretlock.MockSecretLock{
					ValEncrypt: "",
					ValDecrypt: "",
				},
			},
			expectedError: errPrefix + "invalid keyset data",
		},
		{
			tcName: "call importKeySet with bad secretLock Encrypt",
			kmsProvider: &mockProvider{
				storage: newInMemoryKMSStore(),
				secretLock: &mocksecretlock.MockSecretLock{
					ValEncrypt: "",
					ErrEncrypt: fmt.Errorf("bad encryption"),
					ValDecrypt: "",
				},
			},
			expectedError: errPrefix + "encrypted failed: bad encryption",
			ks: &tinkpb.Keyset{
				PrimaryKeyId: 1,
				Key:          nil,
			},
		},
		{
			tcName: "call importKeySet with bad storage getKeySet call",
			kmsProvider: &goMockProvider{
				storage: newInMemoryKMSStore(),
				secretLock: &mocksecretlock.MockSecretLock{
					ValEncrypt: "",
					ValDecrypt: "",
				},
			},
			expectedError: "import private EC key successful but failed to get key from store:",
			ks: &tinkpb.Keyset{
				PrimaryKeyId: 1,
				Key:          nil,
			},
		},
	}

	for _, tt := range flagTests {
		tc := tt
		t.Run(tc.tcName, func(t *testing.T) {
			k, err := New(testMasterKeyURI, tc.kmsProvider)
			require.NoError(t, err)

			_, _, err = k.importKeySet(tc.ks)
			if tc.tcName == "call importKeySet with bad storage getKeySet call" {
				require.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestGetKeysetInfoInvalid(t *testing.T) {
	_, err := getKeysetInfo(nil)
	require.EqualError(t, err, "keyset is nil")

	_, err = getKeysetInfo(&tinkpb.Keyset{
		PrimaryKeyId: 1,
		Key:          []*tinkpb.Keyset_Key{nil},
	})
	require.EqualError(t, err, "keyset key is nil")
}

func TestValidECPrivateKey(t *testing.T) {
	err := validECPrivateKey(&ecdsa.PrivateKey{})
	require.EqualError(t, err, "private key's public key is missing x coordinate")

	err = validECPrivateKey(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			X: new(big.Int),
		},
	})
	require.EqualError(t, err, "private key's public key is missing y coordinate")

	err = validECPrivateKey(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			X: new(big.Int),
			Y: new(big.Int),
		},
	})
	require.EqualError(t, err, "private key data is missing")

	err = validECPrivateKey(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			X: new(big.Int),
			Y: new(big.Int),
		},
		D: new(big.Int),
	},
	)
	require.NoError(t, err)
}

func TestValidPSMSPrivateKey(t *testing.T) {
	var nattr uint16
	nattr = 4
	pk, sk, _ := psms.New().GenerateKeyPair([]byte("randomseed2"), int(nattr))
	lkm := createKMS(t)
	id, key, err := lkm.ImportPrivateKey(sk, kms.BLS12381G1Type)
	require.NoError(t, err)
	require.IsType(t, &psms2.MockHandle{}, key)
	require.NotEmpty(t, id)
	skBytes, err := sk.Marshal()
	require.NoError(t, err)
	pkBytes, err := pk.Marshal()
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(skBytes, key.(*psms2.MockHandle).PrivateKey))
	require.True(t, reflect.DeepEqual(pkBytes, key.(*psms2.MockHandle).PublicKey))
	require.Equal(t, key.(*psms2.MockHandle).Nattr, nattr)
}

func createKMS(t *testing.T) *CustomLocalKMS {
	t.Helper()

	p, err := mockkms.NewProviderForKMS(mockstorage.NewMockStoreProvider(), &noop.NoLock{})
	require.NoError(t, err)

	k, err := New(testMasterKeyURI, p)
	require.NoError(t, err)

	return k
}

// mockProvider mocks a provider for KMS storage.
type goMockProvider struct {
	storage    kms.Store
	secretLock secretlock.Service
}

func (m *goMockProvider) StorageProvider() kms.Store {
	return m.storage
}

func (m *goMockProvider) SecretLock() secretlock.Service {
	return m.secretLock
}
