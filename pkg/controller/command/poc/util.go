package poc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"unsafe"
	vcwalletc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vcwallet"
	vdrc "github.com/hyperledger/aries-framework-go/pkg/controller/command/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk"
	"github.com/hyperledger/aries-framework-go/pkg/kms"

	"os"
)

const (
	ed25519VerificationKey2018 = "Ed25519VerificationKey2018"
	ed25519VerificationKey2020 = "Ed25519VerificationKey2020"
	x25519KeyAgreementKey2019  = "X25519KeyAgreementKey2019"
	bls12381G2Key2020          = "Bls12381G2Key2020"
	bls12381G1Key2022          = "Bls12381G1Key2022"
)



func (o *Command) generateCredentialSubject(proofs []IdProof) (map[string]interface{}, bool, error) {
	result := make(map[string]interface{})
	for _, proof := range proofs {
		if !o.ValidateProof(proof) {
			return nil, false, fmt.Errorf("failed proof validity check at proof %s", proof)
		}
		result[proof.AttrName] = proof.AttrValue //Any other error that results from parsing results in nil, true, err
	}
	return result, true, nil
}

func getReader(v interface{}) (io.Reader, error) {
	vcReqBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(vcReqBytes), nil
}

func randStringBytesMaskImprSrcUnsafe(n int, src rand.Source) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func getUnlockToken(b bytes.Buffer) string {
	var response vcwalletc.UnlockWalletResponse
	err := json.NewDecoder(&b).Decode(&response)
	if err != nil {
		return ""
	}
	return response.Token
}

func parsePufOutput(out []byte) (*jwk.JWK, error) {
	trimmed := trim(out)
	//var res *jwk.JWK
	res := &jwk.JWK{}
	err := res.UnmarshalJSON(trimmed)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func trim(out []byte) []byte {
	res := strings.TrimLeftFunc(string(out), func(r rune) bool {
		return r != '{'
	})
	res = strings.TrimRightFunc(res, func(r rune) bool {
		return r != '}'
	})
	return []byte(res)
}

func getContentID(content []byte) (string, error) {
	//XXX UMU This is needed for now to get the Id, but should not be
	// necessary as it seems to be intended to be kept hidden
	var cid contentID
	if err := json.Unmarshal(content, &cid); err != nil {
		return "", fmt.Errorf("failed to read content to be saved : %w", err)
	}

	key := cid.ID
	if strings.TrimSpace(key) == "" {
		// use document hash as key to avoid duplicates if id is missing
		digest := sha256.Sum256(content)
		return hex.EncodeToString(digest[0:]), nil
	}

	return key, nil
}

type contentID struct {
	ID string `json:"id"`
}

func checkAuthKeyPresent(keys []KeyTypePurpose) bool {
	for _, keyPurpose := range keys {
		if keyPurpose.Purpose == "Authentication" {
			return true
		}
	}
	return false
}

func parseKeyType(keyType KeyTypeModel) kms.KeyType {
	switch keyType.Type {
	case ed25519VerificationKey2018:
		return kms.ED25519Type
	case ed25519VerificationKey2020:
		return kms.ED25519Type
	case bls12381G2Key2020:
		return kms.BLS12381G2Type
	case bls12381G1Key2022:
		if keyType.Attrs[0] == "0" {
			return ""
		}
		return kms.BLS12381G1Type
	default:
		return ""
	}
}

func getDID(response vdrc.Document) string {
	var cid contentID
	if err := json.Unmarshal(response.DID, &cid); err != nil {
		return ""
	}
	return cid.ID
}

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, fmt.Errorf("Environment variable is empty")
	}
	return v, nil
}

// generateRandomString creates a random alphanumeric string of specified length.
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length) // Generate a byte slice to hold the random bytes.
	if _, err := rand.Read(bytes); err != nil {
		return "", err // Return an error if the random byte generation fails.
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))] // Map each byte to a character in the charset.
	}
	return string(bytes), nil
}



