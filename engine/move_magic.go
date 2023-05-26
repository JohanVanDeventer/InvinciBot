package main

import "math/rand"

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Magic Bitboards: Goal ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Magic bitboards allow us to look up the rook and bishop moves from a table directly, given a set of blockers.

Code for generating the moves will ultimately be to make a magic struct that we can use based on the square we are on
to generate the pseudo-legal moves of bishops and rooks given the input blockers.

func getRookMovesPseudo(blockers Bitboard, sq int) Bitboard {
	blockers &= magicStructsRooks[sq].mask          // PART 1
	blockers *= magicStructsRooks[sq].magic         // PART 2
	blockers >>= (64 - magicStructsRooks[sq].shift) // PART 3
	return magicRookMovesTable[sq][blockers]        // PART 4
}

func getBishopMovesPseudo(blockers Bitboard, sq int) Bitboard {
	blockers &= magicStructsBishops[sq].mask          // PART 1
	blockers *= magicStructsBishops[sq].magic         // PART 2
	blockers >>= (64 - magicStructsBishops[sq].shift) // PART 3
	return magicBishopMovesTable[sq][blockers]        // PART 4
}

PART 1
------
We start with a mask that contains the "POSSIBLE BLOCKERS" bits for a piece on a square.
Only these squares can block the piece moves to LESS than it's full pseudo-legal moves.

8 | . . . . . . . .
7 | . . . . 1 . . .
6 | . . . . 1 . . .
5 | . 1 1 1 . 1 1 .
4 | . . . . 1 . . .
3 | . . . . 1 . . .
2 | . . . . 1 . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

We "&" this mask with the blockers to obtain only the blockers that influence the piece moves.


PART 2 and 3
------------
Now that we have the blockers that influence the piece moves, we need to transform it into a key.
We use the key we get to look up a precalculated move table that gives the moves based on this key.

To transform it into a key, we need to generate the following:
- Magic number (to multiply by)
- Shift number (to right shift by)

The result is a key used for the table lookup.

*/

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------- Magic Bitboards: Global Tables -----------------------------------------
// --------------------------------------------------------------------------------------------------------------------

type MagicForSq struct {
	mask  Bitboard // PART 1
	magic Bitboard // PART 2
	shift int      // PART 3
}

// store a magic struct for each square
var magicStructsRooks [64]MagicForSq
var magicStructsBishops [64]MagicForSq

// stores moves that are looked up later using a magic struct
// the table is indexed as: [sq][key]
var magicRookMovesTable [64][4096]Bitboard
var magicBishopMovesTable [64][512]Bitboard

/*
// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------- Magic Bitboards: PART 1: Masks -----------------------------------------
// --------------------------------------------------------------------------------------------------------------------
*/

// init the "masks" variable in the magic struct
func initMagicMasks() {

	// ----------------------------------- CLEARING MASKS ------------------------------------
	/*
		We set clearing masks to remove bits on the edges of the board.

		The useful coordinates to set clearing masks are:
		------- Coordinates ----------
		56, 57, 58, 59, 60, 61, 62, 63
		48, 49, 50, 51, 52, 53, 54, 55
		40, 41, 42, 43, 44, 45, 46, 47
		32, 33, 34, 35, 36, 37, 38, 39
		24, 25, 26, 27, 28, 29, 30, 31
		16, 17, 18, 19, 20, 21, 22, 23
		08, 09, 10, 11, 12, 13, 14, 15
		00, 01, 02, 03, 04, 05, 06, 07

	*/

	clearMaskTop := emptyBB
	for i := 56; i < 64; i++ {
		clearMaskTop.setBit(i)
	}
	clearMaskBottom := emptyBB
	for i := 0; i < 8; i++ {
		clearMaskBottom.setBit(i)
	}
	clearMaskLeft := emptyBB
	for i := 0; i < 64; i += 8 {
		clearMaskLeft.setBit(i)
	}
	clearMaskRight := emptyBB
	for i := 7; i < 64; i += 8 {
		clearMaskRight.setBit(i)
	}

	// ----------------------------------- ROOK MOVES ------------------------------------
	// we already have masks that include the outer edges
	// now we just need to remove those bits
	// we remove them if and only if the original square is not already on that mask line

	for sq := 0; sq < 64; sq++ {

		// get the pseudo-legal moves
		pseudoMask := moveRooksTable[sq]

		// mask out the edge bits
		if !clearMaskTop.isBitSet(sq) {
			pseudoMask &= ^clearMaskTop
		}
		if !clearMaskBottom.isBitSet(sq) {
			pseudoMask &= ^clearMaskBottom
		}
		if !clearMaskLeft.isBitSet(sq) {
			pseudoMask &= ^clearMaskLeft
		}
		if !clearMaskRight.isBitSet(sq) {
			pseudoMask &= ^clearMaskRight
		}

		// finally set the magic mask for this square
		magicStructsRooks[sq].mask = pseudoMask
	}

	// ----------------------------------- BISHOP MOVES ------------------------------------
	// we already have masks that include the outer edges
	// now we just need to remove those bits
	// we remove them if and only if the original square is not already on that mask line

	for sq := 0; sq < 64; sq++ {

		// get the pseudo-legal moves
		pseudoMask := moveBishopsTable[sq]

		// mask out the edge bits
		if !clearMaskTop.isBitSet(sq) {
			pseudoMask &= ^clearMaskTop
		}
		if !clearMaskBottom.isBitSet(sq) {
			pseudoMask &= ^clearMaskBottom
		}
		if !clearMaskLeft.isBitSet(sq) {
			pseudoMask &= ^clearMaskLeft
		}
		if !clearMaskRight.isBitSet(sq) {
			pseudoMask &= ^clearMaskRight
		}

		// finally set the magic mask for this square
		magicStructsBishops[sq].mask = pseudoMask
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------- Magic Bitboards: PART 2: Magic Numbers ------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Obtained the magic numbers from: https://github.com/GunshipPenguin/shallow-blue/blob/c6d7e9615514a86533a9e0ffddfc96e058fc9cfd/src/attacks.h#L120
*/
var rookMagics = [64]uint64{
	0xa8002c000108020, 0x6c00049b0002001, 0x100200010090040, 0x2480041000800801, 0x280028004000800,
	0x900410008040022, 0x280020001001080, 0x2880002041000080, 0xa000800080400034, 0x4808020004000,
	0x2290802004801000, 0x411000d00100020, 0x402800800040080, 0xb000401004208, 0x2409000100040200,
	0x1002100004082, 0x22878001e24000, 0x1090810021004010, 0x801030040200012, 0x500808008001000,
	0xa08018014000880, 0x8000808004000200, 0x201008080010200, 0x801020000441091, 0x800080204005,
	0x1040200040100048, 0x120200402082, 0xd14880480100080, 0x12040280080080, 0x100040080020080,
	0x9020010080800200, 0x813241200148449, 0x491604001800080, 0x100401000402001, 0x4820010021001040,
	0x400402202000812, 0x209009005000802, 0x810800601800400, 0x4301083214000150, 0x204026458e001401,
	0x40204000808000, 0x8001008040010020, 0x8410820820420010, 0x1003001000090020, 0x804040008008080,
	0x12000810020004, 0x1000100200040208, 0x430000a044020001, 0x280009023410300, 0xe0100040002240,
	0x200100401700, 0x2244100408008080, 0x8000400801980, 0x2000810040200, 0x8010100228810400,
	0x2000009044210200, 0x4080008040102101, 0x40002080411d01, 0x2005524060000901, 0x502001008400422,
	0x489a000810200402, 0x1004400080a13, 0x4000011008020084, 0x26002114058042,
}

var bishopMagics = [64]uint64{
	0x89a1121896040240, 0x2004844802002010, 0x2068080051921000, 0x62880a0220200808, 0x4042004000000,
	0x100822020200011, 0xc00444222012000a, 0x28808801216001, 0x400492088408100, 0x201c401040c0084,
	0x840800910a0010, 0x82080240060, 0x2000840504006000, 0x30010c4108405004, 0x1008005410080802,
	0x8144042209100900, 0x208081020014400, 0x4800201208ca00, 0xf18140408012008, 0x1004002802102001,
	0x841000820080811, 0x40200200a42008, 0x800054042000, 0x88010400410c9000, 0x520040470104290,
	0x1004040051500081, 0x2002081833080021, 0x400c00c010142, 0x941408200c002000, 0x658810000806011,
	0x188071040440a00, 0x4800404002011c00, 0x104442040404200, 0x511080202091021, 0x4022401120400,
	0x80c0040400080120, 0x8040010040820802, 0x480810700020090, 0x102008e00040242, 0x809005202050100,
	0x8002024220104080, 0x431008804142000, 0x19001802081400, 0x200014208040080, 0x3308082008200100,
	0x41010500040c020, 0x4012020c04210308, 0x208220a202004080, 0x111040120082000, 0x6803040141280a00,
	0x2101004202410000, 0x8200000041108022, 0x21082088000, 0x2410204010040, 0x40100400809000,
	0x822088220820214, 0x40808090012004, 0x910224040218c9, 0x402814422015008, 0x90014004842410,
	0x1000042304105, 0x10008830412a00, 0x2520081090008908, 0x40102000a0a60140,
}

// helper function to get a Bitboard/uint64 with only a few bits set
// we "&"" various random Bitboards to leave only overlapping bits
func getRandomSparseBitboard() Bitboard {
	filledBitboard := fullBB
	return filledBitboard & Bitboard(rand.Uint64()) & Bitboard(rand.Uint64()) & Bitboard(rand.Uint64())
}

func initMagicNumbers() {

	// rook magic numbers
	for sq := 0; sq < 64; sq++ {
		magicStructsRooks[sq].magic = Bitboard(rookMagics[63-sq])
	}

	// bishop magic numbers
	for sq := 0; sq < 64; sq++ {
		magicStructsBishops[sq].magic = Bitboard(bishopMagics[63-sq])
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------- Magic Bitboards: PART 3: Magic Shifts -------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Generate the magic shifts.

Note: the numbers below represent the maximum number of blockers from each square on the board.
For example, a rook on d4 has 10 blockers (3 up, 2 down, 3 right, 2 left) excluding edge bits.
For example, a rook on a1 has 12 blockers (6 up and 6 right).
For example, a bishop on d4 has 9 blockers (2 UL, 3 DR, 2 DL, 2 UR).const

This max amount of blockers represent the shift value.

This is also the reason for the 4096 and 512 in the tables above.
It is the max amount of permutations of the blockers.
2 ^ 12 = 4096.
2 ^ 9 = 512.
*/

// constant shift values for each square
var rookShifts = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var bishopShifts = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

func initMagicShifts() {

	// rook shifts
	for sq := 0; sq < 64; sq++ {
		magicStructsRooks[sq].shift = rookShifts[sq]
	}

	// bishop shifts
	for sq := 0; sq < 64; sq++ {
		magicStructsBishops[sq].shift = bishopShifts[sq]
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------- Magic Bitboards: PART 4: Global Move Table -----------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Generate the global move lookup tables.
*/

// original function used to get rook moves
// now used to fill the magic move tables
func getRookMovesPseudoOriginal(sq int, blockers Bitboard) Bitboard {

	var newBitboard = emptyBB

	// ------------ UP ----------------
	newBitboard |= moveRaysTable[sq][RAY_UP]
	if moveRaysTable[sq][RAY_UP]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_UP] & blockers).getMSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_UP]
	}

	// ------------ RIGHT ----------------
	newBitboard |= moveRaysTable[sq][RAY_RIGHT]
	if moveRaysTable[sq][RAY_RIGHT]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_RIGHT] & blockers).getMSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_RIGHT]
	}

	// ------------ DOWN ----------------
	newBitboard |= moveRaysTable[sq][RAY_DOWN]
	if moveRaysTable[sq][RAY_DOWN]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_DOWN] & blockers).getLSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_DOWN]
	}

	// ------------ LEFT ----------------
	newBitboard |= moveRaysTable[sq][RAY_LEFT]
	if moveRaysTable[sq][RAY_LEFT]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_LEFT] & blockers).getLSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_LEFT]
	}

	return newBitboard

}

// original function used to get bishop moves
// now used to fill the magic move tables
func getBishopMovesPseudoOriginal(sq int, blockers Bitboard) Bitboard {
	var newBitboard = emptyBB

	// ------------ UP LEFT ----------------
	newBitboard |= moveRaysTable[sq][RAY_UL]
	if moveRaysTable[sq][RAY_UL]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_UL] & blockers).getMSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_UL]
	}

	// ------------ UP RIGHT ----------------
	newBitboard |= moveRaysTable[sq][RAY_UR]
	if moveRaysTable[sq][RAY_UR]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_UR] & blockers).getMSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_UR]
	}

	// ------------ DOWN RIGHT ----------------
	newBitboard |= moveRaysTable[sq][RAY_DR]
	if moveRaysTable[sq][RAY_DR]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_DR] & blockers).getLSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_DR]
	}

	// ------------ DOWN LEFT ----------------
	newBitboard |= moveRaysTable[sq][RAY_DL]
	if moveRaysTable[sq][RAY_DL]&blockers != 0 {
		blockerSq := (moveRaysTable[sq][RAY_DL] & blockers).getLSBSq()
		newBitboard &= ^moveRaysTable[blockerSq][RAY_DL]
	}

	return newBitboard
}

// function to get all the permutations of a blocker bitboard (max 12 bits set for rooks and 9 bits set for bishops, but can be less)
func generateBlockerPermutations(num uint64) []uint64 {

	// count the number of bits set to 1
	var count uint64 = 0
	tempNum := num
	for tempNum != 0 {
		count++
		tempNum &= tempNum - 1
	}

	// generate permutations recursively
	result := make([]uint64, 0)
	generatePermutationsRecursively(num, count, 0, 0, &result)
	return result
}

func generatePermutationsRecursively(num uint64, count, index, current uint64, result *[]uint64) {
	if count == 0 {
		*result = append(*result, current)
		return
	}

	// iterate through each bit position
	for i := index; i < 64; i++ {

		// check if the bit is set to 1
		if num&(1<<i) != 0 {

			// generate permutation with this bit set to 0
			generatePermutationsRecursively(num, count-1, i+1, current, result)

			// generate permutation with this bit set to 1
			generatePermutationsRecursively(num, count-1, i+1, current|(1<<i), result)
		}
	}
}

// we now init the magic tables
// by filling in each possible square and key we can get
// with the old move generation bitboard
func initMagicMoveTables() {

	// ------------ rooks -------------

	// for each square on the board
	for sq := 0; sq < 64; sq++ {

		// get the mask with all possible blockers
		blockers := magicStructsRooks[sq].mask

		// get all the permutations from those blockers
		blockerPermutations := generateBlockerPermutations(uint64(blockers))

		// for each permutation, generate the key for the square, and set the moves of that key to the actual moves generated using the old way
		for _, blockerPermutation := range blockerPermutations {
			key := Bitboard(blockerPermutation) // this already is after applying the mask above
			key *= magicStructsRooks[sq].magic
			key >>= (64 - magicStructsRooks[sq].shift)
			magicRookMovesTable[sq][key] = getRookMovesPseudoOriginal(sq, Bitboard(blockerPermutation))
		}

	}

	// ------------ bishops -------------

	// for each square on the board
	for sq := 0; sq < 64; sq++ {

		// get the mask with all possible blockers
		blockers := magicStructsBishops[sq].mask

		// get all the permutations from those blockers
		blockerPermutations := generateBlockerPermutations(uint64(blockers))

		// for each permutation, generate the key for the square, and set the moves of that key to the actual moves generated using the old way
		for _, blockerPermutation := range blockerPermutations {
			key := Bitboard(blockerPermutation) // this already is after applying the mask above
			key *= magicStructsBishops[sq].magic
			key >>= (64 - magicStructsBishops[sq].shift)
			magicBishopMovesTable[sq][key] = getBishopMovesPseudoOriginal(sq, Bitboard(blockerPermutation))
		}
	}
}
