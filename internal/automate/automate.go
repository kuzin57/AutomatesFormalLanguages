package automate

import (
	"fmt"
	"reflect"
	customerrors "workspace/internal/errors"
)

func NewNFA() *nfa {
	return &nfa{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

func NewFA() *fa {
	return &fa{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

func NewFAFromNFA(other Automate) (*fa, error) {
	a, ok := other.(*nfa)
	if !ok {
		return nil, fmt.Errorf("non deterministic automate expected")
	}

	alphabet := make(map[rune]bool)
	getAlphabet(a.startState, &alphabet)
	fmt.Println("alphabet", alphabet)

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

		transitions[setsIndex] = map[rune]int{}
		newSets := make([]*map[*state]bool, len(alphabet))
		for letter := range alphabet {
			newSets[index] = &map[*state]bool{}

			for key := range *(sets[setsIndex]) {
				for _, st := range (*key).next[letter] {
					(*(newSets[index]))[st] = true
				}
			}

			if *(newSets[index]) == nil {
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

		if !newSetAdded {
			break
		}
		setsIndex++
	}

	created := make(map[int]*state)

	res := &fa{}
	res.startState = ConstructState(&transitions, &sets, &created, 0)

	return res, nil
}

func NewNFAFromFA(other Automate) (*nfa, error) {
	realAutomate, ok := other.(*fa)
	if !ok {
		return nil, fmt.Errorf("deterministic automate expected")
	}

	res := &nfa{}
	res.startState = realAutomate.startState
	return res, nil
}

func ConstructState(transitions *map[int]map[rune]int, sets *[]*map[*state]bool, created *map[int]*state, index int) *state {
	res := &state{next: make(map[rune][]*state)}
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

type nfa struct {
	startState *state
	terminals  []*state
}

func getAlphabet(st *state, alphabet *map[rune]bool) {
	for key, val := range st.next {
		fmt.Println("key:", string(key))
		if key != emptyWord {
			(*alphabet)[key] = true
			for _, st := range val {
				getAlphabet(st, alphabet)
			}
		}
	}
}

func deleteEps(st *state, parent *state, letter rune, used *map[*state]bool) error {
	(*used)[st] = true
	for key, val := range st.next {

		for _, s := range val {
			_, ok := (*used)[s]
			if ok {
				continue
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

func (a *nfa) DeleteEps() error {
	used := make(map[*state]bool)
	err := deleteEps(a.startState, nil, emptyWord, &used)
	fmt.Println("haahahahahahahahahah", a.startState.next['f'][0].next)
	return err
}

func (a *nfa) Check() bool {
	return false
}

func (a *nfa) Read(word string) error {
	type stateAndIndex struct {
		state      *state
		stateIndex int
		indexEmpty int
		wordIndex  int
	}

	stackStates := make([]*stateAndIndex, 1)
	stackStates[0] = &stateAndIndex{state: a.startState}
	var cur *stateAndIndex

	for len(stackStates) > 0 {
		cur = stackStates[len(stackStates)-1]

		fmt.Println("word", word, cur.wordIndex)
		if cur.wordIndex == len(word) && cur.state.isTerm {
			return nil
		}

		_, empty := cur.state.next[emptyWord]
		if empty && cur.indexEmpty < len(cur.state.next[emptyWord]) {
			stackStates = append(
				stackStates,
				&stateAndIndex{state: cur.state.next[emptyWord][cur.indexEmpty], wordIndex: cur.wordIndex},
			)
			cur.indexEmpty++
			continue
		}

		var ok bool
		if len(word) > cur.wordIndex {
			_, ok = cur.state.next[rune(word[cur.wordIndex])]
		}

		if len(word) == cur.wordIndex || !ok || cur.stateIndex == len(cur.state.next[rune(word[cur.wordIndex])]) {
			stackStates = stackStates[:len(stackStates)-1]
			continue
		}
		stackStates = append(
			stackStates,
			&stateAndIndex{
				state:     cur.state.next[rune(word[cur.wordIndex])][cur.stateIndex],
				wordIndex: cur.wordIndex + 1,
			},
		)

		cur.stateIndex++
	}
	return customerrors.ErrNoSuchWord
}

func (a *nfa) putStartState(automate Automate) error {
	realAutomate, ok := automate.(*nfa)
	if !ok {
		return customerrors.ErrNoAutomate
	}

	a.startState = realAutomate.startState
	a.terminals = realAutomate.terminals
	return nil
}

func (a *nfa) AddNewWord(word string) error {
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

func (a *nfa) Cycle() error {
	a.startState.isTerm = true

	for _, term := range a.terminals {
		term.next[emptyWord] = append(term.next[emptyWord], a.startState)
	}
	a.terminals = append(a.terminals, a.startState)

	return nil
}

func (a *nfa) Concat(other Automate) error {
	realAutomate, ok := other.(*nfa)
	if !ok {
		return fmt.Errorf("can not concat automates of different types")
	}

	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}

	for _, term := range a.terminals {
		term.isTerm = false
		term.next[emptyWord] = append(term.next[emptyWord], realAutomate.startState)
	}

	a.terminals = realAutomate.terminals
	return nil
}

func (a *nfa) Join(other Automate) error {
	realAutomate, ok := other.(*nfa)
	if !ok {
		return fmt.Errorf("can not join automates of different types")
	}

	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}

	a.startState.next[emptyWord] = append(a.startState.next[emptyWord], realAutomate.startState)
	a.terminals = append(a.terminals, realAutomate.terminals...)
	return nil
}

type fa struct {
	startState *state
}

func (a *fa) DeleteEps() error {
	return customerrors.ErrNotImplemented
}

func (a *fa) Read(line string) error {
	return nil
}

func (a *fa) Check() bool {
	return true
}

func (a *fa) AddNewWord(word string) error {
	return customerrors.ErrNotImplemented
}

func (a *fa) Concat(other Automate) error {
	return customerrors.ErrNotImplemented
}

func (a *fa) Join(other Automate) error {
	return customerrors.ErrNotImplemented
}

func (a *fa) Cycle() error {
	return customerrors.ErrNotImplemented
}

func (a *fa) putStartState(automate Automate) error {
	return customerrors.ErrNotImplemented
}

type state struct {
	next   map[rune][]*state
	isTerm bool
}
