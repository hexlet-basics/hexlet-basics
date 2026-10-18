package courseloader

import (
	"os"
	"path/filepath"
	"strings"
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

	_, err := parseSpec(openRepo(t, dir))

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

	_, err := parseModuleInfos(openRepo(t, dir), ".")

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

	_, err := parseLessonInfos(openRepo(t, dir), ".")

	require.Error(t, err)
	assert.ErrorContains(t, err, "field definition not found")
}

// Every missing required key is named, in file order, so one build failure
// shows the whole fix instead of a random first offender.
func TestParseSpecReportsEveryMissingKeyInOrder(t *testing.T) {
	dir := t.TempDir()
	spec := "language:\n  name: JavaScript\n  extension: js\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "spec.yml"), []byte(spec), 0o644))

	_, err := parseSpec(openRepo(t, dir))

	require.Error(t, err)

	message := err.Error()
	for _, key := range []string{"progress", "learn_as", "docker_image", "exercise_filename", "exercise_test_filename"} {
		assert.Contains(t, message, "language."+key+" is required")
	}

	assert.NotContains(t, message, "language.name")
	assert.Less(t, strings.Index(message, "language.progress"), strings.Index(message, "language.exercise_test_filename"))
}

// Only a reference with a scheme or host is remote: a local image whose name
// merely starts with "http" must still be uploaded.
func TestIsRemoteRef(t *testing.T) {
	assert.True(t, isRemoteRef("https://example.com/a.png"))
	assert.True(t, isRemoteRef("http://example.com/a.png"))
	assert.True(t, isRemoteRef("//cdn.example.com/a.png"))
	assert.True(t, isRemoteRef("data:image/png;base64,AAAA"))
	assert.False(t, isRemoteRef("http_flow.png"))
	assert.False(t, isRemoteRef("./assets/a.png"))
	assert.False(t, isRemoteRef("../assets/a.png"))
}

// openRepo opens dir the way Parse does: confined to it.
func openRepo(t *testing.T, dir string) repository {
	t.Helper()

	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	t.Cleanup(func() { _ = root.Close() })

	return repository{fsys: root.FS(), dir: dir}
}

// A cloned repo is untrusted: a symlink that leaves it must not put the host's
// files into the database.
func TestParseRefusesASymlinkOutOfTheRepo(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "secret.yml")
	require.NoError(t, os.WriteFile(outside, []byte("language:\n  name: Leaked\n"), 0o644))
	dir := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(dir, "spec.yml")))

	_, err := Parse(dir, "js")

	require.Error(t, err)
	assert.NotContains(t, err.Error(), "Leaked")
}
