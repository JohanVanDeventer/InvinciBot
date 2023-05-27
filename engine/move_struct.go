package main

type Move struct {
	fromSq         int // 0-63
	toSq           int // 0-63
	piece          int // PIECE_TYPE
	moveType       int // quiet, capture, castle
	promotionType  int // none, queen, rook, knight, bishop
	moveOrderScore int // for move ordering later
}

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
TODO:

To try and reduce memory, we encode moves as a Bitboard (uint64).

We also don't add too much information, because not all move generation info is always used.
For example at leaf nodes during the search, we only count the number of moves (and ignore whether it was a capture etc.)
So we just set special flags for en-passant, castling and promotions.

Standard Info: 20 bits
----------------------
FromSq: needs 6 bits. 000000 (0-63).
ToSq: needs 6 bits. 000000 (0-63).
PieceType: needs 3 bits. 000 (0-5).
FlagEnPassant: needs 1 bit (set when the move is an en-passant capture).
FlagCastling: needs 1 bit (set when the move is a castling move).
FlagPromotion: needs 3 bits. 000 (set as the promotion type when there is a promotion).

We don't include quiet vs capture encoding (this is calculated as needed in make move and order moves).

Move ordering score: 32 bits
----------------------------
We use the upper 32 bits to add the move ordering score information.
*/

func getEncodedMove(fromSq uint8, toSq uint8, piece uint8, flagEnPassant uint8, flagCastling uint8, flagPromotion uint8) Bitboard {
	var newMove uint64 = 0

	return Bitboard(newMove)
}
