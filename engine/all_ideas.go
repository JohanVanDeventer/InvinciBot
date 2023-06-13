package main

/*

Ideas to consider / implement
=============================

--- In check ---
Don't save killer moves while in check, because they are unlikely to be a good move in sibling nodes.
Because they will likely still be legal, but now not the best move.

--- Bad captures ---
Add killer moves between equal and losing captures.

--- Check extensions and QS ---
Only do check extensions at depth == 0, so we don't enter qs in check but also don't grow the tree too much.

Cannot stand pat in qs in check (tested but this results in an ELO loss?).

--- Quiet and Threat Move Ordering ---
Only order after best moves? Previously did not give an improvement?

--- QS TT ---
Add a small TT specially for QS to fit in the cache.

--- Null move ---
Reduce depth from 6 to 2 or close to that?

--- LMR ---
After we have good move ordering, reduce the depth of later moves.
Don't reduce threat moves.
Also maybe don't reduce quiet pawn moves (changes the structure), only piece moves?

Formula to reduce more the later in the move list?

--- Quiet Move Ordering ---
Add a "next move picker" and not sort all moves up front?

--- TT Startup ---
Takes around 10ms to create a new TT for each search.
Keep TT between searches? - later because it will make debugging harder (first implement other ideas, then test this)

--- Auto tune ---
Add a function to be able to "modify" eval heatmaps and other parameters before engine init.
That way it can be passed from the Python match manager.

--- TT ---
Remove mod operator, replace with something faster.

--- QS depth ---
Increase / decrease qs depth to see effect (is qs too deep/shallow?)
Should not have an effect if the qs pruning works well.

--- Eval hash table ---
If the evaluation takes long, store the eval results in a hash table instead like the TT.

--- Better 3 fold repetition detection ---
After a pawn move or a capture or a change in castling rights, we can never again have a 3 fold repetition with positions before that.
Therefore we don't need to iterate over all previous zobrist hashes, only those since that half move counter was reset.

--- IID ---
Try IID at nodes where there is no hash move available.
Only try at nodes close to the root where better move ordering will have a greater impact.

--- TT Buckets ---
2 entries for each TT slot/index to improve hit rates.
Match the entry size to a cache line size (64 bytes).

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
