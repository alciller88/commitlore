package style

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete an installed style (cannot delete built-ins)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(args[0])
		},
	}
}

func runDelete(name string) error {
	if err := styles.Delete(name); err != nil {
		return err
	}
	fmt.Printf("Style %q deleted.\n", name)
	return nil
}
