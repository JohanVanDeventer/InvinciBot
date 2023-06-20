package main

/*

Ideas to consider / implement
=============================

--- Static null move pruning / reverse futility pruning ---
A cheaper way than doing a null move (which is saying the opponent can have a free move),
is saying if the evaluation less a margin is still above beta, take an early cut.
The margin is higher the higher the depth is, because of deep tactics.

We don't do it in check, or at qs nodes, or when beta is close to checkmate.

if currentDepth >= 1 && !inCheck {
	pos.evalPosAfter()
	var staticScore int

	if pos.isWhiteTurn {
		staticScore = pos.evalMaterial + pos.evalHeatmaps + pos.evalOther
	} else {
		staticScore = 0 - (pos.evalMaterial + pos.evalHeatmaps + pos.evalOther)
	}

	var staticNullMoveMargin int
	if currentDepth == 1 {
		staticNullMoveMargin = VALUE_BISHOP
	} else if currentDepth == 2 {
		staticNullMoveMargin = VALUE_ROOK
	} else if currentDepth == 3 {
		staticNullMoveMargin = VALUE_QUEEN
	}

	if (staticScore - staticNullMoveMargin) >= beta {
		return beta, false
	}
}

--- Magic Bitboards ---
1. Add function to generate own magic numbers.
2. Test speed difference switching back to non-magic generation (less memory intensive for TT hits).

--- Search and draws ---
Try scoring the 1st repetition as a draw and not the 2nd to save search depth?

--- Killer moves ---
In the endgame, sort killer moves before good threat moves?
The cutoff rate for those moves seem higher in the LATE endagme (stage value at most 5 or less).

--- Eval ---
- King safety
- Passed pawn bonus

--- SEE ---
Recursive call line negamax where we input the bitboards,
and only look for recaptures on a particular square.
Once this is done, try MVV-LVA again, where "bad" captured have a SEE test.

--- LMR and Pawns ---
Test whether completely removing LMR from pawn pushes is a gain (especially in the endgame).

--- QS TT ---
Add a small TT specially for QS to fit in the cache.

--- Null move ---
Reduce depth from 6, but don't decrease the depth as much (similar to LMR).
Also, allow null move on shallower depths straight to QS (assumes that a null move is refuted by a simple capture).

--- Delta pruninng (futility pruning in qs) ---
If eval + cature piece value < alpha by a margin, just ignore (likely to fail low).

--- Auto tune ---
Add a function to be able to "modify" eval heatmaps and other parameters before engine init.
That way it can be passed from the Python match manager.

--- TT ---
Remove mod operator, replace with bitboard & operations.
This assumes the TT is a power of 2 size.

--- Eval hash table ---
If the evaluation takes long, store the eval results in a hash table instead like the TT.

--- Better 3 fold repetition detection ---
After a pawn move or a capture or a change in castling rights, we can never again have a 3 fold repetition with positions before that.
Therefore we don't need to iterate over all previous zobrist hashes, only those since that half move counter was reset.

--- Better eval ---
- Doubled pawns?
- Isolated pawns?

--- TT Size ---
Test whether increasing the TT size helps improve play at longer time controls (vs lower cache hits)?

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.
We also want to save time if the best move changes from the previous iteration suddenly.
Use node count to sort moves (longer node count means more difficult to refute).

--- Draw by insufficient material ---
Especially king vs king (eval game stage count == 0: all pieces off, then also just check pawn count == 0), king vs king and knight etc.
Also incorporate this into the evaluation (should be given a drawish score).

--- Book moves ---
Play book moves for the first few moves.

--- Reduce memory ---
Switch to eg. uint8 where possible.

Ideas implemented but failed
=============================
[Try implementhing these again at a later time: might then show an improvement]

--- TT Buckets ---
The other bucket was only used about 0%-2% of the time, did not show a big improvement.

--- IID ---
At normal nodes above say depth 5+, if we don't have a hash move,
use IID to get a best move to search first,
because we don't otherwise have a good estimate of a move to cause a cutoff.

We only got some success if the IID reduction was only 1 ply,
and we only test it in the children of quiet moves,
because we did not have a cutoff up to that point, so we expect it to be an ALL node (about 90% likely),
and children of ALL nodes are more likely to be a CUT node.

But it slowed down the search and did not show an improvement.

--- MVV LVA ---
Try MVV-LVA ordering again later.
It did not show a noticeble improvement over cature value - capturing value.

--- Mobility ---
Smooth the mobility as the average of say the past 2 or 3, not just the last one.
Also add back queen mobility, but at say 1/2 or 1/3 of the actual mobility.

This did not show a noticeable improvement.

--- Move gen ---
Generate pawn moves before other pieces (for quiet move ordering to get avg earlier cutoffs).
Also:
- Castling moves in front of quiet moves
- Pawn moves, then knight, bishop, rook, queen

Was unclear, sometimes better, sometimes worse.

--- QS Checks ---
Rather generate all evasions in qs.
And then cannot stand pat in qs in check.

Tested but this results in an ELO loss, so test again later.

--- QS Move Gen ---
Generate only captures in qs - flag the move generator.

The issue was that mobility cannot be scored in qs because we don't generate all the moves.
So even though move gen was 50% faster, the impact on mobility evaluation was too big.

*/
