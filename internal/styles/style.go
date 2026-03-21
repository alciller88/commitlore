package styles

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed builtin/*.shipstyle
var builtinFS embed.FS

// Style represents a .shipstyle definition.
type Style struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	Author      string    `yaml:"author"`
	LLMPrompt   string    `yaml:"llm_prompt"`
	Templates   Templates `yaml:"templates"`
}

// Templates holds the template strings for each commit type and story sections.
type Templates struct {
	Header   string `yaml:"header"`
	Feature  string `yaml:"feature"`
	Fix      string `yaml:"fix"`
	Breaking string `yaml:"breaking"`
	Footer   string `yaml:"footer"`

	StoryIntro       string `yaml:"story_intro"`
	StoryMilestone   string `yaml:"story_milestone"`
	StoryPeak        string `yaml:"story_peak"`
	StoryContributor string `yaml:"story_contributor"`
	StoryFooter      string `yaml:"story_footer"`
}

var builtinNames = []string{"formal", "patchnotes", "ironic", "epic"}

// Load returns a style by name. Looks up built-in styles first.
func Load(name string) (Style, error) {
	if name == "" {
		name = "formal"
	}
	return loadBuiltin(name)
}

func loadBuiltin(name string) (Style, error) {
	path := "builtin/" + name + ".shipstyle"
	data, err := builtinFS.ReadFile(path)
	if err != nil {
		return Style{}, fmt.Errorf("style %q not found: %w", name, err)
	}
	return parseStyle(data)
}

func parseStyle(data []byte) (Style, error) {
	var s Style
	if err := yaml.Unmarshal(data, &s); err != nil {
		return Style{}, fmt.Errorf("parsing style: %w", err)
	}
	if err := validate(s); err != nil {
		return Style{}, err
	}
	return s, nil
}

func validate(s Style) error {
	if s.Name == "" {
		return fmt.Errorf("style missing required field: name")
	}
	if s.Templates.Header == "" && s.Templates.Feature == "" {
		return fmt.Errorf("style %q has no templates defined", s.Name)
	}
	return nil
}

// ListBuiltin returns the names of all built-in styles.
func ListBuiltin() []string {
	return append([]string{}, builtinNames...)
}

// IsBuiltin returns true if the given name is a built-in style.
func IsBuiltin(name string) bool {
	for _, n := range builtinNames {
		if n == name {
			return true
		}
	}
	return false
}
