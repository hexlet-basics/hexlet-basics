package eventhandlers

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/events"
	"hexletbasics/internal/jobs"
)

type recordingInserter struct {
	args river.JobArgs
}

func (r *recordingInserter) Insert(
	_ context.Context,
	args river.JobArgs,
	_ *river.InsertOpts,
) (*rivertype.JobInsertResult, error) {
	r.args = args
	return &rivertype.JobInsertResult{}, nil
}

func TestLeadCreatedEnqueuesAmoCRMJob(t *testing.T) {
	inserter := &recordingInserter{}
	handler := LeadCreated(inserter)
	event := &events.LeadCreated{LeadID: 42, UserID: 7}

	require.NoError(t, handler.Handle(t.Context(), event))
	job, ok := inserter.args.(jobs.AmoCRMLeadArgs)
	require.True(t, ok)
	assert.Equal(t, *event, job.Event)
	assert.Equal(t, leadCreatedHandlerName, handler.HandlerName())
}
