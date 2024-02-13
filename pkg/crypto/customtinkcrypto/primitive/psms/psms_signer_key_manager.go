/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package psms

import (
	"errors"
)

const (
	psmsSignerKeyVersion = 0
	PsmsSignerKeyTypeURL = "type.hyperledger.org/hyperledger.aries.crypto.tink.PSMSPrivateKey"
)

// common errors.
var (
	errInvalidPSMSSignerKey      = errors.New("psms_signer_key_manager: invalid key")
	errInvalidPSMSignerKeyFormat = errors.New("psms_signer_key_manager: invalid key format")
)

// psmsSignerKeyManager is an implementation of KeyManager interface for PSMS signatures/proofs.
// It generates new PSMSPrivateKeys and produces new instances of PSMSSign subtle.
type psmsSignerKeyManager struct{}

// newPSMSSignerKeyManager creates a new psmsSignerKeyManager.
func newPSMSSignerKeyManager() *psmsSignerKeyManager {
	return new(psmsSignerKeyManager)
}

/* //TODO UMU (crypto) For "proper" integration with tink
// Primitive creates an PSMS Signer subtle for the given serialized PSMSPrivateKey proto.
func (km *psmsSignerKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidPSMSSignerKey
	}

	key := new(psmspb.PSMSPrivateKey)

	err := proto.Unmarshal(serializedKey, key)
	if err != nil {
		return nil, fmt.Errorf(errInvalidPSMSSignerKey.Error()+": invalid proto: %w", err)
	}

	err = km.validateKey(key)
	if err != nil {
		return nil, fmt.Errorf(errInvalidPSMSSignerKey.Error()+": %w", err)
	}

	return psmssubtle.NewBLS12381G2Signer(key.KeyValue), nil
}

// NewKey creates a new key according to the specification of PSMSPrivateKey format.
func (km *psmsSignerKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	if len(serializedKeyFormat) == 0 {
		return nil, errInvalidPSMSSignerKeyFormat
	}

	keyFormat := new(psmspb.PSMSKeyFormat)

	err := proto.Unmarshal(serializedKeyFormat, keyFormat)
	if err != nil {
		return nil, fmt.Errorf(errInvalidPSMSSignerKeyFormat.Error()+": invalid proto: %w", err)
	}

	err = validateKeyFormat(keyFormat)
	if err != nil {
		return nil, fmt.Errorf(errInvalidPSMSSignerKeyFormat.Error()+": %w", err)
	}

	var (
		pubKey  *psms12381g2pub.PublicKey
		privKey *psms12381g2pub.PrivateKey
	)

	// Since psms+ in aries-framework-go only supports BLS12-381 curve on G2, we create keys of this curve and
	// group only. PSMS+ keys with other curves/group field can be added later if needed.
	if keyFormat.Params.Group == psmspb.GroupField_G2 && keyFormat.Params.Curve == psmspb.PSMSCurveType_BLS12_381 {
		seed := make([]byte, 32)

		_, err = rand.Read(seed)
		if err != nil {
			return nil, err
		}

		hFunc := subtle.GetHashFunc(keyFormat.Params.HashType.String())

		pubKey, privKey, err = psms12381g2pub.GenerateKeyPair(hFunc, seed)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errInvalidPSMSSignerKeyFormat
	}

	pubKeyBytes, err := pubKey.Marshal()
	if err != nil {
		return nil, err
	}

	privKeyBytes, err := privKey.Marshal()
	if err != nil {
		return nil, err
	}

	return &psmspb.PSMSPrivateKey{
		Version:  psmsSignerKeyVersion,
		KeyValue: privKeyBytes,
		PublicKey: &psmspb.PSMSPublicKey{
			Version:  psmsSignerKeyVersion,
			Params:   keyFormat.Params,
			KeyValue: pubKeyBytes,
		},
	}, nil
}

// NewKeyData creates a new KeyData according to the specification of ECDHESPrivateKey Format.
// It should be used solely by the key management API.
func (km *psmsSignerKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	key, err := km.NewKey(serializedKeyFormat)
	if err != nil {
		return nil, err
	}

	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, fmt.Errorf("psms_signer_key_manager: Proto.Marshal failed: %w", err)
	}

	return &tinkpb.KeyData{
		TypeUrl:         psmsSignerKeyTypeURL,
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_ASYMMETRIC_PRIVATE,
	}, nil
}

// PublicKeyData returns the enclosed public key data of serializedPrivKey.
func (km *psmsSignerKeyManager) PublicKeyData(serializedPrivKey []byte) (*tinkpb.KeyData, error) {
	privKey := new(psmspb.PSMSPrivateKey)

	err := proto.Unmarshal(serializedPrivKey, privKey)
	if err != nil {
		return nil, errInvalidPSMSSignerKey
	}

	serializedPubKey, err := proto.Marshal(privKey.PublicKey)
	if err != nil {
		return nil, errInvalidPSMSSignerKey
	}

	return &tinkpb.KeyData{
		TypeUrl:         PsmsSignerKeyTypeURL,
		Value:           serializedPubKey,
		KeyMaterialType: tinkpb.KeyData_ASYMMETRIC_PUBLIC,
	}, nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *psmsSignerKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == PsmsSignerKeyTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *psmsSignerKeyManager) TypeURL() string {
	return PsmsSignerKeyTypeURL
}

// validateKey validates the given ECDHPrivateKey and returns the KW curve.
func (km *psmsSignerKeyManager) validateKey(key *psmspb.PSMSPrivateKey) error {
	err := keyset.ValidateKeyVersion(key.Version, psmsSignerKeyVersion)
	if err != nil {
		return fmt.Errorf("psms_signer_key_manager: invalid key: %w", err)
	}

	return validateKeyParams(key.PublicKey.Params)
}

// validateKeyFormat validates the given PSMS curve and Group field.
func validateKeyFormat(format *psmspb.PSMSKeyFormat) error {
	return validateKeyParams(format.Params)
}

func validateKeyParams(params *psmspb.PSMSParams) error {
	switch params.Curve {
	case psmspb.PSMSCurveType_BLS12_381:
	default:
		return fmt.Errorf("bad curve '%s'", params.Curve)
	}

	switch params.Group {
	case psmspb.GroupField_G1, psmspb.GroupField_G2:
	default:
		return fmt.Errorf("bad group field '%s'", params.Group)
	}

	switch params.HashType {
	case commonpb.HashType_SHA256, commonpb.HashType_SHA384, commonpb.HashType_SHA512:
	default:
		return fmt.Errorf("unsupported hash type '%s'", params.HashType)
	}

	return nil
}
*/
