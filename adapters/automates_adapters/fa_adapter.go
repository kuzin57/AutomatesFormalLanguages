package automatesadapters

import (
	"fmt"
	"workspace/adapters"
	"workspace/internal/automate"
	customerrors "workspace/internal/errors"
)

type faAutomateAdapter struct {
	automate automate.Automate
	name     string
}

func NewFAadapter(other adapters.AutomateAdapter, name string) (adapters.AutomateAdapter, error) {
	var err error
	realAutomate, ok := other.(*nfaAutomateAdapter)
	if !ok {
		return nil, fmt.Errorf("nfa adapter expected")
	}

	res := &faAutomateAdapter{name: name}
	res.automate, err = automate.NewFAFromNFA(realAutomate.automate)
	if err != nil {
		return nil, err
	}

	ans := &faAutomateAdapter{}
	ans.automate = res.automate
	ans.name = name

	return ans, nil
}

func (a *faAutomateAdapter) Get() (automate.Automate, error) {
	return nil, customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) GetStates() ([]*automate.State, error) {
	return a.automate.GetStates()
}

func (a *faAutomateAdapter) Create(string, string) error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) AddStar() error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) Join(adapters.AutomateAdapter) error {
	return customerrors.ErrNotImplemented
}

func (a *faAutomateAdapter) Read(word string) error {
	return a.automate.Read(word)
}

func (a *faAutomateAdapter) GetName() string {
	return a.name
}

func (a *faAutomateAdapter) SetName(name string) {
	a.name = name
}

func (a *faAutomateAdapter) DeleteEps() error {
	return customerrors.ErrNotImplemented
}
