package sudoku

import "sort"

type SquareGroup []*SudokuSquare

func (g *SquareGroup) addSquare(s *SudokuSquare) {
	*g = append(*g, s)
}

func (g SquareGroup) filled() bool {
	success := len(g) > 0
	if success {
		values := make([]int, len(g))
		for i := range g {
			values[i] = int(g[i].value)
		}
		sort.Ints(values)
		success = values[0] != int(SquareValueNone)
		for i := 0; i < len(g) && success; i++ {
			success = values[i] == i+1
		}
	}
	return success
}

func (g SquareGroup) contains(value SquareValue) bool {
	found := false
	for i := 0; i < len(g) && !found; i++ {
		found = g[i].value == value
	}
	return found
}
