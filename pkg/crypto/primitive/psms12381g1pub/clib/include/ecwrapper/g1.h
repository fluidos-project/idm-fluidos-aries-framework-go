/**
 * @file g1.h
 * @author Jesús García Rodríguez
 * @brief File with method definitions for the (EC) group 1 on a cryptographic pairing
 *
 * This file defines the type and methods for manipulating elements form the first group (additive notation) of a
 * cryptographic pairing, including generators, serialization, addition, multiplication...
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef G1_H
#define G1_H

#include <Zp.h> 

/**
 * Encapsulated declaration of type G1, which represents elements from G1 of the pairing
 */
typedef struct G1Impl G1;


/**
 * @brief Generates a copy of the group generator. Can be safely modified. Has to be freed after usage
 * 
 * @return 
 */
G1* g1Generator();

/**
 * @brief Generates a copy of group identity (addition operation). Can be safely modified. Has to be freed after usage
 * 
 * @return 
 */
G1* g1Identity();

/**
 * @brief Size of byte representation of G1 element (when computing g1ToBytes)
 */
int g1ByteSize();

/**
 * @brief Represent G1 element as byte array. Array is assumed to be big enough for copy
 * 
 * @param res Byte array where it will be copied
 * @param a G1 element
 */
void g1ToBytes(char *res, const G1 *a);

/**
 * @brief Generates a hashed G1 point from bytes. Has to be freed after usage
 * 
 * @param bytes Byte array 
 * @param n Number of bytes
 * @return 
 */
G1* hashToG1(const char *bytes, int n);

/**
 * @brief Generates a G1 element from bytes (previously serialized). Has to be freed after usage
 * 
 * @param bytes Byte array of serialized element (assumed to be the correct length, i.e. g1ByteSyze())
 */
G1 * g1FromBytes(const char *bytes);

/**
 * @brief Addition within the curve a+b
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g1Add(G1* a, const G1* b);

/**
 * @brief Subtraction within the curve a-b
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g1Sub(G1* a, const G1* b);

/**
 * @brief Multiplication [b]a
 * 
 * @param a EC point, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g1Mul(G1* a, const Zp* b);

/**
 * @brief Multiplication [b] lt[0], with lt being a lookup table for G1 element g with lt[i]=g^(2^(i+1)). Result must be freed.
 * 
 * @param lt being a lookup table for G1 element g with lt[i]=g^(2^(i+1))
 * @param b Second operand
 * @return The resut [b]·lt[1]
 */
G1* g1MulLookup(const G1* lt[], const Zp* b);

/**
 * @brief Multiplication res= [b] lt[0], with lt being a lookup table for G1 element g with lt[i]=g^(2^(i+1)). 
 * 
 * @param res Will contain the result of the computation, should have been allocated before
 * @param lt A lookup table for G1 element g with lt[i]=g^(2^(i+1))
 * @param b Second operand
 */
void g1MulLookupWithoutAllocation(G1* res, const G1* lt[], const Zp* b);

/**
 * @brief Compute lookup table of size n for element g. Result must be freed.
 * @return The lookup table, res[i]=g^(2^(i+1))
 */
G1** g1CompLookupTable(const G1* g, int n);

/**
 * @brief Inverse multiplication [-b]a
 * 
 * @param a EC point, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void g1InvMul(G1* a, const Zp* b);

/**
 * @brief Multi-multiplication (n-multiplication), res=Sigma [b_i]a_i
 * 
 * @param a Array of bases (G1 elements)
 * @param b Array of multipliers (Zp elements)
 * @param n Length of arrays (assumed to be same and valid)
 */
G1* g1Muln(const G1* a[], const Zp* b[],int n);

/**
 * @brief Check if element is identity
 * 
 * @param a EC point
 * @return 1 if identity, 0 if not
 */
int g1IsIdentity(const G1* a);

/**
 * @brief Checks a=b in G1
 * 
 * @param a First element
 * @param b Second element
 * @returns 1 if equal 0 if not
 */
int g1Equals(const G1* a, const G1* b);

/**
 * @brief Copies G1 element. Result has to be freed after usage
 * 
 * @param a G1 element to be copied
 */
G1* g1Copy(const G1* a);

/**
 * @brief Overwrite a value to b (no new memory is reserved)
 * 
 * @param a G1 element, after method a=b
 * @param b G1 element, not modified
 */
void g1CopyValue(G1* a, const G1* b);

/**
 * @brief Output to stdout G1 (for testing)
 * 
 * @param e 
 */
void g1Print(const G1* e);


/**
 * @brief Destroy G1 element, freeing memory (from malloc etc.) as needed. Depending on internal structure might be equivalent to free(e)
 * 
 * @param e 
 */
void g1Free(G1* e);

#endif 