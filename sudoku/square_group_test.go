package sudoku

import (
	"testing"
)

func TestSquareGroupFilled(t *testing.T) {
	g := make(SquareGroup, 9)
	for i := range g {
		g[i] = &SudokuSquare{value: SquareValueNone}
	}
	for i := range g {
		if g.filled() {
			t.Errorf("Expected group to not be filled")
		}
		g[i].SetValue(SquareValue(i + 1))
	}
	if !g.filled() {
		t.Errorf("Expected group to be filled")
	}
	g[1].value = SquareValueNine
	if g.filled() {
		t.Errorf("Expected group to not be filled")
	}
}
