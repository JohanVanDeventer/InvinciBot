package main

// --------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------- Index Functions------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

func sqFromRowAndCol(row int, col int) int {
	return row*8 + col
}

func rowAndColFromSq(sq int) (int, int) {
	row := sq / 8
	col := sq % 8
	return row, col
}

// converts for example "c3" to sq 18
func getSqFromString(inputStr string) int {
	var colInt int = -1
	var rowInt int = -1

	switch inputStr[:1] {
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

	switch inputStr[1:] {
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

	return sqFromRowAndCol(rowInt, colInt)
}

// converts for example sq 18 to "c3"
func getStringFromSq(inputSq int) string {
	var rowStr string
	var colStr string

	row, col := rowAndColFromSq(inputSq)

	switch row {
	case 0:
		rowStr = "1"
	case 1:
		rowStr = "2"
	case 2:
		rowStr = "3"
	case 3:
		rowStr = "4"
	case 4:
		rowStr = "5"
	case 5:
		rowStr = "6"
	case 6:
		rowStr = "7"
	case 7:
		rowStr = "8"
	}

	switch col {
	case 0:
		colStr = "a"
	case 1:
		colStr = "b"
	case 2:
		colStr = "c"
	case 3:
		colStr = "d"
	case 4:
		colStr = "e"
	case 5:
		colStr = "f"
	case 6:
		colStr = "g"
	case 7:
		colStr = "h"
	}

	return colStr + rowStr
}

// converts for example queen promotion to "q"
func getPromotionStringFromType(inputPromotion int) string {
	switch inputPromotion {
	case PROMOTION_NONE:
		return ""
	case PROMOTION_QUEEN:
		return "q"
	case PROMOTION_ROOK:
		return "r"
	case PROMOTION_KNIGHT:
		return "n"
	case PROMOTION_BISHOP:
		return "b"
	}
	return ""
}
