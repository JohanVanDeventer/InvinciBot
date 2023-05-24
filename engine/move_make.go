package main

import (
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Make Move -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// function to make a move on a position
// note: code marked with ^^^ HASH ^^^ and ^^^ EVAL ^^^ are added in for incremental updates, and do not directly relate to making moves

func (pos *Position) makeMove(move Move) {

	start_time := time.Now()

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

	// set sides
	var frSide int
	var enSide int

	if pos.isWhiteTurn {
		frSide = SIDE_WHITE
		enSide = SIDE_BLACK
	} else {
		frSide = SIDE_BLACK
		enSide = SIDE_WHITE
	}

	// get the enemy piece type in case of a capture (remember, cannot capture king)
	var enemyPiece int = 6 // set outside range to catch bugs

	if move.moveType == MOVE_TYPE_CAPTURE { // try order by most likely piece first (most numerous opponents)
		if pos.pieces[enSide][PIECE_PAWN].isBitSet(move.toSq) {
			enemyPiece = PIECE_PAWN
		} else if pos.pieces[enSide][PIECE_KNIGHT].isBitSet(move.toSq) {
			enemyPiece = PIECE_KNIGHT
		} else if pos.pieces[enSide][PIECE_BISHOP].isBitSet(move.toSq) {
			enemyPiece = PIECE_BISHOP
		} else if pos.pieces[enSide][PIECE_ROOK].isBitSet(move.toSq) {
			enemyPiece = PIECE_ROOK
		} else if pos.pieces[enSide][PIECE_QUEEN].isBitSet(move.toSq) {
			enemyPiece = PIECE_QUEEN
		}
	}

	// remove the piece on the "from" square from all friendly bitboards
	pos.piecesAll[SIDE_BOTH].clearBit(move.fromSq)
	pos.piecesAll[frSide].clearBit(move.fromSq)
	pos.pieces[frSide][move.piece].clearBit(move.fromSq)

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash the "from" friendly piece out
	pos.hashOfPos ^= hashTablePieces[move.fromSq][frSide][move.piece]

	// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ no eval yet, no piece is taken

	// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the weighted heatmap value of the "from" square
	// this is done on the "before" stage, because that stage was used to add the value previously (after captures on the previous move)
	evalMidVsEndStageBefore := pos.evalMidVsEndStage
	if evalMidVsEndStageBefore > 24 { // cap to 24
		evalMidVsEndStageBefore = 24
	}
	midValueFriendlyFrom := evalTableCombinedMid[frSide][move.piece][move.fromSq]
	endValueFriendlyFrom := evalTableCombinedEnd[frSide][move.piece][move.fromSq]
	weightedValueFriendlyFrom := ((midValueFriendlyFrom * evalMidVsEndStageBefore) + (endValueFriendlyFrom * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
	pos.evalHeatmaps -= weightedValueFriendlyFrom

	// add the piece on the "to" square on all friendly bitboards
	pos.piecesAll[SIDE_BOTH].setBit(move.toSq)
	pos.piecesAll[frSide].setBit(move.toSq)
	pos.pieces[frSide][move.piece].setBit(move.toSq)

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash the "to" friendly piece in
	pos.hashOfPos ^= hashTablePieces[move.toSq][frSide][move.piece]

	// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ no eval yet, no piece is taken

	// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ wait before we add this, we need to first update the new game stage value (captures, promotions) before it's added

	// now depending on the move type, remove enemy pieces, capture en-passant, or castle
	switch move.moveType {

	case MOVE_TYPE_QUIET:
		// for quiet moves, just place the piece on the new square
		// already done above

		// ^^^^^^^^^ HASH ^^^^^^^^^ nothing extra required

		// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

		// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ we add the friendly piece new heatmap value (use the stage value before, no changes were made to the stage value)
		midValueFriendlyTo := evalTableCombinedMid[frSide][move.piece][move.toSq]
		endValueFriendlyTo := evalTableCombinedEnd[frSide][move.piece][move.toSq]
		weightedValueFriendlyTo := ((midValueFriendlyTo * evalMidVsEndStageBefore) + (endValueFriendlyTo * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
		pos.evalHeatmaps += weightedValueFriendlyTo

	case MOVE_TYPE_CAPTURE:
		// remove the enemy piece
		pos.piecesAll[enSide].clearBit(move.toSq)
		pos.pieces[enSide][enemyPiece].clearBit(move.toSq)

		// ^^^^^^^^^ HASH ^^^^^^^^^ hash the "to" enemy piece out
		pos.hashOfPos ^= hashTablePieces[move.toSq][enSide][enemyPiece]

		// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ remove the captured piece from the material eval and game stage eval
		pos.evalMaterial -= evalTableMaterial[enSide][enemyPiece]
		pos.evalMidVsEndStage -= evalTableGameStage[enemyPiece]

		// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the friendly piece, and remove the enemy piece from the heatmap value
		// we add the friendly piece using the stage after captures
		// but we remove the enemy piece using the stage before captures (that was used to record it's value initially)
		midValueEnemyTo := evalTableCombinedMid[enSide][enemyPiece][move.toSq]
		endValueEnemyTo := evalTableCombinedEnd[enSide][enemyPiece][move.toSq]
		weightedValueEnemyTo := ((midValueEnemyTo * evalMidVsEndStageBefore) + (endValueEnemyTo * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
		pos.evalHeatmaps -= weightedValueEnemyTo

		evalMidVsEndStageAfter := pos.evalMidVsEndStage
		if evalMidVsEndStageAfter > 24 {
			evalMidVsEndStageAfter = 24
		}
		midValueFriendlyTo := evalTableCombinedMid[frSide][move.piece][move.toSq]
		endValueFriendlyTo := evalTableCombinedEnd[frSide][move.piece][move.toSq]
		weightedValueFriendlyTo := ((midValueFriendlyTo * evalMidVsEndStageAfter) + (endValueFriendlyTo * (STAGE_VAL_STARTING - evalMidVsEndStageAfter))) / STAGE_VAL_STARTING
		pos.evalHeatmaps += weightedValueFriendlyTo

	case MOVE_TYPE_EN_PASSANT:
		// remove the en-passant captured pawn
		if pos.isWhiteTurn {
			pos.piecesAll[SIDE_BOTH].clearBit(move.toSq - 8)
			pos.piecesAll[enSide].clearBit(move.toSq - 8)
			pos.pieces[enSide][PIECE_PAWN].clearBit(move.toSq - 8)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the "en-passant" enemy piece out
			pos.hashOfPos ^= hashTablePieces[move.toSq-8][enSide][PIECE_PAWN]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ remove the captured piece from the material eval and game stage eval
			pos.evalMaterial -= evalTableMaterial[enSide][PIECE_PAWN]
			pos.evalMidVsEndStage -= evalTableGameStage[PIECE_PAWN]

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the friendly piece, and remove the enemy piece from the heatmap value
			// we add the friendly piece using the stage after captures
			// but we remove the enemy piece using the stage before captures (that was used to record it's value initially)
			midValueEnemyTo := evalTableCombinedMid[enSide][PIECE_PAWN][move.toSq-8]
			endValueEnemyTo := evalTableCombinedEnd[enSide][PIECE_PAWN][move.toSq-8]
			weightedValueEnemyTo := ((midValueEnemyTo * evalMidVsEndStageBefore) + (endValueEnemyTo * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueEnemyTo

			evalMidVsEndStageAfter := pos.evalMidVsEndStage
			if evalMidVsEndStageAfter > 24 {
				evalMidVsEndStageAfter = 24
			}
			midValueFriendlyTo := evalTableCombinedMid[frSide][PIECE_PAWN][move.toSq]
			endValueFriendlyTo := evalTableCombinedEnd[frSide][PIECE_PAWN][move.toSq]
			weightedValueFriendlyTo := ((midValueFriendlyTo * evalMidVsEndStageAfter) + (endValueFriendlyTo * (STAGE_VAL_STARTING - evalMidVsEndStageAfter))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyTo

		} else {
			pos.piecesAll[SIDE_BOTH].clearBit(move.toSq + 8)
			pos.piecesAll[enSide].clearBit(move.toSq + 8)
			pos.pieces[enSide][PIECE_PAWN].clearBit(move.toSq + 8)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the "en-passant" enemy piece out
			pos.hashOfPos ^= hashTablePieces[move.toSq+8][enSide][PIECE_PAWN]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ remove the captured piece from the material eval and game stage eval
			pos.evalMaterial -= evalTableMaterial[enSide][PIECE_PAWN]
			pos.evalMidVsEndStage -= evalTableGameStage[PIECE_PAWN]

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the friendly piece, and remove the enemy piece from the heatmap value
			// we add the friendly piece using the stage after captures
			// but we remove the enemy piece using the stage before captures (that was used to record it's value initially)
			midValueEnemyTo := evalTableCombinedMid[enSide][PIECE_PAWN][move.toSq+8]
			endValueEnemyTo := evalTableCombinedEnd[enSide][PIECE_PAWN][move.toSq+8]
			weightedValueEnemyTo := ((midValueEnemyTo * evalMidVsEndStageBefore) + (endValueEnemyTo * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueEnemyTo

			evalMidVsEndStageAfter := pos.evalMidVsEndStage
			if evalMidVsEndStageAfter > 24 {
				evalMidVsEndStageAfter = 24
			}
			midValueFriendlyTo := evalTableCombinedMid[frSide][PIECE_PAWN][move.toSq]
			endValueFriendlyTo := evalTableCombinedEnd[frSide][PIECE_PAWN][move.toSq]
			weightedValueFriendlyTo := ((midValueFriendlyTo * evalMidVsEndStageAfter) + (endValueFriendlyTo * (STAGE_VAL_STARTING - evalMidVsEndStageAfter))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyTo
		}

	case MOVE_TYPE_CASTLE:
		if move.toSq == 6 {
			// remove the rook from the original square
			pos.piecesAll[SIDE_BOTH].clearBit(7)
			pos.piecesAll[frSide].clearBit(7)
			pos.pieces[frSide][PIECE_ROOK].clearBit(7)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook out
			pos.hashOfPos ^= hashTablePieces[7][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the value of the removed rook (using value before, no captures were made)
			midValueFriendlyToRemove := evalTableCombinedMid[frSide][PIECE_ROOK][7]
			endValueFriendlyToRemove := evalTableCombinedEnd[frSide][PIECE_ROOK][7]
			weightedValueFriendlyToRemove := ((midValueFriendlyToRemove * evalMidVsEndStageBefore) + (endValueFriendlyToRemove * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueFriendlyToRemove

			// and add to the new square
			pos.piecesAll[SIDE_BOTH].setBit(5)
			pos.piecesAll[frSide].setBit(5)
			pos.pieces[frSide][PIECE_ROOK].setBit(5)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook in
			pos.hashOfPos ^= hashTablePieces[5][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the value of the moved rook (using value before, no captures were made)
			// also add the king value on the "to" square (not done yet)
			midValueFriendlyToAdd := evalTableCombinedMid[frSide][PIECE_ROOK][5]
			endValueFriendlyToAdd := evalTableCombinedEnd[frSide][PIECE_ROOK][5]
			weightedValueFriendlyToAdd := ((midValueFriendlyToAdd * evalMidVsEndStageBefore) + (endValueFriendlyToAdd * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyToAdd
			midValueFriendlyKing := evalTableCombinedMid[frSide][PIECE_KING][move.toSq]
			endValueFriendlyKing := evalTableCombinedEnd[frSide][PIECE_KING][move.toSq]
			weightedValueFriendlyKing := ((midValueFriendlyKing * evalMidVsEndStageBefore) + (endValueFriendlyKing * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyKing
		}

		if move.toSq == 2 {
			// remove the rook from the original square
			pos.piecesAll[SIDE_BOTH].clearBit(0)
			pos.piecesAll[frSide].clearBit(0)
			pos.pieces[frSide][PIECE_ROOK].clearBit(0)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook out
			pos.hashOfPos ^= hashTablePieces[0][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the value of the removed rook (using value before, no captures were made)
			midValueFriendlyToRemove := evalTableCombinedMid[frSide][PIECE_ROOK][0]
			endValueFriendlyToRemove := evalTableCombinedEnd[frSide][PIECE_ROOK][0]
			weightedValueFriendlyToRemove := ((midValueFriendlyToRemove * evalMidVsEndStageBefore) + (endValueFriendlyToRemove * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueFriendlyToRemove

			// and add to the new square
			pos.piecesAll[SIDE_BOTH].setBit(3)
			pos.piecesAll[frSide].setBit(3)
			pos.pieces[frSide][PIECE_ROOK].setBit(3)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook in
			pos.hashOfPos ^= hashTablePieces[3][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the value of the moved rook (using value before, no captures were made)
			// also add the king value on the "to" square (not done yet)
			midValueFriendlyToAdd := evalTableCombinedMid[frSide][PIECE_ROOK][3]
			endValueFriendlyToAdd := evalTableCombinedEnd[frSide][PIECE_ROOK][3]
			weightedValueFriendlyToAdd := ((midValueFriendlyToAdd * evalMidVsEndStageBefore) + (endValueFriendlyToAdd * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyToAdd
			midValueFriendlyKing := evalTableCombinedMid[frSide][PIECE_KING][move.toSq]
			endValueFriendlyKing := evalTableCombinedEnd[frSide][PIECE_KING][move.toSq]
			weightedValueFriendlyKing := ((midValueFriendlyKing * evalMidVsEndStageBefore) + (endValueFriendlyKing * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyKing
		}

		if move.toSq == 62 {
			// remove the rook from the original square
			pos.piecesAll[SIDE_BOTH].clearBit(63)
			pos.piecesAll[frSide].clearBit(63)
			pos.pieces[frSide][PIECE_ROOK].clearBit(63)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook out
			pos.hashOfPos ^= hashTablePieces[63][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the value of the removed rook (using value before, no captures were made)
			midValueFriendlyToRemove := evalTableCombinedMid[frSide][PIECE_ROOK][63]
			endValueFriendlyToRemove := evalTableCombinedEnd[frSide][PIECE_ROOK][63]
			weightedValueFriendlyToRemove := ((midValueFriendlyToRemove * evalMidVsEndStageBefore) + (endValueFriendlyToRemove * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueFriendlyToRemove

			// and add to the new square
			pos.piecesAll[SIDE_BOTH].setBit(61)
			pos.piecesAll[frSide].setBit(61)
			pos.pieces[frSide][PIECE_ROOK].setBit(61)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook in
			pos.hashOfPos ^= hashTablePieces[61][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the value of the moved rook (using value before, no captures were made)
			// also add the king value on the "to" square (not done yet)
			midValueFriendlyToAdd := evalTableCombinedMid[frSide][PIECE_ROOK][61]
			endValueFriendlyToAdd := evalTableCombinedEnd[frSide][PIECE_ROOK][61]
			weightedValueFriendlyToAdd := ((midValueFriendlyToAdd * evalMidVsEndStageBefore) + (endValueFriendlyToAdd * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyToAdd
			midValueFriendlyKing := evalTableCombinedMid[frSide][PIECE_KING][move.toSq]
			endValueFriendlyKing := evalTableCombinedEnd[frSide][PIECE_KING][move.toSq]
			weightedValueFriendlyKing := ((midValueFriendlyKing * evalMidVsEndStageBefore) + (endValueFriendlyKing * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyKing
		}

		if move.toSq == 58 {
			// remove the rook from the original square
			pos.piecesAll[SIDE_BOTH].clearBit(56)
			pos.piecesAll[frSide].clearBit(56)
			pos.pieces[frSide][PIECE_ROOK].clearBit(56)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook out
			pos.hashOfPos ^= hashTablePieces[56][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the value of the removed rook (using value before, no captures were made)
			midValueFriendlyToRemove := evalTableCombinedMid[frSide][PIECE_ROOK][56]
			endValueFriendlyToRemove := evalTableCombinedEnd[frSide][PIECE_ROOK][56]
			weightedValueFriendlyToRemove := ((midValueFriendlyToRemove * evalMidVsEndStageBefore) + (endValueFriendlyToRemove * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps -= weightedValueFriendlyToRemove

			// and add to the new square
			pos.piecesAll[SIDE_BOTH].setBit(59)
			pos.piecesAll[frSide].setBit(59)
			pos.pieces[frSide][PIECE_ROOK].setBit(59)

			// ^^^^^^^^^ HASH ^^^^^^^^^ hash the rook in
			pos.hashOfPos ^= hashTablePieces[59][frSide][PIECE_ROOK]

			// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ nothing extra required

			// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ add the value of the moved rook (using value before, no captures were made)
			// also add the king value on the "to" square (not done yet)
			midValueFriendlyToAdd := evalTableCombinedMid[frSide][PIECE_ROOK][56]
			endValueFriendlyToAdd := evalTableCombinedEnd[frSide][PIECE_ROOK][56]
			weightedValueFriendlyToAdd := ((midValueFriendlyToAdd * evalMidVsEndStageBefore) + (endValueFriendlyToAdd * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyToAdd
			midValueFriendlyKing := evalTableCombinedMid[frSide][PIECE_KING][move.toSq]
			endValueFriendlyKing := evalTableCombinedEnd[frSide][PIECE_KING][move.toSq]
			weightedValueFriendlyKing := ((midValueFriendlyKing * evalMidVsEndStageBefore) + (endValueFriendlyKing * (STAGE_VAL_STARTING - evalMidVsEndStageBefore))) / STAGE_VAL_STARTING
			pos.evalHeatmaps += weightedValueFriendlyKing
		}
	}

	// handle promotions if there are any
	if move.promotionType != PROMOTION_NONE {

		// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ save the "after" game stage used for pawn captures and quiet moves but before promotions (used below)
		evalMidVsEndStageAfterCapturesBeforePromotions := pos.evalMidVsEndStage

		// remove the friendly pawn on that square
		pos.pieces[frSide][PIECE_PAWN].clearBit(move.toSq)

		// ^^^^^^^^^ HASH ^^^^^^^^^ remove the friendly pawn
		pos.hashOfPos ^= hashTablePieces[move.toSq][frSide][PIECE_PAWN]

		// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ remove the pawn from the eval
		pos.evalMaterial -= evalTableMaterial[frSide][PIECE_PAWN]
		pos.evalMidVsEndStage -= evalTableGameStage[PIECE_PAWN]

		// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ do both below after the whole game stage has been updated (remove pawn and add promoted piece)

		// add the promoted piece to the relevant bitboard
		pos.pieces[frSide][move.promotionType].setBit(move.toSq)

		// ^^^^^^^^^ HASH ^^^^^^^^^ add the promoted piece
		pos.hashOfPos ^= hashTablePieces[move.toSq][frSide][move.promotionType]

		// ^^^^^^^^^ EVAL: MATERIAL AND GAME STAGE ^^^^^^^^^ add the promoted piece to the eval
		pos.evalMaterial += evalTableMaterial[frSide][move.promotionType]
		pos.evalMidVsEndStage += evalTableGameStage[move.promotionType]

		// ^^^^^^^^^ EVAL: HEATMAPS ^^^^^^^^^ remove the friendly pawn and add the promoted piece
		// remove the friendly pawn using the stage value after captures and other moves but before updating for promotions
		// add the promoted piece using the stage value after updating for promotions
		midValueFriendlyPawn := evalTableCombinedMid[frSide][PIECE_PAWN][move.toSq]
		endValueFriendlyPawn := evalTableCombinedEnd[frSide][PIECE_PAWN][move.toSq]
		weightedValueFriendlyPawn := ((midValueFriendlyPawn * evalMidVsEndStageAfterCapturesBeforePromotions) + (endValueFriendlyPawn * (STAGE_VAL_STARTING - evalMidVsEndStageAfterCapturesBeforePromotions))) / STAGE_VAL_STARTING
		pos.evalHeatmaps -= weightedValueFriendlyPawn

		midValueFriendlyPiece := evalTableCombinedMid[frSide][move.promotionType][move.toSq]
		endValueFriendlyPiece := evalTableCombinedEnd[frSide][move.promotionType][move.toSq]
		weightedValueFriendlyPiece := ((midValueFriendlyPiece * pos.evalMidVsEndStage) + (endValueFriendlyPiece * (STAGE_VAL_STARTING - pos.evalMidVsEndStage))) / STAGE_VAL_STARTING
		pos.evalHeatmaps += weightedValueFriendlyPiece
	}

	// ^^^^^^^^^ HASH ^^^^^^^^^ store the castling rights before changes
	castlingRightsBefore := pos.castlingRights

	// if the king moves (castle or otherwise), or a rook moves or is captured, remove castling rights
	if move.fromSq == 4 { // if the king moves, cancel both castling rights
		pos.castlingRights[CASTLE_WHITE_KINGSIDE] = false
		pos.castlingRights[CASTLE_WHITE_QUEENSIDE] = false
	}
	if move.fromSq == 7 || move.toSq == 7 { // else, cancel the rook moves on that side only
		pos.castlingRights[CASTLE_WHITE_KINGSIDE] = false
	}
	if move.fromSq == 0 || move.toSq == 0 { // else, cancel the rook moves on that side only
		pos.castlingRights[CASTLE_WHITE_QUEENSIDE] = false
	}

	if move.fromSq == 60 { // if the king moves, cancel both castling rights
		pos.castlingRights[CASTLE_BLACK_KINGSIDE] = false
		pos.castlingRights[CASTLE_BLACK_QUEENSIDE] = false
	}
	if move.fromSq == 63 || move.toSq == 63 { // else, cancel the rook moves on that side only
		pos.castlingRights[CASTLE_BLACK_KINGSIDE] = false
	}
	if move.fromSq == 56 || move.toSq == 56 { // else, cancel the rook moves on that side only
		pos.castlingRights[CASTLE_BLACK_QUEENSIDE] = false
	}

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash changes in castling rights
	for i := 0; i < 4; i++ {
		if castlingRightsBefore[i] != pos.castlingRights[i] { // meaning changes
			pos.hashOfPos ^= hashTableCastling[i]
		}
	}

	// ^^^^^^^^^ HASH ^^^^^^^^^ store the en-passant BB before changes
	enPBBBefore := pos.enPassantTargetBB

	// if the piece is a pawn that moved 2 squares, set en-passant, otherwise clear it
	pos.enPassantTargetBB = emptyBB
	if (move.toSq-move.fromSq) == 16 && move.piece == PIECE_PAWN {
		pos.enPassantTargetBB = bbReferenceArray[move.toSq-8]
	}
	if (move.toSq-move.fromSq) == -16 && move.piece == PIECE_PAWN {
		pos.enPassantTargetBB = bbReferenceArray[move.toSq+8]
	}

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash changes in en-passant targets if there are changes
	if enPBBBefore != pos.enPassantTargetBB {

		// hash out the old target (if there is an en-passant target and it is not blank)
		if enPBBBefore != 0 {
			enPBBBeforeSq := enPBBBefore.popBitGetSq()
			pos.hashOfPos ^= hashTableEnPassant[enPBBBeforeSq]
		}

		// hash in the new target (if there is an en-passant target and it is not blank)
		if pos.enPassantTargetBB != 0 {
			enPBBCopy := pos.enPassantTargetBB
			enPBBCopySq := enPBBCopy.popBitGetSq()
			pos.hashOfPos ^= hashTableEnPassant[enPBBCopySq]
		}
	}

	// increase the ply
	pos.ply += 1

	// if black moves, increase the full move counter
	if !pos.isWhiteTurn {
		pos.fullMoves += 1
	}

	// reset the half-move counter (for 50-move rule) when a pawn moves, or there is a capture/promotion
	// else increment it by 1
	if move.piece == PIECE_PAWN || move.moveType == MOVE_TYPE_CAPTURE || move.moveType == MOVE_TYPE_EN_PASSANT || move.promotionType != PROMOTION_NONE {
		pos.halfMoves = 1
	} else {
		pos.halfMoves += 1
	}

	// reset the king check counter until the next move generation call
	pos.kingChecks = 0

	// finally, switch the side to play
	pos.isWhiteTurn = !pos.isWhiteTurn

	// ^^^^^^^^^ HASH ^^^^^^^^^ hash the side to move
	pos.hashOfPos ^= hashTableSideToMove[0]

	// also, reset the move counter because no moves have been generated for the new position yet
	pos.availableMovesCounter = 0

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_MAKE_MOVE].addTime(int(duration_time))
}
