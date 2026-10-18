package transition_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/transition"
)

type state string

func TestATransitionAcceptsOnlyItsSourceStates(t *testing.T) {
	t.Parallel()

	publish := transition.Transition[state]{To: "published", From: []state{"draft", "archived"}}

	require.NoError(t, publish.Check("draft"))
	require.NoError(t, publish.Check("archived"))

	err := publish.Check("published")
	illegal, ok := transition.IsIllegal(err)
	require.True(t, ok, "repeating the current state is refused")
	assert.Equal(t, transition.State("published"), illegal.From)
	assert.Equal(t, transition.State("published"), illegal.To)
	assert.Error(t, err)
}

func TestAnUnknownStateIsRefused(t *testing.T) {
	t.Parallel()

	archive := transition.Transition[state]{To: "archived", From: []state{"draft"}}

	_, ok := transition.IsIllegal(archive.Check("something"))
	assert.True(t, ok)
	_, ok = transition.IsIllegal(nil)
	assert.False(t, ok)
}

func TestThePublicationStatesHaveOneTableForEveryResource(t *testing.T) {
	t.Parallel()

	moves := map[string]transition.Transition[transition.PublicationState]{
		"publish": transition.Publish(), "unpublish": transition.Unpublish(), "archive": transition.Archive(),
	}

	legal := map[string][]transition.PublicationState{
		"publish":   {transition.Draft, transition.Archived},
		"unpublish": {transition.Published},
		"archive":   {transition.Draft, transition.Published},
	}
	for name, move := range moves {
		for _, from := range []transition.PublicationState{transition.Draft, transition.Published, transition.Archived} {
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
	assert.Equal(t, transition.Draft, transition.CurrentOf(nil))

	published := "published"
	assert.Equal(t, transition.Published, transition.CurrentOf(&published))
}
