package main

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------- Bitboard Move Lookup Tables -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

/*
The purpose is to initialize tables that give the possible move squares of a piece on a particular square,
given the remaining pieces on the board.

------------------------------------------------ Single Movers (Normal Bitboards) -------------------------------------
Kings, knights and pawns do not have rays where pieces can potentially block these rays.
They just have normal lookup tables that list possible moves from a given square.

Pawns have separate tables for moving and attacking, whereas kings and knights have the same table for moves and attacks.

---------------------------------------------- Ray Movers (Classic Bitboard Approach) ---------------------------------
Ray moves are complex, because of blocking pieces.
The classic approach is followed (not magic bitboards) as explained below.

This is a blank board:

8 | . . . . . . . .
7 | . . . . . . . .
6 | . . . . . . . .
5 | . . . . . . . .
4 | . . . . . . . .
3 | . . . . . . . .
2 | . . . . . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

Say there is a bishop on C3. The normal moves for the bishop to the top right is:

8 | . . . . . . . 1
7 | . . . . . . 1 .
6 | . . . . . 1 . .
5 | . . . . 1 . . .
4 | . . . 1 . . . .
3 | . . 1 . . . . .
2 | . . . . . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

Then say we have the following blockers:

8 | . 1 . 1 . . . 1
7 | . . . . . . . .
6 | 1 . . . 1 1 1 .
5 | . . . . . . . .
4 | . . . . . . . .
3 | . . B . . . . .
2 | 1 . . 1 . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

For this ray of the bishop, we "&" these bitboards to get the blockers intersecting the ray:

8 | . . . . . . . 1
7 | . . . . . . . .
6 | . . . . . 1 . .
5 | . . . . . . . .
4 | . . . . . . . .
3 | . . . . . . . .
2 | . . . . . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

We call this the "masked blockers".

We now need to identify the "first" blocker in the ray, because that is the blocker that will actually block the ray.
We can use the leading zeroes function to determine that bit (scanning from 0 to 63).
This will return the square index of that blocker (in this case 45, which is F6).

We then use that index to cast a further ray out from the point on F6 in the original direction.
We can use our initial lookup table for that. We then get:

8 | . . . . . . . 1
7 | . . . . . . 1 .
6 | . . . . . . . .
5 | . . . . . . . .
4 | . . . . . . . .
3 | . . . . . . . .
2 | . . . . . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

These are therefore the bits that need to be masked out from the initial attacking squares.
The mask to apply is therefore the inverse of this bitboard:

8 | 1 1 1 1 1 1 1 .
7 | 1 1 1 1 1 1 . 1
6 | 1 1 1 1 1 1 1 1
5 | 1 1 1 1 1 1 1 1
4 | 1 1 1 1 1 1 1 1
3 | 1 1 1 1 1 1 1 1
2 | 1 1 1 1 1 1 1 1
1 | 1 1 1 1 1 1 1 1
   ----------------
    a b c d e f g h

We do this for each of the 4 directions, to at the end get an attack map of:

8 | . . . . . . . .
7 | . . . . . . . .
6 | . . . . . 1 . .
5 | 1 . . . 1 . . .
4 | . 1 . 1 . . . .
3 | . . B . . . . .
2 | . 1 . 1 . . . .
1 | 1 . . . . . . .
   ----------------
    a b c d e f g h

NOTE!
The above move/attack map includes the blockers themselves.
This is now the same as the maps for the kings, knights and pawns.
Both these sets will include friendly pieces as blockers.
Therefore it is important that the blockers needs to be filtered for friendly pieces later.

*/

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------- Lookup Tables --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// single movement lookup tables - kings, knights, pawns
var moveKingsTable [64]Bitboard          // kings (colour does not matter)
var moveKnightsTable [64]Bitboard        // knights (colour does not matter)
var moveOnlyPawnsTable [64][2]Bitboard   // pawn moves only, not attacks (colour matters)
var moveAttackPawnsTable [64][2]Bitboard // pawn attacks only, not moves (colour matters)

// ray movement tables - queens, rooks, bishops
var moveRaysTable [64][8]Bitboard // moves up to and including the edge squares, for each of 8 individual directions
var moveRooksTable [64]Bitboard   // combines the up, right, down and left directions
var moveBishopsTable [64]Bitboard // combines the UL, UR, DR, DL directions

// pawn double move masks
var movePawnDoubleMasks [64][2]Bitboard // masks that have the 2 bits in front of pawns set to check for blockers

// castling is clear masks
var moveCastlingIsClearMasks [4]Bitboard // correspond to the 4 sides of castling, to check for blockers in those squares

// pawns attacking king masks
var movePawnsAttackingKingMasks [64][2]Bitboard // from a given king position, which enemy pawns can attack the king

// pinned piece masks
var movePinnedMasksTable [64][4]Bitboard // masks for each of the 4 types of pins

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------- Directions ---------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// direction struct make it easy to work with rows and columns
type Direction struct {
	rowChange int
	colChange int
}

// set up the directions for various directions and pieces
var (
	DIR_UP    = Direction{1, 0}
	DIR_DOWN  = Direction{-1, 0}
	DIR_LEFT  = Direction{0, -1}
	DIR_RIGHT = Direction{0, 1}

	DIR_UL = Direction{1, -1}
	DIR_UR = Direction{1, 1}
	DIR_DL = Direction{-1, -1}
	DIR_DR = Direction{-1, 1}

	DIR_UUL = Direction{2, -1}
	DIR_UUR = Direction{2, 1}
	DIR_RRU = Direction{1, 2}
	DIR_RRD = Direction{-1, 2}
	DIR_DDR = Direction{-2, 1}
	DIR_DDL = Direction{-2, -1}
	DIR_LLD = Direction{-1, -2}
	DIR_LLU = Direction{1, -2}

	DIR_UU = Direction{2, 0}
	DIR_DD = Direction{-2, 0}
)

const (
	RAY_UP    = 0
	RAY_RIGHT = 1
	RAY_DOWN  = 2
	RAY_LEFT  = 3
	RAY_UL    = 4
	RAY_UR    = 5
	RAY_DR    = 6
	RAY_DL    = 7

	PIN_UD    = 0
	PIN_LR    = 1
	PIN_ULtDR = 2
	PIN_DLtUR = 3
)

// set up an array of directions to loop over for each piece type
var (
	directionsQueen = [8]Direction{DIR_UP, DIR_RIGHT, DIR_DOWN, DIR_LEFT, DIR_UL, DIR_UR, DIR_DR, DIR_DL}

	directionsKnight = [8]Direction{DIR_UUL, DIR_UUR, DIR_RRU, DIR_RRD, DIR_DDR, DIR_DDL, DIR_LLD, DIR_LLU}
	directionsKing   = [8]Direction{DIR_UP, DIR_UR, DIR_RIGHT, DIR_DR, DIR_DOWN, DIR_DL, DIR_LEFT, DIR_UL}

	directionsWhitePawnCaptures = [2]Direction{DIR_UL, DIR_UR}
	directionsBlackPawnCaptures = [2]Direction{DIR_DL, DIR_DR}

	//directionsWhitePawnMovesStart = [2]Direction{DIR_UP, DIR_UU}
	directionsWhitePawnMovesAfter = [1]Direction{DIR_UP}
	//directionsBlackPawnMovesStart = [2]Direction{DIR_DOWN, DIR_DD}
	directionsBlackPawnMovesAfter = [1]Direction{DIR_DOWN}
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------- Move Lookup Table (Pawns) ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

func initMoveLookupTablePawns() {
	for colour := 0; colour < 2; colour++ {
		for sq := 0; sq < 64; sq++ {
			row, col := rowAndColFromSq(sq)
			moveOnlyPawnsTable[sq][colour] = moveLookupForPawn(row, col, colour)
			moveAttackPawnsTable[sq][colour] = attackLookupForPawn(row, col, colour)
			movePawnDoubleMasks[sq][colour] = doubleMoveMaskLookupForPawn(row, col, colour)
		}
	}

}

func moveLookupForPawn(row int, col int, colour int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	if colour == SIDE_WHITE {

		if row >= 1 && row <= 6 {
			for _, value := range directionsWhitePawnMovesAfter {
				adjustedRow := row + value.rowChange
				adjustedCol := col + value.colChange

				if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
					newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
				}
			}
		}
	}

	if colour == SIDE_BLACK {

		if row >= 1 && row <= 6 {
			for _, value := range directionsBlackPawnMovesAfter {
				adjustedRow := row + value.rowChange
				adjustedCol := col + value.colChange

				if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
					newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
				}
			}
		}
	}

	return newEmptyBB
}

func attackLookupForPawn(row int, col int, colour int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	if colour == SIDE_WHITE {

		for _, value := range directionsWhitePawnCaptures {
			adjustedRow := row + value.rowChange
			adjustedCol := col + value.colChange

			if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
				newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
			}
		}
	}

	if colour == SIDE_BLACK {

		for _, value := range directionsBlackPawnCaptures {
			adjustedRow := row + value.rowChange
			adjustedCol := col + value.colChange

			if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
				newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
			}
		}
	}

	return newEmptyBB
}

func doubleMoveMaskLookupForPawn(row int, col int, colour int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	if colour == SIDE_WHITE {
		if row == 1 {
			newEmptyBB.setBit(sqFromRowAndCol(row+1, col))
			newEmptyBB.setBit(sqFromRowAndCol(row+2, col))
		}
	}

	if colour == SIDE_BLACK {
		if row == 6 {
			newEmptyBB.setBit(sqFromRowAndCol(row-1, col))
			newEmptyBB.setBit(sqFromRowAndCol(row-2, col))
		}
	}

	return newEmptyBB
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Move Lookup Table (Kings) ------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

func initMoveLookupTableKings() {
	for sq := 0; sq < 64; sq++ {
		moveKingsTable[sq] = moveLookupForKing(sq)
	}
}

func moveLookupForKing(sq int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	for _, value := range directionsKing {
		rowCh := value.rowChange
		colCh := value.colChange

		adjustedRow, adjustedCol := rowAndColFromSq(sq)
		adjustedRow += rowCh
		adjustedCol += colCh

		if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
			newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
		}
	}

	return newEmptyBB
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Move Lookup Table (Knights) ----------------------------------------
// --------------------------------------------------------------------------------------------------------------------

func initMoveLookupTableKnights() {
	for sq := 0; sq < 64; sq++ {
		moveKnightsTable[sq] = moveLookupForKnight(sq)
	}
}

func moveLookupForKnight(sq int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	for _, value := range directionsKnight {
		rowCh := value.rowChange
		colCh := value.colChange

		adjustedRow, adjustedCol := rowAndColFromSq(sq)
		adjustedRow += rowCh
		adjustedCol += colCh

		if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
			newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
		}
	}

	return newEmptyBB
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------- Ray Lookup Table (Queens, Rooks, Bishops) ---------------------------------
// --------------------------------------------------------------------------------------------------------------------

func initMoveLookupTableRays() {

	// get moves for each of the 8 directions
	for sq := 0; sq < 64; sq++ {
		for ray_dir := 0; ray_dir < 8; ray_dir++ {
			moveRaysTable[sq][ray_dir] = moveLookupForRays(sq, ray_dir)
		}
	}

	// combine these into one table for rooks and bishops
	for sq := 0; sq < 64; sq++ {

		// rooks
		var newEmptyBBRooks Bitboard = emptyBB
		for ray_dir := 0; ray_dir < 4; ray_dir++ {
			newEmptyBBRooks |= moveRaysTable[sq][ray_dir]
		}
		moveRooksTable[sq] = newEmptyBBRooks

		// bishops
		var newEmptyBBBishops Bitboard = emptyBB
		for ray_dir := 4; ray_dir < 8; ray_dir++ {
			newEmptyBBBishops |= moveRaysTable[sq][ray_dir]
		}
		moveBishopsTable[sq] = newEmptyBBBishops
	}
}

func moveLookupForRays(sq int, ray_dir int) Bitboard {
	var newEmptyBB Bitboard = emptyBB

	current_direction := directionsQueen[ray_dir]

	rowCh := current_direction.rowChange
	colCh := current_direction.colChange

	adjustedRow, adjustedCol := rowAndColFromSq(sq)

	inBounds := true
	for inBounds {
		adjustedRow += rowCh
		adjustedCol += colCh
		if 0 <= adjustedRow && adjustedRow <= 7 && 0 <= adjustedCol && adjustedCol <= 7 {
			newEmptyBB.setBit(sqFromRowAndCol(adjustedRow, adjustedCol))
		} else {
			inBounds = false
		}
	}

	return newEmptyBB
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Castling Masks ------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// set the bits inbetween the king and the rook to later determine if that path is clear
func initMoveCastlingMasks() {

	// castle K
	moveCastlingIsClearMasks[CASTLE_WHITE_KINGSIDE].setBit(5)
	moveCastlingIsClearMasks[CASTLE_WHITE_KINGSIDE].setBit(6)

	// castle Q
	moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE].setBit(1)
	moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE].setBit(2)
	moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE].setBit(3)

	// castle k
	moveCastlingIsClearMasks[CASTLE_BLACK_KINGSIDE].setBit(61)
	moveCastlingIsClearMasks[CASTLE_BLACK_KINGSIDE].setBit(62)

	// castle q
	moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE].setBit(57)
	moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE].setBit(58)
	moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE].setBit(59)
}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------- Pawns Attacking King Masks -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// given a certain king position, which pawns of the given colour can attack that king

func initMovePawnAttackingKingMasks() {
	for colour := 0; colour < 2; colour++ {
		for sq := 0; sq < 64; sq++ {
			row, col := rowAndColFromSq(sq)

			// the king is attacked from down right and down left for white pawns
			if colour == SIDE_WHITE {

				// attacks from down right
				if col <= 6 && row >= 1 {
					movePawnsAttackingKingMasks[sq][SIDE_WHITE].setBit(sqFromRowAndCol(row-1, col+1))
				}

				// attacks from down left
				if col >= 1 && row >= 1 {
					movePawnsAttackingKingMasks[sq][SIDE_WHITE].setBit(sqFromRowAndCol(row-1, col-1))
				}

				// the king is attacked from up right and up left for black pawns
			} else {

				// attacks from up right
				if col <= 6 && row <= 6 {
					movePawnsAttackingKingMasks[sq][SIDE_BLACK].setBit(sqFromRowAndCol(row+1, col+1))
				}

				// attacks from up left
				if col >= 1 && row <= 6 {
					movePawnsAttackingKingMasks[sq][SIDE_BLACK].setBit(sqFromRowAndCol(row+1, col-1))
				}
			}
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Pinned Pieces Masks ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// combine the UD, LR and both diagonal rays from the previous table
// this will get the movement in each of the 4 main directions

func initMovePinnedPiecesMasks() {
	for sq := 0; sq < 64; sq++ {

		// UD
		newBitboardUD := emptyBB
		newBitboardUD |= moveRaysTable[sq][RAY_UP]
		newBitboardUD |= moveRaysTable[sq][RAY_DOWN]
		movePinnedMasksTable[sq][PIN_UD] = newBitboardUD

		// LR
		newBitboardLR := emptyBB
		newBitboardLR |= moveRaysTable[sq][RAY_LEFT]
		newBitboardLR |= moveRaysTable[sq][RAY_RIGHT]
		movePinnedMasksTable[sq][PIN_LR] = newBitboardLR

		// ULtDR
		newBitboardULtDR := emptyBB
		newBitboardULtDR |= moveRaysTable[sq][RAY_UL]
		newBitboardULtDR |= moveRaysTable[sq][RAY_DR]
		movePinnedMasksTable[sq][PIN_ULtDR] = newBitboardULtDR

		// DLtUR
		newBitboardDLtUR := emptyBB
		newBitboardDLtUR |= moveRaysTable[sq][RAY_DL]
		newBitboardDLtUR |= moveRaysTable[sq][RAY_UR]
		movePinnedMasksTable[sq][PIN_DLtUR] = newBitboardDLtUR
	}
}
