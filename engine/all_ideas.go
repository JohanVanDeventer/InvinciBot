package main

/*

Ideas to consider / implement
=============================

--- Legal moves gen stop ---
If at a leaf node and need to determine checkmate/stalemate, we just need to generate up to the first available move.
Once we have ONE move, we know it's not checkmate or stalemate.

--- TT Buckets ---
2 entries for each TT slot/index to improve hit rates.

--- Faster move gen ---
Only calculate pinned pieces where sliders are in line with the enemy king.

--- Branching factor ---
Log the branching factor during search to roughly measure improvements.

--- TT Lookup Qs ---
Don't save Qs nodes but look up in the TT for Qs nodes?
Test whether there is an improvement.

--- Better eval ---
Note: needs to check both sides.

Doubled pawns?
Isolated pawns?
Mobility? Simple pseudo legal moves masked with all blockers (don't go into legal moves only)?
Bishop pair bonus?

--- TT Size ---
Test whether increasing the TT size helps improve play.

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.

--- Draw by insufficient material ---
Especially king vs king (eval game stage count == 0: all pieces off, then also just check pawn count == 0)

--- 50 move rule hash ---
Check that the TT handles the 50 move rule correctly.

--- Book moves ---
Play book moves for the first few moves.

--- Reduce memory ---
Switch to eg. uint8 where possible.

-- Heatmap Tuning ---
Tune heatmap values automatically.

*/
