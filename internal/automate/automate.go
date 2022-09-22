package automate

type Automate interface {
	DeleteEps()
	Read(string) (string, error)
	Check() bool
}

type nfa struct {
	startState *state
}

func (a *nfa) DeleteEps() {

}

func (a *nfa) Check() bool {
	return false
}

func (a *nfa) Read(line string) (string, error) {
	return "", nil
}

type fa struct {
	startState *state
}

func (a *fa) DeleteEps() {

}

func (a *fa) Read(line string) (string, error) {
	return "", nil
}

func (a *fa) Check() bool {
	return true
}

type state struct {
	transiotions []transition
	isTerm       bool
}

type transition struct {
	to   state
	word string
}
