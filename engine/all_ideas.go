package main

/*

Ideas to consider / implement
=============================

--- QS depth ---
Increase / decrease qs depth to see effect (is qs too deep?)

--- Eval hash table ---
If the evaluation takes long, store the eval results in a table instead like the TT.

--- 3 fold repetition detection ---
After a pawn move or a capture or a change in castling rights, we can never again have a 3 fold repetition with positions before that.
Therefore we don't need to iterate over all previous zobrist hashes, only those since that half move counter was reset.

-- Hash Move ---
Store best moves in the TT.
Apply internal iterative deepening if no hash move is available.

--- Quiet moves ---
Move the quiet move ordering (killers etc.) only after threat moves have been looped over?
Can't do because we need quiet moves to sort the best move from the previous iteration.
Unless we later put the root node call in a separate function.

--- Move in TT ---
Also store best move if available in each TT entry.

--- TT Buckets ---
2 entries for each TT slot/index to improve hit rates.

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
