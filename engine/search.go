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
- Other techniques included in the comments below
*/

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------- Quiescence Limits ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Depening on the initial depth we call the search from, we have a certain quiescence depth.
Deeper depth searches have longer quiescence depths.
*/

const (
	MAX_DEPTH    int = 100 // we set a max depth for the search (otherwise messes with assigning a best move)
	MIN_QS_DEPTH int = -8  // we set the minimum qs depth as a limit on qs

	// remember, with check extensions (and possibly generating check evasions in qs), the ply can be 2x as big as depth
	// because with check extensions the ply keeps on increasing but the depth remains the same for an extra ply
	KILLER_TABLE_SIZE int = (MAX_DEPTH-MIN_QS_DEPTH)*2 + 1
)

var qsDepthLimitTable [MAX_DEPTH + 1]int

func initQSDepthLimits() {
	for depth := 0; depth <= MAX_DEPTH; depth++ {

		// for finding at least 1 best move, only depths <= 2 have no qs
		// all other depths have qs
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
			qsDepthLimitTable[depth] = MIN_QS_DEPTH

		} else if depth <= 16 {
			qsDepthLimitTable[depth] = MIN_QS_DEPTH

		} else {
			qsDepthLimitTable[depth] = MIN_QS_DEPTH
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

	pos.logTime.allLogTypes[LOG_ONCE_SEARCH_STARTUP].start()

	// reset the search statistics
	pos.logSearch = getNewSearchLogger()

	// set the starting time of the search
	pos.logSearch.start()

	// reset the depth, we start searching at depth 2
	depth := 1

	// reset the position's best move
	pos.bestMove = BLANK_MOVE
	pos.bestMoveSoFar = BLANK_MOVE

	// reset time management nodes
	pos.timeNodesCount = 0
	pos.timeStartingTime = time.Now()
	pos.timeTotalAllowedTime = timeLimitMs

	// reset the killer moves table
	pos.resetKillerMoveTable()

	// create a new transposition table for the search
	tt := getNewTT()

	// create a new history table for the search
	ht := getNewHistoryTable()

	pos.logTime.allLogTypes[LOG_ONCE_SEARCH_STARTUP].stop()

	// do an iterative deepening search
	for depth < MAX_DEPTH {

		// increase the depth
		depth += 1

		// get the qs depth for this depth
		qsDepth := 0
		if depth < 99 {
			qsDepth = qsDepthLimitTable[depth]
		}

		// set the starting ply
		ply := 0

		// do the search
		_, terminated := pos.negamax(depth, depth, 0-INFINITY, INFINITY, tt, ht, qsDepth, ply, false)

		// store the best move from the search only after each iteration, and continue with the next iteration
		// in case of terminated searches in the middle of a search, we can't use that move, and exit immediately
		// we will definitely hit at least one iteration (say about 400 nodes) at depth 2 with 0 quiescence depth
		// so we will have one best move before the time node limit is checked
		if !terminated {
			pos.bestMove = pos.bestMoveSoFar
			pos.bestMoveSoFar = BLANK_MOVE

			pos.logSearch.depth = depth
			pos.logSearch.qsDepth = qsDepth
			pos.logSearch.logIteration()

		} else {
			break
		}
	}

	// finally, log the time taken
	pos.logSearch.stop()
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Search: Negamax -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	WHITE_WIN_VALUE int = 10000000  // win value is 10mil (arbitrarily large)
	BLACK_WIN_VALUE int = -10000000 // win value is 10mil  (arbitrarily large)

	PLY_PENALTY int = 5000 // ply penalty is to get the shortest checkmate path: queen is 900 value so 5k is enough

	NODES_BEFORE_CHECK_INTERRUPT int = 5000 // after this many nodes, we check whether we need to stop the search

	MIN_CHECKMATE = BLACK_WIN_VALUE + (1000 * PLY_PENALTY) // below this score is just checkmate scores
	MAX_CHECKMATE = WHITE_WIN_VALUE - (1000 * PLY_PENALTY) // above this score is just checkmate scores
)

// return the score, along with a flag for whether the search was aborted
func (pos *Position) negamax(initialDepth int, currentDepth int, alpha int, beta int, tt *TranspositionTable, ht *HistoryTable, qsDepth int, ply int, parentWasNull bool) (int, bool) {

	// -------------------------------------------------------- Ply -----------------------------------------------------
	// ply is the depth since the root, independent of any changes to currentDepth (such as check extensions)
	// ply is used to make sure the side to move and actual depth from the root stays correct, for example for killer moves
	ply += 1

	// -------------------------------------------------- Time Management -----------------------------------------------
	// count the nodes searched for time management checks
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
	}

	// ---------------------------------------------------- Node Statistics -----------------------------------------------
	// count the nodes searched for statistics
	nodeType := NODE_TYPE_QS
	if currentDepth > 0 {
		nodeType = NODE_TYPE_NORMAL
	}
	pos.logSearch.depthLogs[nodeType].nodes++

	// ---------------------------------------------------- TT Lookup -----------------------------------------------
	// the TT lookup is done before moves are generated, eval is done etc. because the TT might already have a hit for the position
	// also, store alpha before the search. if alpha is not increased (no better move found), then we know this is an upper bound
	alphaOriginal := alpha

	// set a variable to catch the TT move if present
	hashMove := BLANK_MOVE

	// we only check the TT for non-quiescence nodes
	// because we only save non-quiescence nodes in the TT
	// this can be changed later if needed
	if currentDepth > 0 {

		pos.logTime.allLogTypes[LOG_SEARCH_TT_PROBE].start()
		pos.logSearch.depthLogs[nodeType].ttProbe++

		ttEntry, success := tt.getTTEntry(pos.hashOfPos)
		if success { // if there is a node in the TT
			if ttEntry.depth >= uint8(currentDepth) { // and if the depth is at least as deep as the current search

				// if the flag is EXACT and the value is within our current bounds, we can use it
				if ttEntry.flag == TT_FLAG_EXACT {
					if int(ttEntry.value) > alpha && int(ttEntry.value) < beta {
						pos.logSearch.depthLogs[nodeType].ttHitExact++
						return int(ttEntry.value), false
					}

					// else if the flag was a LOWERBOUND, and it is already higher than out current UPPERBOUND (beta), we can use it
				} else if ttEntry.flag == TT_FLAG_LOWERBOUND {
					if int(ttEntry.value) >= beta {
						pos.logSearch.depthLogs[nodeType].ttHitLower++
						return beta, false
					}

					// else if the flag was an UPPERBOUND, and it is already lower than our current LOWERBOUND (alpha), we can use it
				} else {
					if int(ttEntry.value) <= alpha {
						pos.logSearch.depthLogs[nodeType].ttHitUpper++
						return alpha, false
					}
				}
			}

			// set the hash move if we found a TT Entry but did not get an early cutoff
			hashMove = ttEntry.move
			pos.logSearch.depthLogs[nodeType].ttRetrievedHashMove++
		}

		pos.logTime.allLogTypes[LOG_SEARCH_TT_PROBE].stop()
	}

	// ---------------------------------------------------- Legal Moves and Game Over --------------------------------------------------

	// _____________________________ Move Generation ______________________________
	// if there is not a TT hit, we need to start with work on the current node
	// first, we generate all legal moves
	// we can then determine if the game is over (no legal moves is checkmate or stalemate)
	pos.generateLegalMoves()
	pos.logSearch.depthLogs[nodeType].generatedLegalMovesFull++

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

		case STATE_DRAW_STALEMATE, STATE_DRAW_50_MOVE_RULE:
			return 0, false

		case STATE_DRAW_3_FOLD_REPETITION:
			return 0, false
		}
	}

	// _____________________________ Check Extensions ______________________________
	// if we are in check in the current node and the game is not over,
	// we extend the search by 1 ply to better search the impact of the check
	// we only do this at nodes at least 2 below the initial depth,
	// because increasing at 1 below would put us back at the root and mess up the root-specific code
	// additionally, we only extend in iterative deepening depth > 2 (depth 2 should just be a quick search)
	inCheck := pos.kingChecks > 0

	if inCheck && currentDepth <= (initialDepth-2) && initialDepth > 2 {
		currentDepth += 1
		pos.logSearch.depthLogs[nodeType].checkExtensions++
	}

	// ------------------------------------------------------------- Evaluation --------------------------------------------------------
	// we now evaluate the position at qs nodes

	// ________________________________ QS LEAF NODES _______________________________
	// if the game is not over, and we are at the leaf nodes, we return the eval score
	if currentDepth <= qsDepth {

		pos.evalPosAfter()
		var nodeEval int
		pos.logSearch.depthLogs[nodeType].evalNode++

		if pos.isWhiteTurn {
			nodeEval = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
		} else {
			nodeEval = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
		}

		pos.logSearch.depthLogs[nodeType].qsLeafNodes++
		return nodeEval, false

		// _____________________________ QS OTHER NODES _____________________________
		// if we are at a quiescence node but not a leaf node, we use the evaluation as a floor for alpha
		// this is done for the case that there is no threat/capture moves, so we at least return the evaluation
		// this also allows beta cutoffs before we loop over all the moves
	} else if currentDepth <= 0 {

		pos.evalPosAfter()
		var nodeEval int
		pos.logSearch.depthLogs[nodeType].evalNode++

		if pos.isWhiteTurn {
			nodeEval = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
		} else {
			nodeEval = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
		}

		pos.logSearch.depthLogs[nodeType].qsOtherNodes++

		// beta (UPPERBOUND) is not changed in this node, so if it is already above that, return beta
		if nodeEval >= beta {
			pos.logSearch.depthLogs[nodeType].qsStandPatBetaCuts++
			return beta, false
		}

		// else, set alpha (LOWERBOUND) to be at least the evaluation score
		// because we assume captures can either improve the position, otherwise we won't make the capture
		if alpha < nodeEval {
			pos.logSearch.depthLogs[nodeType].qsStandPatAlphaRaises++
			alpha = nodeEval
		}
	}

	// ------------------------------------------- Static Null Move Pruning / Reverse Futility Pruning ---------------------------------------
	// this is basically a cheaper version of null move pruning below
	// if our static evaluation is so good that even if we give ourselves a big hit materially,
	// and we are still above beta, we assume this node will fail high and we can prune it
	// we don't do this while in check to limit tactical weakness
	// we also don't test this if beta is close to checkmate
	// we do this for all nodes that are not qs, however the margin increases for each depth

	if currentDepth > 0 && !inCheck && initialDepth > 2 && beta < MAX_CHECKMATE && currentDepth != initialDepth {

		// evaluate the node
		pos.evalPosAfter()
		var nodeEval int
		pos.logSearch.depthLogs[nodeType].evalNode++

		if pos.isWhiteTurn {
			nodeEval = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
		} else {
			nodeEval = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
		}

		// set the pruning margin
		pruningMargin := currentDepth * VALUE_PAWN

		// check whether we can prune
		if nodeEval-pruningMargin > beta {
			pos.logSearch.depthLogs[nodeType].staticNullMovePrunes++
			return beta, false
		}
		pos.logSearch.depthLogs[nodeType].staticNullMoveNonPrunes++

	}

	// ----------------------------------------------------------- Null Move Pruning ---------------------------------------------------------

	// ___ Idea ___
	// a null move is a move where we do nothing and just pass the turn to the opponent
	// we do a search on a tree with reduced depth after the null move
	// we also search using a zero width window (because it gives enough information about the null move: we only want to see if there is a beta cut or not)
	// if the score of the search is still high enough after this (score >= beta), we can prune the whole branch

	// ___ Null Move Code  ___
	// we implemented special code that allows making a null move
	// the undo move code will still work with null moves, no need for a special undo null move function

	// ___ Restrictions ___
	// we cannot try a null move when we are in check (because then the position would be illegal after the null move)
	// a null move will also not work in zugzwang positions (because there a null move is actually better than all other moves)
	// zugzwang positions are most common in endgames, so we only try a null move when we are not in the endgame as determined by the game stage
	// we also don't do null moves while in QS
	// we also don't try null moves at the root
	// we also don't allow CONSECUTIVE null moves (two null moves after each other)

	// ___ Restrictions: Eval ___
	// we also only try null moves if the eval is above beta
	// because we then say in this case our eval is generally good enough that we can give our opponent a free move
	// this will also prevent beta mate scores from influencing the result, because eval is much smaller than the mate score

	// ___ Null Move Pruning ___
	if currentDepth >= 4 && currentDepth != initialDepth {

		// evaluate the node
		var nodeEval int
		pos.logSearch.depthLogs[nodeType].evalNode++

		if pos.isWhiteTurn {
			nodeEval = pos.evalMaterial + pos.evalHeatmaps
		} else {
			nodeEval = 0 - (pos.evalMaterial + pos.evalHeatmaps)
		}

		// check whether we can do a null move
		if !inCheck && pos.evalMidVsEndStage >= 6 && nodeEval >= (beta-30) && !parentWasNull {

			nullMoveReduction := 2
			if currentDepth >= 5 {
				nullMoveReduction = 2 + (currentDepth / 4)
			}

			// make the null move, get the search score, and undo the null move
			pos.makeNullMove()
			nullMoveScore, terminated := pos.negamax(initialDepth, currentDepth-nullMoveReduction-1, 0-beta, 0-beta+1, tt, ht, qsDepth, ply, true)
			nullMoveValue := 0 - nullMoveScore
			pos.undoMove()

			// if the search was terminated, return with a zero value
			if terminated {
				return 0, true
			}

			// if we fail-high on the null move, we can prune this whole branch before starting with the normal move ordering and search
			if nullMoveValue >= beta {
				pos.logSearch.depthLogs[nodeType].nullMoveSuccesses++
				return beta, false
			}
			pos.logSearch.depthLogs[nodeType].nullMoveFailures++

			// if we don't return early, we need to again generate all legal moves in the position (making and undoing moves reset the legal moves)
			pos.generateLegalMoves()
			pos.logSearch.depthLogs[nodeType].generatedLegalMovesFull++
		}
	}

	/* (commenting out for now: did not give an improvement)

	// ----------------------------------------------------------- Razoring --------------------------------------------------------
	// if we are close to the qs depth and the evaluation is very bad for us (below alpha by a large margin),
	// we try to immediately drop to qsearch
	// to confirm the position will likely fail low (i.e. no moves will cause an alpha increase)
	// so in qsearch basically only a capture or promotion above that can save the position
	// if the qsearch confirms a fail-low, we trust it and fail low by returning alpha
	// we don't try this in the endgame, where not pruning quiet moves become more important
	// we also don't try this if the parent was a null move (to allow all replies to the null move)
	// we also don't try this if beta is a checkmate score
	if currentDepth <= 2 && currentDepth > 0 {
		if initialDepth > 2 && !inCheck && pos.evalMidVsEndStage >= 8 && !parentWasNull && beta < MAX_CHECKMATE {

			// get the static eval
			pos.evalPosAfter()
			var staticEval int
			if pos.isWhiteTurn {
				staticEval = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
			} else {
				staticEval = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
			}

			// set the razoring margin
			var razorMargin int
			if currentDepth == 1 {
				razorMargin = VALUE_KNIGHT
			} else {
				razorMargin = VALUE_ROOK
			}

			// if the static eval + the razoring margin is below alpha
			// we start a direct negamax search at depth -1 (to not allow check extensions out of qs),
			// but we increase the qs depth by 1 to compensate
			// we don't do a move, so we don't swap alpha and beta
			// we are still at the same level, so we decrease ply by 1, because it will immediately be increased again
			if (staticEval + razorMargin) < alpha {

				// do the qsearch
				score, terminated := pos.negamax(initialDepth, -1, alpha, beta, tt, qsDepth-1, ply-1, false)
				moveValue := 0 - score

				// if the search was terminated, return with a zero value
				if terminated {
					return 0, true
				}

				// now check whether the value improved alpha
				// if not, we regenerate the legal moves (they are overwritten in qsearch)
				if moveValue < alpha {
					pos.logSearch.depthLogs[nodeType].razorSuccesses++
					return alpha, false
				} else {
					pos.logSearch.depthLogs[nodeType].razorFailures++
					pos.logSearch.depthLogs[nodeType].generatedLegalMovesFull++
					pos.generateLegalMoves()
				}
			}
		}
	}

	*/

	// ----------------------------------------------------------- Order Moves --------------------------------------------------------
	// if we are not at a leaf node, we start ordering moves for the search to try optimise cutoffs
	// we assume there are moves, because if there were no moves, we already would have returned checkmate or stalemate before
	// we try moves in the following order:
	// 1. Best Moves (moves from the TT or the best move from the previous iteration)
	// 2. Good Threat Moves (captures (normal or en-passant captures) or promotions that have a score >= 0)
	// 3. Killer Moves (quiet moves that caused a cutoff in a sibling node)
	// 4. Bad Threat Moves (captures (normal or en-passant captures) or promotions that have a score < 0)
	// 5. Other Quiet Moves (rest of the moves not included above)

	// ___________________________________ THREAT MOVES ___________________________________
	// we always create threat moves

	pos.logTime.allLogTypes[LOG_SEARCH_ORDER_THREAT_MOVES].start()

	// get the ordered good and bad threat moves in the position
	var copyOfGoodThreatMoves []Move
	var copyOfBadThreatMoves []Move

	if currentDepth > 0 {
		copyOfGoodThreatMoves, copyOfBadThreatMoves = pos.getOrderedThreatMovesNormalNodes()
	} else {
		copyOfGoodThreatMoves, copyOfBadThreatMoves = pos.getOrderedThreatMovesQsNodes()
	}

	pos.logTime.allLogTypes[LOG_SEARCH_ORDER_THREAT_MOVES].stop()
	pos.logSearch.depthLogs[nodeType].orderThreatMoves++

	// ___________________________________ QUIET MOVES ___________________________________
	// we only create quiet moves at non-quiescence nodes

	var copyOfQuietMoves []Move
	if currentDepth > 0 { // non-quiescence nodes

		pos.logTime.allLogTypes[LOG_SEARCH_COPY_QUIET_MOVES].start()

		// create a slice with the length of the available moves and copy the unordered moves into the created slice
		copyOfQuietMoves = make([]Move, pos.quietMovesCounter)
		copy(copyOfQuietMoves, pos.quietMoves[:pos.quietMovesCounter])

		pos.logTime.allLogTypes[LOG_SEARCH_COPY_QUIET_MOVES].stop()
		pos.logSearch.depthLogs[nodeType].copyQuietMoves++
	}

	// ___________________________________ BEST MOVES: HASH ___________________________________
	// we have certain guesses for the best move, regardless of the threat vs quiet move split
	// we test these best moves before threat and quiet moves
	// one of these is the hash move from the transposition table

	var copyOfBestMoves []Move
	if hashMove != BLANK_MOVE { // if we found a valid candidate hash move from the transposition table

		pos.logTime.allLogTypes[LOG_SEARCH_ORDER_HASH_MOVES].start()
		pos.logSearch.depthLogs[nodeType].ttTestedHashMove++

		// _____ Good Threat Moves _____
		hashIndexGoodThreatMoves := -1
		for index, move := range copyOfGoodThreatMoves {
			if move == hashMove {
				hashIndexGoodThreatMoves = index
			}
		}

		if hashIndexGoodThreatMoves != -1 {
			// remove the hash move from the original position from threat moves
			copyOfGoodThreatMoves = append(copyOfGoodThreatMoves[:hashIndexGoodThreatMoves], copyOfGoodThreatMoves[hashIndexGoodThreatMoves+1:]...)

			// append the hash move at the start of the list of best moves
			copyOfBestMoves = []Move{hashMove}

			pos.logSearch.depthLogs[nodeType].ttUsedAndOrderedHashMove++

			// _____ Bad Threat Moves _____
		} else {
			hashIndexBadThreatMoves := -1
			for index, move := range copyOfBadThreatMoves {
				if move == hashMove {
					hashIndexBadThreatMoves = index
				}
			}

			if hashIndexBadThreatMoves != -1 {
				// remove the hash move from the original position from threat moves
				copyOfBadThreatMoves = append(copyOfBadThreatMoves[:hashIndexBadThreatMoves], copyOfBadThreatMoves[hashIndexBadThreatMoves+1:]...)

				// append the hash move at the start of the list of best moves
				copyOfBestMoves = []Move{hashMove}

				pos.logSearch.depthLogs[nodeType].ttUsedAndOrderedHashMove++

				// _____ Quiet Moves _____
			} else if currentDepth > 0 {

				hashIndexQuietMoves := -1
				for index, move := range copyOfQuietMoves {
					if move == hashMove {
						hashIndexQuietMoves = index
					}
				}

				if hashIndexQuietMoves != -1 {
					// remove the hash move from the original position from quiet moves
					copyOfQuietMoves = append(copyOfQuietMoves[:hashIndexQuietMoves], copyOfQuietMoves[hashIndexQuietMoves+1:]...)

					// append the hash move at the start of the list of best moves
					copyOfBestMoves = []Move{hashMove}

					pos.logSearch.depthLogs[nodeType].ttUsedAndOrderedHashMove++
				}
			}
		}
		pos.logTime.allLogTypes[LOG_SEARCH_ORDER_HASH_MOVES].stop()
	}

	// ___________________________________ BEST MOVES: ITERATIVE DEEPENING  ___________________________________
	// we also use the previous iterative deepening search's best move first, if we are at the root
	// this code will therefore only run once each time the depth is increased (acceptable because the code takes long)
	// we call this last after other move ordering, because we put the previous iteration's best move first regardless of other move ordering

	if currentDepth == initialDepth { // if we are at the root depth
		bestMovePreviousIteration := pos.bestMove

		// we need to first get a best move from iterative deepening before we can put it at the front (assuming we don't have a hash move already as that move)
		if bestMovePreviousIteration != BLANK_MOVE && bestMovePreviousIteration != hashMove {

			pos.logTime.allLogTypes[LOG_SEARCH_ORDER_PREVIOUS_ITERATION_MOVES].start()
			pos.logSearch.depthLogs[nodeType].orderIterativeDeepeningMove++

			// ________________________ GOOD THREAT MOVES ________________________
			// find the index of the best move
			bestIndexGoodThreatMoves := -1
			for index, move := range copyOfGoodThreatMoves {
				if move == bestMovePreviousIteration {
					bestIndexGoodThreatMoves = index
				}
			}

			if bestIndexGoodThreatMoves != -1 { // best move is a threat move

				// remove the best move from the original position (FROM THREAT MOVES)
				copyOfGoodThreatMoves = append(copyOfGoodThreatMoves[:bestIndexGoodThreatMoves], copyOfGoodThreatMoves[bestIndexGoodThreatMoves+1:]...)

				// append the best move at the start of the list of moves after ordering (TO BEST MOVES)
				copyOfBestMoves = append([]Move{bestMovePreviousIteration}, copyOfBestMoves...)

				// ________________________ BAD THREAT MOVES ________________________
			} else {

				// find the index of the best move
				bestIndexBadThreatMoves := -1
				for index, move := range copyOfBadThreatMoves {
					if move == bestMovePreviousIteration {
						bestIndexBadThreatMoves = index
					}
				}

				if bestIndexBadThreatMoves != -1 { // best move is a threat move

					// remove the best move from the original position (FROM THREAT MOVES)
					copyOfBadThreatMoves = append(copyOfBadThreatMoves[:bestIndexBadThreatMoves], copyOfBadThreatMoves[bestIndexBadThreatMoves+1:]...)

					// append the best move at the start of the list of moves after ordering (TO BEST MOVES)
					copyOfBestMoves = append([]Move{bestMovePreviousIteration}, copyOfBestMoves...)

					// ________________________ QUIET MOVES ________________________
				} else { // best move is a quiet move

					// find the index of the best move
					bestIndexQuietMoves := -1
					for index, move := range copyOfQuietMoves {
						if move == bestMovePreviousIteration {
							bestIndexQuietMoves = index
						}
					}

					// remove the best move from the original position (FROM QUIET MOVES)
					copyOfQuietMoves = append(copyOfQuietMoves[:bestIndexQuietMoves], copyOfQuietMoves[bestIndexQuietMoves+1:]...)

					// append the best move at the start of the list of moves after ordering (TO BEST MOVES)
					copyOfBestMoves = append([]Move{bestMovePreviousIteration}, copyOfBestMoves...)
				}

				// ________________________ HASH MOVES ________________________
				// if the best move is not in threat or quiet moves, we assume it is already in best moves as the hash move

				pos.logTime.allLogTypes[LOG_SEARCH_ORDER_PREVIOUS_ITERATION_MOVES].stop()
			}
		}
	}

	// ------------------------------------------------------- Main Search: Setup -------------------------------------------------------
	// we create a bestMove variable to catch the best move to store in the TT
	bestMove := BLANK_MOVE
	pos.logSearch.depthLogs[nodeType].loopedOverMoves++

	// get the side for storing history moves
	side := SIDE_BLACK
	if pos.isWhiteTurn {
		side = SIDE_WHITE
	}

	// ------------------------------------------------------- Main Search: Best Moves --------------------------------------------------
	// start the search and iterate over each move
	bestMovesTried := 0
	for _, move := range copyOfBestMoves {
		bestMovesTried++

		// ___________________________________________ Make and Undo Move ___________________________________
		// play the move, get the score of the node, and undo the move again
		pos.makeMove(move)
		score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
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

		// ___________ KILLER AND HISTORY MOVES ___________
		// for best moves, we don't know if the move was a quiet move or a threat move
		// so first test for that before storing it as a killer move or in the history table
		moveType := move.getMoveType()
		promotionType := move.getPromotionType()
		isQuiet := false
		if (moveType == MOVE_TYPE_QUIET || moveType == MOVE_TYPE_CASTLE) && (promotionType == PROMOTION_NONE) {
			isQuiet = true
		}

		// beta cutoff
		if moveValue >= beta {

			if isQuiet {
				// ___________ KILLER MOVES ___________
				if move != pos.killerMoves[ply][0] { // if we have a unique new killer move
					pos.killerMoves[ply][1] = pos.killerMoves[ply][0] // move the previous new move to the old move slot
					pos.killerMoves[ply][0] = move                    // save the current killer move in the new move slot
				}

				// ___________ HISTORY MOVES ___________
				ht.goodBetaMove(move, currentDepth, side)
			}

			// ___________ NORMAL CODE ___________
			if currentDepth > 0 {
				// store TT entries for non-quiescence nodes
				// if we have a beta cut, this node failed high
				// so beta is the lowest bound for next searches
				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

				tt.storeNewTTEntry(pos.hashOfPos, move, int32(beta), uint8(currentDepth), TT_FLAG_LOWERBOUND)

				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
				pos.logSearch.depthLogs[nodeType].ttStoreLower++
			}

			pos.logSearch.depthLogs[nodeType].bestMovesCutoffs++
			pos.logSearch.depthLogs[nodeType].bestMovesTriedBeforeCuts += bestMovesTried
			return beta, false
		}

		// improvement of alpha
		if moveValue > alpha {

			// ___________ NORMAL CODE ___________
			alpha = moveValue
			bestMove = move

			if isQuiet {
				// ___________ HISTORY MOVES ___________
				ht.goodAlphaMove(move, currentDepth, side)
			}

		} else {
			if isQuiet {
				// ___________ HISTORY MOVES ___________
				ht.badAlphaMove(move, currentDepth, side)
			}
		}
	}

	if bestMovesTried > 0 {
		pos.logSearch.depthLogs[nodeType].bestMovesNoCutoffs++
		pos.logSearch.depthLogs[nodeType].bestMovesTriedWhenNoCuts += bestMovesTried
	}

	// ------------------------------------------------------- Main Search: Good Threat Moves --------------------------------------------------
	// start the search and iterate over each move
	threatGoodMovesTried := 0
	for _, move := range copyOfGoodThreatMoves {
		threatGoodMovesTried++

		// ___________________________________________ Make and Undo Move ___________________________________
		// play the move, get the score of the node, and undo the move again
		pos.makeMove(move)
		score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
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
				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

				tt.storeNewTTEntry(pos.hashOfPos, move, int32(beta), uint8(currentDepth), TT_FLAG_LOWERBOUND)

				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
				pos.logSearch.depthLogs[nodeType].ttStoreLower++
			}

			pos.logSearch.depthLogs[nodeType].threatGoodMovesCutoffs++
			pos.logSearch.depthLogs[nodeType].threatGoodMovesTriedBeforeCuts += threatGoodMovesTried
			return beta, false
		}

		// improvement of alpha
		if moveValue > alpha {
			alpha = moveValue
			bestMove = move
		}
	}

	if threatGoodMovesTried > 0 {
		pos.logSearch.depthLogs[nodeType].threatGoodMovesNoCutoffs++
		pos.logSearch.depthLogs[nodeType].threatGoodMovesTriedWhenNoCuts += threatGoodMovesTried
	}

	// ------------------------------------------------------- Order Moves: Killer Moves --------------------------------------------------
	// we now order quiet moves to search the best quiet moves first
	var copyOfKillerMoves []Move
	if currentDepth > 0 {

		// _____________ Killer Moves ____________
		// killer moves try to improve the move ordering of quiet moves
		// we store the move that caused a beta-cutoff as a killer move to be tried in sibling nodes
		// we identify sibling nodes using the killerMoves[ply][entry] table
		// we store 2 killer moves in the table for each ply (new and old killer move)
		// we replace the old move, and keep the new move for each depth as we find killer moves
		// after having previously stored killer moves,
		// we now check whether there are killer moves that can help the move ordering of quiet moves in other sibling nodes
		// we get the killer moves, and then loop over them to check whether we have current moves that are the same
		// if found, we then move killer moves to the front of the quiet move list
		// note: we sort killer 2 first, so that killer 1 will be at the 1st position

		// _____ Killer 2 _____
		killer2Move := pos.killerMoves[ply][1]

		// we only loop if we previously stored a killer move
		if killer2Move != BLANK_MOVE {

			pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_2].start()

			killer2Index := -1

			for index, move := range copyOfQuietMoves {
				if move == killer2Move {
					killer2Index = index
				}
			}

			if killer2Index != -1 {
				// remove the killer move from the original position
				copyOfQuietMoves = append(copyOfQuietMoves[:killer2Index], copyOfQuietMoves[killer2Index+1:]...)

				// append the killer move at the start of the list
				copyOfKillerMoves = append(copyOfKillerMoves, killer2Move)
			}

			pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_2].stop()
			pos.logSearch.depthLogs[nodeType].orderKiller2++
		}

		// _____ Killer 1 _____
		killer1Move := pos.killerMoves[ply][0]

		// we only loop if we previously stored a killer move
		if killer1Move != BLANK_MOVE {

			pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_1].start()

			killer1Index := -1

			for index, move := range copyOfQuietMoves {
				if move == killer1Move {
					killer1Index = index
				}
			}

			if killer1Index != -1 {
				// remove the killer move from the original position
				copyOfQuietMoves = append(copyOfQuietMoves[:killer1Index], copyOfQuietMoves[killer1Index+1:]...)

				// append the killer move at the start of the list
				copyOfKillerMoves = append(copyOfKillerMoves, killer1Move)
			}

			pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_1].stop()
			pos.logSearch.depthLogs[nodeType].orderKiller1++
		}
	}

	// ------------------------------------------------------- Main Search: Killer Moves --------------------------------------------------
	// start the search and iterate over each move

	// ___________________________________________ Quiescence Moves ___________________________________
	// at the depth of zero or lower, we only consider threat moves (captures, en-passant and promotions)
	// we therefore only iterate over quiet moves at depths > 0 (non-quiescence nodes)
	if currentDepth > 0 {

		quietKillerMovesTried := 0
		for _, move := range copyOfKillerMoves {
			quietKillerMovesTried++

			// ___________________________________________ Make and Undo Move ___________________________________
			// play the move, get the score of the node, and undo the move again
			pos.makeMove(move)
			score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
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

				// ___________ KILLER MOVES ___________
				if move != pos.killerMoves[ply][0] { // if we have a unique new killer move
					pos.killerMoves[ply][1] = pos.killerMoves[ply][0] // move the previous new move to the old move slot
					pos.killerMoves[ply][0] = move                    // save the current killer move in the new move slot
				}

				// ___________ HISTORY MOVES ___________
				ht.goodBetaMove(move, currentDepth, side)

				// ___________ NORMAL CODE ___________
				if currentDepth > 0 {
					// store TT entries for non-quiescence nodes
					// if we have a beta cut, this node failed high
					// so beta is the lowest bound for next searches
					pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

					tt.storeNewTTEntry(pos.hashOfPos, move, int32(beta), uint8(currentDepth), TT_FLAG_LOWERBOUND)

					pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
					pos.logSearch.depthLogs[nodeType].ttStoreLower++
				}

				pos.logSearch.depthLogs[nodeType].quietKillerMovesCutoffs++
				pos.logSearch.depthLogs[nodeType].quietKillerMovesTriedBeforeCuts += quietKillerMovesTried
				return beta, false
			}

			// improvement of alpha
			if moveValue > alpha {
				// ___________ NORMAL CODE ___________
				alpha = moveValue
				bestMove = move

				// ___________ HISTORY MOVES ___________
				ht.goodAlphaMove(move, currentDepth, side)

			} else {

				// ___________ HISTORY MOVES ___________
				ht.badAlphaMove(move, currentDepth, side)
			}
		}

		if quietKillerMovesTried > 0 {
			pos.logSearch.depthLogs[nodeType].quietKillerMovesNoCutoffs++
			pos.logSearch.depthLogs[nodeType].quietKillerMovesTriedWhenNoCuts += quietKillerMovesTried
		}
	}

	// ------------------------------------------------------- Main Search: Bad Threat Moves --------------------------------------------------
	// start the search and iterate over each move
	threatBadMovesTried := 0
	for _, move := range copyOfBadThreatMoves {
		threatBadMovesTried++

		// ___________________________________________ Make and Undo Move ___________________________________
		// play the move, get the score of the node, and undo the move again
		pos.makeMove(move)
		score, terminated := pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
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
				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

				tt.storeNewTTEntry(pos.hashOfPos, move, int32(beta), uint8(currentDepth), TT_FLAG_LOWERBOUND)

				pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
				pos.logSearch.depthLogs[nodeType].ttStoreLower++
			}

			pos.logSearch.depthLogs[nodeType].threatBadMovesCutoffs++
			pos.logSearch.depthLogs[nodeType].threatBadMovesTriedBeforeCuts += threatBadMovesTried
			return beta, false
		}

		// improvement of alpha
		if moveValue > alpha {
			alpha = moveValue
			bestMove = move
		}
	}

	if threatBadMovesTried > 0 {
		pos.logSearch.depthLogs[nodeType].threatBadMovesNoCutoffs++
		pos.logSearch.depthLogs[nodeType].threatBadMovesTriedWhenNoCuts += threatBadMovesTried
	}

	// ------------------------------------------------------- Order Moves: Other Quiet Moves --------------------------------------------------
	// we order other quiet moves based on their history score
	if currentDepth > 0 && currentDepth != initialDepth {
		pos.logTime.allLogTypes[LOG_SEARCH_ORDER_OTHER_QUIET_MOVES].start()

		pos.orderQuietHistoryMoves(copyOfQuietMoves, ht)

		pos.logTime.allLogTypes[LOG_SEARCH_ORDER_OTHER_QUIET_MOVES].stop()
		pos.logSearch.depthLogs[nodeType].orderOtherQuietMoves++
	}

	// ------------------------------------------------------- Main Search: Other Quiet Moves --------------------------------------------------
	// start the search and iterate over each move

	// ___________________________________________ Quiescence Moves ___________________________________
	// at the depth of zero or lower, we only consider threat moves (captures, en-passant and promotions)
	// we therefore only iterate over quiet moves at depths > 0 (non-quiescence nodes)
	if currentDepth > 0 {

		// _____________ Late Move Reductions _____________

		// ___ Idea ___
		// at the moment other quiet moves represent at most about 1%-2% of the cut nodes on average in the middlegame
		// because we have good move ordering before we reach this point
		// we therefore don't want to waste too much time searching these "late" moves in the move ordering scheme
		// we therefore apply a reduction in the search depth to these moves, and search with a zero width window
		// if a move is able to improve alpha, we know the move is better than expected,
		// and we do a re-search with the full window and the full depth
		// if alpha is not improved, we saved time in doing a reduced depth search

		// ___ Restrictions ___
		// we don't reduce any threat moves (captures or promotions)
		// we also don't reduce any moves when we are in check (as we need to resolve the tactical implications)
		// we also don't reduce while we are at the root
		// we also don't reduce if we are too close to qs (we need at least a bit of time to check all replies before qs)
		// for now we don't reduce killer moves (to test later whether reducing killers also would be an improvement)
		canDoLMR := !inCheck && currentDepth < initialDepth && currentDepth > 2

		// ___ Depth Reduction ___
		// we use a formula to calculate the depth reduction based on the remaining depth
		// remember: LMR is cumulative (reduce child nodes despite the fact that the parent node was also reduced)
		// so reductions need to be not too aggressive, because cumulative reductions will stack quickly
		lmr := 0
		if canDoLMR {
			lmr = 1 + currentDepth/5
		}

		// ___ Endgame LMR ___
		// we can't reduce as much in the endgame because other quiet moves go up to a much as 10% of cut types
		// because of factors such as king opposition, pawn pushes to get close to promotions etc.
		// we therefore reduce lmr the closer we are to the endgame
		if canDoLMR {
			if pos.evalMidVsEndStage <= 4 {
				canDoLMR = false
			} else if pos.evalMidVsEndStage <= 8 {
				if lmr > 1 {
					lmr = 1
				}
			} else if pos.evalMidVsEndStage <= 12 {
				if lmr > 2 {
					lmr = 2
				}
			}
		}

		quietOtherMovesTried := 0
		for _, move := range copyOfQuietMoves {
			quietOtherMovesTried++

			// ___________________________________________ Make and Undo Move ___________________________________
			// play the move, get the score of the node, and undo the move again
			var score int
			var terminated bool
			var moveValue int

			if canDoLMR {
				// ___________ LMR: Zero-Width Window ___________
				pos.makeMove(move)
				score, terminated = pos.negamax(initialDepth, currentDepth-lmr-1, 0-(alpha+1), 0-alpha, tt, ht, qsDepth, ply, false)
				moveValue = 0 - score
				pos.undoMove()

				pos.logSearch.depthLogs[nodeType].lmrReducedNodes++

				// ___________ LMR: Fail-High Research ___________
				if moveValue > alpha {
					pos.makeMove(move)
					score, terminated = pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
					moveValue = 0 - score
					pos.undoMove()

					pos.logSearch.depthLogs[nodeType].lmrReducedNodesFailures++
				}

			} else {
				// ___________ No LMR ___________
				pos.makeMove(move)
				score, terminated = pos.negamax(initialDepth, currentDepth-1, 0-beta, 0-alpha, tt, ht, qsDepth, ply, false)
				moveValue = 0 - score
				pos.undoMove()

				pos.logSearch.depthLogs[nodeType].lmrNonReducedNodes++
			}

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

				// ___________ KILLER MOVES ___________
				if move != pos.killerMoves[ply][0] { // if we have a unique new killer move
					pos.killerMoves[ply][1] = pos.killerMoves[ply][0] // move the previous new move to the old move slot
					pos.killerMoves[ply][0] = move                    // save the current killer move in the new move slot
				}

				// ___________ HISTORY MOVES ___________
				ht.goodBetaMove(move, currentDepth, side)

				// ___________ NORMAL CODE ___________
				if currentDepth > 0 {
					// store TT entries for non-quiescence nodes
					// if we have a beta cut, this node failed high
					// so beta is the lowest bound for next searches
					pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

					tt.storeNewTTEntry(pos.hashOfPos, move, int32(beta), uint8(currentDepth), TT_FLAG_LOWERBOUND)

					pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
					pos.logSearch.depthLogs[nodeType].ttStoreLower++
				}

				pos.logSearch.depthLogs[nodeType].quietOtherMovesCutoffs++
				pos.logSearch.depthLogs[nodeType].quietOtherMovesTriedBeforeCuts += quietOtherMovesTried
				return beta, false
			}

			// improvement of alpha
			if moveValue > alpha {
				// ___________ NORMAL CODE ___________
				alpha = moveValue
				bestMove = move

				// ___________ HISTORY MOVES ___________
				ht.goodAlphaMove(move, currentDepth, side)

			} else {
				// ___________ HISTORY MOVES ___________
				ht.badAlphaMove(move, currentDepth, side)
			}
		}

		if quietOtherMovesTried > 0 {
			pos.logSearch.depthLogs[nodeType].quietOtherMovesNoCutoffs++
			pos.logSearch.depthLogs[nodeType].quietOtherMovesTriedWhenNoCuts += quietOtherMovesTried
		}
	}

	// ---------------------------------------------------- TT Store Entry -----------------------------------------------
	// after iterating over all the moves, we store the node in the TT
	// we only store TT entries for non-quiescence nodes because they are fully searched

	if currentDepth > 0 {

		pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].start()

		if alpha > alphaOriginal {
			// if alpha increased in the search, we know the exact value of the node, because:
			// we did not fail high, because we already would have had a beta cut before this code
			tt.storeNewTTEntry(pos.hashOfPos, bestMove, int32(alpha), uint8(currentDepth), TT_FLAG_EXACT)
			pos.logSearch.depthLogs[nodeType].ttStoreExact++

		} else {
			// if alpha did not increase in the search, this node failed low
			// it did not fail high, because no beta cut was found
			// this node value is therefore the upper bound for next searches
			tt.storeNewTTEntry(pos.hashOfPos, BLANK_MOVE, int32(alpha), uint8(currentDepth), TT_FLAG_UPPERBOUND)
			pos.logSearch.depthLogs[nodeType].ttStoreUpper++
		}

		pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].stop()
	}

	// ---------------------------------------------------- Return Final Value -----------------------------------------------
	return alpha, false
}
