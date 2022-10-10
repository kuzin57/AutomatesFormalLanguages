package automata_test

import (
	"testing"
	"workspace/internal/entities/automata"
	"workspace/internal/entities/parser"
	"workspace/internal/fabric"

	"github.com/stretchr/testify/assert"
)

var (
	regularExpressions = []string{
		"(a.a+b.b+(a.b+b.a).(b.b+a.a)*.(b.a+a.b))*",
		"a.c",
		"a*",
		"a.b.c",
		"(a.a.b+a+a.b.a+b.b.a)*",
		"(a*.b)*",
		"(a.b+b.a+a.a.b)*",
	}
)

func SetupAutomata(regularExpression string) (*automata.Automata, error) {
	parser := parser.NewParser(regularExpression, nil)
	parser.Parse()

	fabric, err := fabric.NewAutomataFabric(parser)
	if err != nil {
		return nil, err
	}

	automata, err := fabric.Create()
	if err != nil {
		return nil, err
	}
	return automata, nil
}

func TestNewAutomata(t *testing.T) {
	automata := automata.NewAutomata()
	assert.NotNil(t, automata)
}

func TestRead(t *testing.T) {
	var (
		testCasesSuccess = []string{"aabb", "aaaabbbb", "abababab", "aa", ""}
		testCasesFail    = []string{"a", "b", "aab", "aaabb", "ababab"}
	)

	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, automata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, automata.Read(testCase))
	}
}

func TestDeleteEps(t *testing.T) {
	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	assert.NoError(t, automata.DeleteEps())
	assert.NoError(t, automata.CheckNoEpsilon())
}

func TestDetermine(t *testing.T) {
	var (
		testCasesSuccess = []string{"aa", "aabb", "abba", "bbbb", "abababababab"}
		testCasesFail    = []string{"aab", "bbbaaaa", "aaa", "abababa"}
	)

	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	err = automata.DeleteEps()
	assert.NoError(t, err)

	detAutomata, err := automata.Determine()
	assert.NoError(t, err)
	assert.NotNil(t, detAutomata)

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, automata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, automata.Read(testCase))
	}

	assert.NoError(t, detAutomata.CheckDetermine())
}

func TestGetStates(t *testing.T) {
	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)
	assert.NoError(t, automata.CheckStates())

	states, err := automata.GetStates()
	assert.NotEmpty(t, states)
	assert.NoError(t, err)
}

func TestConcat(t *testing.T) {
	var (
		testCasesSuccess = []string{"aabbaaac", "ababababac", "ac", "ababac"}
		testCasesFail    = []string{"aabac", "aac", "abababa", "abaababc"}
	)

	firstAutomata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, firstAutomata)

	secondAutomata, err := SetupAutomata(regularExpressions[1])
	assert.NoError(t, err)
	assert.NotNil(t, secondAutomata)

	assert.NoError(t, firstAutomata.Concat(secondAutomata))

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, firstAutomata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, firstAutomata.Read(testCase))
	}
}

func TestJoin(t *testing.T) {
	var (
		testCasesSuccess = []string{"abababab", "aabb", "aa", "abbbba", "aaaaa", "aaa", "a", "abba"}
		testCasesFail    = []string{"abbb", "aaab", "aaaaab", "bbbbbaa"}
	)

	firstAutomata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, firstAutomata)

	secondAutomata, err := SetupAutomata(regularExpressions[2])
	assert.NoError(t, err)
	assert.NotNil(t, secondAutomata)

	assert.NoError(t, firstAutomata.Join(secondAutomata))

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, firstAutomata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, firstAutomata.Read(testCase))
	}
}

func TestCycle(t *testing.T) {
	var (
		testCasesSuccess = []string{"abcabcabc", "abc", "", "abcabc", "abcabcabcabcabcabcabcabcabcabcabcabcabc", "abcabcabcabcabc"}
		testCasesFail    = []string{"ab", "cccc", "aaaaab", "bbbbbaa"}
	)

	automata, err := SetupAutomata(regularExpressions[3])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	assert.NoError(t, automata.Cycle())

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, automata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, automata.Read(testCase))
	}
}

func TestFull(t *testing.T) {
	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	assert.NoError(t, automata.DeleteEps())

	detAutomata, err := automata.Determine()
	assert.NoError(t, err)
	assert.NotNil(t, detAutomata)

	assert.NoError(t, detAutomata.Full())
	assert.NoError(t, detAutomata.CheckFull())
}

func TestInvert(t *testing.T) {
	var (
		testCasesSuccess = []string{"aaa", "abbba", "b", "aaabbbb", "abababa"}
		testCasesFail    = []string{"aaaa", "abba", "aaaabbbb", "aabb"}
	)

	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	assert.NoError(t, automata.DeleteEps())
	assert.NoError(t, automata.Invert())

	for _, testCase := range testCasesSuccess {
		assert.NoError(t, automata.Read(testCase))
	}

	for _, testCase := range testCasesFail {
		assert.Error(t, automata.Read(testCase))
	}
}

func TestMinimize(t *testing.T) {
	automata, err := SetupAutomata(regularExpressions[0])
	assert.NoError(t, err)
	assert.NotNil(t, automata)

	assert.NoError(t, automata.DeleteEps())

	detAutomata, err := automata.Determine()
	assert.NoError(t, err)
	assert.NotNil(t, detAutomata)

	assert.NoError(t, detAutomata.Full())
	assert.NoError(t, detAutomata.Minimize())

	assert.NoError(t, detAutomata.CheckDetermine())
	assert.NoError(t, detAutomata.CheckFull())

	states, err := detAutomata.GetStates()
	assert.NoError(t, err)

	assert.Equal(t, 4, len(states))
}

func TestGetRegularExpression(t *testing.T) {
	testCases := []string{
		regularExpressions[0],
		regularExpressions[4],
		regularExpressions[5],
		regularExpressions[6],
	}
	for _, tC := range testCases {
		t.Run(tC, func(t *testing.T) {
			automata, err := SetupAutomata(tC)
			assert.NoError(t, err)
			assert.NotNil(t, automata)

			assert.NoError(t, automata.DeleteEps())

			detAutomata, err := automata.Determine()
			assert.NoError(t, err)
			assert.NotNil(t, detAutomata)

			assert.NoError(t, detAutomata.Full())
			assert.NoError(t, detAutomata.Minimize())

			statesBeforeRegularExpression, err := detAutomata.GetStates()
			assert.NoError(t, err)
			assert.NotEmpty(t, statesBeforeRegularExpression)

			regularExpression, err := detAutomata.GetRegularExpression()
			assert.NotEmpty(t, regularExpression)
			assert.NoError(t, err)

			newAutomata, err := SetupAutomata(regularExpression)
			assert.NoError(t, err)
			assert.NotNil(t, newAutomata)

			assert.NoError(t, newAutomata.DeleteEps())

			newDetAutomata, err := newAutomata.Determine()
			assert.NoError(t, err)
			assert.NotNil(t, newDetAutomata)

			assert.NoError(t, newDetAutomata.Full())
			assert.NoError(t, newDetAutomata.Minimize())

			statesAfterRegularExpression, err := newDetAutomata.GetStates()
			assert.NoError(t, err)
			assert.NotEmpty(t, statesAfterRegularExpression)

			assert.Equal(t, len(statesBeforeRegularExpression), len(statesAfterRegularExpression))
		})
	}
}
