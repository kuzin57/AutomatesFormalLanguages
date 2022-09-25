package makeactions

import (
	"workspace/adapters"
	automatesadapters "workspace/adapters/automates_adapters"
	"workspace/internal/actions"
	"workspace/internal/config"
)

type MakeNFAParams struct {
	Expr string
	Name string
}

type MakeNFAAction struct {
	actions.BaseAction

	params *MakeNFAParams
	result *MakeNFAResult
}

type MakeNFAResult struct {
	Adapter adapters.AutomateAdapter
}

func NewMakeNFAAction(params *MakeNFAParams, adapter adapters.AutomateAdapter) (*MakeNFAAction, error) {
	return &MakeNFAAction{BaseAction: actions.NewBaseAction(nil), params: params, result: &MakeNFAResult{}}, nil
}

func (a *MakeNFAAction) Do() {
	automateAdapter := automatesadapters.NewAutomateAdapter(config.MakeAdaptersConfig(false))
	automateAdapter.SetName(a.params.Name)

	newAutomateAdapter := automatesadapters.NewAutomateAdapter(config.MakeAdaptersConfig(false))

	newAutomateAdapter.Create(a.params.Name, a.params.Expr)
	if a.Error = automateAdapter.Join(newAutomateAdapter); a.CheckErr() {
		return
	}

	a.result = &MakeNFAResult{Adapter: automateAdapter}
}

func (a *MakeNFAAction) Result() *MakeNFAResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
