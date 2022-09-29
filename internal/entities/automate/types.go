package automate

import (
	customerrors "workspace/internal/errors"
)

type inverter struct{}

func (i *inverter) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	terminals, ok := args[2].(*[]*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	switch st.isTerm {
	case true:
		st.isTerm = false
	case false:
		st.isTerm = true
		(*terminals) = append((*terminals), st)
	}

	return nil
}

type deleterEps struct{}

func (d *deleterEps) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	delete(st.next, emptyWord)

	return nil
}

type fuller struct {
	alphabet *map[string]bool
}

func (f *fuller) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	stock, ok := args[1].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	for key := range *(f.alphabet) {
		_, ok := st.next[key]
		if !ok {
			st.next[key] = []*state{stock}
		}
	}
	return nil
}

type getterAlphabet struct {
	alphabet *map[string]bool
}

func (g *getterAlphabet) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	for key := range st.next {
		if key != emptyWord {
			(*(g.alphabet))[key] = true
		}
	}

	return nil
}

type getterTerminals struct{}

func (g *getterTerminals) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	classes, ok := args[3].(*map[*state]int)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	switch st.isTerm {
	case true:
		(*classes)[st] = 0
	case false:
		(*classes)[st] = 1
	}

	return nil
}

type getterStates struct{}

func (g *getterStates) Execute(args ...any) error {
	st, ok := args[0].(*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	states, ok := args[2].(*[]*state)
	if !ok {
		return customerrors.ErrInvalidArgument
	}

	*states = append(*states, st)
	return nil
}
