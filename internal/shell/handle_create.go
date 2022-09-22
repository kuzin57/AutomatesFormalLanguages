package shell

import "github.com/spf13/cobra"

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
}

type createAutomateHandler struct {
	shell *Shell
}

func (h *createAutomateHandler) RunE(cmd *cobra.Command, args []string) error {
	return nil
}
