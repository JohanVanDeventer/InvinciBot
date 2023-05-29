package main

import (
	"sort"
	"time"
)

const (
	MOVE_ORD_NONQUIET_BONUS int = 2000 // any non-quiet moves are evaluated before quiet moves
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Score and Order Moves -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
We have less intensive move ordering at other nodes than the root.
*/

// returns a slice of moves ordered from best to worst
func (pos *Position) getScoredAndOrderedThreatMoves() []Move {

	start_time := time.Now()

	// create a copy of the available moves
	moves := make([]Move, pos.threatMovesCounter)
	copy(moves, pos.threatMoves[:pos.threatMovesCounter])

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

		// ------------------------------------------- BONUS FOR CAPTURES AND PROMOTIONS -------------------------------------
		// add a bonus for captures, en-passant and promotions (threat moves)
		moveOrderScore += MOVE_ORD_NONQUIET_BONUS

		// ------------------------------------------- PROMOTIONS: VALUE GAINED -------------------------------------
		// if there is a promotion, add that promoted piece's value less the pawn value
		// the score is done from the white side (positve score), because for move ordering this is the point
		if promotionType != PROMOTION_NONE {
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][promotionType] - VALUE_PAWN)
		}

		// ------------------------------------------- CAPTURES: VALUE GAINED -------------------------------------
		// MOVE_TYPE_EN_PASSANT: no bonus because pawn traded for pawn = 0 incremental value
		// MOVE_TYPE_CAPTURE: evaluate below
		if moveType == MOVE_TYPE_CAPTURE {
			var enemyPiece int = 6 // set to 6 to catch bugs
			var enSide int = SIDE_WHITE
			if pos.isWhiteTurn {
				enSide = SIDE_BLACK
			}
			toSq := move.getToSq()
			piece := move.getPiece()

			if pos.pieces[enSide][PIECE_PAWN].isBitSet(toSq) {
				enemyPiece = PIECE_PAWN
			} else if pos.pieces[enSide][PIECE_KNIGHT].isBitSet(toSq) {
				enemyPiece = PIECE_KNIGHT
			} else if pos.pieces[enSide][PIECE_BISHOP].isBitSet(toSq) {
				enemyPiece = PIECE_BISHOP
			} else if pos.pieces[enSide][PIECE_ROOK].isBitSet(toSq) {
				enemyPiece = PIECE_ROOK
			} else if pos.pieces[enSide][PIECE_QUEEN].isBitSet(toSq) {
				enemyPiece = PIECE_QUEEN
			}

			// add the difference between the captured and friendly piece
			// therefore lower piece value captures higher piece value is evaluated first
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][enemyPiece]) - (evalTableMaterial[SIDE_WHITE][piece])
		}

		// ------------------------------------------- SAVE THE SCORE -------------------------------------
		// finally set the updated move score if the score is not zero
		moves[i].setMoveOrderingScore(moveOrderScore)
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return moves[i].getMoveOrderingScore() > moves[j].getMoveOrderingScore() })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES_NOT_AT_ROOT].addTime(int(duration_time))

	return moves
}

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Score and Order Moves -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
We have less intensive move ordering at other nodes than the root.

// returns a slice of moves ordered from best to worst
func (pos *Position) getOrderedMoves() []Move {

	start_time := time.Now()

	// create a copy of the available moves
	moves := make([]Move, pos.availableMovesCounter)
	copy(moves, pos.availableMoves[:pos.availableMovesCounter])

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

		// ------------------------------------------- BONUS FOR CAPTURES AND PROMOTIONS -------------------------------------
		// add a bonus for captures, en-passant and promotions (non-quiet moves)
		if moveType == MOVE_TYPE_CAPTURE || moveType == MOVE_TYPE_EN_PASSANT || promotionType != PROMOTION_NONE {
			moveOrderScore += MOVE_ORD_NONQUIET_BONUS
		}

		// ------------------------------------------- PROMOTIONS: VALUE GAINED -------------------------------------
		// if there is a promotion, add that promoted piece's value less the pawn value
		// the score is done from the white side (positve score), because for move ordering this is the point
		if promotionType != PROMOTION_NONE {
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][promotionType] - VALUE_PAWN)
		}

		// ------------------------------------------- CAPTURES: VALUE GAINED -------------------------------------
		// MOVE_TYPE_EN_PASSANT: no bonus because pawn traded for pawn = 0 incremental value
		// MOVE_TYPE_CAPTURE: evaluate below
		if moveType == MOVE_TYPE_CAPTURE {
			var enemyPiece int = 6 // set to 6 to catch bugs
			var enSide int = SIDE_WHITE
			if pos.isWhiteTurn {
				enSide = SIDE_BLACK
			}
			toSq := move.getToSq()
			piece := move.getPiece()

			if pos.pieces[enSide][PIECE_PAWN].isBitSet(toSq) {
				enemyPiece = PIECE_PAWN
			} else if pos.pieces[enSide][PIECE_KNIGHT].isBitSet(toSq) {
				enemyPiece = PIECE_KNIGHT
			} else if pos.pieces[enSide][PIECE_BISHOP].isBitSet(toSq) {
				enemyPiece = PIECE_BISHOP
			} else if pos.pieces[enSide][PIECE_ROOK].isBitSet(toSq) {
				enemyPiece = PIECE_ROOK
			} else if pos.pieces[enSide][PIECE_QUEEN].isBitSet(toSq) {
				enemyPiece = PIECE_QUEEN
			}

			// add the difference between the captured and friendly piece
			// therefore lower piece value captures higher piece value is evaluated first
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][enemyPiece]) - (evalTableMaterial[SIDE_WHITE][piece])
		}

		// ------------------------------------------- SAVE THE SCORE -------------------------------------
		// finally set the updated move score if the score is not zero
		moves[i].setMoveOrderingScore(moveOrderScore)
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return moves[i].getMoveOrderingScore() > moves[j].getMoveOrderingScore() })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES_NOT_AT_ROOT].addTime(int(duration_time))

	return moves
}
*/

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------------- SEE ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
SEE (static exchange evaluation).
*/

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Order Moves: At Root --------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
We have more detailed move ordering at the root than for other nodes, because it is called much less.
At root nodes, we score and sort the whole list upfront before starting the search.

// returns a slice of moves ordered from best to worst
func (pos *Position) getOrderedMovesAtRoot() []Move {

	start_time := time.Now()

	// create a copy of the available moves
	moves := make([]Move, pos.availableMovesCounter)
	copy(moves, pos.availableMoves[:pos.availableMovesCounter])

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		fromSq := move.getFromSq()
		toSq := move.getToSq()
		piece := move.getPiece()
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

		// ------------------------------------------- BONUS FOR CAPTURES AND PROMOTIONS -------------------------------------
		// add a bonus for captures, en-passant and promotions (non-quiet moves)
		if moveType == MOVE_TYPE_CAPTURE || moveType == MOVE_TYPE_EN_PASSANT || promotionType != PROMOTION_NONE {
			moveOrderScore += MOVE_ORD_NONQUIET_BONUS
		}

		// ------------------------------------------- PROMOTIONS: VALUE GAINED -------------------------------------
		// if there is a promotion, add that promoted piece's value less the pawn value
		// the score is done from the white side (positve score), because for move ordering this is the point
		if promotionType != PROMOTION_NONE {
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][promotionType] - VALUE_PAWN)
		}

		// ------------------------------------------- CAPTURES: VALUE GAINED -------------------------------------
		// MOVE_TYPE_EN_PASSANT: no bonus because pawn traded for pawn = 0 incremental value
		// MOVE_TYPE_CAPTURE: evaluate below
		if moveType == MOVE_TYPE_CAPTURE {
			var enemyPiece int = 6 // set to 6 to catch bugs
			var enSide int = SIDE_WHITE
			if pos.isWhiteTurn {
				enSide = SIDE_BLACK
			}
			toSq := move.getToSq()
			piece := move.getPiece()

			if pos.pieces[enSide][PIECE_PAWN].isBitSet(toSq) {
				enemyPiece = PIECE_PAWN
			} else if pos.pieces[enSide][PIECE_KNIGHT].isBitSet(toSq) {
				enemyPiece = PIECE_KNIGHT
			} else if pos.pieces[enSide][PIECE_BISHOP].isBitSet(toSq) {
				enemyPiece = PIECE_BISHOP
			} else if pos.pieces[enSide][PIECE_ROOK].isBitSet(toSq) {
				enemyPiece = PIECE_ROOK
			} else if pos.pieces[enSide][PIECE_QUEEN].isBitSet(toSq) {
				enemyPiece = PIECE_QUEEN
			}

			// add the difference between the captured and friendly piece
			// therefore lower piece value captures higher piece value is evaluated first
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][enemyPiece]) - (evalTableMaterial[SIDE_WHITE][piece])
		}

		// ------------------------------------------- ALL MOVES: HEATMAP VALUE GAINED -------------------------------------
		// we add the heatmap value gained for all moves
		var heatmapValBefore int
		var heatmapValAfter int
		evalStage := pos.evalMidVsEndStage

		if pos.isWhiteTurn {
			midBefore := evalTableCombinedMid[SIDE_WHITE][piece][fromSq]
			endBefore := evalTableCombinedEnd[SIDE_WHITE][piece][fromSq]
			midAfter := evalTableCombinedMid[SIDE_WHITE][piece][toSq]
			endAfter := evalTableCombinedEnd[SIDE_WHITE][piece][toSq]

			heatmapValBefore = ((midBefore * evalStage) + (endBefore * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING
			heatmapValAfter = ((midAfter * evalStage) + (endAfter * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING

		} else {
			midBefore := 0 - evalTableCombinedMid[SIDE_BLACK][piece][fromSq]
			endBefore := 0 - evalTableCombinedEnd[SIDE_BLACK][piece][fromSq]
			midAfter := 0 - evalTableCombinedMid[SIDE_BLACK][piece][toSq]
			endAfter := 0 - evalTableCombinedEnd[SIDE_BLACK][piece][toSq]

			heatmapValBefore = ((midBefore * evalStage) + (endBefore * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING
			heatmapValAfter = ((midAfter * evalStage) + (endAfter * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING
		}

		moveOrderScore += heatmapValAfter - heatmapValBefore

		// ------------------------------------------- SAVE THE SCORE -------------------------------------
		// finally set the updated move score if the score is not zero
		moves[i].setMoveOrderingScore(moveOrderScore)
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return moves[i].getMoveOrderingScore() > moves[j].getMoveOrderingScore() })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES_AT_ROOT].addTime(int(duration_time))

	return moves
}

// scores the moves in a move list, but does not actually sort them yet
func (pos *Position) scoreMovesNotAtRoot(moves []Move) {

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

		// ------------------------------------------- BONUS FOR CAPTURES AND PROMOTIONS -------------------------------------
		// add a bonus for captures, en-passant and promotions (non-quiet moves)
		if moveType == MOVE_TYPE_CAPTURE || moveType == MOVE_TYPE_EN_PASSANT || promotionType != PROMOTION_NONE {
			moveOrderScore += MOVE_ORD_NONQUIET_BONUS
		}

		// ------------------------------------------- PROMOTIONS: VALUE GAINED -------------------------------------
		// if there is a promotion, add that promoted piece's value less the pawn value
		// the score is done from the white side (positve score), because for move ordering this is the point
		if promotionType != PROMOTION_NONE {
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][promotionType] - VALUE_PAWN)
		}

		// ------------------------------------------- CAPTURES: VALUE GAINED -------------------------------------
		// MOVE_TYPE_EN_PASSANT: no bonus because pawn traded for pawn = 0 incremental value
		// MOVE_TYPE_CAPTURE: evaluate below
		if moveType == MOVE_TYPE_CAPTURE {
			var enemyPiece int = 6 // set to 6 to catch bugs
			var enSide int = SIDE_WHITE
			if pos.isWhiteTurn {
				enSide = SIDE_BLACK
			}
			toSq := move.getToSq()
			piece := move.getPiece()

			if pos.pieces[enSide][PIECE_PAWN].isBitSet(toSq) {
				enemyPiece = PIECE_PAWN
			} else if pos.pieces[enSide][PIECE_KNIGHT].isBitSet(toSq) {
				enemyPiece = PIECE_KNIGHT
			} else if pos.pieces[enSide][PIECE_BISHOP].isBitSet(toSq) {
				enemyPiece = PIECE_BISHOP
			} else if pos.pieces[enSide][PIECE_ROOK].isBitSet(toSq) {
				enemyPiece = PIECE_ROOK
			} else if pos.pieces[enSide][PIECE_QUEEN].isBitSet(toSq) {
				enemyPiece = PIECE_QUEEN
			}

			// add the difference between the captured and friendly piece
			// therefore lower piece value captures higher piece value is evaluated first
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][enemyPiece]) - (evalTableMaterial[SIDE_WHITE][piece])
		}

		// ------------------------------------------- SAVE THE SCORE -------------------------------------
		// finally set the updated move score if the score is not zero
		moves[i].setMoveOrderingScore(moveOrderScore)
	}
}

// function to get the next best move from a list of moves
// it searches for the best score, and moves that move to the front of the move list
// finally it returns the best move
func (pos *Position) getNextBestMove(moves []Move, movesTriedBefore int) Move {
	bestScore := -1000000
	bestIndex := -99

	// loop over the remaining moves and get the best score
	for i, move := range moves[movesTriedBefore:] {
		moveScore := move.getMoveOrderingScore()
		if moveScore > bestScore {
			bestScore = moveScore
			bestIndex = i
		}
	}

	// get the best move
	bestMove := moves[bestIndex]

	// now move the best move to the front of the move list
	moves = append(moves[:bestIndex], moves[bestIndex+1:]...) // remove the best move from the original position
	moves = append([]Move{bestMove}, moves...)                // append the best move at the start of the list

	// finally, return the best move
	return bestMove
}
*/
