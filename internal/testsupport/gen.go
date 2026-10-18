package testsupport

// Mocks of the true external boundaries tests cannot run for real: the git
// remote course content is fetched from, the LLM provider (lessonreview.Completer for free-text completions,
// relatedcourse.JSONCompleter for schema-constrained ones) and the Docker exercise
// runner. moq generates them from the
// production interfaces so a changed signature breaks the build. Tests only set
// a mock's Func to say what the boundary answers; they never assert on what it
// was called with. Everything else the handlers touch is the real
// implementation over the test transaction (see services.go).
//go:generate go tool moq -stub -out mock_assistant_gen.go -pkg testsupport ../assistant Streamer
//go:generate go tool moq -stub -out mock_courseloader_gen.go -pkg testsupport ../courseloader Fetcher
//go:generate go tool moq -stub -out mock_lessonreviews_gen.go -pkg testsupport ../lessonreview Completer
//go:generate go tool moq -stub -out mock_relatedcourses_gen.go -pkg testsupport ../relatedcourse JSONCompleter
//go:generate go tool moq -stub -out mock_progress_gen.go -pkg testsupport ../progress ExerciseRunner
