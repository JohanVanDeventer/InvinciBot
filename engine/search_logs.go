package main

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Search: Log Details ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// logs details about a search

type LogSearch struct {
	timeMs int // stores the time of the last search in milliseconds

	depth   int // stores the depth for the last search
	qsDepth int // stores the qs depth of the last search

	nodesAtDepth1Plus int // non-leaf nodes
	nodesAtDepth0     int // nodes at depth zero (either qs starting nodes or eval nodes when there is no qs)
	nodesAtDepth1Min  int // quiescence-only nodes

	nodesTTProbe int // number of times the TT was checked for a position
	nodesTTHit   int // number of times a TT entry was used
	nodesTTStore int // number of times a TT entry was stored

	moveOrderedNodes   int // nodes where moves were ordered
	moveUnorderedNodes int // nodes where moves were unordered

	nodesGeneratedLegalMovesFull int // nodes where full legal moves were generated
	nodesGeneratedLegalMovesPart int // nodes where legal moves were generated until at least one is found

	checkExtensions int // nodes where the depth was extended due to a check
}

func (log *LogSearch) resetLog() {
	log.timeMs = 0

	log.depth = 0
	log.qsDepth = 0

	log.nodesAtDepth1Plus = 0
	log.nodesAtDepth0 = 0
	log.nodesAtDepth1Min = 0

	log.nodesTTProbe = 0
	log.nodesTTHit = 0
	log.nodesTTStore = 0

	log.moveOrderedNodes = 0
	log.moveUnorderedNodes = 0

	log.nodesGeneratedLegalMovesFull = 0
	log.nodesGeneratedLegalMovesPart = 0

	log.checkExtensions = 0
}

func (log *LogSearch) getTotalNodes() int {
	return log.nodesAtDepth1Plus + log.nodesAtDepth0 + log.nodesAtDepth1Min
}
