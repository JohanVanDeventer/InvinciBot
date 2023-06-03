package main

// positions to test whether the computer can play the best move
var bestMovePositions []string

func initBestMovePositions() {
	bestMovePositions = append(bestMovePositions, "8/k7/3p4/p2P1p2/P2P1P2/8/8/K7 w - - 0 1") // deep search winning endgame for white
	bestMovePositions = append(bestMovePositions, "6Q1/8/8/5K1k/8/8/8/8 w - - 0 1")          // checkmate in 1
	bestMovePositions = append(bestMovePositions, "8/8/4k3/8/8/8/8/2KQ4 w - - 0 1")          // king vs king and queen
	bestMovePositions = append(bestMovePositions, "8/8/4k3/8/8/8/8/2KR4 w - - 0 1")          // king vs king and rook

	bestMovePositions = append(bestMovePositions, "5B2/6P1/1p6/8/1N6/kP6/2K5/8 w - - 0 1") // white to mate in 3 moves

	bestMovePositions = append(bestMovePositions, "7K/8/k1P5/7p/8/8/8/8 w - - 0 1")         // endgame: draw for white with deep play
	bestMovePositions = append(bestMovePositions, "8/8/7p/3KNN1k/2p4p/8/3P2p1/8 w - - 0 1") // behting study: draw for white with deep play
	bestMovePositions = append(bestMovePositions, "4k3/5ppp/8/8/8/8/PPP5/3K4 w - - 0 1")    // 3 pawns vs 3 pawns
}

func playBestMoveGames(timeLimitMs int) {
	for _, fen := range bestMovePositions {
		newPos := Position{}
		newPos.initPositionFromFen(fen)
		newPos.startGameLoopTerminalGUI()
	}
}
