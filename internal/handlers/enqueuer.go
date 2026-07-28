package handlers

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// JobEnqueuer is the slice of the river client the handlers need: inserting a
// background job. Depending on this interface (rather than *river.Client
// directly) keeps handler tests off the pgx pool river requires — the test
// harness supplies a recording fake, while production wires the real client,
// which satisfies this signature exactly.
type JobEnqueuer interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}
