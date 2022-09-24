package shell

import (
	makeactions "workspace/internal/actions/make_actions"

	"github.com/spf13/cobra"
)

func registerCreateSubcommands(shell *Shell) {
	makeCreateAutomateCommand(shell)
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

type createAutomateHandler struct {
	shell       *Shell
	regularExpr string
	name        string
}

func (h *createAutomateHandler) RunE(cmd *cobra.Command, args []string) error {
	params := &makeactions.MakeNFAParams{Expr: h.regularExpr, Name: h.name}
	action, err := makeactions.NewMakeNFAAction(params, nil)
	if err != nil {
		return err
	}

	action.Do()
	if action.CheckErr() {
		return action.Error
	}

	h.shell.Automates = append(h.shell.Automates, action.Result().Adapter)

	return nil
}
