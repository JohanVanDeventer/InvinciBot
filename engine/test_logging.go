package main

import (
	"fmt"
)

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Logged Type ---------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// logged type struct to log times taken for individual functions

type LoggedType struct {
	name  string
	count int
	total int
}

// log a new entry and the time it took
func (log *LoggedType) addTime(time int) {
	log.count += 1
	log.total += time
}

// print the logged details to the terminal
func (log *LoggedType) printLoggedDetails() {
	average := 0
	if log.count > 0 {
		average = log.total / log.count
	}
	fmt.Printf("%v. Took %v ns for %v calls. Avg of %v per call.\n", log.name, log.total, log.count, average)
}

// get the average nanoseconds per call
func (log *LoggedType) getAverageNsPerCall() int {
	average := 0
	if log.count > 0 {
		average = log.total / log.count
	}
	return average
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Logging Manager -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (

	// move generation
	LOG_MOVES_KING   int = 0
	LOG_MOVES_QUEEN  int = 1
	LOG_MOVES_ROOK   int = 2
	LOG_MOVES_KNIGHT int = 3
	LOG_MOVES_BISHOP int = 4
	LOG_MOVES_PAWN   int = 5

	LOG_MOVES_KING_ATTACKS int = 6
	LOG_MOVES_PINS         int = 7
	LOG_MOVES_EN_PASSANT   int = 8
	LOG_MOVES_CASTLING     int = 9

	// make move
	LOG_MAKE_MOVE int = 10

	// undo move
	LOG_UNDO_MOVE int = 11

	// hashing
	LOG_HASHING int = 12

	// evaluation
	LOG_EVAL int = 13

	// game state
	LOG_GAME_STATE int = 14

	// order moves
	LOG_ORDER_MOVES int = 15

	// search logging
	LOG_TT_GET   int = 16
	LOG_TT_STORE int = 17

	LOG_ITER_DEEP_MOVE_FIRST int = 18
)

// logging manager to manage all the LoggedType structs for a position
type LogOther struct {
	allLogTypes [19]LoggedType
}

// get a new logging manager when a position is initialized
func getLoggingManager() LogOther {
	var newLoggingManager LogOther

	// move generation
	newLoggingManager.allLogTypes[LOG_MOVES_KING] = LoggedType{"King Moves       ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_QUEEN] = LoggedType{"Queen Moves      ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_ROOK] = LoggedType{"Rook Moves       ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_KNIGHT] = LoggedType{"Knight Moves     ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_BISHOP] = LoggedType{"Bishop Moves     ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_PAWN] = LoggedType{"Pawn Moves       ", 0, 0}

	newLoggingManager.allLogTypes[LOG_MOVES_KING_ATTACKS] = LoggedType{"King Attacks     ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_PINS] = LoggedType{"Pins             ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_EN_PASSANT] = LoggedType{"En Passant       ", 0, 0}
	newLoggingManager.allLogTypes[LOG_MOVES_CASTLING] = LoggedType{"Castling         ", 0, 0}

	// make move
	newLoggingManager.allLogTypes[LOG_MAKE_MOVE] = LoggedType{"Make Move        ", 0, 0}

	// undo move
	newLoggingManager.allLogTypes[LOG_UNDO_MOVE] = LoggedType{"Undo Move        ", 0, 0}

	// hashing
	newLoggingManager.allLogTypes[LOG_HASHING] = LoggedType{"Hash Position    ", 0, 0}

	// evaluation
	newLoggingManager.allLogTypes[LOG_EVAL] = LoggedType{"Eval Position    ", 0, 0}

	// game state check
	newLoggingManager.allLogTypes[LOG_GAME_STATE] = LoggedType{"Get Game State   ", 0, 0}

	// order moves
	newLoggingManager.allLogTypes[LOG_ORDER_MOVES] = LoggedType{"Order Moves      ", 0, 0}

	// search logging
	newLoggingManager.allLogTypes[LOG_TT_GET] = LoggedType{"Get TT Entry     ", 0, 0}
	newLoggingManager.allLogTypes[LOG_TT_STORE] = LoggedType{"Store TT Entry   ", 0, 0}
	newLoggingManager.allLogTypes[LOG_ITER_DEEP_MOVE_FIRST] = LoggedType{"IterDeep Ordering", 0, 0}

	return newLoggingManager
}

func (lm *LogOther) printLoggedDetails() {
	for _, logType := range lm.allLogTypes {
		logType.printLoggedDetails()
	}
}
