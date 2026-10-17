// Package transitions is the legal-move table of a Publication State, written
// once for every resource that has one (blog posts, then banners, reviews,
// landing pages, ...). A state is not an update: each change of state is its own
// contract operation (publish, unpublish, archive), and each of those declares
// the states it may start from. The table lives here, in code the handlers call,
// so the website and the MCP server are bound by the same rules.
package transitions

import (
	"errors"
	"fmt"
	"slices"
)

// Transition is one operation's rule: it moves a resource to To, and only from
// one of the From states.
//
// Repeating the current state is refused, not a no-op (ADR-0017, state
// transitions addendum): To is never in From, so "publish a published post" is
// an illegal move like any other. A client that retries after a lost answer
// learns that the state is already there, and the Audit Record never claims a
// change that did not happen.
type Transition[S ~string] struct {
	To   S
	From []S
}

// Illegal is the typed error of a refused move: the state the resource is in and
// the one the operation would have set. Handlers map it to a 409 with a
// localized explanation.
type Illegal struct {
	From, To string
}

func (e *Illegal) Error() string {
	return fmt.Sprintf("illegal state transition from %q to %q", e.From, e.To)
}

// IsIllegal reports whether err is a refused move, and returns it.
func IsIllegal(err error) (*Illegal, bool) {
	var illegal *Illegal
	if errors.As(err, &illegal) {
		return illegal, true
	}
	return nil, false
}

// Check says whether a resource in the current state may take the transition.
func (t Transition[S]) Check(current S) error {
	if slices.Contains(t.From, current) {
		return nil
	}
	return &Illegal{From: string(current), To: string(t.To)}
}
