package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// ai_messages.yml seeds these three.
const totalAssistantMessages = 3

func TestAdminListMessages(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListMessages(ctx, api.AdminListMessagesParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalAssistantMessages), page.Total)
	assert.Len(t, page.Items, totalAssistantMessages)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}

	// Business facts from the fixtures, keyed by content (not raw ids).
	byContent := map[string]api.LessonAssistantMessage{}
	for _, it := range page.Items {
		byContent[it.Content] = it
	}

	question, ok := byContent["What is a variable?"]
	require.True(t, ok, "user question in the variables chat not found")
	assert.Equal(t, "user", question.Role)
	assert.False(t, question.UserId.Null, "a user question carries its author")
	assert.Equal(t, "javascript", question.CourseSlug)
	assert.Equal(t, "variables", question.CourseLessonSlug)
	assert.Equal(t, "Variables", question.CourseLessonName)
	assert.False(t, question.CreatedAt.IsZero())

	reply, ok := byContent["Use console.log to print a string."]
	require.True(t, ok, "assistant reply in the hello-world chat not found")
	assert.Equal(t, "assistant", reply.Role)
	assert.True(t, reply.UserId.Null, "an assistant reply has no author user")
	assert.Equal(t, "hello-world", reply.CourseLessonSlug)
	assert.Equal(t, "Hello, World!", reply.CourseLessonName)
}

func TestAdminListMessagesPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListMessages(ctx, api.AdminListMessagesParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(2),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalAssistantMessages), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1) // 3 rows, page 2 of size 2 -> the last row
}
