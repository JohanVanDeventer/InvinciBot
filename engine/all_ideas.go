package main

/*

Ideas to consider / implement
=============================

--- Stats ---
Log the node count at each depth to get a better overall feel for what depths to set various settings at.

--- History and Heatmaps ---
For any quiet history move sorting, if a move has a good history, give a constant large bonus (say 5000).
Then for any moves with a history score of 0, use the heatmap gain/loss to sort those.
Make sure to also give a bonus of say 1000 for those moves, so that the max heatmap loss will not cause below 0.
Remember to also scale using game stage.

--- MVV LVA ---
Just for the queen, if the queen is captured, give another 1000 point bonus.
It can never be a bad capture, and will possibly reduce qs queen plunder raids.

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.
We also want to save time if the best move changes from the previous iteration suddenly.
Use node count to sort moves (longer node count means more difficult to refute).

--- History Ordering and LMR ---
Don't do LMR on say the first x other quiet moves once we have history ordering.

--- Eval Normalizing ---
The eval terms should be "normalized" around zero.
For example, don't give a mobility bonus for each knight move,
but for each knight move above average (say 4).

This should apply to all eval terms. Above normal is a bonus, below normal is a penalty.

--- Null Move Pruning Eval ---
Use full eval and not just simple eval for null moves, now that other eval is more important.

--- Draw Detection ---
Add detection of drawn endgames, even though we are up material.
Split between strict draw (draw according to the rules), or drawish (tendency to be a draw).

Strict draws (cannot be checkmate):
___________________________________
- k vs k
- kn v k
- kb v k

Drawish (tends to be a draw):
_____________________________
(single queens or rooks)
- kq vs kq
- kr vs kr
(minors only):
- kn vs kn
- kn vs kb
- kb vs kb
(queen vs minors):
- kq vs krr
- kq vs kbb
- kq vs knn
(rook vs minors):
- krb vs kr
- krn vs kr

--- Eval bishop pair ---
Bonus for having the bishop pair.

--- Index speedup ---
Change the /8 and %8 indexing to just look up the row and column directly (precomputed), if it is slow.

--- Magic Bitboards ---
1. Add function to generate own magic numbers.
2. Test speed difference switching back to non-magic generation (less memory intensive for TT hits).

--- Search and draws ---
Try scoring the 1st repetition as a draw and not the 2nd to save search depth?

--- Killer moves in endgame ---
In the endgame, sort killer moves before good threat moves?
The cutoff rate for those moves seem higher in the LATE endagme (stage value at most 5 or less).

--- SEE ---
Recursive call line negamax where we input the bitboards,
and only look for recaptures on a particular square.
Once this is done, try MVV-LVA again, where "bad" captured have a SEE test.

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

--- TT Size ---
Test whether increasing the TT size helps improve play at longer time controls (vs lower cache hits)?

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

--- QS Checks ---
Rather generate all evasions in qs.
And then cannot stand pat in qs in check.

Tested but this results in an ELO loss, so test again later.

--- QS Move Gen ---
Generate only captures in qs - flag the move generator.

The issue was that mobility cannot be scored in qs because we don't generate all the moves.
So even though move gen was 50% faster, the impact on mobility evaluation was too big.

--- Eval King Safety ---
Use the move generation done at each turn, and count the number of squares next to the king that are attacked,
(we already calculate this for move generation).
Give a penalty for each square that is attacked.
Additionally, count the king checks as part of that score (to de-incentivise allowing checks).
Scale this down towards the endgame.

This did not show an improvement, so come back to this later.

*/
