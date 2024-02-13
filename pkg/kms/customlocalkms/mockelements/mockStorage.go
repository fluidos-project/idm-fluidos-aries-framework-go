package mockelements

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/crypto/customtinkcrypto/primitive/psms"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

type MockStorage struct {
	storedKeys kms.Store
}

func NewMockStorage(p kms.Store) (*MockStorage, error) {
	return &MockStorage{
		storedKeys: p,
	}, nil
}

func (ms *MockStorage) StoreKeySet(kh *psms.MockHandle, kt kms.KeyType) (string, error) {
	var (
		kid string
		err error
	)
	if kt == kms.BLS12381G1Type {
		kid, err = psms.GenerateKIDWrapper(kh, kt)
		if err != nil {
			return "", fmt.Errorf("storeKeySet: failed to generate kid: %w", err)
		}
	} else {
		return "", fmt.Errorf("mockStoreKeySet: cannot call mock method with different key type")
	}
	_, err = ms.WriteToStore(kid, kh)
	return kid, err
}

func (ms *MockStorage) WriteToStore(id string, kh *psms.MockHandle) (bool, error) {
	existing := false
	if _, err := ms.storedKeys.Get(id); err == nil {
		existing = true
	}
	val, err := kh.ToBytes()
	if err != nil {
		return false, err
	}
	err = ms.storedKeys.Put(id, val)
	if err != nil {
		return false, err
	}
	return existing, nil
}

// Retrieve Returns the stored value if present, along with a boolean that indicates whether it was present or not
func (ms *MockStorage) Retrieve(id string) (*psms.MockHandle, bool) {
	val, err := ms.storedKeys.Get(id)
	if err != nil {
		return nil, false
	}
	dval, err := psms.MockHandleFromBytes(val)
	if err != nil {
		return nil, false
	}
	return dval, true
}

func (ms *MockStorage) Delete(id string) {
	ms.storedKeys.Delete(id)
}

func (ms *MockStorage) WriteImportedKey(key *psms12381g1pub.PrivateKey, kt kms.KeyType, opts ...kms.PrivateKeyOpts) (string, *psms.MockHandle, error) {
	if kt != kms.BLS12381G1Type {
		return "", nil, fmt.Errorf("import private PSMS key failed: invalid key type")
	}
	pubKey, err := key.Public()
	if err != nil {
		return "", nil, fmt.Errorf("import private PSMS key failed: cannot retrieve public")
	}
	prKeyBytes, err := key.Marshal()
	if err != nil {
		return "", nil, fmt.Errorf("import private PSMS key failed: error marshalling")
	}
	kh := &psms.MockHandle{
		PrivateKey: prKeyBytes,
		PublicKey:  pubKey,
		Nattr:      key.Nattrs(),
	}
	pOpts := kms.NewOpt()
	for _, opt := range opts {
		opt(pOpts)
	}
	var kid string
	if pOpts.KsID() != "" {
		if _, err := ms.storedKeys.Get(pOpts.KsID()); err == nil {
			return "", nil, fmt.Errorf("storeKeySet: key Id already existed %s", pOpts.KsID())
		}
		kid = pOpts.KsID()
	} else {
		kid, err = psms.GenerateKIDWrapper(kh, kt)
		if err != nil {
			return "", nil, fmt.Errorf("storeKeySet: failed to generate kid: %w", err)
		}
	}
	_, err = ms.WriteToStore(kid, kh)
	return kid, kh, err
}
