package handlers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// Every non-2xx answer is RFC 9457 Problem Details, whatever the operation
// declared (ADR-0028). Each case is one status the clients meet, reached through
// a real operation, and read back through the generated client.

func TestAMissingRecordAnswersProblemDetails404(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.GetCourse(t.Context(), api.GetCourseParams{Slug: "no-such-course"})

	problem := problemWithStatus(t, err, http.StatusNotFound)
	assert.Equal(t, "about:blank", problem.Response.Type)
	assert.NotEmpty(t, problem.Response.Title)
	assert.False(t, problem.Response.Errors.IsSet(), "only a 422 carries field errors")
}

func TestASignedOutCallAnswersProblemDetails401(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.GetProfile(t.Context())

	problem := problemWithStatus(t, err, http.StatusUnauthorized)
	assert.NotEmpty(t, problem.Response.Title)
}

func TestARefusalAnswersProblemDetails403WithADetail(t *testing.T) {
	h := testsupport.NewHarness(t)

	err := h.Client.AdminDeleteUser(t.Context(), api.AdminDeleteUserParams{ID: int32(aliceID())})

	problem := problemWithStatus(t, err, http.StatusForbidden)
	assert.NotEmpty(t, problem.Response.Detail.Value, "the refusal is explained in the caller's language")
}

func TestBadInputAnswersProblemDetails422WithFieldErrors(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)
	taken := h.DB.User.Query().FirstX(t.Context())

	_, err := h.Client.CreateUser(t.Context(), &api.SignUpInput{
		Email:    *taken.Email,
		Password: "correct horse battery staple",
	})

	problem := problemWithStatus(t, err, http.StatusUnprocessableEntity)
	assert.Len(t, problem.Response.Errors.Value["email"], 1)
}

func TestAUniqueViolationOnAnUndeclaredFieldAnswersProblemDetails409(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminCreateCourseCategory(t.Context(), &api.CourseCategoryInput{
		Name:        "Programming", // the name the fixtures already use
		Header:      "Some other header",
		Slug:        "some-other-slug",
		Description: api.NilString{Null: true},
	})

	problem := problemWithStatus(t, err, http.StatusConflict)
	assert.False(t, problem.Response.Errors.IsSet(), "a conflict names no field")
}
