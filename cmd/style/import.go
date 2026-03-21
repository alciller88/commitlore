package style

import (
	"fmt"
	"strings"

	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <path-or-url>",
		Short: "Import a style from a local file or URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(args[0])
		},
	}
}

func runImport(source string) error {
	s, err := importStyle(source)
	if err != nil {
		return err
	}
	warnIfUntrustedLLMPrompt(s)
	fmt.Printf("Style %q imported successfully.\n", s.Name)
	return nil
}

func importStyle(source string) (styles.Style, error) {
	if isURL(source) {
		return styles.ImportFromURL(source)
	}
	return styles.ImportFromPath(source)
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func warnIfUntrustedLLMPrompt(s styles.Style) {
	if s.LLMPrompt == "" {
		return
	}
	fmt.Printf("Warning: this style contains an llm_prompt field. Imported styles are\n")
	fmt.Printf("untrusted — review it with \"commitlore style show %s\" before using --llm.\n", s.Name)
}
