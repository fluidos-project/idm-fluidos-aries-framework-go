/**
 * @file Zp.h
 * @brief File with method definitions for integers modulo p, with p prime (used in schemes that take advantage of EC crypto)
 *
 * @details This file defines the type and methods for manipulating integers modulo p, with p prime which will be decided at compile time
 * depending on the instantiation of this wrapper library (depends on the needs of the pariring-friendly-EC based scheme)
 * 
 * @author Jesús García Rodríguez
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef ZP_H
#define ZP_H
#include <utils.h>

/**
 * Encapsulated declaration of type Zp, which represents integers modulo p, p prime 
 * The internal management of modular expressions is up to implementation, but it must conceptually work as if everything is always being computed within the field Zp
 */
typedef struct ZpImpl Zp;


/**
 * @brief Size of byte representation of Zp element (when computing zpToBytes)
 */
int zpByteSize();

/**
 * @brief Represent Zp as byte array. Array is assumed to be big enough for copy
 * 
 * @param res Byte array where it will be copied
 * @param a Zp element
 */
void zpToBytes(char *res, const Zp *a);

/**
 * @brief Generates a Zp element from bytes (previously serialized). Has to be freed after usage
 * 
 * @param bytes Byte array of serialized element (assumed to be the correct length, i.e. zpByteSyze())
 */
Zp * zpFromBytes(const char *bytes);

//TODO Add if possible generic ZpFromBigIntegerRepresentation or establish a method for doing this in documentation

/**
 * @brief Construct a new Zp element, hashing from bytes. Has to be freed after usage
 * 
 * @param bytes Bytes for hash
 * @param nBytes Number of bytes
 */
Zp *hashToZp(const char * bytes,int nBytes);

/**
 * @brief Zp from integer. Has to be freed after usage
 * 
 * @param a Integer
 * @return Zp Resulting Zp
 */
Zp* zpFromInt(int a);

/**
 * @brief Generate random Zp. Has to be freed after usage
 * 
 * @param rg Random Generator
 * @return 
 */
Zp* zpRandom(ranGen *rg);

/**
 * @brief Change value of Zp to random value
 * 
 * @param rg Random Generator
 * @param r Zp element, will be set to random value
 */
void zpRandomValue(ranGen *rg, Zp *r);

/**
 * @brief Addition a+=b (mod p)
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void zpAdd(Zp* a, const Zp* b);

/**
 * @brief Addition a-=b (mod p)
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void zpSub(Zp* a, const Zp* b);

/**
 * @brief Multiplication a*=b (mod p)
 * 
 * @param a First operand, will be modified to the resulting value
 * @param b Second operand, not modified
 */
void zpMul(Zp* a, const Zp* b);


/**
 * @brief Negate modulo p (-a) (mod p)
 * 
 * @param a First operand, will be modified to the resulting value
 */
void zpNeg(Zp* a);

/**
 * @brief Return the number of bits (ignoring left-zeros, i.e. position of most significant bit counting from the least significant bit as 1) of a
 * 
 * For instance: For a=7 return 3, for a=8 return 4.
 */
int zpNbits(Zp* a);

/**
 * @brief Return the parity (0 even, 1 odd) of a
 */
int zpParity(Zp* a);

/**
 * @brief Double a (usually through bit-shift)
 */
void zpDouble(Zp* a);

/**
 * @brief Halve a, return reminder (usually through bit-shift)
 */
int zpHalf(Zp* a);

/**
 * @brief Checks a=b in Zp
 * 
 * @param a First element
 * @param b Second element
 * @returns 1 if equal 0 if not
 */
int zpEquals(const Zp* a, const Zp* b);

/**
 * @brief Checks a=0 in Zp
 * 
 * @returns 1 if zero, 0 if not
 */
int zpIsZero(const Zp* a);

/**
 * @brief Copies Zp element. Result has to be freed after usage
 * 
 * @param a Zp element to be copied
 */
Zp* zpCopy(const Zp* a);


/**
 * @brief Overwrite a value to b (no new memory is reserved)
 * 
 * @param a Zp element, after method a=b
 * @param b Zp element, not modified
 */
void zpCopyValue(Zp* a, const Zp* b);

/**
 * @brief Output to stdout Zp (for testing)
 * 
 * @param e 
 */
void zpPrint(const Zp* e);


/**
 * @brief Destroy Zp element, freeing memory (from malloc etc.) as needed. Depending on internal structure might be equivalent to free(e)
 * 
 * @param e 
 */
void zpFree(Zp* e); 

#endif 