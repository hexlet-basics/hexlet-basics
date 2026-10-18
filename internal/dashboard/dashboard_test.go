package dashboard_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/dashboard"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/progress"
	"hexletbasics/internal/testsupport"
)

// A course without a catalog card cannot be drawn, so it is left off the
// dashboard; the gap is a content problem nobody sees on the page, so it is
// logged on the server, naming the course and the learner.
func TestCoursesLogsTheCourseWithoutACard(t *testing.T) {
	db, transactor := testsupport.NewClientWithTransactor(t)

	var logs bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelWarn}))
	board := dashboard.New(db, landingpage.NewCatalog(db), progress.New(db, transactor, nil, nil), logger)
	userID := fixture.ID(fixture.UserLearnerHiddenCards)
	cardless := fixture.ID(fixture.CourseSingleLesson)

	courses, err := board.Courses(t.Context(), userID, "ru")
	require.NoError(t, err)

	assert.Len(t, courses.Started, 1, "the course with a card is still listed")
	assert.Contains(t, logs.String(), "level=WARN")
	assert.Contains(t, logs.String(), fmt.Sprintf("course_id=%d", cardless))
	assert.Contains(t, logs.String(), fmt.Sprintf("user_id=%d", userID))
}
