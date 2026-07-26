package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListCourses checks the catalog ordering, which is pushed entirely into
// SQL: by the Course `order` (NULLS LAST), then Course id, then landing page id.
// The fixtures set course orders 1/2/null and give course 10 two listed pages so
// the landing-id tie-breaker is exercised; the non-listed page is excluded.
func TestListCourses(t *testing.T) {
	srv := newServer(t)

	items, err := srv.ListCourses(context.Background())
	require.NoError(t, err)

	ids := make([]int32, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	// 100: course 12 (order 1) · 101,104: course 10 (order 2, id tie-break) ·
	// 102: course 11 (order null, last). 103 is not listed, so it is absent.
	assert.Equal(t, []int32{100, 101, 104, 102}, ids)
}
