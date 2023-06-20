package main

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------- Legal Move Generation -----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*
Legal move generation works in the following stages:
1. Generate legal moves of each piece (captures, quiet moves and promotions).
2. Generate en-passant moves.
3. Generate castling moves.

To generate legal moves of each piece, we generally follow these steps:
1.1 Get the pseudo-legal moves of each piece (all moves filtered out for blockers).
1.2 Mask the moves with a king in check mask (if the king is in check, only certain squares can remove the check).
1.3 Mask the moves with a pin mask (pieces that are pinned can only move along certain rays).
The final result is only legal moves.
*/

// generate all the legal moves for a position
func (pos *Position) generateLegalMoves() {

	pos.logTime.allLogTypes[LOG_MOVE_GEN_TOTAL].start()

	// ------------------------------------------------- Setup ---------------------------------------------
	// reset the moves counter
	pos.totalMovesCounter = 0
	pos.threatMovesCounter = 0
	pos.quietMovesCounter = 0

	// set a mobility bonus counter for moves we want to give a mobility bonus to
	mobilityBonusCounter := 0

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

	piecesAttKingBB, piecesAndSqAttKingBB := getAttacksOnKing(
		kingSq, pos.piecesAll[SIDE_BOTH], pos.piecesAll[enSide], pos.piecesAll[frSide], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
		pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.isWhiteTurn)

	// ------------------------------------------------- King Moves ---------------------------------------------
	// get king pseudo-legal moves
	// filter out attacked squares
	// the rest are legal moves

	// get the pseudo legal moves of the piece on that square
	kingMovesPseudo := getKingMovesPseudo(kingSq)

	// mask out moves to friendly pieces
	kingMovesPseudo &= ^pos.piecesAll[frSide]

	// check the remaining moves for legality
	for kingMovesPseudo != 0 {

		// get the next move square
		nextMoveSq := kingMovesPseudo.popBitGetSq()

		// king can only move to non-threatened squares
		if !isSqAttacked(
			nextMoveSq, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
			pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
			pos.isWhiteTurn) {
			if pos.piecesAll[enSide]&bbReferenceArray[nextMoveSq] != 0 { // capture
				pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(kingSq, nextMoveSq, PIECE_KING, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.threatMovesCounter += 1
			} else { // quiet move
				pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(kingSq, nextMoveSq, PIECE_KING, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.quietMovesCounter += 1
			}
		}
	}

	// ------------------------------------------------- Pins ---------------------------------------------
	// get pinned pieces bitboards

	pinsUD, pinsLR, pinsULtDR, pinsDLtUR := getPinnedPieces(
		kingSq, pos.piecesAll[SIDE_BOTH], pos.piecesAll[frSide], pos.pieces[frSide][PIECE_PAWN], pos.piecesAll[enSide], pos.pieces[enSide][PIECE_QUEEN],
		pos.pieces[enSide][PIECE_ROOK], pos.pieces[enSide][PIECE_BISHOP], pos.enPassantTargetBB, pos.isWhiteTurn)
	pinsCombined := pinsUD | pinsLR | pinsULtDR | pinsDLtUR

	// ------------------------------------------------- King Attacks ---------------------------------------------
	// count the attacks on the king
	kingChecks := piecesAttKingBB.countBits()

	// store the number of checks for detecting checkmate/stalemate later
	pos.kingChecks = kingChecks

	// if the number of checks is two, no other moves are possible (already generated king moves above)
	if kingChecks >= 2 {
		return
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

	// while there are pieces left
	for frQueens != 0 {

		// get the square of the piece
		nextQueenOriginSq := frQueens.popBitGetSq()

		// get the pseudo legal moves of the piece on that square
		nextQueenMoves := getQueenMovesPseudo(nextQueenOriginSq, pos.piecesAll[SIDE_BOTH])

		// mask out moves to friendly pieces
		nextQueenMoves &= ^pos.piecesAll[frSide]

		// mask the moves with the king check mask
		nextQueenMoves &= kingInCheckMask

		// if pinned, mask the moves with the pins mask
		if pinsCombined != 0 {
			if bbReferenceArray[nextQueenOriginSq]&pinsUD != 0 {
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_UD]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsLR != 0 {
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_LR]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsULtDR != 0 {
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextQueenOriginSq]&pinsDLtUR != 0 {
				nextQueenMoves &= movePinnedMasksTable[nextQueenOriginSq][PIN_DLtUR]
			}
		}

		// finally save the remaining moves
		for nextQueenMoves != 0 {
			nextQueenTargetSq := nextQueenMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextQueenTargetSq] != 0 { // capture
				pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextQueenOriginSq, nextQueenTargetSq, PIECE_QUEEN, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.threatMovesCounter += 1
			} else { // quiet move
				pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(nextQueenOriginSq, nextQueenTargetSq, PIECE_QUEEN, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.quietMovesCounter += 1
			}
		}
	}

	// ------------------------------------------------- Rook Moves ---------------------------------------------
	// rooks

	// while there are pieces left
	for frRooks != 0 {

		// get the square of the piece
		nextRookOriginSq := frRooks.popBitGetSq()

		// get the pseudo legal moves of the piece on that square
		nextRookMoves := getRookMovesPseudo(nextRookOriginSq, pos.piecesAll[SIDE_BOTH])

		// mask out moves to friendly pieces
		nextRookMoves &= ^pos.piecesAll[frSide]

		// mask the moves with the king check mask
		nextRookMoves &= kingInCheckMask

		// if pinned, mask the moves with the pins mask
		if pinsCombined != 0 {
			if bbReferenceArray[nextRookOriginSq]&pinsUD != 0 {
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_UD]
			} else if bbReferenceArray[nextRookOriginSq]&pinsLR != 0 {
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_LR]
			} else if bbReferenceArray[nextRookOriginSq]&pinsULtDR != 0 {
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextRookOriginSq]&pinsDLtUR != 0 {
				nextRookMoves &= movePinnedMasksTable[nextRookOriginSq][PIN_DLtUR]
			}
		}

		// finally save the remaining moves
		for nextRookMoves != 0 {
			nextRookTargetSq := nextRookMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextRookTargetSq] != 0 { // capture
				pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextRookOriginSq, nextRookTargetSq, PIECE_ROOK, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.threatMovesCounter += 1
				mobilityBonusCounter++
			} else { // quiet move
				pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(nextRookOriginSq, nextRookTargetSq, PIECE_ROOK, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.quietMovesCounter += 1
				mobilityBonusCounter++
			}
		}
	}

	// ------------------------------------------------- Bishop Moves ---------------------------------------------
	// bishops

	// while there are pieces left
	for frBishops != 0 {

		// get the square of the piece
		nextBishopOriginSq := frBishops.popBitGetSq()

		// get the pseudo legal moves of the piece on that square
		nextBishopMoves := getBishopMovesPseudo(nextBishopOriginSq, pos.piecesAll[SIDE_BOTH])

		// mask out moves to friendly pieces
		nextBishopMoves &= ^pos.piecesAll[frSide]

		// mask the moves with the king check mask
		nextBishopMoves &= kingInCheckMask

		// if pinned, mask the moves with the pins mask
		if pinsCombined != 0 {
			if bbReferenceArray[nextBishopOriginSq]&pinsUD != 0 {
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_UD]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsLR != 0 {
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_LR]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsULtDR != 0 {
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextBishopOriginSq]&pinsDLtUR != 0 {
				nextBishopMoves &= movePinnedMasksTable[nextBishopOriginSq][PIN_DLtUR]
			}
		}

		// finally save the remaining moves
		for nextBishopMoves != 0 {
			nextBishopTargetSq := nextBishopMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextBishopTargetSq] != 0 { // capture
				pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextBishopOriginSq, nextBishopTargetSq, PIECE_BISHOP, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.threatMovesCounter += 1
				mobilityBonusCounter++
			} else { // quiet move
				pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(nextBishopOriginSq, nextBishopTargetSq, PIECE_BISHOP, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.quietMovesCounter += 1
				mobilityBonusCounter++
			}
		}
	}

	// ------------------------------------------------- Knight Moves ---------------------------------------------
	// knights

	// while there are pieces left
	for frKnights != 0 {

		// get the square of the piece
		nextKnightOriginSq := frKnights.popBitGetSq()

		// get the pseudo legal moves of the piece on that square
		nextKnightMoves := getKnightMovesPseudo(nextKnightOriginSq)

		// mask out moves to friendly pieces
		nextKnightMoves &= ^pos.piecesAll[frSide]

		// mask the moves with the king check mask
		nextKnightMoves &= kingInCheckMask

		// if pinned, mask the moves with the pins mask
		if pinsCombined != 0 {
			if bbReferenceArray[nextKnightOriginSq]&pinsUD != 0 {
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_UD]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsLR != 0 {
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_LR]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsULtDR != 0 {
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextKnightOriginSq]&pinsDLtUR != 0 {
				nextKnightMoves &= movePinnedMasksTable[nextKnightOriginSq][PIN_DLtUR]
			}
		}

		// finally save the remaining moves
		for nextKnightMoves != 0 {
			nextKnightTargetSq := nextKnightMoves.popBitGetSq()
			if pos.piecesAll[enSide]&bbReferenceArray[nextKnightTargetSq] != 0 { // capture
				pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextKnightOriginSq, nextKnightTargetSq, PIECE_KNIGHT, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.threatMovesCounter += 1
				mobilityBonusCounter++
			} else { // quiet move
				pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(nextKnightOriginSq, nextKnightTargetSq, PIECE_KNIGHT, MOVE_TYPE_QUIET, PROMOTION_NONE)
				pos.totalMovesCounter += 1
				pos.quietMovesCounter += 1
				mobilityBonusCounter++
			}
		}
	}

	// ------------------------------------------------- Pawn Moves ---------------------------------------------
	// pawns

	// while there are pieces left
	for frPawns != 0 {

		// get the square of the piece
		nextPawnOriginSq := frPawns.popBitGetSq()

		// get the pseudo legal moves of the piece on that square
		var nextPawnMoves Bitboard
		if pos.isWhiteTurn {
			nextPawnMoves = getPawnMovesWhitePseudo(
				nextPawnOriginSq,
				pos.piecesAll[SIDE_BOTH],
				pos.piecesAll[enSide])
		} else {
			nextPawnMoves = getPawnMovesBlackPseudo(
				nextPawnOriginSq,
				pos.piecesAll[SIDE_BOTH],
				pos.piecesAll[enSide])
		}

		// mask out moves to friendly pieces
		nextPawnMoves &= ^pos.piecesAll[frSide]

		// mask the moves with the king check mask
		nextPawnMoves &= kingInCheckMask

		// if pinned, mask the moves with the pins mask
		if pinsCombined != 0 {
			if bbReferenceArray[nextPawnOriginSq]&pinsUD != 0 {
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_UD]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsLR != 0 {
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_LR]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsULtDR != 0 {
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_ULtDR]
			} else if bbReferenceArray[nextPawnOriginSq]&pinsDLtUR != 0 {
				nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_DLtUR]
			}
		}

		// finally save the remaining moves
		for nextPawnMoves != 0 {
			nextPawnTargetSq := nextPawnMoves.popBitGetSq()

			if pos.piecesAll[enSide]&bbReferenceArray[nextPawnTargetSq] != 0 { // capture

				if nextPawnTargetSq >= 56 || nextPawnTargetSq <= 7 { // if there is a promotion
					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_QUEEN)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_ROOK)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_KNIGHT)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_BISHOP)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

				} else { // if there is not a promotion
					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_CAPTURE, PROMOTION_NONE)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1
				}

			} else { // quiet move

				if nextPawnTargetSq >= 56 || nextPawnTargetSq <= 7 { // if there is a promotion
					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_QUEEN)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_ROOK)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_KNIGHT)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_BISHOP)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1

				} else { // if there is not a promotion
					pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnTargetSq, PIECE_PAWN, MOVE_TYPE_QUIET, PROMOTION_NONE)
					pos.totalMovesCounter += 1
					pos.quietMovesCounter += 1
				}
			}
		}
	}

	// ------------------------------------------------- En-Passant Moves ---------------------------------------------
	// captures en-passant (includes checking for en-passant pawn on pin bitmask)
	// separate check based on the 2 pawn squares that can attack the en-passant target

	// special rule start
	// the king in check mask does not include pawns checking the king that can be captured en-passant
	// therefore, for en-passant captures only, add it to the check mask if a pawn that is giving check can be captured en-passant

	// create a copy of the king check mask to adjust if needed
	enPassantKingCheckMask := kingInCheckMask

	// if there is a possible en-passant capture
	if pos.enPassantTargetBB != 0 {

		// get the en-passant bitboard and get its square
		enPassantTarget := pos.enPassantTargetBB
		enPassantTargetSq := enPassantTarget.popBitGetSq()

		// get the enemy pawns are giving check
		pawnsCheckingKing := movePawnsAttackingKingMasks[kingSq][enSide] & pos.pieces[enSide][PIECE_PAWN]

		// for each of them
		for pawnsCheckingKing != 0 {

			// get the square of the pawn checking
			nextCheckerSq := pawnsCheckingKing.popBitGetSq()

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

	// mask with allowable moves when the king is in check
	enPasTargetMasked := pos.enPassantTargetBB & enPassantKingCheckMask

	// if an en-passant capture can be made
	if enPasTargetMasked != 0 {

		// get the bitboard of the pawn that can be captured
		var enPassantCapturedPieceSqBB Bitboard
		if pos.isWhiteTurn {
			enPassantCapturedPieceSqBB = enPasTargetMasked << 8
		} else {
			enPassantCapturedPieceSqBB = enPasTargetMasked >> 8
		}

		// get the square of the pawn that can be captured
		enPassantCapturedPieceSq := enPassantCapturedPieceSqBB.popBitGetSq()

		// only if the en passant captured pawn is not pinned, allow the en-passant
		if bbReferenceArray[enPassantCapturedPieceSq]&pinsCombined == 0 {

			// which pawns can capture
			pawnsCanCapture := pos.pieces[frSide][PIECE_PAWN] & movePawnsAttackingKingMasks[enPasTargetMasked.popBitGetSq()][frSide]

			// if there are pawns that can capture
			for pawnsCanCapture != 0 {

				// get the origin of the pawn that can capture
				nextPawnOriginSq := pawnsCanCapture.popBitGetSq()

				// get the target of the pawn that can capture
				nextPawnMoves := pos.enPassantTargetBB

				// now need to check if the CAPTURING pawn is pinned (already checked for CAPTURED pawn pins above)
				// if pinned, mask the moves with the pins mask
				if bbReferenceArray[nextPawnOriginSq]&pinsUD != 0 {
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_UD]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsLR != 0 {
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_LR]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsULtDR != 0 {
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_ULtDR]
				} else if bbReferenceArray[nextPawnOriginSq]&pinsDLtUR != 0 {
					nextPawnMoves &= movePinnedMasksTable[nextPawnOriginSq][PIN_DLtUR]
				}

				if nextPawnMoves != 0 { // if there are still moves remaining
					pos.threatMoves[pos.threatMovesCounter] = getEncodedMove(nextPawnOriginSq, nextPawnMoves.popBitGetSq(), PIECE_PAWN, MOVE_TYPE_EN_PASSANT, PROMOTION_NONE)
					pos.totalMovesCounter += 1
					pos.threatMovesCounter += 1
				}
			}
		}
	}

	// ------------------------------------------------- Castling Moves ---------------------------------------------
	// castling moves

	// if we are allowed to generate castling moves
	if generateCastlingMoves {

		// white castling
		if pos.isWhiteTurn {

			// if castling is available
			if pos.castlingRights[CASTLE_WHITE_KINGSIDE] {

				// get the mask
				castlingSquares := moveCastlingIsClearMasks[CASTLE_WHITE_KINGSIDE]

				// if there are no pieces on those squares
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]

				// check if those squares are attacked
				if castlingMasked == 0 {
					AttSq1 := 5
					AttSq2 := 6
					if !isSqAttacked(
						AttSq1, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(4, 6, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.totalMovesCounter += 1
						pos.quietMovesCounter += 1
					}
				}
			}

			// if castling is available
			if pos.castlingRights[CASTLE_WHITE_QUEENSIDE] {

				// get the mask
				castlingSquares := moveCastlingIsClearMasks[CASTLE_WHITE_QUEENSIDE]

				// if there are no pieces on those squares
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]

				// check if those squares are attacked
				if castlingMasked == 0 {
					AttSq1 := 3
					AttSq2 := 2
					if !isSqAttacked(
						AttSq1, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(4, 2, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.totalMovesCounter += 1
						pos.quietMovesCounter += 1
					}
				}
			}

			// black castling
		} else {

			// if castling is available
			if pos.castlingRights[CASTLE_BLACK_KINGSIDE] {

				// get the mask
				castlingSquares := moveCastlingIsClearMasks[CASTLE_BLACK_KINGSIDE]

				// if there are no pieces on those squares
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]

				// check if those squares are attacked
				if castlingMasked == 0 {
					AttSq1 := 61
					AttSq2 := 62
					if !isSqAttacked(
						AttSq1, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(60, 62, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.totalMovesCounter += 1
						pos.quietMovesCounter += 1
					}
				}
			}

			// if castling is available
			if pos.castlingRights[CASTLE_BLACK_QUEENSIDE] {

				// get the mask
				castlingSquares := moveCastlingIsClearMasks[CASTLE_BLACK_QUEENSIDE]

				// if there are no pieces on those squares
				castlingMasked := castlingSquares & pos.piecesAll[SIDE_BOTH]

				// check if those squares are attacked
				if castlingMasked == 0 {
					AttSq1 := 59
					AttSq2 := 58
					if !isSqAttacked(
						AttSq1, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) && !isSqAttacked(
						AttSq2, pos.piecesAll[SIDE_BOTH], pos.pieces[frSide][PIECE_KING], pos.pieces[enSide][PIECE_QUEEN], pos.pieces[enSide][PIECE_ROOK],
						pos.pieces[enSide][PIECE_KNIGHT], pos.pieces[enSide][PIECE_BISHOP], pos.pieces[enSide][PIECE_PAWN], pos.pieces[enSide][PIECE_KING],
						pos.isWhiteTurn) {

						pos.quietMoves[pos.quietMovesCounter] = getEncodedMove(60, 58, PIECE_KING, MOVE_TYPE_CASTLE, PROMOTION_NONE)
						pos.totalMovesCounter += 1
						pos.quietMovesCounter += 1
					}
				}
			}
		}
	}

	// ------------------------------------------------- Eval Mobility ---------------------------------------------
	// before we end the function, we store the mobility bonus counter
	// we don't update when in check to remove wild fluctuations
	if kingChecks == 0 {
		if pos.isWhiteTurn {
			pos.evalWhiteMobility = mobilityBonusCounter
		} else {
			pos.evalBlackMobility = mobilityBonusCounter
		}
	}

	pos.logTime.allLogTypes[LOG_MOVE_GEN_TOTAL].stop()
}

// --------------------------------------------------------------------------------------------------------------------
// ---------------- Generate Pseudo-Legal Moves (excluding castling, promoting, pinned pieces and checks) -------------
// --------------------------------------------------------------------------------------------------------------------
// gets the moves of a piece from a square filtered for blockers

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
