package automata

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
