package main

import (
	"fmt"
	"math/bits"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Bitwise Operations -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

Explanation of the common types of bit operations:
&	(AND) Sets 1 in the positions where BOTH bitboards have 1 in those positions
|	(OR) Sets 1 in the positions where AT LEAST 1 bitboard have 1 in those positions
^	(XOR) Sets 1 in the positions where EXACTLY 1 bitboard has 1 in those positions and the other 0

*/

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------- Setup ------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

type Bitboard uint64

// a full and an empty bitboard
const fullBB Bitboard = 0xffffffffffffffff
const emptyBB Bitboard = 0x0

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Bit Reference Table ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// a reference table where each square on the board has that specific bit set, and only that bit
// therefore lookup can be done for a specific square from 0 - 63, and used in relevant functions
var bbReferenceArray [64]Bitboard

func initBBReferenceArray() {
	// sets the right most bit to 1
	var startingBit = emptyBB + 1

	// set each square's bit in the table
	for sq := 0; sq < 64; sq++ {
		bbReferenceArray[63-sq] = startingBit
		startingBit = startingBit << 1
	}
}

// just a test function to see if the initialized table is correct
func printBBReferenceArray() {
	for _, num := range bbReferenceArray {
		fmt.Printf("%064b\n", num)
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Bitboard Functions -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// functions to manipulate specific bitboard bits, or get information about them

func (bb *Bitboard) setBit(sq int) {
	*bb |= bbReferenceArray[sq]
}

func (bb *Bitboard) clearBit(sq int) {
	*bb &= ^bbReferenceArray[sq]
}

func (bb Bitboard) isBitSet(sq int) bool {
	return (bb & bbReferenceArray[sq]) != 0
}

func (bb Bitboard) countBits() int {
	return bits.OnesCount64(uint64(bb))
}

// returns the sq index of the most significant bit
// for example: 0000000000000000000000000000000000000000000000000000000000000001 is 63
// will return 64 if all zeros
func (bb Bitboard) getMSBSq() int {
	return bits.LeadingZeros64(uint64(bb))
}

func (bb Bitboard) getLSBSq() int {
	return 63 - bits.TrailingZeros64(uint64(bb))
}

// returns the sq index of the next bit, and clears that bit to zero
// will return 64 if there are no more bits (i.e. 64 zeros)
func (bb *Bitboard) popBitGetSq() int {
	sq := bb.getMSBSq()
	bb.clearBit(sq)
	return sq
}

func (bb Bitboard) printBitboardSingleLine() {
	fmt.Printf("%064b\n", bb)
}

func (bb Bitboard) printBitboardFancy8x8() {
	for i := uint(0); i < 8; i++ {
		row := uint64(bb >> (i * 8) & 0xff)
		fmt.Printf("%08b\n", row)
	}
}
