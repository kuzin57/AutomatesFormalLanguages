package makeactions

import (
	"fmt"
	"strings"
	"workspace/adapters"
	"workspace/internal/actions"
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
}

func NewMakeNFAAction(params *MakeNFAParams) (*MakeNFAAction, error) {
	return &MakeNFAAction{BaseAction: actions.NewBaseAction(), params: params, result: &MakeNFAResult{}}, nil
}

func (a *MakeNFAAction) Do() {
	parts := strings.Split(a.params.Expr, "+")
	automateAdapter := adapters.NewNFAAdapter()

	for _, part := range parts {
		fmt.Println("last: ", part[len(part)-2:len(part)-1])
		part = part[1:]
		part = part[:len(part)-1]
		smallerParts := strings.Split(part, ",")
		automateAdapter.Create(a.params.Name, smallerParts)
	}
}
