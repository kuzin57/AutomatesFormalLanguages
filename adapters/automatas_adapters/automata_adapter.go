package automataadapter

import (
	"workspace/adapters"
	automata "workspace/internal/entities/automata"
	"workspace/internal/entities/parser"
	customerrors "workspace/internal/errors"
	"workspace/internal/fabric"
)

type automataAdapter struct {
	automata *automata.Automata
	name     string
}

func (a *automataAdapter) Get() (*automata.Automata, error) {
	if a.automata == nil {
		return nil, customerrors.ErrNotImplemented
	}
	return a.automata, nil
}

func (a *automataAdapter) Create(name string, parser *parser.Parser) (err error) {
	automataFabric, err := fabric.NewAutomataFabric(parser)
	if err != nil {
		return
	}

	a.automata, err = automataFabric.Create()
	if err != nil {
		return
	}

	return
}

func (a *automataAdapter) Determine() error {
	res, err := a.automata.Determine()
	if err != nil {
		return err
	}

	a.automata = res
	return nil
}

func (a *automataAdapter) AddStar() error {
	err := a.automata.Cycle()
	if err != nil {
		return err
	}
	return nil
}

func (a *automataAdapter) Join(other adapters.AutomataAdapter) error {
	realNFAAdapter := other.(*automataAdapter)
	if a.automata == nil || realNFAAdapter.automata == nil {
		return customerrors.ErrNoAutomate
	}

	return a.automata.Join(realNFAAdapter.automata)
}

func (a *automataAdapter) Read(word string) error {
	return a.automata.Read(word)
}

func (a *automataAdapter) SetName(name string) {
	a.name = name
}

func (a *automataAdapter) GetName() string {
	return a.name
}

func (a *automataAdapter) DeleteEps() error {
	return a.automata.DeleteEps()
}

func (a *automataAdapter) GetStates() ([]*automata.State, error) {
	return a.automata.GetStates()
}

func (a *automataAdapter) MakeFull() error {
	return a.automata.Full()
}

func (a *automataAdapter) Invert() error {
	return a.automata.Invert()
}

func (a *automataAdapter) Minimize() error {
	return a.automata.Minimize()
}

func (a *automataAdapter) GetRegularExpression() (string, error) {
	return a.automata.GetRegularExpression()
}
