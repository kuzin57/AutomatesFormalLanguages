package fabric

import (
	"strings"
	automata "workspace/internal/entities/automata"
	"workspace/internal/entities/parser"
)

type automataFabric struct {
	parser *parser.Parser
}

func NewAutomataFabric(parser *parser.Parser) (AutomataFabric, error) {
	return &automataFabric{parser: parser}, nil
}

func (f *automataFabric) Create() (*automata.Automata, error) {
	return create(f.parser)
}

func create(parser *parser.Parser) (*automata.Automata, error) {
	if parser == nil {
		return nil, nil
	}

	if !strings.Contains(parser.Token, "+") && !strings.Contains(parser.Token, ".") && !strings.Contains(parser.Token, "*") {
		ans := automata.NewAutomata()
		ans.AddNewWord(parser.Token)
		return ans, nil
	}

	leftAutomata, err := create(parser.ChildLeft)
	if err != nil {
		return nil, err
	}

	rightAutomata, err := create(parser.ChildRight)
	if err != nil {
		return nil, err
	}

	switch parser.Operation {
	case '+':
		leftAutomata.Join(rightAutomata)
	case '.':
		leftAutomata.Concat(rightAutomata)
	case '*':
		leftAutomata.Cycle()
	}

	return leftAutomata, nil
}
