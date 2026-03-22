package main

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/styles"
)

// StyleApp exposes style management to the frontend.
type StyleApp struct{}

// NewStyleApp creates a new StyleApp instance.
func NewStyleApp() *StyleApp {
	return &StyleApp{}
}

// ListStyles returns all available styles (built-in + user) as JSON.
func (s *StyleApp) ListStyles() (string, error) {
	all, err := styles.ListAll()
	if err != nil {
		return "", cleanError(err)
	}

	type styleInfo struct {
		Name    string `json:"name"`
		BuiltIn bool   `json:"builtIn"`
		Desc    string `json:"description"`
		Author  string `json:"author"`
		Version string `json:"version"`
	}

	var result []styleInfo
	for _, name := range all {
		st, err := styles.Load(name)
		if err != nil {
			continue
		}
		result = append(result, styleInfo{
			Name:    st.Name,
			BuiltIn: styles.IsBuiltin(name),
			Desc:    st.Description,
			Author:  st.Author,
			Version: st.Version,
		})
	}

	return toJSON(result)
}

// StyleTheme holds the theme fields for frontend CSS variable injection.
type StyleTheme struct {
	Primary     string         `json:"primary"`
	Secondary   string         `json:"secondary"`
	Background  string         `json:"background"`
	Surface     string         `json:"surface"`
	Text        string         `json:"text"`
	Accent      string         `json:"accent"`
	Border      string         `json:"border"`
	FontFamily  string         `json:"fontFamily"`
	FontSize    string         `json:"fontSize"`
	Mode        string         `json:"mode"`
	Logo        string         `json:"logo"`
	WinDefault  string         `json:"winDefault"`
	WinClose    string         `json:"winClose"`
	WinMinimize string         `json:"winMinimize"`
	WinMaximize string         `json:"winMaximize"`
	UILabels    UILabelsDetail `json:"uiLabels"`
}

// UILabelsDetail holds navigation label overrides for the frontend.
type UILabelsDetail struct {
	Dashboard    string `json:"dashboard"`
	Generate     string `json:"generate"`
	Story        string `json:"story"`
	History      string `json:"history"`
	Contributors string `json:"contributors"`
	Styles       string `json:"styles"`
	Settings     string `json:"settings"`
}

// GetStyleTheme returns only the theme fields for a style, with fallback defaults.
func (s *StyleApp) GetStyleTheme(name string) (StyleTheme, error) {
	st, err := styles.Load(name)
	if err != nil {
		return StyleTheme{}, cleanError(err)
	}
	return buildStyleTheme(st), nil
}

func buildStyleTheme(st styles.Style) StyleTheme {
	return StyleTheme{
		Primary:     withDefault(st.Theme.Colors.Primary, "#58a6ff"),
		Secondary:   withDefault(st.Theme.Colors.Secondary, "#8b949e"),
		Background:  withDefault(st.Theme.Colors.Background, "#0d1117"),
		Surface:     withDefault(st.Theme.Colors.Surface, "#161b22"),
		Text:        withDefault(st.Theme.Colors.Text, "#e6edf3"),
		Accent:      withDefault(st.Theme.Colors.Accent, "#58a6ff"),
		Border:      withDefault(st.Theme.Colors.Border, "#30363d"),
		FontFamily:  withDefault(st.Theme.Typography.FontFamily, "system-ui, sans-serif"),
		FontSize:    withDefault(st.Theme.Typography.FontSizeBase, "14px"),
		Mode:        withDefault(st.Theme.Mode, "dark"),
		Logo:        st.Theme.Logo,
		WinDefault:  withDefault(st.Theme.WindowControls.Default, "#666666"),
		WinClose:    withDefault(st.Theme.WindowControls.Close, "#FF5F57"),
		WinMinimize: withDefault(st.Theme.WindowControls.Minimize, "#FEBC2E"),
		WinMaximize: withDefault(st.Theme.WindowControls.Maximize, "#28C840"),
		UILabels:    buildUILabels(st.UILabels),
	}
}

func buildUILabels(l styles.UILabels) UILabelsDetail {
	return UILabelsDetail{
		Dashboard:    withDefault(l.Dashboard, "Dashboard"),
		Generate:     withDefault(l.Generate, "Generate"),
		Story:        withDefault(l.Story, "Story"),
		History:      withDefault(l.History, "History"),
		Contributors: withDefault(l.Contributors, "Contributors"),
		Styles:       withDefault(l.Styles, "Styles"),
		Settings:     withDefault(l.Settings, "Settings"),
	}
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// ShowStyle returns the full style definition as JSON.
func (s *StyleApp) ShowStyle(name string) (string, error) {
	st, err := styles.Load(name)
	if err != nil {
		return "", cleanError(err)
	}
	return toJSON(st)
}

// ImportStyle imports a style from a local file path.
func (s *StyleApp) ImportStyle(path string) (string, error) {
	st, err := styles.ImportFromPath(path)
	if err != nil {
		return "", cleanError(err)
	}
	return st.Name, nil
}

// ExportStyle exports a style to the given output path.
func (s *StyleApp) ExportStyle(name, output string) error {
	return cleanError(styles.Export(name, output))
}

// DeleteStyle removes a user-installed style.
func (s *StyleApp) DeleteStyle(name string) error {
	return cleanError(styles.Delete(name))
}

// CreateStyle creates a new user style with the given parameters.
func (s *StyleApp) CreateStyle(name, description, author string) error {
	st := styles.Style{
		Name:        name,
		Description: description,
		Author:      author,
		Version:     "1.0.0",
	}
	return cleanError(styles.Save(st))
}

// StyleDetail holds all .shipstyle fields for the editor UI.
type StyleDetail struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	LLMPrompt   string            `json:"llmPrompt"`
	Tags        []string          `json:"tags"`
	PreviewURL  string            `json:"previewUrl"`
	Homepage    string            `json:"homepage"`
	Templates   TemplatesDetail   `json:"templates"`
	Vocabulary  map[string]string `json:"vocabulary"`
	Theme       ThemeDetail       `json:"theme"`
	Terminal    TerminalDetail    `json:"terminal"`
	UILabels    UILabelsDetail    `json:"uiLabels"`
}

// TemplatesDetail holds all template strings.
type TemplatesDetail struct {
	Header           string `json:"header"`
	Feature          string `json:"feature"`
	Fix              string `json:"fix"`
	Breaking         string `json:"breaking"`
	Footer           string `json:"footer"`
	StoryIntro       string `json:"storyIntro"`
	StoryMilestone   string `json:"storyMilestone"`
	StoryPeak        string `json:"storyPeak"`
	StoryContributor string `json:"storyContributor"`
	StoryFooter      string `json:"storyFooter"`
}

// ThemeDetail holds all theme fields.
type ThemeDetail struct {
	Mode           string           `json:"mode"`
	Colors         ColorsDetail     `json:"colors"`
	Typography     TypoDetail       `json:"typography"`
	HeaderImage    string           `json:"headerImage"`
	Logo           string           `json:"logo"`
	CardStyle      string           `json:"cardStyle"`
	Animations     bool             `json:"animations"`
	CustomCSS      string           `json:"customCss"`
	WindowControls WinControlDetail `json:"windowControls"`
}

// WinControlDetail holds window control colors.
type WinControlDetail struct {
	Default  string `json:"default"`
	Close    string `json:"close"`
	Minimize string `json:"minimize"`
	Maximize string `json:"maximize"`
}

// ColorsDetail holds the color palette.
type ColorsDetail struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	Accent     string `json:"accent"`
	Border     string `json:"border"`
}

// TypoDetail holds typography fields.
type TypoDetail struct {
	FontFamily   string `json:"fontFamily"`
	FontSizeBase string `json:"fontSizeBase"`
	FontSizeH    string `json:"fontSizeHeader"`
	FontSizeCode string `json:"fontSizeCode"`
}

// TerminalDetail holds terminal visual identity.
type TerminalDetail struct {
	Colors     TermColorsDetail `json:"colors"`
	Decorators DecorDetail      `json:"decorators"`
	Density    string           `json:"density"`
}

// TermColorsDetail holds terminal ANSI color names.
type TermColorsDetail struct {
	Header   string `json:"header"`
	Feature  string `json:"feature"`
	Fix      string `json:"fix"`
	Breaking string `json:"breaking"`
	Footer   string `json:"footer"`
}

// DecorDetail holds terminal decorator strings.
type DecorDetail struct {
	Separator string `json:"separator"`
	Bullet    string `json:"bullet"`
	Indent    string `json:"indent"`
}

// GetStyleDetail returns all fields of a .shipstyle for the editor UI.
func (s *StyleApp) GetStyleDetail(name string) (StyleDetail, error) {
	st, err := styles.Load(name)
	if err != nil {
		return StyleDetail{}, cleanError(err)
	}
	return styleToDetail(st), nil
}

// SaveStyleDetail saves all fields of a user style. Built-in styles are rejected.
func (s *StyleApp) SaveStyleDetail(detail StyleDetail) error {
	if styles.IsBuiltin(detail.Name) {
		return fmt.Errorf("cannot modify built-in style %q", detail.Name)
	}
	st := detailToStyle(detail)
	return cleanError(styles.Save(st))
}

// IsStyleBuiltIn returns true if the given name is a built-in style.
func (s *StyleApp) IsStyleBuiltIn(name string) bool {
	return styles.IsBuiltin(name)
}

func styleToDetail(st styles.Style) StyleDetail {
	return StyleDetail{
		Name: st.Name, Version: st.Version, Description: st.Description,
		Author: st.Author, LLMPrompt: st.LLMPrompt,
		Tags: st.Tags, PreviewURL: st.PreviewURL, Homepage: st.Homepage,
		Vocabulary: st.Vocabulary,
		Templates:  templatesToDetail(st.Templates),
		Theme:      themeToDetail(st.Theme),
		Terminal:   terminalToDetail(st.Terminal),
		UILabels:   buildUILabels(st.UILabels),
	}
}

func templatesToDetail(t styles.Templates) TemplatesDetail {
	return TemplatesDetail{
		Header: t.Header, Feature: t.Feature, Fix: t.Fix,
		Breaking: t.Breaking, Footer: t.Footer,
		StoryIntro: t.StoryIntro, StoryMilestone: t.StoryMilestone,
		StoryPeak: t.StoryPeak, StoryContributor: t.StoryContributor,
		StoryFooter: t.StoryFooter,
	}
}

func themeToDetail(t styles.Theme) ThemeDetail {
	return ThemeDetail{
		Mode: t.Mode, HeaderImage: t.HeaderImage, Logo: t.Logo,
		CardStyle: t.CardStyle, Animations: t.Animations, CustomCSS: t.CustomCSS,
		Colors: ColorsDetail{
			Primary: t.Colors.Primary, Secondary: t.Colors.Secondary,
			Background: t.Colors.Background, Surface: t.Colors.Surface,
			Text: t.Colors.Text, Accent: t.Colors.Accent, Border: t.Colors.Border,
		},
		Typography: TypoDetail{
			FontFamily: t.Typography.FontFamily, FontSizeBase: t.Typography.FontSizeBase,
			FontSizeH: t.Typography.FontSizeH, FontSizeCode: t.Typography.FontSizeCode,
		},
		WindowControls: WinControlDetail{
			Default: t.WindowControls.Default, Close: t.WindowControls.Close,
			Minimize: t.WindowControls.Minimize, Maximize: t.WindowControls.Maximize,
		},
	}
}

func terminalToDetail(t styles.Terminal) TerminalDetail {
	return TerminalDetail{
		Density: t.Density,
		Colors: TermColorsDetail{
			Header: t.Colors.Header, Feature: t.Colors.Feature,
			Fix: t.Colors.Fix, Breaking: t.Colors.Breaking, Footer: t.Colors.Footer,
		},
		Decorators: DecorDetail{
			Separator: t.Decorators.Separator, Bullet: t.Decorators.Bullet,
			Indent: t.Decorators.Indent,
		},
	}
}

func detailToStyle(d StyleDetail) styles.Style {
	return styles.Style{
		Name: d.Name, Version: d.Version, Description: d.Description,
		Author: d.Author, LLMPrompt: d.LLMPrompt,
		Tags: d.Tags, PreviewURL: d.PreviewURL, Homepage: d.Homepage,
		Vocabulary: d.Vocabulary,
		UILabels: styles.UILabels{
			Dashboard: d.UILabels.Dashboard, Generate: d.UILabels.Generate,
			Story: d.UILabels.Story, History: d.UILabels.History,
			Contributors: d.UILabels.Contributors, Styles: d.UILabels.Styles,
			Settings: d.UILabels.Settings,
		},
		Templates: styles.Templates{
			Header: d.Templates.Header, Feature: d.Templates.Feature,
			Fix: d.Templates.Fix, Breaking: d.Templates.Breaking,
			Footer: d.Templates.Footer, StoryIntro: d.Templates.StoryIntro,
			StoryMilestone: d.Templates.StoryMilestone, StoryPeak: d.Templates.StoryPeak,
			StoryContributor: d.Templates.StoryContributor, StoryFooter: d.Templates.StoryFooter,
		},
		Theme: styles.Theme{
			Mode: d.Theme.Mode, HeaderImage: d.Theme.HeaderImage,
			Logo: d.Theme.Logo, CardStyle: d.Theme.CardStyle,
			Animations: d.Theme.Animations, CustomCSS: d.Theme.CustomCSS,
			Colors: styles.Colors{
				Primary: d.Theme.Colors.Primary, Secondary: d.Theme.Colors.Secondary,
				Background: d.Theme.Colors.Background, Surface: d.Theme.Colors.Surface,
				Text: d.Theme.Colors.Text, Accent: d.Theme.Colors.Accent,
				Border: d.Theme.Colors.Border,
			},
			Typography: styles.Typography{
				FontFamily: d.Theme.Typography.FontFamily, FontSizeBase: d.Theme.Typography.FontSizeBase,
				FontSizeH: d.Theme.Typography.FontSizeH, FontSizeCode: d.Theme.Typography.FontSizeCode,
			},
			WindowControls: styles.WindowControls{
				Default: d.Theme.WindowControls.Default, Close: d.Theme.WindowControls.Close,
				Minimize: d.Theme.WindowControls.Minimize, Maximize: d.Theme.WindowControls.Maximize,
			},
		},
		Terminal: styles.Terminal{
			Density: d.Terminal.Density,
			Colors: styles.TerminalColors{
				Header: d.Terminal.Colors.Header, Feature: d.Terminal.Colors.Feature,
				Fix: d.Terminal.Colors.Fix, Breaking: d.Terminal.Colors.Breaking,
				Footer: d.Terminal.Colors.Footer,
			},
			Decorators: styles.TerminalDecorators{
				Separator: d.Terminal.Decorators.Separator, Bullet: d.Terminal.Decorators.Bullet,
				Indent: d.Terminal.Decorators.Indent,
			},
		},
	}
}
