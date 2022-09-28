package showactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
	"workspace/internal/display"
	"workspace/internal/entities/automate"
)

type ShowStatesParams struct {
	Adapter adapters.AutomateAdapter
}

type ShowStatesAction struct {
	actions.BaseAction
}

func NewShowStatesAction(params *ShowStatesParams) (*ShowStatesAction, error) {
	return &ShowStatesAction{BaseAction: actions.NewBaseAction(params.Adapter)}, nil
}

func (a *ShowStatesAction) Do() {
	var states []*automate.State

	if states, a.Error = a.Adapter.GetStates(); a.CheckErr() {
		return
	}

	if a.Error = display.DisplayGraph(states, a.Adapter.GetName()); a.CheckErr() {
		return
	}
}
