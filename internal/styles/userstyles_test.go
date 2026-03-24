package styles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleStyle() Style {
	return Style{
		Name:        "teststyle",
		Version:     "1.0.0",
		Description: "A test style",
		Author:      "tester",
		Templates: Templates{
			Header:  "# Test",
			Feature: "- {{.Message}}",
		},
	}
}

func setTestConfigDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func TestSaveAndLoadUser(t *testing.T) {
	setTestConfigDir(t)

	s := sampleStyle()
	err := Save(s)
	require.NoError(t, err)

	loaded, err := LoadUser("teststyle")
	require.NoError(t, err)
	assert.Equal(t, "teststyle", loaded.Name)
	assert.Equal(t, "A test style", loaded.Description)
}

func TestListUser_empty(t *testing.T) {
	setTestConfigDir(t)

	names, err := ListUser()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestListUser_withStyles(t *testing.T) {
	setTestConfigDir(t)

	require.NoError(t, Save(sampleStyle()))

	names, err := ListUser()
	require.NoError(t, err)
	assert.Contains(t, names, "teststyle")
}

func TestListAll_includesBothTypes(t *testing.T) {
	setTestConfigDir(t)

	require.NoError(t, Save(sampleStyle()))

	all, err := ListAll()
	require.NoError(t, err)
	assert.Contains(t, all, "formal")
	assert.Contains(t, all, "teststyle")
}

func TestDelete_userStyle(t *testing.T) {
	setTestConfigDir(t)

	require.NoError(t, Save(sampleStyle()))
	err := Delete("teststyle")
	require.NoError(t, err)

	_, err = LoadUser("teststyle")
	assert.Error(t, err)
}

func TestDelete_builtinFails(t *testing.T) {
	err := Delete("formal")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete built-in")
}

func TestDelete_nonexistent(t *testing.T) {
	setTestConfigDir(t)

	err := Delete("nonexistent")
	assert.Error(t, err)
}

func TestImportFromPath(t *testing.T) {
	setTestConfigDir(t)

	srcPath := filepath.Join(t.TempDir(), "imported.shipstyle")
	writeTestStyle(t, srcPath)

	s, err := ImportFromPath(srcPath)
	require.NoError(t, err)
	assert.Equal(t, "imported", s.Name)

	loaded, err := LoadUser("imported")
	require.NoError(t, err)
	assert.Equal(t, "imported", loaded.Name)
}

func TestImportFromPath_invalidFile(t *testing.T) {
	_, err := ImportFromPath("/nonexistent/file.shipstyle")
	assert.Error(t, err)
}

func TestExport_builtin(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "exported.shipstyle")
	err := Export("formal", outputPath)
	require.NoError(t, err)

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "formal")
}

func TestExport_nonexistent(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "out.shipstyle")
	err := Export("nonexistent", outputPath)
	assert.Error(t, err)
}

func TestLoadFromFile_invalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.shipstyle")
	require.NoError(t, os.WriteFile(path, []byte(":::invalid"), 0644))

	_, err := LoadFromFile(path)
	assert.Error(t, err)
}

func TestLoad_prefersBuiltinOverUser(t *testing.T) {
	setTestConfigDir(t)

	custom := sampleStyle()
	custom.Name = "formal"
	custom.Description = "custom formal"
	require.NoError(t, Save(custom))

	loaded, err := Load("formal")
	require.NoError(t, err)
	assert.NotEqual(t, "custom formal", loaded.Description)
}

func TestListUser_excludesLanguageVariants(t *testing.T) {
	setTestConfigDir(t)

	require.NoError(t, Save(sampleStyle()))

	// Write a language variant file directly
	dir, err := UserStylesDir()
	require.NoError(t, err)
	variantContent := `name: teststyle
language: "es"
version: "1.0.0"
templates:
  header: "# Prueba"
  feature: "- {{.Message}}"
`
	variantPath := filepath.Join(dir, "teststyle.es.shipstyle")
	require.NoError(t, os.WriteFile(variantPath, []byte(variantContent), 0644))

	names, err := ListUser()
	require.NoError(t, err)
	assert.Contains(t, names, "teststyle")
	assert.NotContains(t, names, "teststyle.es", "language variant should not appear as separate style")
}

func TestValidateName_valid(t *testing.T) {
	assert.NoError(t, ValidateName("my-style"))
	assert.NoError(t, ValidateName("style_v2"))
	assert.NoError(t, ValidateName("CamelCase123"))
}

func TestValidateName_rejectsPathTraversal(t *testing.T) {
	err := ValidateName(`..\..\Windows\malicious`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid style name")
}

func TestValidateName_rejectsDots(t *testing.T) {
	assert.Error(t, ValidateName("../etc/passwd"))
	assert.Error(t, ValidateName("foo.bar"))
	assert.Error(t, ValidateName("a/b"))
}

func TestValidateName_rejectsEmpty(t *testing.T) {
	assert.Error(t, ValidateName(""))
}

func TestSave_rejectsInvalidName(t *testing.T) {
	setTestConfigDir(t)
	s := sampleStyle()
	s.Name = "../malicious"
	err := Save(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid style name")
}

func TestDelete_rejectsInvalidName(t *testing.T) {
	err := Delete("../../etc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid style name")
}

func TestValidateOutputPath_allowsNormal(t *testing.T) {
	assert.NoError(t, ValidateOutputPath("/tmp/style.shipstyle"))
	assert.NoError(t, ValidateOutputPath("output/formal.shipstyle"))
}

func TestValidateOutputPath_rejectsGitDir(t *testing.T) {
	err := ValidateOutputPath("/repo/.git/hooks/style.shipstyle")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ".git directory")
}

func TestValidateOutputPath_rejectsGitSubdir(t *testing.T) {
	assert.Error(t, ValidateOutputPath("project/.git/objects/out"))
	assert.Error(t, ValidateOutputPath(`.git\config`))
}

func TestExport_rejectsGitPath(t *testing.T) {
	err := Export("formal", "/repo/.git/style.shipstyle")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ".git directory")
}

func writeTestStyle(t *testing.T, path string) {
	t.Helper()
	content := `name: imported
version: 1.0.0
description: "An imported style"
author: "importer"
templates:
  header: "# Imported"
  feature: "- {{.Message}}"
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}
