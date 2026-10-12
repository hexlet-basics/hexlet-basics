// Package curriculum reads a Course Version's learning program: its modules
// with their lessons in course order. The course page and the Yandex feed both
// render exactly this tree, so the eager-load shape and its ordering live once
// here rather than as a copy in each reader.
package curriculum

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/ent/courselessonversion"
	"hexletbasics/ent/coursemoduletranslation"
	"hexletbasics/ent/coursemoduleversion"
)

// Modules loads the module translations of the given Course Versions in the
// locale, each with its module version's lesson versions and their lessons.
//
// Modules are ordered by their version row's order, as the loader built them,
// with the translation id breaking ties. A module's lessons are ordered by
// natural order (unnumbered last, id breaking ties), the order a flat lesson
// list uses, so a module's slugs read in the same order the list shows them.
// Callers group the rows by CourseVersionID when they asked for several.
func Modules(ctx context.Context, db *ent.Client, locale string, versionIDs ...int) ([]*ent.CourseModuleTranslation, error) {
	if len(versionIDs) == 0 {
		return nil, nil
	}

	rows, err := db.CourseModuleTranslation.Query().
		Where(
			coursemoduletranslation.LocaleEQ(locale),
			coursemoduletranslation.CourseVersionIDIn(versionIDs...),
		).
		WithVersion(func(q *ent.CourseModuleVersionQuery) {
			q.WithLessonVersions(func(q *ent.CourseLessonVersionQuery) {
				q.WithLesson().Order(
					courselessonversion.ByNaturalOrder(sql.OrderNullsLast()),
					courselessonversion.ByID(),
				)
			})
		}).
		Order(
			coursemoduletranslation.ByVersionField(coursemoduleversion.FieldOrder, sql.OrderNullsLast()),
			// Qualify the id: the edge ordering joins language_module_versions,
			// so a bare `id` is ambiguous.
			func(s *sql.Selector) { s.OrderBy(s.C(coursemoduletranslation.FieldID)) },
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("load modules of versions %v: %w", versionIDs, err)
	}
	return rows, nil
}
