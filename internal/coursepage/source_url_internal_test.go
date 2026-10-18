package coursepage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessonSourceURL(t *testing.T) {
	const base = "https://github.com/hexlet-basics"

	ptr := func(s string) *string { return &s }

	for _, tt := range []struct {
		name   string
		base   string
		path   *string
		locale string
		want   string
	}{
		{"modules become the blob path", base, ptr("exercises-javascript/modules/10-basics/20-variables"), "en",
			"https://github.com/hexlet-basics/exercises-javascript/blob/main/modules/10-basics/20-variables/en/README.md"},
		{"the locale picks the readme", base, ptr("exercises-go/modules/10-basics/10-hello"), "ru",
			"https://github.com/hexlet-basics/exercises-go/blob/main/modules/10-basics/10-hello/ru/README.md"},
		{"a trailing slash on the base is harmless", base + "/", ptr("exercises-go/modules/10-a/10-b"), "en",
			"https://github.com/hexlet-basics/exercises-go/blob/main/modules/10-a/10-b/en/README.md"},
		{"a path with no modules is linked as it is", base, ptr("exercises-go/10-a"), "en",
			"https://github.com/hexlet-basics/exercises-go/10-a/en/README.md"},
		{"no path, no link", base, nil, "en", ""},
		{"an empty path, no link", base, ptr(""), "en", ""},
		{"an unparsable base, no link", "://bad", ptr("a/modules/b"), "en", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lessonSourceURL(tt.base, tt.path, tt.locale))
		})
	}
}
