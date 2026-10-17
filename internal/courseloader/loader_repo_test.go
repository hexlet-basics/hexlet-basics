package courseloader_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// miniRepo is a one-lesson course repo whose lesson locale files the test
// supplies, so a build can be pointed at exactly the content under test.
func miniRepo(t *testing.T, data, theory string) string {
	t.Helper()
	dir := t.TempDir()
	put := func(rel, content string) {
		path := filepath.Join(dir, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	put("spec.yml", `language:
  name: Mini
  progress: completed
  learn_as: first_language
  docker_image: hexletbasics/exercises-mini
  extension: js
  exercise_filename: exercise.js
  exercise_test_filename: test.js
`)
	put("modules/10-basics/description.en.yml", "name: Basics\ndescription: Start\n")
	lesson := "modules/10-basics/10-hello/"
	put(lesson+"exercise.js", "x\n")
	put(lesson+"test.js", "x\n")
	put(lesson+"en/data.yml", data)
	put(lesson+"en/README.md", theory)
	put(lesson+"en/EXERCISE.md", "do it")

	// The same picture the fixture course ships, for the theory to reference.
	picture, err := os.ReadFile(filepath.Join(fixtureRepo(t), "modules", "10-basics", "10-hello-world", "en", "assets", "dart.png"))
	require.NoError(t, err)
	put(lesson+"en/assets/a.png", string(picture))
	return dir
}

func TestLoaderStoresTipsAndDefinitionsAndLeavesAnAbsentDescriptionNull(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := t.Context()
	repo := miniRepo(t, "name: Hello\ntips:\n  - first tip\ndefinitions:\n  - name: term\n    description: its meaning\n", "theory")
	version := db.CourseVersion.GetX(ctx, fixtures.ID(fixtures.VersionLoadableCreated))

	require.NoError(t, newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID))

	info := db.CourseLessonTranslation.Query().
		Where(courselessontranslation.CourseVersionID(version.ID)).OnlyX(ctx)
	assert.Contains(t, derefStr(info.Tips), "first tip")
	assert.Contains(t, derefStr(info.Definitions), "its meaning")
	assert.Nil(t, info.Description)
}

func TestLoaderRewritesEveryLocalTheoryImageAndLeavesTheRestAlone(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := t.Context()
	theory := "one ![a](./assets/a.png) two ![remote](https://example.com/r.png) " +
		"three ![again](./assets/a.png) four ![ref][pic]\n\n[pic]: ./assets/a.png\n"
	repo := miniRepo(t, "name: Hello\ndescription: d\n", theory)
	version := db.CourseVersion.GetX(ctx, fixtures.ID(fixtures.VersionLoadableCreated))

	require.NoError(t, newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID))

	info := db.CourseLessonTranslation.Query().
		Where(courselessontranslation.CourseVersionID(version.ID)).OnlyX(ctx)
	got := derefStr(info.Theory)
	assert.Equal(t, 2, strings.Count(got, "http://localhost:3001/api/storage/"))
	assert.Contains(t, got, "![remote](https://example.com/r.png)")
	assert.Contains(t, got, "![ref][pic]")
	assert.NotContains(t, got, "(./assets/a.png)")
}

func TestLoaderFailsTheBuildWhenATheoryImageIsMissing(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := t.Context()
	repo := miniRepo(t, "name: Hello\n", "![gone](./assets/missing.png)")
	version := db.CourseVersion.GetX(ctx, fixtures.ID(fixtures.VersionLoadableCreated))

	err := newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID)

	require.Error(t, err)
	assert.ErrorContains(t, err, "open image")
	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "failed", derefStr(version.State))
	assert.Contains(t, derefStr(version.Result), "missing.png")
}

func TestLoaderFailsAVersionWhoseRepoCannotBeParsed(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := t.Context()
	repo := t.TempDir() // no spec.yml
	course := db.Course.GetX(ctx, fixtures.ID(fixtures.CourseLoadable))
	version := db.CourseVersion.GetX(ctx, fixtures.ID(fixtures.VersionLoadableCreated))

	err := newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID)

	require.Error(t, err)
	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "failed", derefStr(version.State))
	assert.Contains(t, derefStr(version.Result), "spec.yml")
	// The course is not promoted to a version that never built.
	assert.Nil(t, db.Course.GetX(ctx, course.ID).CurrentVersionID)
}

func TestLoaderFailsACourseWithoutASlug(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := t.Context()
	version := db.CourseVersion.GetX(ctx, fixtures.ID(fixtures.VersionLoadableWithoutSlugCreated))

	err := newLoaderWith(t, db, txStore, unreachableFetcher()).Run(ctx, version.ID)

	require.Error(t, err)
	assert.ErrorContains(t, err, "no slug")
	assert.Equal(t, "failed", derefStr(db.CourseVersion.GetX(ctx, version.ID).State))
}

func TestLoaderRejectsAVersionThatDoesNotExist(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)

	err := newLoaderWith(t, db, txStore, unreachableFetcher()).Run(t.Context(), 987654321)

	require.Error(t, err)
	assert.ErrorContains(t, err, "load course version")
}
