package automate

import (
	"reflect"
)

func NewAutomate() *Automate {
	return &Automate{startState: &state{isTerm: false, next: make(map[string][]*state)}}
}

func getSetsTranisitions(a *Automate) (*[]*map[*state]bool, *map[int]map[string]int, error) {
	alphabet := make(map[string]bool)
	used := make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, nil, &used, &getterAlphabet{alphabet: &alphabet}); err != nil {
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

func (a *Automate) Determine() (*Automate, error) {
	sets, transitions, err := getSetsTranisitions(a)
	if err != nil {
		return nil, err
	}

	created := make(map[int]*state)
	res := &Automate{}
	res.startState = res.ConstructState(transitions, sets, &created, 0)

	return res, nil
}

func (a *Automate) ConstructState(
	transitions *map[int]map[string]int,
	sets *[]*map[*state]bool,
	created *map[int]*state,
	index int,
) *state {
	res := &state{next: make(map[string][]*state)}
	(*created)[index] = res

	for key := range *((*sets)[index]) {
		if key.isTerm {
			res.isTerm = true
			a.terminals = append(a.terminals, res)
			break
		}
	}

	for key, val := range (*transitions)[index] {
		set, ok := (*created)[val]
		if ok {
			res.next[key] = make([]*state, 1)
			res.next[key][0] = set
		} else {
			newState := a.ConstructState(transitions, sets, created, val)
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

func (a *Automate) DeleteEps() error {
	used := make(map[*state]bool)

	if err := iterateStatesDeleteEps(a.startState, &used); err != nil {
		return err
	}

	used = make(map[*state]bool)
	if err := traversal(a.startState, nil, nil, nil, &used, &deleterEps{}); err != nil {
		return err
	}

	return nil
}

func (a *Automate) GetStates() (ans []*State, err error) {
	var (
		states = make([]*state, 0)
		used   = make(map[*state]bool)
	)

	if err := traversal(a.startState, nil, &states, nil, &used, &getterStates{}); err != nil {
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
		strLetter := string(letter)
		_, ok := curState.next[strLetter]
		switch ok {
		case true:
			curState = curState.next[strLetter][0]
		case false:
			curState.next[strLetter] = append(
				curState.next[strLetter],
				&state{isTerm: false, next: make(map[string][]*state)},
			)
			curState = curState.next[strLetter][0]
		}
	}

	curState.isTerm = true
	a.terminals = append(a.terminals, curState)
	return nil
}

func (a *Automate) Cycle() error {
	newStartState := &state{next: make(map[string][]*state)}
	newStartState.next[emptyWord] = []*state{a.startState}
	newStartState.isTerm = true

	for _, term := range a.terminals {
		term.next[emptyWord] = append(term.next[emptyWord], newStartState)
	}

	a.terminals = append(a.terminals, newStartState)
	a.startState = newStartState

	return nil
}

func (a *Automate) Concat(other *Automate) error {
	if len(a.startState.next) == 0 {
		return a.putStartState(other)
	}
	for _, term := range a.terminals {
		term.isTerm = false
		term.next[emptyWord] = append(term.next[emptyWord], other.startState)
	}

	a.terminals = other.terminals
	return nil
}

func (a *Automate) Join(other *Automate) error {
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

func (a *Automate) Full() error {
	stock := &state{next: make(map[string][]*state)}

	alphabet := make(map[string]bool)
	used := make(map[*state]bool)
	if err := traversal(a.startState, nil, nil, nil, &used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}

	for letter := range alphabet {
		stock.next[letter] = []*state{stock}
	}

	used = make(map[*state]bool)

	return traversal(a.startState, stock, nil, nil, &used, &fuller{alphabet: &alphabet})
}

func (a *Automate) Invert() error {
	newTerminals := make([]*state, 0)
	used := make(map[*state]bool)

	err := traversal(a.startState, nil, &newTerminals, nil, &used, &inverter{})
	if err != nil {
		return err
	}

	a.terminals = newTerminals
	return nil
}

func (a *Automate) Minimize() error {
	prevClasses := make(map[*state]int)
	curClasses := make(map[*state]int)
	alphabet := make(map[string]bool)
	used := make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, &prevClasses, &used, &getterTerminals{}); err != nil {
		return err
	}

	used = make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, nil, &used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}

	alphabetArray := make([]string, 0)
	for letter := range alphabet {
		alphabetArray = append(alphabetArray, letter)
	}

	states := make([]*state, 0)
	used = make(map[*state]bool)
	if err := traversal(a.startState, nil, &states, nil, &used, &getterStates{}); err != nil {
		return err
	}

	var counter int
	for {
		counter = 0

		for i, st := range states {
			index := -1
			for j := 0; j < i; j++ {
				if prevClasses[states[j]] != prevClasses[st] {
					continue
				}
				index = j
				for _, letter := range alphabetArray {
					if prevClasses[states[j].next[letter][0]] != prevClasses[st.next[letter][0]] {
						index = -1
					}
				}
				if index != -1 {
					break
				}
			}
			switch index {
			case -1:
				curClasses[st] = counter
				counter++
			default:
				curClasses[st] = curClasses[states[index]]
			}
		}

		isEqual := true
		for k, v := range prevClasses {
			if v != curClasses[k] {
				isEqual = false
			}
		}

		if len(curClasses) == len(prevClasses) && isEqual {
			break
		}

		prevClasses = make(map[*state]int)
		for k, v := range curClasses {
			prevClasses[k] = v
		}
	}

	newStates := make([]*state, counter)
	for i := 0; i < counter; i++ {
		newStates[i] = &state{next: make(map[string][]*state)}
	}

	terminals := make(map[*state]bool)
	for key, val := range curClasses {
		if key.isTerm {
			newStates[val].isTerm = true

			_, ok := terminals[newStates[val]]
			if !ok {
				terminals[newStates[val]] = true
			}
		}

		for k, v := range key.next {
			for _, st := range v {
				addNewTransition(newStates[val], newStates[curClasses[st]], k)
			}
		}
	}

	a.startState = newStates[curClasses[a.startState]]

	a.terminals = make([]*state, 0)
	for k := range terminals {
		a.terminals = append(a.terminals, k)
	}

	return nil
}

func addExpr(to string, first string, second string) string {
	switch {
	case first == "", len(first) == 1:
		to += first
	default:
		to += "("
		to += first
		to += ")"
		to += second
	}

	return to
}

func (a *Automate) GetRegularExpr() (string, error) {
	newTerm := &state{isTerm: true}

	for _, term := range a.terminals {
		term.isTerm = false
		term.next[""] = append(term.next[""], newTerm)
	}

	used := make(map[*state]bool)
	states := make([]*state, 0)
	err := traversal(a.startState, nil, &states, nil, &used, &getterStates{})
	if err != nil {
		return "", err
	}

	for {
		var (
			exit  = true
			st    *state
			index int
		)

		for i, s := range states {
			if !s.isTerm && a.startState != s {
				st = s
				exit = false
				index = i
				break
			}
		}

		for _, s := range states {
			err := removeDublicates(s)
			if err != nil {
				return "", err
			}
		}

		if exit || len(states) == 2 {
			break
		}

		expr, err := proccessSelfTransitions(st)
		if err != nil {
			return "", err
		}

		for _, stat := range states {
			err := removeDublicates(stat)
			if err != nil {
				return "", err
			}

			toAdd := make(map[string][]*state)
			for key, val := range stat.next {
				toDelete := -1
				for _, v := range val {
					if v == st {
						for k, node := range v.next {
							for _, n := range node {
								var ex string
								ex = addExpr(ex, key, "")
								ex = addExpr(ex, expr, "*")
								ex = addExpr(ex, k, "")

								_, ok := toAdd[ex]
								switch ok {
								case true:
									toAdd[ex] = append(toAdd[ex], n)
								case false:
									toAdd[ex] = []*state{n}
								}
							}
						}
					}
				}

				for j, n := range stat.next[key] {
					if n == st {
						toDelete = j
					}
				}

				if toDelete != -1 {
					if len(stat.next[key]) == 1 {
						stat.next[key] = nil
					} else {
						stat.next[key] = append(stat.next[key][:toDelete], stat.next[key][(toDelete+1):]...)
					}

					if len(stat.next[key]) == 0 {
						delete(stat.next, key)
					}
				}
			}

			for exp, node := range toAdd {
				for _, nn := range node {
					_, ok := stat.next[exp]
					if !ok {
						stat.next[exp] = []*state{nn}
					}
				}
			}
		}

		states = append(states[:index], states[(index+1):]...)
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

	var ans string
	for key, val := range a.startState.next {
		for _, v := range val {
			if v == a.startState && key != "" {
				ans += "(" + key + ")" + "*"
			} else {
				ans += key
			}
		}
	}

	return ans, nil
}

type state struct {
	next   map[string][]*state
	isTerm bool
}
