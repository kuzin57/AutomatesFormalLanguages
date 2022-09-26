package automate

import (
	"fmt"
	customerrors "workspace/internal/errors"
)

func NewNFA() *nfa {
	return &nfa{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

func NewFA() *fa {
	return &fa{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

type nfa struct {
	startState *state
	terminals  []*state
}

func (a *nfa) DeleteEps() error {
	type stateAndIndex struct {
		state      *state
		indexEmpty int
	}
	stackStates := make([]*stateAndIndex, 1)
	stackStates[0] = &stateAndIndex{state: a.startState}
	var cur *stateAndIndex

	for len(stackStates) > 0 {
		cur = stackStates[len(stackStates)-1]

		for key, val := range cur.state.next {
			if key == emptyWord {
				continue
			}

			for i := 0; i < len(stackStates)-1; i++ {
				stackStates[i].state.next[key] = append(stackStates[i].state.next[key], val...)
			}

			if cur.state.isTerm {
				for i := 0; i < len(stackStates); i++ {
					stackStates[i].state.isTerm = true
				}
			}
		}

		stateWithEmpty, empty := cur.state.next[emptyWord]
		if !empty || cur.indexEmpty == len(stateWithEmpty) {
			stackStates = stackStates[:len(stackStates)-1]
			delete(cur.state.next, emptyWord)
			continue
		}

		stackStates = append(stackStates, &stateAndIndex{state: stateWithEmpty[cur.indexEmpty]})
		cur.indexEmpty++
	}
	return nil
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
