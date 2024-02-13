/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms

import (
	"fmt"

	psmsapi "github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto/primitive/psms/api"
	subtlepsms "github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto/primitive/psms/subtle"

	"github.com/google/tink/go/core/cryptofmt"
	"github.com/google/tink/go/core/primitiveset"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/keyset"
	tinkpb "github.com/google/tink/go/proto/tink_go_proto"
)

func NewSignerMock(kh *MockHandle) (*subtlepsms.BLS12381G1Signer, error) {
	if kh == nil || kh.PrivateKey == nil || kh.Nattr == 0 {
		return nil, fmt.Errorf("bad key handle")
	}
	return subtlepsms.NewBLS12381G1Signer(kh.PrivateKey), nil
}

// NewSigner returns a PSMS Signer primitive from the given keyset handle.
func NewSigner(h *keyset.Handle) (psmsapi.Signer, error) {
	return NewSignerWithKeyManager(h, nil)
}

// NewSignerWithKeyManager returns a PSMS Signer primitive from the given keyset handle and custom key manager.
func NewSignerWithKeyManager(h *keyset.Handle, km registry.KeyManager) (psmsapi.Signer, error) {
	ps, err := h.PrimitivesWithKeyManager(km)
	if err != nil {
		return nil, fmt.Errorf("psms_sign_factory: cannot obtain primitive set: %w", err)
	}

	return newWrappedSigner(ps)
}

// wrappedSigner is a PSMS Signer implementation that uses the underlying primitive set for psms signing.
type wrappedSigner struct {
	ps *primitiveset.PrimitiveSet
}

// newWrappedSigner constructor creates a new wrappedSigner and checks primitives in ps are all of PSMS Signer type.
func newWrappedSigner(ps *primitiveset.PrimitiveSet) (*wrappedSigner, error) {
	if _, ok := (ps.Primary.Primitive).(psmsapi.Signer); !ok {
		return nil, fmt.Errorf("psms_signer_factory: not a PSMS Signer primitive")
	}

	for _, primitives := range ps.Entries {
		for _, p := range primitives {
			if _, ok := (p.Primitive).(psmsapi.Signer); !ok {
				return nil, fmt.Errorf("psms_signer_factory: not a PSMS Signer primitive")
			}
		}
	}

	ret := new(wrappedSigner)
	ret.ps = ps

	return ret, nil
}

// Sign signs the given messages and returns the signature concatenated with the identifier of the primary primitive.
func (ws *wrappedSigner) Sign(messages [][]byte) ([]byte, error) {
	primary := ws.ps.Primary

	signer, ok := (primary.Primitive).(psmsapi.Signer)
	if !ok {
		return nil, fmt.Errorf("psms_signer_factory: not a PSMS Signer primitive")
	}

	var dataToSign [][]byte
	if primary.PrefixType == tinkpb.OutputPrefixType_LEGACY {
		dataToSign = append(dataToSign, messages...)
		dataToSign = append(dataToSign, []byte{cryptofmt.LegacyStartByte})
	} else {
		dataToSign = append(dataToSign, messages...)
	}

	signature, err := signer.Sign(dataToSign)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, 0, len(primary.Prefix)+len(signature))
	ret = append(ret, primary.Prefix...)
	ret = append(ret, signature...)

	return ret, nil
}
