//go:build !ACAPyInterop
// +build !ACAPyInterop

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import diddoc "github.com/hyperledger/aries-framework-go/pkg/doc/did"

// interopPreprocess in a !ACAPyInterop build, this is a noop.
func interopPreprocess(doc *diddoc.Doc) *diddoc.Doc {
	return doc
}
