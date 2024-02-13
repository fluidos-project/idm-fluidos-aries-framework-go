/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignatureproof2022_test

import (
	_ "embed"
	"encoding/base64"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignatureproof2022"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/internal/ldtestutil"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	//go:embed testdata/selective_disclosure_vc.jsonld
	vcDoc string
	//go:embed testdata/expected_canonical_doc.rdf
	expectedCanonicalDoc string
	//go:embed testdata/expected_digest_doc.rdf
	expectedDigestDoc string
)

func TestSuite(t *testing.T) {
	blsVerifier := &testVerifier{}

	blsSuite := psmsblssignatureproof2022.New(suite.WithCompactProof(), suite.WithVerifier(blsVerifier))

	//nolint:lll
	pkBase64 := "AAUFBBA4Mlps/67GkyzVz0DLoim02qDqIoYO7EKuW4WhS2gZqKaWSYEiXSP7aWcOKLkHtwccBK4qcRkxVi88o+OwFyXuYPbVEtgyzSwZU9OWDyMy2wbmhcvcR1jybjaiDMgoxAQXJ8Xt/IS7fXtnUOzFY+lTaDkzzM+JjMI61JShGY4Eus9lrQAoncWuWUjELDqkhTUEs6Af6wQBl5sqPg8qkgqJSy7hqWu/MYfiPI72QvtsJsxmtQ/A3cGsQhHmW3r1hXMEDJCf3DR3/K9ACL5H54d7mrwUJSYnyjl00POu1Lj/ta4PltriEoW68xgaBJIAFQSMBJct6bZ5RreA/MBIUOxdt6baELPPh5cseayRf+nNn2FTqZqdyQw/5bheYPMR1cIsBBiFGvSlCoygiDXjzPBemUdoOQV72tst6fq+eS8EnlEPP+InwfQu8Dn1FdgqJi8JwxUhAVpjKYT1zfQcqsd1GO9Dj1hXmQ2cPvXGuiX1S3+wqAkaA6+EzQdHKHds7LoUrQQKs1qGalBfKK9RG7VxvHF/wwruTQK0m68pHblbMcEC36vkqIdvvBM83E1UNsqX83QMxZ+I3HEBqkAcoalkxwpwQ6uVEW+mRejQ+EBG/uK97L99k04y3j0fxe/EvM2jO/0EFUlyvjVYIVwwBQggFqgSPqFoAw4byd6CytslbjT6LqvRMEIPgoVnSecfZXAUwJSUBUt/6F4dKuWQ4xw6lED7rfM+TJ1v1T486gpU64tig7L3Ur3govGa91RKNsVVaLD7BAUyKU6nn1mZvJBVnOtMC0Z4s4f5gGLAhVw85fxaldN+3KWTg5vcQ7uGEo+GGaw6xwPfoXeHEnQvoT0J+5BRijPIwcPj2/lqswgpSYNWXkk4psXpJdniC0hyAJxjGQx3EwQGcy3VPsJWO7QA/MDQe4aEu3kq3v4ZhcoHtp1mZ2daZlVoiSv2KOroT/AHXzUQgOcZCV4GwGUuc0FeE2a6R/rfQkYKLbu38tzf3KF1s+i3C8hmfAKICqASz0OqDQDOuzQ"
	pkBytes, err := base64.RawStdEncoding.DecodeString(pkBase64)
	require.NoError(t, err)

	blsPSMSPublicKey := &verifier.PublicKey{
		Type:  "Bls12381G1Key2022",
		Value: pkBytes,
	}

	v, err := verifier.New(&testKeyResolver{
		publicKey: blsPSMSPublicKey,
	}, blsSuite)
	require.NoError(t, err)
	err = v.Verify([]byte(vcDoc), ldtestutil.WithDocumentLoader(t))
	require.NoError(t, err)
	require.Equal(t, expectedDigestDoc, blsVerifier.doc)
}

func TestSignatureSuite_GetDigest(t *testing.T) {
	digest := psmsblssignatureproof2022.New().GetDigest([]byte(expectedCanonicalDoc))
	require.NotNil(t, digest)
	require.Equal(t, []byte(expectedDigestDoc), digest)
}

func TestSignatureSuite_Accept(t *testing.T) {
	ss := psmsblssignatureproof2022.New()
	accepted := ss.Accept("PsmsBlsSignatureProof2022")
	require.True(t, accepted)

	accepted = ss.Accept("RsaSignature2018")
	require.False(t, accepted)
}

type testVerifier struct {
	err error
	doc string
}

func (v *testVerifier) Verify(_ *verifier.PublicKey, doc, _ []byte) error {
	v.doc = string(doc)
	return v.err
}

type testKeyResolver struct {
	publicKey *verifier.PublicKey
	variants  map[string]*verifier.PublicKey
	err       error
}

func (r *testKeyResolver) Resolve(id string) (*verifier.PublicKey, error) {
	if r.err != nil {
		return nil, r.err
	}

	if len(r.variants) > 0 {
		return r.variants[id], nil
	}

	return r.publicKey, r.err
}
