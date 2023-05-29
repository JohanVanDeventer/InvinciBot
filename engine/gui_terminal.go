package main

import (
	"fmt"
	"math"
	"strconv"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Terminal GUI -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// simple homemade gui code for playing the game in the console

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
	//if pos.availableMovesCounter <= 0 {
	if pos.totalMovesCounter <= 0 {
		fmt.Println("Error: No available moves.")
		return false
	}

	//allMoves := make([]Move, pos.availableMovesCounter)
	//copy(allMoves, pos.availableMoves[:pos.availableMovesCounter])
	allMoves := make([]Move, pos.totalMovesCounter)
	copy(allMoves, pos.threatMoves[:pos.threatMovesCounter])
	copy(allMoves[pos.threatMovesCounter:], pos.quietMoves[:pos.quietMovesCounter])

	// loop over moves
	// and where the input matches the move, play that move
	var playedMove Move
	foundMove := false
	for _, move := range allMoves {
		//if move.fromSq == fromSq && move.toSq == toSq && move.promotionType == promoteType {
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
		fmt.Printf("======================\n")
	} else {
		fmt.Printf("======================\n")
	}

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
			if pos.isWhiteTurn {
				fmt.Printf("Turn: White. Game Status: %v. Current move: %v. 50-move counter: %v.\n", gameStateToText[pos.gameState], pos.fullMoves, pos.halfMoves)
			} else {
				fmt.Printf("Turn: Black. Game Status: %v. Current move: %v. 50-move counter: %v.\n", gameStateToText[pos.gameState], pos.fullMoves, pos.halfMoves)
			}

		case 6:
			totalEval := pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
			fmt.Printf("Evaluation: %v (material: %v, heatmaps: %v, other: %v).\n", totalEval, pos.evalMaterial, pos.evalHeatmaps, pos.evalOther)

		case 5:
			knps := math.Round((float64(pos.logSearch.getTotalNodes()) / (float64(pos.logSearch.timeMs) / 1000)) / 1000)
			totalNodes := pos.logSearch.getTotalNodes()

			nodesPlus1Percent := 0
			if totalNodes > 0 {
				nodesPlus1Percent = int((float64(pos.logSearch.nodesAtDepth1Plus) / float64(totalNodes)) * 100)
			}

			nodes0Percent := 0
			if totalNodes > 0 {
				nodes0Percent = int((float64(pos.logSearch.nodesAtDepth0) / float64(totalNodes)) * 100)
			}

			nodesMinus1Percent := 0
			if totalNodes > 0 {
				nodesMinus1Percent = int((float64(pos.logSearch.nodesAtDepth1Min) / float64(totalNodes)) * 100)
			}

			fmt.Printf("Search depth: %v. QS Depth: %v. Knps: %v. Total nodes: %v (+1: %v%%  0:%v%%  -1:%v%%). Check extension nodes: %v.\n",
				pos.logSearch.depth, pos.logSearch.qsDepth, knps,
				totalNodes, nodesPlus1Percent, nodes0Percent, nodesMinus1Percent,
				pos.logSearch.checkExtensions)

		case 4:
			totalNodes := pos.logSearch.getTotalNodes()

			ttProbeRate := 0
			if pos.logSearch.nodesTTProbe > 0 {
				ttProbeRate = int((float64(pos.logSearch.nodesTTProbe) / float64(totalNodes)) * 100)
			}

			ttHitRate := 0
			if pos.logSearch.nodesTTProbe > 0 {
				ttHitRate = int((float64(pos.logSearch.nodesTTHit) / float64(pos.logSearch.nodesTTProbe)) * 100)
			}

			ttStoreRate := 0
			if totalNodes > 0 {
				ttStoreRate = int((float64(pos.logSearch.nodesTTStore) / float64(totalNodes)) * 100)
			}

			fmt.Printf("Probed TT at %v%% of nodes. Hit valid TT entry in %v%% of probed nodes. Stored %v%% of nodes in the TT.\n",
				ttProbeRate, ttHitRate, ttStoreRate)

		case 3:
			totalNodes := pos.logSearch.getTotalNodes()
			nodesCreatedMovesCopy := pos.logSearch.moveOrderedNodes + pos.logSearch.moveUnorderedNodes

			genFullLegalMovesRate := 0
			if totalNodes > 0 {
				genFullLegalMovesRate = int((float64(pos.logSearch.nodesGeneratedLegalMovesFull) / float64(totalNodes)) * 100)
			}

			genPartLegalMovesRate := 0
			if totalNodes > 0 {
				genPartLegalMovesRate = int((float64(pos.logSearch.nodesGeneratedLegalMovesPart) / float64(totalNodes)) * 100)
			}

			createMovesCopyRate := 0
			if totalNodes > 0 {
				createMovesCopyRate = int((float64(nodesCreatedMovesCopy) / float64(totalNodes)) * 100)
			}

			orderedNodesRate := 0
			if totalNodes > 0 {
				orderedNodesRate = int((float64(pos.logSearch.moveOrderedNodes) / float64(totalNodes)) * 100)
			}

			fmt.Printf("Fully generated moves at %v%% of nodes. Partially generated moves at %v%% of nodes. Created moves copies at %v%% of nodes. Ordered the moves of %v%% of nodes.\n",
				genFullLegalMovesRate, genPartLegalMovesRate, createMovesCopyRate, orderedNodesRate)

		case 2:
			fmt.Printf("\n")

		case 1:
			avgMoveGen := pos.logOther.allLogTypes[LOG_MOVE_GEN].getAverageNsPerCall()
			avgMakeMove := pos.logOther.allLogTypes[LOG_MAKE_MOVE].getAverageNsPerCall()
			avgUndoMove := pos.logOther.allLogTypes[LOG_UNDO_MOVE].getAverageNsPerCall()
			avgStoreMoves := pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].getAverageNsPerCall()
			avgEval := pos.logOther.allLogTypes[LOG_EVAL].getAverageNsPerCall()
			avgGameState := pos.logOther.allLogTypes[LOG_GAME_STATE].getAverageNsPerCall()

			fmt.Printf("<<Average ns>> Move Gen: %v. Make Move: %v. Undo Move: %v. Store move: %v. Eval: %v. Game State: %v.\n",
				avgMoveGen, avgMakeMove, avgUndoMove, avgStoreMoves, avgEval, avgGameState)

		case 0:
			avgOrderMovesNotAtRoot := pos.logOther.allLogTypes[LOG_ORDER_MOVES_NOT_AT_ROOT].getAverageNsPerCall()
			avgTTGet := pos.logOther.allLogTypes[LOG_TT_PROBE].getAverageNsPerCall()
			avgTTStore := pos.logOther.allLogTypes[LOG_TT_STORE].getAverageNsPerCall()
			avgIterDeepOrder := pos.logOther.allLogTypes[LOG_ITER_DEEP_MOVE_FIRST].getAverageNsPerCall()
			avgCreateMoveSlice := pos.logOther.allLogTypes[LOG_CREATE_MOVE_SLICE].getAverageNsPerCall()
			avgCopyIntoMoveSlice := pos.logOther.allLogTypes[LOG_COPY_INTO_MOVE_SLICE].getAverageNsPerCall()

			fmt.Printf("<<Average ns>> TT Probe: %v. TT Store: %v. Order moves: %v. Create move slice: %v. Copy into move slice: %v. IterDeep Ordering: %v.\n",
				avgTTGet, avgTTStore, avgOrderMovesNotAtRoot, avgCreateMoveSlice, avgCopyIntoMoveSlice, avgIterDeepOrder)
		}
	}
	fmt.Println("======================")
	fmt.Println("    a b c d e f g h")
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
		fmt.Printf("Enter the time the computer is allowed per move in seconds: ")
		fmt.Scanln(&timeInput)
		timeInt, _ := strconv.Atoi(timeInput)
		timePerMoveMs = timeInt * 1000
	}

	// start the loop
	for {

		// print the board to the terminal
		pos.printBoardToTerminal()

		// if the game is over, break the loop
		if pos.gameState != STATE_ONGOING {
			pos.logOther.printLoggedDetails()
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
