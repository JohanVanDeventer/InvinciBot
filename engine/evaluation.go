package main

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Background ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Evaluation is rounded to the nearest centipawn.
Pawns are worth 100 centipawns as reference.
*/

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Heatmap Tables --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Heatmap tables for each piece on each square, but blended later for the current game stage.
*/

// not used
var evalTableBLANK [8][8]int = [8][8]int{
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
}

// kings
var evalTableKingsMid8x8 [8][8]int = [8][8]int{
	{-30, -40, -40, -40, -40, -40, -40, -30},
	{-30, -40, -50, -50, -50, -50, -40, -30},
	{-30, -30, -40, -50, -50, -40, -40, -30},
	{-30, -30, -40, -50, -50, -40, -30, -30},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-10, -20, -20, -30, -30, -20, -20, -10},
	{+10, +10, -10, -20, -20, -10, +10, +10},
	{+15, +30, +10, -10, -10, -10, +40, +15},
}

var evalTableKingsEnd8x8 [8][8]int = [8][8]int{
	{-40, -30, -20, -20, -20, -20, -30, -40},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-20, 000, +10, +10, +10, +10, 000, -20},
	{-20, 000, +10, +20, +20, +10, 000, -20},
	{-20, 000, +10, +20, +20, +10, 000, -20},
	{-20, 000, +10, +10, +10, +10, 000, -20},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-40, -30, -20, -20, -20, -20, -30, -40},
}

// queens
var evalTableQueensMid8x8 [8][8]int = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, 000, 000, 000, 000, -10},
	{-20, -10, -10, +10, -10, -10, -10, -20},
}

var evalTableQueensEnd8x8 [8][8]int = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

// rooks
var evalTableRooksMid8x8 [8][8]int = [8][8]int{
	{000, 000, 000, 000, 000, 000, 000, 000},
	{+05, +10, +10, +10, +10, +10, +10, +05},
	{-05, 000, 000, 000, 000, 000, 000, -05},
	{-05, 000, 000, 000, 000, 000, 000, -05},
	{-05, 000, 000, 000, 000, 000, 000, -05},
	{-05, 000, 000, 000, 000, 000, 000, -05},
	{-05, 000, 000, 000, 000, 000, 000, -05},
	{000, 000, 000, +10, +10, 000, 000, 000},
}

var evalTableRooksEnd8x8 [8][8]int = [8][8]int{
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
}

// bishops
var evalTableBishopsMid8x8 [8][8]int = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-10, 000, 000, +05, +05, 000, 000, -10},
	{-10, +05, +05, +10, +10, +05, +05, -10},
	{-10, +05, +10, +10, +10, +10, +05, -10},
	{-10, +10, +10, +05, +05, +10, +10, -10},
	{-10, +10, 000, +05, +05, 000, +10, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

var evalTableBishopsEnd8x8 [8][8]int = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, +10, +10, +10, +10, 000, -10},
	{-10, 000, 000, 000, 000, 000, 000, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

// knights
var evalTableKnightsMid8x8 [8][8]int = [8][8]int{
	{-50, -30, -20, -20, -20, -20, -30, -50},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-20, 000, +10, +15, +15, +10, 000, -20},
	{-20, 000, +15, +20, +20, +15, 000, -20},
	{-20, 000, +15, +20, +20, +15, 000, -20},
	{-20, 000, +10, +15, +15, +10, 000, -20},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-50, -25, -20, -20, -20, -20, -25, -50},
}

var evalTableKnightsEnd8x8 [8][8]int = [8][8]int{
	{-50, -30, -20, -20, -20, -20, -30, -50},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-20, 000, +15, +15, +15, +15, 000, -20},
	{-20, 000, +15, +15, +15, +15, 000, -20},
	{-20, 000, +15, +15, +15, +15, 000, -20},
	{-20, 000, +15, +15, +15, +15, 000, -20},
	{-30, -10, 000, 000, 000, 000, -10, -30},
	{-50, -30, -20, -20, -20, -20, -30, -50},
}

// pawns
var evalTablePawnsMid8x8 [8][8]int = [8][8]int{
	{000, 000, 000, 000, 000, 000, 000, 000},
	{+40, +40, +40, +40, +40, +40, +40, +40},
	{000, 000, +15, +25, +25, 000, 000, 000},
	{000, 000, +15, +25, +25, 000, 000, -10},
	{000, 000, +10, +25, +25, 000, -10, -20},
	{000, 000, +05, +10, +10, -10, 000, 000},
	{000, 000, -10, -20, -20, +20, +10, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
}

var evalTablePawnsEnd8x8 [8][8]int = [8][8]int{
	{000, 000, 000, 000, 000, 000, 000, 000},
	{+60, +60, +60, +60, +60, +60, +60, +60},
	{+35, +35, +35, +35, +35, +35, +35, +35},
	{+25, +25, +25, +25, +25, +25, +25, +25},
	{+15, +15, +15, +15, +15, +15, +15, +15},
	{+10, +10, +10, +10, +10, +10, +10, +10},
	{000, 000, 000, 000, 000, 000, 000, 000},
	{000, 000, 000, 000, 000, 000, 000, 000},
}

// converted eval tables from above tables
var evalTableCombinedMid [2][6][64]int // for each piece, for each side, for each square
var evalTableCombinedEnd [2][6][64]int // for each piece, for each side, for each square

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Heatmap Tables Init -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// combine the above tables into a useable format for later lookup
func initEvalTables() {

	// kings: mid
	for rowIndex, row := range evalTableKingsMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_KING][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_KING][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// kings: end
	for rowIndex, row := range evalTableKingsEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_KING][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_KING][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// queens: mid
	for rowIndex, row := range evalTableQueensMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_QUEEN][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_QUEEN][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// queens: end
	for rowIndex, row := range evalTableQueensEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_QUEEN][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_QUEEN][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// rooks: mid
	for rowIndex, row := range evalTableRooksMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_ROOK][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_ROOK][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// rooks: end
	for rowIndex, row := range evalTableRooksEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_ROOK][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_ROOK][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// knights: mid
	for rowIndex, row := range evalTableKnightsMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_KNIGHT][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_KNIGHT][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// knights: end
	for rowIndex, row := range evalTableKnightsEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_KNIGHT][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_KNIGHT][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// bishops: mid
	for rowIndex, row := range evalTableBishopsMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_BISHOP][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_BISHOP][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// bishops: end
	for rowIndex, row := range evalTableBishopsEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_BISHOP][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_BISHOP][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// pawns: mid
	for rowIndex, row := range evalTablePawnsMid8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedMid[SIDE_WHITE][PIECE_PAWN][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedMid[SIDE_BLACK][PIECE_PAWN][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}

	// pawns: end
	for rowIndex, row := range evalTablePawnsEnd8x8 {
		for colIndex, value := range row {
			correctRowIndex := 7 - rowIndex

			// white side
			evalTableCombinedEnd[SIDE_WHITE][PIECE_PAWN][sqFromRowAndCol(correctRowIndex, colIndex)] = value

			// black side: invert rows but not columns, and also invert values (+ score for white is - score for black in absolute terms)
			evalTableCombinedEnd[SIDE_BLACK][PIECE_PAWN][sqFromRowAndCol(rowIndex, colIndex)] = -value
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ Eval: Basics ------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
The eval is always done from the white side (absolute value).
Positive is good for white, negative is good for black.

Some eval is done incrementally, some for each position from the start.
*/

const (
	// direct material value in centipawns
	VALUE_PAWN   int = 100
	VALUE_KNIGHT int = 400
	VALUE_BISHOP int = 420
	VALUE_ROOK   int = 600
	VALUE_QUEEN  int = 1200

	// stage of the game (mid vs end) value
	STAGE_VAL_QUEEN    int = 4
	STAGE_VAL_ROOK     int = 2
	STAGE_VAL_KNIGHT   int = 1
	STAGE_VAL_BISHOP   int = 1
	STAGE_VAL_STARTING int = STAGE_VAL_QUEEN*2 + STAGE_VAL_ROOK*4 + STAGE_VAL_KNIGHT*4 + STAGE_VAL_BISHOP*4 // normally 24

)

var evalTableMaterial [2][6]int // maps the side and piece type to their material values (black has negative values, kings are 0)
var evalTableGameStage [6]int   // maps the piece type to their game stage values (kings and pawns are 0)

func initEvalMaterialAndStageTables() {
	evalTableMaterial[SIDE_WHITE][PIECE_KING] = 0
	evalTableMaterial[SIDE_WHITE][PIECE_QUEEN] = VALUE_QUEEN
	evalTableMaterial[SIDE_WHITE][PIECE_ROOK] = VALUE_ROOK
	evalTableMaterial[SIDE_WHITE][PIECE_KNIGHT] = VALUE_KNIGHT
	evalTableMaterial[SIDE_WHITE][PIECE_BISHOP] = VALUE_BISHOP
	evalTableMaterial[SIDE_WHITE][PIECE_PAWN] = VALUE_PAWN

	evalTableMaterial[SIDE_BLACK][PIECE_KING] = 0
	evalTableMaterial[SIDE_BLACK][PIECE_QUEEN] = 0 - VALUE_QUEEN
	evalTableMaterial[SIDE_BLACK][PIECE_ROOK] = 0 - VALUE_ROOK
	evalTableMaterial[SIDE_BLACK][PIECE_KNIGHT] = 0 - VALUE_KNIGHT
	evalTableMaterial[SIDE_BLACK][PIECE_BISHOP] = 0 - VALUE_BISHOP
	evalTableMaterial[SIDE_BLACK][PIECE_PAWN] = 0 - VALUE_PAWN

	evalTableGameStage[PIECE_KING] = 0
	evalTableGameStage[PIECE_QUEEN] = STAGE_VAL_QUEEN
	evalTableGameStage[PIECE_ROOK] = STAGE_VAL_ROOK
	evalTableGameStage[PIECE_KNIGHT] = STAGE_VAL_KNIGHT
	evalTableGameStage[PIECE_BISHOP] = STAGE_VAL_BISHOP
	evalTableGameStage[PIECE_PAWN] = 0

}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------- Eval: Material, Game Stage, and Heatmaps -----------------------------------
// --------------------------------------------------------------------------------------------------------------------

// evaluate a fresh starting position
// therefore incremental updates during make move will be correct
func (pos *Position) evalPosAtStart() {

	pos.logTime.allLogTypes[LOG_ONCE_EVAL].start()

	// start with a zero eval for all eval variables
	pos.evalMaterial = 0
	pos.evalHeatmaps = 0
	pos.evalOther = 0
	pos.evalMidVsEndStage = 0

	// ------------------ MATERIAL VALUE + GAME STAGE VALUE -----------------
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {
			pieces := pos.pieces[side][pieceType]
			pieceCount := pieces.countBits()

			pos.evalMaterial += evalTableMaterial[side][pieceType] * pieceCount
			pos.evalMidVsEndStage += evalTableGameStage[pieceType] * pieceCount
		}
	}

	// ----------------------------- HEATMAP VALUE --------------------------
	// we only do this after the game stage value is determined above
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {

			// get the pieces bitboard
			pieces := pos.pieces[side][pieceType]
			for pieces != 0 {

				// get the next piece square
				nextPieceSq := pieces.popBitGetSq()

				// add the heatmap value of that piece on that square to the eval
				evalStage := pos.evalMidVsEndStage
				if evalStage > STAGE_VAL_STARTING { // cap to the max stage value
					evalStage = STAGE_VAL_STARTING
				}
				midValue := evalTableCombinedMid[side][pieceType][nextPieceSq]
				endValue := evalTableCombinedEnd[side][pieceType][nextPieceSq]
				pos.evalHeatmaps += ((midValue * evalStage) + (endValue * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING
			}
		}
	}

	pos.logTime.allLogTypes[LOG_ONCE_EVAL].stop()
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Eval: Doubled Pawns Setup ------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// setup for evaluating doubled pawns

const (
	DOUBLED_PAWN_PENALTY int = 10 // penalty for a pawn if there are other pawns on that column
)

var columnMasks [8]Bitboard // masks where all the bits for that column only is set

// create masks for each column on the board
func initEvalColumnMasks() {
	for col := 0; col < 8; col++ {
		newBitboard := emptyBB
		for row := 0; row < 8; row++ {
			newBitboard.setBit(sqFromRowAndCol(row, col))
		}
		columnMasks[col] = newBitboard
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------- Eval: Other --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// evalue a position after for non-incremental evaluations
func (pos *Position) evalPosAfter() {

	pos.logTime.allLogTypes[LOG_EVAL].start()

	/*
		// reset other evaluation scores
		pos.evalOther = 0

		// ------------------------------------------------- DOUBLED PAWNS --------------------------------------------------
		// white pawns
		for col := 0; col < 8; col++ {
			colMask := columnMasks[col]                                 // get the column mask
			maskedPawns := colMask & pos.pieces[SIDE_WHITE][PIECE_PAWN] // get the pawns on that mask
			pawnsOnColCount := maskedPawns.countBits()                  // count the pawns
			if pawnsOnColCount > 1 {                                    // if there are more than 1 pawn, we have doubled pawns
				pos.evalOther -= DOUBLED_PAWN_PENALTY * pawnsOnColCount
			}
		}

		// black pawns
		for col := 0; col < 8; col++ {
			colMask := columnMasks[col]                                 // get the column mask
			maskedPawns := colMask & pos.pieces[SIDE_BLACK][PIECE_PAWN] // get the pawns on that mask
			pawnsOnColCount := maskedPawns.countBits()                  // count the pawns
			if pawnsOnColCount > 1 {                                    // if there are more than 1 pawn, we have doubled pawns
				pos.evalOther += DOUBLED_PAWN_PENALTY * pawnsOnColCount
			}
		}
	*/

	pos.logTime.allLogTypes[LOG_EVAL].stop()

}
