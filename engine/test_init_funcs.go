package main

import "fmt"

// debug function to test whether engine tables have initialized correctly
// not used
func printInitTestResults() {

	fmt.Println(" ")
	fmt.Println("------------------------ Table Initialization Results --------------------------")

	// bb reference array
	fmt.Println("---- bbReferenceArray ----")
	printBBReferenceArray()

	// move tables
	fmt.Println("---- Move Table: King on 27, 56 and 15 ----")
	moveKingsTable[27].printBitboardFancy8x8()
	fmt.Println(" ")
	moveKingsTable[56].printBitboardFancy8x8()
	fmt.Println(" ")
	moveKingsTable[15].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Move Table: Knight on 27, 56 and 15 ----")
	moveKnightsTable[27].printBitboardFancy8x8()
	fmt.Println(" ")
	moveKnightsTable[56].printBitboardFancy8x8()
	fmt.Println(" ")
	moveKnightsTable[15].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Move Table: Pawn on 12, 20, 52 (WHITE)----")
	moveOnlyPawnsTable[12][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")
	moveOnlyPawnsTable[20][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")
	moveOnlyPawnsTable[52][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Attack Table: Pawn on 12, 20, 52 (WHITE)----")
	moveAttackPawnsTable[12][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")
	moveAttackPawnsTable[20][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")
	moveAttackPawnsTable[52][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Move Table: Pawn on 44, 52, 12 (BLACK) ----")
	moveOnlyPawnsTable[44][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")
	moveOnlyPawnsTable[52][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")
	moveOnlyPawnsTable[12][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Attack Table: Pawn on 44, 52, 12 (BLACK) ----")
	moveAttackPawnsTable[44][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")
	moveAttackPawnsTable[52][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")
	moveAttackPawnsTable[12][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Double Move: Pawn on 12 (WHITE) ----")
	movePawnDoubleMasks[12][SIDE_WHITE].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Double Move: Pawn on 52 (BLACK) ----")
	movePawnDoubleMasks[52][SIDE_BLACK].printBitboardFancy8x8()
	fmt.Println(" ")

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 UP ----")
	moveRaysTable[27][RAY_UP].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 RIGHT ----")
	moveRaysTable[27][RAY_RIGHT].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 DOWN ----")
	moveRaysTable[27][RAY_DOWN].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 LEFT ----")
	moveRaysTable[27][RAY_LEFT].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 UL ----")
	moveRaysTable[27][RAY_UL].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 UR ----")
	moveRaysTable[27][RAY_UR].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 DR ----")
	moveRaysTable[27][RAY_DR].printBitboardFancy8x8()

	fmt.Println("---- Move Table Single Ray: Ray Mover 27 DL ----")
	moveRaysTable[27][RAY_DL].printBitboardFancy8x8()

	fmt.Println("---- Move Table Combined: Rook on 27 ----")
	moveRooksTable[27].printBitboardFancy8x8()

	fmt.Println("---- Move Table Combined: Rook on 0 ----")
	moveRooksTable[0].printBitboardFancy8x8()

	fmt.Println("---- Move Table Combined: Bishop on 27 ----")
	moveBishopsTable[27].printBitboardFancy8x8()

	fmt.Println("---- Move Table Combined: Bishop on 0 ----")
	moveBishopsTable[0].printBitboardFancy8x8()

	// castling mask tables
	fmt.Println("---- Castling Mask: K ----")
	moveCastlingIsClearMasks[CASTLE_WHITE_KINGSIDE].printBitboardFancy8x8()

	fmt.Println("---- Castling Mask: Q ----")
	moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE].printBitboardFancy8x8()

	fmt.Println("---- Castling Mask: k ----")
	moveCastlingIsClearMasks[CASTLE_BLACK_KINGSIDE].printBitboardFancy8x8()

	fmt.Println("---- Castling Mask: q ----")
	moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE].printBitboardFancy8x8()

	// pawns attacking king masks
	fmt.Println("---- White Pawns Attacking King on 27 Mask ----")
	movePawnsAttackingKingMasks[27][SIDE_WHITE].printBitboardFancy8x8()

	fmt.Println("---- Black Pawns Attacking King on 27 Mask ----")
	movePawnsAttackingKingMasks[27][SIDE_BLACK].printBitboardFancy8x8()

	// pinned pice movement masks
	fmt.Println("---- Pinned Pice on 27: UP DOWN Mask ----")
	movePinnedMasksTable[27][PIN_UD].printBitboardFancy8x8()

	fmt.Println("---- Pinned Pice on 27: LEFT RIGHT Mask ----")
	movePinnedMasksTable[27][PIN_LR].printBitboardFancy8x8()

	fmt.Println("---- Pinned Pice on 27: ULtDR Mask ----")
	movePinnedMasksTable[27][PIN_ULtDR].printBitboardFancy8x8()

	fmt.Println("---- Pinned Pice on 27: DLtUR Mask ----")
	movePinnedMasksTable[27][PIN_DLtUR].printBitboardFancy8x8()

	// hash tables
	fmt.Println("---- Hash Table: hashTablePieces ----")
	fmt.Printf("%v\n", hashTablePieces)

	fmt.Println("---- Hash Table: hashTableCastling ----")
	fmt.Printf("%v\n", hashTableCastling)

	fmt.Println("---- Hash Table: hashTableSideToMove ----")
	fmt.Printf("%v\n", hashTableSideToMove)

	fmt.Println("---- Hash Table: hashTableEnPassant ----")
	fmt.Printf("%v\n", hashTableEnPassant)

	fmt.Println("---- Hash Table: startingHash ----")
	fmt.Printf("%v\n", startingHash)

	fmt.Println("---- Hash Table: Hash Collision Check Table ----")
	fmt.Printf("%v\n", hashCollisionsStack)

	// eval tables: kings
	fmt.Println("---- Combined Eval Table: Mid: White King: Sq 4 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_KING][4])

	fmt.Println("---- Combined Eval Table: Mid: Black King: Sq 4 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_KING][4])

	fmt.Println("---- Combined Eval Table: End: White King: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_KING][60])

	fmt.Println("---- Combined Eval Table: End: Black King: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_KING][60])

	// eval tables: queens
	fmt.Println("---- Combined Eval Table: Mid: White Queen: Sq 4 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_QUEEN][4])

	fmt.Println("---- Combined Eval Table: Mid: Black Queen: Sq 4 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_QUEEN][4])

	fmt.Println("---- Combined Eval Table: End: White Queen: Sq 63 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_QUEEN][63])

	fmt.Println("---- Combined Eval Table: End: Black Queen: Sq 63 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_QUEEN][63])

	// eval tables: rooks
	fmt.Println("---- Combined Eval Table: Mid: White Rook: Sq 24 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_ROOK][24])

	fmt.Println("---- Combined Eval Table: Mid: Black Rook: Sq 24 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_ROOK][24])

	fmt.Println("---- Combined Eval Table: End: White Rook: Sq 27 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_ROOK][27])

	fmt.Println("---- Combined Eval Table: End: Black Rook: Sq 27 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_ROOK][27])

	// eval tables: knights
	fmt.Println("---- Combined Eval Table: Mid: White Knight: Sq 16 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_KNIGHT][16])

	fmt.Println("---- Combined Eval Table: Mid: Black Knight: Sq 16 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_KNIGHT][16])

	fmt.Println("---- Combined Eval Table: End: White Knight: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_KNIGHT][60])

	fmt.Println("---- Combined Eval Table: End: Black Knight: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_KNIGHT][60])

	// eval tables: bishops
	fmt.Println("---- Combined Eval Table: Mid: White Bishop: Sq 27 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_BISHOP][27])

	fmt.Println("---- Combined Eval Table: Mid: Black Bishop: Sq 27 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_BISHOP][27])

	fmt.Println("---- Combined Eval Table: End: White Bishop: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_BISHOP][60])

	fmt.Println("---- Combined Eval Table: End: Black Bishop: Sq 60 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_BISHOP][60])

	// eval tables: pawns
	fmt.Println("---- Combined Eval Table: Mid: White Pawn: Sq 28 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_WHITE][PIECE_PAWN][28])

	fmt.Println("---- Combined Eval Table: Mid: Black Pawn: Sq 28 ----")
	fmt.Printf("%v\n", evalTableCombinedMid[SIDE_BLACK][PIECE_PAWN][28])

	fmt.Println("---- Combined Eval Table: End: White Pawn: Sq 33 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_WHITE][PIECE_PAWN][33])

	fmt.Println("---- Combined Eval Table: End: Black Pawn: Sq 33 ----")
	fmt.Printf("%v\n", evalTableCombinedEnd[SIDE_BLACK][PIECE_PAWN][33])
}
