/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignature2022

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk"
	sigverifier "github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
)

//nolint:lll
func TestNewG1PublicKeyVerifier(t *testing.T) {
	verifier := NewG1PublicKeyVerifier()

	pkBase64 := "AAUFBBA4Mlps/67GkyzVz0DLoim02qDqIoYO7EKuW4WhS2gZqKaWSYEiXSP7aWcOKLkHtwccBK4qcRkxVi88o+OwFyXuYPbVEtgyzSwZU9OWDyMy2wbmhcvcR1jybjaiDMgoxAQXJ8Xt/IS7fXtnUOzFY+lTaDkzzM+JjMI61JShGY4Eus9lrQAoncWuWUjELDqkhTUEs6Af6wQBl5sqPg8qkgqJSy7hqWu/MYfiPI72QvtsJsxmtQ/A3cGsQhHmW3r1hXMEDJCf3DR3/K9ACL5H54d7mrwUJSYnyjl00POu1Lj/ta4PltriEoW68xgaBJIAFQSMBJct6bZ5RreA/MBIUOxdt6baELPPh5cseayRf+nNn2FTqZqdyQw/5bheYPMR1cIsBBiFGvSlCoygiDXjzPBemUdoOQV72tst6fq+eS8EnlEPP+InwfQu8Dn1FdgqJi8JwxUhAVpjKYT1zfQcqsd1GO9Dj1hXmQ2cPvXGuiX1S3+wqAkaA6+EzQdHKHds7LoUrQQKs1qGalBfKK9RG7VxvHF/wwruTQK0m68pHblbMcEC36vkqIdvvBM83E1UNsqX83QMxZ+I3HEBqkAcoalkxwpwQ6uVEW+mRejQ+EBG/uK97L99k04y3j0fxe/EvM2jO/0EFUlyvjVYIVwwBQggFqgSPqFoAw4byd6CytslbjT6LqvRMEIPgoVnSecfZXAUwJSUBUt/6F4dKuWQ4xw6lED7rfM+TJ1v1T486gpU64tig7L3Ur3govGa91RKNsVVaLD7BAUyKU6nn1mZvJBVnOtMC0Z4s4f5gGLAhVw85fxaldN+3KWTg5vcQ7uGEo+GGaw6xwPfoXeHEnQvoT0J+5BRijPIwcPj2/lqswgpSYNWXkk4psXpJdniC0hyAJxjGQx3EwQGcy3VPsJWO7QA/MDQe4aEu3kq3v4ZhcoHtp1mZ2daZlVoiSv2KOroT/AHXzUQgOcZCV4GwGUuc0FeE2a6R/rfQkYKLbu38tzf3KF1s+i3C8hmfAKICqASz0OqDQDOuzQ"
	pkBytes, err := base64.RawStdEncoding.DecodeString(pkBase64)
	require.NoError(t, err)

	sigBase64 := "BA1JICdnrQVt9wEoBpbAWrUW6yj5JH/5VxpVsvJqpu9WIwT8Eoco8X/BRodb+P90hw5PO3AoebeHcjbRfW42VQ/0isnJ08dCBCQSUow6Vn51i9+0QNd3aPk3Ort+Ui5CPBIR6RFjSzCHdbpH9Kh6uLfCw2JqrhFtpEjKCctGb8Ns/h5nZRY0ZsR4SzIJneAENwSqrHfsIWJWsubLDyWG9n8YQf1GdgD5OT/ewazaKLKscWbNE8lFhJ7K3tcQyNjrhQQAvhYB8dJfHBHH5XqtNDjsCLCNco9+tihMQUMk3wFDZ9xkc6SfbXUCrdwtYltTHhYEKK2sTFZSl9hLlLP+FrOxQImzYtLgg3+e4FbUZugtHphLLrSlx30841BKRIsYt2oNjWcHggLXuEdDB7AnpZW4l/SHhst6XjhzcpNvhs7z2IeZaqpEbMJqda/LdnDWdGELqZNIvzDt4x/TxXZDwX+ixPtSb6RNK39qRiDI/b9L1n7NGKzS6PsQNk7tSqsvg5pk03Mu8famOPDUJ8ZnbniLCS+Inx7sx8JfjCQyJTZzAGHmTyZaSrYDE+uLSEGFvOY="
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	msg := `
epoch123123123123
test message1
test message2
test message3
test message4
test message5
`

	//key, err := psms12381g1pub.UnmarshalPublicKey(pkBytes)
	//require.NoError(t, err)
	//fmt.Println(key.Nattrs())

	err = verifier.Verify(&sigverifier.PublicKey{
		Type:  "Bls12381G1Key2022",
		Value: pkBytes,
	}, []byte(msg), sigBytes)
	require.NoError(t, err)

	err = verifier.Verify(&sigverifier.PublicKey{
		Type:  "NotBls12381G1Key2022",
		Value: pkBytes,
	}, []byte(msg), sigBytes)
	require.Error(t, err)
	require.EqualError(t, err, "a type of public key is not 'Bls12381G1Key2022'")

	// Success as we now support JWK for Bls12381G2Key2020.
	err = verifier.Verify(&sigverifier.PublicKey{
		Type: "Bls12381G1Key2022",
		JWK: &jwk.JWK{
			Kty: "EC",
			Crv: "BLS12381_G1",
		},
		Value: pkBytes,
	}, []byte(msg), sigBytes)
	require.NoError(t, err)
}
