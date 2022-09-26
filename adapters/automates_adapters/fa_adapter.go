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

func NewFAadapter(other adapters.AutomateAdapter, name string) (*nfaAutomateAdapter, error) {
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

	a, _ := automate.NewNFAFromFA(res.automate)
	ans := &nfaAutomateAdapter{}
	ans.name = realAutomate.name
	ans.automate = a

	return ans, nil
}

func (a *faAutomateAdapter) Get() (automate.Automate, error) {
	return nil, customerrors.ErrNotImplemented
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
	return customerrors.ErrNotImplemented
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
