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
	testPos1.fen = startingFen
	testPos1.depthResults = append(testPos1.depthResults, "20")
	testPos1.depthResults = append(testPos1.depthResults, "400")
	testPos1.depthResults = append(testPos1.depthResults, "8902")
	testPos1.depthResults = append(testPos1.depthResults, "197281")
	testPos1.depthResults = append(testPos1.depthResults, "4865609")
	testPos1.depthResults = append(testPos1.depthResults, "119060324")
	testPositions = append(testPositions, testPos1)

	// test position 2: Kiwipete (complex middlegame)
	testPos2 := TestPosition{}
	testPos2.fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	testPos2.depthResults = append(testPos2.depthResults, "48")
	testPos2.depthResults = append(testPos2.depthResults, "2039")
	testPos2.depthResults = append(testPos2.depthResults, "97862")
	testPos2.depthResults = append(testPos2.depthResults, "4085603")
	testPos2.depthResults = append(testPos2.depthResults, "193690690")
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
	testPos3.depthResults = append(testPos3.depthResults, "178633661")
	testPositions = append(testPositions, testPos3)

	// test position 4: another complex middlegame, starting with check
	testPos4 := TestPosition{}
	testPos4.fen = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	testPos4.depthResults = append(testPos4.depthResults, "6")
	testPos4.depthResults = append(testPos4.depthResults, "264")
	testPos4.depthResults = append(testPos4.depthResults, "9467")
	testPos4.depthResults = append(testPos4.depthResults, "422333")
	testPos4.depthResults = append(testPos4.depthResults, "15833292")
	testPositions = append(testPositions, testPos4)

	// test position 5: promotions and castling
	testPos5 := TestPosition{}
	testPos5.fen = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
	testPos5.depthResults = append(testPos5.depthResults, "44")
	testPos5.depthResults = append(testPos5.depthResults, "1486")
	testPos5.depthResults = append(testPos5.depthResults, "62379")
	testPos5.depthResults = append(testPos5.depthResults, "2103487")
	testPos5.depthResults = append(testPos5.depthResults, "89941194")
	testPositions = append(testPositions, testPos5)

	// test position 6: symmetrical italian game
	testPos6 := TestPosition{}
	testPos6.fen = "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	testPos6.depthResults = append(testPos6.depthResults, "46")
	testPos6.depthResults = append(testPos6.depthResults, "2079")
	testPos6.depthResults = append(testPos6.depthResults, "89890")
	testPos6.depthResults = append(testPos6.depthResults, "3894594")
	testPos6.depthResults = append(testPos6.depthResults, "164075551")
	testPositions = append(testPositions, testPos6)

	// test position 7 (custom): custom game with lots of pins, checks, en-passants, start with black, and each piece type is still on the board
	testPos7 := TestPosition{}
	testPos7.fen = "8/1np1p1b1/2Bk3B/q6r/4K3/2Q3R1/2NP1P2/1b6 b - - 0 1"
	testPos7.depthResults = append(testPos7.depthResults, "41")
	testPos7.depthResults = append(testPos7.depthResults, "1375")
	testPos7.depthResults = append(testPos7.depthResults, "43366")
	testPos7.depthResults = append(testPos7.depthResults, "1469509")
	testPos7.depthResults = append(testPos7.depthResults, "46721049")
	testPositions = append(testPositions, testPos7)

	// test position 8 (custom): each: king and 3 pawns on starting squares, and both colored bishops
	testPos8 := TestPosition{}
	testPos8.fen = "8/1kpppbb1/8/8/8/8/1BBPPPK1/8 w - - 0 1"
	testPos8.depthResults = append(testPos8.depthResults, "30")
	testPos8.depthResults = append(testPos8.depthResults, "855")
	testPos8.depthResults = append(testPos8.depthResults, "22805")
	testPos8.depthResults = append(testPos8.depthResults, "594588")
	testPos8.depthResults = append(testPos8.depthResults, "14762722")
	testPos8.depthResults = append(testPos8.depthResults, "359874054")
	testPositions = append(testPositions, testPos8)

	// test position 9 (custom): each: king, queen, rook and 3 pawns
	testPos9 := TestPosition{}
	testPos9.fen = "8/4p3/3pkp1R/5q2/3Q4/r2PKP2/4P3/8 b - - 0 1"
	testPos9.depthResults = append(testPos9.depthResults, "28")
	testPos9.depthResults = append(testPos9.depthResults, "617")
	testPos9.depthResults = append(testPos9.depthResults, "14167")
	testPos9.depthResults = append(testPos9.depthResults, "343514")
	testPos9.depthResults = append(testPos9.depthResults, "8426450")
	testPos9.depthResults = append(testPos9.depthResults, "213938858")
	testPositions = append(testPositions, testPos9)

	// test position 10 (custom): each: king in front of pawns 2 columns of doubled pawns and one rook
	testPos10 := TestPosition{}
	testPos10.fen = "4r3/3pp3/2kpp3/8/4K3/3PP3/3PP3/3R4 w - - 0 1"
	testPos10.depthResults = append(testPos10.depthResults, "11")
	testPos10.depthResults = append(testPos10.depthResults, "140")
	testPos10.depthResults = append(testPos10.depthResults, "1989")
	testPos10.depthResults = append(testPos10.depthResults, "30869")
	testPos10.depthResults = append(testPos10.depthResults, "448878")
	testPos10.depthResults = append(testPos10.depthResults, "7248110")
	testPos10.depthResults = append(testPos10.depthResults, "107129005")
	testPositions = append(testPositions, testPos10)

	// test position 11 (custom): each: king in front of pawns 2 columns of doubled pawns and back rank of bishops
	testPos11 := TestPosition{}
	testPos11.fen = "bbbbbbbb/3pp3/2kpp3/8/8/3PPK2/3PP3/BBBBBBBB w - - 0 1"
	testPos11.depthResults = append(testPos11.depthResults, "30")
	testPos11.depthResults = append(testPos11.depthResults, "733")
	testPos11.depthResults = append(testPos11.depthResults, "18464")
	testPos11.depthResults = append(testPos11.depthResults, "457276")
	testPos11.depthResults = append(testPos11.depthResults, "11778803")
	testPos11.depthResults = append(testPos11.depthResults, "301306691")
	testPositions = append(testPositions, testPos11)

	// test position 12 (custom): each: discovered double checks, revealed checks
	testPos12 := TestPosition{}
	testPos12.fen = "8/4q3/4r3/2bbkp2/3n1N2/2BBKP2/4Q3/4R3 w - - 0 1"
	testPos12.depthResults = append(testPos12.depthResults, "37")
	testPos12.depthResults = append(testPos12.depthResults, "1209")
	testPos12.depthResults = append(testPos12.depthResults, "44574")
	testPos12.depthResults = append(testPos12.depthResults, "1544090")
	testPos12.depthResults = append(testPos12.depthResults, "58148575")
	testPositions = append(testPositions, testPos12)

	// test position 13 (ccc): rook captures remove castling rights
	testPos13 := TestPosition{}
	testPos13.fen = "r3k2r/8/8/8/3pPp2/8/8/R3K1RR b KQkq e3 0 1"
	testPos13.depthResults = append(testPos13.depthResults, "29")
	testPos13.depthResults = append(testPos13.depthResults, "829")
	testPos13.depthResults = append(testPos13.depthResults, "20501")
	testPos13.depthResults = append(testPos13.depthResults, "624871")
	testPos13.depthResults = append(testPos13.depthResults, "15446339")
	testPositions = append(testPositions, testPos13)

	// test position 14 (ccc): middlegame with en-passant checks, checkmates
	testPos14 := TestPosition{}
	testPos14.fen = "8/7p/p5pb/4k3/P1pPn3/8/P5PP/1rB2RK1 b - d3 0 28"
	testPos14.depthResults = append(testPos14.depthResults, "5")
	testPos14.depthResults = append(testPos14.depthResults, "117")
	testPos14.depthResults = append(testPos14.depthResults, "3293")
	testPos14.depthResults = append(testPos14.depthResults, "67197")
	testPos14.depthResults = append(testPos14.depthResults, "1881089")
	testPos14.depthResults = append(testPos14.depthResults, "38633283")
	testPositions = append(testPositions, testPos14)

	// test position 15 (ccc): promotions
	testPos15 := TestPosition{}
	testPos15.fen = "n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1"
	testPos15.depthResults = append(testPos15.depthResults, "24")
	testPos15.depthResults = append(testPos15.depthResults, "496")
	testPos15.depthResults = append(testPos15.depthResults, "9483")
	testPos15.depthResults = append(testPos15.depthResults, "182838")
	testPos15.depthResults = append(testPos15.depthResults, "3605103")
	testPositions = append(testPositions, testPos15)

}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------- Perft ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// count nodes visited
func (pos *Position) runPerft(initialDepth int, currentDepth int, bulkCounting bool) int {

	// check for depth limit
	if currentDepth == 0 {
		return 1
	}

	// reset the node count
	totalNodeCount := 0

	// generate legal moves
	pos.generateLegalMoves(false)

	// check for bulk-counting enhacements
	if bulkCounting && currentDepth == 1 {
		return pos.totalMovesCounter
	}

	// if there are legal moves, iterate over them
	//if pos.availableMovesCounter > 0 {
	if pos.totalMovesCounter > 0 {
		legalMoves := make([]Move, pos.totalMovesCounter)
		copy(legalMoves, pos.threatMoves[:pos.threatMovesCounter])
		copy(legalMoves[pos.threatMovesCounter:], pos.quietMoves[:pos.quietMovesCounter])

		for _, move := range legalMoves {
			pos.makeMove(move)
			currentNodeCount := pos.runPerft(initialDepth, currentDepth-1, bulkCounting)
			totalNodeCount += currentNodeCount
			pos.undoMove()

			//if initialDepth == currentDepth {
			//	fmt.Printf("For move (from:%v to:%v promote:%v) the node count is: %v\n", move.getFromSq(), move.getToSq(), move.getPromotionType(), currentNodeCount)
			//}
		}
	}

	// return the nodeCount
	return totalNodeCount
}

// loop over each perft position and print the results to the terminal
func printPerftTestResults() {

	bulkCounting := true

	fmt.Println(" ")

	if bulkCounting {
		fmt.Println("------------------------------------ Perft Test Results (Bulk Counting Enabled) ---------------------------------------")
	} else {
		fmt.Println("------------------------------------ Perft Test Results (Bulk Counting Disabled) ---------------------------------------")
	}

	for _, testPosition := range testPositions {

		fmt.Printf(" ------- New Position: %v -------\n", testPosition.fen)

		// create a new position from the test fen string
		newPos := Position{}
		newPos.initPositionFromFen(testPosition.fen)

		// start the time
		start_time := time.Now()
		totalNodes := 0

		// set a flag to catch errors
		moveGenSuccess := true

		// test each depth and print the results
		for depth, depthResults := range testPosition.depthResults {
			resultNodes := newPos.runPerft(depth+1, depth+1, bulkCounting) // run the perft
			totalNodes += resultNodes
			//fmt.Printf("Depth: %v. Correct test nodes: %v. My nodes: %v.\n", depth+1, depthResults, resultNodes)
			if depthResults != strconv.Itoa(resultNodes) {
				moveGenSuccess = false

				// code to debug the legal moves generated
				newPos.generateLegalMoves(false)
				fmt.Printf("Failure in move gen. Number of moves: %v. Generated moves:\n", newPos.totalMovesCounter)
				allMoves := make([]Move, newPos.totalMovesCounter)
				copy(allMoves, newPos.threatMoves[:newPos.threatMovesCounter])
				copy(allMoves[newPos.threatMovesCounter:], newPos.quietMoves[:newPos.quietMovesCounter])
				for _, move := range allMoves {
					fmt.Printf("<<MOVE>> From: %v. To: %v. Promote: %v.\n", move.getFromSq(), move.getToSq(), move.getPromotionType())
				}
			}
		}

		// record the time
		duration_time_sec := time.Since(start_time).Seconds()
		duration_time_ms := time.Since(start_time).Milliseconds()
		mnps := math.Round((float64(totalNodes) / (float64(duration_time_ms) / 1000)) / 1000000)

		// print the logging details
		// newPos.logOther.printLoggedDetails()

		// print the perft speed
		fmt.Printf("Completed perft in %v seconds. Speed: %v mnps.\n", duration_time_sec, mnps)

		// print the final result
		if moveGenSuccess {
			fmt.Printf("[SUCCESS!]\n")
		} else {
			fmt.Printf("[FAILURE!]\n")
		}
	}
}
