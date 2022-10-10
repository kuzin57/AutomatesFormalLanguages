package automata

type Executor interface {
	Execute(*state) error
}

type Checker interface {
	Check(*state) error
}
