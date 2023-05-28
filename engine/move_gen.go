package main

import "time"

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------- Legal Move Generation -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
-------------------------------------------------------- Move Struct ----------------------------------------------------
A move needs to have the following information:
1. From square (easy to obtain when generating moves from that square).
2. To square (easy to obtain by popping bits returned from move generation).
3. Move type flag of quiet, capture, castle and promotion:
	- Castle and promotion is easier, those are separate move generation, so a flag can be set in those functions.
	- Quiet or capture needs an extra check (combine move bitboard with enemy bitboard including en-passant).
4. Promotion type (easy to set in promotion move generation instead).

------------------------------------------------- Pseudo Legal vs Legal Moves ---------------------------------------------
Pseudo-legal move generation is fast.

More expensive to now implement is:
- Pinned pieces need to have limited movement.
- When the king is in check, pieces can only move to the attacking ray to block, or attack the piece giving check.
- However in double check only the king can move.

Therefore, start off each move generation by casting rays from the king position to identify direct enemy attackers (direct checks).
Through this we can get a bitboard containing the direct number of enemy attackers. Then the number of checks can be counted.
We can also get a bitboard containing the enemy attackers and squares inbetween them and the king.
We can then limit piece movement to this mask.
Tip: set the mask to all 1's when there are no checks and combining the check mask to the other piece movements (removes if statements later).

Check the number of king attacks:
- if this is one, we need to limit pieces attacks to the attacking or inbetween squares using the mask obtained above.
- if this is two, we only generate king moves.
- if it is one or more, we also do not generate castling moves.
- if it is one, we can generate promotion moves, if the promotion will block an attack (i.e. fall in the mask of squares attacking the king).

We also need a map of squares attacked by the enemy pieces when generating king moves:
- the king cannot move to an attacked square.
- however, remember that a king can only move 1 squares in 8 directions.
- therefore we only need to check whether an enemy piece attacks each of the 8 squares around the king.
- we can put this in a function "isSqAttacked".
- then after generating pseudo-legal moves for the king (which already filters out any blocking friendly pieces, limiting calls to "isSqAttacked"),
- we check whether each of the remaining squares are attacked, and filter those out from the king move generation.

This does not yet address pinned pieces.

For pinned pieces we:
- create 4 types of pin masks and store them in a table (pinned UD, LR, ULtDR, DLtUR).
- pieces can be pinned in these 4 directions.
- we have a initialized move table for that, which for a specific king square, give the rays in those directions.
- we therefore need to determine the "pin type" for pieces.
- we have a function that casts 8 rays from the king, to intialize 4 Bitboards containing pinned piece locations.
- then given the pin type, we can mask their normal moves with the pinned pieces moves from the table.
- the result will give just movement while pinned.

Note: for en-passant, check whether the pawn we would take is part of the pin mask (only case where enemy piece is in pinmask).

At the end, valid moves will be a combination of:
1. Pseudo-legal moves
2. Inbetween squares attacking king mask
3. Pinned pieces mask

*/

// generate all the legal moves for a position
// we set a flag "atLeafCheckForOneMove": at leaf nodes we only need to have one move to determine it is not checkmate or stalemate
func (pos *Position) generateLegalMoves(atLeafCheckForOneMove bool) {

	start_time := time.Now()

	// ------------------------------------------------- Setup ---------------------------------------------
	// reset the moves counter
	pos.availableMovesCounter = 0

	// assign the friendly and enemy pieces and sides
	var frKing Bitboard
	var frQueens Bitboard
	var frRooks Bitboard
	var frKnights Bitboard
	var frBishops Bitboard
	var frPawns Bitboard
	var frSide int
	var enSide int

	if pos.isWhiteTurn {
		frKing = pos.pieces[SIDE_WHITE][PIECE_KING]
		frQueens = pos.pieces[SIDE_WHITE][PIECE_QUEEN]
		frRooks = pos.pieces[SIDE_WHITE][PIECE_ROOK]
		frKnights = pos.pieces[SIDE_WHITE][PIECE_KNIGHT]
		frBishops = pos.pieces[SIDE_WHITE][PIECE_BISHOP]
		frPawns = pos.pieces[SIDE_WHITE][PIECE_PAWN]
		frSide = SIDE_WHITE
		enSide = SIDE_BLACK
	} else {
		frKing = pos.pieces[SIDE_BLACK][PIECE_KING]
		frQueens = pos.pieces[SIDE_BLACK][PIECE_QUEEN]
		frRooks = pos.pieces[SIDE_BLACK][PIECE_ROOK]
		frKnights = pos.pieces[SIDE_BLACK][PIECE_KNIGHT]
		frBishops = pos.pieces[SIDE_BLACK][PIECE_BISHOP]
		frPawns = pos.pieces[SIDE_BLACK][PIECE_PAWN]
		frSide = SIDE_BLACK
		enSide = SIDE_WHITE
	}
	kingSq := frKing.popBitGetSq()

	// ------------------------------------------------- King Attacks ---------------------------------------------
	// get attacks on the king
	//start_1 := time.Now()

	piecesAttKingBB, piecesAndSqAttKingBB := getAttacksOnKing(
		kingSq,
		pos.piecesAll[SIDE_BOTH],
		pos.piecesAll[enSide],
		pos.piecesAll[frSide],
		pos.pieces[enSide][PIECE_QUEEN],
		pos.pieces[enSide][PIECE_ROOK],
		pos.pieces[enSide][PIECE_KNIGHT],
		pos.pieces[enSide][PIECE_BISHOP],
		pos.pieces[enSide][PIECE_PAWN],
		pos.isWhiteTurn)

	//duration_1 := time.Since(start_1).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_KING_ATTACKS].addTime(int(duration_1))

	// ------------------------------------------------- King Moves ---------------------------------------------
	// get king pseudo-legal moves
	// filter out attacked squares
	// the rest are legal moves
	//start_2 := time.Now()
	kingMovesPseudo := getKingMovesPseudo(kingSq)
	kingMovesPseudo &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

	for kingMovesPseudo != 0 {
		nextMoveSq := kingMovesPseudo.popBitGetSq()
		if !isSqAttacked(
			nextMoveSq,
			pos.piecesAll[SIDE_BOTH],
			pos.pieces[frSide][PIECE_KING],
			pos.pieces[enSide][PIECE_QUEEN],
			pos.pieces[enSide][PIECE_ROOK],
			pos.pieces[enSide][PIECE_KNIGHT],
			pos.pieces[enSide][PIECE_BISHOP],
			pos.pieces[enSide][PIECE_PAWN],
			pos.pieces[enSide][PIECE_KING],
			pos.isWhiteTurn) {
			if pos.piecesAll[enSide]&bbReferenceArray[nextMoveSq] != 0 { // capture
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(kingSq, nextMoveSq, PIECE_KING, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			} else { // quiet move
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(kingSq, nextMoveSq, PIECE_KING, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			}
		}
	}
	//duration_2 := time.Since(start_2).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_KING].addTime(int(duration_2))

	// ------------------------------------------------- Pins ---------------------------------------------
	// get pinned pieces bitboards
	//start_3 := time.Now()
	pinsUD, pinsLR, pinsULtDR, pinsDLtUR := getPinnedPieces(
		kingSq,
		pos.piecesAll[SIDE_BOTH],
		pos.piecesAll[frSide],
		pos.pieces[frSide][PIECE_PAWN],
		pos.piecesAll[enSide],
		pos.pieces[enSide][PIECE_QUEEN],
		pos.pieces[enSide][PIECE_ROOK],
		pos.pieces[enSide][PIECE_BISHOP],
		pos.enPassantTargetBB,
		pos.isWhiteTurn)
	pinsCombined := pinsUD | pinsLR | pinsULtDR | pinsDLtUR

	//duration_3 := time.Since(start_3).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_PINS].addTime(int(duration_3))

	// count the attacks on the king
	kingChecks := piecesAttKingBB.countBits()

	// store the number of checks for detecting checkmate/stalemate later
	pos.kingChecks = kingChecks

	// if the number of checks is two, no other moves are possible (already generated king moves above)
	if kingChecks >= 2 {
		return
	}

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	// NOTE: this is the earlies we can check for this, because we need to store the number of king checks for the game state evaluation
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// for single checks, generate moves masked with the king attacked sq mask
	// otherwise the kingInCheckMask is all squares set (i.e. no influence)
	// also, if the king is in check, no castling is allowed (but can do promotions)
	generateCastlingMoves := true
	kingInCheckMask := fullBB
	if kingChecks == 1 {
		kingInCheckMask = piecesAndSqAttKingBB
		generateCastlingMoves = false
	}

	// ------------------------------------------------- Queen Moves ---------------------------------------------
	// now we are ready to generate the other piece moves
	// queens
	//start_4 := time.Now()
	for frQueens != 0 { // while there are pieces left
		nextQueenOriginSq := frQueens.popBitGetSq() // get the square of the piece

		nextQueenMoves := getQueenMovesPseudo( // get the pseudo legal moves of the piece on that square
			nextQueenOriginSq,
			pos.piecesAll[SIDE_BOTH])
		nextQueenMoves &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

		nextQueenMoves &= kingInCheckMask // mask the moves with the king check mask

		if pinsCombined != 0 {
			if bbReferenceArray[nextQueenOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_UD]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_LR]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_DLtUR]
			}
		}

		for nextQueenMoves != 0 {
			nextQueenTargetSq := nextQueenMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextQueenTargetSq] != 0 { // capture
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextQueenOriginSq, nextQueenTargetSq, PIECE_QUEEN, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			} else { // quiet move
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextQueenOriginSq, nextQueenTargetSq, PIECE_QUEEN, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			}
		}
	}
	//duration_4 := time.Since(start_4).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_QUEEN].addTime(int(duration_4))

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// ------------------------------------------------- Rook Moves ---------------------------------------------
	// rooks
	//start_5 := time.Now()
	for frRooks != 0 { // while there are pieces left
		nextRookOriginSq := frRooks.popBitGetSq() // get the square of the piece

		nextRookMoves := getRookMovesPseudo( // get the pseudo legal moves of the piece on that square
			nextRookOriginSq,
			pos.piecesAll[SIDE_BOTH])
		nextRookMoves &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

		nextRookMoves &= kingInCheckMask // mask the moves with the king check mask

		if pinsCombined != 0 {
			if bbReferenceArray[nextRookOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_UD]
			} else if bbReferenceArray[nextRookOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_LR]
			} else if bbReferenceArray[nextRookOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextRookOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_DLtUR]
			}
		}

		for nextRookMoves != 0 {
			nextRookTargetSq := nextRookMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextRookTargetSq] != 0 { // capture
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextRookOriginSq, nextRookTargetSq, PIECE_ROOK, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			} else { // quiet move
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextRookOriginSq, nextRookTargetSq, PIECE_ROOK, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			}
		}
	}
	//duration_5 := time.Since(start_5).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_ROOK].addTime(int(duration_5))

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// ------------------------------------------------- Bishop Moves ---------------------------------------------
	// bishops
	//start_6 := time.Now()
	for frBishops != 0 { // while there are pieces left
		nextBishopOriginSq := frBishops.popBitGetSq() // get the square of the piece

		nextBishopMoves := getBishopMovesPseudo( // get the pseudo legal moves of the piece on that square
			nextBishopOriginSq,
			pos.piecesAll[SIDE_BOTH])
		nextBishopMoves &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

		nextBishopMoves &= kingInCheckMask // mask the moves with the king check mask

		if pinsCombined != 0 {
			if bbReferenceArray[nextBishopOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_UD]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_LR]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_DLtUR]
			}
		}

		for nextBishopMoves != 0 {
			nextBishopTargetSq := nextBishopMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextBishopTargetSq] != 0 { // capture
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextBishopOriginSq, nextBishopTargetSq, PIECE_BISHOP, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			} else { // quiet move
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextBishopOriginSq, nextBishopTargetSq, PIECE_BISHOP, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			}
		}
	}
	//duration_6 := time.Since(start_6).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_BISHOP].addTime(int(duration_6))

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// ------------------------------------------------- Knight Moves ---------------------------------------------
	// knights
	//start_7 := time.Now()
	for frKnights != 0 { // while there are pieces left
		nextKnightOriginSq := frKnights.popBitGetSq() // get the square of the piece

		nextKnightMoves := getKnightMovesPseudo( // get the pseudo legal moves of the piece on that square
			nextKnightOriginSq)
		nextKnightMoves &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

		nextKnightMoves &= kingInCheckMask // mask the moves with the king check mask

		if pinsCombined != 0 {
			if bbReferenceArray[nextKnightOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_UD]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_LR]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_DLtUR]
			}
		}

		for nextKnightMoves != 0 {
			nextKnightTargetSq := nextKnightMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextKnightTargetSq] != 0 { // capture
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextKnightOriginSq, nextKnightTargetSq, PIECE_KNIGHT, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			} else { // quiet move
				//start_time_store_move := time.Now()
				pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextKnightOriginSq, nextKnightTargetSq, PIECE_KNIGHT, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.availableMovesCounter += 1
				//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
				//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
			}
		}
	}
	//duration_7 := time.Since(start_7).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_KNIGHT].addTime(int(duration_7))

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// ------------------------------------------------- Pawn Moves ---------------------------------------------
	// pawns
	//start_8 := time.Now()
	for frPawns != 0 { // while there are pieces left

		nextPawnOriginSq := frPawns.popBitGetSq() // get the square of the piece
		var nextPawnMoves Bitboard

		if pos.isWhiteTurn {
			nextPawnMoves = getPawnMovesWhitePseudo( // get the pseudo legal moves of the piece on that square
				nextPawnOriginSq,
				pos.piecesAll[SIDE_BOTH],
				pos.piecesAll[enSide])
		} else {
			nextPawnMoves = getPawnMovesBlackPseudo( // get the pseudo legal moves of the piece on that square
				nextPawnOriginSq,
				pos.piecesAll[SIDE_BOTH],
				pos.piecesAll[enSide])
		}
		nextPawnMoves &= ^pos.piecesAll[frSide] // mask out moves to friendly pieces

		nextPawnMoves &= kingInCheckMask // mask the moves with the king check mask

		if pinsCombined != 0 {
			if bbReferenceArray[nextPawnOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_UD]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_LR]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_DLtUR]
			}
		}

		for nextPawnMoves != 0 {
			nextPawnTargetSq := nextPawnMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextPawnTargetSq] != 0 { // capture
				if nextPawnTargetSq >= 56 || nextPawnTargetSq <= 7 { // if there is a promotion
					//start_time_store_move1 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_QUEEN)
					pos.availableMovesCounter += 1
					//duration_time_store_move1 := time.Since(start_time_store_move1).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move1))

					//start_time_store_move2 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_ROOK)
					pos.availableMovesCounter += 1
					//duration_time_store_move2 := time.Since(start_time_store_move2).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move2))

					//start_time_store_move3 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_KNIGHT)
					pos.availableMovesCounter += 1
					//duration_time_store_move3 := time.Since(start_time_store_move3).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move3))

					//start_time_store_move4 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_BISHOP)
					pos.availableMovesCounter += 1
					//duration_time_store_move4 := time.Since(start_time_store_move4).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move4))
				} else { // if there is not a promotion
					//start_time_store_move := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
					pos.availableMovesCounter += 1
					//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
				}
			} else { // quiet move
				if nextPawnTargetSq >= 56 || nextPawnTargetSq <= 7 { // if there is a promotion
					//start_time_store_move1 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_QUEEN)
					pos.availableMovesCounter += 1
					//duration_time_store_move1 := time.Since(start_time_store_move1).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move1))

					//start_time_store_move2 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_ROOK)
					pos.availableMovesCounter += 1
					//duration_time_store_move2 := time.Since(start_time_store_move2).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move2))

					//start_time_store_move3 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_KNIGHT)
					pos.availableMovesCounter += 1
					//duration_time_store_move3 := time.Since(start_time_store_move3).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move3))

					//start_time_store_move4 := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_BISHOP)
					pos.availableMovesCounter += 1
					//duration_time_store_move4 := time.Since(start_time_store_move4).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move4))
				} else { // if there is not a promotion
					//start_time_store_move := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_NONE)
					pos.availableMovesCounter += 1
					//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
				}
			}
		}
	}
	//duration_8 := time.Since(start_8).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_PAWN].addTime(int(duration_8))

	// if we are generating moves for leaf nodes, we stop after finding one valid move
	if atLeafCheckForOneMove {
		if pos.availableMovesCounter > 0 {
			return
		}
	}

	// ------------------------------------------------- En-Passant Moves ---------------------------------------------
	// captures en-passant (includes checking for en-passant pawn on pin bitmask)
	// separate check based on the 2 pawn squares that can attack the en-passant target
	//start_9 := time.Now()

	// special rule start >>>
	// the king in check mask does not include pawns checking the king that can be captured en-passant
	// therefore, for en-passant captures only, add it to the check mask if a pawn that is giving check can be captured en-passant

	enPassantKingCheckMask := kingInCheckMask // create a copy of the king check mask to adjust if needed
	if pos.enPassantTargetBB != 0 {           // if there is a possible en-passant capture
		enPassantTarget := pos.enPassantTargetBB           // get the en-passant bitboard
		enPassantTargetSq := enPassantTarget.popBitGetSq() // and get its square

		pawnsCheckingKing := movePawnsAttackingKingMasks[kingSq][enSide] & pos.pieces[enSide][PIECE_PAWN] // get the enemy pawns are giving check

		for pawnsCheckingKing != 0 { // for each of them
			nextCheckerSq := pawnsCheckingKing.popBitGetSq() // get the square of the pawn checking

			if pos.isWhiteTurn {
				if nextCheckerSq+8 == enPassantTargetSq { // if the checking pawn can be captured en-passant
					enPassantTargetUnpopped := pos.enPassantTargetBB  // get a fresh copy for the end
					enPassantKingCheckMask |= enPassantTargetUnpopped // combine with the check mask
				}
			} else {
				if nextCheckerSq-8 == enPassantTargetSq { // if the checking pawn can be captured en-passant
					enPassantTargetUnpopped := pos.enPassantTargetBB  // get a fresh copy for the end
					enPassantKingCheckMask |= enPassantTargetUnpopped // combine with the check mask
				}
			}
		}
	}
	// special rule end <<<

	enPasTargetMasked := pos.enPassantTargetBB & enPassantKingCheckMask // mask with allowable moves when the king is in check

	if enPasTargetMasked != 0 { // if an en-passant capture can be made

		var enPassantCapturedPieceSqBB Bitboard
		if pos.isWhiteTurn {
			enPassantCapturedPieceSqBB = enPasTargetMasked << 8 // get the bitboard of the pawn that can be captured
		} else {
			enPassantCapturedPieceSqBB = enPasTargetMasked >> 8 // get the bitboard of the pawn that can be captured
		}

		enPassantCapturedPieceSq := enPassantCapturedPieceSqBB.popBitGetSq() // get the square of the pawn that can be captured
		if bbReferenceArray[enPassantCapturedPieceSq]&pinsCombined == 0 {    // only if the en passant captured pawn is not pinned, allow the en-passant
			pawnsCanCapture := pos.pieces[frSide][PIECE_PAWN] & movePawnsAttackingKingMasks[enPasTargetMasked.popBitGetSq()][frSide] // which pawns can capture
			for pawnsCanCapture != 0 {                                                                                               // if there are pawns that can capture
				nextPawnOriginSq := pawnsCanCapture.popBitGetSq() // get the origin of the pawn that can capture
				nextPawnMoves := pos.enPassantTargetBB            // get the target of the pawn that can capture

				// now need to check if the CAPTURING pawn is pinned (already checked for CAPTURED pawn pins above)
				if bbReferenceArray[nextPawnOriginSq]&pinsUD != 0 { // if pinned, mask the moves with the pins mask
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_UD]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsLR != 0 { // if pinned, mask the moves with the pins mask
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_LR]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsULtDR != 0 { // if pinned, mask the moves with the pins mask
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_ULtDR]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsDLtUR != 0 { // if pinned, mask the moves with the pins mask
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_DLtUR]
				}

				if nextPawnMoves != 0 { // if there are still moves remaining
					//start_time_store_move := time.Now()
					pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnMoves.popBitGetSq(), PIECE_PAWN, MOVE_TYPE_EN_PASSANT, PROMOTION_NONE)
					pos.availableMovesCounter += 1
					//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
					//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
				}
			}
		}
	}
	//duration_9 := time.Since(start_9).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_EN_PASSANT].addTime(int(duration_9))

	// ------------------------------------------------- Castling Moves ---------------------------------------------
	// castling moves
	//start_10 := time.Now()
	if generateCastlingMoves {
		if pos.isWhiteTurn {
			if pos.castlingRights[CASTLE_WHITE_KINGSIDE] { // if castling is available
				castlingSquares := moveCastlingIsClearMasks[CASTLE_WHITE_KINGSIDE] // get the mask
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]       // if there are no pieces on those squares
				if castlingMasked == 0 {                                           // check if those squares are attacked
					AttSq1 := 5
					AttSq2 := 6
					if !isSqAttacked(
						AttSq1,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						//start_time_store_move := time.Now()
						pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(4, 6, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.availableMovesCounter += 1
						//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
						//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
					}
				}
			}

			if pos.castlingRights[CASTLE_WHITE_QUEENSIDE] { // if castling is available
				castlingSquares := moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE] // get the mask
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]        // if there are no pieces on those squares
				if castlingMasked == 0 {                                            // check if those squares are attacked
					AttSq1 := 3
					AttSq2 := 2
					if !isSqAttacked(
						AttSq1,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						//start_time_store_move := time.Now()
						pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(4, 2, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.availableMovesCounter += 1
						//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
						//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
					}
				}
			}
		} else {
			if pos.castlingRights[CASTLE_BLACK_KINGSIDE] { // if castling is available
				castlingSquares := moveCastlingIsClearMasks[CASTLE_BLACK_KINGSIDE] // get the mask
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]       // if there are no pieces on those squares
				if castlingMasked == 0 {                                           // check if those squares are attacked
					AttSq1 := 61
					AttSq2 := 62
					if !isSqAttacked(
						AttSq1,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						//start_time_store_move := time.Now()
						pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(60, 62, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.availableMovesCounter += 1
						//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
						//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
					}
				}
			}

			if pos.castlingRights[CASTLE_BLACK_QUEENSIDE] { // if castling is available
				castlingSquares := moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE] // get the mask
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]        // if there are no pieces on those squares
				if castlingMasked == 0 {                                            // check if those squares are attacked
					AttSq1 := 59
					AttSq2 := 58
					if !isSqAttacked(
						AttSq1,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2,
						pos.piecesAll[SIDE_BOTH],
						pos.pieces[frSide][PIECE_KING],
						pos.pieces[enSide][PIECE_QUEEN],
						pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT],
						pos.pieces[enSide][PIECE_BISHOP],
						pos.pieces[enSide][PIECE_PAWN],
						pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						//start_time_store_move := time.Now()
						pos.availableMoves[pos.availableMovesCounter] = getEncodedMove(60, 58, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.availableMovesCounter += 1
						//duration_time_store_move := time.Since(start_time_store_move).Nanoseconds()
						//pos.logOther.allLogTypes[LOG_STORE_MOVE_TIME].addTime(int(duration_time_store_move))
					}
				}
			}
		}
	}
	//duration_10 := time.Since(start_10).Nanoseconds()
	//pos.logOther.allLogTypes[LOG_MOVES_CASTLING].addTime(int(duration_10))

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_MOVE_GEN].addTime(int(duration_time))
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------- Generate Pseudo-Legal Moves (excluding castling, promoting, pinned pieces and checks) -------------
// --------------------------------------------------------------------------------------------------------------------
// gets the moves of a piece from a square filtered for blockers (meaning including blockers but not further down each ray)
// the result is then also filtered for any friendly pieces (cannot capture friendly pieces)

func getRookMovesPseudo(sq int, blockers Bitboard) Bitboard {

	blockers &= magicStructsRooks[sq].mask
	blockers *= magicStructsRooks[sq].magic
	blockers >>= (64 - magicStructsRooks[sq].shift)
	return magicRookMovesTable[sq][blockers]

}

func getBishopMovesPseudo(sq int, blockers Bitboard) Bitboard {

	blockers &= magicStructsBishops[sq].mask
	blockers *= magicStructsBishops[sq].magic
	blockers >>= (64 - magicStructsBishops[sq].shift)
	return magicBishopMovesTable[sq][blockers]

}

func getQueenMovesPseudo(sq int, blockers Bitboard) Bitboard {
	var newBitboard = emptyBB

	newBitboard |= getRookMovesPseudo(sq, blockers)
	newBitboard |= getBishopMovesPseudo(sq, blockers)

	return newBitboard
}

func getKingMovesPseudo(sq int) Bitboard {
	return moveKingsTable[sq]
}

func getKnightMovesPseudo(sq int) Bitboard {
	return moveKnightsTable[sq]
}

func getPawnMovesWhitePseudo(sq int, blockers Bitboard, enPieces Bitboard) Bitboard {
	var newBitboard = emptyBB

	// move 1 square forward - filtered for all blockers
	newBitboard |= (moveOnlyPawnsTable[sq][SIDE_WHITE] & ^blockers)

	// add captures if the pawn capture bitboard intersects with enemy pieces
	newBitboard |= (moveAttackPawnsTable[sq][SIDE_WHITE]) & enPieces

	// move 2 squares forward if not blocked
	if sq >= 8 && sq <= 15 {
		if blockers&movePawnDoubleMasks[sq][SIDE_WHITE] == 0 {
			newBitboard |= bbReferenceArray[sq+16]
		}
	}

	return newBitboard
}

func getPawnMovesBlackPseudo(sq int, blockers Bitboard, enPieces Bitboard) Bitboard {
	var newBitboard = emptyBB

	// move 1 square forward - filtered for all blockers
	newBitboard |= (moveOnlyPawnsTable[sq][SIDE_BLACK] & ^blockers)

	// add captures if the pawn capture bitboard intersects with enemy pieces
	newBitboard |= moveAttackPawnsTable[sq][SIDE_BLACK] & enPieces

	// move 2 squares forward if not blocked
	if sq >= 48 && sq <= 55 {
		if blockers&movePawnDoubleMasks[sq][SIDE_BLACK] == 0 {
			newBitboard |= bbReferenceArray[sq-16]
		}
	}

	return newBitboard
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------- Check for Attacks on King -------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// return a bitboard with the bits set that are directly attacking the king
// also return a bitboard of direct attackers and all inbetween squares attacking king using some of the same results
func getAttacksOnKing(
	kingSq int, blockers Bitboard, enPieces Bitboard, frPieces Bitboard, enQ Bitboard, enR Bitboard, enKn Bitboard, enB Bitboard, enP Bitboard, isWhiteTurn bool) (Bitboard, Bitboard) {

	var kingDirectAttackers = emptyBB       // only pieces directly attacking king
	var kingAttackersAndInbetween = emptyBB // pieces directly attacking king and squares inbetween them and the king

	// -------- Knights --------
	// knight attacks on king are just knight moves from the king position masked with the enemy knights bitboard
	knightAttacks := getKnightMovesPseudo(kingSq) & enKn

	kingDirectAttackers |= knightAttacks
	kingAttackersAndInbetween |= knightAttacks

	// -------- Rooks, Bishops, Queens --------
	// set masks for pieces that can attack like a rook or bishop
	enRAndQ := enR | enQ
	enBAndQ := enB | enQ

	// ------------ UP ----------------
	rayAttacksUP := moveRaysTable[kingSq][RAY_UP]             // get the ray
	rayBlockersUP := moveRaysTable[kingSq][RAY_UP] & blockers // get the blockers in the ray
	if rayBlockersUP != 0 {                                   // if there are blockers
		rayBlockerSqUP := rayBlockersUP.getMSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyUP := bbReferenceArray[rayBlockerSqUP] & enRAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyUP != 0 {                                     // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskUP := ^moveRaysTable[rayBlockerSqUP][RAY_UP] // get the mask of the ray after the 1st blocker
			rayFinalUP := rayAfterMaskUP & rayAttacksUP              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyUP         // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalUP                  // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ RIGHT ----------------
	rayAttacksRIGHT := moveRaysTable[kingSq][RAY_RIGHT]             // get the ray
	rayBlockersRIGHT := moveRaysTable[kingSq][RAY_RIGHT] & blockers // get the blockers in the ray
	if rayBlockersRIGHT != 0 {                                      // if there are blockers
		rayBlockerSqRIGHT := rayBlockersRIGHT.getMSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyRIGHT := bbReferenceArray[rayBlockerSqRIGHT] & enRAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyRIGHT != 0 {                                        // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskRIGHT := ^moveRaysTable[rayBlockerSqRIGHT][RAY_RIGHT] // get the mask of the ray after the 1st blocker
			rayFinalRIGHT := rayAfterMaskRIGHT & rayAttacksRIGHT              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyRIGHT               // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalRIGHT                        // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ DOWN ----------------
	rayAttacksDOWN := moveRaysTable[kingSq][RAY_DOWN]             // get the ray
	rayBlockersDOWN := moveRaysTable[kingSq][RAY_DOWN] & blockers // get the blockers in the ray
	if rayBlockersDOWN != 0 {                                     // if there are blockers
		rayBlockerSqDOWN := rayBlockersDOWN.getLSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyDOWN := bbReferenceArray[rayBlockerSqDOWN] & enRAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyDOWN != 0 {                                       // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskDOWN := ^moveRaysTable[rayBlockerSqDOWN][RAY_DOWN] // get the mask of the ray after the 1st blocker
			rayFinalDOWN := rayAfterMaskDOWN & rayAttacksDOWN              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyDOWN             // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalDOWN                      // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ LEFT ----------------
	rayAttacksLEFT := moveRaysTable[kingSq][RAY_LEFT]             // get the ray
	rayBlockersLEFT := moveRaysTable[kingSq][RAY_LEFT] & blockers // get the blockers in the ray
	if rayBlockersLEFT != 0 {                                     // if there are blockers
		rayBlockerSqLEFT := rayBlockersLEFT.getLSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyLEFT := bbReferenceArray[rayBlockerSqLEFT] & enRAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyLEFT != 0 {                                       // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskLEFT := ^moveRaysTable[rayBlockerSqLEFT][RAY_LEFT] // get the mask of the ray after the 1st blocker
			rayFinalLEFT := rayAfterMaskLEFT & rayAttacksLEFT              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyLEFT             // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalLEFT                      // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ UL ----------------
	rayAttacksUL := moveRaysTable[kingSq][RAY_UL]             // get the ray
	rayBlockersUL := moveRaysTable[kingSq][RAY_UL] & blockers // get the blockers in the ray
	if rayBlockersUL != 0 {                                   // if there are blockers
		rayBlockerSqUL := rayBlockersUL.getMSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyUL := bbReferenceArray[rayBlockerSqUL] & enBAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyUL != 0 {                                     // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskUL := ^moveRaysTable[rayBlockerSqUL][RAY_UL] // get the mask of the ray after the 1st blocker
			rayFinalUL := rayAfterMaskUL & rayAttacksUL              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyUL         // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalUL                  // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ UR ----------------
	rayAttacksUR := moveRaysTable[kingSq][RAY_UR]             // get the ray
	rayBlockersUR := moveRaysTable[kingSq][RAY_UR] & blockers // get the blockers in the ray
	if rayBlockersUR != 0 {                                   // if there are blockers
		rayBlockerSqUR := rayBlockersUR.getMSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyUR := bbReferenceArray[rayBlockerSqUR] & enBAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyUR != 0 {                                     // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskUR := ^moveRaysTable[rayBlockerSqUR][RAY_UR] // get the mask of the ray after the 1st blocker
			rayFinalUR := rayAfterMaskUR & rayAttacksUR              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyUR         // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalUR                  // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ DR ----------------
	rayAttacksDR := moveRaysTable[kingSq][RAY_DR]             // get the ray
	rayBlockersDR := moveRaysTable[kingSq][RAY_DR] & blockers // get the blockers in the ray
	if rayBlockersDR != 0 {                                   // if there are blockers
		rayBlockerSqDR := rayBlockersDR.getLSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyDR := bbReferenceArray[rayBlockerSqDR] & enBAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyDR != 0 {                                     // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskDR := ^moveRaysTable[rayBlockerSqDR][RAY_DR] // get the mask of the ray after the 1st blocker
			rayFinalDR := rayAfterMaskDR & rayAttacksDR              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyDR         // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalDR                  // add the final attacks to the list of attackers and squares
		}
	}

	// ------------ DL ----------------
	rayAttacksDL := moveRaysTable[kingSq][RAY_DL]             // get the ray
	rayBlockersDL := moveRaysTable[kingSq][RAY_DL] & blockers // get the blockers in the ray
	if rayBlockersDL != 0 {                                   // if there are blockers
		rayBlockerSqDL := rayBlockersDL.getLSBSq()                              // get the square of the 1st blocker in the ray
		rayFirstBlockerAndEnemyDL := bbReferenceArray[rayBlockerSqDL] & enBAndQ // combine the 1st blocker square with enemy pieces that can attack the king in this ray
		if rayFirstBlockerAndEnemyDL != 0 {                                     // if the blocker is an enemy piece that can attack in this way
			rayAfterMaskDL := ^moveRaysTable[rayBlockerSqDL][RAY_DL] // get the mask of the ray after the 1st blocker
			rayFinalDL := rayAfterMaskDL & rayAttacksDL              // combine the initial ray with the mask to get the final attacks
			kingDirectAttackers |= rayFirstBlockerAndEnemyDL         // add the blocker to the list of direct attacks
			kingAttackersAndInbetween |= rayFinalDL                  // add the final attacks to the list of attackers and squares
		}
	}

	// -------- Pawns --------
	// use the initialized table of pawns that can attack the king
	if isWhiteTurn {
		blackPawnsAttacking := movePawnsAttackingKingMasks[kingSq][SIDE_BLACK] & enP

		kingDirectAttackers |= blackPawnsAttacking
		kingAttackersAndInbetween |= blackPawnsAttacking

	} else {
		whitePawnsAttacking := movePawnsAttackingKingMasks[kingSq][SIDE_WHITE] & enP

		kingDirectAttackers |= whitePawnsAttacking
		kingAttackersAndInbetween |= whitePawnsAttacking
	}

	return kingDirectAttackers, kingAttackersAndInbetween
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------- Check Whether a Square is Attacked ---------------------------------------
// --------------------------------------------------------------------------------------------------------------------

func isSqAttacked(sq int, blockers Bitboard, frKing Bitboard, enQ Bitboard, enR Bitboard, enKn Bitboard, enB Bitboard, enP Bitboard, enK Bitboard, isWhiteTurn bool) bool {

	// mask out the friendly king from blockers, because enemy rays should ignore the king (in instances the king steps 1 square back)
	blockers &= ^frKing

	// now from the square, cast out rays using rook rays, bishop rays, knight rays, and pawn rays
	// knights
	knightsAttacking := getKnightMovesPseudo(sq) & enKn
	if knightsAttacking != 0 {
		return true
	}

	// rooks
	rooksAttacking := getRookMovesPseudo(sq, blockers) & (enR | enQ)
	if rooksAttacking != 0 {
		return true
	}

	// bishops
	bishopsAttacking := getBishopMovesPseudo(sq, blockers) & (enB | enQ)
	if bishopsAttacking != 0 {
		return true
	}

	// pawns
	if isWhiteTurn {
		pawnsAttacking := movePawnsAttackingKingMasks[sq][SIDE_BLACK] & enP
		if pawnsAttacking != 0 {
			return true
		}
	} else {
		pawnsAttacking := movePawnsAttackingKingMasks[sq][SIDE_WHITE] & enP
		if pawnsAttacking != 0 {
			return true
		}
	}

	// kings: also check that the king cannot move to a square attacked by the enemy king
	kingsAttacking := getKingMovesPseudo(sq) & enK
	return kingsAttacking != 0
}

// --------------------------------------------------------------------------------------------------------------------
// ----------------------------------------------- Check for Pinned Pieces --------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// return 4 bitboards, each with the pieces pinned in the 4 main directions (vertical, horizontal, and both diagonals)
func getPinnedPieces(
	kingSq int, blockers Bitboard, frPieces Bitboard, frPawns Bitboard, enPieces Bitboard, enQ Bitboard, enR Bitboard, enB Bitboard, enPTarget Bitboard, isWhiteTurn bool) (
	Bitboard, Bitboard, Bitboard, Bitboard) {

	// create 4 new pin bitboards
	UDpins := emptyBB
	LRpins := emptyBB
	ULtDRpins := emptyBB
	DLtURpins := emptyBB

	// create the enemy pinning pieces lists
	enRAndQ := enR | enQ
	enBAndQ := enB | enQ

	// special rule for UP-DOWN rays
	frPiecesUDOnly := frPieces

	// we need to add the pawn represented by the en-passant target to the friendly pieces list (ignore UD rays)
	// so it can be included in the pin mask
	// because this is the only case that an enemy piece can be pinned
	// then when creating en-passant moves, check if the pawn captured by en-passant is part of the relevant pin mask
	// this will then filter out en-passant revealed checks

	if isWhiteTurn {
		frPieces |= (enPTarget << 8) // the actual pawn is 8 down from the en-passant target sq
	} else {
		frPieces |= (enPTarget >> 8) // the actual pawn is 8 up from the en-passant target sq
	}

	// UP
	// special rule: pawns on the UD rays are not pinned, because the pawn will be replaced by the capturing pawn
	// therefore for UD rays, use the normal friendly pieces and not the adjusted friendly pieces
	pinRayUP := moveRaysTable[kingSq][RAY_UP] // get the ray
	pinBlockersUP := pinRayUP & blockers      // get the blockers in the ray
	if pinBlockersUP != 0 {                   // if there are blockers
		pinBlockerSq1UP := pinBlockersUP.getMSBSq()                                   // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyUP := bbReferenceArray[pinBlockerSq1UP] & frPiecesUDOnly // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyUP != 0 {                                             // if the piece is a friendly piece, continue the ray
			rayFromFriendlyUP := moveRaysTable[pinBlockerSq1UP][RAY_UP] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersUP := rayFromFriendlyUP & blockers   // check for further blockers
			if rayFromFriendlyBlockersUP != 0 {                         // if there are further blockers
				pinBlockerSq2UP := rayFromFriendlyBlockersUP.getMSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyUP := bbReferenceArray[pinBlockerSq2UP] & enRAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyUP != 0 {                                      // if a enemy pinning piece is found
					UDpins.setBit(pinBlockerSq1UP) // set the friendly piece found as pinned
				}
			}
		}
	}

	// DOWN
	// special rule: pawns on the UD rays are not pinned, because the pawn will be replaced by the capturing pawn
	// therefore for UD rays, use the normal friendly pieces and not the adjusted friendly pieces
	pinRayDOWN := moveRaysTable[kingSq][RAY_DOWN] // get the ray
	pinBlockersDOWN := pinRayDOWN & blockers      // get the blockers in the ray
	if pinBlockersDOWN != 0 {                     // if there are blockers
		pinBlockerSq1DOWN := pinBlockersDOWN.getLSBSq()                                   // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyDOWN := bbReferenceArray[pinBlockerSq1DOWN] & frPiecesUDOnly // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyDOWN != 0 {                                               // if the piece is a friendly piece, continue the ray
			rayFromFriendlyDOWN := moveRaysTable[pinBlockerSq1DOWN][RAY_DOWN] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersDOWN := rayFromFriendlyDOWN & blockers     // check for further blockers
			if rayFromFriendlyBlockersDOWN != 0 {                             // if there are further blockers
				pinBlockerSq2DOWN := rayFromFriendlyBlockersDOWN.getLSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyDOWN := bbReferenceArray[pinBlockerSq2DOWN] & enRAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyDOWN != 0 {                                        // if a enemy pinning piece is found
					UDpins.setBit(pinBlockerSq1DOWN) // set the friendly piece found as pinned
				}
			}
		}
	}

	// LEFT
	// special rule: remove friendly pawns able to capture en passant from the friendly piece ray - they also get out of the way
	frPiecesAdjustedLEFT := frPieces
	blockersAdjustedLEFT := blockers
	enPTargetCopyLEFT := enPTarget
	if enPTarget != 0 { // if there is an en-passant possible
		var pawnsCanCaptureEnPassantLEFT Bitboard
		if isWhiteTurn {
			pawnsCanCaptureEnPassantLEFT = movePawnsAttackingKingMasks[enPTargetCopyLEFT.popBitGetSq()][SIDE_WHITE] & frPawns // get the possible pawns that can capture en-passant
		} else {
			pawnsCanCaptureEnPassantLEFT = movePawnsAttackingKingMasks[enPTargetCopyLEFT.popBitGetSq()][SIDE_BLACK] & frPawns // get the possible pawns that can capture en-passant
		}

		// counter-rule: if there are 2 pawns that can capture, the pawn is not pinned,
		if pawnsCanCaptureEnPassantLEFT.countBits() < 2 {
			frPiecesAdjustedLEFT &= ^pawnsCanCaptureEnPassantLEFT // mask out those pawns from friendly pieces
			blockersAdjustedLEFT &= ^pawnsCanCaptureEnPassantLEFT // mask out those pawns from blockers
		}
	}

	// normal rules:
	pinRayLEFT := moveRaysTable[kingSq][RAY_LEFT]        // get the ray
	pinBlockersLEFT := pinRayLEFT & blockersAdjustedLEFT // get the blockers in the ray
	if pinBlockersLEFT != 0 {                            // if there are blockers
		pinBlockerSq1LEFT := pinBlockersLEFT.getLSBSq()                                         // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyLEFT := bbReferenceArray[pinBlockerSq1LEFT] & frPiecesAdjustedLEFT // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyLEFT != 0 {                                                     // if the piece is a friendly piece, continue the ray
			rayFromFriendlyLEFT := moveRaysTable[pinBlockerSq1LEFT][RAY_LEFT]         // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersLEFT := rayFromFriendlyLEFT & blockersAdjustedLEFT // check for further blockers
			if rayFromFriendlyBlockersLEFT != 0 {                                     // if there are further blockers
				pinBlockerSq2LEFT := rayFromFriendlyBlockersLEFT.getLSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyLEFT := bbReferenceArray[pinBlockerSq2LEFT] & enRAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyLEFT != 0 {                                        // if a enemy pinning piece is found
					LRpins.setBit(pinBlockerSq1LEFT) // set the friendly piece found as pinned
				}
			}
		}
	}

	// RIGHT
	// special rule: remove white/black pawns able to capture en passant from the friendly piece ray - they also get out of the way
	frPiecesAdjustedRIGHT := frPieces
	blockersAdjustedRIGHT := blockers
	enPTargetCopyRIGHT := enPTarget
	if enPTarget != 0 { // if there is an en-passant possible
		var pawnsCanCaptureEnPassantRIGHT Bitboard
		if isWhiteTurn {
			pawnsCanCaptureEnPassantRIGHT = movePawnsAttackingKingMasks[enPTargetCopyRIGHT.popBitGetSq()][SIDE_WHITE] & frPawns // get the possible pawns that can capture en-passant
		} else {
			pawnsCanCaptureEnPassantRIGHT = movePawnsAttackingKingMasks[enPTargetCopyRIGHT.popBitGetSq()][SIDE_BLACK] & frPawns // get the possible pawns that can capture en-passant
		}

		// counter-rule: if there are 2 pawns that can capture, the pawn is not pinned,
		if pawnsCanCaptureEnPassantRIGHT.countBits() < 2 {
			frPiecesAdjustedRIGHT &= ^pawnsCanCaptureEnPassantRIGHT // mask out those pawns from friendly pieces
			blockersAdjustedRIGHT &= ^pawnsCanCaptureEnPassantRIGHT // mask out those pawns from blockers
		}
	}

	// normal rules:
	pinRayRIGHT := moveRaysTable[kingSq][RAY_RIGHT]         // get the ray
	pinBlockersRIGHT := pinRayRIGHT & blockersAdjustedRIGHT // get the blockers in the ray
	if pinBlockersRIGHT != 0 {                              // if there are blockers
		pinBlockerSq1RIGHT := pinBlockersRIGHT.getMSBSq()                                          // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyRIGHT := bbReferenceArray[pinBlockerSq1RIGHT] & frPiecesAdjustedRIGHT // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyRIGHT != 0 {                                                       // if the piece is a friendly piece, continue the ray
			rayFromFriendlyRIGHT := moveRaysTable[pinBlockerSq1RIGHT][RAY_RIGHT]         // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersRIGHT := rayFromFriendlyRIGHT & blockersAdjustedRIGHT // check for further blockers
			if rayFromFriendlyBlockersRIGHT != 0 {                                       // if there are further blockers
				pinBlockerSq2RIGHT := rayFromFriendlyBlockersRIGHT.getMSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyRIGHT := bbReferenceArray[pinBlockerSq2RIGHT] & enRAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyRIGHT != 0 {                                         // if a enemy pinning piece is found
					LRpins.setBit(pinBlockerSq1RIGHT) // set the friendly piece found as pinned
				}
			}
		}
	}

	// UL
	pinRayUL := moveRaysTable[kingSq][RAY_UL] // get the ray
	pinBlockersUL := pinRayUL & blockers      // get the blockers in the ray
	if pinBlockersUL != 0 {                   // if there are blockers
		pinBlockerSq1UL := pinBlockersUL.getMSBSq()                             // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyUL := bbReferenceArray[pinBlockerSq1UL] & frPieces // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyUL != 0 {                                       // if the piece is a friendly piece, continue the ray
			rayFromFriendlyUL := moveRaysTable[pinBlockerSq1UL][RAY_UL] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersUL := rayFromFriendlyUL & blockers   // check for further blockers
			if rayFromFriendlyBlockersUL != 0 {                         // if there are further blockers
				pinBlockerSq2UL := rayFromFriendlyBlockersUL.getMSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyUL := bbReferenceArray[pinBlockerSq2UL] & enBAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyUL != 0 {                                      // if a enemy pinning piece is found
					ULtDRpins.setBit(pinBlockerSq1UL) // set the friendly piece found as pinned
				}
			}
		}
	}

	// DR
	pinRayDR := moveRaysTable[kingSq][RAY_DR] // get the ray
	pinBlockersDR := pinRayDR & blockers      // get the blockers in the ray
	if pinBlockersDR != 0 {                   // if there are blockers
		pinBlockerSq1DR := pinBlockersDR.getLSBSq()                             // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyDR := bbReferenceArray[pinBlockerSq1DR] & frPieces // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyDR != 0 {                                       // if the piece is a friendly piece, continue the ray
			rayFromFriendlyDR := moveRaysTable[pinBlockerSq1DR][RAY_DR] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersDR := rayFromFriendlyDR & blockers   // check for further blockers
			if rayFromFriendlyBlockersDR != 0 {                         // if there are further blockers
				pinBlockerSq2DR := rayFromFriendlyBlockersDR.getLSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyDR := bbReferenceArray[pinBlockerSq2DR] & enBAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyDR != 0 {                                      // if a enemy pinning piece is found
					ULtDRpins.setBit(pinBlockerSq1DR) // set the friendly piece found as pinned
				}
			}
		}
	}

	// DL
	pinRayDL := moveRaysTable[kingSq][RAY_DL] // get the ray
	pinBlockersDL := pinRayDL & blockers      // get the blockers in the ray
	if pinBlockersDL != 0 {                   // if there are blockers
		pinBlockerSq1DL := pinBlockersDL.getLSBSq()                             // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyDL := bbReferenceArray[pinBlockerSq1DL] & frPieces // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyDL != 0 {                                       // if the piece is a friendly piece, continue the ray
			rayFromFriendlyDL := moveRaysTable[pinBlockerSq1DL][RAY_DL] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersDL := rayFromFriendlyDL & blockers   // check for further blockers
			if rayFromFriendlyBlockersDL != 0 {                         // if there are further blockers
				pinBlockerSq2DL := rayFromFriendlyBlockersDL.getLSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyDL := bbReferenceArray[pinBlockerSq2DL] & enBAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyDL != 0 {                                      // if a enemy pinning piece is found
					DLtURpins.setBit(pinBlockerSq1DL) // set the friendly piece found as pinned
				}
			}
		}
	}

	// UR
	pinRayUR := moveRaysTable[kingSq][RAY_UR] // get the ray
	pinBlockersUR := pinRayUR & blockers      // get the blockers in the ray
	if pinBlockersUR != 0 {                   // if there are blockers
		pinBlockerSq1UR := pinBlockersUR.getMSBSq()                             // get the square of the 1st blocker in the ray
		pinBlockerSq1FriendlyUR := bbReferenceArray[pinBlockerSq1UR] & frPieces // combine the 1st blocker square with friendly pieces
		if pinBlockerSq1FriendlyUR != 0 {                                       // if the piece is a friendly piece, continue the ray
			rayFromFriendlyUR := moveRaysTable[pinBlockerSq1UR][RAY_UR] // get the ray from the point of the friendly piece
			rayFromFriendlyBlockersUR := rayFromFriendlyUR & blockers   // check for further blockers
			if rayFromFriendlyBlockersUR != 0 {                         // if there are further blockers
				pinBlockerSq2UR := rayFromFriendlyBlockersUR.getMSBSq()             // get the square for the 2nd blocker in the ray
				pinBlockerSq2EnemyUR := bbReferenceArray[pinBlockerSq2UR] & enBAndQ // combine this with the specific enemy mask
				if pinBlockerSq2EnemyUR != 0 {                                      // if a enemy pinning piece is found
					DLtURpins.setBit(pinBlockerSq1UR) // set the friendly piece found as pinned
				}
			}
		}
	}

	return UDpins, LRpins, ULtDRpins, DLtURpins
}
