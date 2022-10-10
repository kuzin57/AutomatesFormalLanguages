package parser

import (
	"strings"
	"unicode"
)

type Parser struct {
	Token      string
	Operation  rune
	Parent     *Parser
	ChildLeft  *Parser
	ChildRight *Parser
}

func NewParser(expr string, parent *Parser) *Parser {
	return &Parser{Token: expr, Parent: parent}
}

func (t *Parser) Parse() {
	var (
		index       int
		minPriority operationPriority
		balance     int
	)

	if !strings.Contains(t.Token, "+") &&
		!strings.Contains(t.Token, ".") &&
		!strings.Contains(t.Token, "*") {
		return
	}

	minPriority = unspecifiedPriority

	for i, char := range t.Token {
		switch char {
		case '(':
			balance++
			continue
		case ')':
			balance--
			continue
		}

		priority, ok := priorities[char]
		if !unicode.IsLetter(char) &&
			ok && minPriority > priority && balance == 0 {
			index = i
			minPriority = priority
		}
	}

	switch minPriority {
	case plusPriority, concatPriority:
		t.ChildLeft = NewParser(t.Token[:index], t)
		t.ChildRight = NewParser(t.Token[(index+1):], t)

		t.ChildLeft.Parse()
		t.ChildRight.Parse()
	case starPriority:
		t.ChildLeft = NewParser(t.Token[:index], t)
		t.ChildLeft.Parse()
	case unspecifiedPriority:
		t.ChildLeft = NewParser(t.Token[1:len(t.Token)-1], t)
		t.ChildLeft.Parse()
	}

	t.Operation = operations[minPriority]
}

func (t *Parser) Print() {
	if t == nil {
		return
	}
	if !strings.Contains(t.Token, "+") &&
		!strings.Contains(t.Token, ".") &&
		!strings.Contains(t.Token, "*") {
		return
	}

	t.ChildLeft.Print()
	t.ChildRight.Print()
}
