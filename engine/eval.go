package main

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Background ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Evaluation is rounded to the nearest centipawn.
Pawns are worth 100 centipawns as reference.

The eval is always done from the white side (absolute value).
Positive is good for white, negative is good for black.

Some eval is done incrementally, some for each position from the start.
*/

const (
	// direct material value in centipawns
	VALUE_PAWN   int = 100
	VALUE_KNIGHT int = 400
	VALUE_BISHOP int = 420
	VALUE_ROOK   int = 600
	VALUE_QUEEN  int = 1200

	// stage of the game (mid vs end) value
	STAGE_VAL_QUEEN    int = 4
	STAGE_VAL_ROOK     int = 2
	STAGE_VAL_KNIGHT   int = 1
	STAGE_VAL_BISHOP   int = 1
	STAGE_VAL_STARTING int = STAGE_VAL_QUEEN*2 + STAGE_VAL_ROOK*4 + STAGE_VAL_KNIGHT*4 + STAGE_VAL_BISHOP*4 // normally 24

)

var evalTableMaterial [2][6]int // maps the side and piece type to their material values (black has negative values, kings are 0)
var evalTableGameStage [6]int   // maps the piece type to their game stage values (kings and pawns are 0)

func initEvalMaterialAndStageTables() {
	evalTableMaterial[SIDE_WHITE][PIECE_KING] = 0
	evalTableMaterial[SIDE_WHITE][PIECE_QUEEN] = VALUE_QUEEN
	evalTableMaterial[SIDE_WHITE][PIECE_ROOK] = VALUE_ROOK
	evalTableMaterial[SIDE_WHITE][PIECE_KNIGHT] = VALUE_KNIGHT
	evalTableMaterial[SIDE_WHITE][PIECE_BISHOP] = VALUE_BISHOP
	evalTableMaterial[SIDE_WHITE][PIECE_PAWN] = VALUE_PAWN

	evalTableMaterial[SIDE_BLACK][PIECE_KING] = 0
	evalTableMaterial[SIDE_BLACK][PIECE_QUEEN] = 0 - VALUE_QUEEN
	evalTableMaterial[SIDE_BLACK][PIECE_ROOK] = 0 - VALUE_ROOK
	evalTableMaterial[SIDE_BLACK][PIECE_KNIGHT] = 0 - VALUE_KNIGHT
	evalTableMaterial[SIDE_BLACK][PIECE_BISHOP] = 0 - VALUE_BISHOP
	evalTableMaterial[SIDE_BLACK][PIECE_PAWN] = 0 - VALUE_PAWN

	evalTableGameStage[PIECE_KING] = 0
	evalTableGameStage[PIECE_QUEEN] = STAGE_VAL_QUEEN
	evalTableGameStage[PIECE_ROOK] = STAGE_VAL_ROOK
	evalTableGameStage[PIECE_KNIGHT] = STAGE_VAL_KNIGHT
	evalTableGameStage[PIECE_BISHOP] = STAGE_VAL_BISHOP
	evalTableGameStage[PIECE_PAWN] = 0

}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------- Eval: Material, Game Stage, and Heatmaps -----------------------------------
// --------------------------------------------------------------------------------------------------------------------

// evaluate a fresh starting position and set up the material and heatmap starting values
// so that incremental updates during make move will have the correct starting point
func (pos *Position) evalPosAtStart() {

	pos.logTime.allLogTypes[LOG_ONCE_EVAL].start()

	// start with a zero eval for all eval variables
	pos.evalMaterial = 0
	pos.evalHeatmaps = 0
	pos.evalOther = 0
	pos.evalMidVsEndStage = 0

	// ------------------ MATERIAL VALUE + GAME STAGE VALUE -----------------
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {
			pieces := pos.pieces[side][pieceType]
			pieceCount := pieces.countBits()

			pos.evalMaterial += evalTableMaterial[side][pieceType] * pieceCount
			pos.evalMidVsEndStage += evalTableGameStage[pieceType] * pieceCount
		}
	}

	// ----------------------------- HEATMAP VALUE --------------------------
	// we only do this after the game stage value is determined above
	for side := 0; side < 2; side++ {
		for pieceType := 0; pieceType < 6; pieceType++ {

			// get the pieces bitboard
			pieces := pos.pieces[side][pieceType]
			for pieces != 0 {

				// get the next piece square
				nextPieceSq := pieces.popBitGetSq()

				// add the heatmap value of that piece on that square to the eval
				evalStage := pos.evalMidVsEndStage
				if evalStage > STAGE_VAL_STARTING { // cap to the max stage value
					evalStage = STAGE_VAL_STARTING
				}
				midValue := evalTableCombinedMid[side][pieceType][nextPieceSq]
				endValue := evalTableCombinedEnd[side][pieceType][nextPieceSq]
				pos.evalHeatmaps += ((midValue * evalStage) + (endValue * (STAGE_VAL_STARTING - evalStage))) / STAGE_VAL_STARTING
			}
		}
	}

	pos.logTime.allLogTypes[LOG_ONCE_EVAL].stop()
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------- Eval: Other --------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	MOBILITY_BONUS int = 3 // bonus for each available knight, bishop and rook move
)

const (
	DOUBLED_PAWN_PENALTY  int = -5  // penalty for a pawn if there are other friendly pawns on that column
	ISOLATED_PAWN_PENALTY int = -15 // penalty for a pawn without other friendly pawns on the 2 columns besides
	PASSED_PAWN_BONUS     int = 15  // bonus for a pawn that has no enemy pawns in front of it in the left, middle, and right columns
	PROTECTED_PAWN_BONUS  int = 5   // bonus for a pawn that is also directly protected by a friendly pawn behind it
)

// evalue a position after for non-incremental evaluations
func (pos *Position) evalPosAfter() {

	pos.logTime.allLogTypes[LOG_EVAL].start()

	// reset other evaluation scores
	pos.evalOther = 0

	// ------------------------------------------------- MOBILITY --------------------------------------------------
	// during move generation we get a counter for the number of mobility moves for each side the last time move gen was done
	// at the moment we only give a bonus to knight, bishop and rook mobility
	// queen mobility is not scored to not incentivise bringing out the queen early
	// pawn mobility is not scored, because the pawn structure will have separate scores
	// and pawn moves will normally dominate in the opening, where we really want other piece mobility
	// king mobility is not scored, because normally we want pieces to surround the king to protect it

	pos.evalOther += pos.evalWhiteMobility * MOBILITY_BONUS
	pos.evalOther -= pos.evalBlackMobility * MOBILITY_BONUS

	// ------------------------------------------------- PAWN STRUCTURE --------------------------------------------------
	// we give penalties and bonuses for good and bad pawn structures
	// we also store the results of the eval in a hash table specifically for pawn structure

	whitePawns := pos.pieces[SIDE_WHITE][PIECE_PAWN]
	blackPawns := pos.pieces[SIDE_BLACK][PIECE_PAWN]

	// we look up the pawn structure hash in a small table
	// if the pawn structure is the same, we just use the last evaluation, else we re-calculate the pawn structure eval
	// and overwrite it in the hash table
	pawnHash := (whitePawns | blackPawns) % PAWN_HASH_TABLE_SIZE_BB
	if whitePawns != pos.evalPawnHashTable[pawnHash].whitePawns || blackPawns != pos.evalPawnHashTable[pawnHash].blackPawns {

		// set the starting eval
		pawnStructureEval := 0

		// white pawns
		whitePawnsPop := pos.pieces[SIDE_WHITE][PIECE_PAWN]
		for whitePawnsPop != 0 {
			pawnSq := whitePawnsPop.popBitGetSq()
			_, pawnCol := rowAndColFromSq(pawnSq)

			// doubled pawns (if more than 1 friendly pawn on the same col)
			friendlyPawnsOnCol := (whitePawns & pawnColumnMasks[pawnCol]).countBits()
			if friendlyPawnsOnCol > 1 {
				pawnStructureEval += DOUBLED_PAWN_PENALTY
			}

			// isolated pawns (if exactly 1 friendly pawn in the mask)
			friendlyPawnsOn3Col := (whitePawns & pawnIsolatedMasks[pawnCol]).countBits()
			if friendlyPawnsOn3Col == 1 {
				pawnStructureEval += ISOLATED_PAWN_PENALTY
			}

			// passed pawns (if no enemy pawns on the 3 columns in front of the pawn)
			enemyPawnsInFront := (blackPawns & pawnPassedMasks[SIDE_WHITE][pawnSq]).countBits()
			if enemyPawnsInFront == 0 {
				pawnStructureEval += PASSED_PAWN_BONUS
			}

			// protected pawns (if the pawn is directly protected by a friendly pawn)
			friendlyPawnsProtecting := (whitePawns & movePawnsAttackingKingMasks[pawnSq][SIDE_WHITE]).countBits()
			if friendlyPawnsProtecting > 0 {
				pawnStructureEval += PROTECTED_PAWN_BONUS
			}
		}

		// black pawns
		blackPawnsPop := pos.pieces[SIDE_BLACK][PIECE_PAWN]
		for blackPawnsPop != 0 {
			pawnSq := blackPawnsPop.popBitGetSq()
			_, pawnCol := rowAndColFromSq(pawnSq)

			// doubled pawns (if more than 1 friendly pawn on the same col)
			friendlyPawnsOnCol := (blackPawns & pawnColumnMasks[pawnCol]).countBits()
			if friendlyPawnsOnCol > 1 {
				pawnStructureEval -= friendlyPawnsOnCol * DOUBLED_PAWN_PENALTY
			}

			// isolated pawns (if exactly 1 friendly pawn in the mask)
			friendlyPawnsOn3Col := (blackPawns & pawnIsolatedMasks[pawnCol]).countBits()
			if friendlyPawnsOn3Col == 1 {
				pawnStructureEval -= ISOLATED_PAWN_PENALTY
			}

			// passed pawns (if no enemy pawns on the 3 columns in front of the pawn)
			enemyPawnsInFront := (whitePawns & pawnPassedMasks[SIDE_BLACK][pawnSq]).countBits()
			if enemyPawnsInFront == 0 {
				pawnStructureEval -= PASSED_PAWN_BONUS
			}

			// protected pawns (if the pawn is directly protected by a friendly pawn)
			friendlyPawnsProtecting := (blackPawns & movePawnsAttackingKingMasks[pawnSq][SIDE_BLACK]).countBits()
			if friendlyPawnsProtecting > 0 {
				pawnStructureEval -= PROTECTED_PAWN_BONUS
			}
		}

		// finally, save the results for use next time
		pos.evalPawnHashTable[pawnHash].whitePawns = whitePawns
		pos.evalPawnHashTable[pawnHash].blackPawns = blackPawns
		pos.evalPawnHashTable[pawnHash].value = pawnStructureEval
	}

	// finally add the pawn structure eval
	pos.evalOther += pos.evalPawnHashTable[pawnHash].value

	pos.logTime.allLogTypes[LOG_EVAL].stop()
}
