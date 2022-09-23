package customerrors

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNoAutomate     = errors.New("no automate")
)
