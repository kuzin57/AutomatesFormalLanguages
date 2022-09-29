package automate

import (
	customerrors "workspace/internal/errors"
)

type State struct {
	Number      int
	Transitions map[rune][]int
	IsTerminal  bool
}

func mapStates(start *state, used *map[*state]int, ans *[]*State, counter *int) error {
	(*counter)++
	(*used)[start] = *counter

	newState := &State{Number: *counter, Transitions: make(map[rune][]int), IsTerminal: start.isTerm}
	(*ans) = append((*ans), newState)

	for key, val := range start.next {
		newState.Transitions[key] = make([]int, 0)
		for _, v := range val {
			_, ok := (*used)[v]
			if !ok {
				e := mapStates(v, used, ans, counter)
				if e != nil {
					return e
				}
			}
			newState.Transitions[key] = append(newState.Transitions[key], (*used)[v])
		}
	}
	return nil
}

func deleteEps(st *state, cur *state, used *map[*state]bool) (err error) {
	(*used)[cur] = true
	if cur.isTerm {
		st.isTerm = true
	}

	for key, val := range cur.next {
		switch {
		case key == emptyWord:
			for _, v := range val {
				_, ok := (*used)[v]
				if !ok {
					err = deleteEps(st, v, used)
					if err != nil {
						return
					}
				}
			}

		default:
			if st != cur {
				st.next[key] = append(st.next[key], val...)
			}
		}
	}
	return
}

func iterateStatesDeleteEps(st *state, used *map[*state]bool) (err error) {
	(*used)[st] = true

	visited := make(map[*state]bool)
	err = deleteEps(st, st, &visited)
	if err != nil {
		return err
	}

	for _, val := range st.next {
		for _, v := range val {
			_, ok := (*used)[v]
			if !ok {
				err = iterateStatesDeleteEps(v, used)
				if err != nil {
					return
				}
			}
		}
	}

	return nil
}

func checkTerminal(start *state, used *map[*state]bool) *state {
	if start.isTerm {
		return start
	}

	(*used)[start] = true
	for key, val := range start.next {
		if key == emptyWord {
			for _, st := range val {
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
		m := make(map[*state]bool)
		res := checkTerminal(start, &m)
		if res != nil {
			return nil
		} else {
			return customerrors.ErrNoSuchWord
		}
	}

	(*used)[start] = index
	for key, val := range start.next {
		switch key {
		case rune(word[index]):
			for _, st := range val {
				e := readWord(st, word, index+1, used)
				if e == nil {
					return nil
				}
			}
		case emptyWord:
			for _, st := range val {
				v, ok := (*used)[st]
				if !ok || v != index {
					e := readWord(st, word, index, used)
					if e == nil {
						return nil
					}
				}
			}
		}
	}
	return customerrors.ErrNoSuchWord
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
