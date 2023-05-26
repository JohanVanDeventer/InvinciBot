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
