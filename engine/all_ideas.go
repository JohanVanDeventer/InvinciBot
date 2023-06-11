package main

/*

Ideas to consider / implement
=============================

--- Log search ---
Log details by depth and not in total only.

--- Auto tune ---
Add a function to be able to "modify" eval heatmaps and other parameters before engine init.
That way it can be passed from the Python match manager.

--- TT ---
Remove mod operator, replace with something faster.

--- Asserts ---
Add various asserts in the program as tests.
Check that the incremental heatmap eval is the same for the same position but different paths.

--- QS depth ---
Increase / decrease qs depth to see effect (is qs too deep?)

--- Eval hash table ---
If the evaluation takes long, store the eval results in a table instead like the TT.

--- Better 3 fold repetition detection ---
After a pawn move or a capture or a change in castling rights, we can never again have a 3 fold repetition with positions before that.
Therefore we don't need to iterate over all previous zobrist hashes, only those since that half move counter was reset.

--- IID ---
Try IID at nodes where there is no hash move available.

--- TT Buckets ---
2 entries for each TT slot/index to improve hit rates.
Match the entry size to a cache line size (64 bytes).

--- Branching factor ---
Log the branching factor during search to roughly measure improvements.

--- TT Lookup QS ---
Don't save QS nodes but look up in the TT for QS nodes?
Test whether there is an improvement.

--- Better eval ---
Note: needs to check both sides.

Doubled pawns?
Isolated pawns?
Mobility? Simple pseudo legal moves masked with all blockers (don't go into legal moves only)?

--- TT Size ---
Test whether increasing the TT size helps improve play.

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.

--- Draw by insufficient material ---
Especially king vs king (eval game stage count == 0: all pieces off, then also just check pawn count == 0), king vs king and knight etc.

--- 50 move rule hash ---
Check that the TT handles the 50 move rule correctly.

--- Book moves ---
Play book moves for the first few moves.

--- Reduce memory ---
Switch to eg. uint8 where possible.

*/
