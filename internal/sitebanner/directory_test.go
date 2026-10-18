package sitebanner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/sitebanner"
	"hexletbasics/internal/testsupport"
)

func TestActiveIsAPublishedBannerOfTheLocaleWhoseWindowIsOpen(t *testing.T) {
	directory := sitebanner.NewDirectory(testsupport.NewClient(t))

	row, err := directory.Active(t.Context(), "ru", testsupport.Today)

	require.NoError(t, err)
	assert.Equal(t, "published", row.State)
	assert.Equal(t, "ru", row.Locale)

	if row.StartsAt != nil {
		assert.False(t, row.StartsAt.After(testsupport.Today))
	}

	if row.FinishesAt != nil {
		assert.False(t, row.FinishesAt.Before(testsupport.Today))
	}
}

func TestActiveWithNoBannerInTheLocaleIsNotFound(t *testing.T) {
	directory := sitebanner.NewDirectory(testsupport.NewClient(t))

	// The fixtures hold banners in ru and en only: another locale shows nothing.
	_, err := directory.Active(t.Context(), "es", testsupport.Today)

	require.ErrorIs(t, err, domain.ErrNotFound)
}
