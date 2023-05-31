package main

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------- Transposition Table: Background ---------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

The transposition table stores node information during a search using a zobrist key derived from the position.
Each node stored in the TT is a struct and needs to have information about:
- Depth int (the remaining depth at which the node was searched)
- Flag int (exact, lower bound, upper bound)
- Value int (the value from the previous search)

The TT has a max number of nodes it can store.
We index the TT through: (ZobristHash) mod (TT Size).
Therefore each unique Zobrist key can map to a specific TT entry based on the TT size.

This will cause some key collisions (different zobrist hashes map to the same key).
In that case, we simply overwrite the old entry (other approaches may be implemented later).

We want to keep the TT size manageable to try and fit inside the L3 cache at least (can be later changed through testing).
So we calculate a max size for the TT based on the size of each node stored.

Each node contains:
- Zobrist hash: 1 x 64bit uint = 1 x 8 = 8 bytes.
- Depth: 1 x uint8 = 1 x 1 = 1 byte.
- Flag: 1 x uint8 = 1 x 1 = 1 byte.
- Value: 1 x int = 1 x 8 = 8 bytes.
Therefore each node is about 18 bytes.

The TT index key also needs to be stored:
- 1 x uint32 (max value of 4bil) = 1 x 4 = 4 bytes

So the total per TT entry is 22 bytes.

*/

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------------ TT Entry -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	TT_FLAG_EXACT      uint8 = 0
	TT_FLAG_UPPERBOUND uint8 = 1
	TT_FLAG_LOWERBOUND uint8 = 2
)

type TTEntry struct {
	zobristHash Bitboard
	depth       uint8
	flag        uint8
	value       int32
	move        Move
}

/*
type TTEntry struct {
	zobristHash Bitboard
	depth       uint8
	flag        uint8
	value       int32
}
*/

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------------------- TT -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	TT_SIZE_IN_MB              int      = 10
	TT_SIZE_PER_ENTRY_IN_BYTES int      = 22
	TT_SIZE_MAX                int      = TT_SIZE_IN_MB * 1024 * 1024 / TT_SIZE_PER_ENTRY_IN_BYTES
	TT_SIZE_MAX_BB             Bitboard = Bitboard(TT_SIZE_MAX)
)

type TTKey uint32

type TranspositionTable struct {
	entries map[TTKey]TTEntry
}

// returns a newly created TT with a pre-allocated maximum size
func getNewTT() TranspositionTable {
	newTT := TranspositionTable{make(map[TTKey]TTEntry, TT_SIZE_MAX)}
	return newTT
}

// this will give TT index keys from 0 (inclusive) to TT_SIZE_MAX (exclusive)
func getTTKeyFromPosHash(posHash Bitboard) TTKey {
	return TTKey(posHash % TT_SIZE_MAX_BB)
}

// this will store a new TT entry with the provided values (overwrite the old value if present)
func (tt *TranspositionTable) storeNewTTEntry(zobristHash Bitboard, depth uint8, flag uint8, value int32, move Move) {
	ttKey := getTTKeyFromPosHash(zobristHash)
	newTTEntry := TTEntry{zobristHash, depth, flag, value, move}
	tt.entries[ttKey] = newTTEntry
}

// this will search the TT for a given hash, and return the TT Entry and whether it exists or not
// the TT Entry will have zero values if "success" is false, so only use the entry when "success" is true
func (tt *TranspositionTable) getTTEntry(zobristHash Bitboard) (TTEntry, bool) {
	ttLookup, lookupSuccess := tt.entries[getTTKeyFromPosHash(zobristHash)]

	// in this case the key is the same, but we need to determine whether the actual hash is the same (possible key collision)
	entrySuccess := false
	if lookupSuccess {
		if ttLookup.zobristHash == zobristHash {
			entrySuccess = true
		}
	}

	return ttLookup, entrySuccess
}
