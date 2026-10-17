package jobs_test

import (
	"log/slog"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/jobs"
	"hexletbasics/internal/testsupport"
)

// Each job kind lands on the queue that isolates its class of work.
func TestJobKindsLandOnTheirQueues(t *testing.T) {
	_, tx := testsupport.NewClientWithTx(t)
	client, err := jobs.NewInsertOnlyClient(nil, slog.New(slog.DiscardHandler), nil, nil)
	require.NoError(t, err)

	for args, queue := range map[river.JobArgs]string{
		jobs.AccountEmailArgs{Purpose: "magic_link", UserID: 1}: jobs.QueueMail,
		jobs.ExerciseLoaderArgs{VersionID: 1}:                   jobs.QueueBuilds,
		jobs.ReapStuckVersionBuildsArgs{}:                       jobs.QueueBuilds,
		jobs.ReviewLessonArgs{LessonInfoID: 1}:                  jobs.QueueLLM,
		jobs.SuggestRelatedCoursesArgs{BlogPostID: 1}:           jobs.QueueLLM,
	} {
		res, err := client.InsertTx(t.Context(), tx, args, nil)
		require.NoError(t, err)
		require.Equal(t, queue, res.Job.Queue, args.Kind())
	}
}
