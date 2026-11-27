package handlers

import (
	"context"

	"hexletbasics/internal/transition"
)

// A state is not an update (ADR-0020, state transitions addendum): publishing,
// unpublishing and archiving are operations of their own, and the legal moves
// are a transition.Transition each. A refused move is a *transition.IllegalError, a
// domain conflict, so every resource that follows the pattern answers the same
// way: 409 with a localized `detail` (detailFor words it in the vocabulary of
// the state it names).

// moveState runs one transition through the generated state writer of the
// resource (Move<Field>, ent/template/statemove.tmpl) and echoes the resource
// back through read. The same two steps serve every resource that has a state:
// the Publication State and the Course Readiness alike.
func moveState[S ~string, T any](
	ctx context.Context,
	id int,
	write func(context.Context, int, transition.Transition[S]) error,
	move transition.Transition[S],
	read func(context.Context, int) (*T, error),
) (*T, error) {
	err := write(ctx, id, move)
	if err != nil {
		return nil, err
	}

	return read(ctx, id)
}
