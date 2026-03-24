package styles

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// UserStylesDir returns the path to the user styles directory.
// Creates the directory if it does not exist.
func UserStylesDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	stylesDir := filepath.Join(dir, "commitlore", "styles")
	if err := os.MkdirAll(stylesDir, 0750); err != nil {
		return "", fmt.Errorf("creating styles dir: %w", err)
	}
	return stylesDir, nil
}

func configDir() (string, error) {
	if runtime.GOOS == "windows" {
		return windowsConfigDir()
	}
	return unixConfigDir()
}

func windowsConfigDir() (string, error) {
	dir := os.Getenv("APPDATA")
	if dir == "" {
		return "", fmt.Errorf("APPDATA not set")
	}
	return dir, nil
}

func unixConfigDir() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir: %w", err)
	}
	return filepath.Join(home, ".config"), nil
}
