package feeds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/course"
	"hexletbasics/internal/feeds"
	"hexletbasics/internal/testsupport"
)

// A course whose current version is not built yet has no modules to list, so
// its landing page is no offer, while the other courses' offers stay.
func TestYandexSkipsACourseWithoutACurrentVersion(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := t.Context()
	whole, err := feeds.NewYandex(db, "code-basics.com").Build(ctx)
	require.NoError(t, err)
	require.Contains(t, string(whole), "<offer ", "the fixtures offer one course")

	db.Course.Update().Where(course.SlugEQ("javascript")).ClearCurrentVersionID().ExecX(ctx)
	got, err := feeds.NewYandex(db, "code-basics.com").Build(ctx)

	require.NoError(t, err)
	assert.NotContains(t, string(got), "<offer ", "the only offered course lost its version")
}

// With no course built at all the feed is still a valid document: the shop
// header and the categories, and an empty list of offers.
func TestYandexWithNoBuiltCourseHasNoOffers(t *testing.T) {
	db := testsupport.NewClient(t)
	ctx := t.Context()
	db.Course.Update().ClearCurrentVersionID().ExecX(ctx)

	got, err := feeds.NewYandex(db, "code-basics.com").Build(ctx)

	require.NoError(t, err)
	assert.Contains(t, string(got), "<yml_catalog")
	assert.NotContains(t, string(got), "<offer ")
}
