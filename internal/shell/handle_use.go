package shell

import (
	"fmt"

	showactions "workspace/internal/actions/show_actions"

	"github.com/spf13/cobra"
)

func registerUseSubcommands(shell *Shell) {
	makeUseReadCommand(shell)
	makeUseShowCommand(shell)
}

func makeUseReadCommand(shell *Shell) {
	handler := &useReadHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "read",
		Short: "read word",
		RunE:  handler.RunE,
	}

	useCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
	cmd.Flags().StringVarP(&handler.word, "word", "w", "", "word for reading")
}

func makeUseShowCommand(shell *Shell) {
	handler := &useShowHandler{shell: shell}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show states",
		RunE:  handler.RunE,
	}

	useCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&handler.name, "name", "n", "", "name of automate")
}

type useReadHandler struct {
	shell *Shell
	name  string
	word  string
}

func (h *useReadHandler) RunE(cmd *cobra.Command, args []string) error {
	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			err := adapter.Read(h.word)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("word exists!")
	return nil
}

type useShowHandler struct {
	shell *Shell
	name  string
}

func (h *useShowHandler) RunE(cmd *cobra.Command, args []string) (err error) {
	params := &showactions.ShowStatesParams{}

	for _, adapter := range h.shell.Automates {
		if adapter.GetName() == h.name {
			params.Adapter = adapter
			break
		}
	}

	action, err := showactions.NewShowStatesAction(params)
	if err != nil {
		return err
	}

	action.Do()
	if err = action.Error; err != nil {
		return err
	}

	return nil
}
