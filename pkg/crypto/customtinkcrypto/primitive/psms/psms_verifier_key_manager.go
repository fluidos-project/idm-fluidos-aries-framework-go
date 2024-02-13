/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms

import (
	"errors"
)

const (
	psmsVerifierKeyVersion = 0
	psmsVerifierKeyTypeURL = "type.hyperledger.org/hyperledger.aries.crypto.tink.PSMSPublicKey"
)

// common errors.
var errInvalidPSMSVerifierKey = errors.New("psms_verifier_key_manager: invalid key")

// psmsVerifierKeyManager is an implementation of KeyManager interface for PSMS signature/proof verification.
// It doesn't support key generation.
type psmsVerifierKeyManager struct{}

/*
// newPSMSVerifierKeyManager creates a new psmsVerifierKeyManager.
func newPSMSVerifierKeyManager() *psmsVerifierKeyManager {
	return new(psmsVerifierKeyManager)
}

// Primitive creates an PSMS Verifier subtle for the given serialized PSMSPublicKey proto.
func (km *psmsVerifierKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidPSMSVerifierKey
	}

	psmsPubKey := new(psmspb.PSMSPublicKey)

	err := proto.Unmarshal(serializedKey, psmsPubKey)
	if err != nil {
		return nil, errInvalidPSMSVerifierKey
	}

	err = km.validateKey(psmsPubKey)
	if err != nil {
		return nil, errInvalidPSMSVerifierKey
	}

	return subtle.NewBLS12381G2Verifier(psmsPubKey.KeyValue), nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *psmsVerifierKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == psmsVerifierKeyTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *psmsVerifierKeyManager) TypeURL() string {
	return psmsVerifierKeyTypeURL
}

// NewKey is not implemented for public key manager.
func (km *psmsVerifierKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	return nil, errors.New("psms_verifier_key_manager: NewKey not implemented")
}

// NewKeyData is not implemented for public key manager.
func (km *psmsVerifierKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	return nil, errors.New("psms_verifier_key_manager: NewKeyData not implemented")
}

// validateKey validates the given EcdhAeadPublicKey.
func (km *psmsVerifierKeyManager) validateKey(key *psmspb.PSMSPublicKey) error {
	err := keyset.ValidateKeyVersion(key.Version, psmsVerifierKeyVersion)
	if err != nil {
		return fmt.Errorf("psms_verifier_key_manager: invalid key: %w", err)
	}

	return nil //validateKeyParams(key.Params)
}


*/
