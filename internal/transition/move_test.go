package transition_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
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

// Two moves racing on one row. The second has read the row as a draft and is
// about to write while the first, which published it, has not committed yet: a
// hook holds the second's mutation until the first commits. A check followed
// by an unconditional write would now pass, and both would claim draft to
// published. The write carries the condition instead, so the second finds the
// row published and is refused.
func TestOfTwoRacingMovesOnlyOneWins(t *testing.T) {
	t.Parallel()
	db := store.NewClient(testdb.Clone(t))
	ctx := context.Background()
	id := fixture.ID(fixture.BlogPostSecondPost)

	first, err := db.Tx(ctx)
	require.NoError(t, err)
	require.NoError(t, first.BlogPost.MoveState(ctx, id, transition.Publish()))

	reached, committed := make(chan struct{}), make(chan struct{})

	var once sync.Once

	db.BlogPost.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			once.Do(func() { close(reached) })
			<-committed

			return next.Mutate(ctx, m)
		})
	})

	second := make(chan error, 1)

	go func() {
		tx, err := db.Tx(ctx)
		if err != nil {
			second <- err

			return
		}

		second <- tx.BlogPost.MoveState(ctx, id, transition.Publish())

		_ = tx.Rollback()
	}()

	<-reached
	require.NoError(t, first.Commit())
	close(committed)

	err = <-second
	illegal, ok := transition.IsIllegal(err)
	require.True(t, ok, "the second move is refused: %v", err)
	assert.Equal(t, transition.State(transition.Published), illegal.From)
	assert.Equal(t, transition.Published, db.BlogPost.GetX(ctx, id).State)
}
