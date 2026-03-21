package cmd

import (
	"github.com/alciller88/commitlore/cmd/style"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newStyleCmd())
}

func newStyleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "style",
		Short: "Manage the modular style system",
	}
	style.Register(cmd)
	return cmd
}
