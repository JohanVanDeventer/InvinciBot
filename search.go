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

*/

// helper functions
func getMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func getMin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------- Quiescence Limits ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Depening on the initial depth we call the search from, we have a certain quiescence depth.
Therefore deeper depth searches have longer quiescence depths.
*/

var (
	qsDepthLimitTable [100]int
)

func initQSDepthLimits() {
	for depth := 0; depth < 100; depth++ {

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
	BLANK_MOVE Move = Move{-1, -1, -1, -1, -1, -1, -1} // dummy move to catch errors
)

// function to initiate a search on the current position and store the best move
func (pos *Position) searchForBestMove(timeLimitMs int) {

	/*
		// ----------- LOG FILE ------------
		// set up to add logs
		// Open the file in append mode. If the file doesn't exist, it will be created.
		file, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// ----------- LOG FILE ------------
	*/

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
	for {

		// increase the depth and search
		depth += 1
		qsDepth := qsDepthLimitTable[depth]
		score, terminated := pos.negamax(depth, depth, 0-INFINITY, INFINITY, &tt, qsDepth)
		if score >= 0 {
		}

		// store the best move from the search only after each iteration, and continue with the next iteration
		// in case of terminated searches in the middle of a search, we can't use that move, and exit immediately
		// we will definitely hit at least one iteration (400 nodes) at depth 2
		// so we will have one best move before the time node limit is checked
		if !terminated {
			pos.bestMove = pos.bestMoveSoFar
			pos.bestMoveSoFar = BLANK_MOVE

			pos.logSearch.depth = depth

		} else {
			pos.logSearch.timeMs = int(time.Since(pos.timeStartingTime).Milliseconds())

			/*
				// ----------- LOG FILE ------------
				// Write a new game to the file
				var searchResult string = "Total time in ms for search: " + strconv.Itoa(timeLimitMs) + ". Depth: " + strconv.Itoa(pos.logSearch.depth) + ". Time Ms: " + strconv.Itoa(pos.logSearch.timeMs) + ". Nodes: " + strconv.Itoa(pos.logSearch.getTotalNodes()) + ". TT Hits: " + strconv.Itoa(pos.logSearch.nodesTTHit) + ". TT Stores: " + strconv.Itoa(pos.logSearch.nodesTTStore) + "."
				_, err = fmt.Fprintln(file, searchResult)
				if err != nil {
					log.Fatal(err)
				}
				// ----------- LOG FILE ------------
			*/

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

	PLY_PENALTY int = 5000 // ply penalty is to get the shortest checkmate path: queen is 900 so 5k is enough

	NODES_BEFORE_CHECK_INTERRUPT = 5000 // after this many nodes, we check whether we need to stop the search
)

// return the evaluation, along with a flag for whether the search aborted
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
		// we only check the TT for non-quiescence nodes
		// because we only save non-quiescence nodes
		// this can be changed later if needed

		ttEntry, success := tt.getTTEntry(pos.hashOfPos)
		if success { // if there is a node in the TT
			if ttEntry.depth >= uint8(currentDepth) { // and if the depth is at least as deep as the current search

				if ttEntry.flag == TT_FLAG_EXACT {
					if int(ttEntry.value) > alpha && int(ttEntry.value) < beta {
						pos.logSearch.nodesTTHit += 1
						return int(ttEntry.value), false
					}

				} else if ttEntry.flag == TT_FLAG_LOWERBOUND {
					if int(ttEntry.value) >= beta {
						pos.logSearch.nodesTTHit += 1
						return beta, false
					}

				} else {
					if int(ttEntry.value) <= alpha {
						pos.logSearch.nodesTTHit += 1
						return alpha, false
					}
				}
			}
		}
	}

	// -------------------------------------------------------- Game Over ------------------------------------------------------
	// if there is not a TT hit, we need to start with work on the current node.
	// first, we generate all legal moves.
	// we can then determine if the game is over (no legal moves is checkmate or stalemate etc.)
	pos.generateLegalMoves()
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

	// ---------------------------------------------------- Normal Evaluation -----------------------------------------------
	// if the game is not over, and we are at the leaf nodes, we return the search score
	if currentDepth <= qsDepth {
		pos.evalPosAfter()
		if pos.isWhiteTurn {
			return pos.evalMaterial + pos.evalHeatmaps + pos.evalOther, false
		} else {
			return 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther), false
		}
	}

	// ------------------------------------------------------- Order Moves --------------------------------------------------
	// we assume there are moves, because if there were no moves, we already would have returned checkmate or stalemate before
	// also, we order the moves before searching to optimise cutoffs
	// move ordering is expensive, so we only sort moves certain number of plies away from the leaf nodes
	copyOfMoves := make([]Move, pos.availableMovesCounter)

	if currentDepth >= qsDepth+1 {
		copy(copyOfMoves, pos.getOrderedMoves())
	} else {
		copy(copyOfMoves, pos.availableMoves[:pos.availableMovesCounter])
	}

	// ***** <<< SPECIAL CODE TO USE ITERATIVE DEEPENING BEST MOVE >>> ***** START
	// we also use the previous iterative deepening search's best move first, if we are at the root
	// this code will therefore only run once each time the depth is increased (acceptable because the code takes long)
	if currentDepth == initialDepth { // if we are at the root depth
		if pos.bestMove != BLANK_MOVE { // we need to first get a best move from iterative deepening before we can put it at the front

			// find the index of the best move
			bestIndex := -1
			for index, move := range copyOfMoves {
				if move == pos.bestMove {
					bestIndex = index
				}
			}

			// remove the best move from the original position
			copyOfMoves = append(copyOfMoves[:bestIndex], copyOfMoves[bestIndex+1:]...)

			// append the best move at the start of the list of moves after ordering
			copyOfMoves = append([]Move{pos.bestMove}, copyOfMoves...)
		}
	}
	// ***** <<< SPECIAL CODE TO USE ITERATIVE DEEPENING BEST MOVE >>> ***** END

	// ------------------------------------------------------- Main Search --------------------------------------------------
	// iterate over each move

	// <<< QUIESCENCE SEARCH >>> Start
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

		if standPat >= beta { // beta is not changed in the search, so if it is already above that, return beta
			return beta, false
		}

		if alpha < standPat { // else, set alpha to be at least the evaluation score
			alpha = standPat
		}
	}
	// <<< QUIESCENCE SEARCH >>> End

	for _, move := range copyOfMoves {

		// <<< QUIESCENCE SEARCH >>> Start
		// at the depth of zero or lower, we only consider captures and promotions, and skip other moves
		if currentDepth <= 0 {
			if (move.moveType != MOVE_TYPE_CAPTURE) && (move.moveType != MOVE_TYPE_EN_PASSANT) && (move.promotionType == PROMOTION_NONE) {
				continue
			}
		}
		// <<< QUIESCENCE SEARCH >>> End

		pos.makeMove(move)
		score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, qsDepth)
		moveValue := 0 - score
		pos.undoMove()

		// if the search is terminated, return with a zero value
		if terminated {
			return 0, true
		}

		// <<< SPECIAL CODE TO STORE BEST MOVE >>> START
		if currentDepth == initialDepth { // if we are at the root
			//fmt.Printf("Move: %v. Move value: %v.\n", move, moveValue)
			if moveValue > alpha { // if the move is the best move so far
				pos.bestMoveSoFar = move
			}
		}
		// <<< SPECIAL CODE TO STORE BEST MOVE >>> END

		// <<< ALPHA-BETA >>> START
		if moveValue >= beta {

			if currentDepth > 0 {
				// store TT entries for non-quiescence nodes
				// if we have a beta cut, this node failed high
				// so beta is the lowest bound for next searches
				tt.storeNewTTEntry(pos.hashOfPos, uint8(currentDepth), TT_FLAG_LOWERBOUND, int32(beta))
				pos.logSearch.nodesTTStore += 1
			}

			return beta, false
		}

		if moveValue > alpha {
			alpha = moveValue
		}
		// <<< ALPHA-BETA >>> END

	}

	// ---------------------------------------------------- TT Store Entry -----------------------------------------------
	/*
		// save the current position and search results in the TT

		var ttValue int32 = int32(value)

		var ttFlag uint8 = 99
		if value <= alphaOriginal {
			ttFlag = TT_FLAG_UPPER
		} else if value >= beta {
			ttFlag = TT_FLAG_LOWER
		} else {
			ttFlag = TT_FLAG_EXACT
		}

		var ttDepth uint8 = uint8(currentDepth)

		tt.storeNewTTEntry(pos.hashOfPos, ttDepth, ttFlag, ttValue)
	*/
	if currentDepth > 0 { // store TT entries for non-quiescence nodes because they are fully searched
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

	}

	// ---------------------------------------------------- Return Final Value -----------------------------------------------

	return alpha, false
}
