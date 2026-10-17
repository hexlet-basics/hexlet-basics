package localization_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/internal/localization"
)

// Cross-reference: the frontend checks the same rule with PREFIX_CASES in
// src/lib/locale-path.test.ts. Keep the (locale, prefix) table identical on
// both sides: English, the default, carries no prefix, every other locale
// carries `/<locale>`.
func TestPathPrefix(t *testing.T) {
	cases := []struct{ locale, prefix string }{
		{"en", ""},
		{"ru", "/ru"},
		{"es", "/es"},
	}
	for _, c := range cases {
		t.Run(c.locale, func(t *testing.T) {
			assert.Equal(t, c.prefix, localization.PathPrefix(c.locale))
		})
	}
}
