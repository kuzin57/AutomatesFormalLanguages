package modifyactions

import (
	"workspace/adapters"
	automatesadapters "workspace/adapters/automates_adapters"
	"workspace/internal/actions"
)

type DeterminizeParams struct {
	NFA  adapters.AutomateAdapter
	Name string
}

type DeterminizeAction struct {
	actions.BaseAction

	params *DeterminizeParams
	result *DeterminizeResult
}

type DeterminizeResult struct {
	FA adapters.AutomateAdapter
}

func NewDeterminizeAction(params *DeterminizeParams) (*DeterminizeAction, error) {
	return &DeterminizeAction{BaseAction: actions.NewBaseAction(nil), params: params, result: &DeterminizeResult{}}, nil
}

func (a *DeterminizeAction) Do() {
	a.result.FA, a.Error = automatesadapters.NewFAadapter(a.params.NFA, a.params.Name)
}

func (a *DeterminizeAction) Result() *DeterminizeResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
