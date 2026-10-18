package progress

import (
	"context"

	"github.com/samber/lo"

	"hexletbasics/ent/course"
	"hexletbasics/ent/lessonprogress"
)

// StateNotStarted is the condition of a learner with no record: the glossary's
// Not started is the absence of a row, so it exists only in this read model and
// is never stored.
const StateNotStarted = "not_started"

// LearnerCourse is a signed-in learner's position in one Course told in the
// glossary's three conditions. The contract's CourseProgress says only whether
// each Lesson is finished; a reader that must tell a Lesson the learner has
// begun from one they have never opened (the MCP learner tools) needs Started
// too, which is why this read model exists next to CourseState.
type LearnerCourse struct {
	CourseSlug string
	// State is the Enrollment's: StateNotStarted without one.
	State          string
	Completion     int
	NextLessonSlug *string
	Lessons        []LearnerLesson
}

// LearnerLesson is one Current Lesson with the learner's Lesson Progress.
type LearnerLesson struct {
	Slug     string
	Position int
	// State is StateFinished, StateStarted or StateNotStarted.
	State string
}

// LearnerCourse reads where the user stands in the Course with this slug. The
// subject is always the explicit user id; there is no guest variant, which is
// what a bearer token's user is. An unknown slug is an ent.NotFoundError.
//
// It builds on CourseStates, so Completion and the Next Lesson are the numbers
// the course page and the dashboard show; it only adds the Started rows.
func (p *Progress) LearnerCourse(ctx context.Context, userID int, courseSlug string) (*LearnerCourse, error) {
	crs, err := p.db.Course.Query().Where(course.SlugEQ(courseSlug)).Only(ctx)
	if err != nil {
		return nil, errs.With("course_slug", courseSlug).Wrapf(err, "load course")
	}

	states, err := p.CourseStates(ctx, userID, []int{crs.ID})
	if err != nil {
		return nil, err
	}
	state := states[crs.ID]

	lessons, err := currentLessons(ctx, p.db, crs)
	if err != nil {
		return nil, err
	}
	// Any row that is not finished counts as started, NULL included, the same
	// reading finishCourse takes of a legacy row without a state.
	rows, err := p.db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(userID),
			lessonprogress.CourseID(crs.ID),
			lessonprogress.Or(lessonprogress.StateNEQ(StateFinished), lessonprogress.StateIsNil()),
		).
		Select(lessonprogress.FieldLessonID).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load started lessons")
	}
	slugByID := make(map[int]string, len(lessons))
	for _, l := range lessons {
		slugByID[l.lessonID] = l.slug
	}
	started := make(map[string]bool, len(rows))
	for _, row := range rows {
		started[slugByID[row.LessonID]] = true
	}

	out := &LearnerCourse{
		CourseSlug:     courseSlug,
		State:          lo.Ternary(state.State == "", StateNotStarted, state.State),
		Completion:     state.Completion,
		NextLessonSlug: state.NextLessonSlug,
		Lessons:        make([]LearnerLesson, 0, len(state.Lessons)),
	}
	for _, l := range state.Lessons {
		lesson := LearnerLesson{Slug: l.Slug, Position: l.Position, State: StateNotStarted}
		switch {
		case l.Finished:
			lesson.State = StateFinished
		case started[l.Slug]:
			lesson.State = StateStarted
		}
		out.Lessons = append(out.Lessons, lesson)
	}
	return out, nil
}
