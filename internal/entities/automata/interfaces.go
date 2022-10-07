package automata

type Executor interface {
	Execute(*state) error
}
