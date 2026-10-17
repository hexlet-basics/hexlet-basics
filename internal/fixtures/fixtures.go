// Package fixtures loads the committed starting data in fixtures/ and names its
// rows. A row is written with a label, Rails-style: the YAML says
// `{{ id "user_removed" }}` where the row's primary key (or a reference to
// another row) goes, and a Go test asks for the same row with
// fixtures.ID("user_removed"). The label is the only handle a test needs, so
// no test depends on a literal id or on how many rows a table happens to hold.
//
// One loader serves every consumer (the Go test containers, `mise run
// db:fixtures` and the end-to-end database), because the template functions
// below are what make the YAML valid; the stock testfixtures CLI cannot
// evaluate them.
package fixtures

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"sync"
	"text/template"

	"github.com/go-testfixtures/testfixtures/v3"
	"golang.org/x/crypto/bcrypt"
)

// Password is the cleartext behind every `{{ password }}` digest: a fixture
// user with a password signs in with it.
const Password = "password"

// maxID is Rails' ActiveRecord::FixtureSet::MAX_ID. Using the same function as
// the legacy Rails fixtures keeps every id that was exported from them
// (`javascript` is 82481401, as it was) stable when its row gets a label.
const maxID = 1<<30 - 1

// ID is the primary key of the fixture row labelled label. It is the same in
// every table, so a reference to a row is written with the row's label.
func ID(label string) int {
	return int(crc32.ChecksumIEEE([]byte(label)) % maxID)
}

var (
	digestOnce sync.Once
	digest     string
	digestErr  error
)

// PasswordDigest is the bcrypt digest of Password at the lowest cost, computed
// once: every fixture user that has a password shares it, and a cheap digest
// keeps sign-in tests fast.
func PasswordDigest() (string, error) {
	digestOnce.Do(func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.MinCost)
		digest, digestErr = string(hash), err
	})
	return digest, digestErr
}

func funcs() template.FuncMap {
	return template.FuncMap{
		"id":       ID,
		"password": PasswordDigest,
	}
}

// Load replaces the contents of the tables in dir with the fixture rows. It
// refuses a database whose name does not contain "test" unless
// allowAnyDatabase is set, because loading wipes the fixture tables first.
func Load(db *sql.DB, dir string, allowAnyDatabase bool) error {
	// The template options come before Directory: it reads and renders the files
	// while the options are applied, in order.
	options := []func(*testfixtures.Loader) error{
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateFuncs(funcs()),
		testfixtures.Directory(dir),
		testfixtures.SkipTableChecksumComputation(),
	}
	if allowAnyDatabase {
		options = append(options, testfixtures.DangerousSkipTestDatabaseCheck())
	}
	loader, err := testfixtures.New(options...)
	if err != nil {
		return fmt.Errorf("create fixture loader: %w", err)
	}
	if err := loader.Load(); err != nil {
		return fmt.Errorf("load fixtures: %w", err)
	}
	return nil
}
