package psms12381g1pub

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -I./clib/include -I./clib/include/ecwrapper
#cgo LDFLAGS: -L./clib/ -ldpabc_psms_bundled
#include <Dpabc.h>
#include <Zp.h>
#include <stdlib.h>
#include <stdio.h>


//Malloc methods for pks,signs... arrays might be summarized in the following method, but not much is gained and forces casting
//void** mallocArrayOfPointers(int nmessages){
//	return (void **) malloc(nmessages*sizeof(void*));
//}

Zp** mallocZpMessages(int nmessages){
	return (Zp **) malloc(nmessages*sizeof(Zp*));
}

publicKey** mallocPks(int nkeys){
	return (publicKey **) malloc(nkeys*sizeof(publicKey*));
}

signature** mallocSignatures(int nsigns){
	return (signature **) malloc(nsigns*sizeof(signature*));
}

Zp* getZpval(Zp** array,int i){
	return array[i];
}

Zp** getZpvalsfrom(Zp** array, int i){
	return &(array[i]);
}

void freeAttrs(Zp** array,int len){
	for(int i=0;i<len;i++)
		zpFree(array[i]);
	free(array);
}

//Assign methods for pks,signs... arrays might be summarized in the following method, but not much is gained and forces casting
//void assignValArray(void** array,void* val,int i){
//	array[i]=val;
//}
void assignZpArray(Zp** array,Zp* val,int i){
	array[i]=val;
}

void assignPkArray(publicKey** array,publicKey* val,int i){
	array[i]=val;
}

void assignSignatureArray(signature** array,signature* val,int i){
	array[i]=val;
}

void freePkIndex(publicKey** array,int index){
	dpabcPkFree(array[index]);
}

void freeSignatureIndex(signature** array,int index){
	dpabcSignFree(array[index]);
}
*/
import "C"

//UMU Maybe organize in multiple files

// PSMSG1Pub defines PSMS multi-signature scheme where public key is a point in the field of G1.
type PSMSG1Pub struct{}

// New creates a new PSMSG1Pub.
func New() *PSMSG1Pub {
	return &PSMSG1Pub{}
}

// PublicKey defines PSMS Public Key Wrapper.
type PublicKey struct {
	key    *C.publicKey
	nattrs uint16
}

// PrivateKey defines PSMS Private Key Wrapper.
type PrivateKey struct {
	key    *C.secretKey
	nattrs uint16
}

func (k *PrivateKey) Marshal() ([]byte, error) {
	return MarshalPrivateKey(k), nil
}

func (pk *PublicKey) Marshal() ([]byte, error) {
	return MarshalPublicKey(pk), nil
}

func (pk *PublicKey) Nattrs() uint16 {
	return pk.nattrs
}

func (k *PrivateKey) Public() ([]byte, error) {
	pk := PublicKey{
		nattrs: k.nattrs,
		key:    C.dpabcSkToPk(k.key)}
	return pk.Marshal()
}

func (k *PrivateKey) PublicKey() *PublicKey {
	pk := PublicKey{
		nattrs: k.nattrs,
		key:    C.dpabcSkToPk(k.key)}
	return &pk
}

func (k *PrivateKey) Nattrs() uint16 {
	return k.nattrs
}

func (psms *PSMSG1Pub) GenerateKeyPair(seed []byte, nattr int) (*PublicKey, *PrivateKey, error) {
	C.seedRng((*C.char)(unsafe.Pointer(&seed[0])), C.int(len(seed)))
	C.changeNattr(C.int(nattr))
	var pk *C.publicKey
	var sk *C.secretKey
	C.keyGen(&sk, &pk)
	return &PublicKey{
			key:    pk,
			nattrs: uint16(nattr),
		}, &PrivateKey{
			key:    sk,
			nattrs: uint16(nattr),
		},
		nil
}

func (psms *PSMSG1Pub) AggregatePublicKeys(pubKeyBytes [][]byte) (*PublicKey, error) {
	pks := C.mallocPks(C.int(len(pubKeyBytes)))
	defer C.free(unsafe.Pointer(pks))
	var nattr uint16
	for i := range pubKeyBytes {
		pubKey, err := UnmarshalPublicKey(pubKeyBytes[i])
		if err != nil {
			return nil, fmt.Errorf("unmarshal public key: %w", err)
		}
		if i != 0 && nattr != pubKey.nattrs {
			return nil, fmt.Errorf("number of attributes does not match: %d,%d,%d", nattr, pubKey.nattrs, i)
		}
		nattr = pubKey.nattrs
		C.assignPkArray(pks, pubKey.key, C.int(i))
		defer C.freePkIndex(pks, C.int(i))
	}
	res := C.keyAggr(pks, C.int(len(pubKeyBytes)))
	return &PublicKey{
		key:    res,
		nattrs: uint16(nattr),
	}, nil
}

// Sign signs the one or more messages (first message is assumed to be epoch). Messages must be serialized forms of elements
// used by this signature scheme. using the serialized private key. Result is serialized.
func (psms *PSMSG1Pub) Sign(messages [][]byte, privKeyBytes []byte) ([]byte, error) {
	privKey, err := UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal private key: %w", err)
	}
	defer C.dpabcSkFree(privKey.key)
	if len(messages) <= 1 {
		return nil, errors.New("messages are not defined")
	}
	if uint16(len(messages)-1) != privKey.nattrs {
		return nil, fmt.Errorf("mismatched number of attributes for pk:%d,messages:%d", privKey.nattrs, len(messages)-1)
	}
	messagesZp := C.mallocZpMessages(C.int(len(messages)))
	for i := range messages {
		C.assignZpArray(messagesZp, ParseMessage(messages[i]), C.int(i))
	}
	defer C.freeAttrs(messagesZp, C.int(len(messages)))
	//TODO UMU (crypto/verifiable) What to do with epoch? For now, assumed to always be first attr
	signature := C.sign(privKey.key, C.getZpval(messagesZp, 0), C.getZpvalsfrom(messagesZp, 1))
	return MarshalSignature(signature), nil
}

func (psms *PSMSG1Pub) AggregateSignatures(pubKeyBytes, signaturesBytes [][]byte) ([]byte, error) {
	if len(pubKeyBytes) != len(signaturesBytes) {
		return nil, fmt.Errorf("number of signaturesBytes and keys must be the same")
	}
	pks := C.mallocPks(C.int(len(pubKeyBytes)))
	defer C.free(unsafe.Pointer(pks))
	var nattr uint16
	for i := range pubKeyBytes {
		pubKey, err := UnmarshalPublicKey(pubKeyBytes[i])
		if err != nil {
			return nil, fmt.Errorf("unmarshal public key: %w", err)
		}
		if i != 0 && nattr != pubKey.nattrs {
			return nil, fmt.Errorf("number of attributes does not match: %d,%d,%d", nattr, pubKey.nattrs, i)
		}
		nattr = pubKey.nattrs
		C.assignPkArray(pks, pubKey.key, C.int(i))
		defer C.freePkIndex(pks, C.int(i))
	}
	signs := C.mallocSignatures(C.int(len(signaturesBytes)))
	defer C.free(unsafe.Pointer(signs))
	for i := range signaturesBytes {
		sign, err := UnmarshalSignature(signaturesBytes[i])
		if err != nil {
			return nil, fmt.Errorf("unmarshal signature: %w", err)
		}
		C.assignSignatureArray(signs, sign, C.int(i))
		defer C.freeSignatureIndex(signs, C.int(i))
		//Should be fine, want to free these when the aggregation method ends
	}
	res := C.combine(pks, signs, C.int(len(pubKeyBytes)))
	return MarshalSignature(res), nil
}

func (psms *PSMSG1Pub) Verify(messages [][]byte, sigBytes, pubKeyBytes []byte) error {
	pubKey, err := UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("unmarshal public key: %w", err)
	}
	defer C.dpabcPkFree(pubKey.key)
	sign, err := UnmarshalSignature(sigBytes)
	if err != nil {
		return fmt.Errorf("unmarshal signature: %w", err)
	}
	defer C.dpabcSignFree(sign)
	if len(messages) <= 1 {
		return errors.New("messages are not defined")
	}
	if uint16(len(messages)-1) != pubKey.nattrs {
		return fmt.Errorf("mismatched number of attributes for pk:%d,messages:%d", pubKey.nattrs, len(messages)-1)
	}
	messagesZp := C.mallocZpMessages(C.int(len(messages)))
	for i := range messages {
		C.assignZpArray(messagesZp, ParseMessage(messages[i]), C.int(i))
	}
	defer C.freeAttrs(messagesZp, C.int(len(messages)))
	if uint16(len(messages)-1) != pubKey.nattrs {
		return fmt.Errorf("mismatched number of attributes for pk:%d,messages:%d", pubKey.nattrs, len(messages)-1)
	}
	res := C.verify(pubKey.key, sign, C.getZpval(messagesZp, 0), C.getZpvalsfrom(messagesZp, 1))
	if res == 1 {
		return nil
	}
	return errors.New("invalid PSMS BLS12-381 signature")
}

func (psms *PSMSG1Pub) DeriveProof(messages [][]byte, sigBytes, nonce, pubKeyBytes []byte,
	revealedIndexes []int) ([]byte, error) {
	if len(nonce) == 0 {
		return nil, fmt.Errorf("nil nonce is not allowed")
	}
	pubKey, err := UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal public key: %w", err)
	}
	defer C.dpabcPkFree(pubKey.key)
	if len(messages) < len(revealedIndexes) {
		return nil, fmt.Errorf("invalid size: %d revealed indexes is larger than %d messages", len(revealedIndexes),
			len(messages))
	}
	translatedRevealIndexes, sizeTranslatedRevealIndexes, err := checkAndTranslateRevealedIndexes(revealedIndexes)
	if err != nil { //TODO UMU (crypto) For now, we assume revealed messages must include epoch (0), and the remaining indexes have to be "translated". Checking validity of epoch?
		return nil, fmt.Errorf("error in revealed indexes, %w", err)
	}
	defer C.free(unsafe.Pointer(translatedRevealIndexes))
	if uint16(len(messages)-1) != pubKey.nattrs {
		return nil, fmt.Errorf("mismatched number of attributes for pk:%d,messages:%d", pubKey.nattrs, len(messages)-1)
	}
	sign, err := UnmarshalSignature(sigBytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal signature: %w", err)
	}
	defer C.dpabcSignFree(sign)
	if len(messages) <= 1 {
		return nil, errors.New("messages are not defined")
	}
	messagesZp := C.mallocZpMessages(C.int(len(messages)))
	for i := range messages {
		C.assignZpArray(messagesZp, ParseMessage(messages[i]), C.int(i))
	}
	defer C.freeAttrs(messagesZp, C.int(len(messages)))
	if psms.Verify(messages, sigBytes, pubKeyBytes) != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}
	zktok := C.presentZkToken(pubKey.key, sign, C.getZpval(messagesZp, 0), C.getZpvalsfrom(messagesZp, 1), translatedRevealIndexes,
		sizeTranslatedRevealIndexes, (*C.char)(unsafe.Pointer(&nonce[0])), C.int(len(nonce)))
	defer C.dpabcZkFree(zktok)
	zkPayload := newZkTokPayload(revealedIndexes, zktok)
	serialZkTok, err := zkPayload.toBytes()
	if err != nil {
		return nil, fmt.Errorf("serialize zktok: %w", err)
	}
	return serialZkTok, nil
}

func (psms *PSMSG1Pub) VerifyProof(messagesBytes [][]byte, proof, nonce, pubKeyBytes []byte) error {
	pubKey, err := UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("unmarshal public key: %w", err)
	}
	defer C.dpabcPkFree(pubKey.key)
	tok, err := parseZkTokPayload(proof)
	if err != nil {
		return fmt.Errorf("unmarshal tok: %w", err)
	}
	defer C.dpabcZkFree(tok.proof)
	revealIndexes, sizeTranslatedRevealIndexes, err := checkAndTranslateRevealedIndexes(tok.revealedIndexes)
	if err != nil {
		return fmt.Errorf("error in revealed indexes, %w", err)
	}
	defer C.free(unsafe.Pointer(revealIndexes))
	if len(messagesBytes) <= 1 {
		return errors.New("messages are not defined")
	}
	messagesZp := C.mallocZpMessages(C.int(len(messagesBytes)))
	for i := range messagesBytes {
		C.assignZpArray(messagesZp, ParseMessage(messagesBytes[i]), C.int(i))
	}
	defer C.freeAttrs(messagesZp, C.int(len(messagesBytes)))
	res := C.verifyZkToken(tok.proof, pubKey.key, C.getZpval(messagesZp, 0), C.getZpvalsfrom(messagesZp, 1),
		revealIndexes, sizeTranslatedRevealIndexes, (*C.char)(unsafe.Pointer(&nonce[0])), C.int(len(nonce)))
	if int(res) == 1 {
		return nil
	}
	return errors.New("invalid PSMS BLS12-381 zk proof")
}
