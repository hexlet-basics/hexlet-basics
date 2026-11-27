package courseloader

import "strings"

// repositoryPrefix starts the name of every course's exercises repository.
const repositoryPrefix = "exercises-"

// RepositoryName is the exercises repository of the course with slug: each
// course's content lives in `exercises-<slug>` under the content organisation,
// and its build is extracted to a directory of the same name. The one place the
// convention is written; the fetcher, the parser, the GitHub webhook and the
// course's repository link all read it from here.
func RepositoryName(slug string) string {
	return repositoryPrefix + slug
}

// CourseOfRepository is the slug of the course whose exercises repository is
// name; false for a repository that is no course's.
func CourseOfRepository(name string) (string, bool) {
	slug, ok := strings.CutPrefix(name, repositoryPrefix)

	return slug, ok && slug != ""
}
