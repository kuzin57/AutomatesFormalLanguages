package actions

type BaseAction struct {
	Error error
}

func NewBaseAction() BaseAction {
	return BaseAction{}
}

func (a BaseAction) CheckErr() bool {
	return a.Error != nil
}
