package qnaitem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixture"
	"hexletbasics/internal/qnaitem"
	"hexletbasics/internal/testsupport"
)

func TestAnItemUnderAnotherCategoryIsNotFound(t *testing.T) {
	db := testsupport.NewClient(t)
	items := qnaitem.NewCategory(db)
	item := db.CategoryQnaItem.GetX(t.Context(), fixture.ID(fixture.CategoryQnaLayoutingRu))
	other := db.CategoryQnaItem.GetX(t.Context(), fixture.ID(fixture.CategoryQnaProgrammingEn)).CourseCategoryID
	require.NotEqual(t, item.CourseCategoryID, other, "the fixtures hold items of two categories")
	fields := qnaitem.Fields{Question: "Q?", Answer: "A."}

	_, err := items.Update(t.Context(), other, item.ID, fields)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.ErrorIs(t, items.Delete(t.Context(), other, item.ID), domain.ErrNotFound)

	untouched := db.CategoryQnaItem.GetX(t.Context(), item.ID)
	assert.Equal(t, item.Question, untouched.Question, "a refused update writes nothing")
}

func TestAnItemOfTheCategoryIsRewrittenAndRemoved(t *testing.T) {
	db := testsupport.NewClient(t)
	items := qnaitem.NewCategory(db)
	item := db.CategoryQnaItem.GetX(t.Context(), fixture.ID(fixture.CategoryQnaLayoutingRu))

	updated, err := items.Update(t.Context(), item.CourseCategoryID, item.ID, qnaitem.Fields{Question: "New?", Answer: "Yes."})
	require.NoError(t, err)
	require.NotNil(t, updated.Question)
	assert.Equal(t, "New?", *updated.Question)

	require.NoError(t, items.Delete(t.Context(), item.CourseCategoryID, item.ID))
	rest, err := items.List(t.Context(), item.CourseCategoryID)
	require.NoError(t, err)
	for _, left := range rest {
		assert.NotEqual(t, item.ID, left.ID)
	}
}

func TestAListIsOldestFirstAndOnlyTheParents(t *testing.T) {
	db := testsupport.NewClient(t)
	items := qnaitem.NewLandingPage(db)
	first := db.LandingPageQnaItem.GetX(t.Context(), fixture.ID(fixture.LandingPageQnaGoRu))

	listed, err := items.List(t.Context(), first.CourseLandingPageID)
	require.NoError(t, err)

	require.NotEmpty(t, listed)
	for i, item := range listed {
		assert.Equal(t, first.CourseLandingPageID, item.CourseLandingPageID)
		if i > 0 {
			assert.Greater(t, item.ID, listed[i-1].ID)
		}
	}
}
