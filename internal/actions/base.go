package actions

import (
	"workspace/adapters"
)

type BaseAction struct {
	Adapter adapters.AutomateAdapter
	Error   error
}

func NewBaseAction(adapter adapters.AutomateAdapter) BaseAction {
	return BaseAction{Adapter: adapter}
}

func (a BaseAction) CheckErr() bool {
	return a.Error != nil
}
