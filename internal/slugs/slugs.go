// Package slugs builds URL-safe slugs.
package slugs

import "github.com/gosimple/slug"

// Make converts s into a URL-safe slug.
func Make(s string) string {
	return slug.Make(s)
}
