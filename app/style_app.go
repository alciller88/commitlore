package app

import (
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
