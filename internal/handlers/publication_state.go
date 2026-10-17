package handlers

import (
	"context"
	"net/http"

	"hexletbasics/internal/localization"
	"hexletbasics/internal/transitions"
)

// A state is not an update (ADR-0017, state transitions addendum): publishing,
// unpublishing and archiving are operations of their own, and the legal moves
// are a transitions.Transition each. This file maps a refused move to the same
// kind of answer the self-protection refusals give, so every resource that
// follows the pattern answers the same way: 409 with a localized `detail`.

// publicationStateNames are the localized names of the Publication States
// (GLOSSARY), shared by every resource that has them.
var publicationStateNames = map[string]localization.Message{
	"draft":     localization.PublicationStateDraft,
	"published": localization.PublicationStatePublished,
	"archived":  localization.PublicationStateArchived,
}

// refuseTransition turns a transitions.Illegal into a 409 refusal naming both
// states in the caller's language; any other error (nil included) passes
// through untouched.
func refuseTransition(err error) error {
	illegal, ok := transitions.IsIllegal(err)
	if !ok {
		return err
	}
	return &refusal{
		status:  http.StatusConflict,
		message: localization.IllegalStateTransition,
		data: map[string]any{
			"From": stateName(illegal.From),
			"To":   stateName(illegal.To),
		},
		cause: err,
	}
}

// stateName is the localized name of a state, or the raw value of one the
// registry does not know (a legacy row with a state of its own).
func stateName(state string) any {
	if message, ok := publicationStateNames[state]; ok {
		return message
	}
	return state
}

// moveState runs one of a module's state moves and echoes the resource back
// through read; a refused move is mapped to its 409. The same three steps serve
// every resource that has a Publication State.
func moveState[T any](ctx context.Context, id int, move func(context.Context, int) error, read func(context.Context, int) (*T, error)) (*T, error) {
	if err := move(ctx, id); err != nil {
		return nil, refuseTransition(err)
	}
	return read(ctx, id)
}
