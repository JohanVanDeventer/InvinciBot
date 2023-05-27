package main

func main() {

	// --- tests ---
	initEngine()
	//playBestMoveGames(1000)
	//printInitTestResults()
	printPerftResults()

	// create a new position for the engine to start with
	pos := Position{}

	// start either uci loop waiting for gui input, or start a gui loop for playing in the terminal
	startAsUCI := false

	if startAsUCI {
		pos.startUCIInputLoop()
	} else {
		initEngine()
		pos.step1InitFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		pos.step2InitRest()
		pos.startGameLoopTerminalGUI()
	}
}
