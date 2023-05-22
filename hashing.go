package main

import (
	"math/rand"
	"time"
)

// a hash is the same as a bitboard (uint64)
// initialize

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------- Init Hash Tables ------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
/*

Initialize unique hash values for unique position identifiers:
1. Every piece type on every square
2. Side to move
3. Castling rights
4. En-Passant Target

*/

// pre-initialized hash tables
var hashTablePieces [64][2][6]Bitboard // for each square, for each side, for each piece type
var hashTableCastling [4]Bitboard      // KQkq
var hashTableSideToMove [1]Bitboard    // if white to move
var hashTableEnPassant [64]Bitboard    // en-passant for the 3rd and 6th row (unused indexes are just for easier lookup later)
var startingHash Bitboard              // random starting hash to work off from

// hash collisions checker
var hashCollisionsStack []Bitboard // all previously generated random numbers to check for collisions

// function to get a random number already checked for collisions
func getRandomUint64() Bitboard {
	var newRandNo Bitboard

	// try to get a new unique number
	for {
		// get a random number
		newRandNo = Bitboard(rand.Uint64())

		// check if the number is not zero or not full
		if newRandNo != emptyBB && newRandNo != fullBB {

			// check that the number has not been previously used
			foundCollision := false
			for _, previousNum := range hashCollisionsStack {
				if newRandNo == previousNum {
					foundCollision = true
				}
			}

			// if the number is unique, add it to the used stack and return it
			if !foundCollision {
				hashCollisionsStack = append(hashCollisionsStack, newRandNo)
				return newRandNo
			}
		}
	}
}

// function to initialize the hash tables
func initHashTables() {

	// seed the random number generator
	rand.Seed(time.Now().Unix())

	// get numbers for each square, side and piece
	for sq := 0; sq < 64; sq++ {
		for side := 0; side < 2; side++ {
			for piece := 0; piece < 6; piece++ {
				hashTablePieces[sq][side][piece] = getRandomUint64()
			}
		}
	}

	// get numbers for castling
	for castlingSide := 0; castlingSide < 4; castlingSide++ {
		hashTableCastling[castlingSide] = getRandomUint64()
	}

	// get a number for the side to move
	hashTableSideToMove[0] = getRandomUint64()

	// get en-passant numbers
	for sq := 0; sq < 64; sq++ {
		hashTableEnPassant[sq] = getRandomUint64()
	}

	// get the starting hash
	startingHash = getRandomUint64()
}

// --------------------------------------------------------------------------------------------------------------------
// ------------------------------------------------- Hash Position ----------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------
// hash the current position and store the value

func (pos *Position) hashPosAndStore() {

	start_time := time.Now()

	// start the hash
	positionHash := startingHash

	// hash all the pieces
	for side := 0; side < 2; side++ {
		for piece := 0; piece < 6; piece++ {
			piecesCopy := pos.pieces[side][piece]
			for piecesCopy != 0 {
				nextSq := piecesCopy.popBitGetSq()
				positionHash ^= hashTablePieces[nextSq][side][piece]
			}
		}
	}

	// hash the castling rights
	for i := 0; i < 4; i++ {
		if pos.castlingRights[i] {
			positionHash ^= hashTableCastling[i]
		}
	}

	// hash the side to move
	if pos.isWhiteTurn {
		positionHash ^= hashTableSideToMove[0]
	}

	// hash the en-passant square
	if pos.enPassantTargetBB != 0 {
		enPBB := pos.enPassantTargetBB
		enPSq := enPBB.popBitGetSq()
		positionHash ^= hashTableEnPassant[enPSq]
	}

	duration_time := time.Since(start_time).Nanoseconds()
	pos.logOther.allLogTypes[LOG_HASHING].addTime(int(duration_time))

	pos.hashOfPos = positionHash
}
