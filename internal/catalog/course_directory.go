package catalog

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/domain"
)

// CourseDirectory is the administrators' read side of the courses.
type CourseDirectory struct {
	db *ent.Client
}

// NewCourseDirectory builds the directory over an ent client.
func NewCourseDirectory(db *ent.Client) *CourseDirectory {
	return &CourseDirectory{db: db}
}

// Option is a course as a select offers it.
type Option struct {
	ID   int
	Name string
}

// Options serves selects: every course as id and name, with no page, so a
// picker cannot silently lose the courses past a page boundary. The order is
// oldest first, a stable order for a list people scan by name; a course without
// a name is listed by its slug, the one label every course has.
func (d *CourseDirectory) Options(ctx context.Context) ([]Option, error) {
	rows, err := d.db.Course.Query().
		Select(course.FieldID, course.FieldSlug, course.FieldName).
		Order(ent.Asc(course.FieldID)).
		All(ctx)
	if err != nil {
		return nil, domain.FromStore(err)
	}

	options := make([]Option, 0, len(rows))
	for _, row := range rows {
		name := ""
		if row.Slug != nil {
			name = *row.Slug
		}

		if row.Name != nil && *row.Name != "" {
			name = *row.Name
		}

		options = append(options, Option{ID: row.ID, Name: name})
	}

	return options, nil
}
