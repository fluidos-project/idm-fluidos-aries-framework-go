/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psmsblssignature2022

// Package psmsblssignature2022 implements the PSMS Signature Suite 2022 signature suite
//  in conjunction with the signing and verification algorithms of the
// Linked Data Proofs.
// It uses the RDF Dataset Normalization Algorithm to transform the input document into its canonical form. //TODO change?
// It uses SHA-256 [RFC6234] as the statement digest algorithm. //TODO change?
// It uses PSMS signature algorithm
// It uses BLS12-381 pairing-friendly curve (https://tools.ietf.org/html/draft-irtf-cfrg-pairing-friendly-curves-03).

import (
	"regexp"
	"strings"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
)

var excluded = map[string]bool{"http://purl.org/dc/terms/created": true,
	"http://www.w3.org/1999/02/22-rdf-syntax-ns#type": true,
	"https://w3id.org/security#proofPurpose":          true,
	"https://w3id.org/security#verificationMethod":    true}

// Suite implements PsmsBlsSignature2022 signature suite.
type Suite struct {
	suite.SignatureSuite
	jsonldProcessor *jsonld.Processor
}

const (
	// SignatureType is the PsmsBlsSignature2022 type string.
	SignatureType = "PsmsBlsSignature2022"
	rdfDataSetAlg = "URDNA2015"
)

// New an instance of Linked Data Signatures for JWS suite.
func New(opts ...suite.Opt) *Suite {
	s := &Suite{jsonldProcessor: jsonld.NewProcessor(rdfDataSetAlg)}

	suite.InitSuiteOptions(&s.SignatureSuite, opts...)

	return s
}

// GetCanonicalDocument will return normalized/canonical version of the document.
// PsmsBlsSignature2022 signature suite uses RDF Dataset Normalization as canonicalization algorithm.
func (s *Suite) GetCanonicalDocument(doc map[string]interface{}, opts ...jsonld.ProcessorOpts) ([]byte, error) {
	return s.jsonldProcessor.GetCanonicalDocument(doc, opts...)
}

// GetDigest returns the doc itself as we would process N-Quads statements as messages to be signed/verified.
func (s *Suite) GetDigest(doc []byte) []byte {
	// TODO UMU ZK (verifiable) Fully implement verifiable vision for PSMS. Current is an approximation

	//XXX UMU Canonical/Digest algorithms (verifiable)
	// Get attributes from subjects. Get only values: Use values directly (if this does not let you distinguish later, maybe add a type byte as first byte or something like that)
	// Current solution idea: Extra required attribute, always included in keys (similar to epoch?) that signs the concatenation of the attribute names included
	// Metadata is not signed (except expiration date, considered to be "epoch" for now I guess), nattr of key only includes identity attrs (no metadata)
	// Original ideas:
	// Can we use the canonicalized form (or digested form) of attribute values without extra verification?
	// Example issue: receive document for salary+name, change salary to another attrname (points) that results in the same order,
	// effectively forging a signature on points+name.
	// Issue: if we include signed attrnames (as coming from canonicalization), we cannot do things like range proofs
	// If verification keys have the signed set of attributes, there is no issue anymore I guess
	// Best way to integrate might be to parse canonical doc and return something like bytes of "<epoch:blah><attr0:blah>"
	// Look at separate lines for format
	// How to keep types?
	// Take into account that this will be called with "options" (coming from proof) and document without proofs, in that order
	// It could also be a point (in the future) to actually do the attribute transformation into Zp as we may know here the type definitions
	// Another extra option could be to add an extra attribute to every key (maybe a sub-version of the signature suite?????) that signs a hash of the attribute names
	// that were used, and set it to required revealed.
	// For additional metadata we could follow a similar approach, with a specific attribute for revealing all metadata? Privacy?
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
	//TODO UMU (verifiable) Consider which ones we are using, plus idea of all signed attribute names concatenated as something
	// that is revealed to keep attribute value types for signing
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
				allattrs = append(allattrs, attr[2]+" "+recursedValues) //TODO UMU Added identifier to attribute parsing
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
	res, err := regexp.MatchString("(^_:\\S+$)|<(urn:bnid:_:\\S+?>$)", att) // \\S == [^\t\n\f\r ] //TODO UMU (issue/verifiable) Check how to distinguish non-blank node identifiers from "identifier attributes (e.g., issuer id)"
	if err != nil {
		return false
	}
	return res
}

// Accept will accept only PsmsBlsSignature2022 signature type.
func (s *Suite) Accept(t string) bool {
	return t == SignatureType
}

func splitMessageIntoLines(msg string) []string {
	rows := strings.Split(msg, "\n")

	msgs := make([]string, 0, len(rows))

	for i := range rows {
		if strings.TrimSpace(rows[i]) != "" {
			msgs = append(msgs, rows[i])
		}
	}

	return msgs
}
