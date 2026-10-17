package mcpserver

import "testing"

// AllowForTest replaces the allowlist for the test that follows, to exercise the
// pipeline with an operation the registry does not offer yet. It must run before
// the server is built, which reads the list once.
func AllowForTest(t *testing.T, operationIDs ...string) {
	t.Helper()
	was := adminOperations
	t.Cleanup(func() { adminOperations = was })
	adminOperations = operationIDs
}

// SessionCookie is the cookie name the dispatcher puts the placeholder under.
const SessionCookie = sessionCookie
