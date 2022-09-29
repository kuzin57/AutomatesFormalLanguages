package modifyactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
)

type MinimizeParams struct {
	Adapter adapters.AutomateAdapter
}

type MinimizeAction struct {
	actions.BaseAction
}

func NewMinimizeAction(params *MinimizeParams) (*MinimizeAction, error) {
	return &MinimizeAction{BaseAction: actions.NewBaseAction(params.Adapter)}, nil
}

func (a *MinimizeAction) Do() {
	a.Error = a.Adapter.Minimize()
}
