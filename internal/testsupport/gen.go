package testsupport

// Mocks of the two true external boundaries the harness cannot run for real:
// the LLM provider and the Docker exercise runner. moq generates them from the
// production interfaces so a changed signature breaks the build. Tests only set
// a mock's Func to say what the boundary answers; they never assert on what it
// was called with. Everything else the handlers touch is the real
// implementation over the test transaction (see services.go).
//go:generate go tool moq -stub -out mock_assistant_gen.go -pkg testsupport ../assistant Streamer
//go:generate go tool moq -stub -out mock_progress_gen.go -pkg testsupport ../progress ExerciseRunner
