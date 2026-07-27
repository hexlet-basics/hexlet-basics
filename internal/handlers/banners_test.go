package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// bannerByBody resolves a fixture banner by its body text (a stable business
// key), reading its id back rather than hard-coding it — the same pattern the
// course-category suite uses for slugs.
func bannerByBody(t *testing.T, h *testsupport.Harness, body string) api.Banner {
	t.Helper()
	page, err := h.Client.AdminListBanners(context.Background(), api.AdminListBannersParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	for _, it := range page.Items {
		if it.Body == body {
			return it
		}
	}
	t.Fatalf("banner with body %q not found", body)
	return api.Banner{}
}

// banners.yml seeds these two.
const totalBanners = 2

func TestAdminListBanners(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListBanners(ctx, api.AdminListBannersParams{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, int32(totalBanners), page.Total)
	assert.Len(t, page.Items, totalBanners)
	// Newest first: ids strictly descending.
	for i := 1; i < len(page.Items); i++ {
		assert.Greater(t, page.Items[i-1].ID, page.Items[i].ID)
	}
}

func TestAdminListBannersPaginated(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListBanners(ctx, api.AdminListBannersParams{
		Page:    api.NewOptInt32(2),
		PerPage: api.NewOptInt32(1),
	})
	require.NoError(t, err)

	assert.Equal(t, int32(totalBanners), page.Total)
	assert.Equal(t, int32(2), page.Page)
	assert.Len(t, page.Items, 1) // 2 rows, page 2 of size 1 -> the 2nd row
}

func TestAdminGetBanner(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	want := bannerByBody(t, h, "Скидка 50% на все курсы")

	got, err := h.Client.AdminGetBanner(ctx, api.AdminGetBannerParams{ID: want.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, api.BannerLocaleRu, got.Locale)
	assert.Equal(t, api.BannerStatePublished, got.State)
	assert.Equal(t, "https://ru.hexlet.io/pricing", got.URL.Value)
}

func TestAdminGetBannerNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetBanner(ctx, api.AdminGetBannerParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateBanner(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	starts := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	created, err := h.Client.AdminCreateBanner(ctx, &api.BannerInput{
		State:      api.BannerStatePublished,
		Background: api.BannerBackgroundBlue,
		Locale:     api.BannerLocaleEn,
		Body:       "New course launch",
		URL:        api.NewNilString("https://hexlet.io"),
		StartsAt:   api.NewNilDateTime(starts),
		FinishesAt: api.NilDateTime{Null: true},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	assert.NotZero(t, created.ID)
	assert.Equal(t, "New course launch", created.Body)
	assert.Equal(t, api.BannerBackgroundBlue, created.Background)
	assert.True(t, created.FinishesAt.Null, "unset finishesAt should stay null")
	assert.False(t, created.CreatedAt.IsZero())

	total, err := h.DB.Banner.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalBanners+1, total)
}

// TestAdminCreateBannerEmptyBody is the parity proof for Rails' `presence: true`
// on body: the contract's `@minLength(1)` rejects an empty body at decode, which
// the central handler surfaces as 400 (not a 500). The generated client does not
// validate the request body, so the empty body reaches the server.
func TestAdminCreateBannerEmptyBody(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminCreateBanner(ctx, &api.BannerInput{
		State:      api.BannerStateDraft,
		Background: api.BannerBackgroundCtaGradient,
		Locale:     api.BannerLocaleRu,
		Body:       "",
		URL:        api.NilString{Null: true},
		StartsAt:   api.NilDateTime{Null: true},
		FinishesAt: api.NilDateTime{Null: true},
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, h.LastStatus())
}

func TestAdminUpdateBanner(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// The ru fixture has a url set; the update clears it (null) and changes body.
	target := bannerByBody(t, h, "Скидка 50% на все курсы")

	updated, err := h.Client.AdminUpdateBanner(ctx, &api.BannerInput{
		State:      api.BannerStateArchived,
		Background: target.Background,
		Locale:     target.Locale,
		Body:       "Промоакция завершена",
		URL:        api.NilString{Null: true}, // clears the column
		StartsAt:   api.NilDateTime{Null: true},
		FinishesAt: api.NilDateTime{Null: true},
	}, api.AdminUpdateBannerParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Промоакция завершена", updated.Body)
	assert.Equal(t, api.BannerStateArchived, updated.State)
	assert.True(t, updated.URL.Null)

	// A null url clears the column in the database.
	row, err := h.DB.Banner.Get(ctx, int(target.ID))
	require.NoError(t, err)
	assert.Equal(t, "Промоакция завершена", row.Body)
	assert.Nil(t, row.URL, "null url should clear the column")
}

func TestAdminDeleteBanner(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := bannerByBody(t, h, "Black Friday sale")

	err := h.Client.AdminDeleteBanner(ctx, api.AdminDeleteBannerParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.Banner.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalBanners-1, total)
}
