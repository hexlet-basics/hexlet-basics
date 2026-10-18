package handlers

import (
	"context"
	"fmt"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
	"hexletbasics/internal/transition"
)

// A state is not an update (ADR-0020, state transitions addendum): publishing,
// unpublishing and archiving are operations of their own, and the legal moves
// are a transition.Transition each. This file maps a refused move to the same
// kind of answer the self-protection refusals give, so every resource that
// follows the pattern answers the same way: 409 with a localized `detail`.

// publicationStateNames are the localized names of the Publication States
// (GLOSSARY), shared by every resource that has them.
var publicationStateNames = map[string]localization.Message{
	"draft":     localization.PublicationStateDraft,
	"published": localization.PublicationStatePublished,
	"archived":  localization.PublicationStateArchived,
}

// readinessNames are the localized names of the Course Readiness values
// (GLOSSARY): a Course's Readiness is a state with a vocabulary of its own.
var readinessNames = map[string]localization.Message{
	"draft":          localization.ReadinessDraft,
	"in_development": localization.ReadinessInDevelopment,
	"completed":      localization.ReadinessCompleted,
}

// refuseTransitionNamed turns a transition.Illegal into a 409 refusal naming
// both states in the caller's language, the names from the map (the Publication
// States' or the Course Readiness'); any other error (nil included) passes
// through untouched.
func refuseTransitionNamed(err error, names map[string]localization.Message) error {
	illegal, ok := transition.IsIllegal(err)
	if !ok {
		return err
	}
	return &refusal{
		message: localization.IllegalStateTransition,
		data: map[string]any{
			"From": stateName(names, illegal.From),
			"To":   stateName(names, illegal.To),
		},
		cause: fmt.Errorf("%w: %w", domain.ErrConflict, err),
	}
}

// stateName is the localized name of a state, or the raw value of one the
// registry does not know (a legacy row with a state of its own).
func stateName(names map[string]localization.Message, state string) any {
	if message, ok := names[state]; ok {
		return message
	}
	return state
}

// moveState runs one of a module's state moves and echoes the resource back
// through read; a refused move is mapped to its 409. The same three steps serve
// every resource that has a Publication State.
func moveState[T any](ctx context.Context, id int, move func(context.Context, int) error, read func(context.Context, int) (*T, error)) (*T, error) {
	return moveNamedState(ctx, id, move, read, publicationStateNames)
}

// moveReadiness is moveState for the Course Readiness, which names its values
// in a vocabulary of its own.
func moveReadiness[T any](ctx context.Context, id int, move func(context.Context, int) error, read func(context.Context, int) (*T, error)) (*T, error) {
	return moveNamedState(ctx, id, move, read, readinessNames)
}

func moveNamedState[T any](ctx context.Context, id int, move func(context.Context, int) error, read func(context.Context, int) (*T, error), names map[string]localization.Message) (*T, error) {
	if err := move(ctx, id); err != nil {
		return nil, refuseTransitionNamed(err, names)
	}
	return read(ctx, id)
}
