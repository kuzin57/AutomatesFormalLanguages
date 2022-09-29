package automate

type Executor interface {
	Execute(args ...any) error
}
