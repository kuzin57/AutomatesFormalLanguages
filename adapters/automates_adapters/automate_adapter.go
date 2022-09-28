package automateadapter

import (
	"strings"
	"workspace/adapters"
	"workspace/internal/entities/automate"
	customerrors "workspace/internal/errors"
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

func (a *automateAdapter) Create(name string, expr string) (err error) {
	a.automate = automate.NewAutomate()

	switch {
	case !a.check(expr, '+') && a.check(expr, '.'):
		var regulars []string
		regulars, err = a.parseExpr(expr, '.')
		if err != nil {
			return
		}

		auto := automate.NewAutomate()
		for _, regular := range regulars {
			adapter := automateAdapter{automate: automate.NewAutomate()}
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

		newAuto := automate.NewAutomate()
		for _, word := range words {
			newAdapter := automateAdapter{automate: automate.NewAutomate()}
			newAdapter.Create("", word)
			newAuto.Join(newAdapter.automate)
		}
		a.automate.Join(newAuto)

	case len(expr) >= 2 && expr[len(expr)-2] == '*':
		newAdapter := automateAdapter{automate: automate.NewAutomate()}
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

func (a *automateAdapter) parseExpr(expr string, sep rune) ([]string, error) {
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

func (a *automateAdapter) check(expr string, sep rune) bool {
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
