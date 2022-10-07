package automata

type Executor interface {
	Execute(args ...any) error
}
