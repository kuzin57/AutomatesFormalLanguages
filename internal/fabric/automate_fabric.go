package fabric

import (
	"strings"
	"workspace/internal/entities/automate"
	"workspace/internal/entities/parser"
)

type automateFabric struct {
	parser *parser.Parser
}

func NewAutomateFabric(parser *parser.Parser) (AutomateFabric, error) {
	return &automateFabric{parser: parser}, nil
}

func (f *automateFabric) Create() (*automate.Automate, error) {
	return create(f.parser)
}

func create(parser *parser.Parser) (*automate.Automate, error) {
	if parser == nil {
		return nil, nil
	}

	if !strings.Contains(parser.Lec, "+") && !strings.Contains(parser.Lec, ".") && !strings.Contains(parser.Lec, "*") {
		ans := automate.NewAutomate()
		ans.AddNewWord(parser.Lec)
		return ans, nil
	}

	leftAutomate, err := create(parser.ChildLeft)
	if err != nil {
		return nil, err
	}

	rightAutomate, err := create(parser.ChildRight)
	if err != nil {
		return nil, err
	}

	switch parser.Operation {
	case '+':
		leftAutomate.Join(rightAutomate)
	case '.':
		leftAutomate.Concat(rightAutomate)
	case '*':
		leftAutomate.Cycle()
	}

	return leftAutomate, nil
}
