package parser

var priorities = map[rune]int{
	'+': 1,
	'.': 2,
	'*': 3,
}

var operations = map[int]rune{
	1: '+',
	2: '.',
	3: '*',
}
