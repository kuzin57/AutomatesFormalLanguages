package automate

import (
	customerrors "workspace/internal/errors"
)

type State struct {
	Number      int
	Transitions map[string][]int
	IsTerminal  bool
}

func mapStates(states []*state) ([]*State, error) {
	ans := make([]*State, 0)
	statesNumbers := make(map[*state]int)
	for i, st := range states {
		statesNumbers[st] = i + 1
	}

	for i, st := range states {
		newState := &State{Transitions: make(map[string][]int), Number: i + 1, IsTerminal: st.isTerm}
		for key, val := range st.next {
			for _, v := range val {
				_, ok := newState.Transitions[key]
				switch ok {
				case true:
					newState.Transitions[key] = append(newState.Transitions[key], statesNumbers[v])
				case false:
					newState.Transitions[key] = []int{statesNumbers[v]}
				}
			}
		}
		ans = append(ans, newState)
	}
	return ans, nil
}

func addNewTransition(from *state, to *state, key string) {
	needAdd := true
	for _, node := range from.next[key] {
		if node == to {
			needAdd = false
		}
	}

	if needAdd {
		from.next[key] = append(from.next[key], to)
	}
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
				for _, v := range val {
					addNewTransition(st, v, key)
				}
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
		case string(word[index]):
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

func proccessSelfTransitions(st *state) (string, error) {
	res := ""
	for key, val := range st.next {
		toDelete := -1
		for i, v := range val {
			if v == st {
				if len(res) > 0 {
					res += "+"
				}
				res += key
				toDelete = i
			}
		}

		if toDelete != -1 {
			st.next[key] = append(st.next[key][:toDelete], st.next[key][(toDelete+1):]...)
		}
	}

	return res, nil
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

func removeDublicates(st *state) error {
	regulars := make(map[*state]string)

	for key, val := range st.next {
		for _, v := range val {
			_, ok := regulars[v]
			if !ok {
				regulars[v] = key
			} else {
				regulars[v] += "+"
				regulars[v] += key
			}
		}
	}

	st.next = make(map[string][]*state)
	for key, val := range regulars {
		_, ok := st.next[val]
		if !ok {
			st.next[val] = append(st.next[val], key)
		}
	}

	return nil
}
