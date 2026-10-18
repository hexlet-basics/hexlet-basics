// Package transition is the legal-move table of a Publication State, written
// once for every resource that has one (blog posts, then banners, reviews,
// landing pages, ...). A state is not an update: each change of state is its own
// contract operation (publish, unpublish, archive), and each of those declares
// the states it may start from. The table lives here, in code the handlers call,
// so the website and the MCP server are bound by the same rules.
package transition

import (
	"errors"
	"slices"

	"hexletbasics/internal/domain"
)

// Transition is one operation's rule: it moves a resource to To, and only from
// one of the From states.
//
// Repeating the current state is refused, not a no-op (ADR-0020, state
// transitions addendum): To is never in From, so "publish a published post" is
// an illegal move like any other. A client that retries after a lost answer
// learns that the state is already there, and the Audit Record never claims a
// change that did not happen.
type Transition[S ~string] struct {
	To   S
	From []S
}

// State is the raw value of a state, as stored in its column.
type State string

// Vocabulary names a set of states an adapter words by name. Each state type
// belongs to one, and the adapter's display names are keyed by it, so the two
// cannot drift apart: States lists what a vocabulary holds.
type Vocabulary string

// The vocabularies. A new state type adds its name here and its states to
// vocabularyStates.
const (
	PublicationStates Vocabulary = "publication_state"
	CourseReadiness   Vocabulary = "course_readiness"
)

// vocabularyStates is every state of each vocabulary.
func vocabularyStates() map[Vocabulary][]State {
	return map[Vocabulary][]State{
		PublicationStates: {State(Draft), State(Published), State(Archived)},
		CourseReadiness:   {State(ReadinessDraft), State(ReadinessInDevelopment), State(ReadinessCompleted)},
	}
}

// States lists the states of the vocabulary.
func (v Vocabulary) States() []State { return slices.Clone(vocabularyStates()[v]) }

// Named is implemented by a state type that belongs to a vocabulary.
type Named interface{ Vocabulary() Vocabulary }

// IllegalError is the typed error of a refused move: the state the resource is in and
// the one the operation would have set. It is a domain.ErrConflict, and the
// HTTP adapter words it from the states' names in their vocabulary. The message
// is static; the states are its fields.
type IllegalError struct {
	From, To State
	// Vocabulary names the set of states From and To belong to, for a state that
	// declares one (see Named), so an adapter picks the right display names.
	Vocabulary Vocabulary
}

// Unwrap makes a refused move a conflict: the resource's state forbids it,
// whoever asks.
func (e *IllegalError) Unwrap() error { return domain.ErrConflict }

func (e *IllegalError) Error() string { return "illegal state transition" }

// IsIllegal reports whether err is a refused move, and returns it.
func IsIllegal(err error) (*IllegalError, bool) {
	if illegal, ok := errors.AsType[*IllegalError](err); ok {
		return illegal, true
	}

	return nil, false
}

// Check says whether a resource in the current state may take the transition.
func (t Transition[S]) Check(current S) error {
	if slices.Contains(t.From, current) {
		return nil
	}

	illegal := &IllegalError{
		From:       State(current),
		To:         State(t.To),
		Vocabulary: "",
	}
	if named, ok := any(current).(Named); ok {
		illegal.Vocabulary = named.Vocabulary()
	}

	return illegal
}
