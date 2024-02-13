/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package verifiable

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignature2022"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignatureproof2022"
	jsonutil "github.com/hyperledger/aries-framework-go/pkg/doc/util/json"
	"github.com/stretchr/testify/require"
)

//
//nolint:lll
func TestCredential_GeneratePSMSSelectiveDisclosure(t *testing.T) {
	seed := []byte("seeedForKey")
	seedAnother := []byte("seeedForAnotherKey")
	nattr := 10

	vcJSON := `
			{
			 "@context": [
			   "https://www.w3.org/2018/credentials/v1",
			   "https://w3id.org/citizenship/v1",
			   "https://ssiproject.inf.um.es/security/psms/v1"
			 ],
			 "id": "https://issuer.oidp.uscis.gov/credentials/83627465",
			 "type": [
			   "VerifiableCredential",
			   "PermanentResidentCard"
			 ],
			 "issuer": "did:example:489398593",
			 "identifier": "83627465",
			 "name": "Permanent Resident Card",
			 "description": "Government of Example Permanent Resident Card.",
			 "issuanceDate": "2019-12-03T12:19:52Z",
			 "expirationDate": "2029-12-03T12:19:52Z",
			 "credentialSubject": {
			   "id": "did:example:b34ca6cd37bbf23",
			   "type": [
			     "PermanentResident",
			     "Person"
			   ],
			   "givenName": "JOHN",
			   "familyName": "SMITH",
			   "gender": "Male",
			   "image": "data:image/png;base64,iVBORw0KGgokJggg==",
			   "residentSince": "2015-01-01",
			   "lprCategory": "C09",
			   "lprNumber": "999-999-999",
			   "commuterClassification": "C1",
			   "birthCountry": "Bahamas",
			   "birthDate": "1958-07-17"
			 }
			}
			`

	pubKey, privKey, err := psms12381g1pub.New().GenerateKeyPair(seed, nattr)
	require.NoError(t, err)

	pubKeyBytes, err := pubKey.Marshal()
	require.NoError(t, err)

	vc, err := parseTestCredential(t, []byte(vcJSON))
	require.NoError(t, err)
	require.Len(t, vc.Proofs, 0)

	signVCWithPSMS(t, privKey, pubKeyBytes, vc)
	signVCWithEd25519(t, vc)

	a, _ := vc.MarshalJSON()
	fmt.Println(string(a))
	fmt.Println(base58.Encode(pubKeyBytes))
	revealJSON := `
		{
		  "@context": [
		    "https://www.w3.org/2018/credentials/v1",
		    "https://w3id.org/citizenship/v1",
		    "https://ssiproject.inf.um.es/security/psms/v1"
		  ],
		  "type": ["VerifiableCredential", "PermanentResidentCard"],
		  "@explicit": true,
		  "identifier": {},
		  "issuer": {},
		  "issuanceDate": {},
		  "credentialSubject": {
		    "@explicit": true,
		    "type": ["PermanentResident", "Person"],
		    "givenName": {},
		    "familyName": {},
		    "gender": {}
		  }
		}
		`

	revealDoc, err := jsonutil.ToMap(revealJSON)
	require.NoError(t, err)

	nonce := []byte("nonce")

	vcOptions := []CredentialOpt{WithJSONLDDocumentLoader(createTestDocumentLoader(t)), WithPublicKeyFetcher(
		SingleKey(pubKeyBytes, "Bls12381G1Key2022"))}

	vcWithSelectiveDisclosure, err := vc.GeneratePSMSSelectiveDisclosure(revealDoc, nonce, vcOptions...)
	require.NoError(t, err)
	require.NotNil(t, vcWithSelectiveDisclosure)
	require.Len(t, vcWithSelectiveDisclosure.Proofs, 1)

	vcSelectiveDisclosureBytes, err := json.Marshal(vcWithSelectiveDisclosure)
	require.NoError(t, err)

	sigSuite := psmsblssignatureproof2022.New(
		suite.WithCompactProof(),
		suite.WithVerifier(psmsblssignatureproof2022.NewG1PublicKeyVerifier(nonce)))

	vcVerified, err := parseTestCredential(t, vcSelectiveDisclosureBytes,
		WithEmbeddedSignatureSuites(sigSuite),
		WithPublicKeyFetcher(SingleKey(pubKeyBytes, "Bls12381G1Key2022")),
	)
	require.NoError(t, err)
	require.NotNil(t, vcVerified)

	// error cases
	t.Run("failed generation of selective disclosure", func(t *testing.T) {
		var (
			anotherPubKey      *psms12381g1pub.PublicKey
			anotherPubKeyBytes []byte
		)

		anotherPubKey, _, err = psms12381g1pub.New().GenerateKeyPair(seedAnother, nattr)
		require.NoError(t, err)

		anotherPubKeyBytes, err = anotherPubKey.Marshal()
		require.NoError(t, err)

		vcWithSelectiveDisclosure, err = vc.GeneratePSMSSelectiveDisclosure(revealDoc, nonce,
			WithJSONLDDocumentLoader(createTestDocumentLoader(t)),
			WithPublicKeyFetcher(SingleKey(anotherPubKeyBytes, "Bls12381G1Key2022")))
		require.Error(t, err)
		require.Contains(t, err.Error(), "create VC selective disclosure")
		require.Empty(t, vcWithSelectiveDisclosure)
	})

	t.Run("public key fetcher is not passed", func(t *testing.T) {
		vcWithSelectiveDisclosure, err = vc.GeneratePSMSSelectiveDisclosure(revealDoc, nonce)
		require.Error(t, err)
		require.EqualError(t, err, "public key fetcher is not defined")
		require.Empty(t, vcWithSelectiveDisclosure)
	})

	t.Run("Reveal document with hidden VC mandatory field", func(t *testing.T) {
		revealJSONWithMissingIssuer := `
		{
		  "@context": [
		    "https://www.w3.org/2018/credentials/v1",
		    "https://w3id.org/citizenship/v1",
		    "https://w3id.org/security/bbs/v1"
		  ],
		  "type": ["VerifiableCredential", "PermanentResidentCard"],
		  "@explicit": true,
		  "identifier": {},
		  "issuanceDate": {},
		  "credentialSubject": {
		    "@explicit": true,
		    "type": ["PermanentResident", "Person"],
		    "givenName": {},
		    "familyName": {},
		    "gender": {}
		  }
		}
		`

		revealDoc, err = jsonutil.ToMap(revealJSONWithMissingIssuer)
		require.NoError(t, err)

		vcWithSelectiveDisclosure, err = vc.GeneratePSMSSelectiveDisclosure(revealDoc, nonce, vcOptions...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "issuer is required")
		require.Nil(t, vcWithSelectiveDisclosure)
	})

	t.Run("VC with no embedded proof", func(t *testing.T) {
		vc.Proofs = nil
		vcWithSelectiveDisclosure, err = vc.GeneratePSMSSelectiveDisclosure(revealDoc, nonce, vcOptions...)
		require.Error(t, err)
		require.EqualError(t, err, "expected at least one proof present")
		require.Empty(t, vcWithSelectiveDisclosure)
	})
}

func signVCWithPSMS(t *testing.T, privKey *psms12381g1pub.PrivateKey, pubKeyBytes []byte, vc *Credential) {
	t.Helper()

	psmsSigner, err := newPSMSSigner(privKey)
	require.NoError(t, err)

	sigSuite := psmsblssignature2022.New(
		suite.WithSigner(psmsSigner),
		suite.WithVerifier(psmsblssignature2022.NewG1PublicKeyVerifier()))

	ldpContext := &LinkedDataProofContext{
		SignatureType:           "PsmsBlsSignature2022",
		SignatureRepresentation: SignatureProofValue,
		Suite:                   sigSuite,
		VerificationMethod:      "did:example:123456#key1",
	}

	err = vc.AddLinkedDataProof(ldpContext, jsonld.WithDocumentLoader(createTestDocumentLoader(t)))
	require.NoError(t, err)

	vcSignedBytes, err := json.Marshal(vc)
	require.NoError(t, err)
	require.NotEmpty(t, vcSignedBytes)

	vcVerified, err := parseTestCredential(t, vcSignedBytes,
		WithEmbeddedSignatureSuites(sigSuite),
		WithPublicKeyFetcher(SingleKey(pubKeyBytes, "Bls12381G1Key2022")),
	)
	require.NoError(t, err)
	require.NotEmpty(t, vcVerified)
}

type psmsSigner struct {
	privKeyBytes []byte
}

func newPSMSSigner(privKey *psms12381g1pub.PrivateKey) (*psmsSigner, error) {
	privKeyBytes, err := privKey.Marshal()
	if err != nil {
		return nil, err
	}

	return &psmsSigner{privKeyBytes: privKeyBytes}, nil
}

func (s *psmsSigner) Sign(data []byte) ([]byte, error) {
	msgs := s.textToLines(string(data))

	return psms12381g1pub.New().Sign(msgs, s.privKeyBytes)
}

func (s *psmsSigner) Alg() string {
	return ""
}

func (s *psmsSigner) textToLines(txt string) [][]byte {
	lines := strings.Split(txt, "\n")
	linesBytes := make([][]byte, 0, len(lines))

	for i := range lines {
		if strings.TrimSpace(lines[i]) != "" {
			linesBytes = append(linesBytes, []byte(lines[i]))
		}
	}

	return linesBytes
}
