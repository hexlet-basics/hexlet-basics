// Package coursepage reads the two public pages of a Course: the course page
// and the lesson page. It owns the rules that decide what a visitor is shown —
// which landing page a slug names, where an archived page moved, which page is
// the course's main one, what makes a lesson not found, where a lesson's source
// lives — so they are testable without HTTP, cookies or a request locale.
//
// The reads return internal structures (ent rows plus the learner's progress
// state), never API types: handlers convert them through apiconv.
package coursepage

import (
	"context"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/samber/lo"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/ent/courselandingpage"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/courselessontranslation"
	"hexletbasics/ent/courseversion"
	"hexletbasics/ent/enrollment"
	"hexletbasics/ent/landingpageqnaitem"
	"hexletbasics/internal/curriculum"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/landingpage"
	"hexletbasics/internal/progress"
)

const stateArchived = "archived"

// progressReader is what the pages need from the progress module: where a
// learner stands in a course, in the order the course gives its lessons.
type progressReader interface {
	CourseState(ctx context.Context, learner progress.Learner, courseID int) (*progress.CourseState, error)
}

// Pages reads the course page and the lesson page.
type Pages struct {
	db       *ent.Client
	catalog  *landingpage.Catalog
	progress progressReader
	// sourceBase is the content organisation under which each course's exercises
	// repository lives.
	sourceBase string
}

// New builds the reader. The catalog owns the main-page rule and the progress
// tracker owns the order of a course's lessons, so the page and the learner's
// standing cannot drift apart.
func New(db *ent.Client, catalog *landingpage.Catalog, tracker progressReader, sourceBase string) *Pages {
	return &Pages{db: db, catalog: catalog, progress: tracker, sourceBase: sourceBase}
}

// CoursePage is everything the course page shows.
type CoursePage struct {
	Course *ent.Course
	// Landing is the page the slug names, not the course's main page, so a
	// secondary page keeps its own copy.
	Landing *ent.CourseLandingPage
	// MovedTo is the slug an archived Landing moved to; nil for a live page.
	MovedTo *string
	// State is where the visitor stands; a guest and a learner have the same
	// shape.
	State *progress.CourseState
	// Infos are the lessons of the current version in the locale, each with its
	// Lesson loaded, in course order.
	Infos []*ent.CourseLessonTranslation
	// Modules are the current version's modules in the locale.
	Modules  []*ent.CourseModuleTranslation
	QnaItems []*ent.LandingPageQnaItem
	// Enrollment is nil when the visitor has not started the course (a guest, or
	// a learner who has not begun): the absence of a record, not a state.
	Enrollment *ent.Enrollment
}

// Course reads the course page of the landing page with slug in locale.
//
// The slug is a landing page's, not the course's: that is what legacy's
// `languages#show` resolved and what every public link carries, so existing URLs
// keep working (ADR-0002). A course has one landing page per locale, so the same
// slug under another locale is not found rather than the other locale's copy.
func (p *Pages) Course(ctx context.Context, learner progress.Learner, locale, slug string) (*CoursePage, error) {
	landing, err := p.db.CourseLandingPage.Query().
		Where(courselandingpage.SlugEQ(slug), courselandingpage.LocaleEQ(locale)).
		WithCourse(func(q *ent.CourseQuery) { q.WithCurrentVersion() }).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, errs.With("slug", slug).Wrapf(err, "load landing page")
	}

	crs := landing.Edges.Course

	target, moved, err := p.archivedTarget(ctx, landing, locale)
	if err != nil {
		return nil, err
	}

	var movedTo *string
	if moved {
		movedTo = &target
	}

	// Everyone has a standing: a visitor who has finished nothing stands at the
	// beginning, a fact about the course and not about having an account.
	state, err := p.progress.CourseState(ctx, learner, crs.ID)
	if err != nil {
		return nil, err
	}

	infos, err := p.lessonInfos(ctx, crs, locale, state)
	if err != nil {
		return nil, err
	}

	var modules []*ent.CourseModuleTranslation
	if crs.CurrentVersionID != nil {
		modules, err = curriculum.Modules(ctx, p.db, locale, *crs.CurrentVersionID)
		if err != nil {
			return nil, err
		}
	}

	// Oldest first, as legacy's `qna_items` association returned them.
	qna, err := p.db.LandingPageQnaItem.Query().
		Where(landingpageqnaitem.CourseLandingPageID(landing.ID)).
		Order(landingpageqnaitem.ByID()).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load landing page questions")
	}

	page := &CoursePage{
		Course:   crs,
		Landing:  landing,
		MovedTo:  movedTo,
		State:    state,
		Infos:    infos,
		Modules:  modules,
		QnaItems: qna,
	}
	if !learner.SignedIn() {
		return page, nil
	}

	enrolled, err := p.db.Enrollment.Query().
		Where(enrollment.UserID(learner.UserID), enrollment.CourseID(crs.ID)).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return page, nil
	case err != nil:
		return nil, errs.Wrapf(err, "load enrollment")
	}

	page.Enrollment = enrolled

	return page, nil
}

// Lesson reads the lesson page: the lesson with lessonSlug of the course with
// courseSlug, in the course's current version and locale.
//
// Reading never starts the lesson (ADR-0012): the frontend preloads on hover.
// It answers everyone, including a visitor on a lesson they may not take yet:
// theory is public and is what search engines index.
//
// The lesson is resolved against the CURRENT version in one query on its
// translations: a lesson only a retired build contained is gone from the course,
// and a lesson not translated into the locale is not one this visitor can read.
// Both, like every other way of missing, are domain.ErrNotFound.
func (p *Pages) Lesson(ctx context.Context, learner progress.Learner, locale, courseSlug, lessonSlug string) (*LessonPage, error) {
	translation, err := p.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(locale),
			courselessontranslation.HasLessonWith(courselesson.SlugEQ(lessonSlug)),
			courselessontranslation.HasCourseVersionWith(
				courseversion.HasCurrentCoursesWith(course.SlugEQ(courseSlug)),
			),
		).
		WithLesson().
		WithVersion().
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, errs.With("lesson_slug", lessonSlug, "course_slug", courseSlug).Wrapf(err, "load lesson")
	}

	if translation.Edges.Lesson == nil || translation.Edges.Version == nil {
		return nil, domain.ErrNotFound
	}

	crs, err := p.db.Course.Query().
		Where(course.ID(translation.CourseID)).
		WithCurrentVersion().
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, errs.With("course_slug", courseSlug).Wrapf(err, "load course")
	}

	state, err := p.progress.CourseState(ctx, learner, crs.ID)
	if err != nil {
		return nil, err
	}

	infos, err := p.lessonInfos(ctx, crs, locale, state)
	if err != nil {
		return nil, err
	}

	// The same rule the dashboard card and the archived-page redirect use.
	mains, err := p.catalog.MainPages(ctx, locale, crs.ID)
	if err != nil {
		return nil, err
	}

	version := translation.Edges.Version

	return &LessonPage{
		Course:      crs,
		Lesson:      translation.Edges.Lesson,
		Version:     version,
		Translation: translation,
		SourceURL:   lessonSourceURL(p.sourceBase, version.PathToCode, locale),
		Landing:     mains[crs.ID],
		State:       state,
		Infos:       infos,
	}, nil
}

// archivedTarget is the slug an archived landing page's URL moved to, as
// legacy's redirect_archived_language chose it: the page the archived one names,
// else the course's main page in the locale (the catalog's rule). Nil for a page
// that is not archived (moved is false). An archived page with neither has
// nowhere to go.
func (p *Pages) archivedTarget(ctx context.Context, page *ent.CourseLandingPage, locale string) (string, bool, error) {
	if page.State == nil || *page.State != stateArchived {
		return "", false, nil
	}

	var target *ent.CourseLandingPage

	if page.LandingPageToRedirectID != nil {
		named, err := p.db.CourseLandingPage.Get(ctx, *page.LandingPageToRedirectID)
		if ent.IsNotFound(err) {
			return "", false, domain.ErrNotFound
		}

		if err != nil {
			return "", false, errs.Wrapf(err, "load redirect target")
		}

		target = named
	} else {
		mains, err := p.catalog.MainPages(ctx, locale, page.CourseID)
		if err != nil {
			return "", false, err
		}

		main, ok := mains[page.CourseID]
		if !ok || main.ID == page.ID {
			return "", false, domain.ErrNotFound
		}

		target = main
	}

	return lo.FromPtr(target.Slug), true, nil
}

// LessonPage is everything the lesson player shows.
type LessonPage struct {
	Course      *ent.Course
	Lesson      *ent.CourseLesson
	Version     *ent.CourseLessonVersion
	Translation *ent.CourseLessonTranslation
	// SourceURL is the lesson's README in the content repository; empty when
	// the version carries no path.
	SourceURL string
	// Landing is the course's main landing page in the locale, which titles the
	// page; nil when it has none a visitor may open.
	Landing *ent.CourseLandingPage
	State   *progress.CourseState
	// Infos are the course's lessons in course order, as on the course page.
	Infos []*ent.CourseLessonTranslation
}

// lessonInfos is the translations of the current version's lessons in locale,
// each with its Lesson loaded, in course order.
//
// The order is the progress state's, not a second sort of our own: the state's
// positions are the one decision of where a lesson stands in the course, so a
// lesson row and its checkmark line up by construction. A lesson the state does
// not place (no position yet) follows the placed ones, oldest first.
func (p *Pages) lessonInfos(
	ctx context.Context,
	crs *ent.Course,
	locale string,
	state *progress.CourseState,
) ([]*ent.CourseLessonTranslation, error) {
	if crs.CurrentVersionID == nil {
		return nil, nil
	}

	infos, err := p.db.CourseLessonTranslation.Query().
		Where(
			courselessontranslation.LocaleEQ(locale),
			courselessontranslation.CourseVersionID(*crs.CurrentVersionID),
		).
		WithLesson().
		Order(courselessontranslation.ByID()).
		All(ctx)
	if err != nil {
		return nil, errs.Wrapf(err, "load lesson translations")
	}

	position := make(map[string]int, len(state.Lessons))
	for _, lesson := range state.Lessons {
		position[lesson.Slug] = lesson.Position
	}

	unplaced := len(position) + 1
	rank := func(info *ent.CourseLessonTranslation) int {
		if info.Edges.Lesson == nil {
			return unplaced
		}

		if pos, ok := position[lo.FromPtr(info.Edges.Lesson.Slug)]; ok {
			return pos
		}

		return unplaced
	}

	slices.SortStableFunc(infos, func(a, b *ent.CourseLessonTranslation) int { return rank(a) - rank(b) })

	return infos, nil
}

// lessonSourceURL points at the lesson's README in the content repository, so a
// reader can propose a fix to what they are reading.
//
// `path_to_code` is repository-relative and starts at the course repo, so the
// link is the org base plus that path plus the locale's README — with `modules`
// rewritten to GitHub's blob path, exactly as the legacy helper did. A lesson
// whose version carries no path has no link rather than a broken one.
func lessonSourceURL(base string, pathToCode *string, locale string) string {
	if pathToCode == nil || *pathToCode == "" {
		return ""
	}

	link, err := url.JoinPath(base, *pathToCode, locale, "README.md")
	if err != nil {
		return ""
	}

	return strings.Replace(link, path.Join("/", "modules"), "/blob/main/modules", 1)
}
