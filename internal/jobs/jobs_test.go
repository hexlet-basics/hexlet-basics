package jobs_test

import (
	"hexletbasics/internal/events"
	"log/slog"
	"testing"
	"time"

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
		jobs.AccountEmailArgs{
			Purpose: "magic_link",
			UserID:  1,
			Locale:  "",
		}: jobs.QueueMail,
		jobs.ExerciseLoaderArgs{VersionID: 1}:         jobs.QueueBuilds,
		jobs.ReapStuckVersionBuildsArgs{}:             jobs.QueueBuilds,
		jobs.ReviewLessonArgs{LessonInfoID: 1}:        jobs.QueueLLM,
		jobs.SuggestRelatedCoursesArgs{BlogPostID: 1}: jobs.QueueLLM,
	} {
		res, err := client.InsertTx(t.Context(), tx, args, nil)
		require.NoError(t, err)
		require.Equal(t, queue, res.Job.Queue, args.Kind())
	}
}

// Every kind states how many tries it gets and why, instead of inheriting
// River's 25 (weeks of retries).
func TestJobKindsStateTheirAttempts(t *testing.T) {
	for args, attempts := range map[river.JobArgs]int{
		jobs.AccountEmailArgs{
			Purpose: "",
			UserID:  0,
			Locale:  "",
		}: 5,
		jobs.AmoCRMLeadArgs{
			Event: events.LeadCreated{
				LeadID:      0,
				UserID:      0,
				UserName:    "",
				FirstName:   nil,
				LastName:    nil,
				YMClientID:  nil,
				UTMSource:   nil,
				UTMMedium:   nil,
				UTMCampaign: nil,
				UTMTerm:     nil,
				UTMContent:  nil,
				Email:       nil,
				Phone:       nil,
				Telegram:    nil,
				WhatsApp:    nil,
				LandingPage: nil,
				Referrer:    nil,
				IP:          nil,
				OccurredAt:  time.Time{},
			},
		}: 8,
		jobs.ExerciseLoaderArgs{
			VersionID: 0,
		}: 1,
		jobs.ReviewLessonArgs{
			LessonInfoID: 0,
		}: 3,
		jobs.SuggestRelatedCoursesArgs{
			BlogPostID: 0,
		}: 3,
	} {
		opts := as[river.JobArgsWithInsertOpts](t, args).InsertOpts()
		require.Equal(t, attempts, opts.MaxAttempts, args.Kind())
	}
}
