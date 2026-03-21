package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRemoteRepo_OwnerRepo(t *testing.T) {
	assert.True(t, IsRemoteRepo("owner/repo"))
	assert.True(t, IsRemoteRepo("my-org/my-repo"))
	assert.True(t, IsRemoteRepo("user123/project.name"))
}

func TestIsRemoteRepo_GitHubURL(t *testing.T) {
	assert.True(t, IsRemoteRepo("https://github.com/owner/repo"))
	assert.True(t, IsRemoteRepo("https://github.com/owner/repo.git"))
	assert.True(t, IsRemoteRepo("http://github.com/owner/repo"))
	assert.True(t, IsRemoteRepo("https://github.com/owner/repo/"))
}

func TestIsRemoteRepo_LocalPath(t *testing.T) {
	assert.False(t, IsRemoteRepo("."))
	assert.False(t, IsRemoteRepo("/path/to/repo"))
	assert.False(t, IsRemoteRepo("C:\\Users\\repo"))
	assert.False(t, IsRemoteRepo("relative/path/deep"))
}

func TestIsRemoteRepo_InvalidURLs(t *testing.T) {
	assert.False(t, IsRemoteRepo("https://gitlab.com/owner/repo"))
	assert.False(t, IsRemoteRepo("https://github.com/owner"))
	assert.False(t, IsRemoteRepo("https://github.com/"))
}

func TestParseRepoRef_OwnerRepo(t *testing.T) {
	owner, repo, err := ParseRepoRef("owner/repo")
	assert.NoError(t, err)
	assert.Equal(t, "owner", owner)
	assert.Equal(t, "repo", repo)
}

func TestParseRepoRef_GitHubURL(t *testing.T) {
	owner, repo, err := ParseRepoRef("https://github.com/my-org/my-project")
	assert.NoError(t, err)
	assert.Equal(t, "my-org", owner)
	assert.Equal(t, "my-project", repo)
}

func TestParseRepoRef_GitHubURLWithGit(t *testing.T) {
	owner, repo, err := ParseRepoRef("https://github.com/owner/repo.git")
	assert.NoError(t, err)
	assert.Equal(t, "owner", owner)
	assert.Equal(t, "repo", repo)
}

func TestParseRepoRef_InvalidFormat(t *testing.T) {
	_, _, err := ParseRepoRef("/local/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GitHub reference")
}

func TestParseRepoRef_InvalidURL(t *testing.T) {
	_, _, err := ParseRepoRef("https://gitlab.com/owner/repo")
	assert.Error(t, err)
}

func TestNewClient_CreatesClient(t *testing.T) {
	client := NewClient("owner", "repo")
	assert.NotNil(t, client)
	assert.Equal(t, "owner", client.owner)
	assert.Equal(t, "repo", client.repo)
	assert.NotNil(t, client.gh)
}

func TestHasToken_WithoutEnv(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	assert.False(t, HasToken())
}

func TestHasToken_WithEnv(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "ghp_testtoken123")
	assert.True(t, HasToken())
}
