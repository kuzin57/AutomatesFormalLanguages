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
}

type createAutomateHandler struct {
	shell       *Shell
	regularExpr string
}

func (h *createAutomateHandler) RunE(cmd *cobra.Command, args []string) error {
	params := &makeactions.MakeNFAParams{Expr: h.regularExpr}
	action, err := makeactions.NewMakeNFAAction(params)
	if err != nil {
		return err
	}

	action.Do()

	// exprs := strings.Split(h.regularExpr, "+")

	return nil
}
