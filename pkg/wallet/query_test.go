/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wallet

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/internal/testdata"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/presexch"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/internal/ldtestutil"
	mockvdr "github.com/hyperledger/aries-framework-go/pkg/mock/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"
)

// nolint: lll
const (
	sampleVCFmt = `{
      "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.w3.org/2018/credentials/examples/v1",
		"https://w3id.org/security/bbs/v1"
      ],
     "credentialSchema": [{"id": "%s", "type": "JsonSchemaValidator2018"}],
      "credentialSubject": {
        "degree": {
          "type": "BachelorDegree",
          "university": "MIT"
        },
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "name": "Jayden Doe",
        "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
      },
      "expirationDate": "2020-01-01T19:23:24Z",
      "id": "http://example.edu/credentials/1872",
      "issuanceDate": "2010-01-01T19:23:24Z",
      "issuer": {
        "id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
        "name": "Example University"
      },
      "referenceNumber": 83294847,
      "type": [
        "VerifiableCredential",
        "UniversityDegreeCredential"
      ]
    }`
	samplePRCVC = `{
	 	"@context": [
	   		"https://www.w3.org/2018/credentials/v1",
	   		"https://w3id.org/citizenship/v1",
	   		"https://w3id.org/security/bbs/v1"
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
		"credentialSchema": [],
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
	}`
	sampleBBSVC = `{
            "@context": ["https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1", "https://w3id.org/security/bbs/v1"],
            "credentialSubject": {
                "degree": {"type": "BachelorDegree", "university": "MIT"},
                "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
                "name": "Jayden Doe",
                "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
            },
            "expirationDate": "2020-01-01T19:23:24Z",
            "id": "http://example.edu/credentials/1872",
            "issuanceDate": "2010-01-01T19:23:24Z",
            "issuer": {"id": "did:example:76e12ec712ebc6f1c221ebfeb1f", "name": "Example University"},
            "proof": {
                "created": "2021-03-29T13:27:36.483097-04:00",
                "proofPurpose": "assertionMethod",
                "proofValue": "rw7FeV6K1wimnYogF9qd-N0zmq5QlaIoszg64HciTca-mK_WU4E1jIusKTT6EnN2GZz04NVPBIw4yhc0kTwIZ07etMvfWUlHt_KMoy2CfTw8FBhrf66q4h7Qcqxh_Kxp6yCHyB4A-MmURlKKb8o-4w",
                "type": "BbsBlsSignature2020",
                "verificationMethod": "did:key:zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ#zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ"
            },
            "referenceNumber": 83294847,
            "type": ["VerifiableCredential", "UniversityDegreeCredential"]
        }`
	samplePSMSVC = `{
		"@context":["https://www.w3.org/2018/credentials/v1","https://w3id.org/citizenship/v1","https://ssiproject.inf.um.es/security/psms/v1"],
		"credentialSubject":{
			"birthCountry":"Bahamas",
			"birthDate":"1958-07-17",
			"commuterClassification":"C1",
			"familyName":"SMITH",
			"gender":"Male",
			"givenName":"JOHN",
			"id":"did:example:b34ca6cd37bbf23",
			"image":"data:image/png;base64,iVBORw0KGgokJggg==",
			"lprCategory":"C09",
			"lprNumber":"999-999-999",
			"residentSince":"2015-01-01",
			"type":["PermanentResident","Person"]
		 },
		"description":"Government of Example Permanent Resident Card.",
		"expirationDate":"2029-12-03T12:19:52Z",
		"id":"https://issuer.oidp.uscis.gov/credentials/83627465",
		"issuanceDate":"2019-12-03T12:19:52Z",
		"issuer":"did:key:zfqLqehvuxiwMHJz27XNid1nYWyK3SLWue7gwwdbMxHs6EVwYUEir131KTrQymZJhRrYDZ6tzYexYktXwSHGZ4UboTWinSmg2xtVg4QwzZ8h1cevBU8CYwFC51hfJbCgKyPpwQCGWZv27jzRdHgMT5f41Ffx3GVNKZhDqti1TY7uMrQ8gHxyyjewpE6REMe6toATxf1y8mk9vEq1uQ6AVmqMXQjtUUTBBGzY3HovGWz9vuQZZTYDUutba2aLQeqCda4pvTCUHv7TCgEi9rUAovcsxxkxMCy3LJkpejuy51tRkKHmVKkyVkRr5DssE8zqG8fPqwXq9fmQGDWnJMUpvJJAwf6i9Rn9oBtetuoWPCufrJkVPHSTdC27ZCMQy9WNG3SpLnbc9bH6MrFQW2hmnScbn7ouDAqT5fAmuZX6oPEyKEvPbJS9tudazdgGuSnnfxX8osYJwdR38UECGksEtDxT3njfzHv57DjBSig19ALMfqgTNZdyZdNJzK888v98VWxyCRDy8xHn54FoxptdwAwvupP6aLEdBWYq2cf1fz4pQvC54yLnz24xjG4x7BjphVSSx9YrUkYSgrZsV2b1vmyU5zdRAHMVePzYUjHissU1FN2jCHU4Z5vT2kXxXm4srGJPEWKN8gfC2QQRZa8jUGftmm3JSBmvD82UkmFWcrqYASpEpxy5sWXsjvprqgjP8rrNMcL32yHNWU63WuExkubC49u1GsrkyLnsmsZWzLfrHeBwn4K5vkgXnZ6Bq8Ft5szmhVryMTqtuvHrpjLQkU3LX3ZnBkTr7MxdRpesxik8NSDrGThq3a3Uz71uJF6RtdVHzddQkU4J28XuTM6kmtvovbgDGm9C5oS9x1GEmM4gYVHp3KKNs4niijQHubmVCyyny7HvpHCad1SnRb1niWJwsJiajaKy8eWjbVkdpbufvm6hqhg5q7RZJzibWVT5758xZxmrxScqWCETDgYdaQvTDwf1NYDsuw4jB9YHfdnL6EBbkZvugdA7DLaCipm9dLizQa1KfTz9FUGSH93VSmwghpBXKVFdm4NQ9BjSaizRrtvqgT43aTR51NKkdBjyoJ8pogyQmVoUHoxqMVRHHCxKdAL2HzGot5JzKzXsVy1iAEmQdJHYQsPSPTdKtpK1XFaFrw87KAa9zMFRadteNqKDrtd2pPere8BT3RfM7cgvcHZ2PhsqENzv5UAQ4ZqxPB29gALpxAQTHwhi5zMVwb6JACECVuveUDsAhjXQcV8G9rHRikHDi8F875wHyVn6BxTtedHo7NGpYsnb6t9XooWKYVgUiH9yoSeCPTueKGKaVGqFwW8GWojY7Eowr5rMUmMXQUtLQNeQbyJV4e1pkv9ReCgdgwBQ6kopyrH6mD1JMyVhVzFR34MrxkxuymH9ZvCAgNEhYdgd7AfW7ssu1W2YLHFECAGsjUwuP13dHwwbeM35L1aP9a5aMwtcrWvLH6eRsXqusNAk1Pq4hxAVzmdmw5DMFNQoJEzGGw2bDhiSTpYmEZFXkrrpa9FrBg5oh3DyQrq5pMXx3p4LJ9YcQTGtd4npsWxVAegN6L5EpPcyfY3jfqAE2hvFyNXrFT4SZshKKHChHg5KdFC4hLeE9SGznvYQwdUrYsPpkZS8QjDfJFF2CQMEmZy9pKHYfVFbE1LnQHNfrzVWkP5PxYnmD2bL7wDYtiU4KGEV21rAhx2dfNSb32eFkEDgK4hYJBS7Ue",
		"name":"Permanent Resident Card",
		"proof":{
		   "created":"2023-04-24T16:45:03.115015899+02:00",
		   "proofPurpose":"assertionMethod",
		   "proofValue":"BAYB91H22ZfK53z5IhHZi3TAB0kZ8qLRWgNJnkvuGjH5myfgy6lY361Tnftdq4X0IxJK4tgyNUFWbYBViSZPhQcxvk5hOMG9FqoT0zbXkW7Mrk5KgDEAs09f9KA6DN4m6hMR6xilKqIYZU7BuamtDAIuxjPfjDu1RIsGDqbk2kyb1NP27C-HlwWSeg6U64RsBRD0TA4ANQ_lQYsiXpGluGSMyz4HEHYNMbRPd89K1q8ExnDI0Hmqnw9GuB4f8GNQogQTQV6ns0S75hHbJASxbIpSEAUZT4hmlh6v0lXnEf8VAgcOmHC6LgQB2Cmr_wx5wx4UsYneR3ugr5VthJPNWkvQN05Q_IvOOFCuUsVQaRz8HbbM-fgkzwM0FqZWrrKCh3QNpWatfxOXUeRReiI-X6sBPgxwvIX8WTQtGdSRPA3_M_cLWuQ1MnApzzI8JB7kqD4PeqjfxyPhAwUEd0xE-jelT05uVvomU0K79Cv4TDHUn417gxP1GIF5TrO23nNB5HaHgv4lZNxo2RRkk-QBPRSPRIQMy4PVIuTYJEAqJ_7am51ZLKUepVLgbwN7TmVM59M",
		   "type":"PsmsBlsSignature2022",
		   "verificationMethod":"did:key:zfqLqehvuxiwMHJz27XNid1nYWyK3SLWue7gwwdbMxHs6EVwYUEir131KTrQymZJhRrYDZ6tzYexYktXwSHGZ4UboTWinSmg2xtVg4QwzZ8h1cevBU8CYwFC51hfJbCgKyPpwQCGWZv27jzRdHgMT5f41Ffx3GVNKZhDqti1TY7uMrQ8gHxyyjewpE6REMe6toATxf1y8mk9vEq1uQ6AVmqMXQjtUUTBBGzY3HovGWz9vuQZZTYDUutba2aLQeqCda4pvTCUHv7TCgEi9rUAovcsxxkxMCy3LJkpejuy51tRkKHmVKkyVkRr5DssE8zqG8fPqwXq9fmQGDWnJMUpvJJAwf6i9Rn9oBtetuoWPCufrJkVPHSTdC27ZCMQy9WNG3SpLnbc9bH6MrFQW2hmnScbn7ouDAqT5fAmuZX6oPEyKEvPbJS9tudazdgGuSnnfxX8osYJwdR38UECGksEtDxT3njfzHv57DjBSig19ALMfqgTNZdyZdNJzK888v98VWxyCRDy8xHn54FoxptdwAwvupP6aLEdBWYq2cf1fz4pQvC54yLnz24xjG4x7BjphVSSx9YrUkYSgrZsV2b1vmyU5zdRAHMVePzYUjHissU1FN2jCHU4Z5vT2kXxXm4srGJPEWKN8gfC2QQRZa8jUGftmm3JSBmvD82UkmFWcrqYASpEpxy5sWXsjvprqgjP8rrNMcL32yHNWU63WuExkubC49u1GsrkyLnsmsZWzLfrHeBwn4K5vkgXnZ6Bq8Ft5szmhVryMTqtuvHrpjLQkU3LX3ZnBkTr7MxdRpesxik8NSDrGThq3a3Uz71uJF6RtdVHzddQkU4J28XuTM6kmtvovbgDGm9C5oS9x1GEmM4gYVHp3KKNs4niijQHubmVCyyny7HvpHCad1SnRb1niWJwsJiajaKy8eWjbVkdpbufvm6hqhg5q7RZJzibWVT5758xZxmrxScqWCETDgYdaQvTDwf1NYDsuw4jB9YHfdnL6EBbkZvugdA7DLaCipm9dLizQa1KfTz9FUGSH93VSmwghpBXKVFdm4NQ9BjSaizRrtvqgT43aTR51NKkdBjyoJ8pogyQmVoUHoxqMVRHHCxKdAL2HzGot5JzKzXsVy1iAEmQdJHYQsPSPTdKtpK1XFaFrw87KAa9zMFRadteNqKDrtd2pPere8BT3RfM7cgvcHZ2PhsqENzv5UAQ4ZqxPB29gALpxAQTHwhi5zMVwb6JACECVuveUDsAhjXQcV8G9rHRikHDi8F875wHyVn6BxTtedHo7NGpYsnb6t9XooWKYVgUiH9yoSeCPTueKGKaVGqFwW8GWojY7Eowr5rMUmMXQUtLQNeQbyJV4e1pkv9ReCgdgwBQ6kopyrH6mD1JMyVhVzFR34MrxkxuymH9ZvCAgNEhYdgd7AfW7ssu1W2YLHFECAGsjUwuP13dHwwbeM35L1aP9a5aMwtcrWvLH6eRsXqusNAk1Pq4hxAVzmdmw5DMFNQoJEzGGw2bDhiSTpYmEZFXkrrpa9FrBg5oh3DyQrq5pMXx3p4LJ9YcQTGtd4npsWxVAegN6L5EpPcyfY3jfqAE2hvFyNXrFT4SZshKKHChHg5KdFC4hLeE9SGznvYQwdUrYsPpkZS8QjDfJFF2CQMEmZy9pKHYfVFbE1LnQHNfrzVWkP5PxYnmD2bL7wDYtiU4KGEV21rAhx2dfNSb32eFkEDgK4hYJBS7Ue#zfqLqehvuxiwMHJz27XNid1nYWyK3SLWue7gwwdbMxHs6EVwYUEir131KTrQymZJhRrYDZ6tzYexYktXwSHGZ4UboTWinSmg2xtVg4QwzZ8h1cevBU8CYwFC51hfJbCgKyPpwQCGWZv27jzRdHgMT5f41Ffx3GVNKZhDqti1TY7uMrQ8gHxyyjewpE6REMe6toATxf1y8mk9vEq1uQ6AVmqMXQjtUUTBBGzY3HovGWz9vuQZZTYDUutba2aLQeqCda4pvTCUHv7TCgEi9rUAovcsxxkxMCy3LJkpejuy51tRkKHmVKkyVkRr5DssE8zqG8fPqwXq9fmQGDWnJMUpvJJAwf6i9Rn9oBtetuoWPCufrJkVPHSTdC27ZCMQy9WNG3SpLnbc9bH6MrFQW2hmnScbn7ouDAqT5fAmuZX6oPEyKEvPbJS9tudazdgGuSnnfxX8osYJwdR38UECGksEtDxT3njfzHv57DjBSig19ALMfqgTNZdyZdNJzK888v98VWxyCRDy8xHn54FoxptdwAwvupP6aLEdBWYq2cf1fz4pQvC54yLnz24xjG4x7BjphVSSx9YrUkYSgrZsV2b1vmyU5zdRAHMVePzYUjHissU1FN2jCHU4Z5vT2kXxXm4srGJPEWKN8gfC2QQRZa8jUGftmm3JSBmvD82UkmFWcrqYASpEpxy5sWXsjvprqgjP8rrNMcL32yHNWU63WuExkubC49u1GsrkyLnsmsZWzLfrHeBwn4K5vkgXnZ6Bq8Ft5szmhVryMTqtuvHrpjLQkU3LX3ZnBkTr7MxdRpesxik8NSDrGThq3a3Uz71uJF6RtdVHzddQkU4J28XuTM6kmtvovbgDGm9C5oS9x1GEmM4gYVHp3KKNs4niijQHubmVCyyny7HvpHCad1SnRb1niWJwsJiajaKy8eWjbVkdpbufvm6hqhg5q7RZJzibWVT5758xZxmrxScqWCETDgYdaQvTDwf1NYDsuw4jB9YHfdnL6EBbkZvugdA7DLaCipm9dLizQa1KfTz9FUGSH93VSmwghpBXKVFdm4NQ9BjSaizRrtvqgT43aTR51NKkdBjyoJ8pogyQmVoUHoxqMVRHHCxKdAL2HzGot5JzKzXsVy1iAEmQdJHYQsPSPTdKtpK1XFaFrw87KAa9zMFRadteNqKDrtd2pPere8BT3RfM7cgvcHZ2PhsqENzv5UAQ4ZqxPB29gALpxAQTHwhi5zMVwb6JACECVuveUDsAhjXQcV8G9rHRikHDi8F875wHyVn6BxTtedHo7NGpYsnb6t9XooWKYVgUiH9yoSeCPTueKGKaVGqFwW8GWojY7Eowr5rMUmMXQUtLQNeQbyJV4e1pkv9ReCgdgwBQ6kopyrH6mD1JMyVhVzFR34MrxkxuymH9ZvCAgNEhYdgd7AfW7ssu1W2YLHFECAGsjUwuP13dHwwbeM35L1aP9a5aMwtcrWvLH6eRsXqusNAk1Pq4hxAVzmdmw5DMFNQoJEzGGw2bDhiSTpYmEZFXkrrpa9FrBg5oh3DyQrq5pMXx3p4LJ9YcQTGtd4npsWxVAegN6L5EpPcyfY3jfqAE2hvFyNXrFT4SZshKKHChHg5KdFC4hLeE9SGznvYQwdUrYsPpkZS8QjDfJFF2CQMEmZy9pKHYfVFbE1LnQHNfrzVWkP5PxYnmD2bL7wDYtiU4KGEV21rAhx2dfNSb32eFkEDgK4hYJBS7Ue"
		},
		"type":["VerifiableCredential","PermanentResidentCard"]
	 }`
	sampleQueryByExFmt = `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://www.w3.org/2018/credentials/examples/v1",
								"https://w3id.org/security/bbs/v1"
                            ],
                            "type": ["UniversityDegreeCredential"],
							"trustedIssuer": [
              					{
                					"issuer": "urn:some:required:issuer"
              					},
								{
                					"required": true,
                					"issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f"
              					}
							],
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"	
							},
							"credentialSchema": {
								"id": "%s",
								"type": "JsonSchemaValidator2018"
							}
                        }
                	}`
	sampleQueryByFrame = `{
                    "reason": "Please provide your Passport details.",
                    "frame": {
                        "@context": [
                            "https://www.w3.org/2018/credentials/v1",
                            "https://w3id.org/citizenship/v1",
                            "https://w3id.org/security/bbs/v1"
                        ],
                        "type": ["VerifiableCredential", "PermanentResidentCard"],
                        "@explicit": true,
                        "identifier": {},
                        "issuer": {},
                        "issuanceDate": {},
                        "credentialSubject": {
                            "@explicit": true,
                            "name": {},
                            "spouse": {}
                        }
                    },
                    "trustedIssuer": [
                        {
                            "issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f",
                            "required": true
                        }
                    ],
                    "required": true
                }`
	sampleQueryByFramePsms = `{
                    "reason": "Please provide your Resident card details.",
                    "frame": {
                        "@context": [
                            "https://www.w3.org/2018/credentials/v1",
                            "https://w3id.org/citizenship/v1",
                            "https://ssiproject.inf.um.es/security/psms/v1"
                        ],
						"type": ["VerifiableCredential", "PermanentResidentCard"],
                        "@explicit": true,
                        "issuer": {},
                        "issuanceDate": {},
						"credentialSubject": {
							"@explicit": true,
							"type": ["PermanentResident", "Person"],
							"givenName": {},
							"familyName": {}						}
                    },
                    "trustedIssuer": [
                        {
                            "issuer": "did:key:zfqLqehvuxiwMHJz27XNid1nYWyK3SLWue7gwwdbMxHs6EVwYUEir131KTrQymZJhRrYDZ6tzYexYktXwSHGZ4UboTWinSmg2xtVg4QwzZ8h1cevBU8CYwFC51hfJbCgKyPpwQCGWZv27jzRdHgMT5f41Ffx3GVNKZhDqti1TY7uMrQ8gHxyyjewpE6REMe6toATxf1y8mk9vEq1uQ6AVmqMXQjtUUTBBGzY3HovGWz9vuQZZTYDUutba2aLQeqCda4pvTCUHv7TCgEi9rUAovcsxxkxMCy3LJkpejuy51tRkKHmVKkyVkRr5DssE8zqG8fPqwXq9fmQGDWnJMUpvJJAwf6i9Rn9oBtetuoWPCufrJkVPHSTdC27ZCMQy9WNG3SpLnbc9bH6MrFQW2hmnScbn7ouDAqT5fAmuZX6oPEyKEvPbJS9tudazdgGuSnnfxX8osYJwdR38UECGksEtDxT3njfzHv57DjBSig19ALMfqgTNZdyZdNJzK888v98VWxyCRDy8xHn54FoxptdwAwvupP6aLEdBWYq2cf1fz4pQvC54yLnz24xjG4x7BjphVSSx9YrUkYSgrZsV2b1vmyU5zdRAHMVePzYUjHissU1FN2jCHU4Z5vT2kXxXm4srGJPEWKN8gfC2QQRZa8jUGftmm3JSBmvD82UkmFWcrqYASpEpxy5sWXsjvprqgjP8rrNMcL32yHNWU63WuExkubC49u1GsrkyLnsmsZWzLfrHeBwn4K5vkgXnZ6Bq8Ft5szmhVryMTqtuvHrpjLQkU3LX3ZnBkTr7MxdRpesxik8NSDrGThq3a3Uz71uJF6RtdVHzddQkU4J28XuTM6kmtvovbgDGm9C5oS9x1GEmM4gYVHp3KKNs4niijQHubmVCyyny7HvpHCad1SnRb1niWJwsJiajaKy8eWjbVkdpbufvm6hqhg5q7RZJzibWVT5758xZxmrxScqWCETDgYdaQvTDwf1NYDsuw4jB9YHfdnL6EBbkZvugdA7DLaCipm9dLizQa1KfTz9FUGSH93VSmwghpBXKVFdm4NQ9BjSaizRrtvqgT43aTR51NKkdBjyoJ8pogyQmVoUHoxqMVRHHCxKdAL2HzGot5JzKzXsVy1iAEmQdJHYQsPSPTdKtpK1XFaFrw87KAa9zMFRadteNqKDrtd2pPere8BT3RfM7cgvcHZ2PhsqENzv5UAQ4ZqxPB29gALpxAQTHwhi5zMVwb6JACECVuveUDsAhjXQcV8G9rHRikHDi8F875wHyVn6BxTtedHo7NGpYsnb6t9XooWKYVgUiH9yoSeCPTueKGKaVGqFwW8GWojY7Eowr5rMUmMXQUtLQNeQbyJV4e1pkv9ReCgdgwBQ6kopyrH6mD1JMyVhVzFR34MrxkxuymH9ZvCAgNEhYdgd7AfW7ssu1W2YLHFECAGsjUwuP13dHwwbeM35L1aP9a5aMwtcrWvLH6eRsXqusNAk1Pq4hxAVzmdmw5DMFNQoJEzGGw2bDhiSTpYmEZFXkrrpa9FrBg5oh3DyQrq5pMXx3p4LJ9YcQTGtd4npsWxVAegN6L5EpPcyfY3jfqAE2hvFyNXrFT4SZshKKHChHg5KdFC4hLeE9SGznvYQwdUrYsPpkZS8QjDfJFF2CQMEmZy9pKHYfVFbE1LnQHNfrzVWkP5PxYnmD2bL7wDYtiU4KGEV21rAhx2dfNSb32eFkEDgK4hYJBS7Ue",
                            "required": true
                        }
                    ],
                    "required": true
                }`
)

func TestGetQueryType(t *testing.T) {
	t.Run("test get query type by string", func(t *testing.T) {
		tests := []struct {
			name         string
			typeStr      []string
			expected     QueryType
			expectedName string
			error        string
		}{
			{
				name:         "test for QueryByExample",
				typeStr:      []string{"QueryByExample", "QuerybyExample", "querybyexample"},
				expected:     QueryByExample,
				expectedName: "QueryByExample",
			},
			{
				name:         "test for QueryByFrame",
				typeStr:      []string{"QueryByFrame", "Querybyframe", "querybyframe"},
				expected:     QueryByFrame,
				expectedName: "QueryByFrame",
			},
			{
				name:         "test for PresentationExchange",
				typeStr:      []string{"PresentationExchange", "Presentationexchange", "presentationExchange"},
				expected:     PresentationExchange,
				expectedName: "PresentationExchange",
			},
			{
				name:         "test for DIDAuth",
				typeStr:      []string{"didAuth", "didauth", "DIDAuth", "DIDauth"},
				expected:     DIDAuth,
				expectedName: "DIDAuth",
			},
			{
				name:         "test for invalid types",
				typeStr:      []string{"", "QueryByFram", "QueryByExamples", "DIDAuthorization", "invalid"},
				error:        "unsupported query type",
				expectedName: "",
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.name, func(t *testing.T) {
				for _, str := range tc.typeStr {
					qType, err := GetQueryType(str)
					require.Equal(t, qType, tc.expected)
					if tc.error != "" {
						require.Error(t, err)
						require.Contains(t, err.Error(), tc.error)
					} else {
						require.NoError(t, err)
					}
				}
			})
		}
	})
}

func TestQuery_PerformQuery(t *testing.T) {
	vc1, err := (&verifiable.Credential{
		Context: []string{verifiable.ContextURI},
		Types:   []string{verifiable.VCType},
		ID:      "http://example.edu/credentials/9999",
		CustomFields: map[string]interface{}{
			"first_name": "Jesse",
		},
		Issued: &util.TimeWrapper{
			Time: time.Now(),
		},
		Issuer: verifiable.Issuer{
			ID: "did:example:76e12ec712ebc6f1c221ebfeb1f",
		},
		Subject: uuid.New().String(),
	}).MarshalJSON()
	require.NoError(t, err)

	// presentation exchange query
	pd := &presexch.PresentationDefinition{
		ID: uuid.New().String(),
		InputDescriptors: []*presexch.InputDescriptor{{
			ID: uuid.New().String(),
			Schema: []*presexch.Schema{{
				URI: fmt.Sprintf("%s#%s", verifiable.ContextID, verifiable.VCType),
			}},
			Constraints: &presexch.Constraints{
				Fields: []*presexch.Field{{
					Path: []string{"$.first_name"},
				}},
			},
		}},
	}

	// query by example
	queryByExample := []byte(fmt.Sprintf(sampleQueryByExFmt, verifiable.ContextURI))
	// query by frame
	queryByFrame := []byte(sampleQueryByFrame)
	queryByFramePsms := []byte(sampleQueryByFramePsms)

	pdJSON, err := json.Marshal(pd)
	require.NoError(t, err)
	require.NotEmpty(t, pdJSON)

	udcVC := testdata.SampleUDCVC
	vcForQuery := []byte(fmt.Sprintf(sampleVCFmt, verifiable.ContextURI))
	vcForDerive := []byte(sampleBBSVC)
	vcForDerivePsms := []byte(samplePSMSVC)

	customVDR := &mockvdr.MockVDRegistry{
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
	pubKeyFetcher := verifiable.NewVDRKeyResolver(customVDR).PublicKeyFetcher()

	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	t.Run("test wallet queries", func(t *testing.T) {
		tests := []struct {
			name        string
			query       []*QueryParams
			credentials []json.RawMessage
			resultCount int
			vcCount     map[int]int
			error       string
		}{
			// Presentation Exchange tests
			{
				name: "query by presentation exchange - success",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by presentation exchange - multi query frame - success",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON, pdJSON}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 2,
				vcCount:     map[int]int{0: 1, 1: 1},
			},
			{
				name: "query by presentation exchange - multiple - success",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON, pdJSON}},
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 3,
				vcCount:     map[int]int{0: 1, 1: 1, 2: 1},
			},
			{
				name: "query by presentation exchange - no results",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
				},
				credentials: []json.RawMessage{udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       ErrQueryNoResultFound.Error(),
			},
			{
				name: "query by presentation exchange - invalid definition",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{[]byte(sampleInvalidDIDContent)}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       "input_descriptors is required",
			},
			{
				name: "query by presentation exchange - invalid query frame",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{[]byte("---")}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       "invalid character",
			},
			// QueryByExample tests
			{
				name: "query by example - success",
				query: []*QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by example - success & normalize results",
				query: []*QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample, queryByExample}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by example - invalid query",
				query: []*QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{[]byte(sampleInvalidDIDContent)}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       "failed to parse QueryByExample query",
			},
			// QueryByFrame tests
			{
				name: "query by frame - success",
				query: []*QueryParams{
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			{
				name: "query by frame - multiple results",
				query: []*QueryParams{
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame, queryByFrame}},
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 3},
			},
			{
				name: "query by frame - invalid query",
				query: []*QueryParams{
					{Type: "QueryByFrame", Query: []json.RawMessage{[]byte(sampleInvalidDIDContent)}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       "failed to parse QueryByFrame query",
			},
			{
				name: "query by frame - success PSMS",
				query: []*QueryParams{
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFramePsms}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 1},
			},
			// DIDAuth tests
			{
				name: "didAuth - success",
				query: []*QueryParams{
					{Type: "DIDAuth"},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 1,
				vcCount:     map[int]int{0: 0},
			},

			// Mixed Query Types
			{
				name: "query by PresentationExchange,QueryByExample,QueryByFrame - success",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
					{Type: "DIDAuth"},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 3,
				vcCount:     map[int]int{0: 1, 1: 0, 2: 2},
			},
			{
				name: "query by PresentationExchange,QueryByExample,QueryByFrame - normalized result - success",
				query: []*QueryParams{
					{Type: "PresentationExchange", Query: []json.RawMessage{pdJSON}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
					{Type: "QueryByFrame", Query: []json.RawMessage{queryByFrame}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 2,
				vcCount:     map[int]int{0: 1, 1: 2},
			},

			// Validations
			{
				name: "unsupported query type",
				query: []*QueryParams{
					{Type: "QueryByInvalid", Query: []json.RawMessage{queryByExample}},
				},
				credentials: []json.RawMessage{vc1, udcVC, vcForQuery, vcForDerive, vcForDerivePsms},
				resultCount: 0,
				error:       "unsupported query type",
			},
			{
				name: "empty credentials",
				query: []*QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
				},
				credentials: []json.RawMessage{},
				resultCount: 0,
				error:       ErrQueryNoResultFound.Error(),
			},
			{
				name: "credential parsing error",
				query: []*QueryParams{
					{Type: "QueryByExample", Query: []json.RawMessage{queryByExample}},
				},
				credentials: []json.RawMessage{[]byte("----")},
				resultCount: 0,
				error:       "unmarshal new credential",
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.name, func(t *testing.T) {
				credentials := make(map[string]json.RawMessage)
				for i, v := range tc.credentials {
					credentials[strconv.Itoa(i)] = v
				}

				results, err := NewQuery(pubKeyFetcher, loader, tc.query...).PerformQuery(credentials)

				if tc.error != "" {
					require.Empty(t, results)
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.error)
					require.Len(t, results, tc.resultCount)

					return
				}

				require.NoError(t, err)
				require.Len(t, results, tc.resultCount)

				for i, result := range results {
					require.Len(t, result.Credentials(), tc.vcCount[i])
				}
			})
		}
	})
}

func TestQueryByExample(t *testing.T) {
	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	vc1, err := verifiable.ParseCredential([]byte(fmt.Sprintf(sampleVCFmt, verifiable.ContextURI)),
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	vc2, err := verifiable.ParseCredential([]byte(samplePRCVC), verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	// sample queries
	queryByExampleAll := []byte(fmt.Sprintf(sampleQueryByExFmt, verifiable.ContextURI))

	queryByExampleContext1 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": "VerifiableCredential",
							"trustedIssuer": [],
							"credentialSubject": {},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContext2 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://w3id.org/citizenship/v1"
                            ],
							"type": "VerifiableCredential"
                        }
                	}`

	queryByExampleContextType1 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://w3id.org/citizenship/v1"
                            ],
                            "type": "PermanentResidentCard",
							"trustedIssuer": [],
							"credentialSubject": {},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContextType2 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://www.w3.org/2018/credentials/examples/v1"
                            ],
                            "type": ["UniversityDegreeCredential"],
							"trustedIssuer": [],
							"credentialSubject": {},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContextTypeIssuer1 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://www.w3.org/2018/credentials/examples/v1"
                            ],
                            "type": ["UniversityDegreeCredential"],
							"trustedIssuer": [
              					{
                					"issuer": "urn:some:required:issuer"
              					},
								{
                					"required": true,
                					"issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f"
              					}
							]
                        }
                	}`

	queryByExampleContextTypeIssuer2 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["PermanentResidentCard"],
							"trustedIssuer": [
								{
                					"required": true,
                					"issuer": "did:example:489398593"
              					}
							]
                        }
                	}`

	queryByExampleContextTypeIssuer3 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["PermanentResidentCard"],
							"trustedIssuer": [
								{
                					"issuer": "did:example:489398593"
              					},
								{
                					"required": true,
                					"issuer": "did:example:1234"
              					}
							]
                        }
                	}`

	queryByExampleContextTypeIssuer4 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["PermanentResidentCard"],
							"trustedIssuer": [
								{
                					"issuer": "did:example:7777"
              					},
								{
                					"issuer": "did:example:1234"
              					}
							]
                        }
                	}`

	queryByExampleContextTypeCredSubject1 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["PermanentResidentCard"],
							"trustedIssuer": [
								{
                					"required": true,
                					"issuer": "did:example:489398593"
              					}
							],
							"credentialSubject": {
								"id": "did:example:b34ca6cd37bbf23"	
							},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContextTypeCredSubject2 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["UniversityDegreeCredential"],
							"trustedIssuer": [
								{
                					"required": true,
                					"issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f"
              					}
							],
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"	
							},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContextTypeCredSubject3 := `{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": ["UniversityDegreeCredential"],
							"trustedIssuer": [
								{
                					"required": true,
                					"issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f"
              					}
							],
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec22"	
							},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleContextTypeCredSchema1 := fmt.Sprintf(`{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": "UniversityDegreeCredential",
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"	
							},
							"credentialSchema": {
								"id": "%s",
								"type": "JsonSchemaValidator2018"
							}
                        }
                	}`, verifiable.ContextURI)

	queryByExampleContextTypeCredSchema2 := fmt.Sprintf(`{
                        "reason": "Please present your identity document.",
                        "example": {
                            "@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": "UniversityDegreeCredential",
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"	
							},
							"credentialSchema": {
								"id": "%s",
								"type": "JsonSchemaValidator2020"
							}
                        }
                	}`, verifiable.ContextURI)

	queryByExampleInvalid1 := `{
                        "reason": "Please present your identity document.",
                        "example": {
							"trustedIssuer": [
								{
                					"required": true,
                					"issuer": "did:example:489398593"
              					}
							],
							"credentialSubject": {
								"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"	
							},
							"credentialSchema": {}
                        }
                	}`

	queryByExampleInvalid2 := `{
                        "reason": "Please present your identity document.",
                        "example": {
							"@context": [
								"https://www.w3.org/2018/credentials/v1"
                            ],
                            "type": []
                        }
                	}`

	queryByExampleInvalid3 := `{
                        "reason": "Please present your identity document.",
                        "example": {
							"credentialSubject": "invalid",
							"credentialSchema": true
                        }
                	}`

	vc3 := &verifiable.Credential{
		Context: []string{verifiable.ContextURI},
		Types:   []string{verifiable.VCType},
		ID:      "http://example.edu/credentials/9999",
		CustomFields: map[string]interface{}{
			"first_name": "Jesse",
		},
		Issued: &util.TimeWrapper{
			Time: time.Now(),
		},
		Issuer: verifiable.Issuer{
			ID: "did:example:76e12ec712ebc6f1c221ebfeb1f",
		},
		Subject: uuid.New().String(),
	}

	t.Run("test query", func(t *testing.T) {
		tests := []struct {
			name        string
			credentials []*verifiable.Credential
			example     []json.RawMessage
			resultCount int
			error       string
			skip        bool
		}{
			{
				name:        "QueryByExample all criteria matched - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{queryByExampleAll},
				resultCount: 1,
			},
			{
				name:        "QueryByExample  multiple query frame - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{queryByExampleAll, queryByExampleAll},
				resultCount: 2,
			},
			{
				name:        "QueryByExample context matching #1 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContext1)},
				resultCount: 3,
			},
			{
				name:        "QueryByExample context matching #2 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContext2)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context & type matching #1 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextType1)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context & type matching #2 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextType2)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type & issuer matching #1 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeIssuer1)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type & issuer matching #2 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeIssuer2)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type & issuer matching #3 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeIssuer3)},
				resultCount: 0,
			},
			{
				name:        "QueryByExample context, type & issuer matching #4 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeIssuer4)},
				resultCount: 0,
			},
			{
				name:        "QueryByExample context, type, issuer & subject matching #1 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeCredSubject1)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type, issuer & subject matching #2 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeCredSubject2)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type, issuer & subject matching #3 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeCredSubject3)},
				resultCount: 0,
			},
			{
				name:        "QueryByExample context, type, issuer, schema matching #1 - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeCredSchema1)},
				resultCount: 1,
			},
			{
				name:        "QueryByExample context, type, issuer, schema matching #2- success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleContextTypeCredSchema2)},
				resultCount: 0,
			},
			{
				name:        "QueryByExample invalid query #1",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleInvalid1)},
				resultCount: 0,
				error:       "failed to parse QueryByExample query",
			},
			{
				name:        "QueryByExample invalid query #2",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleInvalid2)},
				resultCount: 0,
				error:       "failed to parse QueryByExample query",
			},
			{
				name:        "QueryByExample invalid query #3",
				credentials: []*verifiable.Credential{vc1, vc2, vc3},
				example:     []json.RawMessage{[]byte(queryByExampleInvalid3)},
				resultCount: 0,
				error:       "failed to parse QueryByExample query",
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.name, func(t *testing.T) {
				if tc.skip {
					return
				}

				results, err := queryByExample(tc.credentials, tc.example...)

				if tc.error != "" {
					require.Empty(t, results)
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.error)
					require.Len(t, results, tc.resultCount)

					return
				}

				require.NoError(t, err)
				require.Len(t, results, tc.resultCount)
			})
		}
	})
}

func TestQueryByFrame(t *testing.T) {
	loader, err := ldtestutil.DocumentLoader()
	require.NoError(t, err)

	vc1, err := verifiable.ParseCredential([]byte(fmt.Sprintf(sampleVCFmt, verifiable.ContextURI)),
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	vc2, err := verifiable.ParseCredential([]byte(samplePRCVC), verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	vc3, err := verifiable.ParseCredential([]byte(sampleBBSVC), verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	tampered := strings.ReplaceAll(sampleBBSVC, `rw7FeV6K1wimnYogF9qd-N0zmq5QlaIoszg64HciTca`, ``)
	vc4, err := verifiable.ParseCredential([]byte(tampered), verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	vc5, err := verifiable.ParseCredential([]byte(samplePSMSVC), verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDDocumentLoader(loader))
	require.NoError(t, err)

	queryByFrameExampleNoIssuer := `{
	                "reason": "Please provide your Passport details.",
	                "frame": {
	                    "@context": [
	                        "https://www.w3.org/2018/credentials/v1",
	                        "https://w3id.org/citizenship/v1",
	                        "https://w3id.org/security/bbs/v1"
	                    ],
	                    "type": ["VerifiableCredential", "PermanentResidentCard"],
	                    "@explicit": true,
	                    "identifier": {},
	                    "issuer": {},
	                    "issuanceDate": {},
	                    "credentialSubject": {
	                        "@explicit": true,
	                        "name": {},
	                        "spouse": {}
	                    }
	                }
	            }`

	queryByFrameExampleIssuerNotMatched := strings.ReplaceAll(sampleQueryByFrame,
		"did:example:76e12ec712ebc6f1c221ebfeb1f", "1234")

	customVDR := &mockvdr.MockVDRegistry{
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
	pubKeyFetcher := verifiable.NewVDRKeyResolver(customVDR).PublicKeyFetcher()

	t.Run("test query", func(t *testing.T) {
		tests := []struct {
			name        string
			credentials []*verifiable.Credential
			example     []json.RawMessage
			resultCount int
			error       string
			skip        bool
		}{
			{
				name:        "QueryByFrame all criteria matched - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3, vc5},
				example:     []json.RawMessage{[]byte(sampleQueryByFrame)},
				resultCount: 1,
			},
			{
				name:        "QueryByFrame multiple query - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3, vc5},
				example: []json.RawMessage{
					[]byte(sampleQueryByFrame), []byte(sampleQueryByFrame),
					[]byte(sampleQueryByFramePsms),
				},
				resultCount: 3,
			},
			{
				name:        "QueryByFrame without issuer criteria - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3, vc5},
				example:     []json.RawMessage{[]byte(queryByFrameExampleNoIssuer)},
				resultCount: 1,
			},
			{
				name:        "QueryByFrame issuer not matched - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc3, vc5},
				example:     []json.RawMessage{[]byte(queryByFrameExampleIssuerNotMatched)},
				resultCount: 0,
			},
			{
				name:        "QueryByFrame NO BBS signature - success",
				credentials: []*verifiable.Credential{vc1, vc2},
				example:     []json.RawMessage{[]byte(queryByFrameExampleIssuerNotMatched)},
				resultCount: 0,
			},
			{
				name:        "QueryByFrame invalid BBS signature - success",
				credentials: []*verifiable.Credential{vc1, vc2, vc4},
				example:     []json.RawMessage{[]byte(sampleQueryByFrame)},
				resultCount: 0,
			},
			{
				name:        "QueryByFrame invalid query - success",
				credentials: []*verifiable.Credential{vc1, vc3},
				example:     []json.RawMessage{[]byte(sampleInvalidDIDContent)},
				error:       "failed to parse QueryByFrame query",
			},
			{
				name:        "QueryByFrame parse query error - success",
				credentials: []*verifiable.Credential{vc1, vc3},
				example:     []json.RawMessage{[]byte("---")},
				error:       "failed to parse QueryByFrame query",
			},
		}

		t.Parallel()

		for _, test := range tests {
			tc := test
			t.Run(tc.name, func(t *testing.T) {
				if tc.skip {
					return
				}

				results, err := queryByFrame(tc.credentials, pubKeyFetcher, loader, tc.example...)

				if tc.error != "" {
					require.Empty(t, results)
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.error)
					require.Len(t, results, tc.resultCount)

					return
				}

				require.NoError(t, err)
				require.Len(t, results, tc.resultCount)
			})
		}
	})
}

func TestUtilFunctions(t *testing.T) {
	require.True(t, isEmpty(""))
	require.True(t, isEmpty([]string{}))
	require.True(t, isEmpty([]interface{}{}))
	require.True(t, isEmpty(nil))

	require.False(t, isEmpty("x"))
	require.False(t, isEmpty([]string{""}))
	require.False(t, isEmpty([]interface{}{&QueryParams{}}))
	require.False(t, isEmpty(&QueryParams{}))

	require.False(t, contains([]string{"a", "b"}, nil))
}
