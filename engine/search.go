package main

import (
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Search: Background ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
The following things are applied during search:
- Negamax (alpha-beta in one function)
- Iterative deepening (searches at increasing depth until the time limit is reached)
- Transposition table (allows lookups to previously encountered positions)
- Quiescence search (search captures only until a position is quiet)
*/

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------- Quiescence Limits ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Depening on the initial depth we call the search from, we have a certain quiescence depth.
Therefore deeper depth searches have longer quiescence depths.
*/

const (
	MAX_DEPTH int = 99 // we set a max depth for the search (otherwise messes with assigning a best move)
)

var (
	qsDepthLimitTable [MAX_DEPTH + 1]int
)

func initQSDepthLimits() {
	for depth := 0; depth <= MAX_DEPTH; depth++ {

		// important for finding at least 1 best move that only depths > 2 have qs
		if depth <= 2 {
			qsDepthLimitTable[depth] = 0

		} else if depth <= 4 {
			qsDepthLimitTable[depth] = -2

		} else if depth <= 6 {
			qsDepthLimitTable[depth] = -4

		} else if depth <= 8 {
			qsDepthLimitTable[depth] = -5

		} else if depth <= 10 {
			qsDepthLimitTable[depth] = -6

		} else if depth <= 12 {
			qsDepthLimitTable[depth] = -7

		} else if depth <= 14 {
			qsDepthLimitTable[depth] = -8

		} else if depth <= 16 {
			qsDepthLimitTable[depth] = -8

		} else {
			qsDepthLimitTable[depth] = -8
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Search: Initial Call --------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	INFINITY int = 1000000000 // 1 bil
)

var (
	BLANK_MOVE Move = getEncodedMove(99, 99, 99, 99, 99) // dummy move to catch errors
)

// function to initiate a search on the current position and store the best move
func (pos *Position) searchForBestMove(timeLimitMs int) {

	// reset the depth, we start searching at depth 2
	depth := 1

	// reset the position's best move
	pos.bestMove = BLANK_MOVE
	pos.bestMoveSoFar = BLANK_MOVE

	// reset time management nodes
	pos.timeNodesCount = 0
	pos.timeStartingTime = time.Now()
	pos.timeTotalAllowedTime = timeLimitMs

	// reset the search statistics
	pos.logSearch.resetLog()

	// create a new transposition table for the search
	tt := getNewTT()

	// do an iterative deepening search
	for depth <= MAX_DEPTH {

		// increase the depth and search
		depth += 1

		qsDepth := 0
		if depth < 99 {
			qsDepth = qsDepthLimitTable[depth]
		}

		_, terminated := pos.negamax(depth, depth, 0-INFINITY, INFINITY, &tt, qsDepth)

		// store the best move from the search only after each iteration, and continue with the next iteration
		// in case of terminated searches in the middle of a search, we can't use that move, and exit immediately
		// we will definitely hit at least one iteration (say about 400 nodes) at depth 2 with 0 quiescence depth
		// so we will have one best move before the time node limit is checked
		if !terminated {
			pos.bestMove = pos.bestMoveSoFar
			pos.bestMoveSoFar = BLANK_MOVE

			pos.logSearch.depth = depth
			pos.logSearch.qsDepth = qsDepth

		} else {
			pos.logSearch.timeMs = int(time.Since(pos.timeStartingTime).Milliseconds())
			return
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Search: Negamax -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	WHITE_WIN_VALUE int = 10000000  // win value is 10mil (arbitrarily large)
	BLACK_WIN_VALUE int = -10000000 // win value is 10mil  (arbitrarily large)

	PLY_PENALTY int = 5000 // ply penalty is to get the shortest checkmate path: queen is 900 value so 5k is enough

	NODES_BEFORE_CHECK_INTERRUPT = 5000 // after this many nodes, we check whether we need to stop the search
)

// return the score, along with a flag for whether the search was aborted
func (pos *Position) negamax(initialDepth int, currentDepth int, alpha int, beta int, tt *TranspositionTable, qsDepth int) (int, bool) {

	// ---------------------------------------------- Time Management and UCI --------------------------------------------
	// count the nodes searched for time management checks and uci input checks
	// if a certain number of nodes have been reached, pause the search and check whether the time has elapsed
	pos.timeNodesCount += 1

	if pos.timeNodesCount >= NODES_BEFORE_CHECK_INTERRUPT {
		// reset the node count ready for the next check
		pos.timeNodesCount = 0

		// calculate time since start
		// if we are over the allowed time, stop the search
		timeSince := time.Since(pos.timeStartingTime).Milliseconds()
		if timeSince >= int64(pos.timeTotalAllowedTime) {
			return 0, true
		}

		// also, process uci inputs, and check whether we need to stop the search
		/*
			if pos.uciMode {
				pos.checkForUCIInputs()
				if pos.uciStopFlag {
					return 0, true
				}
			}
		*/
	}

	// ---------------------------------------------------- Node Statistics -----------------------------------------------
	// count the nodes searched for statistics
	if currentDepth >= 1 {
		pos.logSearch.nodesAtDepth1Plus += 1
	} else if currentDepth == 0 {
		pos.logSearch.nodesAtDepth0 += 1
	} else {
		pos.logSearch.nodesAtDepth1Min += 1
	}

	// ---------------------------------------------------- TT Lookup -----------------------------------------------
	// the TT lookup is done before moves are generated, eval is done etc. because the TT might already have a hit for the position
	alphaOriginal := alpha // store alpha before the search, if alpha is not increased (no better move found), then we know this is an upper bound

	if currentDepth > 0 {
		start_time_tt_get := time.Now()
		pos.logSearch.nodesTTProbe += 1

		// we only check the TT for non-quiescence nodes
		// because we only save non-quiescence nodes in the TT
		// this can be changed later if needed

		ttEntry, success := tt.getTTEntry(pos.hashOfPos)
		if success { // if there is a node in the TT
			if ttEntry.depth >= uint8(currentDepth) { // and if the depth is at least as deep as the current search

				// if the flag is EXACT and the value is within our current bounds, we can use it
				if ttEntry.flag == TT_FLAG_EXACT {
					if int(ttEntry.value) > alpha && int(ttEntry.value) < beta {
						pos.logSearch.nodesTTHit += 1
						return int(ttEntry.value), false
					}

					// else if the flag was a LOWERBOUND, and it is already higher than out current UPPERBOUND (beta), we can use it
				} else if ttEntry.flag == TT_FLAG_LOWERBOUND {
					if int(ttEntry.value) >= beta {
						pos.logSearch.nodesTTHit += 1
						return beta, false
					}

					// else if the flag was an UPPERBOUND, and it is already lower than our current LOWERBOUND (alpha), we can use it
				} else {
					if int(ttEntry.value) <= alpha {
						pos.logSearch.nodesTTHit += 1
						return alpha, false
					}
				}
			}
		}

		duration_time_tt_get := time.Since(start_time_tt_get).Nanoseconds()
		pos.logOther.allLogTypes[LOG_TT_PROBE].addTime(int(duration_time_tt_get))
	}

	// -------------------------------------------------------- Game Over ------------------------------------------------------

	// _____________________________ Move Generation ______________________________
	// if there is not a TT hit, we need to start with work on the current node
	// first, we generate all legal moves, or at least one move when we are at a leaf node
	// we can then determine if the game is over (no legal moves is checkmate or stalemate)
	generatedPartialMoves := false // flag to catch partial move generation

	if currentDepth <= qsDepth { // leaf nodes: partial move generation
		pos.generateLegalMoves(true)
		generatedPartialMoves = true
		pos.logSearch.nodesGeneratedLegalMovesPart += 1

	} else { // other nodes: full move generation
		pos.generateLegalMoves(false)
		pos.logSearch.nodesGeneratedLegalMovesFull += 1
	}

	// _____________________________ Game State ______________________________
	// once we have generated at least some legal moves, we check whether the game is over
	// if it is over, we return with the game over score
	pos.getGameStateAndStore()
	if pos.gameState != STATE_ONGOING {
		switch pos.gameState {

		case STATE_WIN_WHITE:
			if pos.isWhiteTurn {
				return WHITE_WIN_VALUE - (pos.ply * PLY_PENALTY), false
			} else {
				return 0 - (WHITE_WIN_VALUE - (pos.ply * PLY_PENALTY)), false
			}

		case STATE_WIN_BLACK:
			if pos.isWhiteTurn {
				return BLACK_WIN_VALUE + (pos.ply * PLY_PENALTY), false
			} else {
				return 0 - (BLACK_WIN_VALUE + (pos.ply * PLY_PENALTY)), false
			}

		case STATE_DRAW_STALEMATE, STATE_DRAW_3_FOLD_REPETITION, STATE_DRAW_50_MOVE_RULE:
			return 0, false
		}
	}

	// _____________________________ Check Extensions ______________________________
	// if we are in check in the current node and the game is not over,
	// we extend the search by 1 ply to better search the impact of the check
	// we only do this during qsearch, where the position needs to be quiet before evaluating
	// we also now need to fully generate legal moves if we only generated them partially before
	inCheck := pos.kingChecks > 0
	if inCheck && currentDepth < 0 { // in check and in a qs node
		currentDepth += 1
		pos.logSearch.checkExtensions += 1
		if generatedPartialMoves {
			pos.generateLegalMoves(false)
		}
	}

	// ---------------------------------------------------- Normal Evaluation -----------------------------------------------
	// if the game is not over, and we are at the leaf nodes, we return the search score
	// we return the eval relative to the side to move (not an absolute eval)
	if currentDepth <= qsDepth {
		pos.evalPosAfter()
		if pos.isWhiteTurn {
			return pos.evalMaterial + pos.evalHeatmaps + pos.evalOther, false
		} else {
			return 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther), false
		}
	}

	// ------------------------------------------------------- Order Moves: Normal --------------------------------------------------
	// if we are not at a leaf node, we start ordering moves for the search to try optimise cutoffs
	// we assume there are moves, because if there were no moves, we already would have returned checkmate or stalemate before
	// move ordering is expensive, so we only sort moves certain number of plies away from the leaf nodes
	// note, if we order at qsDepth + 1 then we will never hit unordered nodes (because at leaf nodes we just evaluate)

	// create a slice with the length of the available moves
	start_time_create_move_slice := time.Now()

	copyOfMoves := make([]Move, pos.availableMovesCounter)

	duration_time_create_move_slice := time.Since(start_time_create_move_slice).Nanoseconds()
	pos.logOther.allLogTypes[LOG_CREATE_MOVE_SLICE].addTime(int(duration_time_create_move_slice))

	// now copy the moves into the created slice
	if currentDepth >= (qsDepth + 1) { // nodes with move ordering
		copy(copyOfMoves, pos.getOrderedMoves())
		pos.logSearch.moveOrderedNodes += 1

	} else { // nodes without move ordering
		start_time_copy_into_move_slice := time.Now()

		originalMoves := pos.availableMoves[:pos.availableMovesCounter]
		copy(copyOfMoves, originalMoves)
		pos.logSearch.moveUnorderedNodes += 1

		duration_time_copy_into_move_slice := time.Since(start_time_copy_into_move_slice).Nanoseconds()
		pos.logOther.allLogTypes[LOG_COPY_INTO_MOVE_SLICE].addTime(int(duration_time_copy_into_move_slice))

	}

	// ------------------------------------------------------- Order Moves: Iterative Deepening --------------------------------------------------
	// we also use the previous iterative deepening search's best move first, if we are at the root
	// this code will therefore only run once each time the depth is increased (acceptable because the code takes long)
	if currentDepth == initialDepth { // if we are at the root depth
		if pos.bestMove != BLANK_MOVE { // we need to first get a best move from iterative deepening before we can put it at the front

			start_time_iter_deep_ordering := time.Now()

			// find the index of the best move
			bestIndex := -1
			for index, move := range copyOfMoves {
				//if move == pos.bestMove {
				if move.getFromSq() == pos.bestMove.getFromSq() && move.getToSq() == pos.bestMove.getToSq() && move.getPromotionType() == pos.bestMove.getPromotionType() {
					bestIndex = index
				}
			}

			// remove the best move from the original position
			copyOfMoves = append(copyOfMoves[:bestIndex], copyOfMoves[bestIndex+1:]...)

			// append the best move at the start of the list of moves after ordering
			copyOfMoves = append([]Move{pos.bestMove}, copyOfMoves...)

			duration_time_iter_deep_ordering := time.Since(start_time_iter_deep_ordering).Nanoseconds()
			pos.logOther.allLogTypes[LOG_ITER_DEEP_MOVE_FIRST].addTime(int(duration_time_iter_deep_ordering))
		}
	}
	// ------------------------------------------------------- Main Search: Quiescence Eval --------------------------------------------------
	// if we are at a quiescence search node
	// we use a standPat score as a floor on the evaluation for alpha
	// this is done for the case that there is no capture moves, so we at least return the evaluation
	if currentDepth <= 0 {
		pos.evalPosAfter()

		var standPat int
		if pos.isWhiteTurn {
			standPat = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
		} else {
			standPat = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
		}

		// beta (UPPERBOUND) is not changed in the search, so if it is already above that, return beta
		if standPat >= beta {
			return beta, false
		}

		// else, set alpha to be at least the evaluation score
		// because we assume captures can either improve the position, otherwise we won't make the capture
		if alpha < standPat {
			alpha = standPat
		}
	}

	// ------------------------------------------------------- Main Search --------------------------------------------------
	// start the search and iterate over each move
	for _, move := range copyOfMoves {

		// ___________________________________________ Quiescence Moves ___________________________________
		// at the depth of zero or lower, we only consider captures, en-passant and promotions, and skip other moves
		if currentDepth <= 0 {
			moveType := move.getMoveType()
			promotionType := move.getPromotionType()
			if (moveType != MOVE_TYPE_CAPTURE) && (moveType != MOVE_TYPE_EN_PASSANT) && (promotionType == PROMOTION_NONE) {
				continue
			}
		}

		// ___________________________________________ Make and Undo Move ___________________________________
		// play the move, get the score of the node, and undo the move again
		pos.makeMove(move)
		score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, qsDepth)
		moveValue := 0 - score
		pos.undoMove()

		// if the search was terminated, return with a zero value
		if terminated {
			return 0, true
		}

		// ___________________________________________ Store Best Move at Root ___________________________________
		// if we are at the root and the move is the best move so far, we store the move as the best so far
		if currentDepth == initialDepth {
			if moveValue > alpha {
				pos.bestMoveSoFar = move
			}
		}

		// ___________________________________________ Alpha-Beta Cutoffs ___________________________________

		// beta cutoff
		if moveValue >= beta {

			if currentDepth > 0 {
				// store TT entries for non-quiescence nodes
				// if we have a beta cut, this node failed high
				// so beta is the lowest bound for next searches
				start_time_tt_store := time.Now()

				tt.storeNewTTEntry(pos.hashOfPos, uint8(currentDepth), TT_FLAG_LOWERBOUND, int32(beta))
				pos.logSearch.nodesTTStore += 1

				duration_time_tt_store := time.Since(start_time_tt_store).Nanoseconds()
				pos.logOther.allLogTypes[LOG_TT_STORE].addTime(int(duration_time_tt_store))
			}

			return beta, false
		}

		// improvement of alpha
		if moveValue > alpha {
			alpha = moveValue
		}
	}

	// ---------------------------------------------------- TT Store Entry -----------------------------------------------
	// after iteration over all the moves, we store the node in the TT
	// we store TT entries for non-quiescence nodes because they are fully searched

	if currentDepth > 0 {

		start_time_tt_store := time.Now()

		if alpha > alphaOriginal {
			// if alpha increased in the search, we know the exact value of the node, because:
			// we also did not fail high, because we already would have had a beta cut before this code
			tt.storeNewTTEntry(pos.hashOfPos, uint8(currentDepth), TT_FLAG_EXACT, int32(alpha))
			pos.logSearch.nodesTTStore += 1

		} else {
			// if alpha did not increase in the search, this node failed low
			// it did not fail high, because no beta cut was found
			// this node value is therefore the upper bound for next searches
			tt.storeNewTTEntry(pos.hashOfPos, uint8(currentDepth), TT_FLAG_UPPERBOUND, int32(alpha))
			pos.logSearch.nodesTTStore += 1
		}

		duration_time_tt_store := time.Since(start_time_tt_store).Nanoseconds()
		pos.logOther.allLogTypes[LOG_TT_STORE].addTime(int(duration_time_tt_store))

	}

	// ---------------------------------------------------- Return Final Value -----------------------------------------------
	return alpha, false
}
