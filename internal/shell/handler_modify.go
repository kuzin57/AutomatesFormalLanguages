package shell

import (
	"fmt"

	modifyactions "workspace/internal/actions/modify_actions"

	"github.com/spf13/cobra"
)

func registerModifySubcommands(shell *Shell) {
	makeModifyEpsCommand(shell)
	makeModifyDetCommand(shell)
	makeModifyFullCommand(shell)
	makeModifyInvertCommand(shell)
	makeMinimizeCommand(shell)
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
		Short: "Determine automate",
		RunE:  handler.RunE,
	}
	modifyCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

func makeModifyFullCommand(shell *Shell) {
	handler := &modifyMakeFullHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "full",
		Short: "make automate full",
		RunE:  handler.RunE,
	}
	modifyCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

func makeModifyInvertCommand(shell *Shell) {
	handler := &modifyInvertHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "invert",
		Short: "invert automate",
		RunE:  handler.RunE,
	}
	modifyCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

func makeMinimizeCommand(shell *Shell) {
	handler := &modifyMinHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "min",
		Short: "minimize automate",
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
		}
	}
	fmt.Println("success!")
	return nil
}

type modifyDetHandler struct {
	shell *Shell
	name  string
}

func (h *modifyDetHandler) RunE(cmd *cobra.Command, args []string) error {
	params := modifyactions.DetermineParams{}
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.NFA = adapter
			action, err := modifyactions.NewDetermineAction(&params)
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

type modifyMakeFullHandler struct {
	shell *Shell
	name  string
}

func (h *modifyMakeFullHandler) RunE(cmd *cobra.Command, args []string) error {
	params := modifyactions.MakeFullParams{}
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.Adapter = adapter
			action, err := modifyactions.NewMakeFullAction(&params)
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

type modifyInvertHandler struct {
	shell *Shell
	name  string
}

func (h *modifyInvertHandler) RunE(cmd *cobra.Command, args []string) error {
	params := modifyactions.InvertParams{}
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.Adapter = adapter
			action, err := modifyactions.NewInvertAction(&params)
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

type modifyMinHandler struct {
	shell *Shell
	name  string
}

func (h *modifyMinHandler) RunE(cmd *cobra.Command, args []string) error {
	params := modifyactions.MinimizeParams{}
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.Adapter = adapter
			action, err := modifyactions.NewMinimizeAction(&params)
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
