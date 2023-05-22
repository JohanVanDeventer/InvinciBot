package main

// init calls needed for the engine to be able to start generating moves and playing
// a global variable is set so that the engine is only initiated once
// we also don't init everything at the start, only when needed such as when required by uci

var initEngineWasDone bool = false

func initEngine() {
	if !initEngineWasDone {
		initBBReferenceArray()

		initMoveLookupTablePawns()
		initMoveLookupTableKings()
		initMoveLookupTableKnights()
		initMoveLookupTableRays()
		initMoveCastlingMasks()
		initMovePawnAttackingKingMasks()
		initMovePinnedPiecesMasks()

		initTestPositions()
		initBestMovePositions()

		initHashTables()

		initEvalTables()
		initEvalMaterialAndStageTables()

		initGameStateToText()

		initQSDepthLimits()

		initEngineWasDone = true
	}
}
