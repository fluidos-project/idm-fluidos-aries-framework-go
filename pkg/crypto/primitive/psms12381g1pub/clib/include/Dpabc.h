/**
 * @file Dpabc.h
 * @author Jesús García Rodríguez
 * 
 * @brief File with method definitions for a distributed p-ABC scheme.
 *
 * @details This file defines the methods related to a distributed privacy-preserving Attribute-Based Credential scheme, 
 * including key generation, public key aggregation, issuance (signing and combining), verification, 
 * and gereration/verification of zero-knowledge presentations. 
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef DPABC_H
#define DPABC_H

#ifndef NATTRINI
#define NATTRINI 4
#endif 

#include <Dpabc_types.h>

//TODO Avoid (should we?) "DoS" by segmentation faults, etc. (e.g., token says correct n but array is shorter)

//TODO Rethink const for methods. Issue: with underlying libraries not declaring const even if they are, causes warnings
// Issue with double pointers: forcing user to cast to const

/**
 * @brief Change the number of attributes considered for the operations of the scheme
 * 
 * @param n Number of attributes
 */
void changeNattr(int n); //TODO Think about this "initialization" (other options: macros/compile, just extra argument when generating keys...)

/**
 * @brief Seed the secure random number generator used in the scheme
 * 
 * @param seed Seed
 * @param n Seed length
 */
void seedRng(const char* seed, int n);

/**
 * @brief Generate secret and public key pair. They must be freed after usage
 * 
 * @param sk Pointer to a pointer to secret key, so it can be modified (i.e., after method *sk will point to the secret key, or null if something went wrong)
 * @param pk Pointer to a pointer to public key, so it can be modified (i.e., after method *pk will point to the pulic key, or null if something went wrong)
 */
void keyGen(secretKey **sk, publicKey **pk);


/**
 * @brief Aggregate public keys of nkeys signers, result is a new public key that must be freed after usage
 * 
 * @param pks Array of public keys (pointers)
 * @param nkeys Number of keys to aggregate
 * @return The aggregated public key (must be freed after usage), or null if something went wrong  
 */
publicKey* keyAggr(const publicKey *pks[], int nkeys);

/**
 * @brief Generate a signature over a set of attributes (the number of attributes is assumed to be the same as in the secret key) and epoch.
 * Note that the order of the attributes is crucial
 * 
 * @param sk Secret kye
 * @param epoch Epoch
 * @param attributes Attributes to be signed
 * @return signature* Generated signature (must be freed after usage), or null if something went wrong 
 */
signature* sign(const secretKey *sk, const Zp *epoch, const Zp *attributes[]);

/**
 * @brief Combine signatures generated with a set of signing keys. 
 * 
 * @param pks Public keys
 * @param signs Signatures (shares)
 * @param nkeys Number of pks/signatures
 * @return signature*  Resulting signature (must be freed after usage), or null if something went wrong 
 */
signature* combine(const publicKey *pks[], const signature *signs[], int nkeys);

/**
 * @brief Verify a signature over a set of attributes (the number of attributes is assumed to be the same as in the public key) and epoch with respect to a public key
 * Note that the order of the attributes is crucial
 * 
 * @param pk Public key
 * @param sign Signature
 * @param epoch Sigend epoch
 * @param attributes  Signed attributes
 * @return 1 if the signature is valid with respect to the public key, 0 otherwise 
 */
int verify(const publicKey *pk, const signature* sign, const Zp *epoch, const Zp *attributes[]);

/**
 * @brief Do a zero-konwledge proof to obtain a zero-knowledge token that reveals the attributes defined by their indexes (indexReveal)
 * Note that the order of the attributes is crucial
 * 
 * @param pk Public key corresponding to the signature
 * @param sign  Signature for which we are computing a zero-knowledge proof
 * @param epoch Epoch
 * @param attributes Signed attributes
 * @param indexReveal Indexes of the revealed attributes  (with respect to the whole set of attributes, starting from 0). Assumed to be in ascendent order
 * @param nIndexReveal Number of revealed attributes
 * @param message Message that will be signed for generating the zero-knowldege token
 * @param messageSize Size of the message to be signed
 * @return zkToken* Resulting zero-knowledge token (must be freed after usage), or null if something went wrong 
 */
zkToken* presentZkToken(const publicKey * pk, const signature *sign, const Zp *epoch, const Zp *attributes[], 
        const int indexReveal[], int nIndexReveal, const char *message, int messageSize);

/**
 * @brief Verify a zero-knowledge token that reveals the attributes defined by their indexes (indexReveal). Revealed attributes are assumed to be on ascendent order in regards to their indexes.
 * Note that the order of the attributes is crucial
 * 
 * @param token Zero-knowledge token
 * @param pk Public key
 * @param epoch Epoch
 * @param revealed Revealed attributes, assumed to be in ascendent order (e.g., if token reveals attributes 0, 3, and 6, revealed is assumed to contain them in positions 0, 1 and 2 respectively)
 * @param indexReveal Indexes of revealed attributes (with respect to the whole set of attributes, starting from 0). Assumed to be in ascendent order
 * @param nReveal Number of revealed attributes
 * @param message Message that was be signed for generating the zero-knowldege token
 * @param messageSize Size of the signed message 
 * @return int 
 */
int verifyZkToken(const zkToken *token, const publicKey * pk, const Zp *epoch, const Zp *revealed[],
        const int indexReveal[], int nReveal, const char *message, int messageSize);

/**
 * @brief Frees all necessary data associated to the scheme, e.g., rng
 * 
 */
void dpabcFreeStateData();


#endif 