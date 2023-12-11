package sudoku

type SudokuSquare struct {
	hint  SquareValue
	value SquareValue
	index int32
}

func (s SudokuSquare) Index() int32                { return s.index }
func (s SudokuSquare) Value() SquareValue          { return s.value }
func (s SudokuSquare) Hint() SquareValue           { return s.hint }
func (s *SudokuSquare) SetValue(value SquareValue) { s.value = value }
func (s *SudokuSquare) SetHint(hint SquareValue)   { s.hint = hint }
