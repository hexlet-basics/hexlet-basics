package blog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/blog"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/testsupport"
)

func TestEditorCreateStampsTheCreatorAndTheLocale(t *testing.T) {
	editor := blog.NewEditor(testsupport.NewClient(t))
	name, slug := "Заметки", "zametki"

	post, err := editor.Create(t.Context(), fixture.ID(fixture.UserAdminLearner), "ru",
		blog.PostFields{Name: &name, Slug: &slug, RichBody: "<p>Текст</p>"})

	require.NoError(t, err)
	assert.Equal(t, fixture.ID(fixture.UserAdminLearner), post.CreatorID)
	require.NotNil(t, post.Locale)
	assert.Equal(t, "ru", *post.Locale)
	assert.Equal(t, "<p>Текст</p>", post.RichBody)
}

func TestEditorUpdateClearsWhatIsLeftOutAndKeepsTheLocale(t *testing.T) {
	db := testsupport.NewClient(t)
	editor := blog.NewEditor(db)
	id := fixture.ID(fixture.BlogPostHelloWorld)
	before := db.BlogPost.GetX(t.Context(), id)
	name := "Renamed"

	after, err := editor.Update(t.Context(), id, blog.PostFields{Name: &name, RichBody: "<p>New</p>"})

	require.NoError(t, err)
	require.NotNil(t, after.Name)
	assert.Equal(t, "Renamed", *after.Name)
	assert.Nil(t, after.Slug, "a nil field clears the column")
	assert.Nil(t, after.Description)
	assert.Equal(t, before.Locale, after.Locale)
}

func TestEditorUpdateOfAMissingPostIsNotFound(t *testing.T) {
	editor := blog.NewEditor(testsupport.NewClient(t))

	_, err := editor.Update(t.Context(), 0, blog.PostFields{RichBody: "<p>x</p>"})

	require.ErrorIs(t, err, domain.ErrNotFound)
}
