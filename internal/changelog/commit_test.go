package changelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseType_conventionalFeat(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("feat: add new command"))
}

func TestParseType_conventionalFix(t *testing.T) {
	assert.Equal(t, TypeFix, ParseType("fix: resolve nil pointer"))
}

func TestParseType_conventionalChore(t *testing.T) {
	assert.Equal(t, TypeChore, ParseType("chore: update dependencies"))
}

func TestParseType_conventionalDocs(t *testing.T) {
	assert.Equal(t, TypeDocs, ParseType("docs: update README"))
}

func TestParseType_conventionalTest(t *testing.T) {
	assert.Equal(t, TypeTest, ParseType("test: add unit tests for parser"))
}

func TestParseType_conventionalRefactor(t *testing.T) {
	assert.Equal(t, TypeRefactor, ParseType("refactor: extract helper function"))
}

func TestParseType_breakingChangePrefix(t *testing.T) {
	assert.Equal(t, TypeBreaking, ParseType("BREAKING CHANGE: remove deprecated flag"))
}

func TestParseType_breakingBangSuffix(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("feat!: remove old API"))
}

func TestParseType_scopedPrefix(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("feat(api): add new endpoint"))
}

func TestParseType_featureAlias(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("feature: add dark mode"))
}

func TestParseType_bugfixAlias(t *testing.T) {
	assert.Equal(t, TypeFix, ParseType("bugfix: correct date parsing"))
}

func TestParseType_caseInsensitive(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("FEAT: uppercase prefix"))
	assert.Equal(t, TypeFix, ParseType("Fix: capitalized prefix"))
}

func TestParseType_inferFeat(t *testing.T) {
	assert.Equal(t, TypeFeat, ParseType("Add dark mode support"))
}

func TestParseType_inferFix(t *testing.T) {
	assert.Equal(t, TypeFix, ParseType("Fix crash on startup"))
}

func TestParseType_inferRefactor(t *testing.T) {
	assert.Equal(t, TypeRefactor, ParseType("Refactor the main loop"))
}

func TestParseType_inferDocs(t *testing.T) {
	assert.Equal(t, TypeDocs, ParseType("Update the README with examples"))
}

func TestParseType_inferTest(t *testing.T) {
	assert.Equal(t, TypeTest, ParseType("Improve test coverage"))
}

func TestParseType_inferChore(t *testing.T) {
	assert.Equal(t, TypeChore, ParseType("Bump version to 1.2.0"))
}

func TestParseType_inferBreaking(t *testing.T) {
	assert.Equal(t, TypeBreaking, ParseType("Removed legacy endpoint"))
}

func TestParseType_unknownMessage(t *testing.T) {
	assert.Equal(t, TypeOther, ParseType("WIP"))
}

func TestParseType_emptyMessage(t *testing.T) {
	assert.Equal(t, TypeOther, ParseType(""))
}

func TestParseType_noColonNoKeyword(t *testing.T) {
	assert.Equal(t, TypeOther, ParseType("initial commit"))
}
