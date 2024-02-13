/**
 * @file g3.h
 * @author Jesús García Rodríguez
 * @brief File with method definitions for the (field subgroup) group 3 on a cryptographic pairing
 *
 * This file defines the type and methods for manipulating elements form the third group (multiplicative notation) of a
 * cryptographic pairing, including serialization, multiplication...
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef G3_H
#define G3_H

#include <Zp.h> 

/**
 * Encapsulated declaration of type G3, which represents elements from G3 of the pairing
 */
typedef struct G3Impl G3;

/**
 * @brief Generates a copy of the unity. Can be safely modified . Has to be freed after usage
 * 
 * @return 
 */
G3* g3One();

/**
 * @brief Size of byte representation of G3 element (when computing g1ToBytes)
 */
int g3ByteSize();

/**
 * @brief Represent G3 element as byte array. Array is assumed to be big enough for copy
 * 
 * @param res Byte array where it will be copied
 * @param a G3 element
 */
void g3ToBytes(char *res, const G3 *a);

/**
 * @brief Multiplicative operation 
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g3Mul(G3* a, const G3* b);

/**
 * @brief Exponentiation a^b
 * 
 * @param a EC point, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g3Exp(G3* a, const Zp* b);

/**	@brief Tests for equality of two G3 elements
 *
	@param a instance to be compared
	@param b instance to be compared
	@return 1 if a=b, else returns 0
 */
int g3equals(const G3* a, const G3* b);

/**
 * @brief Output to stdout G3 (for testing)
 * 
 * @param e 
 */
void g3Print(const G3* e);


/**
 * @brief Destroy Zp element, freeing memory (from malloc etc.) as needed. Depending on internal structure might be equivalent to free(e)
 * 
 * @param e 
 */
void g3Free(G3* e);

#endif 