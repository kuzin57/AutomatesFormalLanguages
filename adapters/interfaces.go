package adapters

import (
	automata "workspace/internal/entities/automata"
	"workspace/internal/entities/parser"
)

type AutomataAdapter interface {
	Get() (*automata.Automata, error)
	Create(string, *parser.Parser) error
	AddStar() error
	Join(AutomataAdapter) error
	Read(string) error
	GetName() string
	SetName(string)
	DeleteEps() error
	GetStates() ([]*automata.State, error)
	Determine() error
	MakeFull() error
	Invert() error
	Minimize() error
	GetRegularExpression() (string, error)
}
