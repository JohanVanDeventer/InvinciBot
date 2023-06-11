package main

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ Init Engine -------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// init calls needed for the engine to be able to start generating moves and playing
// a global variable is set so that the engine is only initiated once
// we also don't init everything at the start, only when needed such as when required by uci

var initEngineWasDone bool = false

func initEngine() {
	if !initEngineWasDone {

		// general
		initBBReferenceArray()

		// move generation - normal
		initMoveLookupTablePawns()
		initMoveLookupTableKings()
		initMoveLookupTableKnights()
		initMoveLookupTableRays()
		initMoveCastlingMasks()
		initMovePawnAttackingKingMasks()
		initMovePinnedPiecesMasks()

		// move generation - magic
		initMagicMasks()
		initMagicNumbers()
		initMagicShifts()
		initMagicMoveTables()

		// perft and tests
		initTestPositions()
		initIncrementalTestSequences()

		// hashing
		initHashTables()

		// eval
		initEvalTables()
		initEvalMaterialAndStageTables()
		initEvalColumnMasks()

		// terminal gui
		initGameStateToText()

		// search
		initQSDepthLimits()

		initEngineWasDone = true
	}
}
