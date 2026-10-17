package fixtures_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"hexletbasics/internal/fixtures"
)

// The legacy Rails fixtures derived every id from its label the same way, so a
// row exported from them keeps its id when it is given the label it had there.
func TestIDMatchesTheLegacyRailsFixtureIDs(t *testing.T) {
	assert.Equal(t, 82481401, fixtures.ID("javascript"))
	assert.Equal(t, 207281424, fixtures.ID("ruby"))
	assert.Equal(t, 912823128, fixtures.ID("go"))
}

func TestIDDiffersPerLabel(t *testing.T) {
	assert.NotEqual(t, fixtures.ID("user_removed"), fixtures.ID("user_learner"))
}

func TestPasswordIsTheCleartextOfTheFixtureDigest(t *testing.T) {
	digest, err := fixtures.PasswordDigest()
	assert.NoError(t, err)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(digest), []byte(fixtures.Password)))
}
