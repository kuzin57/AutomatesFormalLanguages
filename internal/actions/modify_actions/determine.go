package modifyactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
)

type DetermineParams struct {
	NFA adapters.AutomateAdapter
}

type DetermineAction struct {
	actions.BaseAction
}

func NewDetermineAction(params *DetermineParams) (*DetermineAction, error) {
	return &DetermineAction{BaseAction: actions.NewBaseAction(params.NFA)}, nil
}

func (a *DetermineAction) Do() {
	a.Error = a.Adapter.Determine()
}
