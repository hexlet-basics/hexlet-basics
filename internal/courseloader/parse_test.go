package courseloader_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/courseloader"
)

// fixtureRepo is the committed minimal course (one module, one lesson, three
// locales) shared with the legacy Rails loader test — the canonical example of
// the on-disk format.
func fixtureRepo(t *testing.T) string {
	t.Helper()
	// internal/courseloader -> repo root is two up.
	return filepath.Join("..", "..", "legacy", "test", "fixtures", "files", "exercises")
}

func TestParseFixtureCourse(t *testing.T) {
	course, err := courseloader.Parse(fixtureRepo(t), "javascript")
	require.NoError(t, err)

	// spec.yml
	assert.Equal(t, "JavaScript", course.Spec.Name)
	assert.Equal(t, "completed", course.Spec.Progress)
	assert.Equal(t, "first_language", course.Spec.LearnAs)
	assert.Equal(t, "js", course.Spec.Extension)
	assert.Equal(t, "index.js", course.Spec.ExerciseFilename)
	assert.Equal(t, "test.js", course.Spec.ExerciseTestFilename)
	assert.Equal(t, "ghcr.io/hexlet-basics/exercises-javascript", course.Spec.DockerImage)

	// one module: 10-basics
	require.Len(t, course.Modules, 1)
	mod := course.Modules[0]
	assert.Equal(t, 10, mod.Order)
	assert.Equal(t, "basics", mod.Slug)

	// module infos: en, es, ru (sorted)
	require.Len(t, mod.Infos, 3)
	assert.Equal(t, []string{"en", "es", "ru"}, localesOf(mod.Infos))
	en := moduleInfo(t, mod.Infos, "en")
	assert.Equal(t, "JavaScript basics", en.Name)
	assert.Contains(t, en.Description, "most popular programming languages")

	// one lesson: 10-hello-world, first in the course
	require.Len(t, mod.Lessons, 1)
	lesson := mod.Lessons[0]
	assert.Equal(t, 10, lesson.Order)
	assert.Equal(t, "hello-world", lesson.Slug)
	assert.Equal(t, 1, lesson.NaturalOrder)
	assert.Equal(t, "/exercises-javascript/modules/10-basics/10-hello-world", lesson.PathToCode)

	// code artifacts
	assert.Equal(t, `console.log("Hello, World!");`, strings.TrimSpace(lesson.OriginalCode))
	assert.Contains(t, lesson.TestCode, "hello world")
	// no BEGIN/END markers -> empty template
	assert.Equal(t, "", lesson.PreparedCode)

	// lesson infos: en, es, ru
	require.Len(t, lesson.Infos, 3)
	assert.Equal(t, []string{"en", "es", "ru"}, lessonLocalesOf(lesson.Infos))
	li := lessonInfo(t, lesson.Infos, "en")
	assert.Equal(t, "Hello, World!", li.Name)
	assert.Contains(t, li.Theory, "Hello, World!")
	assert.Contains(t, li.Instructions, "Copy the exact code")
	require.Len(t, li.Tips, 1)
	assert.Contains(t, li.Tips[0], "Hello,")
}

func localesOf(infos []courseloader.ModuleInfo) []string {
	out := make([]string, len(infos))
	for i, in := range infos {
		out[i] = in.Locale
	}
	return out
}

func lessonLocalesOf(infos []courseloader.LessonInfo) []string {
	out := make([]string, len(infos))
	for i, in := range infos {
		out[i] = in.Locale
	}
	return out
}

func moduleInfo(t *testing.T, infos []courseloader.ModuleInfo, locale string) courseloader.ModuleInfo {
	t.Helper()
	for _, in := range infos {
		if in.Locale == locale {
			return in
		}
	}
	t.Fatalf("no module info for locale %q", locale)
	return courseloader.ModuleInfo{}
}

func lessonInfo(t *testing.T, infos []courseloader.LessonInfo, locale string) courseloader.LessonInfo {
	t.Helper()
	for _, in := range infos {
		if in.Locale == locale {
			return in
		}
	}
	t.Fatalf("no lesson info for locale %q", locale)
	return courseloader.LessonInfo{}
}
