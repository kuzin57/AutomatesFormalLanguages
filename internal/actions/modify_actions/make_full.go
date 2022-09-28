package modifyactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
)

type MakeFullParams struct {
	Adapter adapters.AutomateAdapter
}

type MakeFullAction struct {
	actions.BaseAction
}

func NewMakeFullAction(params *MakeFullParams) (*MakeFullAction, error) {
	return &MakeFullAction{BaseAction: actions.NewBaseAction(params.Adapter)}, nil
}

func (a *MakeFullAction) Do() {
	a.Error = a.Adapter.MakeFull()
}
