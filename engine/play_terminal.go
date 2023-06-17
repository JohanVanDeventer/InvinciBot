package main

import (
	"fmt"
	"strconv"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Terminal GUI -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// simple code for playing the game in the terminal

// translate the user move input to a move recognized by the engine
// example: c3d5 looks for a move from sq 18 to sq 35
// also, e7e8q has a queen promotion
// returns false if the input is not recognized
func (pos *Position) playInputMove(input string) bool {

	// split the input to the separate string parts
	var fromStr string
	var toStr string
	var promoteStr string

	if len(input) == 4 { // normal moves, including castle moves
		fromStr = input[0:2]
		toStr = input[2:]
	} else if len(input) == 5 { // promotion moves
		fromStr = input[0:2]
		toStr = input[2:4]
		promoteStr = input[4:]
	} else { // unrecognized move
		fmt.Println("Error: Unrecognized move.")
		return false
	}

	// convert the string parts to engine ints
	fromSq := getSqFromString(fromStr)
	toSq := getSqFromString(toStr)

	var promoteType int
	if len(promoteStr) == 1 {
		switch promoteStr {
		case "q":
			promoteType = PROMOTION_QUEEN
		case "r":
			promoteType = PROMOTION_ROOK
		case "n":
			promoteType = PROMOTION_KNIGHT
		case "b":
			promoteType = PROMOTION_BISHOP
		}
	}

	// get the available moves
	if pos.totalMovesCounter <= 0 {
		fmt.Println("Error: No available moves.")
		return false
	}

	allMoves := make([]Move, pos.totalMovesCounter)
	copy(allMoves, pos.threatMoves[:pos.threatMovesCounter])
	copy(allMoves[pos.threatMovesCounter:], pos.quietMoves[:pos.quietMovesCounter])

	// loop over moves
	// and where the input matches the move, play that move
	var playedMove Move
	foundMove := false
	for _, move := range allMoves {
		if move.getFromSq() == fromSq && move.getToSq() == toSq && move.getPromotionType() == promoteType {
			playedMove = move
			foundMove = true
		}
	}

	if foundMove {
		pos.makeMove(playedMove)
		return true
	}

	fmt.Println("Error: Did not play a move.")
	return false
}

// lookup table for each chess character
var unicodeChessIcons [2][6]string = [2][6]string{
	{"♚", "♛", "♜", "♞", "♝", "♟"},
	{"♔", "♕", "♖", "♘", "♗", "♙"},
}

// prints the board to the terminal
func (pos *Position) printBoardToTerminal() {

	pos.generateLegalMoves(false)
	pos.getGameStateAndStore()
	pos.evalPosAfter()

	if pos.isWhiteTurn {
		fmt.Printf("====================== ")
	} else {
		fmt.Printf("====================== ")
	}

	if pos.isWhiteTurn {
		fmt.Printf("< INFO >     Turn: White. Game Status: %v. Current move: %v. 50-move counter: %v. ",
			gameStateToText[pos.gameState], pos.fullMoves, pos.halfMoves)
	} else {
		fmt.Printf("< INFO >     Turn: Black. Game Status: %v. Current move: %v. 50-move counter: %v. ",
			gameStateToText[pos.gameState], pos.fullMoves, pos.halfMoves)
	}

	totalEval := pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
	fmt.Printf("Evaluation: %v (material: %v, heatmaps: %v, other: %v).\n", totalEval, pos.evalMaterial, pos.evalHeatmaps, pos.evalOther)

	for rowCounter := 7; rowCounter >= 0; rowCounter-- {

		fmt.Printf("%v |", rowCounter+1) // print the row number for each row

		for colCounter := 0; colCounter <= 7; colCounter++ {
			if pos.pieces[SIDE_WHITE][PIECE_KING].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_KING])
			} else if pos.pieces[SIDE_WHITE][PIECE_QUEEN].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_QUEEN])
			} else if pos.pieces[SIDE_WHITE][PIECE_ROOK].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_ROOK])
			} else if pos.pieces[SIDE_WHITE][PIECE_KNIGHT].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_KNIGHT])
			} else if pos.pieces[SIDE_WHITE][PIECE_BISHOP].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_BISHOP])
			} else if pos.pieces[SIDE_WHITE][PIECE_PAWN].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_WHITE][PIECE_PAWN])
			} else if pos.pieces[SIDE_BLACK][PIECE_KING].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_KING])
			} else if pos.pieces[SIDE_BLACK][PIECE_QUEEN].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_QUEEN])
			} else if pos.pieces[SIDE_BLACK][PIECE_ROOK].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_ROOK])
			} else if pos.pieces[SIDE_BLACK][PIECE_KNIGHT].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_KNIGHT])
			} else if pos.pieces[SIDE_BLACK][PIECE_BISHOP].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_BISHOP])
			} else if pos.pieces[SIDE_BLACK][PIECE_PAWN].isBitSet(sqFromRowAndCol(rowCounter, colCounter)) {
				fmt.Printf(" %v", unicodeChessIcons[SIDE_BLACK][PIECE_PAWN])
			} else {
				fmt.Printf(" .")
			}
		}
		fmt.Printf("  | ")

		switch rowCounter {

		case 7:
			fmt.Printf("< ALL >      %v%v\n", pos.logSearch.getOverallSummary(), pos.logSearch.getBranchingFactorSummary())

		case 6:
			fmt.Printf("< ALL >      %v%v%v\n", pos.logSearch.getMoveOrderingSummary(), pos.logSearch.getMoveGenerationSummary(), pos.logSearch.getCheckExtensionsSummary())

		case 5:
			fmt.Printf("< NON-QS >   %v\n", pos.logSearch.getTTNormalSummary())

		case 4:
			fmt.Printf("< QS >       %v\n", pos.logSearch.getTTQsSummary())

		case 3:
			fmt.Printf("< NON-QS >   %v%v\n", pos.logSearch.getNullMoveSummary(), pos.logSearch.getLMRSummary())

		case 2:
			fmt.Printf("< QS >       %v\n", pos.logSearch.getEvalSummary())

		case 1:
			fmt.Printf("< NON-QS >   %v\n", pos.logSearch.getMoveLoopsNormalSummary())

		case 0:
			fmt.Printf("< QS >       %v\n", pos.logSearch.getMoveLoopsQsSummary())

		}
	}
	fmt.Printf("====================== ")

	avgMoveGen := pos.logTime.allLogTypes[LOG_MOVE_GEN_TOTAL].getAverageNsPerCall()
	avgMakeMove := pos.logTime.allLogTypes[LOG_MAKE_MOVE].getAverageNsPerCall()
	avgUndoMove := pos.logTime.allLogTypes[LOG_UNDO_MOVE].getAverageNsPerCall()
	avgMakeNullMove := pos.logTime.allLogTypes[LOG_MAKE_NULLMOVE].getAverageNsPerCall()
	avgEval := pos.logTime.allLogTypes[LOG_EVAL].getAverageNsPerCall()
	avgGameState := pos.logTime.allLogTypes[LOG_GAME_STATE].getAverageNsPerCall()
	avgTTProbe := pos.logTime.allLogTypes[LOG_SEARCH_TT_PROBE].getAverageNsPerCall()
	avgTTStore := pos.logTime.allLogTypes[LOG_SEARCH_TT_STORE].getAverageNsPerCall()

	fmt.Printf("< AVG TIME > Move Gen: %v. Make Move: %v. Undo Move: %v. Make Null Move: %v. Eval: %v. Game State: %v. TT Probe: %v. TT Store: %v.\n",
		avgMoveGen, avgMakeMove, avgUndoMove, avgMakeNullMove, avgEval, avgGameState, avgTTProbe, avgTTStore)

	fmt.Printf("    a b c d e f g h    ")

	avgCopyThreatMoves := pos.logTime.allLogTypes[LOG_SEARCH_COPY_THREAT_MOVES].getAverageNsPerCall()
	avgCopyQuietMoves := pos.logTime.allLogTypes[LOG_SEARCH_COPY_QUIET_MOVES].getAverageNsPerCall()
	avgOrderThreat := pos.logTime.allLogTypes[LOG_SEARCH_ORDER_THREAT_MOVES].getAverageNsPerCall()
	avgOrderKiller1 := pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_1].getAverageNsPerCall()
	avgOrderKiller2 := pos.logTime.allLogTypes[LOG_SEARCH_ORDER_KILLER_2].getAverageNsPerCall()
	avgOrderHash := pos.logTime.allLogTypes[LOG_SEARCH_ORDER_HASH_MOVES].getAverageNsPerCall()
	avgOrderPrevIter := pos.logTime.allLogTypes[LOG_SEARCH_ORDER_PREVIOUS_ITERATION_MOVES].getAverageNsPerCall()

	fmt.Printf("< AVG TIME > Copy Moves (threat: %v, quiet: %v). Order Moves (threat: %v, killer1: %v, killer2: %v, hash: %v, prev iter: %v). ",
		avgCopyThreatMoves, avgCopyQuietMoves, avgOrderThreat, avgOrderKiller1, avgOrderKiller2, avgOrderHash, avgOrderPrevIter)

	avgStartupFen := pos.logTime.allLogTypes[LOG_ONCE_LOAD_FEN].getAverageNsPerCall()
	avgStartupHash := pos.logTime.allLogTypes[LOG_ONCE_HASH].getAverageNsPerCall()
	avgStartupEval := pos.logTime.allLogTypes[LOG_ONCE_EVAL].getAverageNsPerCall()
	avgStartupSearch := pos.logTime.allLogTypes[LOG_ONCE_SEARCH_STARTUP].getAverageNsPerCall()
	avgStartupTotal := avgStartupFen + avgStartupHash + avgStartupEval + avgStartupSearch
	avgStartupTotalMs := avgStartupTotal / 1000000

	fmt.Printf("Startup (%v ms): %v (fen: %v, hash: %v, eval: %v, search: %v).\n\n",
		avgStartupTotalMs, avgStartupTotal, avgStartupFen, avgStartupHash, avgStartupEval, avgStartupSearch)
}

// computer looks for and plays the best move in the position
func (pos *Position) searchAndPlayBestMove(timePerMoveMs int) {
	pos.searchForBestMove(timePerMoveMs)
	pos.makeMove(pos.bestMove)
}

var gameStateToText [6]string

func initGameStateToText() {
	gameStateToText[STATE_ONGOING] = "Game is ongoing"
	gameStateToText[STATE_WIN_WHITE] = "White Wins!"
	gameStateToText[STATE_WIN_BLACK] = "Black Wins!"
	gameStateToText[STATE_DRAW_STALEMATE] = "Draw by stalemate"
	gameStateToText[STATE_DRAW_3_FOLD_REPETITION] = "Draw by 3-fold repetition"
	gameStateToText[STATE_DRAW_50_MOVE_RULE] = "Draw by 50-move rule"
}

// start the terminal GUI loop until the game is over
func (pos *Position) startGameLoopTerminalGUI() {

	// set up the user preferences
	var computerPlaysWhite bool = false
	var computerPlaysBlack bool = false
	var timePerMoveMs int = 0

	var whiteInput string
	fmt.Printf("Press 'y' for the computer to play white, else press any other key: ")
	fmt.Scanln(&whiteInput)
	if whiteInput == "y" {
		computerPlaysWhite = true
	}

	var blackInput string
	fmt.Printf("Press 'y' for the computer to play black, else press any other key: ")
	fmt.Scanln(&blackInput)
	if blackInput == "y" {
		computerPlaysBlack = true
	}

	if computerPlaysWhite || computerPlaysBlack {
		var timeInput string
		fmt.Printf("Enter the time the computer is allowed per move in milliseconds: ")
		fmt.Scanln(&timeInput)
		timeInt, _ := strconv.Atoi(timeInput)
		timePerMoveMs = timeInt
	}

	// start the loop
	for {

		// print the board to the terminal
		pos.printBoardToTerminal()

		// if the game is over, break the loop
		if pos.gameState != STATE_ONGOING {
			pos.logTime.printLoggedDetails()
			pos.printBoardToTerminal()
			break
		}

		// if not game over, play the game
		if pos.isWhiteTurn {

			if computerPlaysWhite { // computer playing
				pos.searchAndPlayBestMove(timePerMoveMs)

			} else { // human playing

				for {
					var userInput string
					fmt.Printf("Enter the move: ")
					fmt.Scanln(&userInput)

					pos.generateLegalMoves(false)
					success := pos.playInputMove(userInput)

					if success {
						break
					}
				}
			}

		} else {

			if computerPlaysBlack { // computer playing
				pos.searchAndPlayBestMove(timePerMoveMs)

			} else { // human playing

				for {
					var userInput string
					fmt.Printf("Enter the move: ")
					fmt.Scanln(&userInput)

					pos.generateLegalMoves(false)
					success := pos.playInputMove(userInput)

					if success {
						break
					}
				}
			}
		}
	}
}
