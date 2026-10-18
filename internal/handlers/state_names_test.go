package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/internal/transition"
)

// Every state a vocabulary holds has a display name, and no name is kept for a
// state the vocabulary dropped, so the table cannot drift from transition.
func TestEveryStateOfEveryVocabularyHasADisplayName(t *testing.T) {
	t.Parallel()
	for _, vocabulary := range []transition.Vocabulary{transition.PublicationStates, transition.CourseReadiness} {
		states := vocabulary.States()
		for _, state := range states {
			assert.Contains(t, stateNames[vocabulary], state, "%s has no name for %s", vocabulary, state)
		}
		assert.Len(t, stateNames[vocabulary], len(states), "%s names a state it does not hold", vocabulary)
	}
}
