package customerrors

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNoAutomate     = errors.New("no automate")
	ErrNoSuchWord     = errors.New("no such word")
	ErrInvalidFormat  = errors.New("invalid format")
)
