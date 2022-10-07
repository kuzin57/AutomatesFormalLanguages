package automata

import (
	customerrors "workspace/internal/errors"
)

type State struct {
	Number      int
	Transitions map[string][]int
	IsTerminal  bool
}

func mapStates(states []*state) ([]*State, error) {
	result := make([]*State, 0)
	statesNumbers := make(map[*state]int)
	for i, st := range states {
		statesNumbers[st] = i + 1
	}

	for i, state := range states {
		newState := &State{Transitions: make(map[string][]int), Number: i + 1, IsTerminal: state.isTerminal}
		for letter, statesTo := range state.next {
			for _, stateTo := range statesTo {
				_, ok := newState.Transitions[letter]
				switch ok {
				case true:
					newState.Transitions[letter] = append(newState.Transitions[letter], statesNumbers[stateTo])
				case false:
					newState.Transitions[letter] = []int{statesNumbers[stateTo]}
				}
			}
		}
		result = append(result, newState)
	}
	return result, nil
}

func addNewTransition(from *state, to *state, letter string) {
	needAdd := true
	for _, node := range from.next[letter] {
		if node == to {
			needAdd = false
		}
	}

	if needAdd {
		from.next[letter] = append(from.next[letter], to)
	}
}

func deleteEps(state *state, currentState *state, used *map[*state]bool) (err error) {
	(*used)[currentState] = true
	if currentState.isTerminal {
		state.isTerminal = true
	}

	for letter, statesTo := range currentState.next {
		switch {
		case letter == emptyWord:
			for _, stateTo := range statesTo {
				_, ok := (*used)[stateTo]
				if !ok {
					err = deleteEps(state, stateTo, used)
					if err != nil {
						return
					}
				}
			}

		default:
			if state != currentState {
				for _, v := range statesTo {
					addNewTransition(state, v, letter)
				}
			}
		}
	}
	return
}

func iterateStatesDeleteEps(currentState *state, used *map[*state]bool) (err error) {
	(*used)[currentState] = true

	visited := make(map[*state]bool)
	err = deleteEps(currentState, currentState, &visited)
	if err != nil {
		return err
	}

	for _, statesTo := range currentState.next {
		for _, stateTo := range statesTo {
			_, ok := (*used)[stateTo]
			if !ok {
				err = iterateStatesDeleteEps(stateTo, used)
				if err != nil {
					return
				}
			}
		}
	}

	return nil
}

func checkTerminal(start *state, used *map[*state]bool) *state {
	if start.isTerminal {
		return start
	}

	(*used)[start] = true
	for letter, statesTo := range start.next {
		if letter == emptyWord {
			for _, st := range statesTo {
				_, ok := (*used)[st]
				if !ok {
					return checkTerminal(st, used)
				}
			}
		}
	}
	return nil
}

func readWord(start *state, word string, index int, used *map[*state]int) error {
	if index == len(word) {
		setStates := make(map[*state]bool)
		res := checkTerminal(start, &setStates)
		if res != nil {
			return nil
		} else {
			return customerrors.ErrNoSuchWord
		}
	}

	(*used)[start] = index
	for letter, states := range start.next {
		switch letter {
		case string(word[index]):
			for _, state := range states {
				err := readWord(state, word, index+1, used)
				if err == nil {
					return nil
				}
			}
		case emptyWord:
			for _, state := range states {
				usedNumber, ok := (*used)[state]
				if !ok || usedNumber != index {
					err := readWord(state, word, index, used)
					if err == nil {
						return nil
					}
				}
			}
		}
	}
	return customerrors.ErrNoSuchWord
}

func proccessSelfTransitions(state *state) (string, error) {
	result := ""
	for letter, statesTo := range state.next {
		toDelete := -1
		for i, stateTo := range statesTo {
			if stateTo == state {
				if len(result) > 0 {
					result += "+"
				}
				result += letter
				toDelete = i
			}
		}

		if toDelete != -1 {
			state.next[letter] = append(
				state.next[letter][:toDelete],
				state.next[letter][(toDelete+1):]...,
			)
		}
	}

	return result, nil
}

func traversal(
	st *state,
	stock *state,
	arrStates *[]*state,
	classes *map[*state]int,
	used *map[*state]bool,
	executor Executor,
) error {
	(*used)[st] = true

	for _, val := range st.next {
		for _, v := range val {
			_, ok := (*used)[v]
			if !ok {
				err := traversal(v, stock, arrStates, classes, used, executor)
				if err != nil {
					return err
				}
			}
		}
	}

	if err := executor.Execute(st, stock, arrStates, classes); err != nil {
		return err
	}

	return nil
}

func removeDublicates(currentState *state) error {
	regulars := make(map[*state]string)

	for letter, statesTo := range currentState.next {
		for _, stateTo := range statesTo {
			_, ok := regulars[stateTo]
			if !ok {
				regulars[stateTo] = letter
			} else {
				regulars[stateTo] += "+"
				regulars[stateTo] += letter
			}
		}
	}

	currentState.next = make(map[string][]*state)
	for state, word := range regulars {
		_, ok := currentState.next[word]
		if !ok {
			currentState.next[word] = append(currentState.next[word], state)
		}
	}

	return nil
}
