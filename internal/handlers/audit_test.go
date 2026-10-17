package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/auditrecord"
	"hexletbasics/ent/blogpost"
	"hexletbasics/internal/api"
	"hexletbasics/internal/audit"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

// The audit trail of ADR-0017: an admin change of a blog post leaves exactly one
// Audit Record, written in the change's own transaction.

// records returns the Audit Records of an operation, oldest first.
func records(t *testing.T, h *testsupport.Harness, operation string) []*ent.AuditRecord {
	t.Helper()
	rows, err := h.DB.AuditRecord.Query().
		Where(auditrecord.Operation(operation)).
		Order(ent.Asc(auditrecord.FieldID)).
		All(t.Context())
	require.NoError(t, err)
	return rows
}

func changesOf(t *testing.T, record *ent.AuditRecord) []audit.Change {
	t.Helper()
	var changes []audit.Change
	require.NoError(t, json.Unmarshal(record.Changes, &changes))
	return changes
}

func fieldOf(changes []audit.Change, name string) (audit.FieldChange, bool) {
	for _, c := range changes {
		for _, f := range c.Fields {
			if f.Name == name {
				return f, true
			}
		}
	}
	return audit.FieldChange{}, false
}

// blogInput is a complete input: the contract requires every field, null or not.
func blogInput(name, slug string) *api.BlogPostInput {
	return &api.BlogPostInput{
		Name:        api.NewNilString(name),
		Slug:        api.NewNilString(slug),
		Description: api.NewNilString("d"),
		RichBody:    "<p>Body</p>",
	}
}

func TestCreatingABlogPostLeavesOneAuditRecord(t *testing.T) {
	h := testsupport.NewHarness(t)
	testsupport.SpeakTo(h, "ru")
	testsupport.ArriveFrom(h, "203.0.113.7")
	testsupport.Identify(h, "req-create-1", "AdminUI/1.0")

	post, err := h.Client.AdminCreateBlogPost(context.Background(), blogInput("Audited", "audited"))
	require.NoError(t, err)

	rows := records(t, h, "AdminCreateBlogPost")
	require.Len(t, rows, 1)
	row := rows[0]
	assert.Equal(t, h.UserID, *row.ActorUserID)
	assert.Equal(t, "web", row.Channel)
	assert.Equal(t, "req-create-1", *row.RequestID)
	assert.Equal(t, "203.0.113.7", *row.Address)
	assert.Equal(t, "AdminUI/1.0", *row.UserAgent)
	assert.Equal(t, "ru", *row.Locale)
	assert.Equal(t, "succeeded", row.Outcome)
	assert.Equal(t, "BlogPost", *row.EntityType)
	assert.Equal(t, strconv.Itoa(int(post.ID)), *row.EntityID)
	assert.Nil(t, row.ConnectedAppID)
	assert.Nil(t, row.TokenID)

	name, ok := fieldOf(changesOf(t, row), "name")
	require.True(t, ok)
	assert.Empty(t, name.Old, "a create has no old value")
	assert.JSONEq(t, `"Audited"`, string(name.New))
}

func TestUpdatingABlogPostRecordsOldAndNewOfEachChangedField(t *testing.T) {
	h := testsupport.NewHarness(t)

	// Second post: slug and state are re-sent unchanged, name and body change,
	// description is cleared. Only the changed three belong in the record.
	before := h.DB.BlogPost.GetX(t.Context(), fixtures.ID(fixtures.BlogPostSecondPost))
	_, err := h.Client.AdminUpdateBlogPost(context.Background(), &api.BlogPostInput{
		Name:        api.NewNilString("Renamed"),
		Slug:        api.NewNilString(*before.Slug),
		Description: api.NilString{Null: true},
		RichBody:    "<p>new</p>",
	}, api.AdminUpdateBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostSecondPost)})
	require.NoError(t, err)

	rows := records(t, h, "AdminUpdateBlogPost")
	require.Len(t, rows, 1)
	assert.Equal(t, strconv.Itoa(fixtures.ID(fixtures.BlogPostSecondPost)), *rows[0].EntityID)

	changes := changesOf(t, rows[0])
	require.Len(t, changes, 1)
	assert.Equal(t, "update", changes[0].Op)

	name, ok := fieldOf(changes, "name")
	require.True(t, ok)
	assert.JSONEq(t, `"`+*before.Name+`"`, string(name.Old))
	assert.JSONEq(t, `"Renamed"`, string(name.New))

	body, ok := fieldOf(changes, "rich_body")
	require.True(t, ok)
	assert.JSONEq(t, `"<p>new</p>"`, string(body.New))

	if before.Description != nil {
		desc, ok := fieldOf(changes, "description")
		require.True(t, ok, "a cleared field is a change")
		assert.JSONEq(t, `"`+*before.Description+`"`, string(desc.Old))
		assert.Empty(t, desc.New)
	}

	_, slugChanged := fieldOf(changes, "slug")
	assert.False(t, slugChanged, "an unchanged field is not recorded")
	_, stamped := fieldOf(changes, "updated_at")
	assert.False(t, stamped, "ent-owned timestamps are not recorded")
}

func TestDeletingABlogPostLeavesOneAuditRecordForThePost(t *testing.T) {
	h := testsupport.NewHarness(t)

	require.NoError(t, h.Client.AdminDeleteBlogPost(context.Background(), api.AdminDeleteBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)}))

	rows := records(t, h, "AdminDeleteBlogPost")
	require.Len(t, rows, 1, "the dependents' deletes are part of this one record, not records of their own")
	assert.Equal(t, "BlogPost", *rows[0].EntityType)
	assert.Equal(t, strconv.Itoa(fixtures.ID(fixtures.BlogPostHelloWorld)), *rows[0].EntityID)

	var entities []string
	for _, c := range changesOf(t, rows[0]) {
		entities = append(entities, c.Entity)
	}
	assert.Contains(t, entities, "BlogPostLike")
	assert.Contains(t, entities, "BlogPost")

	total, err := h.DB.AuditRecord.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, total)
}

func TestAReadLeavesNoAuditRecord(t *testing.T) {
	h := testsupport.NewHarness(t)

	_, err := h.Client.AdminGetBlogPost(context.Background(), api.AdminGetBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostHelloWorld)})
	require.NoError(t, err)

	total, err := h.DB.AuditRecord.Query().Count(t.Context())
	require.NoError(t, err)
	assert.Zero(t, total)
}

func TestAFailedChangeIsUndoneAndRecordedAsFailedWithoutItsChanges(t *testing.T) {
	h := testsupport.NewHarness(t)
	name := *h.DB.BlogPost.GetX(t.Context(), fixtures.ID(fixtures.BlogPostSecondPost)).Name

	// The update itself succeeds, then the handler reads the likes to answer and
	// that read fails: the operation fails after its write.
	breakReads(t, h, "blog_post_likes")
	_, err := h.Client.AdminUpdateBlogPost(context.Background(), blogInput("Never kept", "second-post"), api.AdminUpdateBlogPostParams{ID: fixtures.ID32(fixtures.BlogPostSecondPost)})
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())

	assert.Equal(t, name, *h.DB.BlogPost.GetX(t.Context(), fixtures.ID(fixtures.BlogPostSecondPost)).Name, "the write is undone")

	// The success record went with the rollback; what stays is a record of the
	// attempt, written on its own, which names no field the failed call tried to
	// write (the name above, the new body) and no entity.
	rows := records(t, h, "AdminUpdateBlogPost")
	require.Len(t, rows, 1)
	assert.Equal(t, "failed", rows[0].Outcome)
	assert.Equal(t, h.UserID, *rows[0].ActorUserID)
	assert.Empty(t, rows[0].Changes)
	assert.Nil(t, rows[0].EntityType)
	assert.Nil(t, rows[0].EntityID)
}

func TestAChangeIsNotCommittedWhenItsRecordCannotBeWritten(t *testing.T) {
	h := testsupport.NewHarness(t)

	breakWrites(t, h, "audit_records", "INSERT")
	_, err := h.Client.AdminCreateBlogPost(context.Background(), blogInput("Unaudited", "unaudited"))
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, h.LastStatus())

	exists, err := h.DB.BlogPost.Query().Where(blogpost.Slug("unaudited")).Exist(t.Context())
	require.NoError(t, err)
	assert.False(t, exists, "no change without its record")
}

func TestTheDatabaseRefusesToUpdateOrDeleteAnAuditRecord(t *testing.T) {
	h := testsupport.NewHarness(t)
	_, err := h.Client.AdminCreateBlogPost(context.Background(), blogInput("Kept", "kept"))
	require.NoError(t, err)

	// Each refused statement aborts the transaction, so each runs in a savepoint
	// that is rolled back, leaving the next one a usable transaction.
	for _, statement := range []string{
		`UPDATE audit_records SET outcome = 'failed'`,
		`DELETE FROM audit_records`,
	} {
		exec(t, h, `SAVEPOINT refused`)
		_, execErr := h.Tx.ExecContext(t.Context(), statement)
		require.Error(t, execErr, statement)
		assert.Contains(t, execErr.Error(), "append-only", statement)
		exec(t, h, `ROLLBACK TO SAVEPOINT refused`)
	}

	assert.Len(t, records(t, h, "AdminCreateBlogPost"), 1)
}
