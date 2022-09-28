package automate

import (
	"reflect"
	customerrors "workspace/internal/errors"
)

func NewAutomate() *Automate {
	return &Automate{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

func DetermineAutomate(a Automate) (*Automate, error) {
	alphabet := make(map[rune]bool)
	used := make(map[*state]bool)
	getAlphabet(a.startState, &alphabet, &used)

	sets := make([]*map[*state]bool, 1)
	transitions := make(map[int]map[rune]int, 0)

	firstSet := make(map[*state]bool)
	firstSet[a.startState] = true

	sets[0] = &firstSet
	var setsIndex int

	for {
		var (
			index       int
			needAdd     bool
			newSetAdded bool
		)

		transitions[setsIndex] = make(map[rune]int)
		newSets := make([]*map[*state]bool, len(alphabet))
		for letter := range alphabet {
			newSets[index] = &map[*state]bool{}

			for key := range *(sets[setsIndex]) {
				for _, st := range key.next[letter] {
					(*(newSets[index]))[st] = true
				}
			}

			if *(newSets[index]) == nil || len(*(newSets[index])) == 0 {
				index++
				continue
			}

			needAdd = true
			var indexEqual int
			for i, set := range sets {
				if len(*set) == len(*(newSets[index])) && reflect.DeepEqual(*set, *(newSets[index])) {
					needAdd = false
					indexEqual = i
				}
			}

			if needAdd {
				sets = append(sets, newSets[index])
				transitions[len(sets)-1] = map[rune]int{}
				transitions[setsIndex][letter] = len(sets) - 1
				newSetAdded = true
			} else {
				transitions[setsIndex][letter] = indexEqual
			}

			index++
		}

		setsIndex++
		if !newSetAdded && setsIndex == len(sets) {
			break
		}
	}

	created := make(map[int]*state)

	res := &Automate{}
	res.startState = ConstructState(&transitions, &sets, &created, 0)

	return res, nil
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

func ConstructState(transitions *map[int]map[rune]int, sets *[]*map[*state]bool, created *map[int]*state, index int) *state {
	res := &state{next: make(map[rune][]*state)}
	(*created)[index] = res
	for key := range *((*sets)[index]) {
		if key.isTerm {
			res.isTerm = true
			break
		}
	}

	for key, val := range (*transitions)[index] {
		set, ok := (*created)[val]
		if ok {
			res.next[key] = make([]*state, 1)
			res.next[key][0] = set
		} else {
			newState := ConstructState(transitions, sets, created, val)
			res.next[key] = make([]*state, 1)
			res.next[key][0] = newState
		}
	}
	return res
}

type Automate struct {
	startState *state
	terminals  []*state
}

func getAlphabet(st *state, alphabet *map[rune]bool, used *map[*state]bool) {
	(*used)[st] = true
	for key, value := range st.next {
		if key != emptyWord {
			(*alphabet)[key] = true
		}
		for _, s := range value {
			_, ok := (*used)[s]
			if !ok {
				getAlphabet(s, alphabet, used)
			}
		}
	}
}

func deleteEps(st *state, parent *state, letter rune, used *map[*state]int) error {
	_, ok := (*used)[st]
	switch ok {
	case false:
		(*used)[st] = 1
	case true:
		(*used)[st]++
	}

	for key, val := range st.next {
		for _, s := range val {
			v, ok := (*used)[s]
			if ok && v == 2 {
				continue
			}
			if ok && key == emptyWord {
				for l, to := range s.next {
					if l != emptyWord {
						st.next[l] = append(st.next[l], to...)
					}
				}
			}
			deleteEps(s, st, key, used)
		}
	}

	delete(st.next, emptyWord)
	if letter == emptyWord && parent != nil {
		if st.isTerm {
			parent.isTerm = true
		}
		for k, v := range st.next {
			parent.next[k] = append(parent.next[k], v...)
		}
	}

	return nil
}

func (a *Automate) DeleteEps() error {
	used := make(map[*state]int)
	err := deleteEps(a.startState, nil, emptyWord, &used)
	return err
}

func (a *Automate) GetStates() ([]*State, error) {
	used := make(map[*state]int)
	ans := make([]*State, 0)

	var counter int
	err := mapStates(a.startState, &used, &ans, &counter)

	return ans, err
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

func (a *Automate) Read(word string) error {
	used := make(map[*state]int)
	return readWord(a.startState, word, 0, &used)
}

func (a *Automate) putStartState(auto *Automate) error {
	a.startState = auto.startState
	a.terminals = auto.terminals
	return nil
}

func (a *Automate) AddNewWord(word string) error {
	curState := a.startState
	for _, letter := range word {
		_, ok := curState.next[letter]
		switch ok {
		case true:
			curState = curState.next[letter][0]
		case false:
			curState.next[letter] = append(
				curState.next[letter],
				&state{isTerm: false, next: make(map[rune][]*state)},
			)
			curState = curState.next[letter][0]
		}
	}

	curState.isTerm = true
	a.terminals = append(a.terminals, curState)
	return nil
}

func (a *Automate) Cycle() error {
	newStartState := &state{next: make(map[rune][]*state)}
	a.startState.isTerm = true

	newStartState.next[emptyWord] = []*state{a.startState}

	for _, term := range a.terminals {
		term.next[emptyWord] = append(term.next[emptyWord], newStartState)
	}

	a.terminals = append(a.terminals, a.startState)
	a.startState = newStartState

	return nil
}

func addEdgesToState(start *state, to *state, used *map[*state]bool) {
	(*used)[start] = true
	for key, val := range start.next {
		for i, st := range val {
			if st.isTerm {
				start.next[key][i].isTerm = false
				start.next[key] = append(start.next[key], to)
			}
			_, ok := (*used)[st]
			if !ok {
				addEdgesToState(st, to, used)
			}
		}
	}
}

func (a *Automate) Concat(other *Automate) error {
	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}

	used := make(map[*state]bool)
	addEdgesToState(a.startState, other.startState, &used)

	a.terminals = other.terminals
	return nil
}

func (a *Automate) Join(other *Automate) error {
	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}

	a.startState.next[emptyWord] = append(a.startState.next[emptyWord], other.startState)
	a.terminals = append(a.terminals, other.terminals...)
	return nil
}

type state struct {
	next   map[rune][]*state
	isTerm bool
}
