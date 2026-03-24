package styles

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateName checks that a style name is safe for use as a filename.
func ValidateName(name string) error {
	if !validNamePattern.MatchString(name) {
		return fmt.Errorf(
			"invalid style name %q: must contain only letters, numbers, hyphens, and underscores",
			name)
	}
	return nil
}

// ValidateOutputPath checks that a path is not inside a .git directory.
func ValidateOutputPath(path string) error {
	cleaned := filepath.Clean(path)
	for _, segment := range splitPath(cleaned) {
		if segment == ".git" {
			return fmt.Errorf("output path cannot be inside a .git directory")
		}
	}
	return nil
}

func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\' || r == filepath.Separator
	})
}

// LoadUser loads a user-installed style by name from the styles directory.
func LoadUser(name string) (Style, error) {
	if err := ValidateName(name); err != nil {
		return Style{}, err
	}
	dir, err := UserStylesDir()
	if err != nil {
		return Style{}, err
	}
	path := filepath.Join(dir, name+".shipstyle")
	return LoadFromFile(path)
}

// LoadFromFile loads and validates a style from a filesystem path.
func LoadFromFile(path string) (Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Style{}, fmt.Errorf("reading style file %q: %w", path, err)
	}
	return parseStyle(data)
}

// Save writes a style to the user styles directory.
func Save(s Style) error {
	if err := ValidateName(s.Name); err != nil {
		return err
	}
	if err := validate(s); err != nil {
		return err
	}
	dir, err := UserStylesDir()
	if err != nil {
		return err
	}
	return writeStyleFile(s, filepath.Join(dir, s.Name+".shipstyle"))
}

func writeStyleFile(s Style, path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling style: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing style to %s: %w", path, err)
	}
	return nil
}

// ListUser returns the names of user-installed styles.
func ListUser() ([]string, error) {
	dir, err := UserStylesDir()
	if err != nil {
		return nil, err
	}
	return listStyleFiles(dir)
}

func listStyleFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading styles dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if name, ok := styleFileName(e); ok {
			names = append(names, name)
		}
	}
	return names, nil
}

// knownLanguages lists language codes used for variant files.
var knownLanguages = map[string]bool{"es": true}

func styleFileName(e os.DirEntry) (string, bool) {
	if e.IsDir() {
		return "", false
	}
	name := e.Name()
	if !strings.HasSuffix(name, ".shipstyle") {
		return "", false
	}
	base := strings.TrimSuffix(name, ".shipstyle")
	if isLanguageVariantName(base) {
		return "", false
	}
	return base, true
}

// isLanguageVariantName returns true if the name ends with a known
// language suffix (e.g. "cyberpunk.es").
func isLanguageVariantName(name string) bool {
	idx := strings.LastIndex(name, ".")
	if idx < 1 {
		return false
	}
	return knownLanguages[name[idx+1:]]
}

// ListAll returns both built-in and user-installed style names.
func ListAll() ([]string, error) {
	all := ListBuiltin()
	user, err := ListUser()
	if err != nil {
		return all, nil
	}
	return append(all, user...), nil
}

// Delete removes a user-installed style and its language variants.
// Built-in styles cannot be deleted.
func Delete(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if IsBuiltin(name) {
		return fmt.Errorf("cannot delete built-in style %q", name)
	}
	dir, err := UserStylesDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name+".shipstyle")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting style %q: %w", name, err)
	}
	deleteLanguageVariants(dir, name)
	return nil
}

func deleteLanguageVariants(dir, name string) {
	for lang := range knownLanguages {
		path := filepath.Join(dir, name+"."+lang+".shipstyle")
		_ = os.Remove(path)
	}
}

// ImportFromPath imports a .shipstyle file from a local path.
func ImportFromPath(path string) (Style, error) {
	s, err := LoadFromFile(path)
	if err != nil {
		return Style{}, err
	}
	if err := Save(s); err != nil {
		return Style{}, err
	}
	return s, nil
}

// ImportFromURL downloads and imports a .shipstyle file from a URL.
func ImportFromURL(url string) (Style, error) {
	data, err := downloadStyle(url)
	if err != nil {
		return Style{}, err
	}
	s, err := parseStyle(data)
	if err != nil {
		return Style{}, err
	}
	if err := Save(s); err != nil {
		return Style{}, err
	}
	return s, nil
}

func downloadStyle(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec,noctx // user-provided URL, no context needed
	if err != nil {
		return nil, fmt.Errorf("downloading style from %s: %w", url, err)
	}
	data, readErr := readResponse(resp)
	closeErr := resp.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, fmt.Errorf("closing response: %w", closeErr)
	}
	return data, nil
}

func readResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading style: HTTP %d", resp.StatusCode)
	}
	return readLimited(resp.Body)
}

const maxStyleSize = 1024 * 1024 // 1MB

func readLimited(r io.Reader) ([]byte, error) {
	limited := io.LimitReader(r, maxStyleSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading style data: %w", err)
	}
	if len(data) > maxStyleSize {
		return nil, fmt.Errorf("style file exceeds maximum size of 1MB")
	}
	return data, nil
}

// Export writes a style (built-in or user) to the specified output path.
func Export(name, outputPath string) error {
	if err := ValidateOutputPath(outputPath); err != nil {
		return err
	}
	s, err := Load(name)
	if err != nil {
		return err
	}
	return writeStyleFile(s, outputPath)
}
