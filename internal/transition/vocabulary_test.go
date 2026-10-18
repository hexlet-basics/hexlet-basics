package transition_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/internal/transition"
)

// A refused move names its vocabulary, and every state a table can name is
// listed in it, so an adapter keyed by the vocabulary has a name for each.
func TestEveryStateOfATableBelongsToItsVocabulary(t *testing.T) {
	t.Parallel()

	publication := transition.PublicationStates.States()
	for _, move := range []transition.Transition[transition.PublicationState]{transition.Publish, transition.Unpublish, transition.Archive} {
		assert.Contains(t, publication, transition.State(move.To))
		for _, from := range move.From {
			assert.Contains(t, publication, transition.State(from))
		}
	}
	readiness := transition.CourseReadiness.States()
	for _, move := range []transition.Transition[transition.Readiness]{transition.MarkDraft, transition.MarkInDevelopment, transition.MarkCompleted} {
		assert.Contains(t, readiness, transition.State(move.To))
		for _, from := range move.From {
			assert.Contains(t, readiness, transition.State(from))
		}
	}

	illegal, ok := transition.IsIllegal(transition.Publish.Check(transition.Published))
	assert.True(t, ok)
	assert.Equal(t, transition.PublicationStates, illegal.Vocabulary)
}
