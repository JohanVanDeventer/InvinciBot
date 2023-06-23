package main

import (
	"sort"
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Quiet Moves History ---------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
To better sort quiet moves (other than killer moves), we have a history table for each search.
A history table is independent of depth (whereas killer moves are dependent on the depth).
The table stores moves in [side][piece][toSq].

If a move causes either a beta cutoff (very good) or an alpha improvement (quite good cause it was better than any capture or killer so far),
we give that move a bonus in the table.
If the move did not improve alpha, we reduce the current score in the table by depth.

The aim is in general to catch good quiet moves.
We also only do this for depths > 4 (where an earlier cutoff will actually have a big effect).
The time impact will be negligible because at least 90%+ of normal nodes are at depth <= 4.

*/

type HistoryTable struct {
	entries [2][6][64]int
}

func getNewHistoryTable() *HistoryTable {
	newHistoryTable := HistoryTable{}
	return &newHistoryTable
}

// give a big bonus for moves causing a beta cutoff
func (h *HistoryTable) goodBetaMove(move Move, currentDepth int, side int) {
	h.entries[side][move.getPiece()][move.getToSq()] += currentDepth * currentDepth
}

// give a small bonus for moves causing an alpha improvement
func (h *HistoryTable) goodAlphaMove(move Move, currentDepth int, side int) {
	h.entries[side][move.getPiece()][move.getToSq()] += currentDepth
}

// reduce the score of moves searched that did not give an alpha improvement
// we have a floor of 0 on the score
func (h *HistoryTable) badAlphaMove(move Move, currentDepth int, side int) {
	previousValue := h.entries[side][move.getPiece()][move.getToSq()]
	if previousValue > currentDepth {
		h.entries[side][move.getPiece()][move.getToSq()] -= currentDepth
	} else {
		h.entries[side][move.getPiece()][move.getToSq()] = 0
	}
}

// returns a slice of other quiet moves ordered from best to worst based on history scores
func (pos *Position) orderQuietHistoryMoves(moves []Move, historyTable *HistoryTable) {

	// get the side
	side := SIDE_BLACK
	if pos.isWhiteTurn {
		side = SIDE_WHITE
	}

	// loop over moves to score them
	for i, move := range moves {

		// get the relevant move information
		piece := move.getPiece()
		toSq := move.getToSq()

		// we only update the score if the history score is > 0
		historyScore := historyTable.entries[side][piece][toSq]
		if historyScore > 0 {
			moves[i].setMoveOrderingScore(historyScore)
		}
	}

	// now sort the moves
	// define the custom comparator function
	// sort the moves based on the scores using the comparator function
	sort.Slice(moves, func(i, j int) bool { return moves[i].getMoveOrderingScore() > moves[j].getMoveOrderingScore() })

	// finally clear the move ordering scores (to make comparison easier later to other moves)
	for i := range moves {
		moves[i].clearMoveOrderingScore()
	}
}
