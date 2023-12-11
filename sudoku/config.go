package sudoku

type SquareValue = int32

const (
	SquareValueNone SquareValue = iota
	SquareValueOne
	SquareValueTwo
	SquareValueThree
	SquareValueFour
	SquareValueFive
	SquareValueSix
	SquareValueSeven
	SquareValueEight
	SquareValueNine
)

const (
	retryAttemptCount = 50
)
