package adapters

import (
	"workspace/internal/entities/automate"
	"workspace/internal/entities/parser"
)

type AutomateAdapter interface {
	Get() (*automate.Automate, error)
	Create(string, *parser.Parser) error
	AddStar() error
	Join(AutomateAdapter) error
	Read(string) error
	GetName() string
	SetName(string)
	DeleteEps() error
	GetStates() ([]*automate.State, error)
	Determine() error
	MakeFull() error
	Invert() error
	Minimize() error
}
