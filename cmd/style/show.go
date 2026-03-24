package style

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show the definition of a style",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(args[0])
		},
	}
}

func runShow(name string) error {
	s, err := styles.Load(name)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling style: %w", err)
	}
	fmt.Print(string(data))
	return nil
}
