package testsupport

// Mocks of the seams the harness replaces: the LLM, the exercise runner, the
// domain-event transport and the job enqueuers. moq
// generates them from the production interfaces, so a changed signature breaks
// the build instead of drifting from a hand-kept fake. -stub makes an unset Func
// return zero values; the harness sets the defaults tests rely on.
//go:generate go tool moq -stub -out mock_assistant_gen.go -pkg testsupport ../assistant Streamer
//go:generate go tool moq -stub -out mock_progress_gen.go -pkg testsupport ../progress ExerciseRunner
//go:generate go tool moq -stub -out mock_events_gen.go -pkg testsupport ../events TxPublisher StandalonePublisher
//go:generate go tool moq -stub -out mock_handlers_gen.go -pkg testsupport ../handlers VersionBuildStarter LessonReviewEnqueuer RelatedCoursesSuggestionEnqueuer AccountEmailEnqueuer
