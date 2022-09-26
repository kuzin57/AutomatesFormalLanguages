package shell

import (
	"fmt"

	modifyactions "workspace/internal/actions/modify_actions"

	"github.com/spf13/cobra"
)

func registerModifySubcommands(shell *Shell) {
	makeModifyEpsCommand(shell)
	makeModifyDetCommand(shell)
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

func makeModifyDetCommand(shell *Shell) {
	handler := &modifyDetHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "det",
		Short: "determinize automate",
		RunE:  handler.RunE,
	}

	modifyCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.nfaName, "name", "n", "", "name of non det automate")
	cmd.Flags().StringVarP(&handler.name, "detname", "d", "", "name of det automate")
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
		}
	}
	fmt.Println("success!")
	return nil
}

type modifyDetHandler struct {
	shell   *Shell
	name    string
	nfaName string
}

func (h *modifyDetHandler) RunE(cmd *cobra.Command, args []string) error {
	params := modifyactions.DeterminizeParams{}
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.nfaName {
			params.NFA = adapter

			action, err := modifyactions.NewDeterminizeAction(&params)
			if err != nil {
				return err
			}

			action.Do()
			if err = action.Error; err != nil {
				return err
			}

			h.shell.Automates = append(h.shell.Automates, action.Result().FA)
		}
	}

	return nil
}
