package parser

import (
	"strings"
	"unicode"
)

type Parser struct {
	Lec        string
	Operation  rune
	Parent     *Parser
	ChildLeft  *Parser
	ChildRight *Parser
}

func NewParser(expr string, parent *Parser) *Parser {
	return &Parser{Lec: expr, Parent: parent}
}

func (t *Parser) Parse() {
	var (
		index       int
		minPriority int
		balance     int
	)

	if !strings.Contains(t.Lec, "+") &&
		!strings.Contains(t.Lec, ".") &&
		!strings.Contains(t.Lec, "*") {
		return
	}

	minPriority = 4

	for i, char := range t.Lec {
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
	case 1, 2:
		t.ChildLeft = NewParser(t.Lec[:index], t)
		t.ChildRight = NewParser(t.Lec[(index+1):], t)

		t.ChildLeft.Parse()
		t.ChildRight.Parse()
	case 3:
		t.ChildLeft = NewParser(t.Lec[:index], t)
		t.ChildLeft.Parse()
	case 4:
		t.ChildLeft = NewParser(t.Lec[1:len(t.Lec)-1], t)
		t.ChildLeft.Parse()
	}

	t.Operation = operations[minPriority]
}

func (t *Parser) Print() {
	if t == nil {
		return
	}
	if !strings.Contains(t.Lec, "+") &&
		!strings.Contains(t.Lec, ".") &&
		!strings.Contains(t.Lec, "*") {
		return
	}

	t.ChildLeft.Print()
	t.ChildRight.Print()
}
