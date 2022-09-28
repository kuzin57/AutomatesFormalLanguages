package modifyactions

import (
	"workspace/adapters"
	automateadapter "workspace/adapters/automates_adapters"
	"workspace/internal/actions"
)

type DetermineParams struct {
	NFA  adapters.AutomateAdapter
	Name string
}

type DetermineAction struct {
	actions.BaseAction

	params *DetermineParams
	result *DetermineResult
}

type DetermineResult struct {
	FA adapters.AutomateAdapter
}

func NewDetermineAction(params *DetermineParams) (*DetermineAction, error) {
	return &DetermineAction{BaseAction: actions.NewBaseAction(nil), params: params, result: &DetermineResult{}}, nil
}

func (a *DetermineAction) Do() {
	a.result.FA = automateadapter.NewAutomateAdapter(a.params.Name)
}

func (a *DetermineAction) Result() *DetermineResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
