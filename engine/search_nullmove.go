package main

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Null Move -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// makes a null move in the position
func (pos *Position) makeNullMove() {

	// first store the game state for undo later
	pos.previousGameStates[pos.previousGameStatesCounter].pieces = pos.pieces
	pos.previousGameStates[pos.previousGameStatesCounter].piecesAll = pos.piecesAll
	pos.previousGameStates[pos.previousGameStatesCounter].castlingRights = pos.castlingRights
	pos.previousGameStates[pos.previousGameStatesCounter].enPassantTargetBB = pos.enPassantTargetBB
	pos.previousGameStates[pos.previousGameStatesCounter].halfMoves = pos.halfMoves
	pos.previousGameStates[pos.previousGameStatesCounter].kingChecks = pos.kingChecks
	pos.previousGameStates[pos.previousGameStatesCounter].evalMaterial = pos.evalMaterial
	pos.previousGameStates[pos.previousGameStatesCounter].evalHeatmaps = pos.evalHeatmaps
	pos.previousGameStates[pos.previousGameStatesCounter].evalOther = pos.evalOther
	pos.previousGameStates[pos.previousGameStatesCounter].evalMidVsEndStage = pos.evalMidVsEndStage

	pos.previousGameStatesCounter += 1

	// also store the hash for undo later
	pos.previousHashes[pos.previousHashesCounter] = pos.hashOfPos
	pos.previousHashesCounter += 1

	// the pieces stays the same, no changes needed
	// the incremental eval values (material, heatmaps) stays the same, no changes needed
	// the hash stays the same (except the side to play and en-passant square below), because pieces and castling rights are the same

	// ^^^^^^^^^ HASH ^^^^^^^^^ we remove en-passant rights if there were before the null move
	enPBBBefore := pos.enPassantTargetBB
	pos.enPassantTargetBB = emptyBB
	if enPBBBefore != pos.enPassantTargetBB {

		// hash out the old target (if there is an en-passant target and it is not blank)
		if enPBBBefore != 0 {
			enPBBBeforeSq := enPBBBefore.popBitGetSq()
			pos.hashOfPos ^= hashTableEnPassant[enPBBBeforeSq]
		}
	}

	// increase the ply
	pos.ply += 1

	// if black moves, increase the full move counter
	if !pos.isWhiteTurn {
		pos.fullMoves += 1
	}

	// reset the king check counter until the next move generation call
	pos.kingChecks = 0

	// finally, switch the side to play
	pos.isWhiteTurn = !pos.isWhiteTurn

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash the side to move
	pos.hashOfPos ^= hashTableSideToMove[0]

	// also, reset the move counter because no moves have been generated for the new position yet
	//pos.availableMovesCounter = 0
	pos.totalMovesCounter = 0
	pos.threatMovesCounter = 0
	pos.quietMovesCounter = 0

}
