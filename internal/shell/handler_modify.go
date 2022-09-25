package shell

import (
	"fmt"

	"github.com/spf13/cobra"
)

func registerModifySubcommands(shell *Shell) {
	makeModifyEpsCommand(shell)
}

func makeModifyEpsCommand(shell *Shell) {
	handler := &modifyEpsHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "eps",
		Short: "delete eps",
		RunE:  handler.RunE,
	}

	modifyCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

type modifyEpsHandler struct {
	shell *Shell
	name  string
}

func (h *modifyEpsHandler) RunE(cmd *cobra.Command, args []string) error {
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			err := adapter.DeleteEps()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("success!")
	return nil
}
