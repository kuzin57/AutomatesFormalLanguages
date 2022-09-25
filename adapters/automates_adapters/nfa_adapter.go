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

func (a *nfaAutomateAdapter) Create(name string, line string) (err error) {
	a.automate = automate.NewNFA()

	var (
		subexprs     []string
		needCircling bool
	)

	if line[len(line)-1] == '*' {
		needCircling = true
		line = line[:len(line)-1]
	}

	subexprs, err = a.parseLine(line)
	for _, subexpr := range subexprs {
		if !strings.Contains(subexpr, "+") {
			switch {
			case !strings.Contains(subexpr, "*"):
				err = a.automate.AddNewWord(subexpr)
				if err != nil {
					return
				}

			default:
				newAutomate := automate.NewNFA()
				parts := strings.Split(subexpr, "*")
				parts = parts[:len(parts)-1]
				for _, part := range parts {
					circledAutomate := automate.NewNFA()
					part = part[1:]
					part = part[:len(part)-1]

					var subparts []string
					subparts, err = a.parseLine(part)
					if err != nil {
						return
					}

					for _, p := range subparts {
						if err = circledAutomate.AddNewWord(p); err != nil {
							return
						}
					}

					if err = circledAutomate.Cycle(); err != nil {
						return
					}
					newAutomate.Join(circledAutomate)
				}

				a.automate.Join(newAutomate)
			}
			continue
		}
		var words []string
		words, err = a.parseLine(subexpr)
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
	}

	if needCircling {
		err = a.automate.Cycle()
	}
	return
}

func (a *nfaAutomateAdapter) parseLine(line string) ([]string, error) {
	var (
		balance int
		curWord string
		ans     []string
	)

	for i, char := range line {
		switch {
		case char == '(':
			balance++
			if i > 0 {
				curWord += string(char)
			}
		case char == ')':
			balance--
			if i < len(line)-1 {
				curWord += string(char)
			}
		case char == '+' && balance == 1:
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
