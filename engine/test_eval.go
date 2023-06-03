package main

import "fmt"

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Evaluation Tests -------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// we confirm that our incremental eval code works correctly (in make and undo move),
// by playing different move sequences to get the same end position, and check that the evaluations are the same
// we also compare against a position created fresh from the fen string of the final position

type EvalTestSequence struct {
	fenMoves1 []string // first move sequence
	fenMoves2 []string // second move sequence
	fenString string   // fen string of final position
}

var evalTestSequences []EvalTestSequence

func initEvalTestSequences() {

	// test 1
	moves01_1 := []string{"e2e4", "e7e5"}
	moves01_2 := []string{"e2e3", "e7e6", "e3e4", "e6e5"}
	fen01 := "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1"
	test01 := EvalTestSequence{moves01_1, moves01_2, fen01}
	evalTestSequences = append(evalTestSequences, test01)

	// test 2
	moves02_1 := []string{"e2e4", "c7c5"}
	moves02_2 := []string{"e2e3", "c7c6", "e3e4", "c6c5"}
	fen02 := "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1"
	test02 := EvalTestSequence{moves02_1, moves02_2, fen02}
	evalTestSequences = append(evalTestSequences, test02)

	// test 3
	moves03_1 := []string{"e2e4", "e7e5", "g1f3", "g8f6", "f3e5", "f6e4", "e5d7", "e4d2", "d1d2", "d8d7"}
	moves03_2 := []string{"e2e4", "e7e5", "g1f3", "e8e7", "f3e5", "e7e8", "e5d7", "g8f6", "f1b5", "f6e4", "b5f1", "e4d2", "d1d2", "d8d7"}
	fen03 := "rnb1kb1r/pppq1ppp/8/8/8/8/PPPQ1PPP/RNB1KB1R w KQkq - 0 1"
	test03 := EvalTestSequence{moves03_1, moves03_2, fen03}
	evalTestSequences = append(evalTestSequences, test03)

	// test 4
	moves04_1 := []string{"e2e4", "e7e5", "g1f3", "g8f6", "f3e5", "f6e4", "e5d7", "e4d2", "d1d2", "d8d7", "d2d7", "e8d7", "e1d2"}
	moves04_2 := []string{"e2e4", "e7e5", "g1f3", "e8e7", "f3e5", "e7e8", "e5d7", "g8f6", "f1b5", "f6e4", "b5f1", "e4d2", "d1d2", "d8d7", "d2a5", "d7d8", "a5d2", "d8d7", "d2a5", "d7d8", "a5d2", "d8d7", "d2a5", "d7d8", "a5d2", "d8d7", "d2d7", "e8d7", "e1d2"}
	fen04 := "rnb2b1r/pppk1ppp/8/8/8/8/PPPK1PPP/RNB2B1R b - - 0 1"
	test04 := EvalTestSequence{moves04_1, moves04_2, fen04}
	evalTestSequences = append(evalTestSequences, test04)

	// test 5
	moves05_1 := []string{"e2e4", "d7d5", "e4d5"}
	moves05_2 := []string{"e2e3", "d7d6", "e3e4", "d6d5", "e4d5"}
	fen05 := "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"
	test05 := EvalTestSequence{moves05_1, moves05_2, fen05}
	evalTestSequences = append(evalTestSequences, test05)

	// test 6
	moves06_1 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "d8d6", "f3d2", "d6e6", "d1c2", "e6g4", "g2g3", "a8d8", "f1e1", "g4h3", "e3e4", "h3e6", "b2c3", "f6e4", "d3e4", "d5e4", "d2e4", "e6f5", "c2e2", "h7h6", "a1b1", "a7a6", "b1b7", "f5d5", "b7c7", "c6d4", "c3d4", "d5d4", "e2f3", "e8f8", "f3f5", "d4a4", "f5f3", "a4a5", "f3c3", "a5c3", "e4c3", "f8e8", "e1e8", "d8e8", "c3d5", "e8a8", "h2h3", "g7g6", "g1g2", "a6a5", "a3a4", "g8g7", "d5e3", "g7g8", "e3g4", "h6h5", "g4h6", "g8g7", "h6f7", "g7f6", "f2f4", "a8b8", "f7g5", "b8b2", "g2f3", "b2b3", "f3e4", "b3b4", "e4d5", "b4b5", "a4b5", "f6f5", "c7f7"}
	moves06_2 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "d8d6", "f3d2", "d6e6", "d1c2", "e6g4", "g2g3", "a8d8", "f1e1", "g4h3", "e3e4", "h3e6", "b2c3", "f6e4", "d3e4", "d5e4", "d2e4", "e6f5", "c2e2", "h7h6", "a1b1", "a7a6", "b1b7", "f5d5", "b7c7", "c6d4", "c3d4", "d5d4", "e2f3", "e8f8", "f3f5", "d4a4", "f5f3", "a4a5", "f3c3", "a5c3", "e4c3", "f8e8", "e1e8", "d8e8", "c3d5", "e8a8", "h2h3", "g7g6", "g1g2", "a6a5", "a3a4", "g8g7", "d5e3", "g7g8", "e3g4", "h6h5", "g4h6", "g8g7", "h6f7", "g7f6", "f2f4", "a8b8", "f7g5", "b8b2", "g2f3", "b2b3", "f3e4", "b3b4", "e4d5", "b4b5", "a4b5", "f6f5", "c7f7"}
	fen06 := "8/5R2/6p1/pP1K1kNp/5P2/6PP/8/8 b - - 2 47"
	test06 := EvalTestSequence{moves06_1, moves06_2, fen06}
	evalTestSequences = append(evalTestSequences, test06)

	// test 7
	moves07_1 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "f6e4", "f1e1", "d8d6", "f3d2", "e4d2", "d1d2", "d5e4", "d3e4", "e8e4", "d2c3", "d6d5", "a1d1", "e4g4", "g2g3", "g4e4", "h2h3", "a8e8", "c3c2", "e8d8", "c2c5", "d5c5", "d4c5", "d8d1", "e1d1", "e4e7", "d1c1", "e7e6", "g1g2", "b7b6", "b2d4", "b6c5", "c1c5", "c6d4", "e3d4", "e6e7", "d4d5", "f7f5", "c5a5", "e7d7", "a5a7", "d7d5", "a7c7", "h7h5", "h3h4", "d5a5", "c7c3", "g8f7", "f2f4", "g7g6", "g2f3", "f7e6", "f3e2", "e6d5", "e2e3", "a5a6", "e3d3", "a6a4", "d3e3", "a4a6", "c3b3", "d5c4", "b3d3", "a6e6", "e3d2", "e6b6", "d3c3", "c4d4", "c3d3", "d4c4", "d3c3", "c4d4", "d2c2", "b6a6", "c3d3", "d4c4", "d3e3", "c4d4", "c2d2", "a6b6", "d2e2", "d4c4", "e2f3", "b6b3", "f3e2", "b3c3", "a3a4", "c3e3", "e2e3", "c4c5", "a4a5", "c5b5", "e3d4", "b5a5", "d4e5", "a5b4", "e5f6", "b4c4", "f6g6", "c4d5", "g6f5", "d5d4", "f5g6", "d4e4", "f4f5", "e4f3", "g6h5", "f3g3", "f5f6", "g3f4", "f6f7", "f4e5", "f7f8q", "e5e4", "h5g4", "e4d5", "h4h5", "d5e5", "h5h6", "e5d5", "h6h7", "d5c4", "h7h8q", "c4d5", "f8e8", "d5c5", "h8e5", "c5b6", "e8b5", "b6a7", "e5a1"}
	moves07_2 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "f6e4", "f1e1", "d8d6", "f3d2", "e4d2", "d1d2", "d5e4", "d3e4", "e8e4", "d2c3", "d6d5", "a1d1", "e4g4", "g2g3", "g4e4", "h2h3", "a8e8", "c3c2", "e8d8", "c2c5", "d5c5", "d4c5", "d8d1", "e1d1", "e4e7", "d1c1", "e7e6", "g1g2", "b7b6", "b2d4", "b6c5", "c1c5", "c6d4", "e3d4", "e6e7", "d4d5", "f7f5", "c5a5", "e7d7", "a5a7", "d7d5", "a7c7", "h7h5", "h3h4", "d5a5", "c7c3", "g8f7", "f2f4", "g7g6", "g2f3", "f7e6", "f3e2", "e6d5", "e2e3", "a5a6", "e3d3", "a6a4", "d3e3", "a4a6", "c3b3", "d5c4", "b3d3", "a6e6", "e3d2", "e6b6", "d3c3", "c4d4", "c3d3", "d4c4", "d3c3", "c4d4", "d2c2", "b6a6", "c3d3", "d4c4", "d3e3", "c4d4", "c2d2", "a6b6", "d2e2", "d4c4", "e2f3", "b6b3", "f3e2", "b3c3", "a3a4", "c3e3", "e2e3", "c4c5", "a4a5", "c5b5", "e3d4", "b5a5", "d4e5", "a5b4", "e5f6", "b4c4", "f6g6", "c4d5", "g6f5", "d5d4", "f5g6", "d4e4", "f4f5", "e4f3", "g6h5", "f3g3", "f5f6", "g3f4", "f6f7", "f4e5", "f7f8q", "e5e4", "h5g4", "e4d5", "h4h5", "d5e5", "h5h6", "e5d5", "h6h7", "d5c4", "h7h8q", "c4d5", "f8e8", "d5c5", "h8e5", "c5b6", "e8b5", "b6a7", "e5a1"}
	fen07 := "8/k7/8/1Q6/6K1/8/8/Q7 b - - 8 78"
	test07 := EvalTestSequence{moves07_1, moves07_2, fen07}
	evalTestSequences = append(evalTestSequences, test07)

	// test 8
	moves08_1 := []string{"e2e4", "d7d5", "e4d5", "d8d5", "b1c3", "d5f5", "g1f3", "b8c6", "d2d4", "g8f6", "d4d5", "c6b4", "f1d3", "b4d3", "d1d3", "f5d3", "c2d3", "e7e6", "d5e6", "c8e6", "c1e3", "f8d6", "a2a4", "e8g8", "e1g1", "c7c5", "c3b5", "d6e7", "b5c7", "a8d8", "c7e6", "f7e6", "f3g5", "f6g4", "g5e6", "g4e3", "f2e3", "f8f1", "a1f1", "d8d3", "e3e4", "d3b3", "f1f5", "b7b6", "e4e5", "b3b2", "e6c7", "c5c4", "g1f1", "b2a2", "c7d5", "e7c5", "e5e6", "a2a1", "f1e2", "a1a2", "e2d1", "a2a1", "d1c2", "a1e1", "c2c3", "e1e6", "c3c4", "e6e2", "f5g5", "e2e4", "c4b3", "g8f7", "h2h3", "a7a5", "g5f5", "f7e6", "f5h5", "h7h6", "d5c3", "e4b4", "b3c2", "c5d4", "h5b5", "g7g5", "b5b4", "a5b4", "c3a2", "h6h5", "g2g4", "h5g4", "h3g4", "d4c5", "c2b3", "e6d5", "a2b4", "d5d4", "b4c6", "d4e4", "b3c4", "c5g1", "c4b5", "e4f4", "c6a7", "f4g4", "a7c8", "g4f5", "c8b6", "g1b6", "b5b6", "g5g4", "a4a5", "g4g3", "b6b7", "g3g2", "a5a6", "g2g1q", "a6a7", "g1b1", "b7c7", "b1e4", "c7b8", "e4b4", "b8c8", "b4a4", "c8b7", "a4b5", "b7a8", "b5a6", "a8b8", "a6b6", "b8a8", "b6a6", "a8b8", "a6b6", "b8a8"}
	moves08_2 := []string{"e2e4", "d7d5", "e4d5", "d8d5", "b1c3", "d5f5", "g1f3", "b8c6", "d2d4", "g8f6", "d4d5", "c6b4", "f1d3", "b4d3", "d1d3", "f5d3", "c2d3", "e7e6", "d5e6", "c8e6", "c1e3", "f8d6", "a2a4", "e8g8", "e1g1", "c7c5", "c3b5", "d6e7", "b5c7", "a8d8", "c7e6", "f7e6", "f3g5", "f6g4", "g5e6", "g4e3", "f2e3", "f8f1", "a1f1", "d8d3", "e3e4", "d3b3", "f1f5", "b7b6", "e4e5", "b3b2", "e6c7", "c5c4", "g1f1", "b2a2", "c7d5", "e7c5", "e5e6", "a2a1", "f1e2", "a1a2", "e2d1", "a2a1", "d1c2", "a1e1", "c2c3", "e1e6", "c3c4", "e6e2", "f5g5", "e2e4", "c4b3", "g8f7", "h2h3", "a7a5", "g5f5", "f7e6", "f5h5", "h7h6", "d5c3", "e4b4", "b3c2", "c5d4", "h5b5", "g7g5", "b5b4", "a5b4", "c3a2", "h6h5", "g2g4", "h5g4", "h3g4", "d4c5", "c2b3", "e6d5", "a2b4", "d5d4", "b4c6", "d4e4", "b3c4", "c5g1", "c4b5", "e4f4", "c6a7", "f4g4", "a7c8", "g4f5", "c8b6", "g1b6", "b5b6", "g5g4", "a4a5", "g4g3", "b6b7", "g3g2", "a5a6", "g2g1q", "a6a7", "g1b1", "b7c7", "b1e4", "c7b8", "e4b4", "b8c8", "b4a4", "c8b7", "a4b5", "b7a8", "b5a6", "a8b8", "a6b6", "b8a8", "b6a6", "a8b8", "a6b6", "b8a8"}
	fen08 := "K7/P7/1q6/5k2/8/8/8/8 b - - 18 66"
	test08 := EvalTestSequence{moves08_1, moves08_2, fen08}
	evalTestSequences = append(evalTestSequences, test08)

	// test 9
	moves09_1 := []string{"e2e4", "d7d5", "e4d5", "d8d5", "g1e2", "b8c6", "b1c3", "d5d8", "d2d4", "g8f6", "c1g5", "e7e6", "c3e4", "f8e7", "e4f6", "e7f6", "g5f6", "d8f6", "c2c3", "e8g8", "e2g3", "e6e5", "d4d5", "c6e7", "f1c4", "f6f4", "b2b3", "c8g4", "g3e2", "g4e2", "d1e2", "f8d8", "e2d3", "c7c6", "d5d6", "b7b5", "d3c2", "d8d6", "c4d3", "d6h6", "g2g3", "f4f6", "e1g1", "e7d5", "a1d1", "f6f3", "f1e1", "a8e8", "a2a3", "a7a5", "c3c4", "b5c4", "b3c4", "d5f6", "c2e2", "f3h5", "e2h5", "h6h5", "d3e2", "h5f5", "d1d6", "f5f2", "d6c6", "f2f5", "c6c5", "e8a8", "e2f1", "a5a4", "c5e5", "f5e5", "e1e5", "h7h6", "c4c5", "g7g5", "c5c6", "g5g4", "e5f5", "g8g7", "f5f4", "h6h5", "f1b5", "f6d5", "f4a4", "a8a4", "b5a4", "f7f5", "a4b3", "d5c7", "h2h4", "g7f6", "g1f2", "f6e7", "b3c2", "e7d6", "c2f5", "c7d5", "f5g6", "d5f6", "f2e3", "d6c6", "e3f4", "c6d5", "f4g5", "d5e5", "g6h5", "f6e4", "g5g4", "e4f6", "g4g5", "f6e4", "g5g6", "e4g3", "h5d1", "g3f5", "h4h5", "f5e7", "g6g5", "e7d5", "h5h6", "d5f6", "d1c2", "e5e6", "c2g6", "e6e7", "g5f5", "f6h7", "g6h7", "e7f7", "f5e5", "f7f8", "e5f6", "f8e8", "h7f5", "e8d8", "h6h7", "d8c7", "h7h8q", "c7c6", "f6e5", "c6b6", "h8b8", "b6c6", "a3a4", "c6c5", "b8b5"}
	moves09_2 := []string{"e2e4", "d7d5", "e4d5", "d8d5", "g1e2", "b8c6", "b1c3", "d5d8", "d2d4", "g8f6", "c1g5", "e7e6", "c3e4", "f8e7", "e4f6", "e7f6", "g5f6", "d8f6", "c2c3", "e8g8", "e2g3", "e6e5", "d4d5", "c6e7", "f1c4", "f6f4", "b2b3", "c8g4", "g3e2", "g4e2", "d1e2", "f8d8", "e2d3", "c7c6", "d5d6", "b7b5", "d3c2", "d8d6", "c4d3", "d6h6", "g2g3", "f4f6", "e1g1", "e7d5", "a1d1", "f6f3", "f1e1", "a8e8", "a2a3", "a7a5", "c3c4", "b5c4", "b3c4", "d5f6", "c2e2", "f3h5", "e2h5", "h6h5", "d3e2", "h5f5", "d1d6", "f5f2", "d6c6", "f2f5", "c6c5", "e8a8", "e2f1", "a5a4", "c5e5", "f5e5", "e1e5", "h7h6", "c4c5", "g7g5", "c5c6", "g5g4", "e5f5", "g8g7", "f5f4", "h6h5", "f1b5", "f6d5", "f4a4", "a8a4", "b5a4", "f7f5", "a4b3", "d5c7", "h2h4", "g7f6", "g1f2", "f6e7", "b3c2", "e7d6", "c2f5", "c7d5", "f5g6", "d5f6", "f2e3", "d6c6", "e3f4", "c6d5", "f4g5", "d5e5", "g6h5", "f6e4", "g5g4", "e4f6", "g4g5", "f6e4", "g5g6", "e4g3", "h5d1", "g3f5", "h4h5", "f5e7", "g6g5", "e7d5", "h5h6", "d5f6", "d1c2", "e5e6", "c2g6", "e6e7", "g5f5", "f6h7", "g6h7", "e7f7", "f5e5", "f7f8", "e5f6", "f8e8", "h7f5", "e8d8", "h6h7", "d8c7", "h7h8q", "c7c6", "f6e5", "c6b6", "h8b8", "b6c6", "a3a4", "c6c5", "b8b5"}
	fen09 := "8/8/8/1Qk1KB2/P7/8/8/8 b - - 2 73"
	test09 := EvalTestSequence{moves09_1, moves09_2, fen09}
	evalTestSequences = append(evalTestSequences, test09)

}

func printEvalTestResults() {

	fmt.Println("-------------------------------- Eval Test Results ----------------------------------")

	for i, testSequence := range evalTestSequences {

		// create new position
		pos1 := Position{}
		pos2 := Position{}
		pos3 := Position{}

		pos1.initPositionFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		pos2.initPositionFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		pos3.initPositionFromFen(testSequence.fenString)

		// play the moves in the first position
		for _, moveStr := range testSequence.fenMoves1 {
			pos1.makeUCIMove(moveStr)
		}

		// play the moves in the second position
		for _, moveStr := range testSequence.fenMoves2 {
			pos2.makeUCIMove(moveStr)
		}

		// evaluate the positions and return the results
		pos1.evalPosAfter()
		pos1MaterialEval := pos1.evalMaterial
		pos1HeatmapEval := pos1.evalHeatmaps
		pos1OtherEval := pos1.evalOther
		pos1EvalTotal := pos1MaterialEval + pos1HeatmapEval + pos1OtherEval

		pos2.evalPosAfter()
		pos2MaterialEval := pos2.evalMaterial
		pos2HeatmapEval := pos2.evalHeatmaps
		pos2OtherEval := pos2.evalOther
		pos2EvalTotal := pos2MaterialEval + pos2HeatmapEval + pos2OtherEval

		pos3.evalPosAfter()
		pos3MaterialEval := pos3.evalMaterial
		pos3HeatmapEval := pos3.evalHeatmaps
		pos3OtherEval := pos3.evalOther
		pos3EvalTotal := pos3MaterialEval + pos3HeatmapEval + pos3OtherEval

		success := pos1EvalTotal == pos2EvalTotal && pos2EvalTotal == pos3EvalTotal
		if success {
			fmt.Printf("Test %v: Success!\n", i+1)
		} else {
			fmt.Printf("Test %v: Failure! Pos 1 eval: %v (material: %v, heatmap: %v, other: %v). Pos 2 eval: %v (material: %v, heatmap: %v, other: %v). Pos 3 eval: %v (material: %v, heatmap: %v, other: %v).\n",
				i+1, pos1EvalTotal, pos1MaterialEval, pos1HeatmapEval, pos1OtherEval, pos2EvalTotal, pos2MaterialEval, pos2HeatmapEval, pos2OtherEval, pos3EvalTotal, pos3MaterialEval, pos3HeatmapEval, pos3OtherEval)
		}
	}
}
