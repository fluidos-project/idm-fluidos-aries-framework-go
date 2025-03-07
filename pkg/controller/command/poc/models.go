package poc

import (
	"encoding/json"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

// Model for puf authentication program output
type PufAuthResult struct {
	Kty string `json:"kty,omitempty"`
	Use string `json:"use,omitempty"`
	Crv string `json:"crv,omitempty"`
	Kid string `json:"kid,omitempty"`
	X   string `json:"x,omitempty"`
}

//rModel for newDID method input

type SignContractArgs struct {
	Contract    json.RawMessage `json:"contract,omitempty"`
	ContractJWT string          `json:"contractJWT,omitempty"`
}

type VerifyContractArgs struct {
	Contract json.RawMessage `json:"contract,omitempty"`
}

type SignJWTContentArgs struct {
	Content json.RawMessage `json:"content,omitempty"`
}

type SignJWTContentResult struct {
	SignedJWTContent string `json:"signedJWTContent,omitempty"`
}

type SignContractResult struct {
	SignedContract string `json:"signedContract,omitempty"`
}

type VerifyContractSignatureResult struct {
	VerifiedChain   bool                   `json:"verifiedChain,omitempty"`
	Signatures      []JWTSignature         `json:"signatures,omitempty"`
	ContractContent map[string]interface{} `json:"contractContent,omitempty"`
}

type VerifyJWTContentArgs struct {
	JWT string `json:"jwt,omitempty"`
}

type VerifyContractSignatureArgs struct {
	Contract string `json:"contract,omitempty"`
}

type verifyJWTContentResult struct {
	Verified bool `json:"verified,omitempty"`
}

type NewDIDArgs struct {
	Keys []KeyTypePurpose `json:"keys,omitempty"`
	Name string           `json:"name,omitempty"`
}
type TestingCallResult struct {
	DoenrolmentMem uint64 `json:"doenrolmentMem,omitempty"`
	GenerateVPMem  uint64 `json:"generateVPMem,omitempty"`
	VerifyMem      uint64 `json:"verifyMem,omitempty"`
}

// Model for newDID method output
type NewDIDResult struct {
	DIDDoc json.RawMessage `json:"didDoc,omitempty"`
}

// Model keytype/purpose pair
type KeyTypePurpose struct {
	KeyType KeyTypeModel `json:"keyType,omitempty"`
	Purpose string       `json:"purpose,omitempty"`
}

// Model keytype
type KeyTypeModel struct {
	Type  string   `json:"keytype,omitempty"`
	Attrs []string `json:"attrs,omitempty"`
}

// Model for enrolDevice method input
type DoDeviceEnrolmentArgs struct {
	Url      string    `json:"url,omitempty"`
	TheirDID string    `json:"theirDID,omitempty"`
	IdProofs []IdProof `json:"idProofs,omitempty"`
}

// Model for enrolDevice method output
type DoDeviceEnrolmentResult struct {
	Credential    json.RawMessage `json:"credential,omitempty"`
	CredStorageId string          `json:"credStorageId,omitempty"`
}

// Model for idProof
type IdProof struct {
<<<<<<< HEAD
	AttrName  string      `json:"attrName,omitempty"`
	AttrValue interface{} `json:"attrValue,omitempty"`
	ProofData string      `json:"proofData,omitempty"`
=======
	AttrName  string          `json:"attrName,omitempty"`
	AttrValue interface{}     `json:"attrValue,omitempty"`
	ProofData string `json:"proofData,omitempty"`
>>>>>>> dev
}

// Model for GenerateVP method input
type GenerateVPArgs struct {
	CredId       string       `json:"credId,omitempty"` //TODO UMU: How do we decide which credential is gonna be presented?
	QueryByFrame QueryByFrame `json:"querybyframe,omitempty"`
}

// Model for GetVCredential method input
type GetVCredentialArgs struct {
	CredId string `json:"credId,omitempty"`
}

type RequestBodyVP struct {
	CredId string       `json:"credId,omitempty"`
	Frame  FrameFluidos `json:"querybyframe,omitempty"`
}

type QueryByFrame struct {
	Frame FrameFluidos `json:"frame,omitempty"`
}

// type FrameFluidos struct {
// 	Context           []string `json:"@context,omitempty"`
// 	Type              []string `json:"type,omitempty"`
// 	Explicit          bool     `json:"@explicit,omitempty"`
// 	Identifier        struct{} `json:"identifier,omitempty"`
// 	Issuer            struct{} `json:"issuer,omitempty"`
// 	IssuanceDate      struct{} `json:"issuanceDate,omitempty"`
// 	CredentialSubject struct {
// 		Explicit   bool     `json:"@explicit,omitempty"`
// 		HolderRole struct{} `json:"holderRole,omitempty"`
// 	} `json:"credentialSubject,omitempty"`
// }

type FrameFluidos struct {
	Context           []string               `json:"@context,omitempty"`
	Type              []string               `json:"type,omitempty"`
	Explicit          bool                   `json:"@explicit,omitempty"`
	Identifier        struct{}               `json:"identifier,omitempty"`
	Issuer            struct{}               `json:"issuer,omitempty"`
	IssuanceDate      struct{}               `json:"issuanceDate,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject,omitempty"`
}

type FrameFluidosVC struct {
	Context           []string               `json:"@context,omitempty"`
	Type              []string               `json:"type,omitempty"`
	Explicit          bool                   `json:"@explicit,omitempty"`
	Identifier        string                 `json:"id,omitempty"`
	Issuer            string                 `json:"issuer,omitempty"`
	IssuanceDate      string                 `json:"issuanceDate,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject,omitempty"`
}

// Model for GenerateVP method output
type GenerateVPResult struct {
	Results []*verifiable.Presentation `json:"results"`
}

type GenerateVPResultCustom struct {
	Results []*json.RawMessage `json:"results,omitempty"`
}

type VerifiablePresentation struct {
	Context              []string         `json:"@context,omitempty"`
	Type                 []string         `json:"type,omitempty"`
	VerifiableCredential []FrameFluidosVC `json:"verifiableCredential,omitempty"`
}

type GetVCredentialResult struct {
	Credential json.RawMessage `json:"credential,omitempty"`
}

// Model for AcceptEnrolment method input
type AcceptEnrolmentArgs struct {
	IdProofs []IdProof `json:"idProofs,omitempty"`
}

// Model for AcceptEnrolment method output
type AcceptEnrolmentResult struct {
	Credential json.RawMessage `json:"credential,omitempty"`
}

// Model for VerfyCredential method input
type VerifyCredentialArgs struct {
	CredentialString string `json:"credential,omitempty"`
	Endpoint         string `json:"endpoint,omitempty"`
	Method           string `json:"method,omitempty"`
}

// Model for VerifyCredential method output
type VerifyCredentialResult struct {
	Result      bool   `json:"result"`
	AccessToken string `json:"accessToken"`

	Error string `json:"error,omitempty"`
}

type GetTrustedIssuerListResult struct {
	TrustedIssuers []TrustedIssuer `json:"trustedIssuers,omitempty"`
}

type TrustedIssuer struct {
	DID       string `json:"did,omitempty"`
	IssuerUrl string `json:"issuerUrl,omitempty"`
}
type decodeJWTResult struct {
	Header  map[string]interface{}
	Payload map[string]interface{}
}

type JWTSignature struct {
	Did      string `json:"did,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}
