package courseloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validSpec = `language:
  name: JavaScript
  progress: completed
  learn_as: first_language
  docker_image: hexletbasics/exercises-javascript
  extension: js
  exercise_filename: exercise.js
  exercise_test_filename: test.js
`

// writeRepo lays out a minimal valid course repo: one module with one lesson in
// en. Tests break one piece of it at a time to reach each parse failure.
func writeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	write(t, filepath.Join(dir, "spec.yml"), validSpec)
	write(t, filepath.Join(dir, "modules", "10-basics", "description.en.yml"), "name: Basics\ndescription: Start\n")
	lesson := filepath.Join(dir, "modules", "10-basics", "10-hello")
	write(t, filepath.Join(lesson, "exercise.js"), "// BEGIN\nsolution();\n// END\n")
	write(t, filepath.Join(lesson, "test.js"), "test()\n")
	write(t, filepath.Join(lesson, "en", "data.yml"), "name: Hello\ndescription: First\ntips:\n  - one\ndefinitions:\n  - name: term\n    description: meaning\n")
	write(t, filepath.Join(lesson, "en", "README.md"), "theory")
	write(t, filepath.Join(lesson, "en", "EXERCISE.md"), "instructions")

	return dir
}

func write(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestParseReadsALessonWithItsStudentTemplate(t *testing.T) {
	course, err := Parse(writeRepo(t), "js")

	require.NoError(t, err)
	require.Len(t, course.Modules, 1)
	lesson := course.Modules[0].Lessons[0]
	assert.Equal(t, "hello", lesson.Slug)
	assert.Equal(t, 1, lesson.NaturalOrder)
	assert.Equal(t, "/exercises-js/modules/10-basics/10-hello", lesson.PathToCode)
	// The solution between BEGIN and END is blanked for the learner.
	assert.Contains(t, lesson.PreparedCode, "BEGIN (write your solution here)")
	assert.NotContains(t, lesson.PreparedCode, "solution();")
	require.Len(t, lesson.Infos, 1)
	assert.Equal(t, []string{"one"}, lesson.Infos[0].Tips)
	assert.Equal(t, "meaning", lesson.Infos[0].Definitions[0].Description)
}

func TestParseOrdersModulesAndLessonsByTheirNumericPrefix(t *testing.T) {
	dir := writeRepo(t)
	// "9-intro" sorts after "10-basics" lexically but before it numerically.
	write(t, filepath.Join(dir, "modules", "9-intro", "description.en.yml"), "name: Intro\ndescription: First\n")
	lesson := filepath.Join(dir, "modules", "9-intro", "2-second")
	write(t, filepath.Join(lesson, "exercise.js"), "x")
	write(t, filepath.Join(lesson, "test.js"), "x")
	write(t, filepath.Join(lesson, "en", "data.yml"), "name: Second\ndescription: d\n")
	write(t, filepath.Join(lesson, "en", "README.md"), "t")
	write(t, filepath.Join(lesson, "en", "EXERCISE.md"), "i")

	course, err := Parse(dir, "js")

	require.NoError(t, err)
	require.Len(t, course.Modules, 2)
	assert.Equal(t, "intro", course.Modules[0].Slug)
	assert.Equal(t, 1, course.Modules[0].Lessons[0].NaturalOrder)
	// Natural order runs across modules, not per module.
	assert.Equal(t, 2, course.Modules[1].Lessons[0].NaturalOrder)
}

func TestParseSkipsFoldersThatAreNotLocales(t *testing.T) {
	dir := writeRepo(t)
	lesson := filepath.Join(dir, "modules", "10-basics", "10-hello")
	write(t, filepath.Join(lesson, "assets", "image.png"), "png")
	write(t, filepath.Join(dir, "modules", "10-basics", "description.notes.txt"), "ignored")
	write(t, filepath.Join(dir, "modules", "10-basics", "description.yml"), "ignored: true")

	course, err := Parse(dir, "js")

	require.NoError(t, err)
	assert.Len(t, course.Modules[0].Infos, 1)
	assert.Len(t, course.Modules[0].Lessons[0].Infos, 1)
}

func TestParseRejectsABrokenRepo(t *testing.T) {
	lesson := func(dir string) string { return filepath.Join(dir, "modules", "10-basics", "10-hello") }

	tests := []struct {
		name    string
		corrupt func(t *testing.T, dir string)
		want    string
	}{
		{"missing spec", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(dir, "spec.yml")))
		}, "spec.yml"},
		{"spec field missing", func(t *testing.T, dir string) {
			t.Helper()
			write(t, filepath.Join(dir, "spec.yml"), "language:\n  name: JavaScript\n")
		}, "is required"},
		{"no modules dir", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.RemoveAll(filepath.Join(dir, "modules")))
		}, "read modules dir"},
		{"module dir without order", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Mkdir(filepath.Join(dir, "modules", "basics"), 0o755))
		}, "not <order>-<slug>"},
		{"module dir with non-numeric order", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Mkdir(filepath.Join(dir, "modules", "x-basics"), 0o755))
		}, "non-numeric order"},
		{"module without description", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(dir, "modules", "10-basics", "description.en.yml")))
		}, "has no description"},
		{"module description with unreadable yaml", func(t *testing.T, dir string) {
			t.Helper()
			write(t, filepath.Join(dir, "modules", "10-basics", "description.en.yml"), ": : :")
		}, "parse"},
		{"lesson dir without order", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Mkdir(filepath.Join(dir, "modules", "10-basics", "hello"), 0o755))
		}, "not <order>-<slug>"},
		{"lesson without test file", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(lesson(dir), "test.js")))
		}, "test file"},
		{"lesson without exercise file", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(lesson(dir), "exercise.js")))
		}, "exercise file"},
		{"lesson without locales", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.RemoveAll(filepath.Join(lesson(dir), "en")))
		}, "has no <locale>/data.yml"},
		{"locale without theory", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(lesson(dir), "en", "README.md")))
		}, "theory"},
		{"locale without instructions", func(t *testing.T, dir string) {
			t.Helper()
			require.NoError(t, os.Remove(filepath.Join(lesson(dir), "en", "EXERCISE.md")))
		}, "instructions"},
		{"locale data with unreadable yaml", func(t *testing.T, dir string) {
			t.Helper()
			write(t, filepath.Join(lesson(dir), "en", "data.yml"), ": : :")
		}, "parse"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := writeRepo(t)
			tc.corrupt(t, dir)

			_, err := Parse(dir, "js")

			require.Error(t, err)
			assert.ErrorContains(t, err, tc.want)
		})
	}
}

func TestPrepareCodeIsEmptyWithoutAnEditableRegion(t *testing.T) {
	assert.Empty(t, prepareCode("no markers here\n"))
}
