package transitions_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/transitions"
)

type state string

func TestATransitionAcceptsOnlyItsSourceStates(t *testing.T) {
	t.Parallel()
	publish := transitions.Transition[state]{To: "published", From: []state{"draft", "archived"}}

	require.NoError(t, publish.Check("draft"))
	require.NoError(t, publish.Check("archived"))

	err := publish.Check("published")
	illegal, ok := transitions.IsIllegal(err)
	require.True(t, ok, "repeating the current state is refused")
	assert.Equal(t, "published", illegal.From)
	assert.Equal(t, "published", illegal.To)
	assert.Error(t, err)
}

func TestAnUnknownStateIsRefused(t *testing.T) {
	t.Parallel()
	archive := transitions.Transition[state]{To: "archived", From: []state{"draft"}}

	_, ok := transitions.IsIllegal(archive.Check("something"))
	assert.True(t, ok)
	_, ok = transitions.IsIllegal(nil)
	assert.False(t, ok)
}

func TestThePublicationStatesHaveOneTableForEveryResource(t *testing.T) {
	t.Parallel()
	moves := map[string]transitions.Transition[transitions.PublicationState]{
		"publish": transitions.Publish, "unpublish": transitions.Unpublish, "archive": transitions.Archive,
	}
	legal := map[string][]transitions.PublicationState{
		"publish":   {transitions.Draft, transitions.Archived},
		"unpublish": {transitions.Published},
		"archive":   {transitions.Draft, transitions.Published},
	}
	for name, move := range moves {
		for _, from := range []transitions.PublicationState{transitions.Draft, transitions.Published, transitions.Archived} {
			if slices.Contains(legal[name], from) {
				assert.NoError(t, move.Check(from), "%s from %s", name, from)
			} else {
				assert.Error(t, move.Check(from), "%s from %s", name, from)
			}
		}
	}
}

func TestARowWithoutAStateIsADraft(t *testing.T) {
	t.Parallel()
	assert.Equal(t, transitions.Draft, transitions.CurrentOf(nil))
	published := "published"
	assert.Equal(t, transitions.Published, transitions.CurrentOf(&published))
}
