package main

import (
	"fmt"
	"time"
)

func main() {

	// engine tests
	runTests := false
	if runTests {

		// init the engine
		startTime := time.Now()
		initEngine()
		durationTime := time.Since(startTime).Milliseconds()

		fmt.Printf("-----------------------------------------\n")
		fmt.Printf("Engine Initialization Time: %v ms.\n", durationTime)
		fmt.Printf("-----------------------------------------\n")

		// main tests
		printPerftTestResults()
		printIncrementalTestResults() // TODO: fix small eval bug: possibly related to rounding black scores vs white scores, but no big impact for now
	}

	// create a new position for the engine to start with
	pos := Position{}

	// start the uci input loop
	pos.startUCIInputLoop()
}
