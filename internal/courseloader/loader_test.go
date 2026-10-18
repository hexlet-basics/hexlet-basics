package courseloader_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob/memblob"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/coursemodule"
	"hexletbasics/ent/coursemoduleversion"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/events"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

// Faults a test injects; each is told apart from the others by identity.
var (
	errGitIsUnreachable = errors.New("git is unreachable")
)

// repoFetcher is the git boundary answering with a local checkout instead of
// cloning, so the loader runs fully offline against a fixture course.
func repoFetcher(dir string) *testsupport.FetcherMock {
	return &testsupport.FetcherMock{
		FetchFunc: func(context.Context, string) (string, func(), error) { return dir, func() {}, nil },
	}
}

// unreachableFetcher is a git boundary that fails. A build that reaches it ends
// in an error, so a test of a version that must not be built passes only if the
// loader never gets that far.
func unreachableFetcher() *testsupport.FetcherMock {
	return &testsupport.FetcherMock{
		FetchFunc: func(context.Context, string) (string, func(), error) {
			return "", nil, errGitIsUnreachable
		},
	}
}

func newLoaderWith(
	t *testing.T,
	db *ent.Client,
	txStore store.Transactor,
	fetcher courseloader.Fetcher,
) *courseloader.Loader {
	t.Helper()

	return newLoaderWithCompletion(t, db, txStore, fetcher,
		events.NewPublisher(txStore, slog.New(slog.DiscardHandler)))
}

// newLoaderWithCompletion builds a loader whose promotion re-evaluates
// completion through the real progress module, recording the facts it publishes.
func newLoaderWithCompletion(
	t *testing.T,
	db *ent.Client,
	txStore store.Transactor,
	fetcher courseloader.Fetcher,
	publisher events.TxPublisher,
) *courseloader.Loader {
	t.Helper()

	bucket := memblob.OpenBucket(nil)

	t.Cleanup(func() { _ = bucket.Close() })

	assets := assetstore.New(db, bucket, "http://localhost:3001")

	// Promotion re-evaluates completion and never runs a submission, so the
	// loader's progress module needs no exercise runner.
	tracker := progress.New(db, txStore, publisher, progress.UnavailableRunner{})

	return courseloader.NewLoader(db, txStore, assets, fetcher, tracker,
		courseloader.WithClock(func() time.Time { return testsupport.Today }))
}

// newLoader builds a loader that fetches the committed fixture course.
func newLoader(t *testing.T, db *ent.Client, txStore store.Transactor) *courseloader.Loader {
	t.Helper()

	return newLoaderWith(t, db, txStore, repoFetcher(fixtureRepo(t)))
}

func TestLoaderBuildsAndPromotesVersion(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()

	// Arrange: a course with a freshly-created version, the state the loader
	// requires. (Its own writes roll back with the test transaction.)
	course := db.Course.GetX(ctx, fixture.ID(fixture.CourseLoadable))
	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))

	loader := newLoader(t, db, txStore)
	require.NoError(t, loader.Run(ctx, version.ID))

	// Version built, spec applied, and promoted to the course's current version.
	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "built", derefStr(version.State))
	assert.Equal(t, "Success", derefStr(version.Result))
	assert.Equal(t, "JavaScript", derefStr(version.Name))
	assert.Equal(t, "js", derefStr(version.Extension))
	assert.Equal(t, 1, version.LessonsCount)

	course = db.Course.GetX(ctx, course.ID)
	require.NotNil(t, course.CurrentVersionID)
	assert.Equal(t, version.ID, *course.CurrentVersionID)
	assert.Equal(t, 1, course.LessonsCount)

	// One module "basics" with 3 locale infos on this version.
	mods := db.CourseModule.Query().Where(coursemodule.CourseID(course.ID)).AllX(ctx)
	require.Len(t, mods, 1)
	assert.Equal(t, "basics", derefStr(mods[0].Slug))

	mvs := db.CourseModuleVersion.Query().Where(coursemoduleversion.CourseVersionID(version.ID)).AllX(ctx)
	require.Len(t, mvs, 1)
	assert.Equal(t, 10, *mvs[0].Order)

	// One lesson version, first in the course, with prepared/original code.
	lvs := db.CourseLessonVersion.Query().Where(courselessonversion.CourseVersionID(version.ID)).AllX(ctx)
	require.Len(t, lvs, 1)
	assert.Equal(t, 1, *lvs[0].NaturalOrder)
	assert.Equal(t, "/exercises-loadable/modules/10-basics/10-hello-world", derefStr(lvs[0].PathToCode))
	assert.Contains(t, derefStr(lvs[0].OriginalCode), "Hello, World!")

	// Three lesson infos (en/es/ru); theory images rewritten to /storage URLs.
	infos := db.CourseLessonTranslation.Query().Where(courselessontranslation.CourseVersionID(version.ID)).AllX(ctx)
	require.Len(t, infos, 3)
	en := lessonInfoByLocale(t, infos, "en")
	assert.Equal(t, "Hello, World!", derefStr(en.Name))
	assert.Contains(t, derefStr(en.Theory), "http://localhost:3001/api/storage/")
	assert.NotContains(t, derefStr(en.Theory), "./assets/")
}

func TestLoaderSkipsNonCreatedVersion(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()

	course := db.Course.GetX(ctx, fixture.ID(fixture.CourseLoadable))
	// A version already 'built' must not be rebuilt or have its result overwritten.
	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableBuilt))

	loader := newLoaderWith(t, db, txStore, unreachableFetcher())
	require.NoError(t, loader.Run(ctx, version.ID))

	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "built", derefStr(version.State)) // unchanged
	assert.Equal(t, "Success", derefStr(version.Result))
	// No module rows were written for this course.
	n := db.CourseModule.Query().Where(coursemodule.CourseID(course.ID)).CountX(ctx)
	assert.Equal(t, 0, n)
}

// A version's job can reach the worker twice. The first run claims the version
// and builds it; the second finds it no longer `created`, so it does not touch
// git, rebuild, or overwrite the result.
func TestLoaderBuildsAVersionOnceWhenItsJobIsDeliveredTwice(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	id := fixture.ID(fixture.VersionLoadableCreated)

	require.NoError(t, newLoader(t, db, txStore).Run(ctx, id))
	built := db.CourseVersion.GetX(ctx, id)
	require.Equal(t, "built", derefStr(built.State))

	// A second delivery that reached git would fail: git is unreachable here.
	require.NoError(t, newLoaderWith(t, db, txStore, unreachableFetcher()).Run(ctx, id))

	again := db.CourseVersion.GetX(ctx, id)
	assert.Equal(t, "built", derefStr(again.State))
	assert.Equal(t, derefStr(built.Result), derefStr(again.Result))
	assert.True(t, built.UpdatedAt.Equal(again.UpdatedAt), "the second delivery wrote nothing")
}

func TestLoaderUpsertsStableLessonAcrossRebuilds(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()

	course := db.Course.GetX(ctx, fixture.ID(fixture.CourseLoadable))
	loader := newLoader(t, db, txStore)

	// First build.
	v1 := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))
	require.NoError(t, loader.Run(ctx, v1.ID))
	module1 := db.CourseModule.Query().Where(coursemodule.CourseID(course.ID)).OnlyX(ctx)
	lesson1 := db.CourseLesson.Query().Where(courselesson.CourseID(course.ID)).OnlyX(ctx)

	// Second build reuses the SAME lesson identity (learner progress FKs it).
	secondVersion := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreatedSecond))
	require.NoError(t, loader.Run(ctx, secondVersion.ID))

	modules := db.CourseModule.Query().Where(coursemodule.CourseID(course.ID)).AllX(ctx)
	require.Len(t, modules, 1)
	assert.Equal(t, module1.ID, modules[0].ID)

	lessons := db.CourseLesson.Query().Where(courselesson.CourseID(course.ID)).AllX(ctx)
	require.Len(t, lessons, 1)
	assert.Equal(t, lesson1.ID, lessons[0].ID)

	// Current version is the newest built one.
	course = db.Course.GetX(ctx, course.ID)
	assert.Equal(t, secondVersion.ID, *course.CurrentVersionID)
}

// TestLoaderFailedRebuildKeepsLiveVersion is the load-bearing invariant (legacy
// ADR-0001): a broken rebuild must NOT take the live course offline. v1 builds and
// goes live; v2 fails to parse (empty repo, no spec.yml); the course must still
// point at v1, and v2 must be recorded as failed.
func TestLoaderFailedRebuildKeepsLiveVersion(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()

	course := db.Course.GetX(ctx, fixture.ID(fixture.CourseLoadable))

	// v1 builds successfully and is promoted live.
	firstVersion := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))
	require.NoError(t, newLoader(t, db, txStore).Run(ctx, firstVersion.ID))
	course = db.Course.GetX(ctx, course.ID)
	require.Equal(t, firstVersion.ID, *course.CurrentVersionID)

	// v2 fails: an empty checkout has no spec.yml, so Parse errors.
	badLoader := newLoaderWith(t, db, txStore, repoFetcher(t.TempDir()))
	secondVersion := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreatedSecond))
	require.Error(t, badLoader.Run(ctx, secondVersion.ID))

	// v2 is failed with the error recorded...
	secondVersion = db.CourseVersion.GetX(ctx, secondVersion.ID)
	assert.Equal(t, "failed", derefStr(secondVersion.State))
	assert.Contains(t, derefStr(secondVersion.Result), "Error")

	// ...and the live version is STILL v1 — the broken rebuild changed nothing.
	course = db.Course.GetX(ctx, course.ID)
	assert.Equal(t, firstVersion.ID, *course.CurrentVersionID)
}

func TestLoaderRejectsUnsupportedTheoryImage(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	repo := t.TempDir()
	require.NoError(t, os.CopyFS(repo, os.DirFS(fixtureRepo(t))))

	imagePath := filepath.Join(
		repo,
		"modules", "10-basics", "10-hello-world", "en", "assets", "dart.png",
	)
	require.NoError(t, os.WriteFile(imagePath, []byte("<svg><script>alert(1)</script></svg>"), 0o644))

	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))
	loader := newLoaderWith(t, db, txStore, repoFetcher(repo))

	err := loader.Run(ctx, version.ID)

	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported asset media type")

	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "failed", derefStr(version.State))
}

func TestLoaderRejectsTheoryImageSymlinkEscapingCheckout(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	repo := t.TempDir()
	require.NoError(t, os.CopyFS(repo, os.DirFS(fixtureRepo(t))))

	imagePath := filepath.Join(
		repo,
		"modules", "10-basics", "10-hello-world", "en", "assets", "dart.png",
	)
	image, err := os.ReadFile(imagePath)
	require.NoError(t, err)

	externalImagePath := filepath.Join(t.TempDir(), "external.png")
	require.NoError(t, os.WriteFile(externalImagePath, image, 0o644)) //nolint:gosec // G703: path is under t.TempDir
	require.NoError(t, os.Remove(imagePath))
	require.NoError(t, os.Symlink(externalImagePath, imagePath))

	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))
	loader := newLoaderWith(t, db, txStore, repoFetcher(repo))

	err = loader.Run(ctx, version.ID)

	require.Error(t, err)
	require.ErrorContains(t, err, "open image")

	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "failed", derefStr(version.State))
}

// The course repos keep a lesson's images in `<lesson>/assets`, shared by every
// locale, and refer to them as `../assets/x.png` from the locale's theory.
func TestLoaderUploadsTheoryImagesFromTheLessonAssetsDir(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	repo := t.TempDir()
	require.NoError(t, os.CopyFS(repo, os.DirFS(fixtureRepo(t))))

	lesson := filepath.Join(repo, "modules", "10-basics", "10-hello-world")
	require.NoError(t, os.Rename(filepath.Join(lesson, "en", "assets"), filepath.Join(lesson, "assets")))
	readme := filepath.Join(lesson, "en", "README.md")
	body, err := os.ReadFile(readme)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(readme, []byte(strings.ReplaceAll(string(body), "./assets/", "../assets/")), 0o644)) //nolint:gosec // G703: path is under t.TempDir

	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))

	require.NoError(t, newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID))

	infos := db.CourseLessonTranslation.Query().Where(courselessontranslation.CourseVersionID(version.ID)).AllX(ctx)
	en := lessonInfoByLocale(t, infos, "en")
	assert.Contains(t, derefStr(en.Theory), "http://localhost:3001/api/storage/")
	assert.NotContains(t, derefStr(en.Theory), "../assets/")
}

func TestLoaderRejectsTheoryImageClimbingOutOfTheLesson(t *testing.T) {
	db, txStore := testsupport.NewClientWithTransactor(t)
	ctx := context.Background()
	repo := t.TempDir()
	require.NoError(t, os.CopyFS(repo, os.DirFS(fixtureRepo(t))))

	lesson := filepath.Join(repo, "modules", "10-basics", "10-hello-world")
	picture, err := os.ReadFile(filepath.Join(lesson, "en", "assets", "dart.png"))
	require.NoError(t, err)
	// A real image, but outside the lesson: only the reference's path is at fault.
	require.NoError(t, os.WriteFile(filepath.Join(repo, "outside.png"), picture, 0o644)) //nolint:gosec // G703: path is under t.TempDir
	readme := filepath.Join(lesson, "en", "README.md")
	body, err := os.ReadFile(readme)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(readme, []byte(strings.ReplaceAll(string(body), "./assets/dart.png", "../../../../outside.png")), 0o644)) //nolint:gosec // G703: path is under t.TempDir

	version := db.CourseVersion.GetX(ctx, fixture.ID(fixture.VersionLoadableCreated))

	err = newLoaderWith(t, db, txStore, repoFetcher(repo)).Run(ctx, version.ID)

	require.Error(t, err)
	require.ErrorContains(t, err, "open image")
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func lessonInfoByLocale(t *testing.T, infos []*ent.CourseLessonTranslation, locale string) *ent.CourseLessonTranslation {
	t.Helper()

	for _, in := range infos {
		if in.Locale != nil && *in.Locale == locale {
			return in
		}
	}

	t.Fatalf("no lesson info for locale %q", locale)

	return nil
}
