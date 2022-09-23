package adapters

import (
	"fmt"
	"strings"
	"workspace/internal/automate"
	"workspace/internal/config"
	customerrors "workspace/internal/errors"
)

type AutomateAdapter interface {
	Get() (automate.Automate, error)
	Create(string, []string) error
	AddStar() error
	Join(AutomateAdapter) error
}

type nfaAutomateAdapter struct {
	automate automate.Automate
}

func NewAdapter(cfg config.AdaptersConfig) AutomateAdapter {
	if cfg.IsDeterministic {
		return &faAutomateAdapter{}
	}
	return &nfaAutomateAdapter{}
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
				fmt.Println("part", part)
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

func (a *nfaAutomateAdapter) Join(other AutomateAdapter) error {
	realNFAAdapter := other.(*nfaAutomateAdapter)
	if a.automate == nil || realNFAAdapter.automate == nil {
		return customerrors.ErrNoAutomate
	}

	a.automate.Join(realNFAAdapter.automate)
	return nil
}

type faAutomateAdapter struct {
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

func (a *faAutomateAdapter) Join(AutomateAdapter) error {
	return customerrors.ErrNotImplemented
}
