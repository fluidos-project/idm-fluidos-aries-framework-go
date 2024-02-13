/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignatureproof2022_test

import (
	"encoding/base64"
	"testing"

	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/psmsblssignatureproof2022"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/stretchr/testify/require"
)

//nolint:lll
func TestNewG1PublicKeyVerifier(t *testing.T) {
	nonce := "G/hn9Ca9bIWZpJGlhnr/41r8RB0OO0TLChZASr3QJVztdri/JzS8Zf/xWJT5jW78zlM="
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	require.NoError(t, err)
	publicKeyVerifier := psmsblssignatureproof2022.NewG1PublicKeyVerifier(nonceBytes)

	pkBase64 := "AAMDBAReybonZuFepTieuaL86VAsUY0XPBYyrkawuEQSHwdty3k0DR4IUlajbyhU2C92aANm+U2pOx8TTAILQz78uME/w3fjP5QjQk1i0Ylz/kzw5vU9NkWGrJxG1MEfvGc7HwQFeBEfQD+G3ft3hP7ocTZ4zQ2KvA7x8vkBgfY1obYvJMCIo+r0J+CnSiKcOE9UePQCVmLvSRUHGe/cjhYUKRHKijbhWUHydJNVq9cuO+052eufu1hDzdfipF9ZvFuT/r0EDGPwjkTMAiknEJ2jhdpuEZA0ihSM4s4bxK94WITW2ip3t98BJFKTpu+P3q3GRf6QA/cdT0/AQdrq1H62CyIRTNdr7LwfxwisEFmU1tnV5QxK+mfs8FZtkakcR95f3zwzBAAAHOcY/op5lu3GwnEI22AWP1oLxpp6hQvxvK9acDCLmAK72AXAbG3dYuGleymQ+w82dDIEo+R89AMxGurtfqO+PgXtYl/1/qCCqGN4vrOuU1lSjF8856mnAWcWtLNz7gQE8hWUkJX5AE87pdu/uZuhI/zjKvV3+bSc0WWsRbnz3zk84nr1N5hkJ/tZ4u/5Tq8Hudjxh09yb4nCE59JV41ODO1dKYQqPhW2IQyNtPyRYPQF4VfupEKSwrMoOq4JSswEBt7kNIL8itoJ0wuRDz2O97+kyOI9wokl3EVoVc0rIeDKUL5OXaSDMvJrMR42mShGFe/UoNIZ8VyQdB459sEWTQvrN9ZyenscGS2BMo3LSOHnNlp6K0ZEOR0O7a6K8Z4+"
	pkBytes, err := base64.RawStdEncoding.DecodeString(pkBase64)
	require.NoError(t, err)

	sigBase64 := "AAIJAgQHG46/tbJImwduuT1OzH0NKkqayc2P9yqyuCpnq1f0ZgRhB08oYgSy8CYPGRlzZLgOphiQa1SMeW8DMRbkuVuwDuPgol7AF1y8OrE6J8Mtbow42WJNIKEWFoyHZgCibUIY8gvV/kE2AJiz4lDPtDyYCi4cmoFOssATA+wIAEul6VOr+TEgniFBdZy+wuC++rYYN0f2C1UEUPDJPWPYMHppHH/z9Irp6cnjnJUiwNT++XiV5dudxpqjoggx2A58LPcEGS0Qc6sMmB/Tv0Y7Pjeso034fRa+DlCp+FuZVuSsUkhLs+xlpV3+oXvkZVsOeCjUGGneifkh4yYMd2qfA8cvbUKLtaIf17WK20OD4P4EMel/BTwkLpYUedfNjW/CE4zsBDy2CQhlnD1gk5dxQ/N3JpNYHrnZgCBCNjA6AYnHu8qk/008zp3G2xzKuT3yLytUFEXORHyLtZlvaUeNi0lsgkJaNAKw90dL6cFgORuNOq5cHgLEt+G/zbtK+ROzts5llT+8p/vFAbHXr7ghsty3VbxPFZWIE87Tzo1DX2XD+wpHRLLMnzdaUr/rI/VFAGkhAAAAAAAAAAAAAAAAAAAAADs4wRZZI7bCyqbhVXbeDnXeN6y5sea/wsSBeDwy8EpZAAAAAAAAAAAAAAAAAAAAAD9Cb8DdS/p3zibf3D90njTJ1dmXlbL6uFOJz98mewtdAAAAAAAAAAAAAAAAAAAAAG9PMKLigdT3FGZPRv53xVh0ORw5yHaLm6AB8nW+f/61AAAAAAAAAAAAAAAAAAAAAEzO56hQsQ+czDoTecSKQhyIh4gnaBQ2F1uW0GTkrDDE"
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	require.NoError(t, err)

	msg := `<http://purl.org/dc/terms/created>"2023-04-04T09:51:39.222650081+02:00"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
http://schema.org/givenName "JOHN" .`

	err = publicKeyVerifier.Verify(&verifier.PublicKey{
		Type:  "Bls12381G1Key2022",
		Value: pkBytes,
	}, []byte(msg), sigBytes)
	require.NoError(t, err)

	err = publicKeyVerifier.Verify(&verifier.PublicKey{
		Type:  "NotBls12381G1Key2022",
		Value: pkBytes,
	}, []byte(msg), sigBytes)
	require.Error(t, err)
	require.EqualError(t, err, "a type of public key is not 'Bls12381G1Key2022'")

	// Failed as we do not support JWK for Bls12381G2Key2020.
	err = publicKeyVerifier.Verify(&verifier.PublicKey{
		Type: "Bls12381G1Key2022",
		JWK: &jwk.JWK{
			Kty: "EC",
			Crv: "BLS12381_G1",
		},
	}, []byte(msg), sigBytes)
	require.Error(t, err)
	require.EqualError(t, err, "verifier does not match JSON Web Key")
}
