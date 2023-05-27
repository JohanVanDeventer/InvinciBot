package main

/*
We encode a move in a single uint64.
*/

type Move uint64

const fullMove Move = 0xffffffffffffffff
const emptyMove Move = 0x0

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Constants -----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

const (
	// specify the type of move, so it's easier in the search to filter later
	MOVE_TYPE_QUIET      int = 0
	MOVE_TYPE_CAPTURE    int = 1
	MOVE_TYPE_CASTLE     int = 2
	MOVE_TYPE_EN_PASSANT int = 3

	// specify the type of promotion
	PROMOTION_NONE   int = 0
	PROMOTION_QUEEN  int = PIECE_QUEEN
	PROMOTION_ROOK   int = PIECE_ROOK
	PROMOTION_KNIGHT int = PIECE_KNIGHT
	PROMOTION_BISHOP int = PIECE_BISHOP

	// specify the type of castling
	CASTLE_WHITE_KINGSIDE  int = 0
	CASTLE_WHITE_QUEENSIDE int = 1
	CASTLE_BLACK_KINGSIDE  int = 2
	CASTLE_BLACK_QUEENSIDE int = 3
)

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Move Encoding ---------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

To try and reduce memory and copy overhead during search and move ordering, we encode moves as a single uint64.

0000000000000000000000000000000000000000000000000000000000000000
                                                    	|------| From Sq: 0
												|------|         To Sq: 8
											|--|                 Piece: 16
										|--|                     Move Type: 20
									|--|                         Promotion Type: 24
								|--|                             Unused Bits: 28+

*/

// constants specifying the location of encoding each move part
const (
	MOVE_SHIFT_FROM           int = 0
	MOVE_SHIFT_TO             int = 8
	MOVE_SHIFT_PIECE          int = 16
	MOVE_SHIFT_MOVE_TYPE      int = 20
	MOVE_SHIFT_PROMOTION_TYPE int = 28
	//MOVE_SHIFT_MOVE_ORDERING_SCORE int = 32
)

// constants specifying masks for retrieving each move part
const (

	// fixed-width masks
	MOVE_BIT_MASK_4_BITS  Move = 0xffffffffffffffff >> (64 - 4)
	MOVE_BIT_MASK_8_BITS  Move = 0xffffffffffffffff >> (64 - 8)
	MOVE_BIT_MASK_32_BITS Move = 0xffffffffffffffff >> (64 - 32)

	// masks set at the specific bits where the move info is encoded
	MOVE_MASK_FROM           = fullMove & (MOVE_BIT_MASK_8_BITS << MOVE_SHIFT_FROM)
	MOVE_MASK_TO             = fullMove & (MOVE_BIT_MASK_8_BITS << MOVE_SHIFT_TO)
	MOVE_MASK_PIECE          = fullMove & (MOVE_BIT_MASK_4_BITS << MOVE_SHIFT_PIECE)
	MOVE_MASK_MOVE_TYPE      = fullMove & (MOVE_BIT_MASK_4_BITS << MOVE_SHIFT_MOVE_TYPE)
	MOVE_MASK_PROMOTION_TYPE = fullMove & (MOVE_BIT_MASK_4_BITS << MOVE_SHIFT_PROMOTION_TYPE)
	//MOVE_MASK_MOVE_ORDERING_SCORE = fullMove & (MOVE_BIT_MASK_32_BITS << MOVE_SHIFT_MOVE_ORDERING_SCORE)
)

// get a single move encoded for all information except for the move ordering score
func getEncodedMove(fromSq int, toSq int, piece int, moveType int, promotionType int) Move {
	return Move(fromSq) | (Move(toSq) << MOVE_SHIFT_TO) | (Move(piece) << MOVE_SHIFT_PIECE) | (Move(moveType) << MOVE_SHIFT_MOVE_TYPE) |
		(Move(promotionType) << MOVE_SHIFT_PROMOTION_TYPE)
}

// --------------------------------------------------------------------------------------------------------------------
// --------------------------------------------- Move Information Retrieval -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Functions to extract information from already-encoded moves.
*/

func (move *Move) getFromSq() int {
	return int((*move & MOVE_MASK_FROM) >> MOVE_SHIFT_FROM)
}

func (move *Move) getToSq() int {
	return int((*move & MOVE_MASK_TO) >> MOVE_SHIFT_TO)
}

func (move *Move) getPiece() int {
	return int((*move & MOVE_MASK_PIECE) >> MOVE_SHIFT_PIECE)
}

func (move *Move) getMoveType() int {
	return int((*move & MOVE_MASK_MOVE_TYPE) >> MOVE_SHIFT_MOVE_TYPE)
}

func (move *Move) getPromotionType() int {
	return int((*move & MOVE_MASK_PROMOTION_TYPE) >> MOVE_SHIFT_PROMOTION_TYPE)
}

//func (move *Move) getMoveOrderingScore() int {
//	return int((*move & MOVE_MASK_MOVE_ORDERING_SCORE) >> MOVE_SHIFT_MOVE_ORDERING_SCORE)
//}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Move Information Update --------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// encodes the move ordering score into the move (assumes those bits were not set previously)
//func (move *Move) setMoveOrderingScore(score int) {
//	*move |= (Move(score) << MOVE_SHIFT_MOVE_ORDERING_SCORE)
//}
