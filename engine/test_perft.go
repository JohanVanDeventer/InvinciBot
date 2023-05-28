package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Perft Positions -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// various test positions used to test move generation

type TestPosition struct {
	fen          string
	depthResults []string
}

var testPositions []TestPosition

func initTestPositions() {

	// test position 1: starting position
	testPos1 := TestPosition{}
	testPos1.fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	testPos1.depthResults = append(testPos1.depthResults, "20")
	testPos1.depthResults = append(testPos1.depthResults, "400")
	testPos1.depthResults = append(testPos1.depthResults, "8902")
	testPos1.depthResults = append(testPos1.depthResults, "197281")
	testPos1.depthResults = append(testPos1.depthResults, "4865609")
	//testPos1.depthResults = append(testPos1.depthResults, "119060324")
	//testPos1.depthResults = append(testPos1.depthResults, "3195901860")
	testPositions = append(testPositions, testPos1)

	// test position 2: Kiwipete (complex middlegame)
	testPos2 := TestPosition{}
	testPos2.fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	testPos2.depthResults = append(testPos2.depthResults, "48")
	testPos2.depthResults = append(testPos2.depthResults, "2039")
	testPos2.depthResults = append(testPos2.depthResults, "97862")
	testPos2.depthResults = append(testPos2.depthResults, "4085603")
	//testPos2.depthResults = append(testPos2.depthResults, "193690690")
	//testPos2.depthResults = append(testPos2.depthResults, "8031647685")
	testPositions = append(testPositions, testPos2)

	// test position 3: en-passant discovered check
	testPos3 := TestPosition{}
	testPos3.fen = "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"
	testPos3.depthResults = append(testPos3.depthResults, "14")
	testPos3.depthResults = append(testPos3.depthResults, "191")
	testPos3.depthResults = append(testPos3.depthResults, "2812")
	testPos3.depthResults = append(testPos3.depthResults, "43238")
	testPos3.depthResults = append(testPos3.depthResults, "674624")
	testPos3.depthResults = append(testPos3.depthResults, "11030083")
	//testPos3.depthResults = append(testPos3.depthResults, "178633661")
	//testPos3.depthResults = append(testPos3.depthResults, "3009794393")
	testPositions = append(testPositions, testPos3)

	// test position 4: another complex middlegame, starting with check
	testPos4 := TestPosition{}
	testPos4.fen = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	testPos4.depthResults = append(testPos4.depthResults, "6")
	testPos4.depthResults = append(testPos4.depthResults, "264")
	testPos4.depthResults = append(testPos4.depthResults, "9467")
	testPos4.depthResults = append(testPos4.depthResults, "422333")
	testPos4.depthResults = append(testPos4.depthResults, "15833292")
	//testPos4.depthResults = append(testPos4.depthResults, "706045033")
	testPositions = append(testPositions, testPos4)

	// test position 5: promotions and castling
	testPos5 := TestPosition{}
	testPos5.fen = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
	testPos5.depthResults = append(testPos5.depthResults, "44")
	testPos5.depthResults = append(testPos5.depthResults, "1486")
	testPos5.depthResults = append(testPos5.depthResults, "62379")
	testPos5.depthResults = append(testPos5.depthResults, "2103487")
	//testPos5.depthResults = append(testPos5.depthResults, "89941194")
	testPositions = append(testPositions, testPos5)

	// test position 6: symmetrical italian game
	testPos6 := TestPosition{}
	testPos6.fen = "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	testPos6.depthResults = append(testPos6.depthResults, "46")
	testPos6.depthResults = append(testPos6.depthResults, "2079")
	testPos6.depthResults = append(testPos6.depthResults, "89890")
	testPos6.depthResults = append(testPos6.depthResults, "3894594")
	//testPos6.depthResults = append(testPos6.depthResults, "164075551")
	//testPos6.depthResults = append(testPos6.depthResults, "6923051137")
	testPositions = append(testPositions, testPos6)

	// test position 7 (custom): custom game with lots of pins, checks, en-passants, start with black, and each piece type is still on the board
	testPos7 := TestPosition{}
	testPos7.fen = "8/1np1p1b1/2Bk3B/q6r/4K3/2Q3R1/2NP1P2/1b6 b - - 0 1"
	testPos7.depthResults = append(testPos7.depthResults, "41")
	testPos7.depthResults = append(testPos7.depthResults, "1375")
	testPos7.depthResults = append(testPos7.depthResults, "43366")
	testPos7.depthResults = append(testPos7.depthResults, "1469509")
	//testPos7.depthResults = append(testPos7.depthResults, "46721049")
	testPositions = append(testPositions, testPos7)

}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------- Perft ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// count nodes visited
func (pos *Position) runPerft(initialDepth int, currentDepth int) int {

	// check for depth limit
	if currentDepth == 0 {
		return 1
	}

	// reset the node count
	totalNodeCount := 0

	// generate legal moves
	pos.generateLegalMoves(false)

	// if there are legal moves, iterate over them
	if pos.availableMovesCounter > 0 {
		legalMoves := make([]Move, pos.availableMovesCounter)
		copy(legalMoves, pos.availableMoves[:pos.availableMovesCounter])

		for _, move := range legalMoves {
			pos.makeMove(move)
			currentNodeCount := pos.runPerft(initialDepth, currentDepth-1)
			totalNodeCount += currentNodeCount
			pos.undoMove()

			if initialDepth == currentDepth {
				//fmt.Printf("For move: %v, the node count is: %v\n", move, currentNodeCount)
			}
		}
	}

	// return the nodeCount
	return totalNodeCount
}

// loop over each perft position and print the results to the terminal
func printPerftResults() {

	fmt.Println(" ")
	fmt.Println("------------------------------------ Perft Test Results ---------------------------------------")

	for _, testPosition := range testPositions {

		fmt.Printf(" ------- New Position: %v -------\n", testPosition.fen)

		// create a new position from the test fen string
		newPos := Position{}
		newPos.step1InitFen(testPosition.fen)
		newPos.step2InitRest()

		// start the time
		start_time := time.Now()
		totalNodes := 0

		// test each depth and print the results
		for depth, depthResults := range testPosition.depthResults {
			resultNodes := newPos.runPerft(depth+1, depth+1) // run the perft
			totalNodes += resultNodes
			fmt.Printf("Depth: %v. Correct test nodes: %v. My nodes: %v.\n", depth+1, depthResults, resultNodes)
			if depthResults != strconv.Itoa(resultNodes) {
				fmt.Println("ERROR IN MOVE GENERATION!!!")
			}
		}

		// record the time
		duration_time_sec := time.Since(start_time).Seconds()
		duration_time_ms := time.Since(start_time).Milliseconds()
		knps := math.Round((float64(totalNodes) / (float64(duration_time_ms) / 1000)) / 1000)

		// print the logging details
		// newPos.logOther.printLoggedDetails()

		fmt.Printf("Completed perft in %v seconds. Speed: %v knps.\n\n", duration_time_sec, knps)

	}
}
