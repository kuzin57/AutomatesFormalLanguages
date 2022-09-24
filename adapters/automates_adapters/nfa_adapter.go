package automatesadapters

import (
	"strings"
	"workspace/adapters"
	"workspace/internal/automate"
	customerrors "workspace/internal/errors"
)

type nfaAutomateAdapter struct {
	automate automate.Automate
	name     string
}

func (a *nfaAutomateAdapter) Get() (automate.Automate, error) {
	if a.automate == nil {
		return nil, customerrors.ErrNotImplemented
	}
	return a.automate, nil
}

func (a *nfaAutomateAdapter) Create(name string, words []string) (err error) {
	a.automate = automate.NewNFA()
	for _, word := range words {
		switch {
		case !strings.Contains(word, "*"):
			err = a.automate.AddNewWord(word)
			if err != nil {
				return
			}

		default:
			newAutomate := automate.NewNFA()
			parts := strings.Split(word, "*")

			for _, part := range parts {
				part = part[1:]
				part = part[:len(part)-1]
				newAutomate.AddNewWord(part)
			}

			a.automate.Join(newAutomate)
		}
	}
	return
}

func (a *nfaAutomateAdapter) AddStar() error {
	err := a.automate.Cycle()
	if err != nil {
		return err
	}
	return nil
}

func (a *nfaAutomateAdapter) Join(other adapters.AutomateAdapter) error {
	realNFAAdapter := other.(*nfaAutomateAdapter)
	if a.automate == nil || realNFAAdapter.automate == nil {
		return customerrors.ErrNoAutomate
	}

	return a.automate.Join(realNFAAdapter.automate)
}

func (a *nfaAutomateAdapter) Read(word string) bool {
	return a.automate.Read(word)
}

func (a *nfaAutomateAdapter) SetName(name string) {
	a.name = name
}

func (a *nfaAutomateAdapter) GetName() string {
	return a.name
}
