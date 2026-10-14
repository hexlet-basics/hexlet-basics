package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/activestorageattachment"
	"hexletbasics/ent/blogpost"
	"hexletbasics/ent/blogpostlike"
	"hexletbasics/ent/blogpostrelatedcourseitem"
	"hexletbasics/ent/categoryqnaitem"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/coursemoduletranslation"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpageqnaitem"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/ent/staffrole"
	"hexletbasics/ent/staffrolepermission"
	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// What an operation answers when the database fails under it. The harness has
// one transaction and no way to take the database down mid-request, so the
// storage is broken inside that transaction instead: a table or column renamed
// away makes the reads that touch it fail, a trigger that raises makes the writes
// to a table fail. DDL is transactional, so the rollback at cleanup undoes it.
//
// A failure aborts the transaction (Postgres 25P02), so a test breaks one thing,
// makes one call and asserts it was answered as a server error — never as a
// not-found or a conflict, which are answers about the request, not the storage.

// breakReads makes every read of the table fail.
func breakReads(t *testing.T, h *testsupport.Harness, table string) {
	t.Helper()
	exec(t, h, fmt.Sprintf(`ALTER TABLE %q RENAME TO %q`, table, table+"_gone"))
}

// breakColumn makes the reads that select the column fail while a count of the
// table's rows still works.
func breakColumn(t *testing.T, h *testsupport.Harness, table, column string) {
	t.Helper()
	exec(t, h, fmt.Sprintf(`ALTER TABLE %q RENAME COLUMN %q TO %q`, table, column, column+"_gone"))
}

// breakWrites makes every INSERT, UPDATE or DELETE (event) on the table fail,
// while its reads keep working.
func breakWrites(t *testing.T, h *testsupport.Harness, table, event string) {
	t.Helper()
	exec(t, h, `CREATE OR REPLACE FUNCTION storage_down() RETURNS trigger LANGUAGE plpgsql
		AS $$ BEGIN RAISE EXCEPTION 'storage is down'; END $$`)
	exec(t, h, fmt.Sprintf(
		`CREATE TRIGGER storage_down BEFORE %s ON %q FOR EACH ROW EXECUTE FUNCTION storage_down()`,
		event, table))
}

func exec(t *testing.T, h *testsupport.Harness, statement string) {
	t.Helper()
	_, err := h.Tx.ExecContext(t.Context(), statement)
	require.NoError(t, err)
}

// failure is one broken storage and the call that must answer it as a server
// error.
type failure struct {
	name string
	// storage breaks what the call depends on.
	storage func(*testing.T, *testsupport.Harness)
	call    func(context.Context, *testsupport.Harness) error
}

func reads(table string) func(*testing.T, *testsupport.Harness) {
	return func(t *testing.T, h *testsupport.Harness) { breakReads(t, h, table) }
}

func writes(table, event string) func(*testing.T, *testsupport.Harness) {
	return func(t *testing.T, h *testsupport.Harness) { breakWrites(t, h, table, event) }
}

func column(table, name string) func(*testing.T, *testsupport.Harness) {
	return func(t *testing.T, h *testsupport.Harness) { breakColumn(t, h, table, name) }
}

func answerServerErrors(t *testing.T, cases []failure) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			c.storage(t, h)
			err := c.call(t.Context(), h)
			require.Error(t, err)
			assert.Equal(t, http.StatusInternalServerError, h.LastStatus())
		})
	}
}

func TestQnaItemsAnswerAStorageFailureAsAServerError(t *testing.T) {
	input := &api.QnaItemInput{Question: "Q?", Answer: "A."}
	category := api.AdminListCategoryQnaItemsParams{CategoryId: qnaCategoryID}
	landing := api.AdminListLandingPageQnaItemsParams{LandingPageId: qnaLandingPageID}
	deleteCategoryItem := func(ctx context.Context, h *testsupport.Harness) error {
		return h.Client.AdminDeleteCategoryQnaItem(ctx,
			api.AdminDeleteCategoryQnaItemParams{CategoryId: qnaCategoryID, ID: 7001})
	}
	deleteLandingItem := func(ctx context.Context, h *testsupport.Harness) error {
		return h.Client.AdminDeleteLandingPageQnaItem(ctx,
			api.AdminDeleteLandingPageQnaItemParams{LandingPageId: qnaLandingPageID, ID: 8001})
	}

	answerServerErrors(t, []failure{
		{"list category items", reads(categoryqnaitem.Table), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminListCategoryQnaItems(ctx, category)
			return err
		}},
		{"create a category item", writes(categoryqnaitem.Table, "INSERT"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminCreateCategoryQnaItem(ctx, input,
				api.AdminCreateCategoryQnaItemParams{CategoryId: qnaCategoryID})
			return err
		}},
		{"update a category item", writes(categoryqnaitem.Table, "UPDATE"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateCategoryQnaItem(ctx, input,
				api.AdminUpdateCategoryQnaItemParams{CategoryId: qnaCategoryID, ID: 7001})
			return err
		}},
		{"find a category item to delete", reads(categoryqnaitem.Table), deleteCategoryItem},
		{"delete a category item", writes(categoryqnaitem.Table, "DELETE"), deleteCategoryItem},
		{"list landing page items", reads(landingpageqnaitem.Table), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminListLandingPageQnaItems(ctx, landing)
			return err
		}},
		{"create a landing page item", writes(landingpageqnaitem.Table, "INSERT"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminCreateLandingPageQnaItem(ctx, input,
				api.AdminCreateLandingPageQnaItemParams{LandingPageId: qnaLandingPageID})
			return err
		}},
		{"update a landing page item", writes(landingpageqnaitem.Table, "UPDATE"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateLandingPageQnaItem(ctx, input,
				api.AdminUpdateLandingPageQnaItemParams{LandingPageId: qnaLandingPageID, ID: 8001})
			return err
		}},
		{"find a landing page item to delete", reads(landingpageqnaitem.Table), deleteLandingItem},
		{"delete a landing page item", writes(landingpageqnaitem.Table, "DELETE"), deleteLandingItem},
	})
}

func TestStaffRolesAnswerAStorageFailureAsAServerError(t *testing.T) {
	const roleID = roleEditorID
	permissions := &api.RolePermissionsInput{Permissions: []api.RolePermissionInput{}}

	answerServerErrors(t, []failure{
		{"update a role", writes(staffrole.Table, "UPDATE"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateRole(ctx, &api.RoleInput{Name: "Renamed"},
				api.AdminUpdateRoleParams{ID: roleID})
			return err
		}},
		{"delete a role", writes(staffrole.Table, "DELETE"), func(ctx context.Context, h *testsupport.Harness) error {
			return h.Client.AdminDeleteRole(ctx, api.AdminDeleteRoleParams{ID: roleID})
		}},
		{"look up the role whose permissions are set", reads(staffrole.Table), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateRolePermissions(ctx, permissions,
				api.AdminUpdateRolePermissionsParams{RoleId: roleID})
			return err
		}},
		{"store the role's permissions", writes(staffrolepermission.Table, "INSERT"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminUpdateRolePermissions(ctx, &api.RolePermissionsInput{
				Permissions: []api.RolePermissionInput{{Resource: "banners", CanIndex: true}},
			}, api.AdminUpdateRolePermissionsParams{RoleId: roleID})
			return err
		}},
	})
}

func TestBlogPostsAnswerAStorageFailureAsAServerError(t *testing.T) {
	list := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.AdminListBlogPosts(ctx, api.AdminListBlogPostsParams{})
		return err
	}
	get := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.AdminGetBlogPost(ctx, api.AdminGetBlogPostParams{ID: 6001})
		return err
	}

	answerServerErrors(t, []failure{
		{"count the posts of a page", reads(blogpost.Table), list},
		{"read the posts of a page", column(blogpost.Table, "name"), list},
		{"read the covers of the listed posts", reads(activestorageattachment.Table), list},
		{"read the likes of the listed posts", reads(blogpostlike.Table), list},
		{"read the related courses of the listed posts", reads(blogpostrelatedcourseitem.Table), list},
		{"read the covers of one post", reads(activestorageattachment.Table), get},
		{"create a post", writes(blogpost.Table, "INSERT"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminCreateBlogPost(ctx, &api.BlogPostInput{
				Name:     api.NewNilString("N"),
				Slug:     api.NewNilString("n"),
				State:    api.NewNilBlogPostState(api.BlogPostStateDraft),
				RichBody: "<p>n</p>",
			})
			return err
		}},
		{"delete the post", writes(blogpost.Table, "DELETE"), func(ctx context.Context, h *testsupport.Harness) error {
			return h.Client.AdminDeleteBlogPost(ctx, api.AdminDeleteBlogPostParams{ID: 6001})
		}},
		{"queue the suggestion of related courses", writes("river_job", "INSERT"), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminSuggestBlogPostRelatedCourses(ctx, api.AdminSuggestBlogPostRelatedCoursesParams{ID: 6001})
			return err
		}},
		{"find the post to delete", reads(blogpost.Table), func(ctx context.Context, h *testsupport.Harness) error {
			return h.Client.AdminDeleteBlogPost(ctx, api.AdminDeleteBlogPostParams{ID: 6001})
		}},
		{"find the post to suggest courses for", reads(blogpost.Table), func(ctx context.Context, h *testsupport.Harness) error {
			_, err := h.Client.AdminSuggestBlogPostRelatedCourses(ctx, api.AdminSuggestBlogPostRelatedCoursesParams{ID: 6001})
			return err
		}},
	})
}

func TestCourseAnswersAStorageFailureAsAServerError(t *testing.T) {
	get := func(ctx context.Context, h *testsupport.Harness) error {
		testsupport.SpeakTo(h, "ru")
		_, err := h.Client.GetCourse(ctx, api.GetCourseParams{Slug: jsLandingSlug})
		return err
	}

	answerServerErrors(t, []failure{
		{"read the course", reads(course.Table), get},
		{"read the landing page", reads(courselandingpage.Table), get},
		{"read the lessons", reads(courselessontranslation.Table), get},
		{"read the modules", reads(coursemoduletranslation.Table), get},
		{"read the landing page's questions", reads(landingpageqnaitem.Table), get},
		{"read the learner's place in the course", reads(lessonprogress.Table), get},
		{"read the learner's enrollment", reads(enrollment.Table), get},
	})
}

func TestLessonAnswersAStorageFailureAsAServerError(t *testing.T) {
	get := func(ctx context.Context, h *testsupport.Harness) error {
		_, err := h.Client.GetCourseLesson(ctx, api.GetCourseLessonParams{
			CourseSlug: jsCourseSlug,
			Slug:       secondLessonSlug,
		})
		return err
	}

	answerServerErrors(t, []failure{
		{"read the course", reads(course.Table), get},
		{"read the lesson", reads(courselesson.Table), get},
		{"read the lesson's version", reads(courselessonversion.Table), get},
		{"read the lesson's translation", reads(courselessontranslation.Table), get},
		{"read the learner's place in the course", reads(lessonprogress.Table), get},
		{"read the landing page", reads(courselandingpage.Table), get},
	})
}
