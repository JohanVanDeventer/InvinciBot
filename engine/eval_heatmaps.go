package main

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
