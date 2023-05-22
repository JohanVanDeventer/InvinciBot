package main

import (
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Previous Game State ------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// information to restore the previous position
// also stores the full bitboard (easier than lots of "if" statements)
// this is easier than adjusting all the info again incrementally

type PreviousState struct {
	pieces            [2][6]Bitboard
	piecesAll         [3]Bitboard
	castlingRights    [4]bool
	enPassantTargetBB Bitboard
	halfMoves         int
	kingChecks        int
	evalMaterial      int
	evalHeatmaps      int
	evalOther         int
	evalMidVsEndStage int
}

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Undo Move -------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// function to undo the last move played
func (pos *Position) undoMove() {

	start_time := time.Now()

	// if the ply is at the start, we cannot undo the previous move
	if pos.ply == 0 {
		return
	}

	// get the last game state
	pos.previousGameStatesCounter -= 1

	// restore pieces and other information
	pos.pieces = pos.previousGameStates[pos.previousGameStatesCounter].pieces
	pos.piecesAll = pos.previousGameStates[pos.previousGameStatesCounter].piecesAll
	pos.castlingRights = pos.previousGameStates[pos.previousGameStatesCounter].castlingRights
	pos.enPassantTargetBB = pos.previousGameStates[pos.previousGameStatesCounter].enPassantTargetBB
	pos.halfMoves = pos.previousGameStates[pos.previousGameStatesCounter].halfMoves
	pos.kingChecks = pos.previousGameStates[pos.previousGameStatesCounter].kingChecks
	pos.evalMaterial = pos.previousGameStates[pos.previousGameStatesCounter].evalMaterial
	pos.evalHeatmaps = pos.previousGameStates[pos.previousGameStatesCounter].evalHeatmaps
	pos.evalOther = pos.previousGameStates[pos.previousGameStatesCounter].evalOther
	pos.evalMidVsEndStage = pos.previousGameStates[pos.previousGameStatesCounter].evalMidVsEndStage

	// also restore the hash
	pos.previousHashesCounter -= 1
	pos.hashOfPos = pos.previousHashes[pos.previousHashesCounter]

	// reset the game state to ongoing (can only be ongoing when undoing a move)
	pos.gameState = STATE_ONGOING

	// if black moved last, decrease the full moves
	if pos.isWhiteTurn {
		pos.fullMoves -= 1
	}

	// decrease the ply
	pos.ply -= 1

	// finally switch the sides
	pos.isWhiteTurn = !pos.isWhiteTurn

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_UNDO_MOVE].addTime(int(duration_time))
}
