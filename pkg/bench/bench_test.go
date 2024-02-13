package bench

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/bbs12381g2pub"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/kmsdidkey"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/internal/ldtestutil"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"

	"github.com/hyperledger/aries-framework-go/pkg/bench/internal/files"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	mockprovider "github.com/hyperledger/aries-framework-go/pkg/mock/provider"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	mockvdr "github.com/hyperledger/aries-framework-go/pkg/mock/vdr"
	w "github.com/hyperledger/aries-framework-go/pkg/wallet"

	//mockvdr "github.com/hyperledger/aries-framework-go/pkg/mock/vdr"
	"github.com/stretchr/testify/require"
)

const (
	samplePassPhrase = "pass"
	sampleUser       = "user"
	sampleUser2      = "user2"
	sampleNonce      = "nonceForZkProof"
	base58Key_5      = "133QDxmJzdAp8d91wCEzXtRiXxG9ukM6bFP9UvnUgMua6bctDh8mtqUZi1JHsQ1ajRQRhSxfACm5URWspCnZJqtraVVszCZnuHoiRZWkZVLbcNVS4ZVThQvbfUgMihmFtkf2mPZVjiGD7yKExmT8W6heZcFxAKkPYXHqgAyQuTLKBknQEgndajX3JStLyiSXRytxHaoRPcoxxSqKHLupMNtQqZ3pKHJ1ckbq6Z8pofqsrTLpxe9zwvMD7UBt73iLnAenLsALWnL2MuR5YnpKvjLUyDAmhdihfFBpnSdhStfxaoP1doX29TBx9A5jA8xbfSGVgYvfcCwaJAQNnJ6oWo8EXM84bVvtUdthFakTgSWtKTypoBKwonqhU3oevj2XMSpMiA9VsbK2iHMF8GHsW55bukAu82nWNhzhJu7EcxfJNw1oGn75iqFM1YGvBAqxSfX85oDJ9LW7kCQHh41wRvWPwshx47BXWj53K6twd4u1yv7DVvtYUrY24Uyv4x8TyGmrsYssjJoQ4jm3"
	base58Key_10     = "12Xdgi83Ua8953cwQG3j25HY5ykAT8k8zxAF6VVG4V8iTELkV9r7LEQawNnDBU9F7ne97VKrC5R53hRxrYGCY4bvWeXD9TZznNW5vn81TQRPt57EmhuPhmuGEzdahCQHVzWNnNakS6gDEUjhHvyJQKLN1hmFiAF2c8F2NeG5FmofVKWTUu5WtJk9B5xAF4K1GaybGSssYUwsCDMMCTL6mBWNKWof91UNQ5hCdeTr6irExzQZE8W3P3nMFgZTZH2stJdfD8xH357QyHLuaXC1sZuCnSYBFiGRrrPwRMirP5msQtLhr7UTPbAXUkwM2Ann5LTWwe3mBvKmhj87GpvikEpKCejuytdFpidwTC8qGwWM4tLfs7DD98tErg4Y78aDjzpyQwFMKm6rNRZ2ph5TnLm8h1iwbjC5fDVFPZ19TVRRG8gcLMimZYZu2oUZyNSTBkYhCb3xYsF6adzHPjPmaebtL9LQQk5RRxYiboZw2nkyqyikuuu7yCxGUQyANsRWerfQ9WGKqTfKNCtNCK17oL1EEvcPCbNRRH3KatupymRgnhRGxVVEBFsatwKgnoac77BwnfPWNR5j66r9orwWSUsHb7zFepXKh1zkHhJD6c1ebgp8Zwubji9z19bNDHiVLZpDf3yUdsyzsQfYtWRot3k1x1qmGRcfA4gxzw7c1615R2kug89pCPwANwRLmv5ETcM4U14L934S6qd6yyQiejSiecB2nFFCEkUsLAMacQ1sg1UHQS3pv1T5JEbitSKFAfm9F7e8ExZbZrwLZp5WT3szMKi644Z3zBJrrPGXYeZXSJA5AqYw9DQnviqUWuHTMpfBzmuZnGUARkWZY1fVjbLi"
	base58BbsKey     = "6gsgGpdx7p1nYoKJ4b5fKt1xEomWdnemg9nJFX6mqNCh"
)

func BenchmarkKeygen(b *testing.B) {
	nattr := 5
	wallet, token := createOpenWallet(b, sampleUser)
	kOpts := []string{strconv.Itoa(int(nattr))}
	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wallet.CreateKeyPair(token, kms.BLS12381G1Type, kms.WithAttrs(kOpts))
	}
}

func BenchmarkKeygenBBS(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)
	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wallet.CreateKeyPair(token, kms.BLS12381G2Type)
	}
}

func BenchmarkIssue(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58Key_10, token, b, w.PsmsBlsSignature2022)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.PsmsBlsSignature2022,
		ProofRepresentation: &proofRepr,
	}
	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
		require.NoError(b, err)
	}
}

func BenchmarkIssueBBS(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58BbsKey, token, b, w.BbsBlsSignature2020)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.BbsBlsSignature2020,
		ProofRepresentation: &proofRepr,
	}
	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
		require.NoError(b, err)
	}
}

func BenchmarkVerify(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58Key_5, token, b, w.PsmsBlsSignature2022)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.PsmsBlsSignature2022,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	wallet.Close()
	wallet2, token2 := createOpenWallet(b, sampleUser2)
	require.NoError(b, err)
	defer wallet2.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := wallet2.Verify(token2, w.WithRawCredentialToVerify(rawVcBytes))
		require.Equal(b, res, true)
		require.NoError(b, err)
	}
}

func BenchmarkVerifyBBS(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58BbsKey, token, b, w.BbsBlsSignature2020)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.BbsBlsSignature2020,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	wallet.Close()
	wallet2, token2 := createOpenWallet(b, sampleUser2)
	require.NoError(b, err)
	defer wallet2.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := wallet2.Verify(token2, w.WithRawCredentialToVerify(rawVcBytes))
		require.Equal(b, res, true)
		require.NoError(b, err)
	}
}

func BenchmarkDeriveZkDisclosure(b *testing.B) {
	print("CAREFUL, TEST if VERIFY IS USED INSIDE ZKDERIVEPROOF TO CHECK FOR PARAMETER ERRORS\n")
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58Key_5, token, b, w.PsmsBlsSignature2022)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.PsmsBlsSignature2022,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	var frameDocPsms map[string]interface{}
	require.NoError(b, json.Unmarshal(files.SampleFramePsms, &frameDocPsms))

	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dervc, err := wallet.Derive(token, w.FromRawCredential(rawVcBytes), &w.DeriveOptions{
			Nonce: sampleNonce,
			Frame: frameDocPsms,
		})
		require.NoError(b, err)
		require.NotEmpty(b, dervc)
		a, _ := dervc.MarshalJSON()
		fmt.Println(string(a))
		//verifyPSMSProof(dervc.Proofs, b)
	}
}

func BenchmarkDeriveZkDisclosureBBS(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58BbsKey, token, b, w.BbsBlsSignature2020)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.BbsBlsSignature2020,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	var frameDocPsms map[string]interface{}
	require.NoError(b, json.Unmarshal(files.SampleFramePsms, &frameDocPsms))

	defer wallet.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dervc, err := wallet.Derive(token, w.FromRawCredential(rawVcBytes), &w.DeriveOptions{
			Nonce: sampleNonce,
			Frame: frameDocPsms,
		})
		require.NoError(b, err)
		require.NotEmpty(b, dervc)
		//verifyPSMSProof(dervc.Proofs, b)
	}
}

func BenchmarkVerifyZkDisclosure(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58Key_5, token, b, w.PsmsBlsSignature2022)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.PsmsBlsSignature2022,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	var frameDocPsms map[string]interface{}
	require.NoError(b, json.Unmarshal(files.SampleFramePsms, &frameDocPsms))
	dervc, err := wallet.Derive(token, w.FromRawCredential(rawVcBytes), &w.DeriveOptions{
		Nonce: sampleNonce,
		Frame: frameDocPsms,
	})
	require.NoError(b, err)
	require.NotEmpty(b, dervc)
	rawDerVcBytes, err := dervc.MarshalJSON()
	require.NoError(b, err)
	wallet.Close()
	wallet2, token2 := createOpenWallet(b, sampleUser2)
	require.NoError(b, err)
	defer wallet2.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := wallet2.Verify(token2, w.WithRawCredentialToVerify(rawDerVcBytes))
		require.Equal(b, res, true)
		require.NoError(b, err)
	}
}

func BenchmarkVerifykDisclosureBBS(b *testing.B) {
	wallet, token := createOpenWallet(b, sampleUser)

	did := addKeyToWallet(base58BbsKey, token, b, w.BbsBlsSignature2020)

	proofRepr := verifiable.SignatureProofValue
	proofOpts := &w.ProofOptions{
		Controller:          did,
		ProofType:           w.BbsBlsSignature2020,
		ProofRepresentation: &proofRepr,
	}
	vc, err := wallet.Issue(token, files.ExampleRawVC, proofOpts)
	require.NoError(b, err)
	rawVcBytes, err := vc.MarshalJSON()
	require.NoError(b, err)
	var frameDocPsms map[string]interface{}
	require.NoError(b, json.Unmarshal(files.SampleFramePsms, &frameDocPsms))

	dervc, err := wallet.Derive(token, w.FromRawCredential(rawVcBytes), &w.DeriveOptions{
		Nonce: sampleNonce,
		Frame: frameDocPsms,
	})
	require.NoError(b, err)
	require.NotEmpty(b, dervc)
	rawDerVcBytes, err := dervc.MarshalJSON()
	require.NoError(b, err)
	wallet.Close()
	wallet2, token2 := createOpenWallet(b, sampleUser2)
	require.NoError(b, err)
	defer wallet2.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := wallet2.Verify(token2, w.WithRawCredentialToVerify(rawDerVcBytes))
		require.Equal(b, res, true)
		require.NoError(b, err)
	}
}

func verifyPSMSProof(proofs []verifiable.Proof, b *testing.B) {
	require.Len(b, proofs, 1)
	require.NotEmpty(b, proofs[0])
	require.Equal(b, proofs[0]["type"], "PsmsBlsSignatureProof2022")
	require.NotEmpty(b, proofs[0]["nonce"])
	require.EqualValues(b, proofs[0]["nonce"], base64.StdEncoding.EncodeToString([]byte(sampleNonce)))
	require.NotEmpty(b, proofs[0]["proofValue"])
}

func addKeyToWallet(b58key, token string, b *testing.B, scheme string) string {
	kmgr, err := w.TestHelperKeymanager(token)
	require.NoError(b, err)
	require.NotEmpty(b, kmgr)

	switch scheme {
	case w.PsmsBlsSignature2022:
		privKeyPSMS, err := psms12381g1pub.UnmarshalPrivateKey(base58.Decode(base58Key_5))
		require.NoError(b, err)
		pkb, _ := privKeyPSMS.PublicKey().Marshal()
		did, _ := kmsdidkey.BuildDIDKeyByKeyType(pkb, kms.BLS12381G1Type)
		parts := strings.SplitN(did, ":", 3)
		//fmt.Println(did)
		// nolint: errcheck, gosec
		_, _, err = (*kmgr).ImportPrivateKey(privKeyPSMS, kms.BLS12381G1Type, kms.WithKeyID(parts[2]))
		require.NoError(b, err)
		return did
	case w.BbsBlsSignature2020:
		privKeyBBS, err := bbs12381g2pub.UnmarshalPrivateKey(base58.Decode(b58key))
		require.NoError(b, err)
		pkb, _ := privKeyBBS.PublicKey().Marshal()
		did, _ := kmsdidkey.BuildDIDKeyByKeyType(pkb, kms.BLS12381G2Type)
		parts := strings.SplitN(did, ":", 3)
		// nolint: errcheck, gosec
		_, _, err = (*kmgr).ImportPrivateKey(privKeyBBS, kms.BLS12381G2Type, kms.WithKeyID(parts[2]))
		require.NoError(b, err)
		return did

	default:
		panic("Should use psms or bbs scheme")
	}

}

func createOpenWallet(b *testing.B, user string) (*w.Wallet, string) {
	mockctx := newMockProvider(b)
	customVDR := &mockvdr.MockVDRegistry{
		ResolveFunc: func(didID string, opts ...vdrapi.DIDMethodOption) (*did.DocResolution, error) {
			if strings.HasPrefix(didID, "did:key:") {
				k := key.New()
				d, e := k.Read(didID)
				if e != nil {
					return nil, e
				}
				return d, nil
			}

			return nil, fmt.Errorf("did not found")
		}}
	mockctx.VDRegistryValue = customVDR
	err := w.CreateProfile(user, mockctx, w.WithPassphrase(samplePassPhrase))
	require.NoError(b, err)

	customCrypto, err := customtinkcrypto.New()
	require.NoError(b, err)

	mockctx.CryptoValue = customCrypto

	wallet, err := w.New(user, mockctx)
	require.NoError(b, err)
	require.NotEmpty(b, wallet)

	token, err := wallet.Open(w.WithUnlockByPassphrase(samplePassPhrase), w.WithUnlockExpiry(50000*time.Millisecond))
	require.NoError(b, err)
	require.NotEmpty(b, token)

	return wallet, token
}

func newMockProvider(t *testing.B) *mockprovider.Provider {
	t.Helper()

	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	return &mockprovider.Provider{
		StorageProviderValue:              mockstorage.NewMockStoreProvider(),
		ProtocolStateStorageProviderValue: mockstorage.NewMockStoreProvider(),
		DocumentLoaderValue:               loader,
	}
}
