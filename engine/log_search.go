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
	orderThreatMoves            int // number of times threat moves were ordered
	copyQuietMoves              int // number of times quiet moves were copied (unordered)
	orderKiller1                int // number of times 1st killer moves were ordered
	orderKiller2                int // number of times 2nd killer moves were ordered
	orderIterativeDeepeningMove int // number of times the best iterative deepening moves were ordered

	// move generation details
	generatedLegalMovesFull int // nodes where full legal moves were generated
	generatedLegalMovesPart int // nodes where legal moves were generated until at least one is found

	// eval details
	evalLeafNodes             int // number of leaf nodes evaluated
	evalQSNodes               int // number of QS nodes evaluated (excluding leaf nodes)
	evalQSStandPatBetaCuts    int // number of beta cuts in quiescence using stand pat
	evalQSStandPatAlphaRaises int // number of alpha raises in quiescence using stand pat

	// special search extensions and cuts
	checkExtensions   int // nodes where the depth was extended due to a check
	nullMoveSuccesses int // tried a null move and got a cutoff
	nullMoveFailures  int // tried a null move and did not get a cutoff

	// details about occurences of each of the main move loops
	searchedThreatMoves int // nodes where all quiet moves were looped over
	searchedQuietMoves  int // nodes where all quiet moves were looped over
	// searchedBestMoves - not needed: we search best moves each time a hash move was successfully ordered

	// details about cutoffs in each of the main move loops
	threatMovesCutoffs int // beta cuts when looping over threat nodes
	quietMovesCutoffs  int // beta cuts when looping over quiet nodes
	bestMovesCutoffs   int // beta cuts when looping over best nodes
	noCutoffs          int // nodes where we returned alpha

	// details about where in the move list the best/cutoff move was found
	movesTriedTotalMoves int // sum of the index of the best moves
	movesTriedCount      int // total best moves found

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

	// loop over the logged branching factors
	cumulativeNodes := 0
	for i := 0; i < log.nodesPerIterationCounter; i++ {
		cumulativeNodes += log.nodesPerIteration[i]

		// if we can still calculate the next branching factor, do it
		if i < (log.nodesPerIterationCounter - 1) {
			branchFactor := float64(log.nodesPerIteration[i+1]-cumulativeNodes) / float64(log.nodesPerIteration[i])
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

func (log *SearchLogger) getTTSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes

	// tt probe
	ttProbeTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttProbe+log.depthLogs[NODE_TYPE_QS].ttProbe)
	summary += "TT Probe: " + strconv.Itoa(ttProbeTotalPercent) + "%. "

	// tt hit
	ttHitExactTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitExact + log.depthLogs[NODE_TYPE_QS].ttHitExact
	ttHitLowerTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitLower + log.depthLogs[NODE_TYPE_QS].ttHitLower
	ttHitUpperTotal := log.depthLogs[NODE_TYPE_NORMAL].ttHitUpper + log.depthLogs[NODE_TYPE_QS].ttHitUpper
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
	ttRetrieveHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttRetrievedHashMove+log.depthLogs[NODE_TYPE_QS].ttRetrievedHashMove)
	summary += "TT Retrieve Hash: " + strconv.Itoa(ttRetrieveHashMoveTotalPercent) + "%. "

	// tt tested hash move
	ttTestedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttTestedHashMove+log.depthLogs[NODE_TYPE_QS].ttTestedHashMove)
	summary += "TT Test Hash: " + strconv.Itoa(ttTestedHashMoveTotalPercent) + "%. "

	// tt used and ordered hash move
	ttUsedAndOrderedHashMoveTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttUsedAndOrderedHashMove+log.depthLogs[NODE_TYPE_QS].ttUsedAndOrderedHashMove)
	summary += "TT Order Hash: " + strconv.Itoa(ttUsedAndOrderedHashMoveTotalPercent) + "%. "

	// tt store
	ttStoreTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].ttStore+log.depthLogs[NODE_TYPE_QS].ttStore)
	summary += "TT Store: " + strconv.Itoa(ttStoreTotalPercent) + "%. "

	return summary
}

func (log *SearchLogger) getMoveOrderingSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.getTotalNodes()

	// order threat moves
	orderThreatMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].orderThreatMoves+log.depthLogs[NODE_TYPE_QS].orderThreatMoves)
	summary += "Order Threat Moves: " + strconv.Itoa(orderThreatMovesPercent) + "%. "

	// copy quiet moves
	copyQuietMovesPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_NORMAL].copyQuietMoves+log.depthLogs[NODE_TYPE_QS].copyQuietMoves)
	summary += "Copy Quiet Moves: " + strconv.Itoa(copyQuietMovesPercent) + "%. "

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

func (log *SearchLogger) getEvalSummary() string {

	// create the summary string
	summary := ""

	// total qs nodes
	totalNodes := log.depthLogs[NODE_TYPE_QS].nodes

	// eval leaf vs other qs nodes
	evalNodesTotalPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].evalLeafNodes+log.depthLogs[NODE_TYPE_QS].evalQSNodes)
	evalNodesLeafPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].evalLeafNodes)
	evalNodesOtherPercent := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].evalQSNodes)

	summary += "Eval Nodes: " + strconv.Itoa(evalNodesTotalPercent) + "% ("
	summary += "leaf: " + strconv.Itoa(evalNodesLeafPercent) + "%, "
	summary += "other qs: " + strconv.Itoa(evalNodesOtherPercent) + "%). "

	// cutoffs
	evalSPBeta := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].evalQSStandPatBetaCuts)
	evalSPAlpha := getPercent(totalNodes, log.depthLogs[NODE_TYPE_QS].evalQSStandPatAlphaRaises)

	summary += "Stand Pat Beta Cuts: " + strconv.Itoa(evalSPBeta) + "%. "
	summary += "Stand Pat Alpha Raises: " + strconv.Itoa(evalSPAlpha) + "%. "

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

func (log *SearchLogger) getNullMoveSummary() string {

	// create the summary string
	summary := ""

	// total normal nodes
	totalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes

	// null moves
	nullMoveTotal := log.depthLogs[NODE_TYPE_NORMAL].nullMoveFailures + log.depthLogs[NODE_TYPE_NORMAL].nullMoveSuccesses

	nullMoveTotalPercent := getPercent(totalNodes, nullMoveTotal)
	nullMoveSuccessPercent := getPercent(nullMoveTotal, log.depthLogs[NODE_TYPE_NORMAL].nullMoveSuccesses)
	nullMoveFailurePercent := getPercent(nullMoveTotal, log.depthLogs[NODE_TYPE_NORMAL].nullMoveFailures)

	summary += "Try Null Moves: " + strconv.Itoa(nullMoveTotalPercent) + "% ("
	summary += "success: " + strconv.Itoa(nullMoveSuccessPercent) + "%, "
	summary += "failure: " + strconv.Itoa(nullMoveFailurePercent) + "%). "

	return summary
}

func (log *SearchLogger) getMoveLoopsNormalSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_NORMAL].nodes
	nodesWithCutoffs := log.depthLogs[NODE_TYPE_NORMAL].movesTriedCount
	nodesWithoutCutoffs := log.depthLogs[NODE_TYPE_NORMAL].noCutoffs
	nodesLoopedOver := nodesWithCutoffs + nodesWithoutCutoffs

	loopedOverNodesPercent := getPercent(totalNodes, nodesLoopedOver)
	nodesWithCutoffsPercent := getPercent(nodesLoopedOver, nodesWithCutoffs)
	nodesWithoutCutoffsPercent := getPercent(nodesLoopedOver, nodesWithoutCutoffs)

	summary += "Looped Over Moves: " + strconv.Itoa(loopedOverNodesPercent) + "% ("
	summary += "cutoffs: " + strconv.Itoa(nodesWithCutoffsPercent) + "%, "
	summary += "no cutoffs: " + strconv.Itoa(nodesWithoutCutoffsPercent) + "%). "

	// best move index
	avgBestMoveIndex := float64(log.depthLogs[NODE_TYPE_NORMAL].movesTriedTotalMoves) / float64(log.depthLogs[NODE_TYPE_NORMAL].movesTriedCount)
	summary += "Best Move Avg Index: " + strconv.FormatFloat(avgBestMoveIndex, 'f', 2, 64) + ", due to cutoffs ("

	// cut type
	bestCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_NORMAL].bestMovesCutoffs)
	threatCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_NORMAL].threatMovesCutoffs)
	quietCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_NORMAL].quietMovesCutoffs)

	summary += "best cutoffs: " + strconv.Itoa(bestCutsPercent) + "%, "
	summary += "threat cutoffs: " + strconv.Itoa(threatCutsPercent) + "%, "
	summary += "quiet cutoffs: " + strconv.Itoa(quietCutsPercent) + "%). "

	return summary
}

func (log *SearchLogger) getMoveLoopsQsSummary() string {

	// create the summary string
	summary := ""

	// total nodes
	totalNodes := log.depthLogs[NODE_TYPE_QS].nodes
	nodesWithCutoffs := log.depthLogs[NODE_TYPE_QS].movesTriedCount
	nodesWithoutCutoffs := log.depthLogs[NODE_TYPE_QS].noCutoffs
	nodesLoopedOver := nodesWithCutoffs + nodesWithoutCutoffs

	loopedOverNodesPercent := getPercent(totalNodes, nodesLoopedOver)
	nodesWithCutoffsPercent := getPercent(nodesLoopedOver, nodesWithCutoffs)
	nodesWithoutCutoffsPercent := getPercent(nodesLoopedOver, nodesWithoutCutoffs)

	summary += "Looped Over Moves: " + strconv.Itoa(loopedOverNodesPercent) + "% ("
	summary += "cutoffs: " + strconv.Itoa(nodesWithCutoffsPercent) + "%, "
	summary += "no cutoffs: " + strconv.Itoa(nodesWithoutCutoffsPercent) + "%). "

	// best move index
	avgBestMoveIndex := float64(log.depthLogs[NODE_TYPE_QS].movesTriedTotalMoves) / float64(log.depthLogs[NODE_TYPE_QS].movesTriedCount)
	summary += "Best Move Avg Index: " + strconv.FormatFloat(avgBestMoveIndex, 'f', 2, 64) + ", due to cutoffs ("

	// cut type
	bestCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_QS].bestMovesCutoffs)
	threatCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_QS].threatMovesCutoffs)
	quietCutsPercent := getPercent(nodesWithCutoffs, log.depthLogs[NODE_TYPE_QS].quietMovesCutoffs)

	summary += "best cutoffs: " + strconv.Itoa(bestCutsPercent) + "%, "
	summary += "threat cutoffs: " + strconv.Itoa(threatCutsPercent) + "%, "
	summary += "quiet cutoffs: " + strconv.Itoa(quietCutsPercent) + "%). "

	return summary
}
