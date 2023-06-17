package main

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------- Transposition Table: Background ---------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

--------------------------------------------- TT Entry ------------------------------------------------
The transposition table stores node information during a search using a zobrist key derived from the position.
Each node stored in the TT is a struct and needs to have information about:
- Hash (the Zobrist hash of the position that is saved in this TT entry)
- Depth (the remaining depth at which the node was searched)
- Flag (specifying the type of bound we got on the search: exact, lower bound, upper bound)
- Value (the negamax value from the previous search)
- Move (the best move at the node from the previous search)

Each TT entry contains:
- Hash: 1 x uint64 = 1 x 8 = 8 bytes.
- Move: 1 x uint64 = 1 x 8 = 8 bytes.
- Value: 1 x int32 = 1 x 4 = 4 bytes.
- Depth: 1 x uint8 = 1 x 1 = 1 byte.
- Flag: 1 x uint8 = 1 x 1 = 1 byte.

The TT index key also needs to be stored: 1 x uint32 (max value of 4bil) = 1 x 4 = 4 bytes.

Therefore each TT entry is about 26 bytes.
*/

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ TT Entry ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	TT_FLAG_EXACT      uint8 = 0
	TT_FLAG_UPPERBOUND uint8 = 1
	TT_FLAG_LOWERBOUND uint8 = 2
)

type TTEntry struct {
	zobristHash Bitboard // zobrist hash of the position
	move        Move     // previous best move at the node
	value       int32    // negamax search value
	depth       uint8    // depth of the search
	flag        uint8    // exact, lower or upperbound
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------------- TT -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	TT_SIZE_IN_MB              int      = 12
	TT_SIZE_PER_ENTRY_IN_BYTES int      = 26
	TT_SIZE_MAX                int      = TT_SIZE_IN_MB * 1024 * 1024 / TT_SIZE_PER_ENTRY_IN_BYTES
	TT_SIZE_MAX_BB             Bitboard = Bitboard(TT_SIZE_MAX)
)

type TTKey uint32

type TranspositionTable struct {
	entries [TT_SIZE_MAX]TTEntry
}

// returns a newly created TT with a pre-allocated maximum size
func getNewTT() *TranspositionTable {
	newTT := TranspositionTable{}
	return &newTT
}

// this will give TT index keys from 0 (inclusive) to TT_SIZE_MAX (exclusive)
func getTTKeyFromPosHash(posHash Bitboard) TTKey {
	return TTKey(posHash % TT_SIZE_MAX_BB)
}

// this will store a new TT entry with the provided values
func (tt *TranspositionTable) storeNewTTEntry(zobristHashToStore Bitboard, move Move, value int32, depth uint8, flag uint8) {
	ttKey := getTTKeyFromPosHash(zobristHashToStore)
	newTTEntry := TTEntry{zobristHashToStore, move, value, depth, flag}
	tt.entries[ttKey] = newTTEntry
}

// this will search the TT for a given hash, and return the TT bucket and success flag
func (tt *TranspositionTable) getTTEntry(zobristHashToGet Bitboard) (TTEntry, bool) {
	ttEntry := tt.entries[getTTKeyFromPosHash(zobristHashToGet)]
	success := ttEntry.zobristHash == zobristHashToGet
	return ttEntry, success
}
