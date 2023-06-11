package main

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ Game State --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Returns the state of the current position, which can be:
- game ongoing
- checkmate (white or black wins)
- draw (stalemate, 3-fold repetition, or 50-move rule)
*/

const (
	STATE_ONGOING                int = 0
	STATE_WIN_WHITE              int = 1
	STATE_WIN_BLACK              int = 2
	STATE_DRAW_STALEMATE         int = 3
	STATE_DRAW_50_MOVE_RULE      int = 4
	STATE_DRAW_3_FOLD_REPETITION int = 5
)

// only call this after generating moves in a position
// because it uses the number of available moves to determine checkmate and stalemate
func (pos *Position) getGameStateAndStore() {

	pos.logTime.allLogTypes[LOG_GAME_STATE].start()
	defer pos.logTime.allLogTypes[LOG_GAME_STATE].stop()

	// if there are no moves remaining, the king is either in checkmate, or it's stalemate
	if pos.totalMovesCounter == 0 {
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

	// game ongoing if noting else applies
	pos.gameState = STATE_ONGOING

}
