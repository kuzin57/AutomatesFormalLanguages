package parser

var priorities = map[rune]operationPriority{
	'+': plusPriority,
	'.': concatPriority,
	'*': starPriority,
}

var operations = map[operationPriority]rune{
	plusPriority:   '+',
	concatPriority: '.',
	starPriority:   '*',
}

type operationPriority uint8

const (
	plusPriority = iota + 1
	concatPriority
	starPriority
	unspecifiedPriority
)
