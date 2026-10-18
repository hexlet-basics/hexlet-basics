package relatedcourse

import (
	"context"

	"hexletbasics/ent"
)

// OfferedCandidates exposes the candidate query to the external tests, which
// cannot be internal because testsupport imports this package.
func OfferedCandidates(ctx context.Context, db *ent.Client, locale string) ([]candidate, error) {
	return offeredCandidates(ctx, db, locale)
}

// Candidate names the unexported candidate type for the external tests.
type Candidate = candidate
