package automate

import (
	"fmt"
	"reflect"
)

func NewAutomate() *Automate {
	return &Automate{startState: &state{isTerm: false, next: make(map[rune][]*state)}}
}

func getSetsTranisitions(a *Automate) (*[]*map[*state]bool, *map[int]map[rune]int, error) {
	alphabet := make(map[rune]bool)
	used := make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, nil, &used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return nil, nil, err
	}

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
	transitions *map[int]map[rune]int,
	sets *[]*map[*state]bool,
	created *map[int]*state,
	index int,
) *state {
	res := &state{next: make(map[rune][]*state)}
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

func (a *Automate) GetStates() ([]*State, error) {
	used := make(map[*state]int)
	ans := make([]*State, 0)

	var counter int
	err := mapStates(a.startState, &used, &ans, &counter)

	return ans, err
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

	newState := &state{next: make(map[rune][]*state)}
	newState.next[emptyWord] = make([]*state, 2)
	newState.next[emptyWord][0] = a.startState
	newState.next[emptyWord][1] = other.startState

	a.startState = newState
	a.terminals = append(a.terminals, other.terminals...)
	return nil
}

func (a *Automate) Full() error {
	stock := &state{next: make(map[rune][]*state)}

	alphabet := make(map[rune]bool)
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
	type tuple struct {
		prevClass        int
		neighbourClasses []int
	}

	cmpTuples := func(first *tuple, second *tuple) bool {
		if first.prevClass != second.prevClass || len(first.neighbourClasses) != len(second.neighbourClasses) {
			return false
		}

		for i := 0; i < len(first.neighbourClasses); i++ {
			if first.neighbourClasses[i] != second.neighbourClasses[i] {
				return false
			}
		}
		return true
	}

	prevClasses := make(map[*state]int)
	curClasses := make(map[*state]int)
	alphabet := make(map[rune]bool)
	used := make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, &prevClasses, &used, &getterTerminals{}); err != nil {
		return err
	}

	used = make(map[*state]bool)

	if err := traversal(a.startState, nil, nil, nil, &used, &getterAlphabet{alphabet: &alphabet}); err != nil {
		return err
	}

	alphabetArray := make([]rune, 0)
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
		tuples := make([]tuple, len(states))
		for i, st := range states {
			tuples[i] = tuple{prevClass: prevClasses[st], neighbourClasses: make([]int, len(alphabetArray))}
			for j, letter := range alphabetArray {
				tuples[i].neighbourClasses[j] = prevClasses[st.next[letter][0]]
			}
		}

		counter = 0
		for i, tuple := range tuples {
			index := -1
			for j := 0; j < i; j++ {
				if cmpTuples(&tuple, &tuples[j]) {
					index = j
					break
				}
			}

			switch index {
			case -1:
				curClasses[states[i]] = counter
				counter++
			default:
				curClasses[states[i]] = index
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

		prevClasses = curClasses
	}

	fmt.Println("classes", curClasses)
	fmt.Println("hehe", prevClasses)

	newStates := make([]*state, counter)
	for i := 0; i < counter; i++ {
		newStates[i] = &state{next: make(map[rune][]*state)}
	}

	for key, val := range curClasses {
		if key.isTerm {
			newStates[val].isTerm = true
		}

		for k, v := range key.next {
			for _, st := range v {
				newStates[val].next[k] = append(newStates[val].next[k], newStates[curClasses[st]])
			}
		}
	}

	a.startState = newStates[curClasses[a.startState]]

	return nil
}

type state struct {
	next   map[rune][]*state
	isTerm bool
}
