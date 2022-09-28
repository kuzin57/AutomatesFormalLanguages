package automate

type State struct {
	Number      int
	Transitions map[rune][]int
	IsTerminal  bool
}
