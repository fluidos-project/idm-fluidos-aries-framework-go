/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package poc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
	"encoding/json"

	"github.com/hyperledger/aries-framework-go/pkg/controller/command/vcwallet"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/didexchange"
	issuecredentialsvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/mediator"
	outofbandSvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofband"
	oobv2 "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/outofbandv2"
	presentproofSvc "github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/presentproof"
	mockoutofbandv2 "github.com/hyperledger/aries-framework-go/pkg/internal/gomocks/client/outofbandv2"
	"github.com/hyperledger/aries-framework-go/pkg/internal/ldtestutil"
	mockdidexchange "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/didexchange"
	mockissuecredential "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/issuecredential"
	mockmediator "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/mediator"
	mockoutofband "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/outofband"
	mockpresentproof "github.com/hyperledger/aries-framework-go/pkg/mock/didcomm/protocol/presentproof"
	mockprovider "github.com/hyperledger/aries-framework-go/pkg/mock/provider"
	mockstore "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
	"github.com/hyperledger/aries-framework-go/internal/testdata"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/tinkcrypto"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	mockvdr "github.com/hyperledger/aries-framework-go/pkg/mock/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"
	"github.com/hyperledger/aries-framework-go/pkg/wallet"
	"github.com/stretchr/testify/assert"

	"testing"
)

const (
	sampleDIDName = "sampleDIDName"
	sampleCredId  = "http://example.edu/credentials/1872"
	sampleDID     = "did:example:123"
	sampleURL     = "http://issuer:9082"
)

func TestNewDID(t *testing.T) {

	// Success: DID successfully created
	t.Run("test newDID method - success", func(t *testing.T) {
		purposeAuth := KeyTypePurpose{Purpose: "Authentication", KeyType: KeyTypeModel{Type: ed25519VerificationKey2018}}
		purposeAssertion := KeyTypePurpose{Purpose: "AssertionMethod", KeyType: KeyTypeModel{Type: bls12381G1Key2022, Attrs: []string{"2"}}}

		newDIDArgs := NewDIDArgs{Keys: []KeyTypePurpose{purposeAuth, purposeAssertion}, Name: sampleDIDName}

		var l bytes.Buffer
		reader, err := getReader(newDIDArgs)

		require.NotNil(t, reader)
		require.NoError(t, err)

		vcwalletCommand := vcwallet.New(newMockProvider(t), &vcwallet.Config{})
		require.NotNil(t, vcwalletCommand)
		require.NoError(t, err)

		vdrCommand, err := vdr.New(&mockprovider.Provider{
			StorageProviderValue: mockstore.NewMockStoreProvider(),
			VDRegistryValue:      &mockvdr.MockVDRegistry{},
		})
		require.NotNil(t, vdrCommand)
		require.NoError(t, err)

		command, err := New(vdrCommand, vcwalletCommand)
		require.NoError(t, err)

		err = command.NewDID(&l, reader)

		require.NoError(t, err)
		require.NotNil(t, command)

		var response NewDIDResult

		err = json.NewDecoder(&l).Decode(&response)
		require.NoError(t, err)

		var didDoc map[string]interface{}

		// Decode DIDDoc content
		err = json.Unmarshal(response.DIDDoc, &didDoc)
		require.NoError(t, err)

		prettyDidDoc, err := json.MarshalIndent(didDoc, "", "  ")
		require.NoError(t, err)

		fmt.Printf("DID Document: %s\n", string(prettyDidDoc))

		fmt.Println()
	})

	// Error case: empty keys
	t.Run("test newDID method - empty keys", func(t *testing.T) {
		newDIDArgs := NewDIDArgs{Keys: nil, Name: sampleDIDName}

		var l bytes.Buffer
		reader, err := getReader(newDIDArgs)
		require.NoError(t, err)

		vcwalletCommand := vcwallet.New(newMockProvider(t), &vcwallet.Config{})
		require.NotNil(t, vcwalletCommand)

		vdrCommand, err := vdr.New(&mockprovider.Provider{
			StorageProviderValue: mockstore.NewMockStoreProvider(),
			VDRegistryValue:      &mockvdr.MockVDRegistry{},
		})
		require.NoError(t, err)

		command, err := New(vdrCommand, vcwalletCommand)
		require.NoError(t, err)

		err = command.NewDID(&l, reader)

		require.Error(t, err)
		require.Contains(t, err.Error(), "keys is mandatory")
	})

	// Error case: invalid key type
	t.Run("test newDID method - invalid key type", func(t *testing.T) {
		purposeAuth := KeyTypePurpose{Purpose: "Authentication", KeyType: KeyTypeModel{Type: "invalidKeyType"}}
		newDIDArgs := NewDIDArgs{Keys: []KeyTypePurpose{purposeAuth}, Name: sampleDIDName}

		var l bytes.Buffer
		reader, err := getReader(newDIDArgs)
		require.NoError(t, err)

		vcwalletCommand := vcwallet.New(newMockProvider(t), &vcwallet.Config{})
		require.NotNil(t, vcwalletCommand)

		vdrCommand, err := vdr.New(&mockprovider.Provider{
			StorageProviderValue: mockstore.NewMockStoreProvider(),
			VDRegistryValue:      &mockvdr.MockVDRegistry{},
		})
		require.NoError(t, err)

		command, err := New(vdrCommand, vcwalletCommand)
		require.NoError(t, err)

		err = command.NewDID(&l, reader)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid key type")
	})
}


func Test_SignVerifyJWTContent(t *testing.T) {

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()

	// Initialize cryptographic provider
	tcrypto, err := tinkcrypto.New()
	require.NoError(t, err)

	mockctx.CryptoValue = tcrypto

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	token, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

	var b bytes.Buffer

	reqCreateKey, err := getReader(&vcwallet.CreateKeyPairRequest{
		WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: token},
		KeyType:    kms.ED25519Type,
	})
	require.NoError(t, err)

	err = command.vcwalletcommand.CreateKeyPair(&b, reqCreateKey)
	require.NoError(t, err)

	var keyPairResponse vcwallet.CreateKeyPairResponse
	require.NoError(t, json.NewDecoder(&b).Decode(&keyPairResponse))

	pubKey, err := base64.RawURLEncoding.DecodeString(keyPairResponse.PublicKey)
	require.NoError(t, err)

	_, didKID := fingerprint.CreateDIDKeyByCode(fingerprint.ED25519PubKeyMultiCodec, pubKey)
	parts := strings.Split(didKID, "#")
	currentDID := parts[0]
	currentKeyID := parts[1]

	b.Reset()
	lock1()

	command.currentDID = currentDID
	command.currentKeyPair = keyPairResponse
	command.currentKeyPair.KeyID = currentKeyID

	content := map[string]interface{}{
		"name":      "John Doe",
		"attrName":  "DID",
		"attrValue": command.currentDID,
	}

	// Success: JWT content successfully signed and verified
	t.Run("test Sign and Verify JWT content - success", func(t *testing.T) {
		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)

		signRequest, err := getReader(SignJWTContentArgs{
			Content: contentBytes,
		})
		require.NoError(t, err)

		var signResponse bytes.Buffer
		err = command.SignJWTContent(&signResponse, signRequest)
		require.NoError(t, err)

		// Verify JWT content
		var jwtSignResponse SignJWTContentResult
		err = json.Unmarshal(signResponse.Bytes(), &jwtSignResponse)
		require.NoError(t, err)

		verifyRequest, err := getReader(VerifyJWTContentArgs{
			JWT: jwtSignResponse.SignedJWTContent,
		})
		require.NoError(t, err)

		var verifyResponse bytes.Buffer
		err = command.VerifyJWTContent(&verifyResponse, verifyRequest)
		require.NoError(t, err)

		var jwtVerifyResponse vcwallet.VerifyJWTResponse
		err = json.Unmarshal(verifyResponse.Bytes(), &jwtVerifyResponse)
		require.NoError(t, err)

		if jwtVerifyResponse.Verified {
			fmt.Println("JWT content verified")
		}
	})

	// Error case: invalid wallet token
	t.Run("test SignJWTContent - invalid wallet token", func(t *testing.T) {
		invalidToken := "incorrect-token"
		reqCreateKey, err := getReader(&vcwallet.CreateKeyPairRequest{
			WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: invalidToken},
			KeyType:    kms.ED25519Type,
		})
		require.NoError(t, err)

		var b bytes.Buffer
		err = command.vcwalletcommand.CreateKeyPair(&b, reqCreateKey)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid auth token")
	})

	// Error case: malformed JWT
	t.Run("test VerifyJWTContent - malformed JWT", func(t *testing.T) {
		malformedJWT := "not-a-jwt-token"
		verifyRequest, err := getReader(VerifyJWTContentArgs{
			JWT: malformedJWT,
		})
		require.NoError(t, err)

		var verifyResponse bytes.Buffer
		err = command.VerifyJWTContent(&verifyResponse, verifyRequest)

		var jwtVerifyResponse vcwallet.VerifyJWTResponse
		err = json.Unmarshal(verifyResponse.Bytes(), &jwtVerifyResponse)

		require.False(t, jwtVerifyResponse.Verified)
	})

	// Error case: tampered JWT
	t.Run("test VerifyJWTContent - tampered JWT", func(t *testing.T) {
		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)

		var signResponse bytes.Buffer
		signRequest, err := getReader(SignJWTContentArgs{
			Content: contentBytes,
		})
		require.NoError(t, err)

		err = command.SignJWTContent(&signResponse, signRequest)
		require.NoError(t, err)

		var jwtSignResponse SignJWTContentResult
		err = json.Unmarshal(signResponse.Bytes(), &jwtSignResponse)
		require.NoError(t, err)

		// Manipulate signed JWT
		tamperedJWT := jwtSignResponse.SignedJWTContent + "tamper"

		verifyRequest, err := getReader(VerifyJWTContentArgs{
			JWT: tamperedJWT,
		})
		require.NoError(t, err)

		var verifyResponse bytes.Buffer
		err = command.VerifyJWTContent(&verifyResponse, verifyRequest)

		var jwtVerifyResponse vcwallet.VerifyJWTResponse
		err = json.Unmarshal(verifyResponse.Bytes(), &jwtVerifyResponse)
		require.NoError(t, err)

		require.False(t, jwtVerifyResponse.Verified)
		require.Contains(t, jwtVerifyResponse.Error, "jwt verification failed")
	})

}


func Test_SignVerifyContract(t *testing.T) {

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()
		
	tcrypto, err := tinkcrypto.New()
	require.NoError(t, err)
	mockctx.CryptoValue = tcrypto

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	// Unlock wallet and create key
	token, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

	var b bytes.Buffer
	createKeyReq, err := getReader(&vcwallet.CreateKeyPairRequest{
		WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: token},
		KeyType:    kms.ED25519Type,
	})
	require.NoError(t, err)

	err = vcwalletCommand.CreateKeyPair(&b, createKeyReq)
	require.NoError(t, err)

	var keyResp vcwallet.CreateKeyPairResponse
	require.NoError(t, json.NewDecoder(&b).Decode(&keyResp))

	pubKey, err := base64.RawURLEncoding.DecodeString(keyResp.PublicKey)
	require.NoError(t, err)

	_, didKID := fingerprint.CreateDIDKeyByCode(fingerprint.ED25519PubKeyMultiCodec, pubKey)
	parts := strings.Split(didKID, "#")
	command.currentDID = parts[0]
	command.currentKeyPair = keyResp
	command.currentKeyPair.KeyID = parts[1]
		
	lock1()
	
	// Contract JSON
	contractJSON := `
	{
	  "apiVersion": "reservation.fluidos.eu/v1alpha1",
	  "kind": "Contract",
	  "metadata": {
	    "creationTimestamp": "2024-05-29T09:24:45Z",
	    "generation": 1,
	    "name": "contract-fluidos.eu-k8s-fluidos-ccbf29bb-b4a7",
	    "namespace": "fluidos",
	    "resourceVersion": "1578",
	    "uid": "22e1136c-98b9-4de3-a5e9-90af7d5aa021"
	  },
	  "spec": {
	    "buyer": {
	      "domain": "fluidos.eu",
	      "ip": "172.18.0.4:30000",
	      "nodeID": "lwuyt2dnxe"
	    },
	    "buyerClusterID": "c47fb461-5bff-4b64-a381-1058fa263235",
	    "expirationTime": "2025-05-29T09:24:45Z",
	    "flavour": {
	      "metadata": {
		"name": "fluidos.eu-k8s-fluidos-ccbf29bb",
		"namespace": "fluidos"
	      },
	      "spec": {
		"characteristics": {
		  "architecture": "amd64",
		  "cpu": "7985105637n",
		  "ephemeral-storage": "0",
		  "gpu": "0",
		  "memory": "32386980Ki",
		  "persistent-storage": "0",
		  "pods": "110"
		},
		"optionalFields": {
		  "availability": true,
		  "workerID": "fluidos-provider-1-worker2"
		},
		"owner": {
		  "domain": "fluidos.eu",
		  "ip": "172.18.0.2:30001",
		  "nodeID": "jgmewzljr9"
		},
		"policy": {
		  "aggregatable": {
		    "maxCount": 0,
		    "minCount": 0
		  },
		  "partitionable": {
		    "cpuMin": "0",
		    "cpuStep": "1",
		    "memoryMin": "0",
		    "memoryStep": "100Mi",
		    "podsMin": "0",
		    "podsStep": "0"
		  }
		},
		"price": {
		  "amount": "",
		  "currency": "",
		  "period": ""
		},
		"providerID": "jgmewzljr9",
		"type": "k8s-fluidos"
	      },
	      "status": {
		"creationTime": "",
		"expirationTime": "",
		"lastUpdateTime": ""
	      }
	    },
	    "partition": {
	      "architecture": "",
	      "cpu": "1",
	      "ephemeral-storage": "0",
	      "gpu": "0",
	      "memory": "1Gi",
	      "pods": "50",
	      "storage": "0"
	    },
	    "seller": {
	      "domain": "fluidos.eu",
	      "ip": "172.18.0.2:30001",
	      "nodeID": "jgmewzljr9"
	    },
	    "sellerCredentials": {
	      "clusterID": "08fcbfd5-a76e-444d-a182-de3b25398e2a",
	      "clusterName": "fluidos-provider-1",
	      "endpoint": "https://172.18.0.2:32197",
	      "token": "secure-token-value"
	    },
	    "transactionID": "b27a019255fa7748c004fb1116ae7281-1716974685039544567"
	  }
	}
`

	var contractRaw json.RawMessage
	err = json.Unmarshal([]byte(contractJSON), &contractRaw)
	require.NoError(t, err)

	// Success: contract successfully signed and signature verified
	t.Run("test Sign Contract and Verify Signature - success", func(t *testing.T) {
		signReqReader, err := getReader(SignContractArgs{
			Contract: contractRaw,
		})
		require.NoError(t, err)

		var signResp bytes.Buffer
		err = command.SignContract(&signResp, signReqReader)
		require.NoError(t, err)

		var result SignContractResult
		err = json.Unmarshal(signResp.Bytes(), &result)
		require.NoError(t, err)

		require.NotEmpty(t, result.SignedContract)
		
		// Verify Contract signature
		verifyReqReader, err := getReader(VerifyContractSignatureArgs{
			Contract: result.SignedContract,
		})
		require.NoError(t, err)

		var verifyResp bytes.Buffer
		err = command.VerifyContractSignature(&verifyResp, verifyReqReader)
		require.NoError(t, err)

		var verifyResult VerifyContractSignatureResult
		err = json.Unmarshal(verifyResp.Bytes(), &verifyResult)
		require.NoError(t, err)

		require.True(t, verifyResult.VerifiedChain, "expected signature chain to be verified")
		require.NotEmpty(t, verifyResult.Signatures, "expected at least one signature")
		for _, sig := range verifyResult.Signatures {
			require.True(t, sig.Verified, "signature should be verified")
		}
	})
	
	// Error case: empty contract
	t.Run("test Sign Contract - empty contract error", func(t *testing.T) {
	    signRequest, err := getReader(SignContractArgs{
		Contract:     json.RawMessage(nil),
		ContractJWT:  "",
	    })
	    require.NoError(t, err)

	    var signResponse bytes.Buffer
	    err = command.SignContract(&signResponse, signRequest)

	    require.Error(t, err)
	    require.Contains(t, err.Error(), "you have to provide a contract in json format or jwt format")
	})

	// Error case: both Contract and ContractJWT provided
	t.Run("test Sign Contract - both Contract and ContractJWT provided", func(t *testing.T) {
	    dummyJWT := "eyJhbGciOiJFZERTQSJ9.eyJmb28iOiJiYXIifQ.XYZ" // false value

	    contentBytes, err := json.Marshal(contractRaw)
	    require.NoError(t, err)

	    signRequest, err := getReader(SignContractArgs{
		Contract:     contentBytes,
		ContractJWT:  dummyJWT,
	    })
	    require.NoError(t, err)

	    var signResponse bytes.Buffer
	    err = command.SignContract(&signResponse, signRequest)

	    require.Error(t, err)
	    require.Contains(t, err.Error(), "Contract and ContractJWT are both provided")
	})
	
	// Error case: empty contract
	t.Run("test Verify Contract Signature - empty contract", func(t *testing.T) {
		verifyReqReader, err := getReader(VerifyContractSignatureArgs{
			Contract: "",
		})
		require.NoError(t, err)

		var verifyResp bytes.Buffer
		err = command.VerifyContractSignature(&verifyResp, verifyReqReader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Contract is mandatory")
	})
	
	// Error case: tampered contract JWT
	t.Run("test Verify Contract Signature - tampered JWT", func(t *testing.T) {
		var signResp bytes.Buffer
		signReqReader, err := getReader(SignContractArgs{
			Contract: contractRaw,
		})
		require.NoError(t, err)

		err = command.SignContract(&signResp, signReqReader)
		require.NoError(t, err)

		var result SignContractResult
		err = json.Unmarshal(signResp.Bytes(), &result)
		require.NoError(t, err)

		// Tamper with the signed JWT
		tamperedJWT := result.SignedContract + "tamper"

		verifyReqReader, err := getReader(VerifyContractSignatureArgs{
			Contract: tamperedJWT,
		})
		require.NoError(t, err)

		var verifyResp bytes.Buffer
		err = command.VerifyContractSignature(&verifyResp, verifyReqReader)
		require.NoError(t, err)

		var verifyResult VerifyContractSignatureResult
		err = json.Unmarshal(verifyResp.Bytes(), &verifyResult)
		require.NoError(t, err)

		require.False(t, verifyResult.VerifiedChain, "expected signature chain to be unverified")
	})

}


func TestDoDeviceEnrolment(t *testing.T) {

	// mockServer to capture the http request (inside DoDeviceEnrolment) to AcceptEnrolment
	mockServer := setupMockEnrolmentServer(t)
	defer mockServer.Close()

	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()

	// Initialize cryptographic provider
	tcrypto, err := tinkcrypto.New()
	require.NoError(t, err)

	mockctx.CryptoValue = tcrypto

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	token, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

	var b bytes.Buffer

	reqCreateKey, err := getReader(&vcwallet.CreateKeyPairRequest{
		WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: token},
		KeyType:    kms.ED25519Type,
	})
	require.NoError(t, err)

	err = command.vcwalletcommand.CreateKeyPair(&b, reqCreateKey)
	require.NoError(t, err)

	var keyPairResponse vcwallet.CreateKeyPairResponse
	require.NoError(t, json.NewDecoder(&b).Decode(&keyPairResponse))

	pubKey, err := base64.RawURLEncoding.DecodeString(keyPairResponse.PublicKey)
	require.NoError(t, err)

	_, didKID := fingerprint.CreateDIDKeyByCode(fingerprint.ED25519PubKeyMultiCodec, pubKey)
	parts := strings.Split(didKID, "#")
	currentDID := parts[0]
	currentKeyID := parts[1]

	b.Reset()
	lock1()

	command.currentDID = currentDID
	command.currentKeyPair = keyPairResponse
	command.currentKeyPair.KeyID = currentKeyID

	// DoDeviceEnrolment sample args
	enrolmentArgs := DoDeviceEnrolmentArgs{
		Url: sampleURL,
		IdProofs: []IdProof{
			{AttrName: "holderName", AttrValue: "FluidosNode"},
			{AttrName: "fluidosRole", AttrValue: "Customer"},
			{AttrName: "deviceType", AttrValue: "Server"},
			{AttrName: "orgIdentifier", AttrValue: "FLUIDOS_id_23241231412"},
			{AttrName: "physicalAddress", AttrValue: "50:80:61:82:ab:c9"},
		},
	}
	enrolmentArgs.Url = mockServer.URL
	
	var requestBody bytes.Buffer
	err = json.NewEncoder(&requestBody).Encode(enrolmentArgs)
	require.NoError(t, err)
	
	var l bytes.Buffer

	// Success: DoDeviceEnrolment
	t.Run("test DoDeviceEnrolment method - success", func(t *testing.T) {
		// DoDeviceEnrolment method
		err = command.DoDeviceEnrolment(&l, &requestBody)
		require.NoError(t, err)

		var response DoDeviceEnrolmentResult
		err = json.NewDecoder(&l).Decode(&response)
		require.NoError(t, err)
		require.NotNil(t, response)

		fmt.Println("Credential storage ID:", response.CredStorageId)

		var cred map[string]interface{}

		err = json.Unmarshal(response.Credential, &cred)
		require.NoError(t, err)

		prettyCred, err := json.MarshalIndent(cred, "", "  ")
		require.NoError(t, err)

		fmt.Printf("Formatted Credential: %s\n", string(prettyCred))

		fmt.Println()
	})
	
	// Error case: malformed JSON request
	t.Run("test DoDeviceEnrolment - malformed JSON request", func(t *testing.T) {
		malformedJSON := strings.NewReader(`{"Url": "sampleURL", "IdProofs": [{{}]}`)

		err = command.DoDeviceEnrolment(&l, malformedJSON)

		require.Error(t, err)
		require.Contains(t, err.Error(), "request decode")
	})

	// Error case: missing URL
	t.Run("test DoDeviceEnrolment - missing URL", func(t *testing.T) {
		enrolmentArgs := DoDeviceEnrolmentArgs{
			IdProofs: []IdProof{
				{AttrName: "holderName", AttrValue: "FluidosNode"},
			},
		}

		var requestBody bytes.Buffer
		err := json.NewEncoder(&requestBody).Encode(enrolmentArgs)
		require.NoError(t, err)

		var l bytes.Buffer
		err = command.DoDeviceEnrolment(&l, &requestBody)

		require.Error(t, err)
		require.Contains(t, err.Error(), "url is mandatory")
	})

	// Error: missing IdProofs
	t.Run("test DoDeviceEnrolment - missing IdProofs", func(t *testing.T) {
		enrolmentArgs := DoDeviceEnrolmentArgs{Url: sampleURL}

		var requestBody bytes.Buffer
		err := json.NewEncoder(&requestBody).Encode(enrolmentArgs)
		require.NoError(t, err)

		var l bytes.Buffer
		err = command.DoDeviceEnrolment(&l, &requestBody)

		require.Error(t, err)
		require.Contains(t, err.Error(), "idProofs is mandatory")
	})
	
}


func TestGetVCredential(t *testing.T) {

	vcwalletCommand := vcwallet.New(newMockProvider(t), &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(&mockprovider.Provider{
		StorageProviderValue: mockstore.NewMockStoreProvider(),
		VDRegistryValue:      &mockvdr.MockVDRegistry{},
	})
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// New command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	// Success: VCredential successfully obtained
	t.Run("test GetVCredential method - success", func(t *testing.T) {
		fmt.Printf("ID of the sought credential: %s\n", sampleCredId)

		// Valid argument for GetVCredential
		getVCredentialArgs := GetVCredentialArgs{CredId: sampleCredId}
		var l bytes.Buffer
		reader, err := getReader(getVCredentialArgs)
		require.NotNil(t, reader)
		require.NoError(t, err)

		// Sample VC
		vc2 := map[string]interface{}{
			"@context": []string{"https://www.w3.org/2018/credentials/v1"},
			"id":       "http://example.edu/credentials/1872",
			"type":     []string{"VerifiableCredential"},
			"issuer":   map[string]interface{}{"id": "did:example:123"},
			"credentialSubject": map[string]interface{}{
				"id":   "did:example:456",
				"name": "John Doe",
			},
		}

		token1, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

		err = command.AddCredentialToWallet(command.walletuid, token1, wallet.Credential, vc2, "")
		require.NoError(t, err)

		lock1()

		err = command.GetVCredential(&l, reader)
		require.NoError(t, err)

		require.NotNil(t, command)

		var response GetVCredentialResult
		err = json.NewDecoder(&l).Decode(&response)
		require.NoError(t, err)

		require.NotNil(t, response)
		fmt.Println(response)

		var didDoc map[string]interface{}

		// Decode DIDDoc content
		err = json.Unmarshal(response.Credential, &didDoc)
		require.NoError(t, err)

		prettyDidDoc, err := json.MarshalIndent(didDoc, "", "  ")
		require.NoError(t, err)

		fmt.Printf("Credential: %s\n", string(prettyDidDoc))

		fmt.Println()
	})

	// Error case: malformed JSON request
	t.Run("test GetVCredential - malformed JSON request", func(t *testing.T) {
		malformedJSON := strings.NewReader("{CredId:}")

		var l bytes.Buffer
		err := command.GetVCredential(&l, malformedJSON)

		require.Error(t, err)
		require.Contains(t, err.Error(), "request decode")
	})

	// Error case: credential not found
	t.Run("test GetVCredential - credential not found", func(t *testing.T) {
		missingCredId := "http://example.edu/credentials/9999"
		getVCredentialArgs := GetVCredentialArgs{CredId: missingCredId}
		reader, err := getReader(getVCredentialArgs)
		require.NoError(t, err)

		var l bytes.Buffer
		err = command.GetVCredential(&l, reader)

		require.Error(t, err)
		require.Contains(t, err.Error(), "data not found")
	})

	// Error case: invalid wallet credentials
	t.Run("test GetVCredential - invalid wallet credentials", func(t *testing.T) {
		command.walletuid = "invalidUserID"
		command.walletpass = "invalidPassphrase"

		var l bytes.Buffer
		getVCredentialArgs := GetVCredentialArgs{CredId: sampleCredId}
		reader, err := getReader(getVCredentialArgs)
		require.NotNil(t, reader)
		require.NoError(t, err)
		err = command.GetVCredential(&l, reader)

		require.Error(t, err)
		require.Contains(t, err.Error(), "profile does not exist")
	})

}

func TestGenerateVP(t *testing.T) {

	// Simulated provider
	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	token1, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

	// Add sample credential to wallet
	var sampleNewUDCVc map[string]interface{}
	err = json.Unmarshal(testdata.SampleUDCVC, &sampleNewUDCVc)
	require.NoError(t, err)

	sampleNewUDCVc["id"] = "http://example.edu/credentials/18722"

	// Add credential to wallet
	err = command.AddCredentialToWallet(command.walletuid, token1, wallet.Credential, sampleNewUDCVc, "")
	require.NoError(t, err)

	// sampleUDCVCWithProofBBS includes a proof object, which is essential for credential
	// verification. This provides information that allows systems to validate that the
	// credential has not been tampered with and that it comes from a legitimate source.
	var sampleNewUDCVProofBBS map[string]interface{}
	err = json.Unmarshal(testdata.SampleUDCVCWithProofBBS, &sampleNewUDCVProofBBS)
	require.NoError(t, err)

	err = command.AddCredentialToWallet(command.walletuid, token1, wallet.Credential, sampleNewUDCVProofBBS, "")
	require.NoError(t, err)

	lock1()

	var queryByFrame QueryByFrame
	err = json.Unmarshal(testdata.SampleWalletQueryByFrame, &queryByFrame)
	require.NoError(t, err)

	// Success: VP successfully generated
	t.Run("test Generate VP - success", func(t *testing.T) {
		request := &GenerateVPArgs{
			CredId:       "http://example.edu/credentials/18722",
			QueryByFrame: queryByFrame,
		}
		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		var b bytes.Buffer
		cmdErr := command.GenerateVP(&b, bytes.NewReader(reqBody))
		require.NoError(t, cmdErr)

		var response GenerateVPResultCustom
		require.NoError(t, json.NewDecoder(&b).Decode(&response))
		require.NotEmpty(t, response)
		require.NotEmpty(t, response.Results)

		t.Log("Prueba de generación de VP exitosa con resultado:", response.Results)
	})

	// Error case: non-existent credential ID
	t.Run("test GenerateVP - non-existent credential ID", func(t *testing.T) {
		request := &GenerateVPArgs{
			CredId:       "http://example.edu/credentials/unknownID",
			QueryByFrame: queryByFrame,
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		var b bytes.Buffer
		cmdErr := command.GenerateVP(&b, bytes.NewReader(reqBody))

		require.Error(t, cmdErr)
		require.Contains(t, cmdErr.Error(), "data not found")
	})

	// Error case: empty QueryByFrame
	t.Run("test GenerateVP - empty QueryByFrame", func(t *testing.T) {
		request := &GenerateVPArgs{
			CredId:       "http://example.edu/credentials/18722",
			QueryByFrame: QueryByFrame{},
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		var b bytes.Buffer
		cmdErr := command.GenerateVP(&b, bytes.NewReader(reqBody))

		require.Error(t, cmdErr)
		require.Contains(t, cmdErr.Error(), "query response not working")
	})

	// Error case: missing CredId in request
	t.Run("test GenerateVP - missing CredId in request", func(t *testing.T) {
		request := &GenerateVPArgs{
			CredId:       "",
			QueryByFrame: queryByFrame,
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		var b bytes.Buffer
		cmdErr := command.GenerateVP(&b, bytes.NewReader(reqBody))

		require.Error(t, cmdErr)
		require.Contains(t, cmdErr.Error(), "data not found")
	})

}

func TestVerifyCredential(t *testing.T) {

	// Simulated provider
	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)

	token1, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)

	// Add sample credential to wallet
	var sampleNewUDCVc map[string]interface{}
	err = json.Unmarshal(testdata.SampleUDCVC, &sampleNewUDCVc)
	require.NoError(t, err)

	sampleNewUDCVc["id"] = "http://example.edu/credentials/18722"

	// Add credential to wallet
	err = command.AddCredentialToWallet(command.walletuid, token1, wallet.Credential, sampleNewUDCVc, "")
	require.NoError(t, err)

	// sampleUDCVCWithProofBBS includes a proof object, which is essential for credential
	// verification. This provides information that allows systems to validate that the
	// credential has not been tampered with and that it comes from a legitimate source.
	var sampleNewUDCVProofBBS map[string]interface{}
	err = json.Unmarshal(testdata.SampleUDCVCWithProofBBS, &sampleNewUDCVProofBBS)
	require.NoError(t, err)

	err = command.AddCredentialToWallet(command.walletuid, token1, wallet.Credential, sampleNewUDCVProofBBS, "")
	require.NoError(t, err)

	lock1()

	// Success: VP successfully issued and verified
	t.Run("test VerifyCredential - success", func(t *testing.T) {
		// Part I: Issue VP
		var queryByFrame QueryByFrame
		err = json.Unmarshal(testdata.SampleWalletQueryByFrame, &queryByFrame)
		require.NoError(t, err)

		request := &GenerateVPArgs{
			CredId:       "http://example.edu/credentials/18722",
			QueryByFrame: queryByFrame,
		}

		reqBody, err := json.Marshal(request)
		require.NoError(t, err)

		var b bytes.Buffer
		cmdErr := command.GenerateVP(&b, bytes.NewReader(reqBody))
		require.NoError(t, cmdErr)

		var res GenerateVPResultCustom
		require.NoError(t, json.NewDecoder(&b).Decode(&res))
		require.NotEmpty(t, res)
		require.NotEmpty(t, res.Results)

		var l bytes.Buffer
		fmt.Println(string(*res.Results[0]))
		reader, err := getReader(VerifyCredentialArgs{CredentialString: string(*res.Results[0])})
		require.NotNil(t, reader)
		require.NoError(t, err)

		// Part II: call VerifyCredential method
		err = command.VerifyCredential(&l, reader)
		require.NoError(t, err)

		// Validate response
		var response VerifyCredentialResult
		err = json.NewDecoder(&l).Decode(&response)
		require.NoError(t, err)

		require.True(t, response.Result)
		fmt.Printf("Credential verification result: %v\n", response.Result)
	})

	// Error case: malformed JSON request
	t.Run("test VerifyCredential - malformed JSON request", func(t *testing.T) {
		malformedJSON := strings.NewReader("{CredentialString:}")

		var l bytes.Buffer
		err := command.VerifyCredential(&l, malformedJSON)

		require.Error(t, err)
		require.Contains(t, err.Error(), "request decode")
	})

	// Error case: invalid credential format
	t.Run("test VerifyCredential - invalid credential format", func(t *testing.T) {
		invalidCredential := map[string]interface{}{
			"id":       "http://example.edu/credentials/invalid",
			"@context": "https://www.w3.org/2018/credentials/v1",
		}

		var l bytes.Buffer
		credBytes, err := json.Marshal(invalidCredential)
		require.NoError(t, err)

		reader := bytes.NewReader(credBytes)
		err = command.VerifyCredential(&l, reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get Verify Request reader")
	})

}


func TestAcceptEnrolment(t *testing.T) {

	// Simulated provider
	mockctx := newMockProvider(t)
	mockctx.VDRegistryValue = getMockDIDKeyVDR()

	// Initialize cryptographic provider
	tcrypto, err := tinkcrypto.New()
	require.NoError(t, err)

	mockctx.CryptoValue = tcrypto

	vcwalletCommand := vcwallet.New(mockctx, &vcwallet.Config{})
	require.NotNil(t, vcwalletCommand)

	vdrCommand, err := vdr.New(mockctx)
	require.NotNil(t, vdrCommand)
	require.NoError(t, err)

	// Command instance
	command, err := New(vdrCommand, vcwalletCommand)
	require.NoError(t, err)
	
	token, lock1 := command.unlockWallet(t, command.walletuid, command.walletpass)
	
	var b bytes.Buffer

	reqCreateKey, err := getReader(&vcwallet.CreateKeyPairRequest{
		WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: token},
		KeyType:    kms.ED25519Type,
	})
	require.NoError(t, err)

	err = command.vcwalletcommand.CreateKeyPair(&b, reqCreateKey)
	require.NoError(t, err)

	var keyPairResponse vcwallet.CreateKeyPairResponse
	require.NoError(t, json.NewDecoder(&b).Decode(&keyPairResponse))

	pubKey, err := base64.RawURLEncoding.DecodeString(keyPairResponse.PublicKey)
	require.NoError(t, err)

	_, didKID := fingerprint.CreateDIDKeyByCode(fingerprint.ED25519PubKeyMultiCodec, pubKey)
	parts := strings.Split(didKID, "#")
	currentDID := parts[0]
	currentKeyID := parts[1]

	b.Reset()

	command.currentDID = currentDID
	command.currentKeyPair = keyPairResponse
	command.currentKeyPair.KeyID = currentKeyID
	
	// Sign JWT (get ProofData)
	request := vcwallet.SignJWTRequest{
		WalletAuth: vcwallet.WalletAuth{UserID: command.walletuid, Auth: token},
		Headers:    nil,
		Claims: map[string]interface{}{
			"attrName":  "DID",
			"attrValue": command.currentDID,
		},
		KID: command.currentDID + "#" + command.currentKeyPair.KeyID,
	}

	reqData, err := json.Marshal(request)
	require.NoError(t, err)
	req := bytes.NewReader(reqData)
	// Capture the output
	var signBuf bytes.Buffer

	// Sign the JWT
	err = vcwalletCommand.SignJWT(&signBuf, req)
	require.NoError(t, err)

	var jwtResponse vcwallet.SignJWTResponse

	err = json.Unmarshal(signBuf.Bytes(), &jwtResponse)
	require.NoError(t, err)

	signedJWT := jwtResponse.JWT
	
	// Add sample DID to VDR (mock VDR content)
	command.AddDIDToVDR(t, currentDID, currentDID+"#"+command.currentKeyPair.KeyID)
	
	// Generate random keys for wallet
	privKeyBase58, pubKeyBase58, err := GenerateRandomBase58Keys()
	require.NoError(t, err)
	
	// Add correct content to wallet (type "Bls12381G1Key2020")
	var sampleWalletKey map[string]interface{}
	err = json.Unmarshal(testdata.SampleWalletContentKeyBase58, &sampleWalletKey)
	require.NoError(t, err)

	sampleWalletKey["privateKeyBase58"] = privKeyBase58
	sampleWalletKey["publicKeyBase58"] = pubKeyBase58
	sampleWalletKey["type"] = "Bls12381G1Key2020"
	sampleWalletKey["controller"] = currentDID
	sampleWalletKey["id"] = command.currentDID + "#" + command.currentKeyPair.KeyID
	
	err = command.AddCredentialToWallet(command.walletuid, token, wallet.Key, sampleWalletKey, "")
	require.NoError(t, err)
	
	lock1()

	// Success: enrolment process
	t.Run("test AcceptEnrolment method - success", func(t *testing.T) {
		// idProofs definition
		idProofs := []IdProof{
			{AttrName: "holderName", AttrValue: "FluidosNode"},
			{AttrName: "fluidosRole", AttrValue: "Customer"},
			{AttrName: "deviceType", AttrValue: "Server"},
			{AttrName: "orgIdentifier", AttrValue: "FLUIDOS_id_23241231412"},
			{
				AttrName:  "DID",
				AttrValue: command.currentDID,
				ProofData: signedJWT,
			},
		}

		var l bytes.Buffer
		acceptEnrolmentArgs := AcceptEnrolmentArgs{IdProofs: idProofs}

		reader, err := getReader(acceptEnrolmentArgs)
		require.NotNil(t, reader)
		require.NoError(t, err)

		err = command.AcceptEnrolment(&l, reader)
		require.NoError(t, err)

		var response AcceptEnrolmentResult
		err = json.NewDecoder(&l).Decode(&response)
		fmt.Println("Raw response output:", l.String())
		require.NoError(t, err)

		require.NotNil(t, response)

		var didDoc map[string]interface{}
		err = json.Unmarshal(response.Credential, &didDoc)
		require.NoError(t, err)

		prettyDidDoc, err := json.MarshalIndent(didDoc, "", "  ")
		require.NoError(t, err)

		fmt.Printf("Credential: %s\n", string(prettyDidDoc))

		fmt.Println("AcceptEnrolment executed successfully.")
	})

	// Error case: empty idProofs
	t.Run("test AcceptEnrolment - empty idProofs", func(t *testing.T) {
		var l bytes.Buffer

		// empty idProofs
		acceptEnrolmentArgs := AcceptEnrolmentArgs{IdProofs: []IdProof{}}

		reader, err := getReader(acceptEnrolmentArgs)
		require.NoError(t, err)

		err = command.AcceptEnrolment(&l, reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "idProofs is mandatory")
	})

	// Error case: DID not found
	t.Run("test AcceptEnrolment - DID not found", func(t *testing.T) {
		var l bytes.Buffer

		// DID not found
		idProofs := []IdProof{
			{AttrName: "holderName", AttrValue: "FluidosNode"},
			{AttrName: "DID", AttrValue: "did:local:abc"},
			{
				AttrName:  "DID",
				AttrValue: command.currentDID,
				ProofData: signedJWT,
			},
		}

		acceptEnrolmentArgs := AcceptEnrolmentArgs{IdProofs: idProofs}

		reader, err := getReader(acceptEnrolmentArgs)
		require.NoError(t, err)

		command.currentDID = "did:key:notfound"
		err = command.AcceptEnrolment(&l, reader)
		require.Nil(t, err)
	})

}


/**
func TestCallOpteeGenerateKey(t *testing.T) {
	t.Run("test CallOpteeGenerateKey - success", func(t *testing.T) {
	
		cmd := Command{}

		// Sample key_id
		inputJSON := `{"key_id": "test_key_1234"}`
		inputReader := bytes.NewBufferString(inputJSON)

		var responseBuffer bytes.Buffer
		err := cmd.CallOpteeGenerateKey(&responseBuffer, inputReader)

		require.NoError(t, err)

		require.NotEmpty(t, responseBuffer.String())
		
		fmt.Println("CallOpteeGenerateKey response:", responseBuffer.String())
	})
}
**/

func readDIDtesting(t *testing.T) {

}

func newMockProvider(t *testing.T) *mockprovider.Provider {
	t.Helper()

	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	serviceMap := map[string]interface{}{
		presentproofSvc.Name:    &mockpresentproof.MockPresentProofSvc{},
		outofbandSvc.Name:       &mockoutofband.MockOobService{},
		didexchange.DIDExchange: &mockdidexchange.MockDIDExchangeSvc{},
		mediator.Coordination:   &mockmediator.MockMediatorSvc{},
		issuecredentialsvc.Name: &mockissuecredential.MockIssueCredentialSvc{},
		oobv2.Name:              &mockoutofbandv2.MockOobService{},
	}

	return &mockprovider.Provider{
		StorageProviderValue:              mockstore.NewMockStoreProvider(),
		ProtocolStateStorageProviderValue: mockstore.NewMockStoreProvider(),
		DocumentLoaderValue:               loader,
		ServiceMap:                        serviceMap,
	}
}

func (o *Command) unlockWallet(t *testing.T, sampleUser string, localKMS string) (string, func()) {
	var b bytes.Buffer

	openReader, err := getReader(&vcwallet.UnlockWalletRequest{
		UserID:             sampleUser,
		LocalKMSPassphrase: localKMS,
	})
	require.NoError(t, err)

	cmdErr := o.vcwalletcommand.Open(&b, openReader)
	require.NoError(t, cmdErr)

	lockReader, err := getReader(&vcwallet.LockWalletRequest{
		UserID: sampleUser,
	})
	require.NoError(t, err)

	return getUnlockToken(b), func() {
		cmdErr = o.vcwalletcommand.Close(&b, lockReader)
		if cmdErr != nil {
			t.Log(t, cmdErr)
		}
	}
}

func (o *Command) AddCredentialToWallet(userID string, walletAuth string, contentType wallet.ContentType, content interface{}, collectionID string) error {

	rawContent, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("error converting content to json.RawMessage: %w", err)
	}

	addReader, err := getReader(&vcwallet.AddContentRequest{
		WalletAuth: vcwallet.WalletAuth{
			UserID: userID,
			Auth:   walletAuth,
		},
		ContentType:  contentType,
		Content:      rawContent,
		CollectionID: collectionID,
	})
	if err != nil {
		return fmt.Errorf("error preparing request: %w", err)
	}

	var addResponse bytes.Buffer
	err = o.vcwalletcommand.Add(&addResponse, addReader)
	if err != nil {
		return fmt.Errorf("error adding content to wallet: %w", err)
	}

	return nil
}

// GenerateRandomBase58Keys generates public and private random keys in Base58 format
func GenerateRandomBase58Keys() (string, string, error) {
	privKey := make([]byte, 32)
	_, err := rand.Read(privKey)
	if err != nil {
		return "", "", err
	}

	pubKey := make([]byte, 48)
	_, err = rand.Read(pubKey)
	if err != nil {
		return "", "", err
	}

	privKeyBase58 := base58.Encode(privKey)
	pubKeyBase58 := base58.Encode(pubKey)

	return privKeyBase58, pubKeyBase58, nil
}

func (o *Command) AddDIDToVDR(t *testing.T, did, keyID string) error {
	pubKey := strings.TrimPrefix(did, "did:key:")

	doc := map[string]interface{}{
		"@context": []string{
			"https://w3id.org/did/v1",
		},
		"id": did,
		"verificationMethod": []map[string]interface{}{
			{
				"controller":      did,
				"id":              keyID,
				"publicKeyBase58": pubKey,
				"type":            "Ed25519VerificationKey2018",
			},
		},
		"service": []map[string]interface{}{
			{
				"id":              did + "#did-communication",
				"priority":        0,
				"recipientKeys":   []string{keyID},
				"serviceEndpoint": "https://agent.example.com/",
				"type":            "did-communication",
			},
		},
		"created": "2024-11-14T18:07:41.062167792+01:00",
	}

	rawContent, err := json.Marshal(doc)
	require.NoError(t, err)

	createDIDReq := &vdr.DIDArgs{
		Document: vdr.Document{DID: rawContent},
		Name:     sampleDIDName,
	}

	reqBytes, err := json.Marshal(createDIDReq)
	require.NoError(t, err)

	var getRW bytes.Buffer
	cmdErr := o.vdrcommand.SaveDID(&getRW, bytes.NewBuffer(reqBytes))
	require.NoError(t, cmdErr)

	return nil
}

func getMockDIDKeyVDR() *mockvdr.MockVDRegistry {
	return &mockvdr.MockVDRegistry{
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
		},
	}
}

func setupMockEnrolmentServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	
		assert.Equal(t, "/fluidos/idm/acceptEnrolment", req.URL.Path)
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var request AcceptEnrolmentArgs
		err := json.NewDecoder(req.Body).Decode(&request)
		require.NoError(t, err)

		// Verify request.IdProofs have expected data
		expectedIdProofs := []IdProof{
			{AttrName: "holderName", AttrValue: "FluidosNode"},
			{AttrName: "fluidosRole", AttrValue: "Customer"},
			{AttrName: "deviceType", AttrValue: "Server"},
			{AttrName: "orgIdentifier", AttrValue: "FLUIDOS_id_23241231412"},
			{AttrName: "physicalAddress", AttrValue: "50:80:61:82:ab:c9"},
		}
		
		require.GreaterOrEqual(t, len(request.IdProofs), len(expectedIdProofs), "insufficient IdProofs returned")
		actualProofsSubset := request.IdProofs[:len(expectedIdProofs)]
		assert.ElementsMatch(t, expectedIdProofs, actualProofsSubset)

		issuanceDate := time.Now().Format(time.RFC3339Nano)
		expirationDate := time.Now().AddDate(0, 0, 1).Format(time.RFC3339Nano)

		// Simulate the enrollment response
		response := map[string]interface{}{
			"credential": map[string]interface{}{
				"@context": []string{
					"https://www.w3.org/2018/credentials/v1",
					"https://www.w3.org/2018/credentials/examples/v1",
					"https://ssiproject.inf.um.es/security/psms/v1",
					"https://ssiproject.inf.um.es/poc/context/v1",
				},
				"credentialSubject": map[string]interface{}{
					"DID":             "did:fabric:DuYKXxLuWnTQzBaa9p1eT1Aya98nirASjVN8dphJjYw",
					"deviceType":      "Server",
					"fluidosRole":     "Customer",
					"holderName":      "FluidosNode",
					"orgIdentifier":   "FLUIDOS_id_23241231412",
					"physicalAddress": "50:80:61:82:ab:c9",
				},
				"expirationDate": expirationDate,
				"id":             "did:fabric:IDDOde0vxswhzfNEQwp05_B209NZ8Ssf3IHHhlrt7Ho3068409",
				"issuanceDate":   issuanceDate,
				"issuer":         "did:fabric:IDDOde0vxswhzfNEQwp05_B209NZ8Ssf3IHHhlrt7Ho",
				"proof": map[string]interface{}{
					"created":            issuanceDate,
					"proofPurpose":       "assertionMethod",
					"proofValue":         "BBVHMPF6Alk2zs934vLFUMg6p83X6TsHI1sE_FJ6GX2AHaFDaXuqZ8PZuQHeRRUKYBIXd2TfGmUvZkec77bXhtEL_yo2wtHiX8vMWUDWQ_fzZ4Y6QG9FYJM2wzaexf43xRiyvYWxiONADhz3sNQHILrHgrSVP1fLyGMIocrrQaGs3xf0-ydEUdfCkpsQNZcFmwMuHh_oUC3MJ5RdkkImP6HIruU-Ke7fM4VYcfnd-Pq7FvwfmSDF33Xbn3Zs0vfp2AQY2WKP1X8IcEuLMea6_0YPHhNRdNn-PA-cUqZXPCDQ2uL40Kud6AdCn7Nms3G5ztUKPD50CXzch8PbyPdVw_mjZYpWEd2xwxc4JIjXeXGPKJo8fjxeQHf5aVyn1HCCg7IIo1UMt-M921Z-hllYZuMrGIOTINKYhPjxBKSXwWw4UTt55k0xzZcrtTM8oFCk0j8Vo8mhjM7albUekUfoQwMaw9xvo9pdYIM9v55B2ooQ8ivVCzNSWoAJIj_YrnyMCfeE0vcYnnYUE_j69Whok_c2FuRODPv1TFnICWk6FGC5bhfxF-xMVT4Pt6uyfWn86os",
					"type":               "PsmsBlsSignature2022",
					"verificationMethod": "did:fabric:IDDOde0vxswhzfNEQwp05_B209NZ8Ssf3IHHhlrt7Ho#xcrh0rw2Gxlx-SZRQQhi5h4YJ9_VEv_hV9X7sV2unlI",
				},
				"type": []string{
					"VerifiableCredential",
					"FluidosCredential",
				},
			},
			"credStorageId": "did:fabric:IDDOde0vxswhzfNEQwp05_B209NZ8Ssf3IHHhlrt7Ho3068409",
		}

		rw.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(rw).Encode(response)
		require.NoError(t, err)
	}))
}
