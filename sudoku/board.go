package sudoku

import (
	"math"
	"math/rand"
	"strings"
)

type SudokuBoard struct {
	rand       *rand.Rand
	squares    [91]SudokuSquare
	groups     [9]SquareGroup
	width      int32
	groupWidth int32
}

func (b SudokuBoard) Width() int32 { return b.width }
func (b *SudokuBoard) Square(index int32) *SudokuSquare {
	return &b.squares[index]
}

func NewBoard() SudokuBoard {
	board := SudokuBoard{
		rand:       rand.New(rand.NewSource(0)),
		width:      9,
		groupWidth: 3,
	}
	for i := range board.groups {
		board.groups[i] = make(SquareGroup, 0, 9)
	}
	section := int32(math.Sqrt(float64(board.width)))
	for y := int32(0); y < board.width; y++ {
		for x, i := int32(0), 0; x < board.width; x++ {
			s := SudokuSquare{
				hint:  SquareValueNone,
				value: SquareValueNone,
				index: int32(i),
			}
			board.squares[s.index] = s
			left := x / section
			top := y / section
			board.groups[top*board.groupWidth+left].addSquare(&board.squares[s.index])
			i++
		}
	}
	return board
}

func (b *SudokuBoard) PossibleValues(index int32) []SquareValue {
	valLen := SquareValueNine + 1
	vals := make([]SquareValue, 0, valLen)
	for i := SquareValue(SquareValueNone); i <= SquareValueNine; i++ {
		vals = append(vals, i)
	}
	x := index % b.width
	y := index / b.width
	gx := x / b.groupWidth
	gy := y / b.groupWidth
	g := b.groups[gy*b.groupWidth+gx]
	for i := range g {
		vals[g[i].value] = SquareValueNone
	}
	size := b.width
	for i := int32(0); i < size; i++ {
		s := b.squares[y*b.width+i]
		vals[s.value] = SquareValueNone
	}
	for i := int32(0); i < size; i++ {
		s := b.squares[i*b.width+x]
		vals[s.value] = SquareValueNone
	}
	for i := range vals {
		if vals[i] == SquareValueNone {
			vals = append(vals[:i], vals[i+1:]...)
			valLen--
			i--
		}
	}
	return vals[:valLen]
}

func GenerateBoard() SudokuBoard {
	board := NewBoard()
	board.generate()
	return board
}

func (b *SudokuBoard) generate() {
	b.reset()
	for !b.run() {
	}
}

func (b *SudokuBoard) reset() {
	for i := range b.squares {
		b.squares[i].value = SquareValueNone
	}
	b.setTopRow()
	b.setGroup(0)
}

func (b *SudokuBoard) setTopRow() {
	for x := int32(0); x < b.width; x++ {
		vals := b.PossibleValues(x)
		b.squares[x].value = vals[b.rand.Intn(len(vals))]
	}
}

func (b *SudokuBoard) setGroup(index int32) bool {
	failed := false
	var undo [9]*SudokuSquare
	undoLen := 0
	group := b.groups[index]
	for i := range group {
		s := group[i]
		if s.value == SquareValueNone {
			vals := b.PossibleValues(s.index)
			if len(vals) > 0 {
				s.SetValue(vals[b.rand.Intn(len(vals))])
				undo[undoLen] = s
				undoLen++
			} else {
				failed = true
				break
			}
		}
	}
	if failed {
		for i := 0; i < undoLen; i++ {
			undo[i].SetValue(SquareValueNone)
		}
	}
	return !failed
}

func (b *SudokuBoard) run() bool {
	for g := int32(0); g < b.width; g++ {
		for i := 0; g <= retryAttemptCount && !b.setGroup(g); i++ {
			if g == retryAttemptCount {
				b.reset()
				return false
			}
		}
	}
	return true
}

func (b *SudokuBoard) CanPlace(index int32, value SquareValue) bool {
	row := index / b.width
	col := index % b.width
	section := (int32)(math.Sqrt(float64(b.width)))
	left := col / section
	top := row / section
	group := top*section + left
	return !b.rowContains(row, value) &&
		!b.columnContains(col, value) &&
		!b.groups[group].contains(value)
}

func (b *SudokuBoard) Validate() bool {
	for i := 0; i < int(b.width); i++ {
		if !b.ValidateRow(i) || !b.ValidateColumn(i) || !b.groups[i].filled() {
			return false
		}
	}
	return true
}

func (b *SudokuBoard) validateLine(idx int, toSquareIdx func(idx, itr int32) int32) bool {
	count := make([]int32, SquareValueNine+1)
	for i := SquareValueNone; i <= SquareValueNine; i++ {
		count[i] = i
	}
	for i := int32(0); i < b.width; i++ {
		count[b.squares[toSquareIdx(int32(idx), i)].value] = 0
	}
	sum := int32(0)
	for i := SquareValueNone; i <= SquareValueNine; i++ {
		sum += count[i]
	}
	return sum == 0
}

func (b *SudokuBoard) ValidateRow(idx int) bool {
	return b.validateLine(idx, func(idx, itr int32) int32 {
		return idx*b.width + itr
	})
}

func (b *SudokuBoard) ValidateColumn(idx int) bool {
	return b.validateLine(idx, func(idx, itr int32) int32 {
		return itr*b.width + idx
	})
}

func (b *SudokuBoard) rowContains(row int32, value SquareValue) bool {
	found := false
	for i := int32(0); i < b.width && !found; i++ {
		found = b.squares[int32(row)*b.width+i].value == value
	}
	return found
}

func (b *SudokuBoard) columnContains(col int32, value SquareValue) bool {
	found := false
	for i := int32(0); i < b.width && !found; i++ {
		found = b.squares[i*b.width+int32(col)].value == value
	}
	return found
}

func (b *SudokuBoard) SolutionCount() int {
	solutions := 0
	// Get a list of all the blank tiles
	blanks := make([]SquareSolver, 0, 81)
	for i := 0; i < len(b.squares); i++ {
		s := &b.squares[i]
		if s.value == SquareValueNone {
			ss := SquareSolver{
				square: s,
				vals:   b.PossibleValues(s.index),
			}
			blanks = append(blanks, ss)
		}
	}
	solutions = b.countSolutions(blanks)
	return solutions
}

func (b *SudokuBoard) countSolutions(blanks []SquareSolver) int {
	solveCount := 0
	for i := 0; i < len(blanks) && solveCount <= 1; i++ {
		s := blanks[i]
		if b.CanPlace(s.square.index, s.vals[i]) {
			s.square.SetValue(s.vals[i])
			if len(blanks) > 1 {
				solveCount += b.countSolutions(blanks[1:])
			} else {
				solveCount++
			}
			s.square.SetValue(SquareValueNone)
		}
	}
	return solveCount
}

func (b *SudokuBoard) ToAscii() string {
	sb := strings.Builder{}
	for i := int32(0); i < int32(len(b.squares)); i++ {
		s := b.squares[i]
		if s.value == SquareValueNone {
			sb.WriteRune(' ')
		} else {
			sb.WriteRune(rune(s.value + '0'))
		}
		if i%b.width == b.width-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
