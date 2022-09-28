package modifyactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
)

type InvertParams struct {
	Adapter adapters.AutomateAdapter
}

type InvertAction struct {
	actions.BaseAction
}

func NewInvertAction(params *InvertParams) (*InvertAction, error) {
	return &InvertAction{BaseAction: actions.NewBaseAction(params.Adapter)}, nil
}

func (a *InvertAction) Do() {
	a.Error = a.Adapter.Invert()
}
