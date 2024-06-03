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

type SignJWTContentArgs struct {
	Content json.RawMessage `json:"content,omitempty"`
}

type SignJWTContentResult struct {
	SignedJWTContent string `json:"signedJWTContent,omitempty"`
}

type VerifyJWTContentArgs struct {
	JWT string `json:"jwt,omitempty"`
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
	GenerateVPMem uint64 `json:"generateVPMem,omitempty"`
	VerifyMem uint64 `json:"verifyMem,omitempty"`
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
	Type  string `json:"keytype,omitempty"`
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
	AttrName  string          `json:"attrName,omitempty"`
	AttrValue interface{}     `json:"attrValue,omitempty"`
	ProofData json.RawMessage `json:"proofData,omitempty"`
}

// Model for GenerateVP method input
type GenerateVPArgs struct {
	CredId string `json:"credId,omitempty"` //TODO UMU: How do we decide which credential is gonna be presented?
	QueryByFrame QueryByFrame `json:"querybyframe,omitempty"`
}

// Model for GetVCredential method input
type GetVCredentialArgs struct {
	CredId string `json:"credId,omitempty"`
}

type RequestBodyVP struct {
	CredId string `json:"credId,omitempty"`
	Frame FrameFluidos `json:"querybyframe,omitempty"`
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
	Context           []string `json:"@context,omitempty"`
	Type              []string `json:"type,omitempty"`
	Explicit          bool     `json:"@explicit,omitempty"`
	Identifier        struct{} `json:"identifier,omitempty"`
	Issuer            struct{} `json:"issuer,omitempty"`
	IssuanceDate      struct{} `json:"issuanceDate,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject,omitempty"`
}

// Model for GenerateVP method output
type GenerateVPResult struct {
	Results []*verifiable.Presentation `json:"results"`
	
}




type GenerateVPResultCustom struct {
	Results []*json.RawMessage `json:"results,omitempty"`
	
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
}

// Model for VerifyCredential method output
type VerifyCredentialResult struct {
	Result bool `json:"result,omitempty"`

	Error string `json:"error,omitempty"`
}

type GetTrustedIssuerListResult struct	{
	TrustedIssuers []TrustedIssuer `json:"trustedIssuers,omitempty"`
}

type TrustedIssuer struct {
	DID string `json:"did,omitempty"`
	IssuerUrl string `json:"issuerUrl,omitempty"`
}