/**
 * @file Dpabc_types.h
 * @author Jesús García Rodríguez
 * @brief Type definitions (encapsulations) for dp-ABC elements and freeing methods
 *
 * This file defines the encapsulated types that will represent the elements srelated to a dp-ABC: keys, 
 * signatures, zero-knowledge tokens...
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef DPABC_TYPES_H
#define DPABC_TYPES_H

#include <Zp.h>
#include <g1.h>
#include <g2.h>
#include <g3.h>

/**
 * @brief Encapsulated definition of the representation of a secret key of the scheme
 * 
 */
typedef struct secretKeyImpl secretKey;

/**
 * @brief Encapsulated definition of the representation of a public key of the scheme
 * 
 */
typedef struct publicKeyImpl publicKey;

/**
 * @brief Encapsulated definition of the representation of a signature of the scheme
 * 
 */
typedef struct signatureImpl signature;

/**
 * @brief Encapsulated definition of the representation of a zero knowledge presentation token of the scheme
 * 
 */
typedef struct zkTokenImpl zkToken;

/**
 * @brief Get public key corresponding to secret key sk 
 */ 
publicKey *dpabcSkToPk(const secretKey *sk);

/**
 * @brief Free memory from publicKey (and all its elements)
 * @param pk 
 */
void dpabcPkFree(publicKey *pk);

/**
 * @brief Size of byte representation of a public key pk (when serializing with dpabcPkToBytes)
 */
int dpabcPkByteSize(const publicKey *pk);

/**
 * @brief Represent PublicKey element as byte array. Array is assumed to be big enough for copy
 * @param res Byte array where it will be copied
 * @param a Public key
 */
void dpabcPkToBytes(char *res, const publicKey *pk);

/**
 * @brief Generate public key from bytes (previously serialized). Has to be freed after usage.
 * @param bytes Byte array of serialized element (assumed to have been correctly generated with dpabcPkToBytes())
 */
publicKey * dpabcPkFromBytes(const char *bytes);

/**
 * @brief Check equiality of public keys. Might be useful in some scenarios for checking validity of keys
 * @return 1 if equals (i.e., all elements are equal), 0 if not
 */
int dpabcPkEquals(const publicKey *pk1, const publicKey *pk2);

/**
 * @brief Free memory from secret key (and all its elements)
 * 
 * @param sk 
 */
void dpabcSkFree(secretKey *sk);

/**
 * @brief Size of byte representation of a secret key sk (when serializing with dpabcSkToBytes)
 */
int dpabcSkByteSize(const secretKey *sk);

/**
 * @brief Represent secret key as byte array. Array is assumed to be big enough for copy
 * @param res Byte array where it will be copied
 * @param a Public key
 */
void dpabcSkToBytes(char *res, const secretKey *sk);

/**
 * @brief Generate secret key from bytes (previously serialized). Has to be freed after usage.
 * @param bytes Byte array of serialized element (assumed to have been correctly generated with dpabcSkToBytes())
 */
secretKey * dpabcSkFromBytes(const char *bytes);

/**
 * @brief Free memory from signature (and all its elements)
 * 
 * @param sign 
 */
void dpabcSignFree(signature *sign);

/**
 * @brief Size of byte representation of a signature (when serializing with dpabcSignToBytes)
 */
int dpabcSignByteSize();

/**
 * @brief Represent signature as byte array. Array is assumed to be big enough for copy
 * @param res Byte array where it will be copied
 * @param a Public key
 */
void dpabcSignToBytes(char *res, const signature *sig);

/**
 * @brief Generate signature from bytes (previously serialized). Has to be freed after usage.
 * @param bytes Byte array of serialized element (assumed to have been correctly generated with dpabcSignToBytes())
 */
signature * dpabcSignFromBytes(const char *bytes);

/**
 * @brief Free memory from secret zk token (and all its elements)
 * 
 * @param zk 
 */
void dpabcZkFree(zkToken *zk);

/**
 * @brief Size of byte representation of a zero knowledge token zk(when serializing with dpabcZkToBytes)
 */
int dpabcZkByteSize(zkToken *zk);

/**
 * @brief Represent zkToken as byte array. Array is assumed to be big enough for copy
 * @param res Byte array where it will be copied
 * @param a Public key
 */
void dpabcZkToBytes(char *res, const zkToken *zk);

/**
 * @brief Generate zkToken from bytes (previously serialized). Has to be freed after usage.
 * @param bytes Byte array of serialized element (assumed to have been correctly generated with dpabcZkToBytes())
 */
zkToken * dpabcZkFromBytes(const char *bytes);


#endif 