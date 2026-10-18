// Package qnaitem edits the questions and answers shown under a course
// category or a course landing page.
//
// An item belongs to exactly one parent, and every operation is scoped to it: an
// item looked up under a different parent is not found, as legacy's
// `parent.qna_items.find(id)` behaved, so an administrator editing one parent
// can never touch another's.
package qnaitem

import (
	"context"

	"hexletbasics/ent"
	"hexletbasics/ent/categoryqnaitem"
	"hexletbasics/ent/landingpageqnaitem"
	"hexletbasics/internal/domain"
)

// Fields is what an administrator writes on an item.
type Fields struct {
	Question string
	Answer   string
}

// Category edits the items of a course category.
type Category struct {
	db *ent.Client
}

// NewCategory builds the editor over an ent client.
func NewCategory(db *ent.Client) *Category {
	return &Category{db: db}
}

// List returns a category's items, oldest first (legacy `qna_items.order(:id)`).
func (c *Category) List(ctx context.Context, categoryID int) ([]*ent.CategoryQnaItem, error) {
	rows, err := c.db.CategoryQnaItem.Query().
		Where(categoryqnaitem.CourseCategoryID(categoryID)).
		Order(ent.Asc(categoryqnaitem.FieldID)).
		All(ctx)

	return rows, domain.FromStore(err)
}

// Create adds an item to the category.
func (c *Category) Create(ctx context.Context, categoryID int, fields Fields) (*ent.CategoryQnaItem, error) {
	row, err := c.db.CategoryQnaItem.Create().
		SetCourseCategoryID(categoryID).
		SetQuestion(fields.Question).
		SetAnswer(fields.Answer).
		Save(ctx)

	return row, domain.FromStore(err)
}

// Update rewrites an item of the category; an id that is missing or under
// another category is domain.ErrNotFound, before any write.
func (c *Category) Update(ctx context.Context, categoryID, id int, fields Fields) (*ent.CategoryQnaItem, error) {
	row, err := c.find(ctx, categoryID, id)
	if err != nil {
		return nil, err
	}

	updated, err := row.Update().SetQuestion(fields.Question).SetAnswer(fields.Answer).Save(ctx)

	return updated, domain.FromStore(err)
}

// Delete removes an item of the category, with the same scoping as Update.
func (c *Category) Delete(ctx context.Context, categoryID, id int) error {
	row, err := c.find(ctx, categoryID, id)
	if err != nil {
		return err
	}

	return domain.FromStore(c.db.CategoryQnaItem.DeleteOne(row).Exec(ctx))
}

func (c *Category) find(ctx context.Context, categoryID, id int) (*ent.CategoryQnaItem, error) {
	row, err := c.db.CategoryQnaItem.Query().
		Where(categoryqnaitem.ID(id), categoryqnaitem.CourseCategoryID(categoryID)).
		Only(ctx)

	return row, domain.FromStore(err)
}

// LandingPage edits the items of a course landing page.
type LandingPage struct {
	db *ent.Client
}

// NewLandingPage builds the editor over an ent client.
func NewLandingPage(db *ent.Client) *LandingPage {
	return &LandingPage{db: db}
}

// List returns a landing page's items, oldest first.
func (l *LandingPage) List(ctx context.Context, landingPageID int) ([]*ent.LandingPageQnaItem, error) {
	rows, err := l.db.LandingPageQnaItem.Query().
		Where(landingpageqnaitem.CourseLandingPageID(landingPageID)).
		Order(ent.Asc(landingpageqnaitem.FieldID)).
		All(ctx)

	return rows, domain.FromStore(err)
}

// Create adds an item to the landing page.
func (l *LandingPage) Create(ctx context.Context, landingPageID int, fields Fields) (*ent.LandingPageQnaItem, error) {
	row, err := l.db.LandingPageQnaItem.Create().
		SetCourseLandingPageID(landingPageID).
		SetQuestion(fields.Question).
		SetAnswer(fields.Answer).
		Save(ctx)

	return row, domain.FromStore(err)
}

// Update rewrites an item of the landing page; an id that is missing or under
// another landing page is domain.ErrNotFound, before any write.
func (l *LandingPage) Update(ctx context.Context, landingPageID, id int, fields Fields) (*ent.LandingPageQnaItem, error) {
	row, err := l.find(ctx, landingPageID, id)
	if err != nil {
		return nil, err
	}

	updated, err := row.Update().SetQuestion(fields.Question).SetAnswer(fields.Answer).Save(ctx)

	return updated, domain.FromStore(err)
}

// Delete removes an item of the landing page, with the same scoping as Update.
func (l *LandingPage) Delete(ctx context.Context, landingPageID, id int) error {
	row, err := l.find(ctx, landingPageID, id)
	if err != nil {
		return err
	}

	return domain.FromStore(l.db.LandingPageQnaItem.DeleteOne(row).Exec(ctx))
}

func (l *LandingPage) find(ctx context.Context, landingPageID, id int) (*ent.LandingPageQnaItem, error) {
	row, err := l.db.LandingPageQnaItem.Query().
		Where(landingpageqnaitem.ID(id), landingpageqnaitem.CourseLandingPageID(landingPageID)).
		Only(ctx)

	return row, domain.FromStore(err)
}
