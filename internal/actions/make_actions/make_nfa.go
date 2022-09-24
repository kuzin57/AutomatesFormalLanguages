package makeactions

import (
	"strings"
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
	parts := strings.Split(a.params.Expr, "+")
	automateAdapter := automatesadapters.NewAutomateAdapter(config.MakeAdaptersConfig(false))

	for _, part := range parts {
		newAutomateAdapter := automatesadapters.NewAutomateAdapter(config.MakeAdaptersConfig(false))

		part = part[1:]
		part = part[:len(part)-1]
		smallerParts := strings.Split(part, ",")

		newAutomateAdapter.Create(a.params.Name, smallerParts)
		if a.Error = automateAdapter.Join(newAutomateAdapter); a.CheckErr() {
			return
		}

	}

	a.result = &MakeNFAResult{Adapter: automateAdapter}
}

func (a *MakeNFAAction) Result() *MakeNFAResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
