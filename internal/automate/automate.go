package automate

import (
	"fmt"
	customerrors "workspace/internal/errors"
)

func NewNFA() *nfa {
	return &nfa{startState: &state{isTerm: false, next: make(map[rune]*state)}}
}

func NewFA() *fa {
	return &fa{startState: &state{isTerm: false, next: make(map[rune]*state)}}
}

type nfa struct {
	startState *state
	terminals  []*state
}

func (a *nfa) DeleteEps() error {
	return customerrors.ErrNotImplemented
}

func (a *nfa) Check() bool {
	return false
}

func (a *nfa) Read(line string) (ok bool) {
	curState := a.startState
	// queueEmpties := make([]*state, 0)

	for _, letter := range line {
		curState, ok = curState.next[letter]
		if !ok {
			return
		}
	}
	return true
}

func (a *nfa) AddNewWord(word string) error {
	curState := a.startState
	for _, letter := range word {
		_, ok := curState.next[letter]
		switch ok {
		case true:
			curState = curState.next[letter]
		case false:
			curState.next[letter] = &state{isTerm: false, next: make(map[rune]*state)}
			curState = curState.next[letter]
		}
	}

	curState.isTerm = true
	a.terminals = append(a.terminals, curState)
	return nil
}

func (a *nfa) Cycle() error {
	for _, term := range a.terminals {
		term.next[emptyWord] = a.startState
	}
	return nil
}

func (a *nfa) Concat(other Automate) error {
	realAutomate, ok := other.(*nfa)
	if !ok {
		return fmt.Errorf("can not concat automates of different types")
	}

	for _, term := range a.terminals {
		term.isTerm = false
		term.next[emptyWord] = realAutomate.startState
	}
	return nil
}

func (a *nfa) Join(other Automate) error {
	realAutomate, ok := other.(*nfa)
	if !ok {
		return fmt.Errorf("can not join automates of different types")
	}

	newStartState := &state{next: make(map[rune]*state)}

	newStartState.next[emptyWord] = a.startState
	newStartState.next[emptyWord] = realAutomate.startState
	a.startState = newStartState

	return nil
}

type fa struct {
	startState *state
}

func (a *fa) DeleteEps() error {
	return customerrors.ErrNotImplemented
}

func (a *fa) Read(line string) bool {
	return false
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

type state struct {
	next   map[rune]*state
	isTerm bool
}
