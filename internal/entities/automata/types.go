package automata

import (
	customerrors "workspace/internal/errors"
)

type walker struct{}

func (w *walker) walk(start *state, used map[*state]bool, executor Executor) error {
	used[start] = true
	for _, val := range start.next {
		for _, v := range val {
			_, ok := used[v]
			if !ok {
				err := w.walk(v, used, executor)
				if err != nil {
					return err
				}
			}
		}
	}

	executor.Execute(start)
	return nil
}

func (w *walker) walkCheck(start *state, used map[*state]bool, checker Checker) error {
	used[start] = true
	for _, val := range start.next {
		for _, v := range val {
			_, ok := used[v]
			if !ok {
				err := w.walkCheck(v, used, checker)
				if err != nil {
					return err
				}
			}
		}
	}

	return checker.Check(start)
}

type inverter struct {
	terminals *[]*state
}

func (i *inverter) Execute(state *state) error {
	switch state.isTerminal {
	case true:
		state.isTerminal = false
	case false:
		state.isTerminal = true
		*(i.terminals) = append(*(i.terminals), state)
	}

	return nil
}

type deleterEps struct{}

func (d *deleterEps) Execute(state *state) error {
	delete(state.next, emptyWord)
	return nil
}

type fuller struct {
	alphabet *map[string]bool
	stock    *state
}

func (f *fuller) Execute(st *state) error {
	for key := range *(f.alphabet) {
		_, ok := st.next[key]
		if !ok {
			st.next[key] = []*state{f.stock}
		}
	}
	return nil
}

type getterAlphabet struct {
	alphabet *map[string]bool
}

func (g *getterAlphabet) Execute(state *state) error {
	for key := range state.next {
		if key != emptyWord {
			(*(g.alphabet))[key] = true
		}
	}

	return nil
}

type getterTerminals struct {
	classes *map[*state]int
}

func (g *getterTerminals) Execute(state *state) error {
	switch state.isTerminal {
	case true:
		(*(g.classes))[state] = 0
	case false:
		(*(g.classes))[state] = 1
	}
	return nil
}

type getterStates struct {
	states *[]*state
}

func (g *getterStates) Execute(state *state) error {
	*(g.states) = append(*(g.states), state)
	return nil
}

// checkers

type determineChecker struct{}

func (d *determineChecker) Check(state *state) error {
	for _, statesTo := range state.next {
		if len(statesTo) >= 2 {
			return customerrors.ErrNotDetermine
		}
	}
	return nil
}

type noEpsilonChecker struct{}

func (n *noEpsilonChecker) Check(state *state) error {
	for letter := range state.next {
		if letter == emptyWord {
			return customerrors.ErrHasEpsilonTransitions
		}
	}
	return nil
}

type fullChecker struct {
	alphabet map[string]bool
}

func (f *fullChecker) Check(state *state) error {
	for letter := range f.alphabet {
		if _, ok := state.next[letter]; !ok {
			return customerrors.ErrNotFull
		}

	}
	return nil
}
