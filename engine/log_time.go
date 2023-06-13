package main

import (
	"fmt"
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Time Log Holder -------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// code to log the time taken by various functions

// struct to log times taken for individual functions
type TimeLogHolder struct {
	count     int
	total     int64
	startTime time.Time
}

// start the timer on a new function call
func (log *TimeLogHolder) start() {
	log.startTime = time.Now()
}

// stop the timer on a new function call, and log the time taken
func (log *TimeLogHolder) stop() {

	// increase the function call count
	log.count += 1

	// and log the time taken
	log.total += time.Since(log.startTime).Nanoseconds()
}

// get the average nanoseconds per call
func (log *TimeLogHolder) getAverageNsPerCall() int {
	average := 0
	if log.count > 0 {
		average = int(log.total) / log.count
	}
	return average
}

// print the logged details to the terminal
func (log *TimeLogHolder) printLoggedDetails() {
	fmt.Printf("Calls: %v. Average ns: %v.\n", log.count, log.getAverageNsPerCall())
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Logging Manager -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (

	// ------------------------------- MAIN FUNCTIONS -------------------------------
	// move generation
	LOG_MOVE_GEN_TOTAL int = 0

	LOG_MOVE_GEN_KING   int = 1
	LOG_MOVE_GEN_QUEEN  int = 2
	LOG_MOVE_GEN_ROOK   int = 3
	LOG_MOVE_GEN_BISHOP int = 4
	LOG_MOVE_GEN_KNIGHT int = 5
	LOG_MOVE_GEN_PAWN   int = 6

	LOG_MOVE_GEN_KING_ATTACKS int = 7
	LOG_MOVE_GEN_PINS         int = 8
	LOG_MOVE_GEN_EN_PASSANT   int = 9
	LOG_MOVE_GEN_CASTLING     int = 10

	// make and undo move
	LOG_MAKE_MOVE     int = 11
	LOG_UNDO_MOVE     int = 12
	LOG_MAKE_NULLMOVE int = 13

	// evaluation
	LOG_EVAL int = 14

	// game state
	LOG_GAME_STATE int = 15

	// -------------------------------- SEARCH FUNCTIONS --------------------------------

	// transposition table
	LOG_SEARCH_TT_PROBE int = 16
	LOG_SEARCH_TT_STORE int = 17

	// order moves
	LOG_SEARCH_COPY_THREAT_MOVES              int = 18
	LOG_SEARCH_COPY_QUIET_MOVES               int = 19
	LOG_SEARCH_ORDER_THREAT_MOVES             int = 20
	LOG_SEARCH_ORDER_KILLER_1                 int = 21
	LOG_SEARCH_ORDER_KILLER_2                 int = 22
	LOG_SEARCH_ORDER_HASH_MOVES               int = 23
	LOG_SEARCH_ORDER_PREVIOUS_ITERATION_MOVES int = 24

	// ------------------------------- ONCE-OFF FUNCTIONS -------------------------------
	// load a position from a fen string
	LOG_ONCE_LOAD_FEN int = 25

	// hashing
	LOG_ONCE_HASH int = 26

	// eval
	LOG_ONCE_EVAL int = 27

	// search startup
	LOG_ONCE_SEARCH_STARTUP int = 28
)

// logging manager to manage all the LoggedType structs for a position
type TimeLogger struct {
	allLogTypes [29]TimeLogHolder
}

// print all logged details
func (tl *TimeLogger) printLoggedDetails() {
	for id, logType := range tl.allLogTypes {
		fmt.Printf("Log ID: %v. ", id)
		logType.printLoggedDetails()
	}
}

// get a new logging manager with clean values
func getNewTimeLogger() TimeLogger {
	var newTimeLogger TimeLogger
	return newTimeLogger
}
