package store_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/testsupport"
)

// collapseScript de-duplicates Enrollments. The migration that adds the unique
// index deliberately does not clean up — it fails loudly instead — so this SQL
// is run by hand against production, once, after the count (#765). These tests
// execute the file rather than a restatement of it, so what is verified is what
// will be run.
var collapseScript = filepath.Join("..", "..", "scripts", "oneoff", "765-duplicate-enrollments", "collapse.sql")

// The winner is the oldest row, it absorbs the losers' Lesson Progress rows and
// their counters, and the losers are gone.
//
// The unique index comes off first: the migrated container already has it, which
// is the only reason duplicates cannot exist here to begin with. The drop, the
// duplicates and the re-creation all vanish with the test's rollback.
func TestCollapseMergesDuplicateEnrollmentsOntoTheOldest(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	existing := firstEnrollmentIn(t, client, enrolledCourseSlug)
	progressBefore := progressCount(t, client, existing.ID)

	execSQL(t, tx, `DROP INDEX "index_language_members_on_user_id_and_language_id"`)

	// Two losers: one plainly newer, one sharing the winner's exact created_at
	// so that only the id tie-break can order it.
	newer := createEnrollment(t, client, existing, existing.CreatedAt.Add(time.Hour), 3)
	createEnrollment(t, client, existing, existing.CreatedAt, 4)
	attachProgress(t, client, newer)

	for _, stmt := range collapseStatements(t) {
		execSQL(t, tx, stmt)
	}

	survivors := client.Enrollment.Query().
		Where(enrollment.UserID(existing.UserID), enrollment.CourseID(existing.CourseID)).
		AllX(t.Context())

	require.Len(t, survivors, 1, "one Enrollment survives per (learner, course)")
	assert.Equal(t, existing.ID, survivors[0].ID, "the oldest row wins, id breaking the created_at tie")
	assert.Equal(t, existing.FinishedLessonsCount+3+4, survivors[0].FinishedLessonsCount,
		"the losers' counters are summed onto the winner")
	assert.Equal(t, progressBefore+1, progressCount(t, client, existing.ID),
		"the loser's Lesson Progress row is re-pointed at the winner")

	// The collapsed data admits the index again — which is exactly what the
	// migration needs of production before it will apply. Asserted last: the
	// conflict aborts the shared transaction.
	execSQL(t, tx, `CREATE UNIQUE INDEX "index_language_members_on_user_id_and_language_id" ON "language_members" ("user_id", "language_id")`)
	_, err := client.Enrollment.Create().
		SetUserID(existing.UserID).
		SetCourseID(existing.CourseID).
		SetState("started").
		Save(t.Context())
	require.Error(t, err, "the pair is unique again once the duplicates are collapsed")
}

// Against a database that already has no duplicates the script is a no-op: the
// CTEs select nothing, so re-running it after a collapse writes nothing.
func TestCollapseIsIdempotentWithoutDuplicates(t *testing.T) {
	client, tx := testsupport.NewClientWithTx(t)
	before := client.Enrollment.Query().CountX(t.Context())
	progressBefore := client.LessonProgress.Query().CountX(t.Context())

	for _, stmt := range collapseStatements(t) {
		execSQL(t, tx, stmt)
	}

	assert.Equal(t, before, client.Enrollment.Query().CountX(t.Context()))
	assert.Equal(t, progressBefore, client.LessonProgress.Query().CountX(t.Context()))
}

// collapseStatements strips the commentary and splits the script into the three
// statements it ships: re-point, sum, delete.
func collapseStatements(t *testing.T) []string {
	t.Helper()
	raw, err := os.ReadFile(collapseScript)
	require.NoError(t, err)

	var code []string
	for line := range strings.SplitSeq(string(raw), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			code = append(code, line)
		}
	}

	var stmts []string
	for part := range strings.SplitSeq(strings.Join(code, "\n"), ";") {
		if strings.TrimSpace(part) != "" {
			stmts = append(stmts, part)
		}
	}
	require.Len(t, stmts, 3, "re-point, sum, delete")
	return stmts
}

func execSQL(t *testing.T, tx *sql.Tx, stmt string) {
	t.Helper()
	_, err := tx.ExecContext(t.Context(), stmt)
	require.NoErrorf(t, err, "exec %.60s", strings.TrimSpace(stmt))
}

// createEnrollment inserts a duplicate of existing with its own created_at and
// counter. Only callable while the unique index is dropped.
func createEnrollment(t *testing.T, client *ent.Client, existing *ent.Enrollment, createdAt time.Time, finished int) *ent.Enrollment {
	t.Helper()
	return client.Enrollment.Create().
		SetUserID(existing.UserID).
		SetCourseID(existing.CourseID).
		SetState("started").
		SetFinishedLessonsCount(finished).
		SetCreatedAt(createdAt).
		SaveX(t.Context())
}

// attachProgress hangs one Lesson Progress row off the given Enrollment, on a
// lesson none of the fixture rows use so the (user, lesson) uniqueness holds.
func attachProgress(t *testing.T, client *ent.Client, owner *ent.Enrollment) {
	t.Helper()
	taken := map[int]bool{}
	for _, row := range client.LessonProgress.Query().
		Where(lessonprogress.UserID(owner.UserID)).
		AllX(t.Context()) {
		taken[row.LessonID] = true
	}

	lessons := client.CourseLesson.Query().
		Where(courselesson.CourseID(owner.CourseID)).
		AllX(t.Context())
	free, found := lo.Find(lessons, func(l *ent.CourseLesson) bool { return !taken[l.ID] })
	require.True(t, found, "fixtures leave a lesson this learner has no progress on")

	client.LessonProgress.Create().
		SetUserID(owner.UserID).
		SetCourseID(owner.CourseID).
		SetEnrollmentID(owner.ID).
		SetLessonID(free.ID).
		SetState("started").
		SaveX(t.Context())
}

func progressCount(t *testing.T, client *ent.Client, enrollmentID int) int {
	t.Helper()
	return client.LessonProgress.Query().
		Where(lessonprogress.EnrollmentID(enrollmentID)).
		CountX(t.Context())
}
