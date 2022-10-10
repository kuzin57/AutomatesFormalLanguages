package automata

import (
	"reflect"
)

func NewAutomata() *Automata {
	return &Automata{startState: &state{isTerminal: false, next: make(map[string][]*state)}}
}

func getSetsTranisitions(a *Automata) (*[]*map[*state]bool, *map[int]map[string]int, error) {
	alphabet := make(map[string]bool)
	used := make(map[*state]bool)

	if err := a.walker.walk(a.startState, used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return nil, nil, err
	}

	sets := make([]*map[*state]bool, 1)
	transitions := make(map[int]map[string]int, 0)

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

		transitions[setsIndex] = make(map[string]int)
		newSets := make([]*map[*state]bool, len(alphabet))
		for letter := range alphabet {
			newSets[index] = &map[*state]bool{}

			for state := range *(sets[setsIndex]) {
				for _, st := range state.next[letter] {
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
				transitions[len(sets)-1] = map[string]int{}
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

	return &sets, &transitions, nil
}

func (a *Automata) Determine() (*Automata, error) {
	sets, transitions, err := getSetsTranisitions(a)
	if err != nil {
		return nil, err
	}

	created := make(map[int]*state)
	result := &Automata{}
	result.startState = result.constructState(transitions, sets, &created, 0)
	return result, nil
}

func (a *Automata) constructState(
	transitions *map[int]map[string]int,
	sets *[]*map[*state]bool,
	created *map[int]*state,
	index int,
) *state {
	result := &state{next: make(map[string][]*state)}
	(*created)[index] = result

	for state := range *((*sets)[index]) {
		if state.isTerminal {
			result.isTerminal = true
			a.terminals = append(a.terminals, result)
			break
		}
	}

	for letter, numberState := range (*transitions)[index] {
		set, ok := (*created)[numberState]
		if ok {
			result.next[letter] = make([]*state, 1)
			result.next[letter][0] = set
		} else {
			newState := a.constructState(transitions, sets, created, numberState)
			result.next[letter] = make([]*state, 1)
			result.next[letter][0] = newState
		}
	}
	return result
}

type Automata struct {
	startState *state
	terminals  []*state
	walker     walker
}

func (a *Automata) DeleteEps() error {
	used := make(map[*state]bool)

	if err := iterateStatesDeleteEps(a.startState, &used); err != nil {
		return err
	}

	used = make(map[*state]bool)
	if err := a.walker.walk(a.startState, used, &deleterEps{}); err != nil {
		return err
	}

	return nil
}

func (a *Automata) GetStates() (ans []*State, err error) {
	var (
		states = make([]*state, 0)
		used   = make(map[*state]bool)
	)

	if err := a.walker.walk(a.startState, used, &getterStates{states: &states}); err != nil {
		return nil, err
	}

	left := 0
	right := len(states) - 1
	for left < right {
		states[left], states[right] = states[right], states[left]
		left++
		right--
	}

	ans, err = mapStates(states)
	return
}

func (a *Automata) Read(word string) error {
	used := make(map[*state]int)
	return readWord(a.startState, word, 0, &used)
}

func (a *Automata) putStartState(auto *Automata) error {
	a.startState = auto.startState
	a.terminals = auto.terminals
	return nil
}

func (a *Automata) AddNewWord(word string) error {
	curState := a.startState
	for _, letter := range word {
		strLetter := string(letter)
		_, ok := curState.next[strLetter]
		switch ok {
		case true:
			curState = curState.next[strLetter][0]
		case false:
			curState.next[strLetter] = append(
				curState.next[strLetter],
				&state{isTerminal: false, next: make(map[string][]*state)},
			)
			curState = curState.next[strLetter][0]
		}
	}

	curState.isTerminal = true
	a.terminals = append(a.terminals, curState)
	return nil
}

func (a *Automata) Cycle() error {
	newStartState := &state{next: make(map[string][]*state)}
	newStartState.next[emptyWord] = []*state{a.startState}
	newStartState.isTerminal = true

	for _, terminal := range a.terminals {
		terminal.next[emptyWord] = append(terminal.next[emptyWord], newStartState)
	}

	a.terminals = append(a.terminals, newStartState)
	a.startState = newStartState
	return nil
}

func (a *Automata) Concat(other *Automata) error {
	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}
	for _, terminal := range a.terminals {
		terminal.isTerminal = false
		terminal.next[emptyWord] = append(terminal.next[emptyWord], other.startState)
	}

	a.terminals = other.terminals
	return nil
}

func (a *Automata) Join(other *Automata) error {
	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}

	newState := &state{next: make(map[string][]*state)}
	newState.next[emptyWord] = make([]*state, 2)
	newState.next[emptyWord][0] = a.startState
	newState.next[emptyWord][1] = other.startState

	a.startState = newState
	a.terminals = append(a.terminals, other.terminals...)
	return nil
}

func (a *Automata) Full() error {
	var (
		stock    = &state{next: make(map[string][]*state)}
		alphabet = make(map[string]bool)
		used     = make(map[*state]bool)
	)

	if err := a.walker.walk(a.startState, used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}
	for letter := range alphabet {
		stock.next[letter] = []*state{stock}
	}

	used = make(map[*state]bool)
	return a.walker.walk(a.startState, used, &fuller{alphabet: &alphabet, stock: stock})
}

func (a *Automata) Invert() error {
	var (
		newTerminals = make([]*state, 0)
		used         = make(map[*state]bool)
	)

	if err := a.walker.walk(a.startState, used, &inverter{terminals: &newTerminals}); err != nil {
		return err
	}

	a.terminals = newTerminals
	return nil
}

func getDifference(
	previousEquivalenceClasses map[*state]int,
	states []*state,
	state *state,
	alphabetArray []string,
	stateNumber int,
) int {
	differenceIndex := -1

	for j := 0; j < stateNumber; j++ {
		if previousEquivalenceClasses[states[j]] != previousEquivalenceClasses[state] {
			continue
		}
		differenceIndex = j
		for _, letter := range alphabetArray {
			if previousEquivalenceClasses[states[j].next[letter][0]] != previousEquivalenceClasses[state.next[letter][0]] {
				differenceIndex = -1
			}
		}
		if differenceIndex != -1 {
			break
		}
	}
	return differenceIndex
}

func getStatesTerminals(currentEquivalenceClasses map[*state]int, classesNumber int) ([]*state, map[*state]bool) {
	newStates := make([]*state, classesNumber)
	for i := 0; i < classesNumber; i++ {
		newStates[i] = &state{next: make(map[string][]*state)}
	}

	terminals := make(map[*state]bool)
	for state, classNumber := range currentEquivalenceClasses {
		if state.isTerminal {
			newStates[classNumber].isTerminal = true

			_, ok := terminals[newStates[classNumber]]
			if !ok {
				terminals[newStates[classNumber]] = true
			}
		}

		for letter, statesTo := range state.next {
			for _, stateTo := range statesTo {
				addNewTransition(newStates[classNumber], newStates[currentEquivalenceClasses[stateTo]], letter)
			}
		}
	}
	return newStates, terminals
}

func (a *Automata) Minimize() error {
	var (
		previousEquivalenceClasses = make(map[*state]int)
		currentEquivalenceClasses  = make(map[*state]int)
		alphabet                   = make(map[string]bool)
		used                       = make(map[*state]bool)
		alphabetArray              = make([]string, 0)
		states                     = make([]*state, 0)
		counterClasses             int
	)

	if err := a.walker.walk(a.startState, used, &getterTerminals{classes: &previousEquivalenceClasses}); err != nil {
		return err
	}

	used = make(map[*state]bool)
	if err := a.walker.walk(a.startState, used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}

	for letter := range alphabet {
		alphabetArray = append(alphabetArray, letter)
	}

	used = make(map[*state]bool)
	if err := a.walker.walk(a.startState, used, &getterStates{states: &states}); err != nil {
		return err
	}

	for {
		counterClasses = 0
		for i, state := range states {
			differenceIndex := getDifference(previousEquivalenceClasses, states, state, alphabetArray, i)
			switch differenceIndex {
			case -1:
				currentEquivalenceClasses[state] = counterClasses
				counterClasses++
			default:
				currentEquivalenceClasses[state] = currentEquivalenceClasses[states[differenceIndex]]
			}
		}

		isEqual := true
		for k, v := range previousEquivalenceClasses {
			if v != currentEquivalenceClasses[k] {
				isEqual = false
			}
		}

		if len(currentEquivalenceClasses) == len(previousEquivalenceClasses) && isEqual {
			break
		}

		previousEquivalenceClasses = make(map[*state]int)
		for k, v := range currentEquivalenceClasses {
			previousEquivalenceClasses[k] = v
		}
	}

	newStates, terminals := getStatesTerminals(currentEquivalenceClasses, counterClasses)

	a.startState = newStates[currentEquivalenceClasses[a.startState]]
	a.terminals = make([]*state, 0)
	for k := range terminals {
		a.terminals = append(a.terminals, k)
	}

	return nil
}

func addExpression(to string, inBrackets string, afterBrackets string) string {
	switch {
	case inBrackets == "":
		to += inBrackets
	case len(inBrackets) == 1:
		to += inBrackets
		to += afterBrackets
	default:
		to += "("
		to += inBrackets
		to += ")"
		to += afterBrackets
	}
	return to
}

func compressTransitions(currentState *state, processedState *state, transtitionsToAdd map[string][]*state, letter string, expression string) {
	for k, node := range currentState.next {
		for _, n := range node {
			var ex string
			ex = addExpression(ex, letter, "")
			if ex != "" && !(expression == "" && k == "") {
				ex += "."
			}
			ex = addExpression(ex, expression, "*")
			if expression != "" && k != "" {
				ex += "."
			}
			ex = addExpression(ex, k, "")
			_, ok := transtitionsToAdd[ex]
			switch ok {
			case true:
				transtitionsToAdd[ex] = append(transtitionsToAdd[ex], n)
			case false:
				transtitionsToAdd[ex] = []*state{n}
			}
		}
	}
}

func deleteStateAddTransitions(stateToDelete *state, processedState *state, expression string) {
	toAdd := make(map[string][]*state)
	for letter, statesTo := range processedState.next {
		for _, stateTo := range statesTo {
			if stateTo == stateToDelete {
				compressTransitions(stateTo, processedState, toAdd, letter, expression)
			}
		}
	}

	for exp, node := range toAdd {
		for _, nn := range node {
			_, ok := processedState.next[exp]
			if !ok {
				processedState.next[exp] = []*state{nn}
			} else {
				processedState.next[exp] = append(processedState.next[exp], nn)
			}
		}
	}

	for letter := range processedState.next {
		toDelete := -1
		for j, n := range processedState.next[letter] {
			if n == stateToDelete {
				toDelete = j
			}
		}

		if toDelete != -1 {
			if len(processedState.next[letter]) == 1 {
				processedState.next[letter] = nil
			} else {
				processedState.next[letter] = append(
					processedState.next[letter][:toDelete],
					processedState.next[letter][(toDelete+1):]...,
				)
			}

			if len(processedState.next[letter]) == 0 {
				delete(processedState.next, letter)
			}
		}
	}
}

func (a *Automata) formExpressionFromStartState() string {
	var fromStartToStart string
	var fromStartToEnd string
	for key, val := range a.startState.next {
		for _, v := range val {
			if v == a.startState && key != "" {
				fromStartToStart = "(" + key + ")" + "*"
			} else {
				fromStartToEnd = key
			}
		}
	}
	if fromStartToEnd == "" {
		return fromStartToStart
	}
	return fromStartToStart + "." + fromStartToEnd
}

func (a *Automata) GetRegularExpression() (string, error) {
	var (
		newTerm = &state{isTerminal: true}
		used    = make(map[*state]bool)
		states  = make([]*state, 0)
	)

	for _, terminal := range a.terminals {
		terminal.isTerminal = false
		terminal.next[""] = append(terminal.next[""], newTerm)
	}

	if err := a.walker.walk(a.startState, used, &getterStates{states: &states}); err != nil {
		return "", err
	}

	for {
		var (
			exit          = true
			stateToDelete *state
			indexToDelete int
		)

		for i, state := range states {
			if !state.isTerminal && a.startState != state {
				stateToDelete = state
				exit = false
				indexToDelete = i
				break
			}
		}

		for _, state := range states {
			err := removeDublicates(state)
			if err != nil {
				return "", err
			}
		}

		if exit || len(states) == 2 {
			break
		}

		expression, err := proccessSelfTransitions(stateToDelete)
		if err != nil {
			return "", err
		}

		for _, state := range states {
			err := removeDublicates(state)
			if err != nil {
				return "", err
			}
			deleteStateAddTransitions(stateToDelete, state, expression)
		}

		states = append(states[:indexToDelete], states[(indexToDelete+1):]...)
	}

	expr, err := proccessSelfTransitions(a.startState)
	if err != nil {
		return "", err
	}

	_, ok := a.startState.next[expr]
	if !ok {
		a.startState.next[expr] = []*state{a.startState}
	} else {
		a.startState.next[expr] = append(a.startState.next[expr], a.startState)
	}

	err = removeDublicates(a.startState)
	if err != nil {
		return "", err
	}

	return a.formExpressionFromStartState(), nil
}

type state struct {
	next       map[string][]*state
	isTerminal bool
}
