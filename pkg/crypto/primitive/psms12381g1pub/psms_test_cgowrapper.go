package psms12381g1pub

/*
#cgo CFLAGS: -I./clib/include -I./clib/include/ecwrapper
#cgo LDFLAGS: -L./clib/ -ldpabc_psms_bundled
#include <Dpabc.h>
#include <Zp.h>
#include <stdlib.h>

*/
import "C"
import "fmt"

func TestPsmsWrapperIndividualFlowCgoWrapper() (error, error, error, error, []byte, error) {
	seed := []byte("randomSeedForTestingWrapper")
	messagesBytes := [][]byte{
		[]byte("epoch0123476386422684282441741263116462036432461864386264"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}
	nattr := len(messagesBytes) - 1

	psms := New()
	pubKey, privKey, err1 := psms.GenerateKeyPair(seed, nattr)
	if err1 != nil {
		return err1, nil, nil, nil, nil, nil
	}
	privKeyBytes := MarshalPrivateKey(privKey)

	signatureBytes, err2 := psms.Sign(messagesBytes, privKeyBytes)
	if err2 != nil {
		return err1, err2, nil, nil, nil, nil
	}
	pubKeyBytes := MarshalPublicKey(pubKey)
	err3 := psms.Verify([][]byte{messagesBytes[0], messagesBytes[1], messagesBytes[2], messagesBytes[3][1:]}, signatureBytes, pubKeyBytes)
	if err3 == nil {
		return err1, err2, fmt.Errorf("error verify should have failed"), nil, nil, nil
	}
	err3 = psms.Verify(messagesBytes, signatureBytes, pubKeyBytes)
	if err3 != nil {
		return err1, err2, err3, nil, nil, nil
	}
	nonce := []byte("nonce")
	revealedIndexes := []int{0, 3, 1}
	proofBytes, err4 := psms.DeriveProof(messagesBytes, signatureBytes, nonce, pubKeyBytes, revealedIndexes)
	if err4 != nil {
		return err1, err2, err3, err4, nil, nil
	}
	revealedMessages := make([][]byte, len(revealedIndexes))
	for i, ind := range revealedIndexes {
		revealedMessages[i] = messagesBytes[ind]
	}
	return err1, err2, err3, err4, proofBytes, psms.VerifyProof(revealedMessages, proofBytes, nonce, pubKeyBytes)
}

func TestPsmsWrapperDistributedFlowCgoWrapper() (error, error, error, error, error, error, []byte, error) {
	seeds := [][]byte{
		[]byte("randomSeedForTestingWrapperOtherSeed"),
		[]byte("secondrandomseedaaaaaaaa"),
		[]byte("thirdRandomSeeedaeasdpf"),
		[]byte("FourBeerCuatroQuatre"),
	}
	nattr := 3
	nSigners := 2
	psms := New()
	sks := make([]*PrivateKey, nSigners)
	pks := make([]*PublicKey, nSigners)
	sksSerial := make([][]byte, nSigners)
	pksSerial := make([][]byte, nSigners)
	var err1 error
	for i := 0; i < nSigners; i++ {
		pks[i], sks[i], err1 = psms.GenerateKeyPair(seeds[i], nattr)
		if err1 != nil {
			return err1, nil, nil, nil, nil, nil, nil, nil
		}
		sksSerial[i] = MarshalPrivateKey(sks[i])
		pksSerial[i] = MarshalPublicKey(pks[i])
	}
	messagesBytes := [][]byte{
		[]byte("epoch0123476386422684282468216462036432461864386264"),
		[]byte("message2"),
		[]byte("message3"),
		[]byte("message4"),
	}
	signatureBytes := make([][]byte, nSigners)
	var err2 error
	for i := 0; i < nSigners; i++ {
		signatureBytes[i], err2 = psms.Sign(messagesBytes, sksSerial[i])
		if err2 != nil {
			return err1, err2, nil, nil, nil, nil, nil, nil
		}
	}
	aggKey, err3 := psms.AggregatePublicKeys(pksSerial)
	if err3 != nil {
		return err1, err2, err3, nil, nil, nil, nil, nil
	}
	combinedSignBytes, err4 := psms.AggregateSignatures(pksSerial, signatureBytes)
	if err4 != nil {
		return err1, err2, err3, err4, nil, nil, nil, nil
	}

	aggKeyBytes := MarshalPublicKey(aggKey)
	err5 := psms.Verify(messagesBytes, combinedSignBytes, aggKeyBytes)
	if err5 != nil {
		return err1, err2, err3, err4, err5, nil, nil, nil
	}
	nonce := []byte("nonce")
	revealedIndexes := []int{0, 3, 1}
	proofBytes, err6 := psms.DeriveProof(messagesBytes, combinedSignBytes, nonce, aggKeyBytes, revealedIndexes)
	if err6 != nil {
		return err1, err2, err3, err4, err5, err6, nil, nil
	}
	revealedMessages := make([][]byte, len(revealedIndexes))
	for i, ind := range revealedIndexes {
		revealedMessages[i] = messagesBytes[ind]
	}
	return err1, err2, err3, err4, err5, err6, proofBytes, psms.VerifyProof(revealedMessages, proofBytes, nonce, aggKeyBytes)
}
