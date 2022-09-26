package shell

import (
	"fmt"

	modifyactions "workspace/internal/actions/modify_actions"

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
	params := modifyactions.DeleteEpsParams{}

	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.Adapter = adapter

			action, err := modifyactions.NewDeleteEpsAction(&params, nil)
			if err != nil {
				return err
			}

			action.Do()
			if err = action.Error; err != nil {
				return err
			}

			// adapter = action.Result().Adapter
		}
	}
	fmt.Println("success!")
	return nil
}
