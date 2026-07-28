package courseloader

import (
	"context"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/samber/oops"
)

// Fetcher checks out a course's exercises repository and returns a local
// directory holding the repo tree, plus a cleanup to remove it. It is an
// interface so the loader can run against a fixture directory in tests without
// touching the network (see fakeFetcher in the loader tests).
type Fetcher interface {
	Fetch(ctx context.Context, slug string) (dir string, cleanup func(), err error)
}

// GitFetcher fetches course content by shallow-cloning the course's repo. A
// depth-1, single-branch clone of the default branch is all the loader needs:
// the parser reads the working tree, and history/other branches would only add
// transfer cost. Replaces the legacy Docker-image pull — the git tree is
// byte-identical to the built image for the content files the parser reads.
type GitFetcher struct {
	// baseURL is the org base (e.g. https://github.com/hexlet-basics); the repo
	// URL is baseURL/exercises-<slug>.git.
	baseURL string
	// token, when non-empty, authenticates the clone as the HTTP basic-auth
	// password with the conventional x-access-token username.
	token string
}

// NewGitFetcher builds a GitFetcher from the repo base URL and an optional token.
func NewGitFetcher(baseURL, token string) *GitFetcher {
	return &GitFetcher{baseURL: baseURL, token: token}
}

// RepoURL is the clone URL for a course slug, exposed so callers (webhook repo
// matching, logging) can reason about the same convention the fetcher uses.
func (f *GitFetcher) RepoURL(slug string) string {
	return f.baseURL + "/exercises-" + slug + ".git"
}

// Fetch shallow-clones the course repo into a fresh temp dir. The cleanup removes
// that dir; callers must always defer it, including on a later parse error.
func (f *GitFetcher) Fetch(ctx context.Context, slug string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "exercises-"+slug+"-")
	if err != nil {
		return "", func() {}, oops.Wrapf(err, "create temp dir for %s", slug)
	}
	cleanup := func() { _ = os.RemoveAll(dir) }

	opts := &git.CloneOptions{
		URL:          f.RepoURL(slug),
		Depth:        1,
		SingleBranch: true,
	}
	if f.token != "" {
		// GitHub accepts a personal-access / installation token as the basic-auth
		// password with any non-empty username; x-access-token is the convention.
		opts.Auth = &http.BasicAuth{Username: "x-access-token", Password: f.token}
	}

	if _, err := git.PlainCloneContext(ctx, dir, false, opts); err != nil {
		cleanup()
		return "", func() {}, oops.Wrapf(err, "clone %s", f.RepoURL(slug))
	}

	return dir, cleanup, nil
}
