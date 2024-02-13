/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignatureproof2022

// Package psmsblssignatureproof2022 implements the BBS+ Signature Proof Suite 2020 signature suite
// (https://w3c-ccg.github.io/ldp-bbs2020) in conjunction with the signing and verification algorithms of the
// Linked Data Proofs.
// It uses the RDF Dataset Normalization Algorithm to transform the input document into its canonical form.
// It uses SHA-256 [RFC6234] as the statement digest algorithm.
// It uses BBS+ signature algorithm (https://mattrglobal.github.io/bbs-signatures-spec/).
// It uses BLS12-381 pairing-friendly curve (https://tools.ietf.org/html/draft-irtf-cfrg-pairing-friendly-curves-03).

import (
	"regexp"
	"strings"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
)

// Suite implements BbsBlsSignatureProof2020 signature suite.
type Suite struct {
	suite.SignatureSuite
	jsonldProcessor *jsonld.Processor
}

const (
	signatureType      = "PsmsBlsSignature2022"
	signatureProofType = "PsmsBlsSignatureProof2022"
	rdfDataSetAlg      = "URDNA2015"
)

var excluded = map[string]bool{"http://purl.org/dc/terms/created": true,
	"http://www.w3.org/1999/02/22-rdf-syntax-ns#type": true,
	"https://w3id.org/security#proofPurpose":          true,
	"https://w3id.org/security#verificationMethod":    true}

// New an instance of Linked Data Signatures for the suite.
func New(opts ...suite.Opt) *Suite {
	s := &Suite{jsonldProcessor: jsonld.NewProcessor(rdfDataSetAlg)}

	suite.InitSuiteOptions(&s.SignatureSuite, opts...)

	return s
}

// GetCanonicalDocument will return normalized/canonical version of the document.
// PSMSBlsSignatureProof2022 signature suite uses RDF Dataset Normalization as canonicalization algorithm.
func (s *Suite) GetCanonicalDocument(doc map[string]interface{}, opts ...jsonld.ProcessorOpts) ([]byte, error) {

	if v, ok := doc["type"]; ok { //XXX UMU (ZK/verifiable) Will not be necessary
		docType, ok := v.(string)

		if ok && strings.HasSuffix(docType, signatureProofType) {
			docType = strings.Replace(docType, signatureProofType, signatureType, 1)
			doc["type"] = docType
		}
	}

	return s.jsonldProcessor.GetCanonicalDocument(doc, opts...)
}

// GetDigest returns the doc itself as we would process N-Quads statements as messages to be signed/verified.
func (s *Suite) GetDigest(doc []byte) []byte {
	metattribute := parseMetadata(string(doc))
	sidattrnames, sidattrvalues := parseSignedIdentityAttributes(string(doc)) //TODO Losing typings for now
	return formatIntoLines(metattribute, sidattrnames, sidattrvalues)
}

func formatIntoLines(metattribute string, sidattrnames []string, sidattrvalues []string) []byte {
	res := metattribute
	if res != "" {
		res += ".\n"
	}
	for i := range sidattrvalues {
		res += sidattrnames[i] + " " + sidattrvalues[i] + " ."
		if i < len(sidattrvalues)-1 {
			res += "\n"
		}
	}
	return []byte(res)
}

func parseMetadata(doc string) string {
	//TODO UMU ZK (verifiable) Consider which ones we are using, plus idea of all signed attribute names concatenated as something
	// that is revealed to keep attribute value types for signing. ISSUE: Composite attributes (i.e. JSON), can be requested to be revealed
	//"partially" if the attribute is parsed in a "variable" way, i.e., similar to a credentialSubject. Right now not supported, careful with tests
	//TODO UMU (verifiable/issue) Two calls to digest, one for proof data (including created), one for rest of credential, so processing
	// all metadata at the same time is not possible? Probably need to mess up with the data.go methods
	res := ""
	createdRegex := regexp.MustCompile("<http://purl.org/dc/terms/created>\\s+(.+?)\\s+\\.\\n")
	//expirationRegex := regexp.MustCompile("<https://www.w3.org/2018/credentials#expirationDate>\\s+(.+?)\\s+\\.\\n")
	created := createdRegex.FindStringSubmatch(doc)
	if created != nil {
		res += "<http://purl.org/dc/terms/created>" + created[1] + " "
	}
	//expired := expirationRegex.FindStringSubmatch(doc)
	//if expired != nil {
	//	res += "<https://www.w3.org/2018/credentials#expirationDate" + expired[1] + " "
	//}
	return res
}

func parseSignedIdentityAttributes(doc string) ([]string, []string) {
	subjectRegex := regexp.MustCompile("<https://www.w3.org/2018/credentials#credentialSubject>\\s+<?([^<>\\t\\n\\f\\r ]+)>?\\s+\\.")
	matches := subjectRegex.FindAllStringSubmatch(doc, 10)
	ids := make([]string, 0)
	for _, match := range matches {
		id := match[1]
		if id != "" {
			ids = append(ids, id)
		}
	}
	allattrs := make([]string, 0)
	attrnames := make([]string, 0)
	for _, id := range ids {
		attrRegex := regexp.MustCompile("<?" + id + ">?\\s+<([^>]+)>\\s+(.+)\\s+\\.\n")
		sattrs := attrRegex.FindAllStringSubmatch(doc, 50)
		for _, attr := range sattrs { //TODO UMU (refactor) Should differentiate between different subjects, as they could have the same attributes
			if isIdentifier(attr[2]) {
				recursedValues := parseJsonAttribute(doc, attr[2])
				attrnames = append(attrnames, attr[1])
				allattrs = append(allattrs, recursedValues)
			} else if !excluded[attr[1]] {
				attrnames = append(attrnames, attr[1])
				allattrs = append(allattrs, attr[2])
			}
		}
	}
	return attrnames, allattrs
}

func parseJsonAttribute(doc string, identifier string) string {
	result := ""
	attrRegex := regexp.MustCompile("<?" + identifier + ">?\\s+<([^>]+)>\\s+(.+)\\s+\\.\n")
	sattrs := attrRegex.FindAllStringSubmatch(doc, 50)
	for _, attr := range sattrs {
		if isIdentifier(attr[1]) {
			recursedValues := parseJsonAttribute(doc, attr[1])
			result += recursedValues + "."
		} else {
			result += attr[1] + " " + attr[2] + "."
		}
	}
	return result
}

func isIdentifier(att string) bool {
	res, err := regexp.MatchString("^_:\\S+$", att)
	if err != nil {
		return false
	}
	return res
}

// Accept will accept only BbsBlsSignatureProof2020 signature type.
func (s *Suite) Accept(t string) bool {
	return t == signatureProofType
}
