/**
 * @file utils.h
 * @author Jesús García Rodríguez
 * @brief File with utilities
 *
 * This file defines utility methods, like secure random number generator
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef UTILS_H
#define UTILS_H

/**
 * Encapsulated declaration of type ranGen, which represents random generator structure (can be initialized with a seed, used to generate random bytes...)
 */
typedef struct ranGen ranGen;

//TODO One or multiple functions to recover information (depending on how we model it) about the concrete implementation

/**
 * @brief Initialize a ranGen (random generator) with seed. Must be freed after usage.
 * 
 * @param seed Bytes for seed
 * @param n Number of bytes provided within the seed
 * @return ranGen* 
 */
ranGen * rgInit(const char *seed,int n);

/**
 * @brief Generate n random bytes from random generator
 * 
 * @param rg Random generator (previously initialized)
 * @param n Number of bytes to be generated
 * @return ranGen* 
 */
char * rgGenBytes(ranGen * rg,int n);

/**
 * @brief Destroy Random Generator, freeing memory (from malloc etc.) as needed. Depending on internal structure might be equivalent to free(rg)
 * 
 * @param rg Random generator to be freed
 */
void rgFree(ranGen* rg);

#endif 