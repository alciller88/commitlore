package style

import "github.com/spf13/cobra"

// Register adds all style subcommands to the parent command.
func Register(parent *cobra.Command) {
	parent.AddCommand(newListCmd())
	parent.AddCommand(newShowCmd())
	parent.AddCommand(newCreateCmd())
	parent.AddCommand(newImportCmd())
	parent.AddCommand(newExportCmd())
	parent.AddCommand(newDeleteCmd())
}
