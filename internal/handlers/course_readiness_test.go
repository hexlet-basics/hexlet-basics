package handlers_test

import (
	"context"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/permissioncatalog"
	"hexletbasics/internal/testsupport"
)

// A state is not an update (ADR-0020, state transitions addendum), for a
// Course's Readiness as for a Publication State: it moves only through the three
// mark-as operations. Readiness is no workflow, so every change between two
// distinct values is legal and only repeating the current value is refused.

var readinessCourseID = fixture.ID(fixture.CourseGo)

var readinessValues = []string{"draft", "in_development", "completed"}

// markCourse calls the mark-as operation of the target readiness and returns the
// readiness the answer shows.
func markCourse(h *testsupport.Harness, to string) (string, error) {
	ctx := context.Background()
	id := int32(readinessCourseID)
	var course *api.Course
	var err error
	switch to {
	case "draft":
		course, err = h.Client.AdminMarkCourseDraft(ctx, api.AdminMarkCourseDraftParams{ID: id})
	case "in_development":
		course, err = h.Client.AdminMarkCourseInDevelopment(ctx, api.AdminMarkCourseInDevelopmentParams{ID: id})
	case "completed":
		course, err = h.Client.AdminMarkCourseCompleted(ctx, api.AdminMarkCourseCompletedParams{ID: id})
	}
	if err != nil {
		return "", err
	}
	return string(course.Readiness.Value), nil
}

func setReadiness(h *testsupport.Harness, readiness string) {
	h.DB.Course.UpdateOneID(readinessCourseID).SetReadiness(readiness).ExecX(context.Background())
}

func readinessOf(t *testing.T, h *testsupport.Harness) *string {
	return h.DB.Course.GetX(t.Context(), readinessCourseID).Readiness
}

func TestEveryChangeOfReadinessBetweenDistinctValuesWorksAndRepeatingIs409(t *testing.T) {
	t.Parallel()
	for _, from := range readinessValues {
		for _, to := range readinessValues {
			t.Run(from+" to "+to, func(t *testing.T) {
				h := testsupport.NewHarness(t)
				setReadiness(h, from)

				shown, err := markCourse(h, to)

				if from != to {
					require.NoError(t, err)
					assert.Equal(t, to, shown, "the answer shows the new readiness")
					assert.Equal(t, to, *readinessOf(t, h))
					return
				}
				require.Error(t, err)
				assert.Equal(t, http.StatusConflict, h.LastStatus())
				assert.Equal(t, from, *readinessOf(t, h), "a refused move changes nothing")
				assert.NotEmpty(t, problemOf(t, err).Detail.Value, "the refusal explains itself")
			})
		}
	}
}

func TestARefusedReadinessChangeExplainsItselfInTheCallersLanguage(t *testing.T) {
	t.Parallel()
	for locale, want := range map[string]string{
		"ru": "«в разработке»",
		"en": "in development to in development",
		"es": "«en desarrollo»",
	} {
		t.Run(locale, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			testsupport.SpeakTo(h, locale)
			setReadiness(h, "in_development")

			_, err := markCourse(h, "in_development")

			require.Error(t, err)
			assert.Contains(t, problemOf(t, err).Detail.Value, want)
		})
	}
}

func TestACourseWithoutReadinessCountsAsADraft(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	h.DB.Course.UpdateOneID(readinessCourseID).ClearReadiness().ExecX(t.Context())

	_, err := markCourse(h, "draft")
	require.Error(t, err, "it already is a draft")
	assert.Equal(t, http.StatusConflict, h.LastStatus())

	shown, err := markCourse(h, "completed")
	require.NoError(t, err)
	assert.Equal(t, "completed", shown)
}

func TestEachReadinessChangeDeclaresItsOwnWritePermission(t *testing.T) {
	t.Parallel()
	for operation, permission := range map[string]string{
		"AdminMarkCourseDraft":         "courses.mark_draft",
		"AdminMarkCourseInDevelopment": "courses.mark_in_development",
		"AdminMarkCourseCompleted":     "courses.mark_completed",
	} {
		got, ok := permissioncatalog.For(operation)
		require.True(t, ok)
		assert.Equal(t, permissioncatalog.Requirement{Permission: permission, Kind: permissioncatalog.Write}, got)
	}
}

func TestAReadinessChangeOfAMissingCourseIs404(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminMarkCourseCompleted(t.Context(), api.AdminMarkCourseCompletedParams{ID: 999999})

	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

// What the readiness controls is read from the column by the public queries, so
// marking a course completed lists it and marking it back withdraws it, with
// nothing copied or cached by the operation.
func TestWhatThePublicSeesFollowsTheReadiness(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	listed := func() bool {
		testsupport.SpeakTo(h, "ru")
		return slices.Contains(categoryLandingSlugs(getCategoryView(t, h, "programming-ru")), "go-ru")
	}
	require.True(t, listed(), "completed in the fixtures")

	_, err := markCourse(h, "in_development")
	require.NoError(t, err)
	assert.False(t, listed(), "a course in development is off the public lists")

	_, err = markCourse(h, "draft")
	require.NoError(t, err)
	assert.False(t, listed(), "so is a draft")

	_, err = markCourse(h, "completed")
	require.NoError(t, err)
	assert.True(t, listed(), "completing it lists it again")
}

// A new course is a draft and the update never touches the readiness, so the
// readiness changes only through the operations.
func TestCreateMakesADraftAndUpdateKeepsTheReadiness(t *testing.T) {
	t.Parallel()
	h := testsupport.NewHarness(t)
	setReadiness(h, "in_development")

	_, err := h.Client.AdminUpdateCourse(t.Context(), &api.CourseInput{
		Slug:                     api.NewNilString("go"),
		LearnAs:                  api.NilCourseLearnAs{Null: true},
		HexletProgramLandingPage: api.NilString{Null: true},
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	}, api.AdminUpdateCourseParams{ID: int32(readinessCourseID)})
	require.NoError(t, err)
	assert.Equal(t, "in_development", *readinessOf(t, h))

	created, err := h.Client.AdminCreateCourse(t.Context(), &api.CourseInput{
		Slug:                     api.NewNilString("haskell"),
		LearnAs:                  api.NilCourseLearnAs{Null: true},
		HexletProgramLandingPage: api.NilString{Null: true},
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	})
	require.NoError(t, err)
	assert.Equal(t, "draft", *h.DB.Course.GetX(t.Context(), int(created.ID)).Readiness)
}

// A role that may edit a course may not change its readiness, and the reverse:
// the state is not an update. Each of the three changes is granted apart too.
func TestUpdateAndTheReadinessChangesArePermissionsOfTheirOwn(t *testing.T) {
	t.Parallel()
	input := &api.CourseInput{
		Slug:                     api.NewNilString("go"),
		LearnAs:                  api.NilCourseLearnAs{Null: true},
		HexletProgramLandingPage: api.NilString{Null: true},
		RepositoryUrl:            api.NilString{Null: true},
		CoverAttachmentId:        api.NilInt32{Null: true},
	}
	update := func(h *testsupport.Harness) error {
		_, err := h.Client.AdminUpdateCourse(context.Background(), input, api.AdminUpdateCourseParams{ID: int32(readinessCourseID)})
		return err
	}

	t.Run("update only", func(t *testing.T) {
		h := testsupport.NewHarnessAs(t, editorLabel)
		actAs(h, fixture.StaffRoleCourseUpdater)

		require.NoError(t, update(h))

		for _, to := range readinessValues {
			_, err := markCourse(h, to)
			require.Error(t, err, to)
			assert.Equal(t, http.StatusForbidden, h.LastStatus(), to)
		}
		assert.Equal(t, "completed", *readinessOf(t, h))
	})

	t.Run("readiness only", func(t *testing.T) {
		h := testsupport.NewHarnessAs(t, editorLabel)
		actAs(h, fixture.StaffRoleCourseReadiness)

		shown, err := markCourse(h, "draft")
		require.NoError(t, err)
		assert.Equal(t, "draft", shown)

		require.Error(t, update(h))
		assert.Equal(t, http.StatusForbidden, h.LastStatus())

		_, err = markCourse(h, "in_development")
		require.Error(t, err, "the role holds mark_draft and mark_completed, not mark_in_development")
		assert.Equal(t, http.StatusForbidden, h.LastStatus())
		assert.Equal(t, "draft", *readinessOf(t, h))
	})
}

func TestEachReadinessChangeLeavesAnAuditRecordWithTheOldAndNewValue(t *testing.T) {
	t.Parallel()
	for operation, to := range map[string]string{
		"AdminMarkCourseDraft":         "draft",
		"AdminMarkCourseInDevelopment": "in_development",
		"AdminMarkCourseCompleted":     "completed",
	} {
		t.Run(operation, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			from := "completed"
			if to == "completed" {
				from = "draft"
			}
			setReadiness(h, from)

			_, err := markCourse(h, to)
			require.NoError(t, err)

			rows := records(t, h, operation)
			require.Len(t, rows, 1)
			assert.Equal(t, "Course", *rows[0].EntityType)
			// The record names the column, which keeps the retired name `progress`.
			readiness, ok := fieldOf(changesOf(t, rows[0]), "progress")
			require.True(t, ok)
			assert.JSONEq(t, `"`+from+`"`, string(readiness.Old))
			assert.JSONEq(t, `"`+to+`"`, string(readiness.New))
		})
	}
}
