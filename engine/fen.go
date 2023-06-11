package main

import (
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------ Load Position From Fen String -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// split each part of the Fen string separated by a space and loads it into the position
func (pos *Position) loadFenIntoPosition(fenString string) {

	pos.logTime.allLogTypes[LOG_ONCE_LOAD_FEN].start()

	/*
		example starting Fen string:
		rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
		1: board state
		2: white or black to move
		3: castling rights remaining
		4: en-passant target square (such as "c3")
		5: halfmove counter
		6: fullmove counter
	*/

	stringParts := strings.Split(fenString, " ")

	// ---------------- Part 1: Load Position Array ---------------------
	rowCounter := 0
	colCounter := 0

	var positionArray [8][8]string
	var sChar string

	for _, char := range stringParts[0] {
		sChar = string(char)
		switch sChar {
		case "/":
			rowCounter += 1
			colCounter = 0
		case "1", "2", "3", "4", "5", "6", "7", "8":
			num, _ := strconv.Atoi(sChar)
			colCounter += num
		default:
			positionArray[rowCounter][colCounter] = sChar
			colCounter += 1
		}
	}

	var sq int
	var pieceStr string
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			sq = sqFromRowAndCol(7-row, col)
			pieceStr = positionArray[row][col]
			switch pieceStr {
			case "K":
				pos.pieces[SIDE_WHITE][PIECE_KING].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "Q":
				pos.pieces[SIDE_WHITE][PIECE_QUEEN].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "R":
				pos.pieces[SIDE_WHITE][PIECE_ROOK].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "N":
				pos.pieces[SIDE_WHITE][PIECE_KNIGHT].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "B":
				pos.pieces[SIDE_WHITE][PIECE_BISHOP].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "P":
				pos.pieces[SIDE_WHITE][PIECE_PAWN].setBit(sq)
				pos.piecesAll[SIDE_WHITE].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "k":
				pos.pieces[SIDE_BLACK][PIECE_KING].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "q":
				pos.pieces[SIDE_BLACK][PIECE_QUEEN].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "r":
				pos.pieces[SIDE_BLACK][PIECE_ROOK].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "n":
				pos.pieces[SIDE_BLACK][PIECE_KNIGHT].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "b":
				pos.pieces[SIDE_BLACK][PIECE_BISHOP].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			case "p":
				pos.pieces[SIDE_BLACK][PIECE_PAWN].setBit(sq)
				pos.piecesAll[SIDE_BLACK].setBit(sq)
				pos.piecesAll[SIDE_BOTH].setBit(sq)
			}
		}
	}

	// ---------------- Part 2: Side to Move ---------------------
	if stringParts[1] == "w" {
		pos.isWhiteTurn = true
	} else {
		pos.isWhiteTurn = false
	}

	// ---------------- Part 3: Castling Rights ---------------------
	for _, char := range stringParts[2] {
		sChar := string(char)
		switch sChar {
		case "K":
			pos.castlingRights[CASTLE_WHITE_KINGSIDE] = true
		case "Q":
			pos.castlingRights[CASTLE_WHITE_QUEENSIDE] = true
		case "k":
			pos.castlingRights[CASTLE_BLACK_KINGSIDE] = true
		case "q":
			pos.castlingRights[CASTLE_BLACK_QUEENSIDE] = true
		}
	}

	// ---------------- Part 4: En-Passant Target ---------------------
	if stringParts[3] != "-" {
		enPassantTargetSq := stringParts[3]
		enPCol := enPassantTargetSq[0:1]
		enPRow := enPassantTargetSq[1:2]

		var colInt int
		var rowInt int

		switch enPCol {
		case "a":
			colInt = 0
		case "b":
			colInt = 1
		case "c":
			colInt = 2
		case "d":
			colInt = 3
		case "e":
			colInt = 4
		case "f":
			colInt = 5
		case "g":
			colInt = 6
		case "h":
			colInt = 7
		}

		switch enPRow {
		case "1":
			rowInt = 0
		case "2":
			rowInt = 1
		case "3":
			rowInt = 2
		case "4":
			rowInt = 3
		case "5":
			rowInt = 4
		case "6":
			rowInt = 5
		case "7":
			rowInt = 6
		case "8":
			rowInt = 7
		}
		pos.enPassantTargetBB.setBit(sqFromRowAndCol(rowInt, colInt))
	}

	// ---------------- Part 5: Half Moves ---------------------
	pos.halfMoves, _ = strconv.Atoi(stringParts[4])

	// ---------------- Part 6: Full Moves ---------------------
	pos.fullMoves, _ = strconv.Atoi(stringParts[5])

	pos.logTime.allLogTypes[LOG_ONCE_LOAD_FEN].stop()
}
