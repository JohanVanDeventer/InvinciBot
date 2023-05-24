package main

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Search: Log Details ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// logs details about a search

type LogSearch struct {
	depth   int // stores the depth for the last search
	qsDepth int // stores the qs depth of the last search

	nodesAtDepth1Plus int // non-leaf nodes
	nodesAtDepth0     int // leaf nodes where quiescence search is started
	nodesAtDepth1Min  int // quiescence-only nodes

	nodesTTHit   int // number of times a TT entry was used
	nodesTTStore int // number of times a TT entry was stored

	timeMs int // stores the time of the last search in milliseconds

	moveOrderedNodes   int // nodes where moves were ordered
	moveUnorderedNodes int // nodes where moves were unordered
}

func (log *LogSearch) resetLog() {
	log.depth = 0
	log.qsDepth = 0

	log.nodesAtDepth1Plus = 0
	log.nodesAtDepth0 = 0
	log.nodesAtDepth1Min = 0

	log.nodesTTHit = 0
	log.nodesTTStore = 0

	log.timeMs = 0

	log.moveOrderedNodes = 0
	log.moveUnorderedNodes = 0
}

func (log *LogSearch) getTotalNodes() int {
	return log.nodesAtDepth1Plus + log.nodesAtDepth0 + log.nodesAtDepth1Min
}
