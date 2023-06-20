package main

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Pawn Evaluation Setup ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// functions and structures to evaluate the pawn structure of a given position

// stores the eval results of the last pawn structure evaluation
// and overwrites the last entry once the pawn structure changes
const (
	PAWN_HASH_TABLE_SIZE = 7777 // size of the hash table for pawn evaluations
	// note: we keep the size an uneven number to increase the spread of hashes over the size
	PAWN_HASH_TABLE_SIZE_BB = Bitboard(PAWN_HASH_TABLE_SIZE) // size of the hash table for pawn evaluations
)

type PawnStructureTable struct {
	whitePawns Bitboard
	blackPawns Bitboard
	value      int
}

var pawnColumnMasks [8]Bitboard // masks where all the bits for that column only is set
var pawnRowMasks [8]Bitboard    // masks where all the bits for that row only is set

var pawnIsolatedMasks [8]Bitboard   // masks that has the column and both side columns set
var pawnPassedMasks [2][64]Bitboard // masks that has 3 columns in front (not side) of the square set

// create masks for each column on the board
func initEvalPawnMasks() {

	// masks for each column only (for doubled pawns checking)
	for col := 0; col < 8; col++ {
		newBitboard := emptyBB
		for row := 0; row < 8; row++ {
			newBitboard.setBit(sqFromRowAndCol(row, col))
		}
		pawnColumnMasks[col] = newBitboard
	}

	// masks for each row only
	for row := 0; row < 8; row++ {
		newBitboard := emptyBB
		for col := 0; col < 8; col++ {
			newBitboard.setBit(sqFromRowAndCol(row, col))
		}
		pawnRowMasks[row] = newBitboard
	}

	// masks that has the column and both side columns set (for isolated pawn checking)
	for col := 0; col < 8; col++ {
		newBitboard := pawnColumnMasks[col]

		// add left col
		colLeft := col - 1
		if colLeft >= 0 {
			newBitboard |= pawnColumnMasks[colLeft]
		}

		// add right col
		colRight := col + 1
		if colRight <= 7 {
			newBitboard |= pawnColumnMasks[colRight]
		}

		// save the mask
		pawnIsolatedMasks[col] = newBitboard
	}

	// masks that has the 3 columns and all rows in front of the square set for the specific side
	for side := 0; side < 2; side++ {
		for sq := 0; sq < 64; sq++ {
			row, col := rowAndColFromSq(sq)

			// set the used columns
			usedCols := pawnIsolatedMasks[col]

			// set the used rows
			usedRows := emptyBB
			for rowUsed := 0; rowUsed < 8; rowUsed++ {
				if side == SIDE_WHITE {
					if rowUsed > row {
						usedRows |= pawnRowMasks[rowUsed]
					}
				}
				if side == SIDE_BLACK {
					if rowUsed < row {
						usedRows |= pawnRowMasks[rowUsed]
					}
				}
			}

			// finally set the combined mask
			finalMask := usedCols & usedRows
			pawnPassedMasks[side][sq] = finalMask
		}
	}
}
