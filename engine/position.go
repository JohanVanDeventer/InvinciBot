package main

import (
	"time"
)

/*

8 | . . . . . . . .
7 | . . . . . . . .
6 | . . . . . . . .
5 | . . . . . . . .
4 | . . . . . . . .
3 | . . . . . . . .
2 | . . . . . . . .
1 | . . . . . . . .
----------------
a b c d e f g h

------- Coordinates ----------
56, 57, 58, 59, 60, 61, 62, 63 (63 is least significant (right most) bit)
48, 49, 50, 51, 52, 53, 54, 55
40, 41, 42, 43, 44, 45, 46, 47
32, 33, 34, 35, 36, 37, 38, 39
24, 25, 26, 27, 28, 29, 30, 31
16, 17, 18, 19, 20, 21, 22, 23
08, 09, 10, 11, 12, 13, 14, 15
00, 01, 02, 03, 04, 05, 06, 07 (0 is the most significant (left most) bit)

position 00: 1000000000000000000000000000000000000000000000000000000000000000
position 63: 0000000000000000000000000000000000000000000000000000000000000001

*/

// constants to make it easier later to lookup position.pieces
const (
	PIECE_KING   int = 0
	PIECE_QUEEN  int = 1
	PIECE_ROOK   int = 2
	PIECE_KNIGHT int = 3
	PIECE_BISHOP int = 4
	PIECE_PAWN   int = 5

	SIDE_WHITE int = 0
	SIDE_BLACK int = 1
	SIDE_BOTH  int = 2
)

type Position struct {

	// variables from the starting Fen for the position
	pieces            [2][6]Bitboard // bitboards: white and black, and K Q R N B P respectively
	isWhiteTurn       bool
	castlingRights    [4]bool  // KQkq ordering
	enPassantTargetBB Bitboard // bitboard where the target square is set for any en-passant captures, otherwise is an empty bitboard
	halfMoves         int      // for 50 move rule: reset to zero after a capture or pawn move
	fullMoves         int      // starts at 1 and increases after black's move

	// additional piece bitboards
	piecesAll [3]Bitboard // all white is 0, all black is 1, all pieces are 2

	// game state info
	ply int // increases by 1 each time white or black moves

	// available moves in the current position
	availableMoves        [256]Move // available moves in the current positiion (max theoretically is 218)
	availableMovesCounter int       // counter points to the number of moves added

	// previous game states
	previousGameStates        [768]PreviousState
	previousGameStatesCounter int

	// hash of position
	hashOfPos             Bitboard
	previousHashes        [768]Bitboard
	previousHashesCounter int

	// game state
	gameState  int
	kingChecks int

	// evaluation of position: split into separate variables to make debugging easier
	evalMaterial      int // pure material count
	evalHeatmaps      int // heatmap count
	evalOther         int // other evaluation metrics (doubled pawns, bishop pair, king to king distance etc.)
	evalMidVsEndStage int // piece value count used for tapered heatmap eval

	// best move search variables
	bestMoveSoFar Move // used to store the best move in the search
	bestMove      Move // store the best move from the search after each iteration

	// search time management variables
	timeNodesCount       int       // increases by 1 at each node, to check time at a certain amount of nodes
	timeStartingTime     time.Time // starts when a search is initiated
	timeTotalAllowedTime int       // in milliseconds, what is the total allowed time for the search

	// log details about time taken, counts etc.
	logSearch LogSearch
	logOther  LogOther
}

// -------------------------------------------------- Step 1: Load the Fen String -----------------------------------------------
// load only the fen string into the position
func (pos *Position) step1InitFen(fen string) {

	// reset the position variables before loading new

	// load the fen string into the position
	pos.loadFenIntoPosition(fen)

}

// ---------------------------------------- Step 2: Initialize the rest of the position -----------------------------------------
// init the rest of the position ready to start
func (pos *Position) step2InitRest() {

	// add the initialized logging manager
	pos.logOther = getLoggingManager()

	// hash the loaded starting position
	pos.hashPosAndStore()

	// store the position starting eval
	pos.evalPosAtStart()
	pos.evalPosAfter()

}

// ----------------------------------------------------------- Reset Position ---------------------------------------------------
// sets the position values to the default to allow a new position to be loaded after a previous position was already loaded
// only some variables need to be reset (for example eval, game state etc are set in other functions already)
func (pos *Position) reset() {

	// reset pieces
	for side := 0; side < 2; side++ {
		for piece := 0; piece < 6; piece++ {
			pos.pieces[side][piece] = emptyBB
		}
	}

	pos.piecesAll[SIDE_WHITE] = emptyBB
	pos.piecesAll[SIDE_BLACK] = emptyBB
	pos.piecesAll[SIDE_BOTH] = emptyBB

	// reset the other fen variables
	pos.isWhiteTurn = false

	pos.castlingRights[CASTLE_WHITE_KINGSIDE] = false
	pos.castlingRights[CASTLE_WHITE_QUEENSIDE] = false
	pos.castlingRights[CASTLE_BLACK_KINGSIDE] = false
	pos.castlingRights[CASTLE_BLACK_QUEENSIDE] = false

	pos.enPassantTargetBB = emptyBB

	pos.halfMoves = 0
	pos.fullMoves = 0

	// reset the other variables
	pos.ply = 0

	pos.availableMovesCounter = 0
	pos.previousGameStatesCounter = 0
	pos.previousHashesCounter = 0

	pos.logSearch = LogSearch{}
	pos.logOther = LogOther{}

}
