// Package courses reads the courses for the administrators' pickers.
package courses

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/course"
	"hexletbasics/internal/domain"
)

// Directory is the administrators' read side of the courses.
type Directory struct {
	db *ent.Client
}

// NewDirectory builds the directory over an ent client.
func NewDirectory(db *ent.Client) *Directory {
	return &Directory{db: db}
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
func (d *Directory) Options(ctx context.Context) ([]Option, error) {
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
