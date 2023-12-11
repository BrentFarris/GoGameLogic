package sudoku

type SquareSolver struct {
	square *SudokuSquare
	vals   []SquareValue
	len    int32
}
