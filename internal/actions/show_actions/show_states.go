package showactions

import (
	"fmt"
	"strings"
	"workspace/adapters"
	"workspace/internal/actions"
	"workspace/internal/automate"

	"github.com/olekukonko/tablewriter"
)

type ShowStatesParams struct {
	Adapter adapters.AutomateAdapter
	Buffer  *strings.Builder
}

type ShowStatesAction struct {
	actions.BaseAction

	buffer *strings.Builder
}

func NewShowStatesAction(params *ShowStatesParams) (*ShowStatesAction, error) {
	return &ShowStatesAction{BaseAction: actions.NewBaseAction(params.Adapter), buffer: params.Buffer}, nil
}

func (a *ShowStatesAction) Do() {
	var states []*automate.State

	if states, a.Error = a.Adapter.GetStates(); a.CheckErr() {
		return
	}

	table := tablewriter.NewWriter(a.buffer)
	table.SetBorder(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetColumnSeparator(":")

	headers := []string{"From", "To", "Letter"}
	table.SetHeader(headers)

	var transitions [][]string
	for _, state := range states {
		for key, val := range state.Transitions {
			for _, v := range val {
				transitions = append(transitions, []string{fmt.Sprint(state.Number), fmt.Sprint(v), string(key)})
			}
		}
	}

	table.AppendBulk(transitions)
	table.Render()
}
