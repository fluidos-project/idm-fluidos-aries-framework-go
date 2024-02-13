/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignatureproof2022_test

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignatureproof2022"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/internal/ldtestutil"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	//go:embed testdata/customcase_vc.jsonld
	customCaseVc string
	//go:embed testdata/customcase_reveal_doc.jsonld
	customCaseRevealDoc string

	//go:embed testdata/doc_with_many_proofs.jsonld
	docWithManyProofsJSON string //nolint:unused // re-enable test that uses this var (#2562)
)

// nolint
func TestSuite_SelectiveDisclosure(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pkBase58 := "17Dy7o2xoHXgYErocwsFLfSeNHYidikmFan1dwpi3T5KHiRkryoZAUF7HvAdG2eHah5QTBNnUUtzfrWJYmG8K7YY39zqdW1vAYE8ZzK1YYK4PXEu9Kwby3xBqVPAjQNsAPbB3xUj4iFfCV8dqqdmBKXzGEcT4gZpnA27ggc6zhwY4GhfZokbBvX9BbANzEEnJqcUx6ZnDR3GahUu6H1A3RXD5CSvkF9KAPpMAk126V8GA85gZGvJnRgHqriqWuL79euMP1yoyCeGuShbfjxxnry4YGtwTvG248Q25DkoEhcNFWTKchTGKrWPQU8mScQ3so6N8p7DoMVTYSFFr9X7RWK6FRynfvAZWBpKrWfrVLA4qXF1ntrVgjUSnTQTX4HgZ54P2isLMayhUWHYuq6dqt14ReaFUjEbNCrajGeuGVANh78uGqvmNbb5KmxuPZM2LbpjLD7U3CPkvztsKWAZGyRoMveg5JWQkg7PyGbeT8JcKHq5R4ueGeXf4g9XzktbApLq4snZAkQ1N5yzeUMgH4Ptq767pfBmmmEvZQhy9kCc8yYHLMkDc2cLDT9qeFe8UsoiCo4bbyiajJEcCaHZ5gMot42d2LEmSXxdmxnfUDmWRu3KK6St7C1NooRmfK6kETzBQxLe3a3o9RMMmpMkw1mXM3RKFvFQZrrKVpWQXgRm7JaHCLkD9Zfadq5f3o8rqa7s2xfDyPHXggqh3ebwejqQo5SH48X6vbULzqdJahCTuBiTercyDkdxh7GrkAv8D6PmJ5BzVDrPbfRXTd4w7Yji6vo6Mb"
	pubKeyBytes := base58.Decode(pkBase58)

	nonce, err := base64.StdEncoding.DecodeString("G/hn9Ca9bIWZpJGlhnr/41r8RB0OO0TLChZASr3QJVztdri/JzS8Zf/xWJT5jW78zlM=")
	require.NoError(t, err)

	docMap := toMap(t, customCaseVc)
	revealDocMap := toMap(t, customCaseRevealDoc)

	s := psmsblssignatureproof2022.New()

	const proofField = "proof"

	pubKeyResolver := &testKeyResolver{
		publicKey: &verifier.PublicKey{
			Type:  "Bls12381G1Key2022",
			Value: pubKeyBytes,
		},
	}

	t.Run("single PSMS+ signature", func(t *testing.T) {
		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMap, revealDocMap, nonce,
			pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.NoError(t, err)
		require.NotEmpty(t, docWithSelectiveDisclosure)
		require.Contains(t, docWithSelectiveDisclosure, proofField)
		proofs, ok := docWithSelectiveDisclosure[proofField].([]map[string]interface{})
		require.True(t, ok)
		//res, err := json.Marshal(docWithSelectiveDisclosure)
		//require.NoError(t, err)
		//fmt.Println(string(res))
		require.Len(t, proofs, 1)
		require.Equal(t, "PsmsBlsSignatureProof2022", proofs[0]["type"])
		require.NotEmpty(t, proofs[0]["proofValue"])
	})
	/*
		t.Run("several proofs including PSMS+ signature", func(t *testing.T) {
			// TODO re-enable (#2562).
			t.Skip()
			docWithSeveralProofsMap := toMap(t, docWithManyProofsJSON)

			pubKeyBytes2 := base58.Decode("tPTWWeUm8yT3aR9HtMvo2pLLvAdyV9Z4nJYZ2ZsyoLVpTupVb7NaRJ3tZePF6YsCN1nw7McqJ38tvpmQxKQxrTbyzjiewUDaj5jbD8gVfpfXJL2SfPBw4TGjYPA6zg6Jrxn")

			compositeResolver := &testKeyResolver{
				variants: map[string]*verifier.PublicKey{
					"did:example:489398593#test": {
						Type:  "Bls12381G2Key2020",
						Value: pubKeyBytes},
					"did:example:123456#key2": {
						Type:  "Bls12381G2Key2020",
						Value: pubKeyBytes2},
				},
			}

			docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docWithSeveralProofsMap, revealDocMap, nonce,
				compositeResolver, ldtestutil.WithDocumentLoader(t))
			require.NoError(t, err)
			require.NotEmpty(t, docWithSelectiveDisclosure)
			require.Contains(t, docWithSelectiveDisclosure, proofField)

			proofs, ok := docWithSelectiveDisclosure[proofField].([]map[string]interface{})
			require.True(t, ok)

			require.Len(t, proofs, 2)
			require.Equal(t, "PsmsBlsSignatureProof2022", proofs[0]["type"])
			require.NotEmpty(t, proofs[0]["proofValue"])
			require.Equal(t, "PsmsBlsSignatureProof2022", proofs[1]["type"])
			require.NotEmpty(t, proofs[1]["proofValue"])
		})*/

	t.Run("malformed input", func(t *testing.T) {
		docMap := make(map[string]interface{})
		docMap["@context"] = "http://localhost/nocontext"
		docMap["bad"] = "example"
		docMap["proof"] = "example"

		_, err := s.SelectiveDisclosure(docMap, revealDocMap, nonce, pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
	})

	t.Run("no proof", func(t *testing.T) {
		docMapWithoutProof := make(map[string]interface{}, len(docMap)-1)

		for k, v := range docMap {
			if k != proofField {
				docMapWithoutProof[k] = v
			}
		}

		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMapWithoutProof, revealDocMap, nonce,
			pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
		require.Contains(t, err.Error(), "document does not have a proof")
		require.Empty(t, docWithSelectiveDisclosure)
	})

	t.Run("invalid proof", func(t *testing.T) {
		docMapWithInvalidProof := make(map[string]interface{}, len(docMap)-1)

		for k, v := range docMap {
			if k != proofField {
				docMapWithInvalidProof[k] = v
			} else {
				docMapWithInvalidProof[k] = "invalid proof"
			}
		}

		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMapWithInvalidProof, revealDocMap, nonce,
			pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
		require.EqualError(t, err, "get BLS proofs: read document proofs: proof is not map or array of maps")
		require.Empty(t, docWithSelectiveDisclosure)
	})

	t.Run("invalid proof value", func(t *testing.T) {
		docMapWithInvalidProofValue := make(map[string]interface{}, len(docMap))

		for k, v := range docMap {
			if k == proofField {
				proofMap := make(map[string]interface{})

				for k1, v1 := range v.(map[string]interface{}) {
					if k1 == "proofValue" {
						proofMap[k1] = "invalid"
					} else {
						proofMap[k1] = v1
					}
				}

				docMapWithInvalidProofValue[proofField] = proofMap
			} else {
				docMapWithInvalidProofValue[k] = v
			}
		}

		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMapWithInvalidProofValue, revealDocMap, nonce,
			pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
		require.EqualError(t, err, "generate signature proof: derive PSMS proof: unmarshal signature: marshalled signature of wrong size 5") //nolint:lll
		require.Empty(t, docWithSelectiveDisclosure)
	})

	t.Run("invalid input PSMS+ proof value", func(t *testing.T) {
		docMapWithInvalidProofType := make(map[string]interface{}, len(docMap)-1)

		for k, v := range docMap {
			if k == proofField {
				proofMap := make(map[string]interface{})

				for k1, v1 := range v.(map[string]interface{}) {
					if k1 == "type" {
						proofMap[k1] = "invalid"
					} else {
						proofMap[k1] = v1
					}
				}

				docMapWithInvalidProofType[proofField] = proofMap
			} else {
				docMapWithInvalidProofType[k] = v
			}
		}

		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMapWithInvalidProofType, revealDocMap, nonce,
			pubKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
		require.EqualError(t, err, "no PsmsBlsSignature2022 proof present")
		require.Empty(t, docWithSelectiveDisclosure)
	})

	t.Run("failed to resolve public key", func(t *testing.T) {
		failingPublicKeyResolver := &testKeyResolver{
			err: errors.New("public key not found"),
		}

		docWithSelectiveDisclosure, err := s.SelectiveDisclosure(docMap, revealDocMap, nonce,
			failingPublicKeyResolver, ldtestutil.WithDocumentLoader(t))
		require.Error(t, err)
		require.EqualError(t, err, "generate signature proof: get public key and signature: resolve public key of PSMS+ signature: public key not found") //nolint:lll
		require.Empty(t, docWithSelectiveDisclosure)
	})
	/*
		t.Run("Case 18 derives into Case 19", func(t *testing.T) {
			case18DocMap := toMap(t, case18VC)
			case18RevealDocMap := toMap(t, case18RevealDoc)

			case19Nonce, err := base64.StdEncoding.DecodeString("lEixQKDQvRecCifKl789TQj+Ii6YWDLSwn3AxR0VpPJ1QV5htod/0VCchVf1zVM0y2E=")
			require.NoError(t, err)

			docWithSelectiveDisclosure, err := s.SelectiveDisclosure(case18DocMap, case18RevealDocMap, case19Nonce,
				pubKeyResolver, ldtestutil.WithDocumentLoader(t))
			require.NoError(t, err)
			require.NotEmpty(t, docWithSelectiveDisclosure)
			require.Contains(t, docWithSelectiveDisclosure, proofField)

			proofs, ok := docWithSelectiveDisclosure[proofField].([]map[string]interface{})
			require.True(t, ok)

			require.Len(t, proofs, 1)
			require.Equal(t, "PsmsBlsSignatureProof2022", proofs[0]["type"])
			require.NotEmpty(t, proofs[0]["proofValue"])

			case18DerivationBytes, err := json.Marshal(docWithSelectiveDisclosure)

			pubKeyFetcher := verifiable.SingleKey(pubKeyBytes, "Bls12381G2Key2020")

			loader, err := ldtestutil.DocumentLoader()
			require.NoError(t, err)

			_, err = verifiable.ParseCredential(case18DerivationBytes, verifiable.WithPublicKeyFetcher(pubKeyFetcher),
				verifiable.WithJSONLDDocumentLoader(loader))
			require.NoError(t, err)
		})
	*/
}

func toMap(t *testing.T, doc string) map[string]interface{} {
	var docMap map[string]interface{}
	err := json.Unmarshal([]byte(doc), &docMap)
	require.NoError(t, err)

	return docMap
}
