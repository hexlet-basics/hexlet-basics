package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentifierIsTheLabelInExportedForm(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "UserRemoved", identifier("user_removed"))
	assert.Equal(t, "CourseLesson2", identifier("course_lesson_2"))
}

func TestRenderNamesEveryLabel(t *testing.T) {
	t.Parallel()
	src, err := render([]string{"user_blank", "user_removed"})

	require.NoError(t, err)
	assert.Contains(t, string(src), `UserBlank   = LabelOf("user_blank")`)
	assert.Contains(t, string(src), `UserRemoved = LabelOf("user_removed")`)
}

func TestRenderRefusesTwoLabelsThatSpellOneVariable(t *testing.T) {
	t.Parallel()
	_, err := render([]string{"user_removed", "user__removed"})

	assert.ErrorContains(t, err, "both become UserRemoved")
}
