package handlers_test

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/testsupport"
)

// A state is not an update (ADR-0019, state transitions addendum), for banners,
// reviews and course landing pages as for blog posts: the Publication State moves
// only through publish, unpublish and archive. The three resources share one
// table of legal moves (internal/transitions), so one set of scenarios runs over
// all of them, and each carries what is its own: how the public sees it.

// stateSubject is one resource under test, named by the operations it has.
type stateSubject struct {
	name string
	// operation is the resource's part of the operation names (AdminPublish<operation>).
	operation string
	// id is the fixture row the scenarios move.
	id int
	// updaterRole may update the resource and not publish it; publisherRole the reverse.
	updaterRole, publisherRole fixtures.Label
	// nullable says the legacy column may hold no state (counts as a draft).
	nullable bool

	setState   func(h *testsupport.Harness, state string)
	clearState func(h *testsupport.Harness)
	stateOf    func(t *testing.T, h *testsupport.Harness) *string
	// move calls publish, unpublish or archive on the row and returns the state the answer shows.
	move func(h *testsupport.Harness, verb string, id int) (string, error)
	// update calls the resource's update operation with an input that changes nothing that matters.
	update func(h *testsupport.Harness) error
	// visible says whether the public sees the row.
	visible func(t *testing.T, h *testsupport.Harness) bool
}

var (
	bannerSubject = stateSubject{
		name: "banner", operation: "Banner", id: fixtures.ID(fixtures.BannerEnPublished),
		updaterRole: fixtures.StaffRoleBannerUpdater, publisherRole: fixtures.StaffRoleBannerPublisher,
		setState: func(h *testsupport.Harness, state string) {
			h.DB.Banner.UpdateOneID(fixtures.ID(fixtures.BannerEnPublished)).SetState(state).ExecX(context.Background())
		},
		stateOf: func(t *testing.T, h *testsupport.Harness) *string {
			state := h.DB.Banner.GetX(t.Context(), fixtures.ID(fixtures.BannerEnPublished)).State
			return &state
		},
		move: func(h *testsupport.Harness, verb string, id int) (string, error) {
			var banner *api.Banner
			var err error
			switch verb {
			case "publish":
				banner, err = h.Client.AdminPublishBanner(context.Background(), api.AdminPublishBannerParams{ID: int32(id)})
			case "unpublish":
				banner, err = h.Client.AdminUnpublishBanner(context.Background(), api.AdminUnpublishBannerParams{ID: int32(id)})
			case "archive":
				banner, err = h.Client.AdminArchiveBanner(context.Background(), api.AdminArchiveBannerParams{ID: int32(id)})
			}
			if err != nil {
				return "", err
			}
			return string(banner.State), nil
		},
		update: func(h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateBanner(context.Background(), &api.BannerInput{
				Background: api.BannerBackgroundBlue,
				Locale:     api.BannerLocaleEn,
				Body:       "Edited",
				URL:        api.NilString{Null: true},
				StartsAt:   api.NilDateTime{Null: true},
				FinishesAt: api.NilDateTime{Null: true},
			}, api.AdminUpdateBannerParams{ID: int32(fixtures.ID(fixtures.BannerEnPublished))})
			return err
		},
		// The page shows the live banner of the request locale: the published, in-window one.
		visible: func(t *testing.T, h *testsupport.Harness) bool {
			testsupport.SpeakTo(h, "en")
			view, err := h.Client.GetActiveBanner(t.Context())
			require.NoError(t, err)
			banner, ok := view.Banner.Get()
			return ok && banner.Body == "Free programming courses — start learning today"
		},
	}

	reviewSubject = stateSubject{
		name: "review", operation: "Review", id: fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished),
		updaterRole: fixtures.StaffRoleReviewUpdater, publisherRole: fixtures.StaffRoleReviewPublisher,
		nullable: true,
		setState: func(h *testsupport.Harness, state string) {
			h.DB.Review.UpdateOneID(fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished)).SetState(state).ExecX(context.Background())
		},
		clearState: func(h *testsupport.Harness) {
			h.DB.Review.UpdateOneID(fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished)).ClearState().ExecX(context.Background())
		},
		stateOf: func(t *testing.T, h *testsupport.Harness) *string {
			return h.DB.Review.GetX(t.Context(), fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished)).State
		},
		move: func(h *testsupport.Harness, verb string, id int) (string, error) {
			var review *api.Review
			var err error
			switch verb {
			case "publish":
				review, err = h.Client.AdminPublishReview(context.Background(), api.AdminPublishReviewParams{ID: int32(id)})
			case "unpublish":
				review, err = h.Client.AdminUnpublishReview(context.Background(), api.AdminUnpublishReviewParams{ID: int32(id)})
			case "archive":
				review, err = h.Client.AdminArchiveReview(context.Background(), api.AdminArchiveReviewParams{ID: int32(id)})
			}
			if err != nil {
				return "", err
			}
			return string(review.State.Value), nil
		},
		update: func(h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateReview(context.Background(), &api.ReviewInput{
				Pinned:    api.NilBool{Null: true},
				CourseId:  api.NilInt32{Null: true},
				UserId:    api.NilInt32{Null: true},
				Body:      api.NewNilString("Edited"),
				FirstName: api.NewNilString("Ivan"),
				LastName:  api.NilString{Null: true},
			}, api.AdminUpdateReviewParams{ID: int32(fixtures.ID(fixtures.ReviewAdminLearnerRubyPublished))})
			return err
		},
		// The public reviews page lists the published reviews of the request locale.
		visible: func(t *testing.T, h *testsupport.Harness) bool {
			testsupport.SpeakTo(h, "ru")
			page, err := h.Client.ListPublicReviews(t.Context(), api.ListPublicReviewsParams{})
			require.NoError(t, err)
			return slices.Contains(reviewNames(page.Items), "Ivan Petrov")
		},
	}

	landingSubject = stateSubject{
		name: "landing page", operation: "CourseLandingPage", id: fixtures.ID(fixtures.LandingPageGoRu),
		updaterRole: fixtures.StaffRoleLandingPageUpdater, publisherRole: fixtures.StaffRoleLandingPagePublisher,
		nullable: true,
		setState: func(h *testsupport.Harness, state string) {
			h.DB.CourseLandingPage.UpdateOneID(fixtures.ID(fixtures.LandingPageGoRu)).SetState(state).ExecX(context.Background())
		},
		clearState: func(h *testsupport.Harness) {
			h.DB.CourseLandingPage.UpdateOneID(fixtures.ID(fixtures.LandingPageGoRu)).ClearState().ExecX(context.Background())
		},
		stateOf: func(t *testing.T, h *testsupport.Harness) *string {
			return h.DB.CourseLandingPage.GetX(t.Context(), fixtures.ID(fixtures.LandingPageGoRu)).State
		},
		move: func(h *testsupport.Harness, verb string, id int) (string, error) {
			var page *api.CourseLandingPage
			var err error
			switch verb {
			case "publish":
				page, err = h.Client.AdminPublishCourseLandingPage(context.Background(), api.AdminPublishCourseLandingPageParams{ID: int32(id)})
			case "unpublish":
				page, err = h.Client.AdminUnpublishCourseLandingPage(context.Background(), api.AdminUnpublishCourseLandingPageParams{ID: int32(id)})
			case "archive":
				page, err = h.Client.AdminArchiveCourseLandingPage(context.Background(), api.AdminArchiveCourseLandingPageParams{ID: int32(id)})
			}
			if err != nil {
				return "", err
			}
			return string(page.State.Value), nil
		},
		update: func(h *testsupport.Harness) error {
			in := fullLandingInput(int32(courseGoID))
			in.Slug = api.NewNilString("go-ru")
			in.Name = api.NewNilString("Edited")
			_, err := h.Client.AdminUpdateCourseLandingPage(context.Background(), &in, api.AdminUpdateCourseLandingPageParams{ID: int32(fixtures.ID(fixtures.LandingPageGoRu))})
			return err
		},
		// A published page is listed on the category page through the category items it
		// already has; nothing is copied when it is published.
		visible: func(t *testing.T, h *testsupport.Harness) bool {
			testsupport.SpeakTo(h, "ru")
			return slices.Contains(categoryLandingSlugs(getCategoryView(t, h, "programming-ru")), "go-ru")
		},
	}

	stateSubjects = []stateSubject{bannerSubject, reviewSubject, landingSubject}
)

// legalMoves is the table as the issue states it; every other pair of (operation,
// current state) is refused, repeating the current state included.
var legalMoves = map[string][]string{
	"publish":   {"draft", "archived"},
	"unpublish": {"published"},
	"archive":   {"draft", "published"},
}

var movesTo = map[string]string{"publish": "published", "unpublish": "draft", "archive": "archived"}

func capitalized(verb string) string {
	return map[string]string{"publish": "Publish", "unpublish": "Unpublish", "archive": "Archive"}[verb]
}

func TestEveryLegalMoveOfBannersReviewsAndLandingPagesWorksAndEveryOtherIs409(t *testing.T) {
	for _, subject := range stateSubjects {
		for verb, legal := range legalMoves {
			for _, from := range []string{"draft", "published", "archived"} {
				t.Run(subject.name+" "+verb+" from "+from, func(t *testing.T) {
					h := testsupport.NewHarness(t)
					subject.setState(h, from)

					shown, err := subject.move(h, verb, subject.id)

					if slices.Contains(legal, from) {
						require.NoError(t, err)
						assert.Equal(t, movesTo[verb], shown, "the answer shows the new state")
						assert.Equal(t, movesTo[verb], *subject.stateOf(t, h))
						return
					}
					require.Error(t, err)
					assert.Equal(t, http.StatusConflict, h.LastStatus())
					assert.Equal(t, from, *subject.stateOf(t, h), "a refused move changes nothing")
					assert.NotEmpty(t, problemOf(t, err).Detail.Value, "the refusal explains itself")
				})
			}
		}
	}
}

func TestARefusedMoveExplainsItselfInTheCallersLanguage(t *testing.T) {
	for _, subject := range stateSubjects {
		for locale, want := range map[string]string{
			"ru": "«опубликовано»",
			"en": "published to published",
			"es": "«publicado»",
		} {
			t.Run(subject.name+" "+locale, func(t *testing.T) {
				h := testsupport.NewHarness(t)
				testsupport.SpeakTo(h, locale)
				subject.setState(h, "published")

				_, err := subject.move(h, "publish", subject.id)

				require.Error(t, err)
				assert.Contains(t, problemOf(t, err).Detail.Value, want)
			})
		}
	}
}

func TestARowWithoutAStateCountsAsADraft(t *testing.T) {
	for _, subject := range stateSubjects {
		if !subject.nullable {
			continue
		}
		t.Run(subject.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			subject.clearState(h)

			shown, err := subject.move(h, "publish", subject.id)

			require.NoError(t, err)
			assert.Equal(t, "published", shown)
		})
	}
}

// Each move is a write operation with a Permission of its own, named with the
// GLOSSARY term of the resource.
func TestEachMoveDeclaresItsOwnWritePermission(t *testing.T) {
	for operation, resource := range map[string]string{"Banner": "banners", "Review": "reviews", "CourseLandingPage": "course_landing_pages"} {
		for verb := range legalMoves {
			got, ok := permissions.For("Admin" + capitalized(verb) + operation)
			require.True(t, ok)
			assert.Equal(t, permissions.Requirement{Permission: resource + "." + verb, Kind: permissions.Write}, got)
		}
	}
}

func TestAMoveOfAMissingRowIs404(t *testing.T) {
	for _, subject := range stateSubjects {
		t.Run(subject.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)

			_, err := subject.move(h, "publish", 999999)

			require.Error(t, err)
			assert.Equal(t, http.StatusNotFound, h.LastStatus())
		})
	}
}

// Publishing is what the public sees follow, and unpublishing or archiving what
// takes it away: nothing else of the old update is lost, because the public
// queries read the state column and nothing is cached or copied.
func TestWhatThePublicSeesFollowsTheState(t *testing.T) {
	for _, subject := range stateSubjects {
		t.Run(subject.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			require.True(t, subject.visible(t, h), "published in the fixtures")

			_, err := subject.move(h, "unpublish", subject.id)
			require.NoError(t, err)
			assert.False(t, subject.visible(t, h), "a draft is not shown")

			_, err = subject.move(h, "publish", subject.id)
			require.NoError(t, err)
			assert.True(t, subject.visible(t, h), "publishing shows it again")

			_, err = subject.move(h, "archive", subject.id)
			require.NoError(t, err)
			assert.False(t, subject.visible(t, h), "an archived one is withdrawn")
		})
	}
}

// A role that may edit may not publish, and the reverse: the state is not an
// update. Unpublishing and archiving are granted apart from publishing too.
func TestUpdatePublishAndTheOtherMovesArePermissionsOfTheirOwn(t *testing.T) {
	for _, subject := range stateSubjects {
		t.Run(subject.name+" update only", func(t *testing.T) {
			h := testsupport.NewHarnessAs(t, editorLabel)
			actAs(h, subject.updaterRole)

			require.NoError(t, subject.update(h))

			subject.setState(h, "draft")
			_, err := subject.move(h, "publish", subject.id)
			require.Error(t, err)
			assert.Equal(t, http.StatusForbidden, h.LastStatus())
			assert.Equal(t, "draft", *subject.stateOf(t, h))
		})
		t.Run(subject.name+" publish only", func(t *testing.T) {
			h := testsupport.NewHarnessAs(t, editorLabel)
			actAs(h, subject.publisherRole)

			subject.setState(h, "draft")
			_, err := subject.move(h, "publish", subject.id)
			require.NoError(t, err)

			require.Error(t, subject.update(h))
			assert.Equal(t, http.StatusForbidden, h.LastStatus())

			_, err = subject.move(h, "unpublish", subject.id)
			require.Error(t, err, "publishing does not include unpublishing")
			assert.Equal(t, http.StatusForbidden, h.LastStatus())
			_, err = subject.move(h, "archive", subject.id)
			require.Error(t, err, "publishing does not include archiving")
			assert.Equal(t, http.StatusForbidden, h.LastStatus())
			assert.Equal(t, "published", *subject.stateOf(t, h))
		})
	}
}

func TestEachMoveLeavesAnAuditRecordWithTheOldAndNewState(t *testing.T) {
	for _, subject := range stateSubjects {
		for verb, legal := range legalMoves {
			from := legal[0]
			t.Run(subject.name+" "+verb, func(t *testing.T) {
				h := testsupport.NewHarness(t)
				subject.setState(h, from)

				_, err := subject.move(h, verb, subject.id)
				require.NoError(t, err)

				rows := records(t, h, "Admin"+capitalized(verb)+subject.operation)
				require.Len(t, rows, 1)
				assert.Equal(t, strconv.Itoa(subject.id), *rows[0].EntityID)
				state, ok := fieldOf(changesOf(t, rows[0]), "state")
				require.True(t, ok)
				assert.JSONEq(t, `"`+from+`"`, string(state.Old))
				assert.JSONEq(t, `"`+movesTo[verb]+`"`, string(state.New))
			})
		}
	}
}
