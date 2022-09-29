package shell

import (
	"fmt"
	makeactions "workspace/internal/actions/make_actions"

	"github.com/spf13/cobra"
)

func registerCreateSubcommands(shell *Shell) {
	makeCreateAutomateCommand(shell)
	makeCreateRegExCommand(shell)
}

func makeCreateAutomateCommand(shell *Shell) {
	handler := &createAutomateHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "automate",
		Short: "make an automate",
		RunE:  handler.RunE,
	}
	createCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.regularExpr, "regular", "r", "", "regular expression")
	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

func makeCreateRegExCommand(shell *Shell) {
	handler := &createRegExHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "regex",
		Short: "make regular expression",
		RunE:  handler.RunE,
	}
	createCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

type createAutomateHandler struct {
	shell       *Shell
	regularExpr string
	name        string
}

func (h *createAutomateHandler) RunE(cmd *cobra.Command, args []string) error {
	params := &makeactions.MakeAutomateParams{Expr: h.regularExpr, Name: h.name}
	action, err := makeactions.NewMakeAutomateAction(params, nil)
	if err != nil {
		return err
	}

	action.Do()
	if action.CheckErr() {
		return action.Error
	}

	h.shell.Automates = append(h.shell.Automates, action.Result().Adapter)

	fmt.Println("success!")

	return nil
}

type createRegExHandler struct {
	shell *Shell
	name  string
}

func (h *createRegExHandler) RunE(cmd *cobra.Command, args []string) (err error) {
	var expr string
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			expr, err = adapter.GetRegularExpr()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println(expr)
	return nil
}
