package main

import "fmt"

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Evaluation Tests -------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// we confirm that our incremental eval and hashing code works correctly (in make and undo move),
// by playing a move sequence to get to a position,
// and compare that to the eval and hash when the final fen string is loaded directly

type IncrementalTestSequence struct {
	fenMoves  []string // fen moves from the starting to final position
	fenString string   // fen string of final position
}

var incrementalTestSequences []IncrementalTestSequence

func initIncrementalTestSequences() {

	// test 1
	moves01 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "d8d6", "f3d2", "d6e6", "d1c2", "e6g4", "g2g3", "a8d8", "f1e1", "g4h3", "e3e4", "h3e6", "b2c3", "e6g4", "e4d5", "e8e1", "a1e1", "c6d4", "c2b2", "d4f3", "d2f3", "g4f3", "c3f6", "f3f6", "b2b7", "f6c3", "e1e3", "c3c1", "d3f1", "c1c5", "b7c6", "c5c6", "d5c6", "g8f8", "f1d3", "g7g6", "g1g2", "f7f5", "d3e2", "a7a5", "e3e5", "d8a8", "e2c4", "a8a7", "e5b5", "a7a8", "b5d5", "a8a7", "d5d7", "h7h6", "a3a4", "g6g5", "c4d3", "f5f4", "g3f4", "g5f4", "g2f3", "a7a8", "f3f4", "a8b8", "d3b5", "b8c8", "h2h4", "h6h5", "f4g5", "f8e8", "b5c4", "c8d8", "c4f7", "e8f8", "g5f6", "d8c8", "d7d1", "c8e8", "d1g1", "e8e6", "f7e6", "f8e8", "g1g8"}
	fen01 := "4k1R1/2p5/2P1BK2/p6p/P6P/8/5P2/8 b - - 2 52"
	test01 := IncrementalTestSequence{moves01, fen01}
	incrementalTestSequences = append(incrementalTestSequences, test01)

	// test 2
	moves02 := []string{"d2d4", "d7d5", "b1c3", "g8f6", "g1f3", "e7e6", "e2e3", "f8b4", "a2a3", "b4c3", "b2c3", "b8c6", "f1d3", "e8g8", "e1g1", "f8e8", "c3c4", "c8d7", "c4d5", "e6d5", "c2c4", "d7e6", "c4d5", "e6d5", "c1b2", "f6e4", "f1e1", "d8d6", "f3d2", "e4d2", "d1d2", "d5e4", "d3e4", "e8e4", "d2c3", "d6d5", "a1d1", "e4g4", "g2g3", "g4e4", "h2h3", "a8e8", "c3c2", "e8d8", "c2c5", "d5c5", "d4c5", "d8d1", "e1d1", "e4e7", "d1c1", "e7e6", "g1g2", "b7b6", "b2d4", "b6c5", "c1c5", "c6d4", "e3d4", "e6e7", "d4d5", "f7f5", "c5a5", "e7d7", "a5a7", "d7d5", "a7c7", "h7h5", "h3h4", "d5a5", "c7c3", "g8f7", "f2f4", "g7g6", "g2f3", "f7e6", "f3e2", "e6d6", "e2e3", "a5a6", "c3d3", "d6c5", "d3c3", "c5d5", "e3d3", "a6a4", "d3e3", "a4e4", "e3f3", "e4a4", "c3d3", "d5c4", "f3e3", "a4a6", "e3d2", "a6b6", "d3c3", "c4d4", "c3d3", "d4c4", "d3c3", "c4d4", "d2c2", "b6a6", "c3b3", "d4e4", "c2d2", "e4d4", "d2e2", "d4c4", "b3e3", "a6b6", "e3e8", "b6b2", "e2f1", "b2b3", "f1g2", "b3a3", "e8e6", "a3a2", "g2h3", "c4d5", "e6g6", "d5e4", "g6g5", "a2a1", "h3g2", "a1a2", "g2h3", "a2a1", "g5h5", "e4f3", "h3h2", "a1a2", "h2h3", "a2a1", "h3h2", "a1a2", "h2h3"}
	fen02 := "8/8/8/5p1R/5P1P/5kPK/r7/8 b - - 8 70"
	test02 := IncrementalTestSequence{moves02, fen02}
	incrementalTestSequences = append(incrementalTestSequences, test02)

	// test 3
	moves03 := []string{"e2e4", "d7d5", "e4d5", "d8d5", "b1c3", "d5f5", "g1f3", "b8c6", "d2d4", "g8f6", "d4d5", "c6b4", "f1d3", "b4d3", "d1d3", "f5d3", "c2d3", "e7e6", "d5e6", "c8e6", "c1e3", "f8d6", "a2a4", "e8g8", "e1g1", "c7c5", "c3b5", "d6e7", "b5c7", "a8d8", "c7e6", "f7e6", "a1d1", "e7d6", "d3d4", "f6d5", "f1e1", "d5e3", "f2e3", "f8e8", "b2b3", "a7a5", "h2h3", "b7b6", "g2g4", "h7h6", "e1f1", "d6c7", "f3e5", "c7e5", "d4e5", "d8d5", "d1d5", "e6d5", "f1f5", "c5c4", "b3c4", "d5c4", "g1f2", "c4c3", "f2e2", "g7g6", "f5f6", "g8g7", "f6b6", "e8e5", "e2d3", "e5c5", "b6b5", "c5c4", "d3c2", "c4a4", "c2c3", "a4a1", "e3e4", "a5a4", "e4e5", "g7f7", "b5a5", "a4a3", "c3b3", "a1h1", "b3a3", "h1a1", "a3b4", "a1a5", "b4a5", "g6g5", "a5a4", "f7e6", "a4b3", "e6e5", "b3c3", "e5f4", "c3d4", "f4g3", "d4e4", "g3h3", "e4f5", "h3g3", "f5g6", "g3g4", "g6h6", "g4f4", "h6g6", "g5g4", "g6f6", "g4g3", "f6e6", "f4e4", "e6d6", "g3g2", "d6c5", "g2g1q", "c5b4", "e4d4", "b4b5", "g1g6", "b5a4", "d4c4", "a4a3", "g6b1", "a3a4", "b1a1"}
	fen03 := "8/8/8/8/K1k5/8/8/q7 w - - 10 63"
	test03 := IncrementalTestSequence{moves03, fen03}
	incrementalTestSequences = append(incrementalTestSequences, test03)

	// test 4
	moves04 := []string{"e2e4", "b8c6", "d2d4", "d7d5", "b1c3", "e7e6", "g1f3", "f8b4", "c1g5", "g8e7", "f1d3", "h7h6", "g5e7", "b4e7", "e1g1", "e8g8", "f1e1", "d5e4", "d3e4", "f8e8", "c3e2", "e7d6", "c2c4", "c8d7", "c4c5", "d6e7", "e2f4", "g7g5", "f4d3", "a7a5", "h2h3", "a8b8", "d3e5", "e7f6", "e5d7", "d8d7", "d1a4", "b8d8", "a1d1", "c6d4", "c5c6", "d4f3", "e4f3", "d7e7", "d1d8", "e8d8", "c6b7", "e7b4", "a4b4", "a5b4", "e1d1", "d8f8", "a2a3", "c7c5", "a3b4", "c5b4", "d1d6", "f6e7", "d6c6", "g8g7", "c6c8", "e7d6", "f3c6", "f7f5", "c6e8", "g7f6", "e8d7", "f6e7", "d7e6", "f5f4", "e6c4", "f8d8", "b2b3", "h6h5", "f2f3", "h5h4", "g1f2", "e7d7", "c4b5", "d7e7", "b5c4", "d6e5", "f2e2", "e5d6"}
	fen04 := "2Rr4/1P2k3/3b4/6p1/1pB2p1p/1P3P1P/4K1P1/8 w - - 8 43"
	test04 := IncrementalTestSequence{moves04, fen04}
	incrementalTestSequences = append(incrementalTestSequences, test04)

	// test 5
	moves05 := []string{"e2e4", "c7c5", "g1f3", "b8c6", "b1c3", "g8f6", "f1b5", "c6d4", "e4e5", "d4b5", "c3b5", "f6e4", "e1g1", "d8a5", "b5a3", "c5c4", "a3c4", "a5c7", "d2d3", "e4c5", "c4e3", "d7d6", "d3d4", "c5d7", "e3d5", "c7c4", "d5e3", "c4c7", "e3d5", "c7c4", "d5f4", "d6e5", "d4e5", "e7e6", "b2b3", "c4c7", "c1b2", "f8c5", "d1d3", "e8g8", "f3g5", "g7g6", "f4e6", "f7e6", "g5e6", "d7e5", "e6c7", "e5d3", "c2d3", "a8b8", "c7d5", "c8f5", "f1c1", "f8c8", "d5f6", "g8f7", "b2e5", "b8a8", "d3d4", "c5e7", "f6d5", "e7a3", "c1c8", "a8c8", "d5e3", "b7b5", "e3f5", "g6f5", "f2f4", "c8c2", "e5b8", "a3c1", "b8a7", "c1e3", "g1h1", "e3f4", "a7c5", "c2d2", "a2a4", "b5a4", "g2g3", "f4g3", "h2g3", "a4b3", "a1b1", "d2d3", "h1g2", "f7f6", "c5b4", "f6e6", "b1b2", "e6d5", "b4c5", "d3c3", "b2b1", "h7h5", "g2f2", "c3d3", "f2g2", "d3c3", "g2f2", "c3d3", "f2g2"}
	fen05 := "8/8/8/2Bk1p1p/3P4/1p1r2P1/6K1/1R6 b - - 7 52"
	test05 := IncrementalTestSequence{moves05, fen05}
	incrementalTestSequences = append(incrementalTestSequences, test05)

	// test 6
	moves06 := []string{"e2e4", "c7c5", "g1f3", "b8c6", "b1c3", "g8f6", "f1b5", "c6d4", "e4e5", "d4b5", "c3b5", "f6e4", "e1g1", "d8a5", "a2a4", "d7d5", "d2d3", "a7a6", "b5a3", "a5a4", "d3e4", "d5e4", "b2b3", "a4d7", "d1d7", "c8d7", "f3d2", "f7f5", "e5f6", "d7c6", "f1e1", "e7f6", "d2e4", "e8c8", "c1f4", "d8d4", "f2f3", "f6f5", "f4e3", "f5e4", "e3d4", "c5d4", "f3e4", "f8b4", "e1e2", "b4c5", "g1h1", "h8e8", "e4e5", "b7b5", "a3b1", "c8b7", "b1d2", "g7g5", "e5e6", "g5g4", "d2e4", "c5e7", "h1g1", "b5b4", "a1a5", "e8g8", "a5f5", "c6e8", "e4f2", "b7b6", "f5f7", "e7c5", "f7h7", "e8g6", "h7h6", "c5e7", "f2g4", "g6c2", "e2c2", "g8g4", "c2c4", "g4g8", "g1f1", "e7f8", "h6h7", "f8g7", "e6e7", "g7h8", "h7h6", "b6a7", "h6d6", "g8e8", "d6d7", "a7b6", "c4b4", "b6c6", "d7d8", "e8e7", "d8h8", "e7e3", "h8h3", "e3h3", "g2h3", "c6c5", "b4a4", "c5d5", "a4a6", "d4d3", "f1e1", "d5e4", "e1d2", "e4f3", "b3b4", "f3e4", "b4b5", "e4e5", "h3h4", "e5d5", "a6h6", "d5c4", "b5b6", "c4d4", "b6b7", "d4d5", "d2d3", "d5e5", "b7b8q", "e5f5", "d3e3", "f5g4", "b8c8"}
	fen06 := "2Q5/8/7R/8/6kP/4K3/7P/8 b - - 4 64"
	test06 := IncrementalTestSequence{moves06, fen06}
	incrementalTestSequences = append(incrementalTestSequences, test06)

	// test 7
	moves07 := []string{"e2e4", "c7c5", "c2c3", "d7d5", "f1b5", "b8c6", "e4d5", "d8d5", "g1f3", "a7a6", "b5c6", "b7c6", "d2d4", "g8f6", "e1g1", "c8f5", "d4c5", "d5c5", "c1e3", "c5b5", "b2b3", "e7e5", "c3c4", "b5b8", "e3g5", "f8e7", "f1e1", "f5g4", "b1d2", "e8g8", "d1c2", "e7d6", "c4c5", "d6c7", "h2h3", "g4f3", "d2f3", "f6d5", "g5c1", "f8e8", "c1b2", "h7h6", "a1d1", "d5b4", "c2c4", "b4a2", "e1e4", "a6a5", "d1d7", "e8f8", "b2e5", "a8a7", "e4f4", "c7e5", "f3e5", "a7d7", "e5d7", "b8b4", "d7f8", "b4c4", "b3c4", "g8f8", "f4f3", "a2c1", "f3a3", "f8e7", "g2g4", "a5a4", "a3a4", "c1d3", "a4a6", "d3c5", "a6c6", "c5e6", "g1g2", "g7g6", "g2g3", "f7f5", "f2f4", "e7f6", "h3h4", "f5g4", "g3g4", "h6h5", "g4f3", "f6f7", "f3e3", "e6g7", "c6c7", "f7f6", "c7c6", "f6f7", "c6c7", "f7f6", "e3e4", "g7f5", "c7c6", "f6g7", "c6b6", "f5h4", "c4c5", "g7f7", "c5c6", "f7e7", "e4e5", "e7d8", "b6b7", "d8c8", "b7g7", "c8d8", "e5e6", "d8c8", "c6c7", "h4f5", "g7g6", "f5e3", "e6d6", "e3f5", "d6e5", "f5e3", "g6g3", "e3f1", "f4f5", "h5h4", "g3g7", "f1g3", "f5f6", "g3h5", "e5d6", "h5g7", "f6g7", "c8b7", "g7g8q", "b7a6", "g8d5", "a6b6", "c7c8q", "b6a7", "d5d4"}
	fen07 := "2Q5/k7/3K4/8/3Q3p/8/8/8 b - - 2 70"
	test07 := IncrementalTestSequence{moves07, fen07}
	incrementalTestSequences = append(incrementalTestSequences, test07)

	// test 8
	moves08 := []string{"e2e4", "c7c5", "c2c3", "d7d5", "f1b5", "b8c6", "e4d5", "d8d5", "g1f3", "a7a6", "b5c6", "b7c6", "d2d4", "g8f6", "e1g1", "c8f5", "d4c5", "d5c5", "c1e3", "c5b5", "b2b3", "e7e5", "c3c4", "b5b8", "e3g5", "f8e7", "f1e1", "f5g4", "b1d2", "e8g8", "d1c2", "e7d6", "c4c5", "d6c7", "g5f6", "g7f6", "h2h3", "g4d7", "d2e4", "b8d8", "a1d1", "g8g7", "c2c3", "d8e7", "c3b4", "a8d8", "b4b7", "c7b8", "b7a6", "g7g8", "d1d7", "d8d7", "a6c6", "g8g7", "e4f6", "d7c7", "f6h5", "g7h8", "c6h6", "f7f6", "h5f4", "h8g8", "f4d5", "e7g7", "h6c1", "c7f7", "c1c4", "g8h8", "a2a4", "f6f5", "a4a5", "f8g8", "g2g3", "e5e4", "f3d4", "b8g3", "g1f1", "g3f2", "f1f2", "g7g2", "f2e3", "f5f4", "d5f4", "g8g3", "d4f3", "g3f3", "e3d4", "f7f4", "d4e5", "g2g5", "e5d6", "f4f6", "d6c7", "g5g7", "c4f7", "g7f7", "c7c8", "f7e7", "e1d1", "f6f8", "d1d8", "f8d8"}
	fen08 := "2Kr3k/4q2p/8/P1P5/4p3/1P3r1P/8/8 w - - 0 52"
	test08 := IncrementalTestSequence{moves08, fen08}
	incrementalTestSequences = append(incrementalTestSequences, test08)

	// test 9
	moves09 := []string{"e2e4", "c7c5", "b1c3", "b8c6", "g2g3", "e7e5", "g1f3", "g8f6", "f1b5", "d7d6", "e1g1", "f8e7", "d2d3", "e8g8", "c3d5", "f6d5", "e4d5", "c6d4", "f3d4", "c5d4", "c2c3", "d8a5", "d1a4", "a5a4", "b5a4", "c8f5", "c3d4", "e5d4", "f1d1", "a7a5", "a4b5", "f8d8", "b5c4", "a5a4", "b2b4", "a4b3", "c4b3", "b7b5", "h2h4", "e7f6", "c1g5", "a8c8", "d1d2", "f6g5", "h4g5", "h7h6", "f2f4", "h6g5", "f4g5", "c8c3", "a1f1", "f5d3", "f1f4", "d3g6", "f4d4", "c3g3", "d2g2", "g3e3", "d4b4", "g6e4", "g1f2", "e3f3", "f2g1", "f3e3", "g2g4", "e4d3", "g4d4", "d8e8", "a2a3", "e8e5", "g1f2", "e3h3", "f2g2", "e5e3", "b3d1", "e3g3", "g2f2", "g3g5", "d1g4", "h3h2", "f2g3", "h2d2", "g3f4", "g5e5", "b4b5", "g7g5", "f4g3", "f7f5", "b5b6", "f5g4", "b6b8", "g8g7", "d4d3", "d2d3", "g3g4", "d3a3", "b8b6", "e5d5", "b6b8", "g7g6", "b8g8", "g6f7", "g4h5", "f7g8", "h5g6", "g5g4", "g6f6", "a3e3", "f6g6", "e3e6"}
	fen09 := "6k1/8/3pr1K1/3r4/6p1/8/8/8 w - - 4 56"
	test09 := IncrementalTestSequence{moves09, fen09}
	incrementalTestSequences = append(incrementalTestSequences, test09)

	// test 4
	moves10 := []string{"e2e4", "c7c5", "b1c3", "b8c6", "g2g3", "e7e5", "g1f3", "g8f6", "f1b5", "d7d6", "b5c6", "b7c6", "e1g1", "f8e7", "d2d3", "c8b7", "f1e1", "e8g8", "a2a4", "a7a5", "b2b3", "d8d7", "c1b2", "a8d8", "c3b1", "d7e6", "b1d2", "f8e8", "d2c4", "d8a8", "b2c3", "e7d8", "d1d2", "b7a6", "c4a5", "d8a5", "c3a5", "d6d5", "a5c7", "f6g4", "e4d5", "e6f6", "d2e2", "c6d5", "f3e5", "f6e7", "e2g4", "e7c7", "g4f5", "a6c8", "f5f4", "f7f6", "e5f3", "e8e1", "a1e1", "c7f4", "g3f4", "c8f5", "g1g2", "g8f7", "h2h3", "d5d4", "e1a1", "a8g8", "a4a5", "g7g5", "f4g5", "f6g5", "a5a6", "g5g4", "f3e5", "f7e6", "e5g4", "f5g4", "h3g4", "g8g4", "g2f3", "g4g8", "a6a7", "g8a8", "b3b4", "c5b4", "f3e4", "e6d6", "e4d4", "d6c6", "f2f3", "c6b7", "a1b1", "a8a7", "c2c3", "b4b3", "b1b3", "b7c6", "f3f4", "h7h5", "f4f5", "h5h4", "f5f6", "a7d7", "d4e4", "d7d6", "b3b8", "d6f6", "b8h8", "f6e6", "e4d4", "e6d6", "d4c4", "c6b6", "h8h4", "d6c6", "c4b4", "c6c8", "c3c4", "b6c7", "c4c5", "c7d7", "b4c4", "d7e6", "d3d4", "e6f5", "h4h1", "c8a8", "c5c6", "a8a7", "d4d5", "a7c7", "c4c5", "c7a7", "d5d6", "a7a5", "c5b6", "a5a8", "d6d7", "a8b8", "b6c5", "b8b2", "d7d8q", "b2c2", "c5b4", "c2b2", "b4c3", "b2f2", "d8f8", "f5e6", "f8f2", "e6e7", "f2e3", "e7d6", "h1h6", "d6c7", "e3a7", "c7c8", "h6h8"}
	fen10 := "2k4R/Q7/2P5/8/8/2K5/8/8 b - - 8 78."
	test10 := IncrementalTestSequence{moves10, fen10}
	incrementalTestSequences = append(incrementalTestSequences, test10)

}

func printIncrementalTestResults() {

	fmt.Println("-------------------------------- Eval Test Results ----------------------------------")

	for i, testSequence := range incrementalTestSequences {

		// create new position
		pos1 := Position{}
		pos2 := Position{}

		pos1.initPositionFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		pos2.initPositionFromFen(testSequence.fenString)

		// play the moves in the first position
		for _, moveStr := range testSequence.fenMoves {
			pos1.makeUCIMove(moveStr)
		}

		// compare the position evaluations
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

		evalSuccess := pos1EvalTotal == pos2EvalTotal

		// compare the position hashes
		pos1Hash := pos1.hashOfPos
		pos2Hash := pos2.hashOfPos

		hashSuccess := pos1Hash == pos2Hash

		// print the test results
		var evalMessage string
		var hashMessage string

		if evalSuccess {
			evalMessage = "EVAL SUCCESS"
		} else {
			evalMessage = "EVAL FAILURE"
		}

		if hashSuccess {
			hashMessage = "HASH SUCCESS"
		} else {
			hashMessage = "HASH FAILURE"
		}

		fmt.Printf("[TEST %v] [%v] [%v] ||| Pos 1 hash: %v. Pos 2 hash: %v. ||| Pos 1 eval: %v (material: %v, heatmap: %v, other: %v). Pos 2 eval: %v (material: %v, heatmap: %v, other: %v).\n",
			i+1, evalMessage, hashMessage, pos1Hash, pos2Hash, pos1EvalTotal, pos1MaterialEval, pos1HeatmapEval, pos1OtherEval, pos2EvalTotal, pos2MaterialEval, pos2HeatmapEval, pos2OtherEval)
	}
}
