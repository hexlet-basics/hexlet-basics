package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"hexletbasics/ent"
	"hexletbasics/internal/oauthserver"
	"hexletbasics/internal/progress"
)

// The learner tools (ADR-0021). getMyDashboard is the operationId of the HTTP
// operation it mirrors; the other two compose what the course and lesson pages
// read and are named for what they answer, because the contract has no single
// operation for them (it cannot tell Started from Not started per Lesson).
const (
	myDashboardTool      = "getMyDashboard"
	myCourseProgressTool = "getMyCourseProgress"
	myLessonProgressTool = "getMyLessonProgress"
)

// addLearnerTools registers the read-only learner tools. Each needs the
// progress:read Scope, counts against the personal-data rate limit and reads
// only the token's user.
//
// What is never offered, in a schema or a result: the learner's solution code
// (the platform does not even store it), the assistant chat history, the email
// and the rest of the profile. The result types below are built for the
// purpose rather than reusing a wide type and trusting a caller to ignore the
// extra fields; a test walks every schema and result for the forbidden names.
func (s *Server) addLearnerTools(server *mcp.Server) error {
	dashboardSchema, err := responseSchema("getMyDashboard")
	if err != nil {
		return errs.With("tool", myDashboardTool).Wrapf(err, "build output schema")
	}

	courseSlug := slugArg("The course's own slug, as in getCourseLesson (course.slug of getCourse's result).")
	lessonState := map[string]any{
		"type":        "string",
		"enum":        []string{progress.StateFinished, progress.StateStarted, progress.StateNotStarted},
		"description": "finished: the exercise was passed; started: begun, not passed yet; not_started: no record.",
	}

	specs := []struct {
		spec toolSpec
		run  func(ctx context.Context, c call) (result, error)
	}{
		{toolSpec{
			name:  myDashboardTool,
			title: "My courses",
			description: "Lists the signed-in learner's own courses: each Enrollment's state, " +
				"Completion and next lesson, with the course card.",
			input:  object(map[string]any{}),
			output: dashboardSchema,
		}, s.myDashboard},
		{toolSpec{
			name:  myCourseProgressTool,
			title: "My progress in a course",
			description: "Where the signed-in learner stands in one course: Enrollment state, Completion, " +
				"Next Lesson and the Lesson Progress of every lesson.",
			input: object(map[string]any{"courseSlug": courseSlug}, "courseSlug"),
			output: object(map[string]any{
				"courseSlug": map[string]any{"type": "string"},
				"state":      lessonState,
				"completion": map[string]any{"type": "integer", "minimum": 0, "maximum": 100},
				"nextLesson": map[string]any{"type": []string{"string", "null"}, "description": "Slug of the Next Lesson; null when every lesson is finished."},
				"lessons": map[string]any{"type": "array", "items": object(map[string]any{
					"slug":     map[string]any{"type": "string"},
					"position": map[string]any{"type": "integer"},
					"state":    lessonState,
				}, "slug", "position", "state")},
			}, "courseSlug", "state", "completion", "nextLesson", "lessons"),
		}, s.myCourseProgress},
		{toolSpec{
			name:  myLessonProgressTool,
			title: "My progress in a lesson",
			description: "The signed-in learner's progress in one lesson: finished, started or not started, " +
				"and passedCheck once a check passed.",
			input: object(map[string]any{
				"courseSlug": courseSlug,
				"lessonSlug": slugArg("The lesson's slug within the course."),
			}, "courseSlug", "lessonSlug"),
			output: object(map[string]any{
				"courseSlug":  map[string]any{"type": "string"},
				"lessonSlug":  map[string]any{"type": "string"},
				"position":    map[string]any{"type": "integer"},
				"state":       lessonState,
				"passedCheck": map[string]any{"type": "boolean"},
			}, "courseSlug", "lessonSlug", "position", "state", "passedCheck"),
		}, s.myLessonProgress},
	}

	for _, t := range specs {
		spec := t.spec
		spec.scope = oauthserver.ScopeProgressRead
		spec.class = ClassPersonal
		s.addTool(server, spec, t.run)
	}
	return nil
}

func (s *Server) myDashboard(ctx context.Context, c call) (result, error) {
	// Nothing is declared, so anything sent (a userId, say) is refused.
	if err := c.decode(&struct{}{}); err != nil {
		return result{}, err
	}
	dashboard, err := s.deps.MyDashboard(ctx, c.userID)
	if err != nil {
		return result{}, err
	}
	return result{
		text: fmt.Sprintf("%d courses in progress, %d finished",
			len(dashboard.Started), len(dashboard.Finished)),
		data: dashboard,
	}, nil
}

type courseArgs struct {
	CourseSlug string `json:"courseSlug"`
}

type lessonProgressArgs struct {
	CourseSlug string `json:"courseSlug"`
	LessonSlug string `json:"lessonSlug"`
}

func (s *Server) myCourseProgress(ctx context.Context, c call) (result, error) {
	var args courseArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}
	course, err := s.learnerCourse(ctx, c.userID, args.CourseSlug)
	if err != nil {
		return result{}, err
	}

	lessons := make([]map[string]any, 0, len(course.Lessons))
	for _, l := range course.Lessons {
		lessons = append(lessons, map[string]any{"slug": l.Slug, "position": l.Position, "state": l.State})
	}
	return result{
		text: fmt.Sprintf("%s: %s, %d%% complete", course.CourseSlug, course.State, course.Completion),
		data: map[string]any{
			"courseSlug": course.CourseSlug,
			"state":      course.State,
			"completion": course.Completion,
			"nextLesson": course.NextLessonSlug,
			"lessons":    lessons,
		},
	}, nil
}

func (s *Server) myLessonProgress(ctx context.Context, c call) (result, error) {
	var args lessonProgressArgs
	if err := c.decode(&args); err != nil {
		return result{}, err
	}
	if args.LessonSlug == "" {
		return result{}, userError{errors.New("invalid arguments: lessonSlug is required")}
	}
	course, err := s.learnerCourse(ctx, c.userID, args.CourseSlug)
	if err != nil {
		return result{}, err
	}
	for _, l := range course.Lessons {
		if l.Slug != args.LessonSlug {
			continue
		}
		return result{
			text: fmt.Sprintf("%s/%s: %s", course.CourseSlug, l.Slug, l.State),
			data: map[string]any{
				"courseSlug":  course.CourseSlug,
				"lessonSlug":  l.Slug,
				"position":    l.Position,
				"state":       l.State,
				"passedCheck": l.State == progress.StateFinished,
			},
		}, nil
	}
	return result{}, errNotFound("lesson " + strconv.Quote(args.LessonSlug))
}

// learnerCourse reads the token's user's position in the course. An unknown
// course is an answer the agent can act on, not a failure of ours.
func (s *Server) learnerCourse(ctx context.Context, userID int, slug string) (*progress.LearnerCourse, error) {
	if slug == "" {
		return nil, userError{errors.New("invalid arguments: courseSlug is required")}
	}
	course, err := s.deps.MyCourseProgress(ctx, userID, slug)
	if ent.IsNotFound(err) {
		return nil, errNotFound("course " + strconv.Quote(slug))
	}
	if err != nil {
		return nil, err
	}
	return course, nil
}
