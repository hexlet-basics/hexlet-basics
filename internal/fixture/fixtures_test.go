package fixture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func write(t *testing.T, dir, name, body string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600))
}

// Ids follow the order rows are written in, across files in name order, so a
// test of "oldest first" can set the order by writing the rows in it.
func TestScanNumbersRowsInTheOrderTheyAreWritten(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	write(t, dir, "a_things.yml", "- id: {{ id \"thing_first\" }}\n  name: x\n- id: {{ id \"thing_second\" }}\n")
	write(t, dir, "b_others.yml", "- id: {{ id \"other_one\" }}\n  thing_id: {{ id \"thing_first\" }}\n")

	got, err := scan(dir)

	require.NoError(t, err)
	assert.Equal(t, labels{"thing_first": firstID, "thing_second": firstID + 1, "other_one": firstID + 2}, got)
}

func TestScanRefusesALabelDefinedTwice(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	write(t, dir, "a.yml", "- id: {{ id \"same\" }}\n- id: {{ id \"same\" }}\n")

	_, err := scan(dir)

	assert.ErrorContains(t, err, `"same" is defined twice`)
}

func TestAReferenceToALabelNoRowDefinesIsAnError(t *testing.T) {
	t.Parallel()
	_, err := labels{"known": firstID}.id("unknown")

	assert.ErrorContains(t, err, `"unknown" is not defined`)
}

func TestIDResolvesALabelOfTheRepositoryFixtures(t *testing.T) {
	t.Parallel()
	assert.NotEqual(t, ID(UserAdminLearner), ID(UserBlank))
	assert.Panics(t, func() { ID(LabelOf("no_such_label")) })
}

func TestPasswordIsTheCleartextOfTheFixtureDigest(t *testing.T) {
	t.Parallel()
	digest, err := PasswordDigest()
	require.NoError(t, err)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(digest), []byte(Password)))
}
