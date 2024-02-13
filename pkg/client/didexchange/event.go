/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package didexchange

import (
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/model"
)

// Event properties related api. This can be used to cast Generic event properties to DID Exchange specific props.
type Event model.Event
