package main

import (
	"sort"
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Order Moves: Background -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
We order moves to try and get earlier cutoffs during the search.
*/

const (
	MOVE_ORD_NONQUIET_BONUS int = 2000 // any non-quiet moves are evaluated before quiet moves
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Order Moves: At Root --------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// we have more intensive move ordering at the root than for other nodes

// returns a slice of moves ordered from best to worst
func (pos *Position) getOrderedMovesAtRoot() []Move {

	start_time := time.Now()

	// create a copy of the available moves
	moves := make([]Move, pos.availableMovesCounter)
	copy(moves, pos.availableMoves[:pos.availableMovesCounter])

	// create a score list
	scores := make([]int, pos.availableMovesCounter)

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		//fromSq := move.getFromSq()
		//toSq := move.getToSq()
		//piece := move.getPiece()
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

			// ------------------------------------------- QUIET MOVES: HEATMAP VALUE GAINED -------------------------------------
			// we add the heatmap value gained for quiet moves
		}

		// finally set the updated move score if the score is not zero
		scores[i] = moveOrderScore
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return scores[i] > scores[j] })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES_AT_ROOT].addTime(int(duration_time))

	return moves
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Order Moves: Not At Root ------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// we have less intensive move ordering at other nodes than the root

// returns a slice of moves ordered from best to worst
func (pos *Position) getOrderedMovesNotAtRoot() []Move {

	start_time := time.Now()

	// create a copy of the available moves
	moves := make([]Move, pos.availableMovesCounter)
	copy(moves, pos.availableMoves[:pos.availableMovesCounter])

	// create a score list
	scores := make([]int, pos.availableMovesCounter)

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

		// finally set the updated move score if the score is not zero
		scores[i] = moveOrderScore
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return scores[i] > scores[j] })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES_NOT_AT_ROOT].addTime(int(duration_time))

	return moves
}
