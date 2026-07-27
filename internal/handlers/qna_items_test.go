package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/api"
	"hexletbasics/internal/testsupport"
)

// Parent ids seeded by the qna fixtures.
const (
	qnaCategoryID       = 272447109 // layouting-ru, 2 items (7001, 7002)
	qnaOtherCategoryID  = 929265035 // programming-en, 1 item (7003)
	qnaLandingPageID    = 59175011  // go-ru, 2 items (8001, 8002)
	qnaOtherLandingPage = 88212034  // ruby-ru, 1 item (8003)
)

func TestAdminListCategoryQnaItems(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	items, err := h.Client.AdminListCategoryQnaItems(ctx, api.AdminListCategoryQnaItemsParams{CategoryId: qnaCategoryID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())

	// Scoped to the parent, ordered by id ascending.
	require.Len(t, items, 2)
	assert.Equal(t, "Что такое вёрстка?", items[0].Question)
	assert.Less(t, items[0].ID, items[1].ID)

	// A different category sees only its own item.
	other, err := h.Client.AdminListCategoryQnaItems(ctx, api.AdminListCategoryQnaItemsParams{CategoryId: qnaOtherCategoryID})
	require.NoError(t, err)
	assert.Len(t, other, 1)
}

func TestAdminCreateCategoryQnaItem(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	created, err := h.Client.AdminCreateCategoryQnaItem(ctx, &api.QnaItemInput{
		Question: "Новый вопрос?",
		Answer:   "Новый ответ.",
	}, api.AdminCreateCategoryQnaItemParams{CategoryId: qnaCategoryID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Новый вопрос?", created.Question)

	items, err := h.Client.AdminListCategoryQnaItems(ctx, api.AdminListCategoryQnaItemsParams{CategoryId: qnaCategoryID})
	require.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestAdminUpdateCategoryQnaItem(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	updated, err := h.Client.AdminUpdateCategoryQnaItem(ctx, &api.QnaItemInput{
		Question: "Изменённый вопрос?",
		Answer:   "Изменённый ответ.",
	}, api.AdminUpdateCategoryQnaItemParams{CategoryId: qnaCategoryID, ID: 7001})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	assert.Equal(t, "Изменённый вопрос?", updated.Question)
}

// TestAdminUpdateCategoryQnaItemWrongParent proves the nested scoping: item 7001
// belongs to qnaCategoryID, so addressing it under a different category is 404.
func TestAdminUpdateCategoryQnaItemWrongParent(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	_, err := h.Client.AdminUpdateCategoryQnaItem(ctx, &api.QnaItemInput{
		Question: "x",
		Answer:   "y",
	}, api.AdminUpdateCategoryQnaItemParams{CategoryId: qnaOtherCategoryID, ID: 7001})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())
}

func TestAdminDeleteCategoryQnaItem(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	err := h.Client.AdminDeleteCategoryQnaItem(ctx, api.AdminDeleteCategoryQnaItemParams{CategoryId: qnaCategoryID, ID: 7002})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())

	items, err := h.Client.AdminListCategoryQnaItems(ctx, api.AdminListCategoryQnaItemsParams{CategoryId: qnaCategoryID})
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestAdminListLandingPageQnaItems(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	items, err := h.Client.AdminListLandingPageQnaItems(ctx, api.AdminListLandingPageQnaItemsParams{LandingPageId: qnaLandingPageID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, h.LastStatus())
	require.Len(t, items, 2)
	assert.Less(t, items[0].ID, items[1].ID)

	other, err := h.Client.AdminListLandingPageQnaItems(ctx, api.AdminListLandingPageQnaItemsParams{LandingPageId: qnaOtherLandingPage})
	require.NoError(t, err)
	assert.Len(t, other, 1)
}

func TestAdminLandingPageQnaItemLifecycle(t *testing.T) {
	h := testsupport.NewHarness(t)
	ctx := context.Background()

	created, err := h.Client.AdminCreateLandingPageQnaItem(ctx, &api.QnaItemInput{
		Question: "Q?",
		Answer:   "A.",
	}, api.AdminCreateLandingPageQnaItemParams{LandingPageId: qnaLandingPageID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, h.LastStatus())

	updated, err := h.Client.AdminUpdateLandingPageQnaItem(ctx, &api.QnaItemInput{
		Question: "Q2?",
		Answer:   "A2.",
	}, api.AdminUpdateLandingPageQnaItemParams{LandingPageId: qnaLandingPageID, ID: created.ID})
	require.NoError(t, err)
	assert.Equal(t, "Q2?", updated.Question)

	// Wrong parent -> 404.
	_, err = h.Client.AdminUpdateLandingPageQnaItem(ctx, &api.QnaItemInput{
		Question: "x", Answer: "y",
	}, api.AdminUpdateLandingPageQnaItemParams{LandingPageId: qnaOtherLandingPage, ID: created.ID})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, h.LastStatus())

	err = h.Client.AdminDeleteLandingPageQnaItem(ctx, api.AdminDeleteLandingPageQnaItemParams{LandingPageId: qnaLandingPageID, ID: created.ID})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, h.LastStatus())
}
