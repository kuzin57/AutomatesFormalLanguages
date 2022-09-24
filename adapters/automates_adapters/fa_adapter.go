package automatesadapters

import (
	"workspace/adapters"
	"workspace/internal/automate"
	customerrors "workspace/internal/errors"
)

type faAutomateAdapter struct {
	automate automate.Automate
	name     string
}

func (a *faAutomateAdapter) Get() (automate.Automate, error) {
	return nil, customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) Create(string, []string) error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) AddStar() error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) Join(adapters.AutomateAdapter) error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) Read(word string) bool {
	return false
}

func (a *faAutomateAdapter) GetName() string {
	return a.name
}

func (a *faAutomateAdapter) SetName(name string) {
	a.name = name
}
