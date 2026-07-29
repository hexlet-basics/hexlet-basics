package courseloader_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob/memblob"

	"hexletbasics/ent"
	"hexletbasics/ent/languagelesson"
	"hexletbasics/ent/languagelessonversion"
	"hexletbasics/ent/languagelessonversioninfo"
	"hexletbasics/ent/languagemodule"
	"hexletbasics/ent/languagemoduleversion"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/courseloader"
	"hexletbasics/internal/testsupport"
)

// fakeFetcher returns a fixed directory instead of cloning, so the loader runs
// fully offline against the committed fixture course.
type fakeFetcher struct{ dir string }

func (f fakeFetcher) Fetch(_ context.Context, _ string) (string, func(), error) {
	return f.dir, func() {}, nil
}

func newLoaderWith(t *testing.T, db *ent.Client, fetcher courseloader.Fetcher) *courseloader.Loader {
	t.Helper()
	bucket := memblob.OpenBucket(nil)
	t.Cleanup(func() { _ = bucket.Close() })
	assets := assetstore.New(db, bucket, "http://localhost:3001")
	return courseloader.NewLoader(db, assets, fetcher)
}

// newLoader builds a loader that fetches the committed fixture course.
func newLoader(t *testing.T, db *ent.Client) *courseloader.Loader {
	t.Helper()
	return newLoaderWith(t, db, fakeFetcher{dir: fixtureRepo(t)})
}

func TestLoaderBuildsAndPromotesVersion(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	// Arrange: a course with a freshly-created version, the state the loader
	// requires. (Its own writes roll back with the test transaction.)
	course := db.Course.Create().SetSlug("loader-test-lang").SetName("Loader Test").SaveX(ctx)
	version := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)

	loader := newLoader(t, db)
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
	mods := db.LanguageModule.Query().Where(languagemodule.LanguageID(course.ID)).AllX(ctx)
	require.Len(t, mods, 1)
	assert.Equal(t, "basics", derefStr(mods[0].Slug))
	mvs := db.LanguageModuleVersion.Query().Where(languagemoduleversion.LanguageVersionID(version.ID)).AllX(ctx)
	require.Len(t, mvs, 1)
	assert.Equal(t, 10, *mvs[0].Order)

	// One lesson version, first in the course, with prepared/original code.
	lvs := db.LanguageLessonVersion.Query().Where(languagelessonversion.LanguageVersionID(version.ID)).AllX(ctx)
	require.Len(t, lvs, 1)
	assert.Equal(t, 1, *lvs[0].NaturalOrder)
	assert.Equal(t, "/exercises-loader-test-lang/modules/10-basics/10-hello-world", derefStr(lvs[0].PathToCode))
	assert.Contains(t, derefStr(lvs[0].OriginalCode), "Hello, World!")

	// Three lesson infos (en/es/ru); theory images rewritten to /storage URLs.
	infos := db.LanguageLessonVersionInfo.Query().Where(languagelessonversioninfo.LanguageVersionID(version.ID)).AllX(ctx)
	require.Len(t, infos, 3)
	en := lessonInfoByLocale(t, infos, "en")
	assert.Equal(t, "Hello, World!", derefStr(en.Name))
	assert.Contains(t, derefStr(en.Theory), "http://localhost:3001/storage/")
	assert.NotContains(t, derefStr(en.Theory), "./assets/")
}

func TestLoaderSkipsNonCreatedVersion(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	course := db.Course.Create().SetSlug("loader-skip-lang").SetName("Skip").SaveX(ctx)
	// A version already 'built' must not be rebuilt.
	version := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("built").SaveX(ctx)

	loader := newLoader(t, db)
	require.NoError(t, loader.Run(ctx, version.ID))

	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "built", derefStr(version.State)) // unchanged
	assert.Contains(t, derefStr(version.Result), "Skipped")
	// No module rows were written for this course.
	n := db.LanguageModule.Query().Where(languagemodule.LanguageID(course.ID)).CountX(ctx)
	assert.Equal(t, 0, n)
}

func TestLoaderUpsertsStableLessonAcrossRebuilds(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	course := db.Course.Create().SetSlug("loader-rebuild-lang").SetName("Rebuild").SaveX(ctx)
	loader := newLoader(t, db)

	// First build.
	v1 := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, v1.ID))
	lesson1 := db.LanguageLesson.Query().Where(languagelesson.LanguageID(course.ID)).OnlyX(ctx)

	// Second build reuses the SAME lesson identity (learner progress FKs it).
	v2 := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, loader.Run(ctx, v2.ID))

	lessons := db.LanguageLesson.Query().Where(languagelesson.LanguageID(course.ID)).AllX(ctx)
	require.Len(t, lessons, 1)
	assert.Equal(t, lesson1.ID, lessons[0].ID)

	// Current version is the newest built one.
	course = db.Course.GetX(ctx, course.ID)
	assert.Equal(t, v2.ID, *course.CurrentVersionID)
}

// TestLoaderFailedRebuildKeepsLiveVersion is the load-bearing invariant (legacy
// ADR-0001): a broken rebuild must NOT take the live course offline. v1 builds and
// goes live; v2 fails to parse (empty repo, no spec.yml); the course must still
// point at v1, and v2 must be recorded as failed.
func TestLoaderFailedRebuildKeepsLiveVersion(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()

	course := db.Course.Create().SetSlug("loader-fail-lang").SetName("Fail").SaveX(ctx)

	// v1 builds successfully and is promoted live.
	v1 := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)
	require.NoError(t, newLoader(t, db).Run(ctx, v1.ID))
	course = db.Course.GetX(ctx, course.ID)
	require.Equal(t, v1.ID, *course.CurrentVersionID)

	// v2 fails: an empty checkout has no spec.yml, so Parse errors.
	badLoader := newLoaderWith(t, db, fakeFetcher{dir: t.TempDir()})
	v2 := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)
	require.Error(t, badLoader.Run(ctx, v2.ID))

	// v2 is failed with the error recorded...
	v2 = db.CourseVersion.GetX(ctx, v2.ID)
	assert.Equal(t, "failed", derefStr(v2.State))
	assert.Contains(t, derefStr(v2.Result), "Error")

	// ...and the live version is STILL v1 — the broken rebuild changed nothing.
	course = db.Course.GetX(ctx, course.ID)
	assert.Equal(t, v1.ID, *course.CurrentVersionID)
}

func TestLoaderRejectsUnsupportedTheoryImage(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := context.Background()
	repo := t.TempDir()
	require.NoError(t, os.CopyFS(repo, os.DirFS(fixtureRepo(t))))

	imagePath := filepath.Join(
		repo,
		"modules", "10-basics", "10-hello-world", "en", "assets", "dart.png",
	)
	require.NoError(t, os.WriteFile(imagePath, []byte("<svg><script>alert(1)</script></svg>"), 0o644))

	course := db.Course.Create().SetSlug("loader-invalid-image-lang").SetName("Invalid Image").SaveX(ctx)
	version := db.CourseVersion.Create().SetLanguageID(course.ID).SetState("created").SaveX(ctx)
	loader := newLoaderWith(t, db, fakeFetcher{dir: repo})

	err := loader.Run(ctx, version.ID)

	require.Error(t, err)
	assert.ErrorContains(t, err, "unsupported asset media type")
	version = db.CourseVersion.GetX(ctx, version.ID)
	assert.Equal(t, "failed", derefStr(version.State))
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func lessonInfoByLocale(t *testing.T, infos []*ent.LanguageLessonVersionInfo, locale string) *ent.LanguageLessonVersionInfo {
	t.Helper()
	for _, in := range infos {
		if in.Locale != nil && *in.Locale == locale {
			return in
		}
	}
	t.Fatalf("no lesson info for locale %q", locale)
	return nil
}
