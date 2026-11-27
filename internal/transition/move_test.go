package transition_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
	"hexletbasics/internal/testsupport/testdb"
	"hexletbasics/internal/transition"
)

// The generated state writer (ent/template/statemove.tmpl) is one template for
// every StateMove field, so it is tested once, here, on blog posts; the
// resources' own tests cover the operations, their Permissions and their Audit
// Records.

func TestALegalMoveWritesTheNewState(t *testing.T) {
	t.Parallel()
	db := testsupport.NewClient(t)
	id := fixture.ID(fixture.BlogPostSecondPost)

	require.NoError(t, db.BlogPost.MoveState(t.Context(), id, transition.Publish()))

	assert.Equal(t, transition.Published, db.BlogPost.GetX(t.Context(), id).State)
}

func TestAnIllegalMoveSaysWhyAndWritesNothing(t *testing.T) {
	t.Parallel()
	db := testsupport.NewClient(t)
	id := fixture.ID(fixture.BlogPostSecondPost)

	err := db.BlogPost.MoveState(t.Context(), id, transition.Unpublish())

	illegal, ok := transition.IsIllegal(err)
	require.True(t, ok, "a draft cannot be unpublished: %v", err)
	assert.Equal(t, transition.State(transition.Draft), illegal.From)
	assert.Equal(t, transition.State(transition.Draft), illegal.To)
	require.ErrorIs(t, err, domain.ErrConflict)
	assert.Equal(t, transition.Draft, db.BlogPost.GetX(t.Context(), id).State)
}

func TestAMoveOfAMissingRowIsNotFound(t *testing.T) {
	t.Parallel()
	db := testsupport.NewClient(t)

	err := db.Course.MoveReadiness(t.Context(), 999999, transition.MarkCompleted())

	require.ErrorIs(t, store.ToDomain(err), domain.ErrNotFound)
}

// Two moves racing on one row: the first holds the row until it commits, and
// the second, written under the condition that the state is still one it
// starts from, finds it published and is refused. Without the condition both
// would pass the check and both would claim draft to published.
func TestOfTwoRacingMovesOnlyOneWins(t *testing.T) {
	t.Parallel()
	db := store.NewClient(testdb.Clone(t))
	ctx := context.Background()
	id := fixture.ID(fixture.BlogPostSecondPost)

	first, err := db.Tx(ctx)
	require.NoError(t, err)
	require.NoError(t, first.BlogPost.MoveState(ctx, id, transition.Publish()))

	second := make(chan error, 1)

	go func() {
		tx, err := db.Tx(ctx)
		if err != nil {
			second <- err

			return
		}

		err = tx.BlogPost.MoveState(ctx, id, transition.Publish())
		second <- err

		_ = tx.Rollback()
	}()

	require.NoError(t, first.Commit())

	err = <-second
	illegal, ok := transition.IsIllegal(err)
	require.True(t, ok, "the second move is refused: %v", err)
	assert.Equal(t, transition.State(transition.Published), illegal.From)
	assert.Equal(t, transition.Published, db.BlogPost.GetX(ctx, id).State)
}
