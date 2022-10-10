package customerrors

import "errors"

var (
	ErrNotImplemented        = errors.New("not implemented")
	ErrNoAutomate            = errors.New("no automate")
	ErrNoSuchWord            = errors.New("no such word")
	ErrInvalidFormat         = errors.New("invalid format")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrIncorrectMapping      = errors.New("incorrect mapping")
	ErrNotFull               = errors.New("automata is not full")
	ErrNotDetermine          = errors.New("automata is not determine")
	ErrHasEpsilonTransitions = errors.New("automata has epsilon transitions")
)
