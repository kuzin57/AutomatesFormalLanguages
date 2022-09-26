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

func (a *nfaAutomateAdapter) Create(name string, expr string) (err error) {
	a.automate = automate.NewNFA()

	switch {
	case !a.check(expr, '+') && a.check(expr, '.'):
		var regulars []string
		regulars, err = a.parseExpr(expr, '.')
		if err != nil {
			return
		}

		auto := automate.NewNFA()
		for _, regular := range regulars {
			adapter := nfaAutomateAdapter{automate: automate.NewNFA()}
			err = adapter.Create("", regular)
			if err != nil {
				return
			}
			auto.Concat(adapter.automate)
		}

		a.automate.Join(auto)

	case a.check(expr, '+'):
		var words []string
		words, err = a.parseExpr(expr, '+')
		if err != nil {
			return
		}

		newAuto := automate.NewNFA()
		for _, word := range words {
			newAdapter := nfaAutomateAdapter{automate: automate.NewNFA()}
			newAdapter.Create("", word)
			newAuto.Join(newAdapter.automate)
		}
		a.automate.Join(newAuto)

	case len(expr) >= 2 && expr[len(expr)-2] == '*':
		newAdapter := nfaAutomateAdapter{automate: automate.NewNFA()}
		expr = expr[1:]
		expr = expr[:len(expr)-2]

		if err = newAdapter.Create("", expr); err != nil {
			return
		}

		if err = newAdapter.automate.Cycle(); err != nil {
			return
		}
		a.automate.Join(newAdapter.automate)

	case !strings.Contains(expr, "+"):
		expr = expr[1:]
		expr = expr[:len(expr)-1]
		err = a.automate.AddNewWord(expr)
		if err != nil {
			return
		}

	}

	return
}

func (a *nfaAutomateAdapter) parseExpr(expr string, sep rune) ([]string, error) {
	var (
		balance int
		curWord string
		ans     []string
	)

	for i, char := range expr {
		switch {
		case char == '(':
			balance++
			if i > 0 {
				curWord += string(char)
			}
		case char == ')':
			balance--
			if i < len(expr)-1 {
				curWord += string(char)
			}
		case char == sep && balance == 1:
			ans = append(ans, curWord)
			curWord = ""
		default:
			curWord += string(char)
		}
	}

	switch {
	case balance == 0:
		ans = append(ans, curWord)
	default:
		return nil, customerrors.ErrInvalidFormat
	}

	return ans, nil
}

func (a *nfaAutomateAdapter) check(expr string, sep rune) bool {
	var balance int
	for _, char := range expr {
		switch {
		case char == '(':
			balance++
		case char == ')':
			balance--
		case char == sep && balance == 1:
			return true
		}
	}
	return false
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

func (a *nfaAutomateAdapter) Read(word string) error {
	return a.automate.Read(word)
}

func (a *nfaAutomateAdapter) SetName(name string) {
	a.name = name
}

func (a *nfaAutomateAdapter) GetName() string {
	return a.name
}

func (a *nfaAutomateAdapter) DeleteEps() error {
	return a.automate.DeleteEps()
}
