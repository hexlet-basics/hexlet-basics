package relatedcourse

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The prompt is the legacy one: the post's plain text, then the candidate
// courses as {id, name} JSON.
func TestSuggestPromptIsTheLegacyMessage(t *testing.T) {
	prompt, err := suggestPrompt("<p>Hello <b>world</b> from\n  the blog</p>", []candidate{
		{ID: 207281424, Name: lo.ToPtr("Курс Ruby")},
	})
	require.NoError(t, err)

	assert.Equal(t, `Текст статьи: Hello world from the blog

Список курсов: [{"id":207281424,"name":"Курс Ruby"}]`, prompt)
}

// Legacy `truncate(2000)`: 1997 characters plus "...", counted in characters,
// not bytes — the posts are Russian.
func TestSuggestPromptTruncatesTheTextInCharacters(t *testing.T) {
	prompt, err := suggestPrompt("<p>"+strings.Repeat("я", 3000)+"</p>", nil)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(prompt, "Текст статьи: "+strings.Repeat("я", 1997)+"...\n\n"))
}
