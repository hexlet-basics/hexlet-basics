package feeds_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/coursecategory"
	"hexletbasics/ent/landingpage"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/testsupport"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current output")

// TestYandexGolden compares the whole feed for the fixtures against
// testdata/yandex_courses.xml. Of the ru main landing pages only JavaScript's
// course has three ru modules, so it is the one offer: Ruby (one module) and
// the rest (none) are skipped. JavaScript's module versions carry ids out of
// `order` sequence, so the plan's order in the golden file is the SQL
// ordering at work. The fixtures leave every legacy category column empty;
// one is set here so `set-ids` is covered too.
//
// Regenerate with `go test ./internal/feeds/ -run TestYandexGolden -update`
// and review the diff.
func TestYandexGolden(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := t.Context()

	frontend := db.CourseCategory.Query().Where(coursecategory.SlugEQ("frontend-ru")).OnlyX(ctx)
	db.LandingPage.Update().
		Where(landingpage.SlugEQ("javascript-ru")).
		SetLanguageCategoryID(frontend.ID).
		ExecX(ctx)

	clock := func() time.Time { return time.Date(2026, 9, 26, 13, 31, 45, 0, time.UTC) }
	got, err := feeds.NewYandex(db, "code-basics.com", feeds.WithClock(clock)).Build(ctx)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "yandex_courses.xml")
	if *update {
		require.NoError(t, os.WriteFile(golden, got, 0o600))
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err)
	assert.Equal(t, string(want), string(got))
}
