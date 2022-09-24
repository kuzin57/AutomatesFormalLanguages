package shell

import "github.com/spf13/cobra"

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create automate",
	}

	useCmd = &cobra.Command{
		Use:   "use",
		Short: "use automate",
	}
)
