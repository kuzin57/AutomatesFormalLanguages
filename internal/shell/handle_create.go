package shell

import (
	"fmt"
	automatesadapters "workspace/adapters/automatas_adapters"
	"workspace/internal/entities/parser"

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

func (h *createAutomateHandler) RunE(cmd *cobra.Command, args []string) (err error) {
	automateAdapter := automatesadapters.NewAutomateAdapter(h.name)
	automateAdapter.SetName(h.name)

	newAutomateAdapter := automatesadapters.NewAutomateAdapter(h.name)

	parser := parser.NewParser(h.regularExpr, nil)
	parser.Parse()

	newAutomateAdapter.Create(h.name, parser)
	if err = automateAdapter.Join(newAutomateAdapter); err != nil {
		return
	}

	h.shell.Automates = append(h.shell.Automates, automateAdapter)
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
			expr, err = adapter.GetRegularExpression()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println(expr)
	return nil
}
