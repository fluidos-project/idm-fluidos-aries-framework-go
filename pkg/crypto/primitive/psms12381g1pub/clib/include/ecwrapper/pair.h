/**
 * @file pair.h
 * @author Jesús García Rodríguez
 * @brief File with method definitions for the pairing operation
 *
 * This file defines the methods for pairing, including a simple pair and optimized multi-pairings
 * 
 * @see https://github.com/JesusGarciaRodriguez/dpabcCimplementation
 */
#ifndef PAIR_H
#define PAIR_H

#include <g1.h>
#include <g2.h> 
#include <g3.h> 


/**
 * @brief Pairing operation result=e(a,b). Result is a generated G3 element that must be freed after use
 * 
 * @param a First operand (G1)
 * @param b Second operand (G2)
 */
G3* pair(const G1 *a, const G2 *b);

/**
 * @brief Multi-pairing operation result=PI(e(a_i,b_i)). Result is a generated G3 element that must be freed after use
 *  
 * 
 * @param a Array of operands from G1
 * @param b Array of operands from G2
 * @param n Number of elements
 */
G3* multipair(const G1 *a[], const G2 *b[], int n); 

/**
 * @brief Double pairing operation result=e(a1,b1)·e(a2,b2). Result is a generated G3 element that must be freed after use
 *  
 * 
 * @param a1 First G1 operand
 * @param a2 Second G1 operand
 * @param b1 First G2 operand
 * @param b2 Second G2 operand
 * @param n Number of elements
 */
G3* doublepair(const G1 *a1,const G1* a2, const G2 *b1,const G2 *b2);


#endif 