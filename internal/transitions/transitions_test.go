package transitions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/transitions"
)

type state string

func TestATransitionAcceptsOnlyItsSourceStates(t *testing.T) {
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
	archive := transitions.Transition[state]{To: "archived", From: []state{"draft"}}

	_, ok := transitions.IsIllegal(archive.Check("something"))
	assert.True(t, ok)
	_, ok = transitions.IsIllegal(nil)
	assert.False(t, ok)
}
