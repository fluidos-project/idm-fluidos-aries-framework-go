package psms

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	tinkpb "github.com/google/tink/go/proto/tink_go_proto"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/primitive/psms12381g1pub"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/jwkkid"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

const (
	ONLY_PUB = iota
	ONLY_PRIV
	BOTH_KEYS
)

type MockHandle struct {
	PrivateKey []byte
	PublicKey  []byte
	Nattr      uint16
}

func NewMockHandle(kt *tinkpb.KeyTemplate) (*MockHandle, error) {
	if kt.TypeUrl != PsmsSignerKeyTypeURL {
		return nil, fmt.Errorf("called mock handle with wrong keytemplate %s", kt)
	}
	seed := make([]byte, 32)
	n, err := rand.Read(seed)
	if err != nil || n != 32 {
		return nil, fmt.Errorf("error generating random value %w", err)
	}
	return NewMockHandleWithSeed(kt, seed)
}

func NewMockHandleWithSeed(kt *tinkpb.KeyTemplate, seed []byte) (*MockHandle, error) {
	if kt.TypeUrl != PsmsSignerKeyTypeURL {
		return nil, fmt.Errorf("called mock handle with wrong keytemplate %s", kt)
	}
	nattr := binary.BigEndian.Uint16(kt.Value)
	pk, sk, err := psms12381g1pub.New().GenerateKeyPair(seed, int(nattr))
	if err != nil {
		return nil, err
	}
	skBytes, err := sk.Marshal()
	if err != nil {
		return nil, err
	}
	pkBytes, err := pk.Marshal()
	if err != nil {
		return nil, err
	}
	return &MockHandle{
		PrivateKey: skBytes,
		PublicKey:  pkBytes,
		Nattr:      nattr,
	}, nil
}

func createKIDWrapper(keyBytes []byte, kt kms.KeyType) (string, error) {
	return jwkkid.CreateKID(keyBytes, kt)
}
func GenerateKIDWrapper(kh *MockHandle, kt kms.KeyType) (string, error) {
	keyBytes, _, err := ExportPubKeyBytes(kh)
	if err != nil {
		return "", fmt.Errorf("GenerateKID: failed to export public key: %w", err)
	}

	return createKIDWrapper(keyBytes, kt)
}

func (kh *MockHandle) Public() (*MockHandle, error) {
	return &MockHandle{
		PrivateKey: nil,
		PublicKey:  kh.PublicKey,
		Nattr:      kh.Nattr,
	}, nil
}

func (kh *MockHandle) ToBytes() ([]byte, error) {
	if kh.PrivateKey == nil && kh.PublicKey == nil {
		return nil, fmt.Errorf("error serializing, both keys are nil")
	}
	kht := BOTH_KEYS
	if kh.PrivateKey == nil {
		kht = ONLY_PUB
	} else if kh.PublicKey == nil {
		kht = ONLY_PRIV
	}
	var bytes []byte
	switch kht {
	case BOTH_KEYS:
		size := len(kh.PrivateKey) + len(kh.PublicKey) + 7
		bytes = make([]byte, size)
		offset := 0
		bytes[offset] = uint8(kht)
		offset += 1
		binary.BigEndian.PutUint16(bytes[offset:], kh.Nattr)
		offset += 2
		binary.BigEndian.PutUint32(bytes[offset:], uint32(len(kh.PrivateKey)))
		offset += 4
		offset += copy(bytes[offset:], kh.PrivateKey)
		offset += copy(bytes[offset:], kh.PublicKey)
		if offset != size {
			return nil, fmt.Errorf("error serializing, incorrect number of bytes written")
		}
	case ONLY_PUB:
		size := len(kh.PublicKey) + 3
		bytes = make([]byte, size)
		offset := 0
		bytes[offset] = uint8(kht)
		offset += 1
		binary.BigEndian.PutUint16(bytes[offset:], kh.Nattr)
		offset += 2
		offset += copy(bytes[offset:], kh.PublicKey)
		if offset != size {
			return nil, fmt.Errorf("error serializing, incorrect number of bytes written")
		}
	case ONLY_PRIV:
		size := len(kh.PrivateKey) + 3
		bytes = make([]byte, size)
		offset := 0
		bytes[offset] = uint8(kht)
		offset += 1
		binary.BigEndian.PutUint16(bytes[offset:], kh.Nattr)
		offset += 2
		offset += copy(bytes[offset:], kh.PrivateKey)
		if offset != size {
			return nil, fmt.Errorf("error serializing, incorrect number of bytes written")
		}
	}
	return bytes, nil
}

func MockHandleFromBytes(val []byte) (*MockHandle, error) {
	if len(val) == 0 {
		return nil, fmt.Errorf("error deserializing, no bytes")
	}
	kht := val[0]
	offset := 1
	switch kht {
	case BOTH_KEYS:
		nattr := binary.BigEndian.Uint16(val[offset : offset+2])
		offset += 2
		skLen := int(binary.BigEndian.Uint32(val[offset : offset+4]))
		offset += 4
		sk := append([]byte(nil), val[offset:offset+skLen]...)
		offset += skLen
		pk := append([]byte(nil), val[offset:]...)
		return &MockHandle{
			PrivateKey: sk,
			PublicKey:  pk,
			Nattr:      nattr,
		}, nil
	case ONLY_PUB:
		nattr := binary.BigEndian.Uint16(val[offset : offset+2])
		offset += 2
		pk := append([]byte(nil), val[offset:]...)
		return &MockHandle{
			PrivateKey: nil,
			PublicKey:  pk,
			Nattr:      nattr,
		}, nil
	case ONLY_PRIV:
		nattr := binary.BigEndian.Uint16(val[offset : offset+2])
		offset += 2
		sk := append([]byte(nil), val[offset:]...)
		return &MockHandle{
			PrivateKey: sk,
			PublicKey:  nil,
			Nattr:      nattr,
		}, nil
	default:
		return nil, fmt.Errorf("error deserializing, unrecognized kht")
	}
}

func ExportPubKeyBytes(kh *MockHandle) ([]byte, kms.KeyType, error) {
	return kh.PublicKey, kms.BLS12381G1Type, nil
}

func MockPublicKeyBytesToHandle(key []byte, kt *tinkpb.KeyTemplate) (*MockHandle, error) {
	if kt.TypeUrl == PsmsSignerKeyTypeURL {
		return &MockHandle{
			PrivateKey: nil,
			PublicKey:  key,
			Nattr:      binary.BigEndian.Uint16(kt.Value),
		}, nil
	}
	return nil, fmt.Errorf("wrong key type/template %s", kt)
}
