package main

import (
	"errors"
	"testing"

	"github.com/samber/oops"
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

	oopsErr, ok := errors.AsType[oops.OopsError](err)
	require.True(t, ok)
	assert.Equal(t, "UserRemoved", oopsErr.Context()["identifier"])
}
