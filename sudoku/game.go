package sudoku

import "math/rand"

func NewGame(removeMin, removeMax int) SudokuBoard {
	rnd := rand.New(rand.NewSource(0))
	board := GenerateBoard()
	removeCount := rnd.Intn(removeMax-removeMin) + removeMin
	removeIdxs := [81]int{}
	removeIdxSkip := [81]int{}
	removeValues := [81]SquareValue{}
	options := [81]int{}
	resetTriggered := false
	resetTimes := 0
	nukeCount := 0
	removeCounter := removeCount
	for removeCounter > 0 {
		optionsLen := 0
		for i := 0; i < 81; i++ {
			if removeIdxSkip[i] == 0 && removeIdxs[i] == 0 {
				options[optionsLen] = i
				optionsLen++
			}
		}
		if (removeCounter < 5 && optionsLen < 15) || optionsLen == 0 {
			clear(removeIdxSkip[:])
			if resetTriggered {
				optionsLen = 0
				for i := 0; i < 81; i++ {
					if removeIdxs[i] == 1 {
						options[optionsLen] = i
						optionsLen++
					}
				}
				resetCount := 5
				resetTimes++
				nukeCount++
				if resetTimes < 2 {
					rnd.Shuffle(optionsLen, func(i, j int) {
						options[i], options[j] = options[j], options[i]
					})
					nukeCount--
				} else if nukeCount == 2 {
					board.generate()
					resetTimes = 0
					nukeCount = 0
					removeCount = rnd.Intn(removeMax-removeMin) + removeMin
					removeCounter = removeCount
					for i := 0; i < 81; i++ {
						removeIdxs[i] = 0
						removeIdxSkip[i] = 0
					}
					// TODO:  Once we get here we are screwed?
					continue
				} else {
					// Nuclear option
					resetCount = optionsLen
					resetTimes = 0
				}
				for i := 0; i < resetCount; i++ {
					v := options[i]
					board.squares[v].SetValue(removeValues[v])
					removeIdxs[v] = 0
					//removeIdxSkip[v] = 1;
					removeCounter++
				}
			}
			resetTriggered = !resetTriggered
			continue
		}
		idx := rnd.Intn(optionsLen)
		idx = options[idx]
		value := board.squares[idx].value
		board.squares[idx].SetValue(SquareValueNone)
		if board.SolutionCount() == 1 {
			removeValues[idx] = value
			removeIdxs[idx] = 1
			removeCounter--
		} else {
			board.squares[idx].SetValue(value)
			removeIdxSkip[idx] = 1
		}
	}
	return board
}
