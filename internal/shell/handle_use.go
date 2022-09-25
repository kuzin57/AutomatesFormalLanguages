package shell

import (
	"fmt"

	"github.com/spf13/cobra"
)

func registerUseSubcommands(shell *Shell) {
	makeUseReadCommand(shell)
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
