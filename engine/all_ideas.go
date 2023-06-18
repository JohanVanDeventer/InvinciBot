package main

/*

Ideas to consider / implement
=============================

--- Mobility ---
Smooth the mobility as the average of say the past 2 or 3, not just the last one.
Also add back queen mobility, but at say 1/2 or 1/3 of the actual mobility.

--- Move gen ---
Generate pawn moves before other pieces (for quiet move ordering to get avg earlier cutoffs).

--- IID ---
At normal nodes above say depth 4, if we don't have a hash move,
use IID to get a best move to search first (otherwise we rely blindly on other moves).

--- LMR and Pawns ---
Test whether completely removing LMR from pawn pushes is a gain (especially in the endgame).

--- QS Checks ---
Rather generate all evasions in qs.

Once that is done:
Cannot stand pat in qs in check (tested but this results in an ELO loss, so test again later).

--- MVV LVA ---
Try MVV-LVA ordering again later (no improvement from current ordering).

--- QS TT ---
Add a small TT specially for QS to fit in the cache.

--- Null move ---
Reduce depth from 6, but don't decrease the depth as much (similar to LMR).

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

--- TT Buckets ---
Not a massive improvement, test again later.

--- Better eval ---
- Doubled pawns?
- Isolated pawns?

--- TT Size ---
Test whether increasing the TT size helps improve play.

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.

--- Draw by insufficient material ---
Especially king vs king (eval game stage count == 0: all pieces off, then also just check pawn count == 0), king vs king and knight etc.
Also incorporate this into the evaluation (should be given a drawish score).

--- Book moves ---
Play book moves for the first few moves.

--- Reduce memory ---
Switch to eg. uint8 where possible.

*/
