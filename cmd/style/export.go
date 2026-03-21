package style

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export a style to a .shipstyle file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(args[0], output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: <name>.shipstyle)")

	return cmd
}

func runExport(name, output string) error {
	if output == "" {
		output = name + ".shipstyle"
	}
	if err := styles.Export(name, output); err != nil {
		return err
	}
	fmt.Printf("Style %q exported to %s\n", name, output)
	return nil
}
