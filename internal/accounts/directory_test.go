package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/accounts"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/fixtures"
	"hexletbasics/internal/testsupport"
)

func TestDirectoryListsNewestFirstAndPagesDoNotOverlap(t *testing.T) {
	directory := accounts.NewDirectory(testsupport.NewClient(t))

	first, total, err := directory.List(t.Context(), "", domain.Page{Offset: 0, Limit: 2})
	require.NoError(t, err)
	second, _, err := directory.List(t.Context(), "", domain.Page{Offset: 2, Limit: 2})
	require.NoError(t, err)

	require.Len(t, first, 2)
	require.NotEmpty(t, second)
	assert.Greater(t, total, 2)
	assert.Greater(t, first[0].ID, first[1].ID, "newest first")
	assert.Greater(t, first[1].ID, second[0].ID, "the next page starts where this one ended")
}

func TestDirectorySearchMatchesNameOrEmailInAnyCase(t *testing.T) {
	directory := accounts.NewDirectory(testsupport.NewClient(t))

	byEmail, err := directory.Search(t.Context(), "SECOND-ADMIN@")
	require.NoError(t, err)
	byName, err := directory.Search(t.Context(), "second")
	require.NoError(t, err)

	require.Len(t, byEmail, 1)
	assert.Equal(t, fixtures.ID(fixtures.UserSecondAdmin), byEmail[0].ID)
	assert.NotEmpty(t, byName)
}

func TestDirectorySearchOfABlankTermIsEmpty(t *testing.T) {
	directory := accounts.NewDirectory(testsupport.NewClient(t))

	found, err := directory.Search(t.Context(), "")

	require.NoError(t, err)
	assert.Empty(t, found)
}

func TestDirectoryGetOfAMissingUserIsNotFound(t *testing.T) {
	directory := accounts.NewDirectory(testsupport.NewClient(t))

	_, err := directory.Get(t.Context(), 0)

	require.ErrorIs(t, err, domain.ErrNotFound)
}
