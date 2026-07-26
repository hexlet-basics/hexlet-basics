package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListCourses checks the catalog ordering, which is pushed entirely into
// SQL: by the Course `order` (NULLS LAST), then Course id, then landing page id.
//
// Fixtures come from the legacy Rails set (crc32 ids), so assertions are on
// business facts, not raw ids. The decisive one: typescript (order 1) sorts
// before javascript (order 2) even though typescript's Course id is the larger
// of the two — proving `order` drives the sort, not the id.
func TestListCourses(t *testing.T) {
	srv := newServer(t)

	items, err := srv.ListCourses(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, items)

	pos := make(map[string]int, len(items))
	for i, it := range items {
		pos[it.Slug] = i
	}

	// typescript is the unique minimum order (1), with a single landing page.
	assert.Equal(t, "typescript-ru", items[0].Slug)

	// Primary key beats the id tie-breaker: order 1 < order 2 even though
	// typescript's Course id > javascript's.
	tsIdx, ok := pos["typescript-ru"]
	require.True(t, ok, "typescript-ru missing from catalog")
	jsIdx, ok := pos["javascript-ru"]
	require.True(t, ok, "javascript-ru missing from catalog")
	assert.Less(t, tsIdx, jsIdx)

	// Tie-breaker: landing pages of the same Course appear in ascending landing
	// id order (and same-Course pages are contiguous, as the ordering groups by
	// Course id first).
	for i := 1; i < len(items); i++ {
		if items[i].Course.ID == items[i-1].Course.ID {
			assert.Less(t, items[i-1].ID, items[i].ID,
				"same-course pages must ascend by landing id")
		}
	}
}
