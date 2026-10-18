package progress

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"time"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/lessonprogress"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/events"
)

// Learner is who is submitting. A signed-in learner's position is rows; a
// guest's is the cookie they carry. Everything above the recording step treats
// the two identically — resolve the Lesson in its Course — which is what
// lets one handler and one client code path serve both.
type Learner struct {
	// UserID is 0 for a guest.
	UserID int
	// Guest is the visitor's cookie progress, empty for a first visit and
	// ignored entirely once UserID is set.
	Guest GuestProgress
}

// SignedIn reports whether the submission belongs to an account.
func (l Learner) SignedIn() bool { return l.UserID != 0 }

// Check is one submission: who is submitting, what they wrote, and the Lesson
// Version they wrote it against.
type Check struct {
	Learner   Learner
	LessonID  int
	VersionID int
	Code      string
	Locale    string
}

// CheckResult is what the check produced and what it changed.
type CheckResult struct {
	Outcome        Outcome
	LessonFinished bool
	CourseFinished bool

	// Guest carries a visitor's updated progress, and is nil whenever there is
	// no cookie to write — for every signed-in learner, and for a guest whose
	// submission finished nothing new. The handler turns a non-nil value into the
	// response cookie; it never has to ask who it is serving.
	Guest *GuestProgress
}

// CheckSolution runs a submitted solution and records what it changed.
//
// The Lesson is resolved BEFORE the run: a submission for a Lesson the Course
// does not have is refused without burning a container and without writing
// anything. The run itself happens outside the transaction — it is a Docker
// call taking seconds, and holding a transaction open across it would tie up a
// connection per submission for no benefit.
//
// Recording does not require the Lesson to have been started. A learner who
// arrived by a direct link and submitted gets their Enrollment and Lesson
// Progress created lazily, so the route they took cannot lose their work.
func (p *Progress) CheckSolution(ctx context.Context, check Check) (*CheckResult, error) {
	open, err := p.resolveLesson(ctx, p.db, check.LessonID)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	version, err := p.submittedVersion(ctx, open.course, check)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	submission, err := p.submissionFor(ctx, version, check)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	outcome, err := p.runner.Run(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("run submission for lesson %d: %w", check.LessonID, err)
	}

	result := &CheckResult{Outcome: outcome}
	err = p.store.WithinTx(ctx, func(tx *sql.Tx, db *ent.Client) error {
		occurredAt := p.now()

		// Published for everyone, guests included: pre-signup activity is exactly
		// what this fact exists to keep visible.
		if err := p.publisher.Publish(ctx, tx, events.SolutionChecked{
			LessonSlug: lo.FromPtr(open.lesson.Slug),
			CourseSlug: lo.FromPtr(open.course.Slug),
			Locale:     check.Locale,
			Passed:     outcome.Passed,
			OccurredAt: occurredAt,
		}); err != nil {
			return fmt.Errorf("publish solution checked: %w", err)
		}

		if !check.Learner.SignedIn() {
			// A guest leaves no rows and produces no enrollment facts: those count
			// transitions on an account, and there is no account. Their whole state
			// is the set of Lessons they finished, so a pass that adds nothing —
			// re-submitting one they already passed — changes nothing and writes
			// no cookie.
			if outcome.Passed && !check.Learner.Guest.Has(lo.FromPtr(open.course.Slug), lo.FromPtr(open.lesson.Slug)) {
				advanced := check.Learner.Guest.Record(lo.FromPtr(open.course.Slug), lo.FromPtr(open.lesson.Slug))
				result.Guest = &advanced
			}
			return nil
		}
		return p.record(ctx, tx, db, check, open, outcome, result, occurredAt)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// submittedVersion resolves the Lesson Version the code was written against. It
// must belong to this Lesson and to the Course's current Version: a submission
// against a superseded build is refused rather than silently run against the
// current one, because the tests it was written for no longer exist.
func (p *Progress) submittedVersion(ctx context.Context, crs *ent.Course, check Check) (*ent.CourseLessonVersion, error) {
	if crs.CurrentVersionID == nil {
		return nil, &ent.NotFoundError{}
	}

	version, err := p.db.CourseLessonVersion.Query().
		Where(
			courselessonversion.ID(check.VersionID),
			courselessonversion.LessonID(check.LessonID),
			courselessonversion.CourseVersionID(*crs.CurrentVersionID),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, &ent.NotFoundError{}
	}
	if err != nil {
		return nil, fmt.Errorf("load lesson version %d: %w", check.VersionID, err)
	}
	return version, nil
}

// submissionFor resolves what running this submission takes: the image the
// Course Version was built as, and the two paths inside it.
//
// The image is required. A Version built before the loader recorded one cannot
// be graded, and running some default image instead would grade the submission
// against the wrong language — so it reads as not found, the same answer as a
// Lesson with no content.
func (p *Progress) submissionFor(
	ctx context.Context,
	version *ent.CourseLessonVersion,
	check Check,
) (Submission, error) {
	built, err := p.db.CourseVersion.Get(ctx, version.CourseVersionID)
	if err != nil {
		return Submission{}, fmt.Errorf("load course version %d: %w", version.CourseVersionID, err)
	}

	image := lo.FromPtr(built.DockerImage)
	testDir := lo.FromPtr(version.PathToCode)
	exerciseFile := lo.FromPtr(built.ExerciseFilename)
	if image == "" || testDir == "" || exerciseFile == "" {
		return Submission{}, &ent.NotFoundError{}
	}

	return Submission{
		Image:        image,
		VersionID:    built.ID,
		TestDir:      testDir,
		ExerciseFile: path.Join(testDir, exerciseFile),
		UserID:       check.Learner.UserID,
		Code:         check.Code,
	}, nil
}

// record writes what a signed-in learner's check changed, and publishes the
// fact for each transition it made.
func (p *Progress) record(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	check Check,
	open *resolvedLesson,
	outcome Outcome,
	result *CheckResult,
	occurredAt time.Time,
) error {
	userID := check.Learner.UserID

	enrolled, createdEnrollment, err := p.enroll(ctx, db, userID, open.course.ID)
	if err != nil {
		return err
	}
	if createdEnrollment {
		if err := p.publishCourseStarted(ctx, tx, db, userID, open.course, check.Locale, occurredAt); err != nil {
			return err
		}
	}

	taken, createdProgress, err := takeLesson(ctx, db, userID, open.course.ID, enrolled.ID, check.LessonID)
	if err != nil {
		return err
	}
	if createdProgress {
		// The learner reached this Lesson without pressing start — a deep link,
		// or a page restored from history. Submitting is as deliberate an act as
		// starting, so it produces the same fact.
		if err := p.publishLessonStarted(ctx, tx, db, enrolled.ID, open, check.Locale, occurredAt); err != nil {
			return err
		}
	}

	if !outcome.Passed {
		// A failed attempt leaves the Lesson started: it counts as work done on
		// it, and nothing beyond the attempt happened.
		return nil
	}

	// Re-submitting to a Lesson already passed transitions nothing and
	// publishes nothing — experimenting must stay safe, and the occurrence
	// counts the CRM consumes count transitions, not submissions.
	if lo.FromPtr(taken.State) != StateFinished {
		if _, err := taken.Update().SetState(StateFinished).Save(ctx); err != nil {
			return fmt.Errorf("finish lesson progress: %w", err)
		}
		result.LessonFinished = true

		count, err := lessonProgressCount(ctx, db, enrolled.ID)
		if err != nil {
			return err
		}
		if err := p.publisher.Publish(ctx, tx, events.LessonFinished{
			OccurrenceCount: count,
			LessonSlug:      lo.FromPtr(open.lesson.Slug),
			CourseSlug:      lo.FromPtr(open.course.Slug),
			Locale:          check.Locale,
			OccurredAt:      occurredAt,
		}); err != nil {
			return fmt.Errorf("publish lesson finished: %w", err)
		}
	}

	// Completion is evaluated on every passing check, not only on one that
	// transitioned a Lesson, exactly as the legacy check did. A learner whose
	// Lessons are all finished while their Enrollment is not — inherited data,
	// or a Version that shrank the Lesson set — closes their Course by passing
	// anything, instead of having no submission left that would notice.
	return p.finishCourse(ctx, tx, db, check, open, enrolled, result, occurredAt)
}

// finishCourse moves the Enrollment to finished when this check left no
// unfinished Lesson in the Course's current Version.
//
// Every Lesson of the current Version must be finished, not merely the ones up
// to the learner's position: a gap — inherited from the old system, which
// allowed any order, or left by a Version that inserted a Lesson behind them —
// keeps the Course open until the learner closes it.
func (p *Progress) finishCourse(
	ctx context.Context,
	tx *sql.Tx,
	db *ent.Client,
	check Check,
	open *resolvedLesson,
	enrolled *ent.Enrollment,
	result *CheckResult,
	occurredAt time.Time,
) error {
	// Anything other than finished counts as in progress, NULL included: the
	// baseline allows it and rows predating AASM have it, and such a learner must
	// not be the one learner who can never complete.
	if lo.FromPtr(enrolled.State) == StateFinished {
		return nil
	}
	if len(open.lessons) == 0 {
		// A Version with no Lessons completes nobody; the alternative reading —
		// everybody has finished all zero of them — is the worse one.
		return nil
	}

	finished, err := finishedLessonIDs(ctx, db, check.Learner.UserID, open.course.ID)
	if err != nil {
		return err
	}
	if !lo.EveryBy(open.lessons, func(lesson currentLesson) bool { return finished[lesson.lessonID] }) {
		return nil
	}

	if _, err := enrolled.Update().SetState(StateFinished).Save(ctx); err != nil {
		return fmt.Errorf("finish enrollment: %w", err)
	}
	result.CourseFinished = true

	// Counted after the transition so the Course just finished is included —
	// the legacy publisher counted the transition itself.
	count, err := db.Enrollment.Query().
		Where(
			enrollment.UserID(check.Learner.UserID),
			enrollment.StateEQ(StateFinished),
		).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("count finished enrollments: %w", err)
	}
	if err := p.publisher.Publish(ctx, tx, events.CourseFinished{
		OccurrenceCount: count,
		Slug:            lo.FromPtr(open.course.Slug),
		Locale:          check.Locale,
		OccurredAt:      occurredAt,
	}); err != nil {
		return fmt.Errorf("publish course finished: %w", err)
	}
	return nil
}

// lessonProgressCount is the number of Lesson Progress rows on an Enrollment,
// which is the occurrence count the lesson facts carry.
func lessonProgressCount(ctx context.Context, db *ent.Client, enrollmentID int) (int, error) {
	count, err := db.LessonProgress.Query().
		Where(lessonprogress.EnrollmentID(enrollmentID)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count lesson progress: %w", err)
	}
	return count, nil
}
