package main

import (
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Heatmap Tables --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

Evaluation is rounded to the nearest centipawn.
Pawns are worth 100 centipawns.

Heatmap tables for each piece on each square, but blended for the current game stage.
Certain pieces have different heatmaps for the endgame and the middlegame (such as the king near the middle in the endgame).

*/

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

// kings: need to castle early, but in the endgame get closer to the middle
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

// queens: same table for mid and endgame, just try to not go near the edge
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

// rooks: favour the 7th rank of opponent, and also incentivize to castle, and generally not favour the sides
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

// bishops: favour the centre and not the sides, and also favout fianchetto a bit
// also favour sitting in front of pawns as protected on the 3rd rank
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

// knights: get the biggest penalties for on the edge, and also a bit bigger bonus to try move out before bishops
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

// pawns: try move 2 squares, and also don't move pawns in front of king; also bonus for close to promotion
// in the endgame try push pawns from the 2nd rank (but not too hard in middlegame)
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
// ------------------------------------------------------ Eval: Summary -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

/*

The eval is always done from the white side (absolute value).
Positive is good for white, negative is good for black.

The following is done during make/undo move:
- Pure material count.
- Game stage count (mid vs endgame).
because it is faster to do incrementally.

However, other evaluations is done for each position as required.

*/

const (
	// direct material value in centipawns
	VALUE_PAWN   int = 100
	VALUE_KNIGHT int = 300
	VALUE_BISHOP int = 315
	VALUE_ROOK   int = 510
	VALUE_QUEEN  int = 900

	// stage of the game (mid vs end) value
	STAGE_VAL_QUEEN    int = 4
	STAGE_VAL_ROOK     int = 2
	STAGE_VAL_KNIGHT   int = 1
	STAGE_VAL_BISHOP   int = 1
	STAGE_VAL_STARTING int = STAGE_VAL_QUEEN*2 + STAGE_VAL_ROOK*4 + STAGE_VAL_KNIGHT*2 + STAGE_VAL_BISHOP*2 // normally 24
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
// ------------------------------------------ Eval: Material Value and Game Stage -------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// evaluate a fresh starting position for things done incrementally during making moves later
func (pos *Position) evalPosAtStart() {

	// start with a zero eval
	pos.evalMaterial = 0
	pos.evalMidVsEndStage = 0

	// get the game stage value and material value
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {
			pieces := pos.pieces[side][pieceType]
			pieceCount := pieces.countBits()

			pos.evalMaterial += evalTableMaterial[side][pieceType] * pieceCount
			pos.evalMidVsEndStage += evalTableGameStage[pieceType] * pieceCount
		}
	}
}

// evalue a position after for non-incremental evaluations
func (pos *Position) evalPosAfter() {

	start_time := time.Now()

	// reset other evaluation scores
	pos.evalHeatmaps = 0
	pos.evalOther = 0

	// evaluation per piece: heatmap value
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {

			// get the pieces bitboard
			pieces := pos.pieces[side][pieceType]
			for pieces != 0 {

				// get the next piece square and get the heatmap values of that square
				nextPieceSq := pieces.popBitGetSq()

				if pos.evalMidVsEndStage >= 24 { // if the game stage is still in the opening, no weighting is needed, use the mid table
					pos.evalHeatmaps += evalTableCombinedMid[side][pieceType][nextPieceSq]

				} else { // else, need to weight based on the game stage
					midValue := evalTableCombinedMid[side][pieceType][nextPieceSq]
					endValue := evalTableCombinedEnd[side][pieceType][nextPieceSq]
					pos.evalHeatmaps += ((midValue * pos.evalMidVsEndStage) + (endValue * (STAGE_VAL_STARTING - pos.evalMidVsEndStage))) / STAGE_VAL_STARTING
				}
			}
		}
	}

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_EVAL].addTime(int(duration_time))
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ Game State --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// returns the state of the current position, which can be:
// - game ongoing
// - checkmate (white or black wins)
// - draw (stalemate, 3-fold repetition, or 50-move rule)

const (
	STATE_ONGOING                int = 0
	STATE_WIN_WHITE              int = 1
	STATE_WIN_BLACK              int = 2
	STATE_DRAW_STALEMATE         int = 3
	STATE_DRAW_50_MOVE_RULE      int = 4
	STATE_DRAW_3_FOLD_REPETITION int = 5
)

// only call this after generating moves in a position
func (pos *Position) getGameStateAndStore() {

	start_time := time.Now()
	defer pos.logOther.allLogTypes[LOG_GAME_STATE].addTime(int(time.Since(start_time).Nanoseconds()))

	// if there are no moves remaining, the king is either in checkmate, or it's a stalemate
	if pos.availableMovesCounter == 0 {
		if pos.kingChecks > 0 {
			if !pos.isWhiteTurn {
				pos.gameState = STATE_WIN_WHITE
				return
			} else {
				pos.gameState = STATE_WIN_BLACK
				return
			}
		} else {
			pos.gameState = STATE_DRAW_STALEMATE
			return
		}
	}

	// if there are moves, it's ongoing, unless the 50-move rule applies, or there is a 3-fold repetition
	// 50-move rule: remember it's 50 full moves, so 100 plies
	if pos.halfMoves > 100 {
		pos.gameState = STATE_DRAW_50_MOVE_RULE
		return
	}

	// 3-fold repetition
	countOfOccurences := 1 // current pos is the 1st occurence
	for _, previousHash := range pos.previousHashes[:pos.previousHashesCounter] {
		if pos.hashOfPos == previousHash {
			countOfOccurences += 1
		}
	}
	if countOfOccurences >= 3 {
		pos.gameState = STATE_DRAW_3_FOLD_REPETITION
		return
	}

	// game ongoing
	pos.gameState = STATE_ONGOING

}
