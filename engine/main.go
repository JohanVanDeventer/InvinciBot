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
		// TODO: fix small eval bug: possibly related to rounding black scores vs white scores (no big impact for now)
		// TODO: fix small mobility difference: only after a few moves played from each side will mobility be accurate (no big impact for now)
		printPerftTestResults()
		printIncrementalTestResults()
	}

	// create a new position for the engine to start with
	pos := Position{}

	// start the uci input loop
	pos.startUCIInputLoop()
}
