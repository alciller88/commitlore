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
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	LLMPrompt   string            `yaml:"llm_prompt"`
	Templates   Templates         `yaml:"templates"`
	Vocabulary  map[string]string `yaml:"vocabulary"`
	Theme       Theme             `yaml:"theme"`
	Terminal    Terminal          `yaml:"terminal"`
	Tags        []string          `yaml:"tags"`
	PreviewURL  string            `yaml:"preview_url"`
	Homepage    string            `yaml:"homepage"`
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

// Theme defines the visual identity for HTML output.
type Theme struct {
	Mode        string     `yaml:"mode"`
	Colors      Colors     `yaml:"colors"`
	Typography  Typography `yaml:"typography"`
	HeaderImage string     `yaml:"header_image"`
	Logo        string     `yaml:"logo"`
	CardStyle   string     `yaml:"card_style"`
	Animations  bool       `yaml:"animations"`
	CustomCSS   string     `yaml:"custom_css"`
}

// Colors holds the color palette for a theme.
type Colors struct {
	Primary    string `yaml:"primary"`
	Secondary  string `yaml:"secondary"`
	Background string `yaml:"background"`
	Surface    string `yaml:"surface"`
	Text       string `yaml:"text"`
	Accent     string `yaml:"accent"`
	Border     string `yaml:"border"`
}

// Typography holds font settings for a theme.
type Typography struct {
	FontFamily   string `yaml:"font_family"`
	FontSizeBase string `yaml:"font_size_base"`
	FontSizeH    string `yaml:"font_size_header"`
	FontSizeCode string `yaml:"font_size_code"`
}

// Terminal defines the visual identity for terminal output.
type Terminal struct {
	Colors     TerminalColors     `yaml:"colors"`
	Decorators TerminalDecorators `yaml:"decorators"`
	Density    string             `yaml:"density"`
}

// TerminalColors maps commit types to ANSI color names.
type TerminalColors struct {
	Header   string `yaml:"header"`
	Feature  string `yaml:"feature"`
	Fix      string `yaml:"fix"`
	Breaking string `yaml:"breaking"`
	Footer   string `yaml:"footer"`
}

// TerminalDecorators defines separator, bullet and indent strings.
type TerminalDecorators struct {
	Separator string `yaml:"separator"`
	Bullet    string `yaml:"bullet"`
	Indent    string `yaml:"indent"`
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
	if err := validateCardStyle(s); err != nil {
		return err
	}
	if err := validateDensity(s); err != nil {
		return err
	}
	return validateMode(s)
}

func validateCardStyle(s Style) error {
	v := s.Theme.CardStyle
	if v == "" {
		return nil
	}
	valid := map[string]bool{"minimal": true, "bordered": true, "glassmorphism": true}
	if !valid[v] {
		return fmt.Errorf("style %q: invalid card_style %q (use minimal, bordered, or glassmorphism)", s.Name, v)
	}
	return nil
}

func validateDensity(s Style) error {
	v := s.Terminal.Density
	if v == "" {
		return nil
	}
	valid := map[string]bool{"compact": true, "normal": true, "verbose": true}
	if !valid[v] {
		return fmt.Errorf("style %q: invalid density %q (use compact, normal, or verbose)", s.Name, v)
	}
	return nil
}

func validateMode(s Style) error {
	v := s.Theme.Mode
	if v == "" {
		return nil
	}
	valid := map[string]bool{"dark": true, "light": true}
	if !valid[v] {
		return fmt.Errorf("style %q: invalid mode %q (use dark or light)", s.Name, v)
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
