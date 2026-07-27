package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// language_landing_pages.yml (exported from legacy) seeds this many rows.
const totalLandingPages = 22

// go-ru landing page and its course, ruby course — all present in the fixtures.
const (
	landingGoRuID = 59175011
	courseRubyID  = 207281424
	courseGoID    = 912823128
)

// fullLandingInput returns a complete CourseLandingPageInput with every nullable
// field set to null. The contract fields are required-nullable, so the client
// always sends them all (an omitted enum would serialize as "" and 400; an
// omitted redirect id as 0 and hit the self-FK). Tests override the few fields
// they exercise.
func fullLandingInput(courseID int32) api.CourseLandingPageInput {
	return api.CourseLandingPageInput{
		CourseId:                  courseID,
		Slug:                      api.NilString{Null: true},
		Name:                      api.NilString{Null: true},
		Main:                      api.NilBool{Null: true},
		Listed:                    api.NilBool{Null: true},
		Footer:                    api.NilBool{Null: true},
		FooterName:                api.NilString{Null: true},
		State:                     api.NilLandingPageState{Null: true},
		Order:                     api.NilString{Null: true},
		LandingPageToRedirectId:   api.NilInt32{Null: true},
		MetaTitle:                 api.NilString{Null: true},
		MetaDescription:           api.NilString{Null: true},
		Header:                    api.NilString{Null: true},
		Description:               api.NilString{Null: true},
		UsedInHeader:              api.NilString{Null: true},
		UsedInDescription:         api.NilString{Null: true},
		OutcomesHeader:            api.NilString{Null: true},
		OutcomesDescription:       api.NilString{Null: true},
		OutcomesImageAttachmentId: api.NilInt32{Null: true},
	}
}

func landingBySlug(t *testing.T, h *testsupport.Harness, slug string) api.CourseLandingPage {
	t.Helper()
	page, err := h.Client.AdminListCourseLandingPages(context.Background(), api.AdminListCourseLandingPagesParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	for _, it := range page.Items {
		if it.Slug == slug {
			return it
		}
	}
	t.Fatalf("landing page with slug %q not found", slug)
	return api.CourseLandingPage{}
}

func TestAdminListCourseLandingPages(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	page, err := h.Client.AdminListCourseLandingPages(ctx, api.AdminListCourseLandingPagesParams{
		PerPage: api.NewOptInt32(100),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, int32(totalLandingPages), page.Total)

	// The embedded course supplies courseSlug (go-ru -> course "go").
	goRu := landingBySlug(t, h, "go-ru")
	assert.Equal(t, "go", goRu.CourseSlug)
	assert.Equal(t, int32(courseGoID), goRu.CourseId)
	assert.Equal(t, api.LandingPageStatePublished, goRu.State.Value)
	assert.True(t, goRu.Main.Value)
	assert.Equal(t, "Курс Go", goRu.Header)
}

func TestAdminGetCourseLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	got, err := h.Client.AdminGetCourseLandingPage(ctx, api.AdminGetCourseLandingPageParams{ID: landingGoRuID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "go-ru", got.Slug)
	assert.Equal(t, "go", got.CourseSlug)
	// Outcomes image is deferred until blob assets are migrated.
	assert.True(t, got.OutcomesImage.Null)
}

func TestAdminGetCourseLandingPageNotFound(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminGetCourseLandingPage(ctx, api.AdminGetCourseLandingPageParams{ID: 999999999})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminCreateCourseLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	in := fullLandingInput(courseRubyID)
	in.Slug = api.NewNilString("ruby-extra")
	in.Name = api.NewNilString("Ruby Extra")
	in.Main = api.NewNilBool(false)
	in.Listed = api.NewNilBool(true)
	in.State = api.NewNilLandingPageState(api.LandingPageStateDraft)
	in.Header = api.NewNilString("Extra header")
	in.MetaTitle = api.NewNilString("Meta")
	in.MetaDescription = api.NewNilString("Meta desc")
	in.Description = api.NewNilString("Desc")

	created, err := h.Client.AdminCreateCourseLandingPage(ctx, &in)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "ruby-extra", created.Slug)
	assert.Equal(t, "ruby", created.CourseSlug) // embedded course reloaded
	assert.Equal(t, api.LandingPageStateDraft, created.State.Value)

	total, err := h.DB.LandingPage.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalLandingPages+1, total)
}

func TestAdminUpdateCourseLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	target := landingBySlug(t, h, "go-ru")

	in := fullLandingInput(target.CourseId)
	in.Slug = api.NewNilString(target.Slug)
	in.Name = api.NewNilString("Go (updated)")
	in.Header = api.NewNilString("Новый заголовок")
	in.Order = api.NilString{Null: true} // clears the column
	in.State = api.NewNilLandingPageState(api.LandingPageStateArchived)
	in.Main = api.NewNilBool(false)

	updated, err := h.Client.AdminUpdateCourseLandingPage(ctx, &in, api.AdminUpdateCourseLandingPageParams{ID: target.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Новый заголовок", updated.Header)
	assert.Equal(t, api.LandingPageStateArchived, updated.State.Value)
	assert.True(t, updated.Order.Null)

	row, err := h.DB.LandingPage.Get(ctx, int(target.ID))
	require.NoError(t, err)
	assert.Nil(t, row.Order, "null order should clear the column")
}

func TestAdminDeleteCourseLandingPage(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	// Create a childless landing page, then delete it — avoids FK conflicts with
	// fixture landing pages that own qna items / category items.
	in := fullLandingInput(courseRubyID)
	in.Slug = api.NewNilString("to-delete")
	in.Name = api.NewNilString("Temp")
	created, err := h.Client.AdminCreateCourseLandingPage(ctx, &in)
	require.NoError(t, err)

	err = h.Client.AdminDeleteCourseLandingPage(ctx, api.AdminDeleteCourseLandingPageParams{ID: created.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	total, err := h.DB.LandingPage.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, totalLandingPages, total)
}
