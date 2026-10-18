package relatedcourse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hexletbasics/internal/domain"
	"hexletbasics/internal/jobs"
	"hexletbasics/internal/relatedcourse"
	"hexletbasics/internal/testsupport"
)

func TestSuggestSchedulesThePickForAPost(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)

	require.NoError(t, relatedcourse.NewScheduler(db, services.Jobs).Suggest(t.Context(), helloWorldPost()))

	assert.Equal(t, helloWorldPost(), testsupport.RequireEnqueued(t, services, jobs.SuggestRelatedCoursesArgs{}).BlogPostID)
}

func TestSuggestForAMissingPostIsNotFoundAndSchedulesNothing(t *testing.T) {
	db, tx, transactor := testsupport.NewClientWithTxAndTransactor(t)
	services := testsupport.NewServices(t, tx, transactor)

	err := relatedcourse.NewScheduler(db, services.Jobs).Suggest(t.Context(), 0)

	require.ErrorIs(t, err, domain.ErrNotFound)
	testsupport.RequireNotEnqueued(t, services, jobs.SuggestRelatedCoursesArgs{})
}
