package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
)

// requireProblem asserts that a generated-client call failed with a Problem
// Details answer of the given status (ADR-0028): every non-2xx is one, whatever
// the operation declared.
func requireProblem(t *testing.T, err error, status int) {
	t.Helper()
	_ = problemWithStatus(t, err, status)
}

// problemWithStatus is requireProblem for a test that reads the rest of the
// answer.
func problemWithStatus(t *testing.T, err error, status int) *api.ProblemDetailsStatusCode {
	t.Helper()
	require.Error(t, err)
	var problem *api.ProblemDetailsStatusCode
	require.ErrorAs(t, err, &problem)
	require.Equal(t, status, problem.StatusCode)
	require.Equal(t, int32(status), problem.Response.Status)
	return problem
}
