package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/ogen-go/ogen/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent/auditrecord"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/eventhandlers"
	"hexletbasics/internal/events"
	"hexletbasics/internal/handlers"
	"hexletbasics/internal/permissions"
	"hexletbasics/internal/testsupport"
)

// A staff read of personal data is recorded asynchronously (ADR-0019): the
// request publishes a fact, and the worker's consumer writes the Audit Record.
// These tests go through the operation, then deliver the fact the way the
// worker does, and read the record.

// processReads delivers the published facts to the consumer the worker runs.
func processReads(t *testing.T, h *testsupport.Harness) {
	t.Helper()
	testsupport.Deliver(t, h.Services, eventhandlers.StaffReadPersonalData(h.DB))
}

type readsOf struct {
	Count int      `json:"count"`
	IDs   []string `json:"ids"`
}

func TestListingUsersLeavesAReadRecordOnceTheFactIsProcessed(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
	testsupport.ArriveFrom(h, "203.0.113.20")
	testsupport.Identify(h, "req-read-1", "AdminUI/3.0")

	page, err := h.Client.AdminListUsers(context.Background(), api.AdminListUsersParams{PerPage: api.NewOptInt32(3)})
	require.NoError(t, err)
	require.Len(t, page.Items, 3)

	assert.Empty(t, records(t, h, "AdminListUsers"), "the record is written by the worker, not by the request")
	require.Len(t, testsupport.PublishedEvents[events.StaffReadPersonalData](t, h.Services), 1)

	processReads(t, h)

	rows := records(t, h, "AdminListUsers")
	require.Len(t, rows, 1)
	row := rows[0]
	assert.Equal(t, h.UserID, *row.ActorUserID)
	assert.Equal(t, "web", row.Channel)
	assert.Equal(t, "succeeded", row.Outcome)
	assert.Equal(t, "req-read-1", *row.RequestID)
	assert.Equal(t, "203.0.113.20", *row.Address)
	assert.Equal(t, "AdminUI/3.0", *row.UserAgent)
	assert.Equal(t, "ru", *row.Locale)
	assert.Equal(t, "User", *row.EntityType)
	assert.Nil(t, row.EntityID, "a list is not filed under one row")
	assert.Empty(t, row.Changes, "a read changed nothing")

	var reads readsOf
	require.NoError(t, json.Unmarshal(row.Reads, &reads))
	assert.Equal(t, 3, reads.Count)
	var want []string
	for _, user := range page.Items {
		want = append(want, strconv.Itoa(int(user.ID)))
	}
	assert.Equal(t, want, reads.IDs)
}

func TestReadingOneUserIsFiledUnderThatUser(t *testing.T) {
	h := testsupport.NewHarness(t)

	user, err := h.Client.AdminGetUser(context.Background(), api.AdminGetUserParams{ID: int32(h.UserID)})
	require.NoError(t, err)
	processReads(t, h)

	rows := records(t, h, "AdminGetUser")
	require.Len(t, rows, 1)
	assert.Equal(t, "User", *rows[0].EntityType)
	assert.Equal(t, strconv.Itoa(int(user.ID)), *rows[0].EntityID)
}

// A record holds identifiers and a count, never what the response said: not the
// email of a user the list showed.
func TestAReadRecordHoldsNoFieldValues(t *testing.T) {
	h := testsupport.NewHarness(t)

	page, err := h.Client.AdminListUsers(context.Background(), api.AdminListUsersParams{})
	require.NoError(t, err)
	require.NotEmpty(t, page.Items)
	processReads(t, h)

	row := records(t, h, "AdminListUsers")[0]
	for _, user := range page.Items {
		assert.NotContains(t, string(row.Reads), user.Email)
	}
	var keys map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(row.Reads, &keys))
	for key := range keys {
		assert.Contains(t, []string{"count", "ids"}, key)
	}
}

func TestSearchingLeadsMessagesAndProgressLeaveRecords(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminSearchUsers(ctx, api.AdminSearchUsersParams{Q: "a"})
	require.NoError(t, err)
	_, err = h.Client.AdminListLeads(ctx, api.AdminListLeadsParams{})
	require.NoError(t, err)
	_, err = h.Client.AdminListMessages(ctx, api.AdminListMessagesParams{})
	require.NoError(t, err)
	_, err = h.Client.AdminListLessonProgress(ctx, api.AdminListLessonProgressParams{})
	require.NoError(t, err)
	processReads(t, h)

	for operation, entity := range map[string]string{
		"AdminSearchUsers":        "User",
		"AdminListLeads":          "Lead",
		"AdminListMessages":       "AiMessage",
		"AdminListLessonProgress": "LessonProgress",
	} {
		rows := records(t, h, operation)
		require.Len(t, rows, 1, operation)
		assert.Equal(t, entity, *rows[0].EntityType, operation)
	}
}

// Reading the catalogue or content is not a read of personal data: no fact, no
// record, even after the consumer ran.
func TestReadingContentLeavesNoRecord(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminListCourses(ctx, api.AdminListCoursesParams{})
	require.NoError(t, err)
	_, err = h.Client.AdminListBlogPosts(ctx, api.AdminListBlogPostsParams{})
	require.NoError(t, err)
	_, err = h.Client.AdminListCourseLessons(ctx, api.AdminListCourseLessonsParams{})
	require.NoError(t, err)
	processReads(t, h)

	assert.Zero(t, testsupport.PublishedCount(t, h.Services))
	count, err := h.DB.AuditRecord.Query().Where(auditrecord.OperationIn(
		"AdminListCourses", "AdminListBlogPosts", "AdminListCourseLessons")).Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, count)
}

// A refused read is recorded once, as a refusal, by the security handler; it
// never reaches the middleware, so no read fact follows it.
func TestARefusedReadIsNotAlsoRecordedAsARead(t *testing.T) {
	h := testsupport.NewHarnessAs(t, editorLabel)

	_, err := h.Client.AdminListUsers(context.Background(), api.AdminListUsersParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, h.LastStatus())
	processReads(t, h)

	rows := records(t, h, "AdminListUsers")
	require.Len(t, rows, 1)
	assert.Equal(t, "forbidden", rows[0].Outcome)
}

// Every operation the contract declares as a personal-data read leaves a record,
// whatever it is. The operations come from the contract, so declaring one in the
// TypeSpec adds it here with no edit, and one the middleware does not reach fails.
func TestEveryPersonalDataReadOfTheContractIsRecorded(t *testing.T) {
	operations := permissions.PersonalDataReads()
	require.NotEmpty(t, operations)

	for _, operation := range operations {
		t.Run(operation, func(t *testing.T) {
			h := testsupport.NewHarness(t)
			read := personalDataMiddleware(t, h, h.Events)

			_, err := read(middleware.Request{Context: t.Context(), OperationName: operation}, answers(`{"items":[{"id":7},{"id":9}],"total":2}`))
			require.NoError(t, err)
			processReads(t, h)

			rows := records(t, h, operation)
			require.Len(t, rows, 1, "%s reads personal data and must be recorded", operation)
			entity, _ := permissions.PersonalDataEntity(operation)
			assert.Equal(t, entity, *rows[0].EntityType)
			assert.Equal(t, string(audit.OutcomeSucceeded), rows[0].Outcome)
			var reads readsOf
			require.NoError(t, json.Unmarshal(rows[0].Reads, &reads))
			assert.Equal(t, readsOf{Count: 2, IDs: []string{"7", "9"}}, reads)
		})
	}
}

// Recording must never fail or change a read: a publisher that cannot take the
// fact leaves the response as the handler gave it.
func TestABrokenRecorderDoesNotFailTheRead(t *testing.T) {
	h := testsupport.NewHarness(t)
	read := personalDataMiddleware(t, h, brokenPublisher{})

	resp, err := read(middleware.Request{Context: t.Context(), OperationName: "AdminListLeads"}, answers(`{"items":[{"id":1}]}`))
	require.NoError(t, err)
	assert.NotNil(t, resp.Type)
	processReads(t, h)
	assert.Empty(t, records(t, h, "AdminListLeads"))
}

// A handler that fails, or answers with an error body, read nothing.
func TestAnErrorAnswerIsNotARead(t *testing.T) {
	h := testsupport.NewHarness(t)
	read := personalDataMiddleware(t, h, h.Events)

	_, err := read(middleware.Request{Context: t.Context(), OperationName: "AdminGetUser"}, func(middleware.Request) (middleware.Response, error) {
		return middleware.Response{}, errors.New("boom")
	})
	require.Error(t, err)
	_, err = read(middleware.Request{Context: t.Context(), OperationName: "AdminGetUser"}, answers(`{"status":404,"title":"Not Found"}`))
	require.NoError(t, err)

	assert.Zero(t, testsupport.PublishedCount(t, h.Services))
}

// A page is far smaller than the cap, but a record that grew with the table
// would turn one read into a bulk write: ids stop at the cap, the count does not.
func TestAReadRecordKeepsAtMostTheCapOfIDs(t *testing.T) {
	h := testsupport.NewHarness(t)
	read := personalDataMiddleware(t, h, h.Events)

	items := make([]map[string]int, audit.MaxReadIDs+50)
	for i := range items {
		items[i] = map[string]int{"id": i + 1}
	}
	body, err := json.Marshal(map[string]any{"items": items})
	require.NoError(t, err)

	_, err = read(middleware.Request{Context: t.Context(), OperationName: "AdminListUsers"}, answers(string(body)))
	require.NoError(t, err)
	processReads(t, h)

	var reads readsOf
	require.NoError(t, json.Unmarshal(records(t, h, "AdminListUsers")[0].Reads, &reads))
	assert.Equal(t, audit.MaxReadIDs+50, reads.Count)
	assert.Len(t, reads.IDs, audit.MaxReadIDs)
}

type brokenPublisher struct{}

func (brokenPublisher) PublishStandalone(context.Context, events.Event) error {
	return errors.New("outbox unavailable")
}

func personalDataMiddleware(t *testing.T, h *testsupport.Harness, publisher events.StandalonePublisher) middleware.Middleware {
	t.Helper()
	translator := testsupport.NewTranslator(t)
	return handlers.NewPersonalDataReadMiddleware(publisher, translator, testsupport.NewAPIErrorHandler(t, translator))
}

// answers is a handler that returns the JSON as its response body.
func answers(body string) middleware.Next {
	return func(middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: json.RawMessage(body)}, nil
	}
}
