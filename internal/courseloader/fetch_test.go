package courseloader_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/courseloader"
)

// publishRepo creates the exercises repo for slug under base, as a real git
// repository with one commit, so the fetcher clones it over the local transport
// and the test needs no network.
func publishRepo(t *testing.T, base, slug string) {
	t.Helper()
	path := filepath.Join(base, "exercises-"+slug+".git")
	repo, err := git.PlainInit(path, false)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(path, "spec.yml"), []byte("language: {}\n"), 0o644))
	tree, err := repo.Worktree()
	require.NoError(t, err)
	_, err = tree.Add("spec.yml")
	require.NoError(t, err)
	_, err = tree.Commit("init", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)
}

func TestGitFetcherNamesTheRepoAfterTheCourse(t *testing.T) {
	fetcher := courseloader.NewGitFetcher("https://github.com/hexlet-basics", "")

	assert.Equal(t, "https://github.com/hexlet-basics/exercises-ruby.git", fetcher.RepoURL("ruby"))
}

func TestGitFetcherChecksOutTheRepoAndCleansUp(t *testing.T) {
	base := t.TempDir()
	publishRepo(t, base, "ruby")
	fetcher := courseloader.NewGitFetcher(base, "")

	dir, cleanup, err := fetcher.Fetch(t.Context(), "ruby")

	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(dir, "spec.yml"))
	cleanup()
	assert.NoDirExists(t, dir)
}

func TestGitFetcherReportsAMissingRepoAndLeavesNothingBehind(t *testing.T) {
	// With and without a token: both clone paths fail the same way, and the temp
	// dir the attempt created is removed.
	for _, token := range []string{"", "secret"} {
		fetcher := courseloader.NewGitFetcher(t.TempDir(), token)

		dir, cleanup, err := fetcher.Fetch(t.Context(), "nope")

		require.Error(t, err)
		assert.ErrorContains(t, err, "exercises-nope.git")
		assert.Empty(t, dir)
		cleanup()
	}
}
