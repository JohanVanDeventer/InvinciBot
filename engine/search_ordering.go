package main

import (
	"sort"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Score and Order Moves -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	// we add this offset to all threat moves to always have a positive score
	// because we encode the score into the move, so having negative numbers will mess that up
	// the max we can gain or lose from a move is 900 + 800 = 1700 (promote a pawn to a queen while capturing a queen)
	MOVE_ORDERING_SCORE_OFFSET int = 4096
)

// returns a slice of threat moves ordered from best to worst
func (pos *Position) getOrderedThreatMoves() []Move {

	// create a copy of the available moves
	moves := make([]Move, pos.threatMovesCounter)
	copy(moves, pos.threatMoves[:pos.threatMovesCounter])

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := MOVE_ORDERING_SCORE_OFFSET

		// get the relevant move information
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

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

	// finally clear the move ordering scores (to make comparison easier later to other moves)
	for _, move := range moves {
		move.clearMoveOrderingScore()
	}

	return moves
}
