package psms12381g1pub

/*
#cgo CFLAGS: -I./clib/include -I./clib/include/ecwrapper
#cgo LDFLAGS: -L./clib/ -ldpabc_psms_bundled
#include <Dpabc.h>
#include <Zp.h>
#include <stdlib.h>
int * mallocRevealedIndexes(int nreveal){
	return (int *) malloc(nreveal*sizeof(int));
}

int * assignIndex(int * array, int i, int index){
	array[i]=index;
}
*/
import "C"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"unsafe"
)

type zkTokPayload struct {
	messagesCount   int
	revealedIndexes []int
	proof           *C.zkToken
}

func MarshalPrivateKey(key *PrivateKey) []byte {
	res := make([]byte, int(C.dpabcSkByteSize(key.key))+2) //XXX UMU unit16 for number of attributes
	binary.BigEndian.PutUint16(res, key.nattrs)
	C.dpabcSkToBytes((*C.char)(unsafe.Pointer(&res[2])), key.key)
	return res

}

func UnmarshalPrivateKey(privKeyBytes []byte) (*PrivateKey, error) {
	return &PrivateKey{
		key:    C.dpabcSkFromBytes((*C.char)(unsafe.Pointer(&privKeyBytes[2]))),
		nattrs: uint16FromBytes(privKeyBytes[0:2]),
	}, nil
}

func MarshalPublicKey(key *PublicKey) []byte {
	res := make([]byte, int(C.dpabcPkByteSize(key.key))+2) //XXX UMU unit16 for number of attributes
	binary.BigEndian.PutUint16(res, key.nattrs)
	C.dpabcPkToBytes((*C.char)(unsafe.Pointer(&res[2])), key.key)
	return res

}

func UnmarshalPublicKey(pubKeyBytes []byte) (*PublicKey, error) {
	return &PublicKey{
		key:    C.dpabcPkFromBytes((*C.char)(unsafe.Pointer(&pubKeyBytes[2]))),
		nattrs: uint16FromBytes(pubKeyBytes[0:2]),
	}, nil
}

func MarshalSignature(signature *C.signature) []byte {
	sigBytes := C.malloc((C.size_t)(C.dpabcSignByteSize() * C.sizeof_char))
	C.dpabcSignToBytes((*C.char)(sigBytes), signature)
	return C.GoBytes(unsafe.Pointer(sigBytes), C.dpabcSignByteSize())

}

func UnmarshalSignature(sigBytes []byte) (*C.signature, error) {
	if int(C.dpabcSignByteSize()) != len(sigBytes) {
		return nil, fmt.Errorf("marshalled signature of wrong size %d", len(sigBytes))
	}
	return C.dpabcSignFromBytes((*C.char)(unsafe.Pointer(&sigBytes[0]))), nil
}

func ParseMessage(message []byte) *C.Zp {
	//TODO UMU (crypto) Messages are assumed to be on "plain" byte array form (for consistency). For now (as the rest of the framework) only revealing proofs actually supported because of this
	return C.hashToZp((*C.char)(unsafe.Pointer(&message[0])), C.int(len(message)))
}

func parseZkTokPayload(bytes []byte) (*zkTokPayload, error) {
	if len(bytes) < 2 {
		return nil, errors.New("invalid size of PoK payload")
	}

	messagesCount := int(uint16FromBytes(bytes[0:2]))
	offset := lenInBytes(messagesCount)

	if len(bytes) < offset {
		return nil, errors.New("invalid size of PoK payload")
	}

	revealed := bitvectorToIndexes(reverseBytes(bytes[2:offset]))

	tok := C.dpabcZkFromBytes((*C.char)(unsafe.Pointer(&bytes[offset])))

	return &zkTokPayload{
		messagesCount:   messagesCount,
		revealedIndexes: revealed,
		proof:           tok,
	}, nil
}

func newZkTokPayload(revealedIndexes []int, proof *C.zkToken) *zkTokPayload {
	return &zkTokPayload{
		messagesCount:   len(revealedIndexes),
		revealedIndexes: revealedIndexes,
		proof:           proof,
	}
}

func (p *zkTokPayload) toBytes() ([]byte, error) {
	bytes := make([]byte, lenInBytes(p.messagesCount)+int(C.dpabcZkByteSize(p.proof)))

	binary.BigEndian.PutUint16(bytes, uint16(p.messagesCount))
	offset := lenInBytes(p.messagesCount)
	bitvector := bytes[2:offset]

	for _, r := range p.revealedIndexes {
		idx := r / 8
		bit := r % 8

		if len(bitvector) <= idx {
			return nil, errors.New("invalid size of PoK payload")
		}

		bitvector[idx] |= 1 << bit
	}

	reverseBytes(bitvector)
	cSize := C.dpabcZkByteSize(p.proof)
	cserial := (*C.char)(C.malloc((C.size_t)(cSize * C.sizeof_char)))
	defer C.free(unsafe.Pointer(cserial))
	C.dpabcZkToBytes(cserial, p.proof)
	serialTok := C.GoBytes(unsafe.Pointer(cserial), cSize)
	copy(bytes[offset:], serialTok)
	return bytes, nil
}

func bitvectorToIndexes(data []byte) []int {
	revealedIndexes := make([]int, 0)
	scalar := 0

	for _, v := range data {
		remaining := 8

		for v > 0 {
			revealed := v & 1
			if revealed == 1 {
				revealedIndexes = append(revealedIndexes, scalar)
			}

			v >>= 1
			scalar++
			remaining--
		}

		scalar += remaining
	}

	return revealedIndexes
}

func uint16FromBytes(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func lenInBytes(messagesCount int) int {
	return 2 + (messagesCount / 8) + 1 //nolint:gomnd
}
func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func checkAndTranslateRevealedIndexes(indexes []int) (*C.int, C.int, error) {
	sort.Ints(indexes)
	uniqueIndexes := removeDuplicateInt(indexes)
	if len(uniqueIndexes) != len(indexes) {
		return nil, C.int(0), fmt.Errorf("invalid token, repeated reveal indexes")
	}
	if uniqueIndexes[0] != 0 {
		return nil, C.int(0), fmt.Errorf("invalid token, does not reveal epoch")
	}
	res := C.mallocRevealedIndexes(C.int(len(indexes) - 1))
	for i := 0; i < len(indexes)-1; i++ {
		C.assignIndex(res, C.int(i), C.int(indexes[i+1]-1))
	}
	return res, C.int(len(indexes) - 1), nil
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
