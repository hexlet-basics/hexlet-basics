package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

// reviews.yml seeds these five.

// reviewByBody resolves a fixture review by its body (a stable business key).
func reviewByBody(t *testing.T, h *testsupport.Harness, body string) api.Review {
	t.Helper()
	page, err := h.Client.AdminListReviews(context.Background(), api.AdminListReviewsParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	for _, it := range page.Items {
		if it.Body.Value == body {
			return it
		}
	}
	t.Fatalf("review with body %q not found", body)
	return api.Review{}
}

func TestAdminListReviews(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalReviews := h.DB.Review.Query().CountX(ctx)

	page, err := h.Client.AdminListReviews(ctx, api.AdminListReviewsParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalReviews), page.Total)
	assert.Len(t, page.Items, totalReviews)
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListReviewsEmbedsCourseAndUser(t *testing.T) {
	h := testsupport.NewHarness(t)

	r := reviewByBody(t, h, "Отличный курс, всё понятно")

	// Embedded course (ruby) is loaded in full.
	assert.Equal(t, "ruby", r.Course.Slug)
	assert.Equal(t, r.CourseId, r.Course.ID)
	// Embedded user (user_admin_learner) as the public projection: id and display name.
	assert.Equal(t, "Admin Learner", r.User.Name.Value)
	// Review's own display name (not the user's).
	assert.Equal(t, "Ivan Petrov", r.FullName.Value)
	assert.Equal(t, api.ReviewStatePublished, r.State.Value)
	assert.Equal(t, api.ReviewLocaleRu, r.Locale.Value)
}

func TestAdminGetReview(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := reviewByBody(t, h, "Great course")

	got, err := h.Client.AdminGetReview(ctx, api.AdminGetReviewParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "John Doe", got.FullName.Value)
	assert.Equal(t, api.ReviewStateDraft, got.State.Value)
	assert.Equal(t, "Staff Editor", got.User.Name.Value)
}

func TestAdminGetReviewNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetReview(ctx, api.AdminGetReviewParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateReview(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalReviews := h.DB.Review.Query().CountX(ctx)

	created, err := h.Client.AdminCreateReview(ctx, &api.ReviewInput{
		Pinned:    api.NewNilBool(false),
		CourseId:  api.NewNilInt32(fixture.ID32(fixture.CourseRuby)),
		UserId:    api.NewNilInt32(fixture.ID32(fixture.UserBlank)),
		Body:      api.NewNilString("New review"),
		FirstName: api.NewNilString("Carol"),
		LastName:  api.NewNilString("Clark"),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "New review", created.Body.Value)
	assert.Equal(t, api.ReviewStateDraft, created.State.Value, "a new review is always a draft")
	assert.Equal(t, "Carol Clark", created.FullName.Value)
	// The response reloads the embedded associations.
	assert.Equal(t, "ruby", created.Course.Slug)
	assert.Equal(t, "Blank User", created.User.Name.Value, "the embedded user is the account, not the name typed on the review")

	total, err := h.DB.Review.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalReviews+1, total)
}

func TestAdminUpdateReview(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := reviewByBody(t, h, "Great course")

	updated, err := h.Client.AdminUpdateReview(ctx, &api.ReviewInput{
		Pinned:    api.NewNilBool(true),
		CourseId:  api.NewNilInt32(target.CourseId),
		UserId:    api.NewNilInt32(target.UserId),
		Body:      api.NewNilString("Edited body"),
		FirstName: api.NewNilString("Jonathan"),
		LastName:  api.NilString{Null: true}, // clears the column
	}, api.AdminUpdateReviewParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Edited body", updated.Body.Value)
	assert.Equal(t, api.ReviewStateDraft, updated.State.Value, "update never changes the state")
	assert.Equal(t, "Jonathan ", updated.FullName.Value) // last name cleared -> trailing space (Ruby join parity)

	row, err := h.DB.Review.Get(ctx, int(target.ID))
	require.NoError(t, err)
	assert.Nil(t, row.LastName, "null lastName should clear the column")
}

func TestAdminDeleteReview(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()
	totalReviews := h.DB.Review.Query().CountX(ctx)

	target := reviewByBody(t, h, "Great course")

	err := h.Client.AdminDeleteReview(ctx, api.AdminDeleteReviewParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.Review.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalReviews-1, total)
}
