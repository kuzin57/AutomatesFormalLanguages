package makeactions

import (
	"workspace/adapters"
	automatesadapters "workspace/adapters/automates_adapters"
	"workspace/internal/actions"
	"workspace/internal/entities/parser"
)

type MakeAutomateParams struct {
	Expr string
	Name string
}

type MakeAutomateAction struct {
	actions.BaseAction

	params *MakeAutomateParams
	result *MakeAutomateResult
}

type MakeAutomateResult struct {
	Adapter adapters.AutomateAdapter
}

func NewMakeAutomateAction(params *MakeAutomateParams, adapter adapters.AutomateAdapter) (*MakeAutomateAction, error) {
	return &MakeAutomateAction{BaseAction: actions.NewBaseAction(nil), params: params, result: &MakeAutomateResult{}}, nil
}

func (a *MakeAutomateAction) Do() {
	automateAdapter := automatesadapters.NewAutomateAdapter(a.params.Name)
	automateAdapter.SetName(a.params.Name)

	newAutomateAdapter := automatesadapters.NewAutomateAdapter(a.params.Name)

	parser := parser.NewParser(a.params.Expr, nil)
	parser.Parse()

	newAutomateAdapter.Create(a.params.Name, parser)
	if a.Error = automateAdapter.Join(newAutomateAdapter); a.CheckErr() {
		return
	}

	a.result = &MakeAutomateResult{Adapter: automateAdapter}
}

func (a *MakeAutomateAction) Result() *MakeAutomateResult {
	if a.CheckErr() {
		return nil
	}
	return a.result
}
