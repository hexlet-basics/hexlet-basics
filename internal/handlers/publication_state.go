package handlers

import (
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

// readinessNames are the localized names of the Course Readiness values
// (GLOSSARY): a Course's Readiness is a state with a vocabulary of its own.
var readinessNames = map[string]localization.Message{
	"draft":          localization.ReadinessDraft,
	"in_development": localization.ReadinessInDevelopment,
	"completed":      localization.ReadinessCompleted,
}

// refuseTransition turns a transitions.Illegal into a 409 refusal naming both
// states in the caller's language; any other error (nil included) passes
// through untouched.
func refuseTransition(err error) error {
	return refuseTransitionNamed(err, publicationStateNames)
}

// refuseTransitionNamed is refuseTransition for a state with a vocabulary of its
// own, the Course Readiness: same 409, same message, the names from the map.
func refuseTransitionNamed(err error, names map[string]localization.Message) error {
	illegal, ok := transitions.IsIllegal(err)
	if !ok {
		return err
	}
	return &refusal{
		status:  http.StatusConflict,
		message: localization.IllegalStateTransition,
		data: map[string]any{
			"From": stateName(names, illegal.From),
			"To":   stateName(names, illegal.To),
		},
		cause: err,
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
