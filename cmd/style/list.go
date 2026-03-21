package style

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available styles (built-in + installed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}
}

func runList() error {
	printBuiltinStyles()
	return printUserStyles()
}

func printBuiltinStyles() {
	fmt.Println("Built-in styles:")
	for _, name := range styles.ListBuiltin() {
		fmt.Printf("  %s (built-in)\n", name)
	}
}

func printUserStyles() error {
	user, err := styles.ListUser()
	if err != nil {
		return err
	}
	if len(user) == 0 {
		fmt.Println("\nNo user-installed styles.")
		return nil
	}
	fmt.Println("\nInstalled styles:")
	for _, name := range user {
		fmt.Printf("  %s\n", name)
	}
	return nil
}
