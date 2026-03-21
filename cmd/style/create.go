package style

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var name, description, author string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new style via flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(name, description, author)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Style name (required)")
	cmd.Flags().StringVar(&description, "description", "", "Style description")
	cmd.Flags().StringVar(&author, "author", "", "Style author")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func runCreate(name, description, author string) error {
	if styles.IsBuiltin(name) {
		return fmt.Errorf("cannot create style with built-in name %q", name)
	}
	s := buildNewStyle(name, description, author)
	if err := styles.Save(s); err != nil {
		return err
	}
	fmt.Printf("Style %q created.\n", name)
	return nil
}

func buildNewStyle(name, description, author string) styles.Style {
	return styles.Style{
		Name:        name,
		Version:     "1.0.0",
		Description: description,
		Author:      author,
		Templates: styles.Templates{
			Header:  "# Changelog",
			Feature: "- {{.Message}} ({{.Hash | short}})",
			Fix:     "- Fixed: {{.Message}} ({{.Hash | short}})",
		},
	}
}
