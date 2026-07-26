package handlers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"hexletbasics/internal/api"
)

// Server implements the generated api.StrictServerInterface backed by Postgres.
type Server struct {
	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{db: db}
}

// ListCourses returns the published course catalog.
//
// URL stays `/languages` for backward-compat; the domain concept is Course.
// Mirrors the legacy Ruby scope: listed landing pages joined to their language,
// ordered by the language display order.
func (s *Server) ListCourses(ctx context.Context, _ api.ListCoursesRequestObject) (api.ListCoursesResponseObject, error) {
	const query = `
		SELECT
			lp.id, lp.slug, lp.header, lp.name, lp.locale,
			l.id, l.slug, l.name, l.learn_as, l.progress,
			l.members_count, l.lessons_count, l.category_id
		FROM language_landing_pages lp
		JOIN languages l ON l.id = lp.language_id
		WHERE lp.listed = true
		ORDER BY l."order" NULLS LAST, l.id
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []api.CourseCatalogItem{}
	for rows.Next() {
		var (
			item    api.CourseCatalogItem
			course  api.Course
			learnAs *string
			progr   *string
		)
		if err := rows.Scan(
			&item.Id, &item.Slug, &item.Header, &item.Name, &item.Locale,
			&course.Id, &course.Slug, &course.Name, &learnAs, &progr,
			&course.MembersCount, &course.LessonsCount, &course.CategoryId,
		); err != nil {
			return nil, err
		}
		if learnAs != nil {
			v := api.CourseLearnAs(*learnAs)
			course.LearnAs = &v
		}
		if progr != nil {
			v := api.CourseProgress(*progr)
			course.Progress = &v
		}
		item.MembersCount = course.MembersCount
		// Legacy: duration = lessons_count * 15 minutes, rendered in hours.
		item.Duration = course.LessonsCount * 15 / 60
		// coverUrl stays nil until course cover assets are re-uploaded.
		item.Course = course
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return api.ListCourses200JSONResponse(items), nil
}
