package courseloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSpecRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	spec := `language:
  name: JavaScript
  progress: completed
  learn_as: first_language
  docker_image: hexletbasics/exercises-javascript
  extension: js
  exercise_filename: exercise.js
  exercise_test_filename: test.js
  exercise_test_filenam: typo.js
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.yml"), []byte(spec), 0o644))

	_, err := parseSpec(dir)

	require.Error(t, err)
	assert.ErrorContains(t, err, "field exercise_test_filenam not found")
}

func TestParseModuleInfosRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	description := `name: Introduction
description: Start here
descriptio: typo
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "description.en.yml"), []byte(description), 0o644))

	_, err := parseModuleInfos(dir)

	require.Error(t, err)
	assert.ErrorContains(t, err, "field descriptio not found")
}

func TestParseLessonInfosRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	localeDir := filepath.Join(dir, "en")
	require.NoError(t, os.Mkdir(localeDir, 0o755))
	data := `name: Hello
description: First lesson
tips: []
definitions: []
definition: typo
`
	require.NoError(t, os.WriteFile(filepath.Join(localeDir, "data.yml"), []byte(data), 0o644))

	_, err := parseLessonInfos(dir)

	require.Error(t, err)
	assert.ErrorContains(t, err, "field definition not found")
}
