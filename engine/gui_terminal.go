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
	if pos.availableMovesCounter <= 0 {
		fmt.Println("Error: No available moves.")
		return false
	}

	allMoves := make([]Move, pos.availableMovesCounter)
	copy(allMoves, pos.availableMoves[:pos.availableMovesCounter])

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

	pos.generateLegalMoves()
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
				fmt.Printf("Turn: White. Game Status: %v.\n", gameStateToText[pos.gameState])
			} else {
				fmt.Printf("Turn: Black. Game Status: %v.\n", gameStateToText[pos.gameState])
			}

		case 6:
			knps := math.Round((float64(pos.logSearch.getTotalNodes()) / (float64(pos.logSearch.timeMs) / 1000)) / 1000)
			fmt.Printf("Search depth: %v. QS Depth: %v. Nodes: %v. Knps: %v.\n", pos.logSearch.depth, pos.logSearch.qsDepth, pos.logSearch.getTotalNodes(), knps)

		case 5:
			fmt.Printf("Nodes+1: %v. Nodes=0: %v. Nodes-1: %v. TT Uses: %v. TT Stores: %v. Ordered move nodes: %v. Unordered move nodes: %v.\n", pos.logSearch.nodesAtDepth1Plus,
				pos.logSearch.nodesAtDepth0, pos.logSearch.nodesAtDepth1Min, pos.logSearch.nodesTTHit, pos.logSearch.nodesTTStore, pos.logSearch.moveOrderedNodes,
				pos.logSearch.moveUnorderedNodes)

		case 4:
			fmt.Printf("Current move: %v. 50-move counter: %v.\n", pos.fullMoves, pos.halfMoves)

		case 3:
			fmt.Printf("Eval material: %v. Eval heatmaps: %v. Eval other: %v.\n", pos.evalMaterial, pos.evalHeatmaps, pos.evalOther)

		case 2:

			// get the average ns per call
			avgMoveGen := 0
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_KING].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_QUEEN].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_ROOK].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_KNIGHT].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_BISHOP].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_PAWN].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_KING_ATTACKS].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_PINS].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_EN_PASSANT].getAverageNsPerCall()
			avgMoveGen += pos.logOther.allLogTypes[LOG_MOVES_CASTLING].getAverageNsPerCall()

			avgMakeMove := pos.logOther.allLogTypes[LOG_MAKE_MOVE].getAverageNsPerCall()
			avgUndoMove := pos.logOther.allLogTypes[LOG_UNDO_MOVE].getAverageNsPerCall()
			avgEval := pos.logOther.allLogTypes[LOG_EVAL].getAverageNsPerCall()
			avgGameState := pos.logOther.allLogTypes[LOG_GAME_STATE].getAverageNsPerCall()

			fmt.Printf("<<Average ns>> Move Gen: %v. Make Move: %v. Undo Move: %v. Eval: %v. Game State: %v.\n", avgMoveGen, avgMakeMove, avgUndoMove, avgEval, avgGameState)

		case 1:
			avgOrderMoves := pos.logOther.allLogTypes[LOG_ORDER_MOVES].getAverageNsPerCall()
			avgTTGet := pos.logOther.allLogTypes[LOG_TT_GET].getAverageNsPerCall()
			avgTTStore := pos.logOther.allLogTypes[LOG_TT_STORE].getAverageNsPerCall()
			avgIterDeepOrder := pos.logOther.allLogTypes[LOG_ITER_DEEP_MOVE_FIRST].getAverageNsPerCall()
			avgStoreMoves := pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].getAverageNsPerCall()

			fmt.Printf("<<Average ns>> Order moves: %v. TT Get: %v. TT Store: %v. IterDeep Ordering: %v. Store move: %v.\n", avgOrderMoves, avgTTGet, avgTTStore, avgIterDeepOrder, avgStoreMoves)

		case 0:
			fmt.Printf("\n")
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

					pos.generateLegalMoves()
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

					pos.generateLegalMoves()
					success := pos.playInputMove(userInput)

					if success {
						break
					}
				}
			}

		}
	}
}
