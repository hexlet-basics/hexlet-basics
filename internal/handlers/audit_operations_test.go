package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/internal/api"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/store"
	"hexletbasics/internal/testsupport"
)

// A delete keeps the old values of the row it removed, so an investigator can
// see what a deleted post said.
func TestDeletingABlogPostRecordsItsOldValues(t *testing.T) {
	h := testsupport.NewHarness(t)
	before := h.DB.BlogPost.GetX(t.Context(), fixtures.ID(fixtures.BlogPostHelloWorld))

	require.NoError(t, h.Client.AdminDeleteBlogPost(context.Background(), api.AdminDeleteBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)}))

	rows := records(t, h, "AdminDeleteBlogPost")
	require.Len(t, rows, 1)
	name, ok := fieldOf(changesOf(t, rows[0]), "name")
	require.True(t, ok, "the deleted post's fields are in the record")
	assert.JSONEq(t, `"`+*before.Name+`"`, string(name.Old))
	assert.Empty(t, name.New)
}

// A job enqueued by an operation commits and rolls back with it: when the
// operation fails after enqueueing (here its record cannot be written), the job
// must not exist, or the worker would act on a change that never happened.
func TestAJobEnqueuedByAnOperationThatFailsIsNotKept(t *testing.T) {
	h := testsupport.NewHarness(t)

	breakWrites(h, ent.TypeAuditRecord)
	_, err := h.Client.AdminSuggestBlogPostRelatedCourses(context.Background(),
		api.AdminSuggestBlogPostRelatedCoursesParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)})
	require.Error(t, err)

	testsupport.RequireNotEnqueued(t, h.Services, jobs.SuggestRelatedCoursesArgs{})
}

func TestAJobEnqueuedBySucceedingOperationIsKeptWithItsRecord(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminSuggestBlogPostRelatedCourses(context.Background(),
		api.AdminSuggestBlogPostRelatedCoursesParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)})
	require.NoError(t, err)

	assert.Equal(t, jobs.SuggestRelatedCoursesArgs{BlogPostID: fixtures.ID(fixtures.BlogPostHelloWorld)},
		testsupport.RequireEnqueued(t, h.Services, jobs.SuggestRelatedCoursesArgs{}))
	require.Len(t, records(t, h, "AdminSuggestBlogPostRelatedCourses"), 1)
}

// The operations behind the administrators-only marker are recorded by what
// they do, with no per-operation code: a user created through the management
// API leaves a record naming the entity.
func TestAnOperationMarkedAdministratorsOnlyIsRecordedEndToEnd(t *testing.T) {
	h := testsupport.NewHarness(t)

	res, err := h.Client.AdminCreateUser(context.Background(), &api.UserInput{
		Email:     "audited-user@example.com",
		FirstName: api.NewNilString("Audited"),
		LastName:  api.NilString{Null: true},
		Admin:     api.NewNilBool(false),
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	rows := records(t, h, "AdminCreateUser")
	require.Len(t, rows, 1)
	assert.Equal(t, "succeeded", rows[0].Outcome)
	assert.Equal(t, "User", *rows[0].EntityType)
}

// The upload writes to blob storage, and a connection held in a transaction
// across that call is what the standalone declaration avoids: no transaction is
// open while the handler runs, and the record is written after it.
func TestAStandaloneOperationRunsItsHandlerWithoutATransactionOpen(t *testing.T) {
	h := testsupport.NewHarness(t)
	tracked := &trackedTransactor{inner: h.Transactor}
	audited := handlers.NewAuditMiddleware(tracked, testsupport.NewTranslator(t),
		testsupport.NewAPIErrorHandler(t, testsupport.NewTranslator(t)))

	var openDuringHandler int32 = -1
	probe := func(middleware.Request) (middleware.Response, error) {
		openDuringHandler = tracked.open.Load()
		return middleware.Response{}, nil
	}

	_, err := audited(middleware.Request{Context: t.Context(), OperationName: "AdminUploadAttachment"}, probe)
	require.NoError(t, err)
	assert.Zero(t, openDuringHandler, "no transaction around the external call")
	require.Len(t, records(t, h, "AdminUploadAttachment"), 1, "the record is still written")

	_, err = audited(middleware.Request{Context: t.Context(), OperationName: "AdminCreateBlogPost"}, probe)
	require.NoError(t, err)
	assert.EqualValues(t, 1, openDuringHandler, "an ordinary write keeps its one transaction")
}

// An upload still leaves its record through the real handler, naming the
// attachment row the standalone write created.
func TestAnUploadIsRecordedWithItsAttachment(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminUploadAttachment(t.Context(), &api.AttachmentUploadFormMultipart{
		File: ht.MultipartFile{Name: "cover.png", File: bytes.NewReader(tinyPNG), Size: int64(len(tinyPNG))},
	})
	require.NoError(t, err)

	rows := records(t, h, "AdminUploadAttachment")
	require.Len(t, rows, 1)
	assert.Equal(t, "succeeded", rows[0].Outcome)
	assert.Equal(t, "Attachment", *rows[0].EntityType)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
}

// A validation answer is a response, not an error, and nothing was changed: it
// is recorded as failed, not as a change that succeeded.
func TestAValidationAnswerIsRecordedAsFailed(t *testing.T) {
	h := testsupport.NewHarness(t)
	audited := handlers.NewAuditMiddleware(h.Transactor, testsupport.NewTranslator(t),
		testsupport.NewAPIErrorHandler(t, testsupport.NewTranslator(t)))

	resp, err := audited(middleware.Request{Context: t.Context(), OperationName: "AdminCreateBlogPost"},
		func(middleware.Request) (middleware.Response, error) {
			return middleware.Response{Type: &api.ValidationError{}}, nil
		})

	require.NoError(t, err, "the caller still gets the validation answer")
	assert.IsType(t, &api.ValidationError{}, resp.Type)
	rows := records(t, h, "AdminCreateBlogPost")
	require.Len(t, rows, 1)
	assert.Equal(t, "failed", rows[0].Outcome)
}

// Reads, public operations and refusals of a signed-out caller leave nothing.
func TestAnUnauthenticatedCallLeavesNoRecord(t *testing.T) {
	h := testsupport.NewAnonymousHarness(t)

	_, err := h.Client.AdminCreateBlogPost(context.Background(), permissionsPostInput("anon"))
	require.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, h.LastStatus())

	total, err := h.DB.AuditRecord.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, total)
}

// trackedTransactor counts the transactions open at a moment, so a test can ask
// whether the handler ran inside one.
type trackedTransactor struct {
	inner store.Transactor
	open  atomic.Int32
}

func (t *trackedTransactor) WithinTx(ctx context.Context, fn func(*sql.Tx, *ent.Client) error) error {
	t.open.Add(1)
	defer t.open.Add(-1)
	return t.inner.WithinTx(ctx, fn)
}
