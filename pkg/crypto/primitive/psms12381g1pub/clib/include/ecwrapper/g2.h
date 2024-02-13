/**
 * @file g2.h
 * @author Jesús García Rodríguez
 * @brief File with method definitions for the (EC) group 2 on a cryptographic pairing
 *
 * This file defines the type and methods for manipulating elements form the second group (additive notation) of a
 * cryptographic pairing, including generators, serialization, addition, multiplication...
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef G2_H
#define G2_H

#include <Zp.h> 

/**
 * Encapsulated declaration of type G2, which represents elements from G2 of the pairing
 */
typedef struct G2Impl G2;


/**
 * @brief Generates a copy of the group generator. Can be safely modified . Has to be freed after usage
 * 
 * @return 
 */
G2* g2Generator();


/**
 * @brief Generates a copy of group identity (addition operation). Can be safely modified. Has to be freed after usage
 * 
 * @return 
 */
G2* g2Identity();

/**
 * @brief Size of byte representation of G2 element (when computing g1ToBytes)
 */
int g2ByteSize();

/**
 * @brief Represent G2 element as byte array. Array is assumed to be big enough for copy
 * 
 * @param res Byte array where it will be copied
 * @param a G2 element
 */
void g2ToBytes(char *res, const G2 *a);

/**
 * @brief Generates a hashed G2 point from bytes. Has to be freed after usage
 * 
 * @param bytes Byte array 
 * @param n Number of bytes
 * @return 
 */
G2* hashToG2(const char *bytes, int n);

/**
 * @brief Generates a G2 element from bytes (previously serialized). Has to be freed after usage
 * 
 * @param bytes Byte array of serialized element (assumed to be the correct length, i.e. g2ByteSyze())
 */
G2 * g2FromBytes(const char *bytes);

/**
 * @brief Addition within the curve a+b
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g2Add(G2* a, const G2* b);

/**
 * @brief Subtraction within the curve a-b
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g2Sub(G2* a, const G2* b);

/**
 * @brief Multiplication [b]a
 * 
 * @param a EC point, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g2Mul(G2* a, const Zp* b);

/**
 * @brief Inverse multiplication [-b]a
 * 
 * @param a EC point, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g2InvMul(G2* a, const Zp* b);

/**
 * @brief Check if element is identity
 * 
 * @param a EC point
 * @return 1 if identity, 0 if not
 */
int g2IsIdentity(const G2* a);

/**
 * @brief Checks a=b in G2
 * 
 * @param a First element
 * @param b Second element
 * @returns 1 if equal 0 if not
 */
int g2Equals(const G2* a, const G2* b);

/**
 * @brief Copies G2 element. Result has to be freed after usage
 * 
 * @param a G2 element to be copied
 */
G2* g2Copy(const G2* a);

/**
 * @brief Overwrite a value to b (no new memory is reserved)
 * 
 * @param a G2 element, after method a=b
 * @param b G2 element, not modified
 */
void g2CopyValue(G2* a, const G2* b);

/**
 * @brief Output to stdout G2 (for testing)
 * 
 * @param e 
 */
void g2Print(const G2* e);


/**
 * @brief Destroy G2 element, freeing memory (from malloc etc.) as needed. Depending on internal structure might be equivalent to free(e)
 * 
 * @param e 
 */
void g2Free(G2* e);

#endif 