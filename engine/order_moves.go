package main

import (
	"sort"
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------ Order Moves: Background -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

Order moves by the expected most valuable to least valuable for earlier search cutoffs.
Order moves by:
- 1. Promotions (direct threat to increase side to move value by normally a queen).
- 2. Most valuable attacker captures least valuable defender (for example pawn takes queen).
- 3. Captures and en-passant before quiet moves (generally expect to need to search more deeply and therefore earlier).
- 4. Castling (could get out of threats).

Other quiet moves are therefore last.

Ordering of moves returns a copy of the current available moves. It does not modify anything in the available moves array.

*/

const (
	MOVE_ORD_NONQUIET_BONUS int = 2000 // any non-quiet moves are evaluated before quiet moves
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Order Moves: Functions -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// returns a slice of moves ordered from best to worst

func (pos *Position) getOrderedMoves() []Move {

	start_time := time.Now()

	// if there are no moves, return an empty slice early
	if pos.availableMovesCounter == 0 {
		var blankMoves []Move
		return blankMoves
	}

	// create a copy of the available moves
	moves := pos.availableMoves[:pos.availableMovesCounter]

	// loop over moves to score them
	for i, move := range moves {

		// set the score of the move to zero
		moveOrderScore := 0

		// get the relevant move information
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()

		// add a bonus for captures, en-passant and promotions (non-quiet moves)
		if moveType == MOVE_TYPE_CAPTURE || moveType == MOVE_TYPE_EN_PASSANT || promotionType != PROMOTION_NONE {
			moveOrderScore += MOVE_ORD_NONQUIET_BONUS
		}

		// if there is a promotion, add that promoted piece's value less the pawn value
		// the score is done from the white side (positve score), because for move ordering this is the point
		if promotionType != PROMOTION_NONE {
			moveOrderScore += (evalTableMaterial[SIDE_WHITE][promotionType] - VALUE_PAWN)
		}

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

		// finally set the updated move score if the scorew is not zero
		if moveOrderScore != 0 {
			moves[i].setMoveOrderingScore(moveOrderScore)
		}
	}

	// now sort the moves
	// define the custom comparator function
	// sort the slice based on the score field using the comparator function
	sort.Slice(moves, func(i, j int) bool { return moves[i].getMoveOrderingScore() > moves[j].getMoveOrderingScore() })

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_ORDER_MOVES].addTime(int(duration_time))

	return moves
}
