package main

import (
	"strconv"
	"time"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Search: Log Details ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// logs details about a search

type SearchDepthLog struct {

	// total nodes visited details
	nodes int // total number of nodes visited at this depth

	// transposition table details
	ttProbe    int // number of times the TT was checked for a position
	ttHitExact int // number of times a TT entry was used as EXACT
	ttHitLower int // number of times a TT entry was used as LOWERBOUND
	ttHitUpper int // number of times a TT entry was used as UPPERBOUND
	ttStore    int // number of times a TT entry was stored

	// hash move from TT details
	ttRetrievedHashMove      int // number of times a hash move was obtained from the TT
	ttTestedHashMove         int // number of times a hash move obtained from the TT was tested for legality
	ttUsedAndOrderedHashMove int // number of times an obtained hash move from the TT was actually used and ordered

	// move ordering details
	copyThreatMoves             int // number of times threat moves were copied (unordered)
	copyQuietMoves              int // number of times quiet moves were copied (unordered)
	orderThreatMoves            int // number of times threat moves were ordered
	orderKiller1                int // number of times 1st killer moves were ordered
	orderKiller2                int // number of times 2nd killer moves were ordered
	orderIterativeDeepeningMove int // number of times the best iterative deepening moves were ordered

	// move generation details
	generatedLegalMovesFull int // nodes where full legal moves were generated
	generatedLegalMovesPart int // nodes where legal moves were generated until at least one is found

	// eval details
	evalNode int // number of nodes evaluated

	// qs details
	qsLeafNodes           int // number of leaf nodes evaluated
	qsOtherNodes          int // number of QS nodes evaluated (excluding leaf nodes)
	qsStandPatBetaCuts    int // number of beta cuts in quiescence using stand pat
	qsStandPatAlphaRaises int // number of alpha raises in quiescence using stand pat

	// special search extensions and cuts
	checkExtensions         int // nodes where the depth was extended due to a check
	nullMoveSuccesses       int // tried a null move and got a cutoff
	nullMoveFailures        int // tried a null move and did not get a cutoff
	staticNullMovePrunes    int // nodes where we had a static null move prune
	staticNullMoveNonPrunes int // nodes where we did not have a static null move prune

	// details about cutoffs in each of the main move loops
	loopedOverMoves int // number of nodes where we actually looped over moves

	bestMovesCutoffs         int // beta cuts when looping over best nodes
	bestMovesNoCutoffs       int // no beta cuts when looping over best nodes
	bestMovesTriedBeforeCuts int // sum of the index of the best moves
	bestMovesTriedWhenNoCuts int // sum of the total moves when there were no cuts

	threatGoodMovesCutoffs         int // beta cuts when looping over threat nodes
	threatGoodMovesNoCutoffs       int // no beta cuts when looping over threat nodes
	threatGoodMovesTriedBeforeCuts int // sum of the index of the best moves
	threatGoodMovesTriedWhenNoCuts int // sum of the index of the best moves

	threatBadMovesCutoffs         int // beta cuts when looping over threat nodes
	threatBadMovesNoCutoffs       int // no beta cuts when looping over threat nodes
	threatBadMovesTriedBeforeCuts int // sum of the index of the best moves
	threatBadMovesTriedWhenNoCuts int // sum of the index of the best moves

	quietKillerMovesCutoffs         int // beta cuts when looping over quiet nodes
	quietKillerMovesNoCutoffs       int // no beta cuts when looping over quiet nodes
	quietKillerMovesTriedBeforeCuts int // sum of the index of the best moves
	quietKillerMovesTriedWhenNoCuts int // sum of the index of the best moves

	quietOtherMovesCutoffs         int // beta cuts when looping over quiet nodes
	quietOtherMovesNoCutoffs       int // no beta cuts when looping over quiet nodes
	quietOtherMovesTriedBeforeCuts int // sum of the index of the best moves
	quietOtherMovesTriedWhenNoCuts int // sum of the index of the best moves

	// LMR
	lmrReducedNodes         int // other quiet nodes where lmr was applied
	lmrReducedNodesFailures int // other quiet nodes where lmr was applied and it was a failure (re-searched)
	lmrNonReducedNodes      int // other quiet nodes where lmr was not applied
}

const (
	NODE_TYPE_NORMAL int = 0
	NODE_TYPE_QS     int = 1
)

type SearchLogger struct {

	// final depth of the last search
	depth   int // reached depth of the last search
	qsDepth int // reached QS depth of the last search

	// overall time taken for the search
	startTime    time.Time // starts when the search starts
	searchTimeMs int64     // stores the total time of the last search in milliseconds

	// logs the nodes per initial depth to later calculate the branching factors of each iteration
	nodesPerIteration        [MAX_DEPTH + 1]int
	nodesPerIterationCounter int

	// logs for normal moves and qs moves
	depthLogs [2]SearchDepthLog
}

// gets a new blank search logger
func getNewSearchLogger() SearchLogger {
	var newSearchLogger SearchLogger
	return newSearchLogger
}

// start the search timer
func (log *SearchLogger) start() {
	log.startTime = time.Now()
}

// stop the search timer
func (log *SearchLogger) stop() {
	log.searchTimeMs = time.Since(log.startTime).Milliseconds()
}

// gets the total nodes
func (log *SearchLogger) getTotalNodes() int {
	return log.depthLogs[NODE_TYPE_NORMAL].nodes + log.depthLogs[NODE_TYPE_QS].nodes
}

// logs a new iteration of the search
func (log *SearchLogger) logIteration() {
	log.nodesPerIteration[log.nodesPerIterationCounter] = log.getTotalNodes()
	log.nodesPerIterationCounter += 1
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------- Search: Log Details Printouts ---------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// functions to clean and present the search log details in an organised way

func getPercent(total int, part int) int {
	percent := 0
	if total > 0 {
		percent = int((float64(part) / float64(total)) * 100)
	}
	return percent
}

func (log *SearchLogger) getOverallSummary() string {

	// get the node split
	totalNodes := log.getTotalNodes()

	normalNodesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].nodes)
	qsNodesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].nodes)

	// get nps
	knps := int((float64(totalNodes) / float64(log.searchTimeMs)))

	// create the summary string
	summary := ""
	summary += "Depth: " + strconv.Itoa(log.depth) + " (" + strconv.Itoa(normalNodesPercent) + "% normal nodes). "
	summary += "QS Depth: " + strconv.Itoa(log.qsDepth) + " (" + strconv.Itoa(qsNodesPercent) + "% qs nodes). "
	summary += "Nodes: " + strconv.Itoa(totalNodes) + ". "
	summary += "Knps: " + strconv.Itoa(knps) + ". "

	return summary
}

func (log *SearchLogger) getBranchingFactorSummary() string {

	// set up the branching factor variables
	var branchingFactorHistory []float64
	var avgTotal float64
	var avgCount int

	// limit the count to prevent displaying too much info on the terminal
	if log.nodesPerIterationCounter > 12 {
		log.nodesPerIterationCounter = 12
	}

	// loop over the logged branching factors
	incrementalNodes := 0
	prevIncrementalNodes := 0

	for i := 0; i < log.nodesPerIterationCounter; i++ {
		if i == 0 {
			prevIncrementalNodes = 1
			incrementalNodes = log.nodesPerIteration[i]
		} else {
			prevIncrementalNodes = incrementalNodes
			incrementalNodes = log.nodesPerIteration[i] - log.nodesPerIteration[i-1]
		}

		// if we can still calculate the next branching factor, do it
		if i > 0 {
			branchFactor := float64(incrementalNodes) / float64(prevIncrementalNodes)
			branchingFactorHistory = append(branchingFactorHistory, branchFactor)
			avgTotal += branchFactor
			avgCount += 1
		}
	}

	// finally calculate the avg branching factor
	avgBranchingFactor := avgTotal / float64(avgCount)

	// create the summary string
	summary := ""
	summary += "Avg Branch Factor: " + strconv.FormatFloat(avgBranchingFactor, 'f', 2, 64) + " (detail: "
	for _, br := range branchingFactorHistory {
		summary += " " + strconv.FormatFloat(br, 'f', 2, 64) + ","
	}
	summary += " )."

	/*
		// detail of cumulative nodes
		summary += " [detail: "
		for _, br := range log.nodesPerIteration[:log.nodesPerIterationCounter] {
			summary += " " + strconv.Itoa(br) + ","
		}
		summary += " ]."
	*/

	return summary
}

func (log *SearchLogger) getTTNormalSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes

	// tt probe
	ttProbeTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttProbe)
	summary += "TT Probe: " + strconv.Itoa(ttProbeTotalPercent) + "%. "

	// tt hit
	ttHitExactTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitExact
	ttHitLowerTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitLower
	ttHitUpperTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitUpper
	ttHitTotal := ttHitExactTotal + ttHitLowerTotal + ttHitUpperTotal

	ttHitExactPercent := getPercent(ttHitTotal, ttHitExactTotal)
	ttHitLowerPercent := getPercent(ttHitTotal, ttHitLowerTotal)
	ttHitUpperPercent := getPercent(ttHitTotal, ttHitUpperTotal)
	ttHitTotalPercent := getPercent(totalNodes, ttHitTotal)

	summary += "TT Hit: " + strconv.Itoa(ttHitTotalPercent) + "% ("
	summary += "exact: " + strconv.Itoa(ttHitExactPercent) + "%, "
	summary += "lower: " + strconv.Itoa(ttHitLowerPercent) + "%, "
	summary += "upper: " + strconv.Itoa(ttHitUpperPercent) + "%). "

	// tt retrieve hash move
	ttRetrieveHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttRetrievedHashMove)
	summary += "TT Retrieve Hash: " + strconv.Itoa(ttRetrieveHashMoveTotalPercent) + "%. "

	// tt tested hash move
	ttTestedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttTestedHashMove)
	summary += "TT Test Hash: " + strconv.Itoa(ttTestedHashMoveTotalPercent) + "%. "

	// tt used and ordered hash move
	ttUsedAndOrderedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttUsedAndOrderedHashMove)
	summary += "TT Order Hash: " + strconv.Itoa(ttUsedAndOrderedHashMoveTotalPercent) + "%. "

	// tt store
	ttStoreTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttStore)
	summary += "TT Store: " + strconv.Itoa(ttStoreTotalPercent) + "%. "

	return summary
}

func (log *SearchLogger) getTTQsSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_QS].nodes

	// tt probe
	ttProbeTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].ttProbe)
	summary += "TT Probe: " + strconv.Itoa(ttProbeTotalPercent) + "%. "

	// tt hit
	ttHitExactTotal := log.depthLogs[NODE_TYPE_QS].ttHitExact
	ttHitLowerTotal := log.depthLogs[NODE_TYPE_QS].ttHitLower
	ttHitUpperTotal := log.depthLogs[NODE_TYPE_QS].ttHitUpper
	ttHitTotal := ttHitExactTotal + ttHitLowerTotal + ttHitUpperTotal

	ttHitExactPercent := getPercent(ttHitTotal, ttHitExactTotal)
	ttHitLowerPercent := getPercent(ttHitTotal, ttHitLowerTotal)
	ttHitUpperPercent := getPercent(ttHitTotal, ttHitUpperTotal)
	ttHitTotalPercent := getPercent(totalNodes, ttHitTotal)

	summary += "TT Hit: " + strconv.Itoa(ttHitTotalPercent) + "% ("
	summary += "exact: " + strconv.Itoa(ttHitExactPercent) + "%, "
	summary += "lower: " + strconv.Itoa(ttHitLowerPercent) + "%, "
	summary += "upper: " + strconv.Itoa(ttHitUpperPercent) + "%). "

	// tt retrieve hash move
	ttRetrieveHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].ttRetrievedHashMove)
	summary += "TT Retrieve Hash: " + strconv.Itoa(ttRetrieveHashMoveTotalPercent) + "%. "

	// tt tested hash move
	ttTestedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].ttTestedHashMove)
	summary += "TT Test Hash: " + strconv.Itoa(ttTestedHashMoveTotalPercent) + "%. "

	// tt used and ordered hash move
	ttUsedAndOrderedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].ttUsedAndOrderedHashMove)
	summary += "TT Order Hash: " + strconv.Itoa(ttUsedAndOrderedHashMoveTotalPercent) + "%. "

	// tt store
	ttStoreTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].ttStore)
	summary += "TT Store: " + strconv.Itoa(ttStoreTotalPercent) + "%. "

	return summary
}

func (log *SearchLogger) getMoveOrderingSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.getTotalNodes()

	// copy threat moves
	copyThreatMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].copyThreatMoves+log.depthLogs[NODE_TYPE_QS].copyThreatMoves)
	summary += "Copy Threat: " + strconv.Itoa(copyThreatMovesPercent) + "%. "

	// order threat moves
	orderThreatMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].orderThreatMoves+log.depthLogs[NODE_TYPE_QS].orderThreatMoves)
	summary += "Order Threat: " + strconv.Itoa(orderThreatMovesPercent) + "%. "

	// copy quiet moves
	copyQuietMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].copyQuietMoves+log.depthLogs[NODE_TYPE_QS].copyQuietMoves)
	summary += "Copy Quiet: " + strconv.Itoa(copyQuietMovesPercent) + "%. "

	// order killer 1
	orderKiller1Percent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].orderKiller1+log.depthLogs[NODE_TYPE_QS].orderKiller1)
	summary += "Order Killer 1: " + strconv.Itoa(orderKiller1Percent) + "%. "

	// order killer 2
	orderKiller2Percent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].orderKiller2+log.depthLogs[NODE_TYPE_QS].orderKiller2)
	summary += "Order Killer 2: " + strconv.Itoa(orderKiller2Percent) + "%. "

	// order iterative deepening move
	orderIterDeepPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].orderIterativeDeepeningMove+log.depthLogs[NODE_TYPE_QS].orderIterativeDeepeningMove)
	summary += "Order Iter Deep: " + strconv.Itoa(orderIterDeepPercent) + "%. "

	return summary
}

func (log *SearchLogger) getMoveGenerationSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.getTotalNodes()

	// generate full moves
	genFullMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].generatedLegalMovesFull+log.depthLogs[NODE_TYPE_QS].generatedLegalMovesFull)
	summary += "Gen Full Moves: " + strconv.Itoa(genFullMovesPercent) + "%. "

	// generate part moves
	genPartMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].generatedLegalMovesPart+log.depthLogs[NODE_TYPE_QS].generatedLegalMovesPart)
	summary += "Gen Part Moves: " + strconv.Itoa(genPartMovesPercent) + "%. "

	return summary
}

func (log *SearchLogger) getQsSummary() string {

	// create the summary string
	summary := ""

	// total qs nodes
	totalNodes := log.depthLogs[NODE_TYPE_QS].nodes
	leafNodes := log.depthLogs[NODE_TYPE_QS].qsLeafNodes
	otherNodes := log.depthLogs[NODE_TYPE_QS].qsOtherNodes

	// eval leaf vs other qs nodes
	leafPercent := getPercent(totalNodes, leafNodes)
	otherPercent := getPercent(totalNodes, otherNodes)

	summary += "Nodes (leaf: " + strconv.Itoa(leafPercent) + "%, "
	summary += "other: " + strconv.Itoa(otherPercent) + "%). "

	// cutoffs
	qsBetaPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].qsStandPatBetaCuts)
	qsAlphaPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].qsStandPatAlphaRaises)

	summary += "Stand Pat (beta cuts: " + strconv.Itoa(qsBetaPercent) + "%, "
	summary += "alpha raises: " + strconv.Itoa(qsAlphaPercent) + "%). "

	return summary
}

func (log *SearchLogger) getCheckExtensionsSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.getTotalNodes()

	// check extensions
	totalCheckExtensions := log.depthLogs[NODE_TYPE_NORMAL].checkExtensions + log.depthLogs[NODE_TYPE_QS].checkExtensions

	totalCheckExtensionsPercent := getPercent(totalNodes, totalCheckExtensions)
	normalCheckExtensionsPercent := getPercent(totalCheckExtensions, log.depthLogs[NODE_TYPE_NORMAL].checkExtensions)
	qsCheckExtensionsPercent := getPercent(totalCheckExtensions, log.depthLogs[NODE_TYPE_QS].checkExtensions)

	summary += "Check Extensions: " + strconv.Itoa(totalCheckExtensionsPercent) + "% ("
	summary += "normal: " + strconv.Itoa(normalCheckExtensionsPercent) + "%, "
	summary += "qs: " + strconv.Itoa(qsCheckExtensionsPercent) + "%). "

	return summary
}

func (log *SearchLogger) getEvalSummary() string {

	// create the summary string
	summary := ""

	// get the total nodes evaluated
	totalNodes := log.getTotalNodes()
	evalNodes := log.depthLogs[NODE_TYPE_NORMAL].evalNode + log.depthLogs[NODE_TYPE_QS].evalNode
	evalPercent := getPercent(totalNodes, evalNodes)

	summary += "Eval: " + strconv.Itoa(evalPercent) + "%. "

	return summary
}

func (log *SearchLogger) getPruningAndReductionSummary() string {

	// create the summary string
	summary := ""

	// null moves
	totalNormalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes
	nullMoveTotal := log.depthLogs[NODE_TYPE_NORMAL].nullMoveFailures + log.depthLogs[NODE_TYPE_NORMAL].nullMoveSuccesses

	nullMoveTotalPercent := getPercent(totalNormalNodes, nullMoveTotal)
	nullMoveSuccessPercent := getPercent(nullMoveTotal, log.depthLogs[NODE_TYPE_NORMAL].nullMoveSuccesses)
	nullMoveFailurePercent := getPercent(nullMoveTotal, log.depthLogs[NODE_TYPE_NORMAL].nullMoveFailures)

	summary += "Try Null Moves: " + strconv.Itoa(nullMoveTotalPercent) + "% ("
	summary += "success: " + strconv.Itoa(nullMoveSuccessPercent) + "%, "
	summary += "failure: " + strconv.Itoa(nullMoveFailurePercent) + "%). "

	// lmr stats
	lmrReduceNodes := log.depthLogs[NODE_TYPE_NORMAL].lmrReducedNodes
	lmrNonReducedNodes := log.depthLogs[NODE_TYPE_NORMAL].lmrNonReducedNodes
	totalLMRNodes := lmrReduceNodes + lmrNonReducedNodes

	lmrTryPercent := getPercent(totalLMRNodes, lmrReduceNodes)

	summary += "LMR on other quiet moves: " + strconv.Itoa(lmrTryPercent) + "% ("

	lmrReduceNodesFailures := log.depthLogs[NODE_TYPE_NORMAL].lmrReducedNodesFailures
	lmrFailurePercent := getPercent(lmrReduceNodes, lmrReduceNodesFailures)
	summary += "re-searches: " + strconv.Itoa(lmrFailurePercent) + "%). "

	// static null move pruning (SNMP)
	snmPrunes := log.depthLogs[NODE_TYPE_NORMAL].staticNullMovePrunes
	snmNonPrunes := log.depthLogs[NODE_TYPE_NORMAL].staticNullMoveNonPrunes
	totalSNMPTries := snmPrunes + snmNonPrunes

	snmTryPercent := getPercent(totalNormalNodes, totalSNMPTries)
	snmTriesSuccessRate := getPercent(totalSNMPTries, snmPrunes)

	summary += "Try StNullMvPrune: " + strconv.Itoa(snmTryPercent) + "% ("
	summary += "success rate: " + strconv.Itoa(snmTriesSuccessRate) + "%). "

	return summary
}

func (log *SearchLogger) getMoveLoopsNormalSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes
	nodesLoopedOver := log.depthLogs[NODE_TYPE_NORMAL].loopedOverMoves
	nodesWithCutoffs := log.depthLogs[NODE_TYPE_NORMAL].bestMovesCutoffs + log.depthLogs[NODE_TYPE_NORMAL].threatGoodMovesCutoffs + log.depthLogs[NODE_TYPE_NORMAL].threatBadMovesCutoffs + log.depthLogs[NODE_TYPE_NORMAL].quietKillerMovesCutoffs + log.depthLogs[NODE_TYPE_NORMAL].quietOtherMovesCutoffs

	loopedOverNodesPercent := getPercent(totalNodes, nodesLoopedOver)
	nodesWithCutoffsPercent := getPercent(nodesLoopedOver, nodesWithCutoffs)
	nodesWithoutCutoffsPercent := 100 - nodesWithCutoffsPercent

	summary += "Loop: " + strconv.Itoa(loopedOverNodesPercent) + "%. "
	summary += "NoCut: " + strconv.Itoa(nodesWithoutCutoffsPercent) + "%. "
	summary += "Cut: " + strconv.Itoa(nodesWithCutoffsPercent) + "% ("

	// cut type split
	bestCuts := log.depthLogs[NODE_TYPE_NORMAL].bestMovesCutoffs
	bestCutsPercent := getPercent(nodesWithCutoffs, bestCuts)

	threatGoodCuts := log.depthLogs[NODE_TYPE_NORMAL].threatGoodMovesCutoffs
	threatGoodCutsPercent := getPercent(nodesWithCutoffs, threatGoodCuts)

	quietKillerCuts := log.depthLogs[NODE_TYPE_NORMAL].quietKillerMovesCutoffs
	quietKillerCutsPercent := getPercent(nodesWithCutoffs, quietKillerCuts)

	threatBadCuts := log.depthLogs[NODE_TYPE_NORMAL].threatBadMovesCutoffs
	threatBadCutsPercent := getPercent(nodesWithCutoffs, threatBadCuts)

	quietOtherCuts := log.depthLogs[NODE_TYPE_NORMAL].quietOtherMovesCutoffs
	quietOtherCutsPercent := getPercent(nodesWithCutoffs, quietOtherCuts)

	summary += "b: " + strconv.Itoa(bestCutsPercent) + "%, "
	summary += "gT: " + strconv.Itoa(threatGoodCutsPercent) + "%, "
	summary += "kQ: " + strconv.Itoa(quietKillerCutsPercent) + "%, "
	summary += "bT: " + strconv.Itoa(threatBadCutsPercent) + "%, "
	summary += "oQ: " + strconv.Itoa(quietOtherCutsPercent) + "%). "

	// cut success rate
	bestNoCuts := log.depthLogs[NODE_TYPE_NORMAL].bestMovesNoCutoffs
	bestTotal := bestCuts + bestNoCuts
	bestCutsRatio := getPercent(bestTotal, bestCuts)
	summary += "Success (b: " + strconv.Itoa(bestCutsRatio) + "%, "

	threatGoodNoCuts := log.depthLogs[NODE_TYPE_NORMAL].threatGoodMovesNoCutoffs
	threatGoodTotal := threatGoodCuts + threatGoodNoCuts
	threatGoodCutsRatio := getPercent(threatGoodTotal, threatGoodCuts)
	summary += "gT: " + strconv.Itoa(threatGoodCutsRatio) + "%, "

	quietKillerNoCuts := log.depthLogs[NODE_TYPE_NORMAL].quietKillerMovesNoCutoffs
	quietKillerTotal := quietKillerCuts + quietKillerNoCuts
	quietKillerCutsRatio := getPercent(quietKillerTotal, quietKillerCuts)
	summary += "kQ: " + strconv.Itoa(quietKillerCutsRatio) + "%, "

	threatBadNoCuts := log.depthLogs[NODE_TYPE_NORMAL].threatBadMovesNoCutoffs
	threatBadTotal := threatBadCuts + threatBadNoCuts
	threatBadCutsRatio := getPercent(threatBadTotal, threatBadCuts)
	summary += "bT: " + strconv.Itoa(threatBadCutsRatio) + "%, "

	quietOtherNoCuts := log.depthLogs[NODE_TYPE_NORMAL].quietOtherMovesNoCutoffs
	quietOtherTotal := quietOtherCuts + quietOtherNoCuts
	quietOtherCutsRatio := getPercent(quietOtherTotal, quietOtherCuts)
	summary += "oQ: " + strconv.Itoa(quietOtherCutsRatio) + "%). "

	// cut move index
	avgBestMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].bestMovesTriedBeforeCuts) / float64(bestCuts)
	avgBestMoveLength := float64(log.depthLogs[NODE_TYPE_NORMAL].bestMovesTriedWhenNoCuts) / float64(bestNoCuts)

	avgGoodThreatMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].threatGoodMovesTriedBeforeCuts) / float64(threatGoodCuts)
	avgGoodThreatMoveLength := float64(log.depthLogs[NODE_TYPE_NORMAL].threatGoodMovesTriedWhenNoCuts) / float64(threatGoodNoCuts)

	avgBadThreatMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].threatBadMovesTriedBeforeCuts) / float64(threatBadCuts)
	avgBadThreatMoveLength := float64(log.depthLogs[NODE_TYPE_NORMAL].threatBadMovesTriedWhenNoCuts) / float64(threatBadNoCuts)

	avgKillerQuietMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].quietKillerMovesTriedBeforeCuts) / float64(quietKillerCuts)
	avgKillerQuietMoveLength := float64(log.depthLogs[NODE_TYPE_NORMAL].quietKillerMovesTriedWhenNoCuts) / float64(quietKillerNoCuts)

	avgOtherQuietMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].quietOtherMovesTriedBeforeCuts) / float64(quietOtherCuts)
	avgOtherQuietMoveLength := float64(log.depthLogs[NODE_TYPE_NORMAL].quietOtherMovesTriedWhenNoCuts) / float64(quietOtherNoCuts)

	summary += "CutIdx/Len (b: " + strconv.FormatFloat(avgBestMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgBestMoveLength, 'f', 2, 64) + ", "
	summary += "gT: " + strconv.FormatFloat(avgGoodThreatMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgGoodThreatMoveLength, 'f', 2, 64) + ", "
	summary += "kQ: " + strconv.FormatFloat(avgKillerQuietMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgKillerQuietMoveLength, 'f', 2, 64) + ", "
	summary += "bT: " + strconv.FormatFloat(avgBadThreatMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgBadThreatMoveLength, 'f', 2, 64) + ", "
	summary += "oQ: " + strconv.FormatFloat(avgOtherQuietMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgOtherQuietMoveLength, 'f', 2, 64) + "). "

	return summary
}

func (log *SearchLogger) getMoveLoopsQsSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_QS].nodes
	nodesLoopedOver := log.depthLogs[NODE_TYPE_QS].loopedOverMoves
	nodesWithCutoffs := log.depthLogs[NODE_TYPE_QS].bestMovesCutoffs + log.depthLogs[NODE_TYPE_QS].threatGoodMovesCutoffs + log.depthLogs[NODE_TYPE_QS].threatBadMovesCutoffs + log.depthLogs[NODE_TYPE_QS].quietKillerMovesCutoffs + log.depthLogs[NODE_TYPE_QS].quietOtherMovesCutoffs

	loopedOverNodesPercent := getPercent(totalNodes, nodesLoopedOver)
	nodesWithCutoffsPercent := getPercent(nodesLoopedOver, nodesWithCutoffs)
	nodesWithoutCutoffsPercent := 100 - nodesWithCutoffsPercent

	summary += "Loop: " + strconv.Itoa(loopedOverNodesPercent) + "%. "
	summary += "NoCut: " + strconv.Itoa(nodesWithoutCutoffsPercent) + "%. "
	summary += "Cut: " + strconv.Itoa(nodesWithCutoffsPercent) + "% ("

	// cut type split
	bestCuts := log.depthLogs[NODE_TYPE_QS].bestMovesCutoffs
	bestCutsPercent := getPercent(nodesWithCutoffs, bestCuts)

	threatGoodCuts := log.depthLogs[NODE_TYPE_QS].threatGoodMovesCutoffs
	threatGoodCutsPercent := getPercent(nodesWithCutoffs, threatGoodCuts)

	quietKillerCuts := log.depthLogs[NODE_TYPE_QS].quietKillerMovesCutoffs
	quietKillerCutsPercent := getPercent(nodesWithCutoffs, quietKillerCuts)

	threatBadCuts := log.depthLogs[NODE_TYPE_QS].threatBadMovesCutoffs
	threatBadCutsPercent := getPercent(nodesWithCutoffs, threatBadCuts)

	quietOtherCuts := log.depthLogs[NODE_TYPE_QS].quietOtherMovesCutoffs
	quietOtherCutsPercent := getPercent(nodesWithCutoffs, quietOtherCuts)

	summary += "b: " + strconv.Itoa(bestCutsPercent) + "%, "
	summary += "gT: " + strconv.Itoa(threatGoodCutsPercent) + "%, "
	summary += "kQ: " + strconv.Itoa(quietKillerCutsPercent) + "%, "
	summary += "bT: " + strconv.Itoa(threatBadCutsPercent) + "%, "
	summary += "oQ: " + strconv.Itoa(quietOtherCutsPercent) + "%). "

	// cut success rate
	bestNoCuts := log.depthLogs[NODE_TYPE_QS].bestMovesNoCutoffs
	bestTotal := bestCuts + bestNoCuts
	bestCutsRatio := getPercent(bestTotal, bestCuts)
	summary += "Success (b: " + strconv.Itoa(bestCutsRatio) + "%, "

	threatGoodNoCuts := log.depthLogs[NODE_TYPE_QS].threatGoodMovesNoCutoffs
	threatGoodTotal := threatGoodCuts + threatGoodNoCuts
	threatGoodCutsRatio := getPercent(threatGoodTotal, threatGoodCuts)
	summary += "gT: " + strconv.Itoa(threatGoodCutsRatio) + "%, "

	quietKillerNoCuts := log.depthLogs[NODE_TYPE_QS].quietKillerMovesNoCutoffs
	quietKillerTotal := quietKillerCuts + quietKillerNoCuts
	quietKillerCutsRatio := getPercent(quietKillerTotal, quietKillerCuts)
	summary += "kQ: " + strconv.Itoa(quietKillerCutsRatio) + "%, "

	threatBadNoCuts := log.depthLogs[NODE_TYPE_QS].threatBadMovesNoCutoffs
	threatBadTotal := threatBadCuts + threatBadNoCuts
	threatBadCutsRatio := getPercent(threatBadTotal, threatBadCuts)
	summary += "bT: " + strconv.Itoa(threatBadCutsRatio) + "%, "

	quietOtherNoCuts := log.depthLogs[NODE_TYPE_QS].quietOtherMovesNoCutoffs
	quietOtherTotal := quietOtherCuts + quietOtherNoCuts
	quietOtherCutsRatio := getPercent(quietOtherTotal, quietOtherCuts)
	summary += "oQ: " + strconv.Itoa(quietOtherCutsRatio) + "%). "

	// cut move index
	avgBestMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].bestMovesTriedBeforeCuts) / float64(bestCuts)
	avgBestMoveLength := float64(log.depthLogs[NODE_TYPE_QS].bestMovesTriedWhenNoCuts) / float64(bestNoCuts)

	avgGoodThreatMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].threatGoodMovesTriedBeforeCuts) / float64(threatGoodCuts)
	avgGoodThreatMoveLength := float64(log.depthLogs[NODE_TYPE_QS].threatGoodMovesTriedWhenNoCuts) / float64(threatGoodNoCuts)

	avgBadThreatMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].threatBadMovesTriedBeforeCuts) / float64(threatBadCuts)
	avgBadThreatMoveLength := float64(log.depthLogs[NODE_TYPE_QS].threatBadMovesTriedWhenNoCuts) / float64(threatBadNoCuts)

	avgKillerQuietMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].quietKillerMovesTriedBeforeCuts) / float64(quietKillerCuts)
	avgKillerQuietMoveLength := float64(log.depthLogs[NODE_TYPE_QS].quietKillerMovesTriedWhenNoCuts) / float64(quietKillerNoCuts)

	avgOtherQuietMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].quietOtherMovesTriedBeforeCuts) / float64(quietOtherCuts)
	avgOtherQuietMoveLength := float64(log.depthLogs[NODE_TYPE_QS].quietOtherMovesTriedWhenNoCuts) / float64(quietOtherNoCuts)

	summary += "CutIdx/Len (b: " + strconv.FormatFloat(avgBestMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgBestMoveLength, 'f', 2, 64) + ", "
	summary += "gT: " + strconv.FormatFloat(avgGoodThreatMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgGoodThreatMoveLength, 'f', 2, 64) + ", "
	summary += "kQ: " + strconv.FormatFloat(avgKillerQuietMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgKillerQuietMoveLength, 'f', 2, 64) + ", "
	summary += "bT: " + strconv.FormatFloat(avgBadThreatMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgBadThreatMoveLength, 'f', 2, 64) + ", "
	summary += "oQ: " + strconv.FormatFloat(avgOtherQuietMoveIndex, 'f', 2, 64) + "/" + strconv.FormatFloat(avgOtherQuietMoveLength, 'f', 2, 64) + "). "

	return summary
}
