package automata

import (
	customerrors "workspace/internal/errors"
)

func checkStates(realStates []*state, states []*State) (err error) {
	for i, realState := range realStates {
		for letter, statesTo := range realState.next {
			if _, ok := states[i].Transitions[letter]; !ok {
				return customerrors.ErrIncorrectMapping
			}
			if len(statesTo) != len(states[i].Transitions[letter]) {
				return customerrors.ErrIncorrectMapping
			}
		}
	}
	return nil
}

func (a *Automata) CheckDetermine() error {
	used := make(map[*state]bool)
	return a.walker.walkCheck(a.startState, used, &determineChecker{})
}

func (a *Automata) CheckNoEpsilon() error {
	used := make(map[*state]bool)
	return a.walker.walkCheck(a.startState, used, &noEpsilonChecker{})
}

func (a *Automata) CheckStates() error {
	used := make(map[*state]bool)
	realStates := make([]*state, 0)
	if err := a.walker.walk(a.startState, used, &getterStates{states: &realStates}); err != nil {
		return err
	}

	states, err := mapStates(realStates)
	if err != nil {
		return err
	}
	return checkStates(realStates, states)
}

func (a *Automata) CheckFull() error {
	alphabet := make(map[string]bool)
	used := make(map[*state]bool)

	if err := a.walker.walk(a.startState, used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}

	used = make(map[*state]bool)
	return a.walker.walkCheck(a.startState, used, &fullChecker{alphabet: alphabet})
}
