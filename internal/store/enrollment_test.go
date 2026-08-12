package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/enrollment"
	"hexletbasics/internal/testsupport"
)

// enrolledCourseSlug is the course the fixture Enrollments belong to;
// otherCourseSlug is any other published course. Both are addressed by slug,
// never by id — fixture ids are legacy crc32 values.
const (
	enrolledCourseSlug = "javascript"
	otherCourseSlug    = "ruby"
)

// The invariant this file guards belongs to PostgreSQL, not to Go: an
// Enrollment is never created twice for the same (learner, Course). The legacy
// find-or-create read before it wrote, so only a real constraint makes the
// later progress tickets safe to write as conflict-tolerant inserts.
func TestSecondEnrollmentForTheSamePairIsRejectedByTheDatabase(t *testing.T) {
	client := testsupport.NewClient(t)
	existing := firstEnrollmentIn(t, client, enrolledCourseSlug)

	_, err := client.Enrollment.Create().
		SetUserID(existing.UserID).
		SetCourseID(existing.CourseID).
		SetState("started").
		Save(t.Context())

	require.Error(t, err)
	assert.True(t, ent.IsConstraintError(err), "want a constraint error, got %v", err)
}

// The same learner in a different Course is not a duplicate — the constraint is
// on the pair, not on the learner.
func TestSecondEnrollmentInAnotherCourseIsAccepted(t *testing.T) {
	client := testsupport.NewClient(t)
	existing := firstEnrollmentIn(t, client, enrolledCourseSlug)
	other := client.Course.Query().
		Where(course.SlugEQ(otherCourseSlug)).
		OnlyX(t.Context())

	created, err := client.Enrollment.Create().
		SetUserID(existing.UserID).
		SetCourseID(other.ID).
		SetState("started").
		Save(t.Context())

	require.NoError(t, err)
	assert.Equal(t, other.ID, created.CourseID)
}

// Lesson Progress rows reach their owning Enrollment through the edge the
// collapse migration re-points, so the mapping is exercised rather than assumed.
func TestLessonProgressResolvesItsEnrollment(t *testing.T) {
	client := testsupport.NewClient(t)

	progress := client.LessonProgress.Query().FirstX(t.Context())
	owner := progress.QueryEnrollment().OnlyX(t.Context())

	assert.Equal(t, progress.UserID, owner.UserID)
	assert.Equal(t, progress.CourseID, owner.CourseID)
}

func firstEnrollmentIn(t *testing.T, client *ent.Client, slug string) *ent.Enrollment {
	t.Helper()
	return client.Enrollment.Query().
		Where(enrollment.HasCourseWith(course.SlugEQ(slug))).
		FirstX(t.Context())
}
