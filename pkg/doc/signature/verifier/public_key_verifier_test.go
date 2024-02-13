/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package verifier

import (
	"crypto"
	"crypto/ed25519"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	gojose "github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/crypto/tinkcrypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/signature"
	kmsapi "github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/kms/localkms"
	mockkms "github.com/hyperledger/aries-framework-go/pkg/mock/kms"
	"github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock/noop"
)

func TestNewPublicKeyVerifier(t *testing.T) {
	var (
		publicKey = &PublicKey{
			Type: "TestType",
			JWK: &jwk.JWK{
				JSONWebKey: gojose.JSONWebKey{
					Algorithm: "alg",
				},
				Kty: "kty",
				Crv: "crv",
			},
		}

		msg    = []byte("message to sign")
		msgSig = []byte("signature")

		signatureVerifier = &testSignatureVerifier{
			baseSignatureVerifier: baseSignatureVerifier{
				keyType:   "kty",
				curve:     "crv",
				algorithm: "alg",
			},
			verifyResult: nil,
		}
	)

	verifier := NewPublicKeyVerifier(signatureVerifier, WithExactPublicKeyType("TestType"))
	require.NotNil(t, verifier)

	err := verifier.Verify(publicKey, msg, msgSig)
	require.NoError(t, err)

	t.Run("check public key type", func(t *testing.T) {
		publicKey.Type = "invalid TestType"

		err = verifier.Verify(publicKey, msg, msgSig)
		require.Error(t, err)
		require.EqualError(t, err, "a type of public key is not 'TestType'")

		publicKey.Type = "TestType"
	})

	t.Run("match JWK key type", func(t *testing.T) {
		publicKey.JWK.Kty = "invalid kty"

		err = verifier.Verify(publicKey, msg, msgSig)
		require.Error(t, err)
		require.EqualError(t, err, "verifier does not match JSON Web Key")

		publicKey.JWK.Kty = "kty"
	})

	t.Run("match JWK curve", func(t *testing.T) {
		publicKey.JWK.Crv = "invalid crv"

		err = verifier.Verify(publicKey, msg, msgSig)
		require.Error(t, err)
		require.EqualError(t, err, "verifier does not match JSON Web Key")

		publicKey.JWK.Crv = "crv"
	})

	t.Run("match JWK algorithm", func(t *testing.T) {
		publicKey.JWK.Algorithm = "invalid alg"

		err = verifier.Verify(publicKey, msg, msgSig)
		require.Error(t, err)
		require.EqualError(t, err, "verifier does not match JSON Web Key")

		publicKey.JWK.Algorithm = "alg"
	})

	signatureVerifier.verifyResult = errors.New("invalid signature")
	err = verifier.Verify(publicKey, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "invalid signature")
}

func TestNewCompositePublicKeyVerifier(t *testing.T) {
	var (
		publicKey = &PublicKey{
			Type: "TestType",
			JWK: &jwk.JWK{
				JSONWebKey: gojose.JSONWebKey{
					Algorithm: "alg",
				},
				Kty: "kty",
				Crv: "crv",
			},
		}

		msg    = []byte("message to sign")
		msgSig = []byte("signature")

		signatureVerifier = &testSignatureVerifier{
			baseSignatureVerifier: baseSignatureVerifier{
				keyType:   "kty",
				curve:     "crv",
				algorithm: "alg",
			},
			verifyResult: nil,
		}
	)

	verifier := NewCompositePublicKeyVerifier([]SignatureVerifier{signatureVerifier},
		WithExactPublicKeyType("TestType"))
	require.NotNil(t, verifier)

	err := verifier.Verify(publicKey, msg, msgSig)
	require.NoError(t, err)

	publicKey.JWK.Kty = "invalid kty"
	err = verifier.Verify(publicKey, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "no matching verifier found")

	publicKey.JWK.Kty = "kty"

	signatureVerifier.verifyResult = errors.New("invalid signature")
	err = verifier.Verify(publicKey, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "invalid signature")
}

func TestNewEd25519SignatureVerifier(t *testing.T) {
	v := NewEd25519SignatureVerifier()
	require.NotNil(t, v)

	signer, err := newCryptoSigner(kmsapi.ED25519Type)
	require.NoError(t, err)

	msg := []byte("test message")
	msgSig, err := signer.Sign(msg)
	require.NoError(t, err)

	pubKey := &PublicKey{
		Type:  kmsapi.ED25519,
		Value: signer.PublicKeyBytes(),
	}

	err = v.Verify(pubKey, msg, msgSig)
	require.NoError(t, err)

	// invalid public key type
	err = v.Verify(&PublicKey{
		Type:  kmsapi.ED25519,
		Value: []byte("invalid-key"),
	}, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "ed25519: invalid key")

	// invalid signature
	err = v.Verify(pubKey, msg, []byte("invalid signature"))
	require.Error(t, err)
	require.EqualError(t, err, "ed25519: invalid signature")
}

func TestNewRSAPS256SignatureVerifier(t *testing.T) {
	v := NewRSAPS256SignatureVerifier()
	require.NotNil(t, v)

	signer, err := newCryptoSigner(kmsapi.RSAPS256Type)
	require.NoError(t, err)

	msg := []byte("test message")

	msgSig, err := signer.Sign(msg)
	require.NoError(t, err)

	pubKey := &PublicKey{
		Type: "JwsVerificationKey2020",
		JWK: &jwk.JWK{
			JSONWebKey: gojose.JSONWebKey{
				Algorithm: "PS256",
			},
			Kty: "RSA",
		},
		Value: signer.PublicKeyBytes(),
	}

	err = v.Verify(pubKey, msg, msgSig)
	require.NoError(t, err)

	// invalid signature
	err = v.Verify(pubKey, msg, []byte("invalid signature"))
	require.Error(t, err)
	require.EqualError(t, err, "rsa: invalid signature")

	// invalid public key
	pubKey.Value = []byte("invalid-key")
	err = v.Verify(pubKey, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "rsa: invalid public key")
}

func TestNewRSARS256SignatureVerifier(t *testing.T) {
	v := NewRSARS256SignatureVerifier()
	require.NotNil(t, v)

	signer, err := newCryptoSigner(kmsapi.RSARS256Type)
	require.NoError(t, err)

	msg := []byte("test message")

	msgSig, err := signer.Sign(msg)
	require.NoError(t, err)

	pubKey := &PublicKey{
		Type: "JsonWebKey2020",
		JWK: &jwk.JWK{
			JSONWebKey: gojose.JSONWebKey{
				Algorithm: "RS256",
			},
			Kty: "RSA",
		},
		Value: signer.PublicKeyBytes(),
	}

	err = v.Verify(pubKey, msg, msgSig)
	require.NoError(t, err)

	// invalid signature
	err = v.Verify(pubKey, msg, []byte("invalid signature"))
	require.Error(t, err)
	require.EqualError(t, err, "crypto/rsa: verification error")

	// invalid public key
	pubKey.Value = []byte("invalid-key")
	err = v.Verify(pubKey, msg, msgSig)
	require.Error(t, err)
	require.EqualError(t, err, "not *rsa.VerificationMethod public key")
}

func TestNewECDSAES256SignatureVerifier(t *testing.T) {
	msg := []byte("test message")

	t.Run("happy path", func(t *testing.T) {
		tests := []struct {
			sVerifier SignatureVerifier
			curve     elliptic.Curve
			curveName string
			algorithm string
			hash      crypto.Hash
		}{
			{
				sVerifier: NewECDSAES256SignatureVerifier(),
				curve:     elliptic.P256(),
				curveName: "P-256",
				algorithm: "ES256",
				hash:      crypto.SHA256,
			},
			{
				sVerifier: NewECDSAES384SignatureVerifier(),
				curve:     elliptic.P384(),
				curveName: "P-384",
				algorithm: "ES384",
				hash:      crypto.SHA384,
			},
			{
				sVerifier: NewECDSAES521SignatureVerifier(),
				curve:     elliptic.P521(),
				curveName: "P-521",
				algorithm: "ES521",
				hash:      crypto.SHA512,
			},
			{
				sVerifier: NewECDSASecp256k1SignatureVerifier(),
				curve:     btcec.S256(),
				curveName: "secp256k1",
				algorithm: "ES256K",
				hash:      crypto.SHA256,
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.curveName, func(t *testing.T) {
				keyType, err := signature.MapECCurveToKeyType(tc.curve)
				require.NoError(t, err)

				signer, err := newCryptoSigner(keyType)
				require.NoError(t, err)

				pubKey := &PublicKey{
					Type:  "JwsVerificationKey2020",
					Value: signer.PublicKeyBytes(),
					JWK: &jwk.JWK{
						JSONWebKey: gojose.JSONWebKey{
							Algorithm: tc.algorithm,
							Key:       signer.PublicKey(),
						},
						Crv: tc.curveName,
						Kty: "EC",
					},
				}

				msgSig, err := signer.Sign(msg)
				require.NoError(t, err)

				err = tc.sVerifier.Verify(pubKey, msg, msgSig)
				require.NoError(t, err)
			})
		}
	})

	v := NewECDSAES256SignatureVerifier()
	require.NotNil(t, v)

	signer, err := newCryptoSigner(kmsapi.ECDSAP256TypeIEEEP1363)
	require.NoError(t, err)
	msgSig, err := signer.Sign(msg)
	require.NoError(t, err)

	t.Run("verify with public key bytes", func(t *testing.T) {
		verifyError := v.Verify(&PublicKey{
			Type:  "JwsVerificationKey2020",
			Value: signer.PublicKeyBytes(),
		}, msg, msgSig)

		require.NoError(t, verifyError)
	})

	t.Run("invalid public key", func(t *testing.T) {
		verifyError := v.Verify(&PublicKey{
			Type:  "JwsVerificationKey2020",
			Value: []byte("invalid public key"),
		}, msg, msgSig)

		require.Error(t, verifyError)
		require.EqualError(t, verifyError, "ecdsa: create JWK from public key bytes: invalid public key")
	})

	t.Run("invalid public key type", func(t *testing.T) {
		ed25519Key := &ed25519.PublicKey{}

		verifyError := v.Verify(&PublicKey{
			Type:  "JwsVerificationKey2020",
			Value: signer.PublicKeyBytes(),
			JWK: &jwk.JWK{
				JSONWebKey: gojose.JSONWebKey{
					Algorithm: "ES256",
					Key:       ed25519Key,
				},
				Crv: "P-256",
				Kty: "EC",
			},
		}, msg, msgSig)
		require.Error(t, verifyError)
		require.EqualError(t, verifyError, "ecdsa: invalid public key type")
	})

	t.Run("invalid signature", func(t *testing.T) {
		pubKey := &PublicKey{
			Type:  "JwsVerificationKey2020",
			Value: signer.PublicKeyBytes(),

			JWK: &jwk.JWK{
				JSONWebKey: gojose.JSONWebKey{
					Algorithm: "ES256",
					Key:       signer.PublicKey(),
				},
				Crv: "P-256",
				Kty: "EC",
			},
		}

		verifyError := v.Verify(pubKey, msg, []byte("signature of invalid size"))
		require.Error(t, verifyError)
		require.EqualError(t, verifyError, "ecdsa: invalid signature size")

		emptySig := make([]byte, 64)
		verifyError = v.Verify(pubKey, msg, emptySig)
		require.Error(t, verifyError)
		require.EqualError(t, verifyError, "ecdsa: invalid signature")
	})
}

func TestTransformFromBlankNodes(t *testing.T) {
	const (
		a  = "<urn:bnid:_:c14n0>"
		ae = "_:c14n0"
		b  = "<urn:bnid:_:c14n0> "
		be = "_:c14n0 "
		c  = "abcd <urn:bnid:_:c14n0> "
		ce = "abcd _:c14n0 "
		d  = "abcd <urn:bnid:_:c14n0> efgh"
		de = "abcd _:c14n0 efgh"
		e  = "abcd <urn:bnid:_:c14n23> efgh"
		ee = "abcd _:c14n23 efgh"
		f  = "abcd <urn:bnid:_:c14n> efgh"
		fe = "abcd _:c14n efgh"
		g  = ""
		ge = ""
	)

	at := transformFromBlankNode(a)
	require.Equal(t, ae, at)

	bt := transformFromBlankNode(b)
	require.Equal(t, be, bt)

	ct := transformFromBlankNode(c)
	require.Equal(t, ce, ct)

	dt := transformFromBlankNode(d)
	require.Equal(t, de, dt)

	et := transformFromBlankNode(e)
	require.Equal(t, ee, et)

	ft := transformFromBlankNode(f)
	require.Equal(t, fe, ft)

	gt := transformFromBlankNode(g)
	require.Equal(t, ge, gt)
}

//
//nolint:lll,goconst
func TestNewBBSG2SignatureVerifier(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pubKeyBase58 := "nEP2DEdbRaQ2r5Azeatui9MG6cj7JUHa8GD7khub4egHJREEuvj4Y8YG8w51LnhPEXxVV1ka93HpSLkVzeQuuPE1mH9oCMrqoHXAKGBsuDT1yJvj9cKgxxLCXiRRirCycki"
	pubKeyBytes := base58.Decode(pubKeyBase58)

	sigBase64 := `qPrB+1BLsVSeOo1ci8dMF+iR6aa5Q6iwV/VzXo2dw94ctgnQGxaUgwb8Hd68IiYTVabQXR+ZPuwJA//GOv1OwXRHkHqXg9xPsl8HcaXaoWERanxYClgHCfy4j76Vudr14U5AhT3v8k8f0oZD+zBIUQ==`
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	// Case 16 (https://github.com/w3c-ccg/vc-http-api/pull/128)
	msg := `
_:c14n0 <http://purl.org/dc/terms/created> "2021-02-23T19:31:12Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
_:c14n0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/security#BbsBlsSignature2020> .
_:c14n0 <https://w3id.org/security#proofPurpose> <https://w3id.org/security#assertionMethod> .
_:c14n0 <https://w3id.org/security#verificationMethod> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2#zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
<did:example:b34ca6cd37bbf23> <http://schema.org/birthDate> "1958-07-17"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<did:example:b34ca6cd37bbf23> <http://schema.org/familyName> "SMITH" .
<did:example:b34ca6cd37bbf23> <http://schema.org/gender> "Male" .
<did:example:b34ca6cd37bbf23> <http://schema.org/givenName> "JOHN" .
<did:example:b34ca6cd37bbf23> <http://schema.org/image> <data:image/png;base64,iVBORw0KGgo...kJggg==> .
<did:example:b34ca6cd37bbf23> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Person> .
<did:example:b34ca6cd37bbf23> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/citizenship#PermanentResident> .
<did:example:b34ca6cd37bbf23> <https://w3id.org/citizenship#birthCountry> "Bahamas" .
<did:example:b34ca6cd37bbf23> <https://w3id.org/citizenship#commuterClassification> "C1" .
<did:example:b34ca6cd37bbf23> <https://w3id.org/citizenship#lprCategory> "C09" .
<did:example:b34ca6cd37bbf23> <https://w3id.org/citizenship#lprNumber> "999-999-999" .
<did:example:b34ca6cd37bbf23> <https://w3id.org/citizenship#residentSince> "2015-01-01"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://schema.org/description> "Government of Example Permanent Resident Card." .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://schema.org/name> "Permanent Resident Card" .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/citizenship#PermanentResidentCard> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://www.w3.org/2018/credentials#VerifiableCredential> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#credentialSubject> <did:example:b34ca6cd37bbf23> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#expirationDate> "2029-12-03T12:19:52Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#issuanceDate> "2019-12-03T12:19:52Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#issuer> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
`

	verifier := NewBBSG2SignatureVerifier()
	err = verifier.Verify(&PublicKey{
		Type:  "Bls12381G2Key2020",
		Value: pubKeyBytes,
	}, []byte(msg), sigBytes)

	require.NoError(t, err)
}

//nolint:lll
func TestNewBBSG2SignatureProofVerifier(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pubKeyBase58 := "nEP2DEdbRaQ2r5Azeatui9MG6cj7JUHa8GD7khub4egHJREEuvj4Y8YG8w51LnhPEXxVV1ka93HpSLkVzeQuuPE1mH9oCMrqoHXAKGBsuDT1yJvj9cKgxxLCXiRRirCycki"
	pubKeyBytes := base58.Decode(pubKeyBase58)

	sigBase64 := "ABgA/wYfjSxZz8DBQHTIuX+F0MmeskKbywg6NSMGHOqJ9LvYrfaakmMaPh+UsJxIK1z5v3NuiRP4OGhIbYgjo0KovKMZzluSzCGwzAyXui2hnFlrySj3RP+WNmWd+6QZQ6bEm+pyhNC6VrEMVDxJ2TH7DShbx6GFQ6RLvuS0Xf38GuOhX26+5RJ9RBs5Qaj4/UKsTfc9AAAAdKGdxxloz3ZJ2QnoFlqicO6MviT8yzeyf5gILHg8YUjNIAVJJNsh26kBqIdQkaROpQAAAAIVX5Y1Jy9hgEQgqUld/aGN2uxOLZAJsri9BRRHoFNWkkcF73EV4BE9+Hs+8fuvX0SNDAmomTVz6vSrq58bjHZ+tmJ5JddwT1tCunHV330hqleI47eAqwGuY9hdeSixzfL0/CGnZ2XoV2YAybVTcupSAAAACw03E8CoLBvqXeMV7EtRTwMpKQmEUyAM5iwC2ZaAkDLnFOt2iHR4P8VExFmOZCl94gt6bqWuODhJ5mNCJXjEO9wmx3RNM5prB7Au5g59mdcuuY/GCKmKNt087BoHYG//dEFi4Q+bRpVE5MKaGv/JZd/LmPAfKfuj5Tr37m0m3hx6HROmIv0yHcakQlNQqM6QuRQLMr2U+nj4U4OFQZfMg3A+f6fVS6T18WLq4xbHc/2L1bYhIw+SjXwkj20cGhEBsmFOqj4oY5AzjN1t4gfzb5itxQNkZFVE2IdBP9v/Ck8rMQLmxs68PDPcp6CAb9dvMS0fX5CTTbJHqG4XEjYRaBVG0Ji5g3vTpGVAA4jqOzpTbxKQawA4SvddV8NUUm4N/zCeWMermi3yRhZRl1AXa8BqGO+mXNI7yAPjn1YDoGliQkoQc5B4CYY/5ldP19XS2hV5Ak16AJtD4tdeqbaX0bo="
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	// Case 17 (https://github.com/w3c-ccg/vc-http-api/pull/128)
	msg := `
_:c14n0 <http://purl.org/dc/terms/created> "2021-02-23T19:31:12Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
_:c14n0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/security#BbsBlsSignature2020> .
_:c14n0 <https://w3id.org/security#proofPurpose> <https://w3id.org/security#assertionMethod> .
_:c14n0 <https://w3id.org/security#verificationMethod> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2#zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
<did:example:b34ca6cd37bbf23> <http://schema.org/birthDate> "1958-07-17"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<did:example:b34ca6cd37bbf23> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Person> .
<did:example:b34ca6cd37bbf23> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/citizenship#PermanentResident> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://schema.org/description> "Government of Example Permanent Resident Card." .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://schema.org/name> "Permanent Resident Card" .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/citizenship#PermanentResidentCard> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://www.w3.org/2018/credentials#VerifiableCredential> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#credentialSubject> <did:example:b34ca6cd37bbf23> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#expirationDate> "2029-12-03T12:19:52Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#issuanceDate> "2019-12-03T12:19:52Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<https://issuer.oidp.uscis.gov/credentials/83627465> <https://www.w3.org/2018/credentials#issuer> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
`

	nonceBytes, err := base64.StdEncoding.DecodeString("G/hn9Ca9bIWZpJGlhnr/41r8RB0OO0TLChZASr3QJVztdri/JzS8Zf/xWJT5jW78zlM=")
	require.NoError(t, err)

	verifier := NewBBSG2SignatureProofVerifier(nonceBytes)
	err = verifier.Verify(&PublicKey{
		Type:  "Bls12381G2Key2020",
		Value: pubKeyBytes,
	}, []byte(msg), sigBytes)

	require.NoError(t, err)
}

//nolint:lll
func TestNewBBSG2SignatureProofVerifierCase19(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pubKeyBase58 := "nEP2DEdbRaQ2r5Azeatui9MG6cj7JUHa8GD7khub4egHJREEuvj4Y8YG8w51LnhPEXxVV1ka93HpSLkVzeQuuPE1mH9oCMrqoHXAKGBsuDT1yJvj9cKgxxLCXiRRirCycki"
	pubKeyBytes := base58.Decode(pubKeyBase58)

	sigBase64 := "AAwP/4nFun/RtaXtUVTppUimMRTcEROs3gbjh9iqjGQAsvD+ne2uzME26gY4zNBcMKpvyLD4I6UGm8ATKLQI4OUiBXHNCQZI4YEM5hWI7AzhFXLEEVDFL0Gzr4S04PvcJsmV74BqST8iI1HUO2TCjdT1LkhgPabP/Zy8IpnbWUtLZO1t76NFwCV8+R1YpOozTNKRQQAAAHSpyGry6Rx3PRuOZUeqk4iGFq67iHSiBybjo6muud7aUyCxd9AW3onTlV2Nxz8AJD0AAAACB3FmuAUcklAj5cdSdw7VY57y7p4VmfPCKaEp1SSJTJRZXiE2xUqDntend+tkq+jjHhLCk56zk5GoZzr280IeuLne4WgpB2kNN7n5dqRpy4+UkS5+kiorLtKiJuWhk+OFTiB8jFlTbm0dH3O3tm5CzQAAAAIhY6I8vQ96tdSoyGy09wEMCdWzB06GElVHeQhWVw8fukq1dUAwWRXmZKT8kxDNAlp2NS7fXpEGXZ9fF7+c1IJp"
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	// Case 19 (https://github.com/w3c-ccg/vc-http-api/pull/128)
	msg := `
_:c14n0 <http://purl.org/dc/terms/created> "2021-02-23T19:37:24Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
_:c14n0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://w3id.org/security#BbsBlsSignature2020> .
_:c14n0 <https://w3id.org/security#proofPurpose> <https://w3id.org/security#assertionMethod> .
_:c14n0 <https://w3id.org/security#verificationMethod> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2#zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
<did:key:z5TcESXuYUE9aZWYwSdrUEGK1HNQFHyTt4aVpaCTVZcDXQmUheFwfNZmRksaAbBneNm5KyE52SdJeRCN1g6PJmF31GsHWwFiqUDujvasK3wTiDr3vvkYwEJHt7H5RGEKYEp1ErtQtcEBgsgY2DA9JZkHj1J9HZ8MRDTguAhoFtR4aTBQhgnkP4SwVbxDYMEZoF2TMYn3s#zUC7LTa4hWtaE9YKyDsMVGiRNqPMN3s4rjBdB3MFi6PcVWReNfR72y3oGW2NhNcaKNVhMobh7aHp8oZB3qdJCs7RebM2xsodrSm8MmePbN25NTGcpjkJMwKbcWfYDX7eHCJjPGM> <https://example.org/examples#degree> <urn:bnid:_:c14n0> .
<http://example.gov/credentials/3732> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://example.org/examples#UniversityDegreeCredential> .
<http://example.gov/credentials/3732> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://www.w3.org/2018/credentials#VerifiableCredential> .
<http://example.gov/credentials/3732> <https://www.w3.org/2018/credentials#credentialSubject> <did:key:z5TcESXuYUE9aZWYwSdrUEGK1HNQFHyTt4aVpaCTVZcDXQmUheFwfNZmRksaAbBneNm5KyE52SdJeRCN1g6PJmF31GsHWwFiqUDujvasK3wTiDr3vvkYwEJHt7H5RGEKYEp1ErtQtcEBgsgY2DA9JZkHj1J9HZ8MRDTguAhoFtR4aTBQhgnkP4SwVbxDYMEZoF2TMYn3s#zUC7LTa4hWtaE9YKyDsMVGiRNqPMN3s4rjBdB3MFi6PcVWReNfR72y3oGW2NhNcaKNVhMobh7aHp8oZB3qdJCs7RebM2xsodrSm8MmePbN25NTGcpjkJMwKbcWfYDX7eHCJjPGM> .
<http://example.gov/credentials/3732> <https://www.w3.org/2018/credentials#issuanceDate> "2020-03-10T04:24:12.164Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
<http://example.gov/credentials/3732> <https://www.w3.org/2018/credentials#issuer> <did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2> .
<urn:bnid:_:c14n0> <http://schema.org/name> "Bachelor of Science and Arts"^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#HTML> .
<urn:bnid:_:c14n0> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://example.org/examples#BachelorDegree> .
`

	nonceBytes, err := base64.StdEncoding.DecodeString("lEixQKDQvRecCifKl789TQj+Ii6YWDLSwn3AxR0VpPJ1QV5htod/0VCchVf1zVM0y2E=")
	require.NoError(t, err)

	verifier := NewBBSG2SignatureProofVerifier(nonceBytes)
	err = verifier.Verify(&PublicKey{
		Type:  "Bls12381G2Key2020",
		Value: pubKeyBytes,
	}, []byte(msg), sigBytes)

	require.NoError(t, err)
}

type testSignatureVerifier struct {
	baseSignatureVerifier

	verifyResult error
}

func (v testSignatureVerifier) Verify(*PublicKey, []byte, []byte) error {
	return v.verifyResult
}

func newCryptoSigner(keyType kmsapi.KeyType) (signature.Signer, error) {
	p, err := mockkms.NewProviderForKMS(storage.NewMockStoreProvider(), &noop.NoLock{})
	if err != nil {
		return nil, err
	}

	localKMS, err := localkms.New("local-lock://custom/main/key/", p)
	if err != nil {
		return nil, err
	}

	tinkCrypto, err := tinkcrypto.New()
	if err != nil {
		return nil, err
	}

	return signature.NewCryptoSigner(tinkCrypto, localKMS, keyType)
}

//nolint:lll,goconst
func TestNewPsmsG1SignatureVerifier(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pubKeyBase58 := "128BJq8sv4ns79z7ULtHR1KjtsNZe4MQputxrUsChM2tBaSfzHh2QkPAGHGuMikt3jpa362kM1CfHqL4aCKbyQWNsnCQRREJQywEXXLv3ynYKzYaSn4CRnKRyRcxtQaBfuskT3AHQbdoDvktS5vmWWTiRHFTBTWeuYA19DaTboNXHfc2MaE7zsoJcqqBvu5GXEpWfbd31XDXP6ay7kHQ7xeTk9HwftQTcX3qycaYu4hGNawghFUBQPr4jR38zMkewCjD2XFs7EiFp4383rCGYEGbDWm7Fo3vzRHxmbkq5fpUmSmfPTMDGYDz67hJZfkzJjBo1ZXgzGd8yGwMKwjXnT2WRLqbjFrCnxSB8eYPRuyBLSgaZr8J1SdihNSv8StyWzTRME6nnZfCwPM1fryM9rYsvP2XDrW1tpYgMm6FFskwGpWo2R4HuiBSuxxay2QhUoA1CKMdinf4UjE9fziUtxCGoJkfMyEbaowMgPNsisw8UomBHZt7TZgNFoBsiagQJYMy8GqsNZ9kt8FLfk9cbEpFAz9gb2aTJYwgwcyBUT9zq81x6NYyrgQxjURMojGu1vWk89kxSM6FUkSv9s7ypB2hxBY4daXncKBDU6mkFLZ972VQjdjXAxdmVoKfH14LY7c6rg1ncXGurs58gVv3vMeyNGNSRApMaeRTknTwv459ftE3E2NfujXZQt843ZakGz6hKEqStVcr3CqsTtWFpiEkmbVRqNYXV8wVnM54xbitB1nikcKHCKtQfTShEv37PfeDYQRuBuTcx1A6VXvoQk6CMvNLChd2ST7HgiowpWcm5pDYGyJdiLidsARX6uhTsDtZ5J6VHaW8nnBjdHFdMKLeqnh4okSz147TJqFnDymunrSPZhq2RJzbjPTd8FL891ggRSdAJX1Npk6hEcQSX93bKdnGT8Vdnd7qGcg5LCMxUMC7gHF2iUhqGYzX2a1C6Ky8momf3xhyei1D2tyTY4rkpmsVMLYahCifKSRZD1ydJRXUPsHmmdZP9f4vebDN21p5zjebnSJaN1dgvcfg44wiFRAJKcdMDiRh7HwM63XZkds7S4DagXPQWmD8w8synX7EiwExMR6ppQYn3HfxcZbSyJEafYEpPPPdVc5AGL1gohvMYFfaj7coFgCEqmCSycavm998kQ4HKSFtPGbLbnHRLjDeAXi86wTJFCZBMzpws6MZxSY3LioGzW83Xq2iacJWU5EuaJaiDhz2V7k1yjm94CBKWSfZmsxkZgcS2yx3BXs3kEKokE5d8s1nb68niFQzeNKU1TJYe2p5LYSF11WwPt568un3dCNKQbUCLoYPkiK6kPwrgc7PuZZS8dcrkybmxdo2GVapmMt7VxnyU45kw7afBaKHWht87sAL3zj4ysDk5cwMxFSo58GmScijXmQCtJMThNJuvekGN8Pq9FA7Mnvbdz4HaVLPbEq5QXhRxoFsSfFXmZEgJvRSzbD2nSaz1MKjJ2eRebqqa5QWF15RiXtPBAvv4bA5FQbky3PVhpqXHbyGLYLQ5D7gjeUgrU3JWfgTWjBNmaebXj8gBunZR6Z9QkmjecfUohhVpNNuyYt36B1gPmXgNDES11JX74cUwkWDk7WQAK2SdfWwaeMmbkAiMkDUiYx38ttUxJyKReXxy9BQ7vxUZ3NwKefqA2UmyAKeD6FmEcyZNwJyEVDBJW4of5ro1hfJNvag5H8SWEYeSht6aABV9fKKwX"
	pubKeyBytes := base58.Decode(pubKeyBase58)

	sigBase64 := `BAaJzkPYA52SuOizzgJD8uz5iAAPXztv4U5WZcH9FxekS3AZBya7SNPm_ciIZFM-hRHjjb05C7Y8ilUN8jyALMp7WvXQZiRYreATByz22DgmOxBWyO8fOO5id5x0e0LtnRmwgAuSp5B82MNLxqVC5t_rJVvZCxxp3nf-fMRon7os5KsQWm8s4uNElYm6qwZSyxjS6CWfluQmufoD5j1N6yzKhwLE348QYZDOe2iC8Ca7w_SwbiznAmH1cpUgV4PDYAQNBS_VYo_4cA9QzJNHFEiJiVYE-ElfLbnE0i50VLWTk3TsSynpuiFxgsnFexEXE2IX4wqSI_p_JMivYdy9srT8gngb5H2Ug45W8AUjwUky7Z3xPM0VGqqqYZnQu2F7fdoJaGY-jwzqFFXLMEJI49Af0EMwZhwZYnW0dG1O5BAJXw9a5tPNfZxjHebDXY3NjjAYxJHwtzTScw5lEBVwWxqeo45Sru-gew9MLCoU6VORo04K18zBBke_6R3rre0jEtJgw9c06-r3y5aEOpWOSGCu3y2OI0vKOv_9PMRIJmIyNCXoaSWlzdFWZb_Jwp596po`
	sigBytes, err := base64.RawURLEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	// Case 16 (https://github.com/w3c-ccg/vc-http-api/pull/128)
	msg := `<http://purl.org/dc/terms/created>"2023-04-19T14:17:10.127385543+02:00"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
http://schema.org/birthDate "1958-07-17"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
http://schema.org/familyName "SMITH" .
http://schema.org/gender "Male" .
http://schema.org/givenName "JOHN" .
http://schema.org/image <data:image/png;base64,iVBORw0KGgo...kJggg==> .
https://w3id.org/citizenship#birthCountry "Bahamas" .
https://w3id.org/citizenship#commuterClassification "C1" .
https://w3id.org/citizenship#lprCategory "C09" .
https://w3id.org/citizenship#lprNumber "999-999-999" .
https://w3id.org/citizenship#residentSince "2015-01-01"^^<http://www.w3.org/2001/XMLSchema#dateTime> .`

	verifier := NewPSMSG1SignatureVerifier()
	err = verifier.Verify(&PublicKey{
		Type:  "Bls12381G1Key2022",
		Value: pubKeyBytes,
	}, []byte(msg), sigBytes)

	require.NoError(t, err)
}

//nolint:lll
func TestNewPsmsg1SignatureProofVerifier(t *testing.T) {
	// pkBase58 from did:key:zUC724vuGvHpnCGFG1qqpXb81SiBLu3KLSqVzenwEZNPoY35i2Bscb8DLaVwHvRFs6F2NkNNXRcPWvqnPDUd9ukdjLkjZd3u9zzL4wDZDUpkPAatLDGLEYVo8kkAzuAKJQMr7N2
	pubKeyBase58 := "17Dy7o2xoHXgYErocwsFLfSeNHYidikmFan1dwpi3T5KHiRkryoZAUF7HvAdG2eHah5QTBNnUUtzfrWJYmG8K7YY39zqdW1vAYE8ZzK1YYK4PXEu9Kwby3xBqVPAjQNsAPbB3xUj4iFfCV8dqqdmBKXzGEcT4gZpnA27ggc6zhwY4GhfZokbBvX9BbANzEEnJqcUx6ZnDR3GahUu6H1A3RXD5CSvkF9KAPpMAk126V8GA85gZGvJnRgHqriqWuL79euMP1yoyCeGuShbfjxxnry4YGtwTvG248Q25DkoEhcNFWTKchTGKrWPQU8mScQ3so6N8p7DoMVTYSFFr9X7RWK6FRynfvAZWBpKrWfrVLA4qXF1ntrVgjUSnTQTX4HgZ54P2isLMayhUWHYuq6dqt14ReaFUjEbNCrajGeuGVANh78uGqvmNbb5KmxuPZM2LbpjLD7U3CPkvztsKWAZGyRoMveg5JWQkg7PyGbeT8JcKHq5R4ueGeXf4g9XzktbApLq4snZAkQ1N5yzeUMgH4Ptq767pfBmmmEvZQhy9kCc8yYHLMkDc2cLDT9qeFe8UsoiCo4bbyiajJEcCaHZ5gMot42d2LEmSXxdmxnfUDmWRu3KK6St7C1NooRmfK6kETzBQxLe3a3o9RMMmpMkw1mXM3RKFvFQZrrKVpWQXgRm7JaHCLkD9Zfadq5f3o8rqa7s2xfDyPHXggqh3ebwejqQo5SH48X6vbULzqdJahCTuBiTercyDkdxh7GrkAv8D6PmJ5BzVDrPbfRXTd4w7Yji6vo6Mb"
	pubKeyBytes := base58.Decode(pubKeyBase58)

	sigBase64 := "AAIJAgQGUrZ1lOaRWt3LUXNukpoRF5V0fSuejyJhuqgWNtF9sDfgI5BVSW03FLJrMOxjU1wQ+YrHJONOUGKHunkDT+qQCnsqV1dF+5cI9qPuXSUOOk6qgqiUhMwjzc5tD0h4RY4U0zTB9rF4l5l5yzhuJWKgnFFXQ4X4NOI71k9yl36fKQk4IEMKpuRLybXUuKaFPc4X2nNgz/DI27YXk7nLUdN+C9+1bbb2ajBopCzzAJu3vcyDjDSvH5tIKZOYE1sd+kIEE98+TYfVx/U0xyHewxozwO/Yvne3sKeXGiujShlXoHcM5Z032aemgIXWm0vvIDHPFv/5C3PSIVOxbA5Nhk1o6EVgQ4HWEtMAQrsUP0JftdtF55RsA2wpasLyWsE4TTuGFgUUAwM57KywCfFzk3C6m14dHzqzXN3Zrcg/IKk9orJsGVxjti4Ow6eEb/xtdUjdGVr5ctUwRCaq1S98tV9KQdBeWa7VjXmky1yY3PV4HOa7Qujp3PAkw0iVgRA5uAIv5hi/AuMIwfCcPxciz/E8tOdmeiAbspoy5/ceRAlzpt32vtB6WcArwqECx8IQlWnsAAAAAAAAAAAAAAAAAAAAADNyN7NsEOFg30qRXvXfcAOpnXkcbQVE5YMeKiEdaQJoAAAAAAAAAAAAAAAAAAAAAGOVgVnU0lvuqFFde2Fl4D2gMX8RHs49Atdwt6c1Gcs1AAAAAAAAAAAAAAAAAAAAAA2WEPjVdiOvNNB2Zd05pSohfwwC7y29x67NvcZ0A27+AAAAAAAAAAAAAAAAAAAAAB7w+x9S2nUBMoNMPPt0vfqf3qJZ3clhXOPZTYiCZHrS"
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	// Case 16
	msg := `<http://purl.org/dc/terms/created>"2023-04-04T09:51:39.222650081+02:00"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
http://schema.org/givenName "JOHN" .`

	nonceBytes, err := base64.StdEncoding.DecodeString("G/hn9Ca9bIWZpJGlhnr/41r8RB0OO0TLChZASr3QJVztdri/JzS8Zf/xWJT5jW78zlM=")
	require.NoError(t, err)

	verifier := NewPSMSG1SignatureProofVerifier(nonceBytes)
	err = verifier.Verify(&PublicKey{
		Type:  "Bls12381G1Key2022",
		Value: pubKeyBytes,
	}, []byte(msg), sigBytes)

	require.NoError(t, err)
}
