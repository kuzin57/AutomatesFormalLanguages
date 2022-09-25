package modifyactions

import (
	"workspace/adapters"
	"workspace/internal/actions"
)

type DeleteEpsParams struct {
	Adapter adapters.AutomateAdapter
}

type DeleteEpsAction struct {
	actions.BaseAction

	params *DeleteEpsParams
	result *DeleteEpsResult
}

type DeleteEpsResult struct {
	Adapter adapters.AutomateAdapter
}

func NewDeleteEpsAction(params *DeleteEpsParams, adapter adapters.AutomateAdapter) (*DeleteEpsAction, error) {
	return &DeleteEpsAction{BaseAction: actions.NewBaseAction(nil), params: params, result: &DeleteEpsResult{}}, nil
}

func (a *DeleteEpsAction) Do() {
	a.Error = a.params.Adapter.DeleteEps()

	a.result.Adapter = a.params.Adapter
}

func (a *DeleteEpsAction) Result() *DeleteEpsResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
