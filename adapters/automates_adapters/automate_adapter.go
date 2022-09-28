package automateadapter

import (
	"workspace/adapters"
	"workspace/internal/entities/automate"
	"workspace/internal/entities/parser"
	customerrors "workspace/internal/errors"
	"workspace/internal/fabric"
)

type automateAdapter struct {
	automate *automate.Automate
	name     string
}

func (a *automateAdapter) Get() (*automate.Automate, error) {
	if a.automate == nil {
		return nil, customerrors.ErrNotImplemented
	}
	return a.automate, nil
}

func (a *automateAdapter) Create(name string, parser *parser.Parser) (err error) {
	automateFabric, err := fabric.NewAutomateFabric(parser)
	if err != nil {
		return
	}

	a.automate, err = automateFabric.Create()
	if err != nil {
		return
	}

	return
}

func (a *automateAdapter) Determine() error {
	res, err := automate.DetermineAutomate(a.automate)
	if err != nil {
		return err
	}

	a.automate = res
	return nil
}

func (a *automateAdapter) AddStar() error {
	err := a.automate.Cycle()
	if err != nil {
		return err
	}
	return nil
}

func (a *automateAdapter) Join(other adapters.AutomateAdapter) error {
	realNFAAdapter := other.(*automateAdapter)
	if a.automate == nil || realNFAAdapter.automate == nil {
		return customerrors.ErrNoAutomate
	}

	return a.automate.Join(realNFAAdapter.automate)
}

func (a *automateAdapter) Read(word string) error {
	return a.automate.Read(word)
}

func (a *automateAdapter) SetName(name string) {
	a.name = name
}

func (a *automateAdapter) GetName() string {
	return a.name
}

func (a *automateAdapter) DeleteEps() error {
	return a.automate.DeleteEps()
}

func (a *automateAdapter) GetStates() ([]*automate.State, error) {
	return a.automate.GetStates()
}

func (a *automateAdapter) MakeFull() error {
	return a.automate.Full()
}

func (a *automateAdapter) Invert() error {
	return a.automate.Invert()
}
