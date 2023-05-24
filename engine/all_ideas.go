package main

/*

Ideas to consider / implement
=============================

--- Incremental Heatmap Eval ---
Adjust the heatmap eval incrementally during make move (undo move just restores the previous value).
Need to then take into account the game stage also during captures (not needed for quiet moves, castling, en-passant), and promotions

--- TT Size ---
Test whether increasing the TT size helps improve play.

--- Root Better Move Ordering ---
Better move ordering for the root, because it is called only once at the start of each search, we can afford more time consuming calcs.

--- Magic Bitboards ---
Slider moves magic bitboards.

--- Draw by insufficient material ---
Especially king vs king (eval game stage count == 0: all pieces off, then also just check pawn count == 0)

--- 50 move rule hash ---
Check that the TT handles the 50 move rule correctly.

--- Branching Factor ---
Log the branchging factor of the engine.

--- Book moves ---
Play book moves for the first few moves.

--- Reduce memory ---
Switch to eg. uint8 where possible.

-- Heatmap Tuning ---
Tune heatmap values automatically.

*/
