// Package fixtures loads the committed starting data in fixtures/ and names its
// rows. A row is written with a label, Rails-style: the YAML says
// `{{ id "user_removed" }}` where the row's primary key (or a reference to
// another row) goes, and a Go test asks for the same row with
// fixtures.ID(fixtures.UserRemoved). The label is the only handle a test needs,
// so no test depends on a literal id or on how many rows a table happens to
// hold.
//
// Every label has a Label variable of its own in labels_gen.go, generated from
// the YAML by `mise run gen:fixtures`, and ID takes a Label rather than a
// string. A typo, or a label whose row was removed, is a compile error in the
// test instead of a panic when it runs.
//
// Ids follow the order the rows are written in. The application reads
// "oldest first" and "newest first" off the id, so a fixture that is written
// before another is older than it, and a test of such an order can say so by
// writing the rows in that order. A label that no row defines is an error, in
// the YAML and in a test alike, so a typo cannot pass quietly.
//
// One loader serves every consumer (the Go test containers, `mise run
// db:fixtures` and the end-to-end database), because the template functions
// below are what make the YAML valid; the stock testfixtures CLI cannot
// evaluate them.
package fixtures

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"text/template"

	"github.com/go-testfixtures/testfixtures/v3"
	"golang.org/x/crypto/bcrypt"
)

// Password is the cleartext behind every `{{ password }}` digest: a fixture
// user with a password signs in with it.
const Password = "password"

// firstID is the id of the first fixture row. Ids count up from it across all
// files, so no two rows share one (a reference that points at the wrong table
// fails instead of matching by accident) and every id stays below the 10000
// the loader resets the sequences to.
const firstID = 1001

// definition is a row's primary key line: `- id: {{ id "label" }}`, written
// exactly so, first in the row. A row spelled any other way is not numbered and
// the template then fails on the first reference to it.
var definition = regexp.MustCompile(`^- id: \{\{\s*id\s+"([^"]+)"\s*\}\}\s*$`)

// labels maps every label to its id.
type labels map[string]int

// scan numbers the rows of every file in dir, files in name order and rows in
// the order they are written.
func scan(dir string) (labels, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, fmt.Errorf("list fixture files: %w", err)
	}
	sort.Strings(files)

	found := labels{}
	next := firstID
	for _, file := range files {
		if err := scanFile(file, found, &next); err != nil {
			return nil, err
		}
	}
	return found, nil
}

// scanFile numbers the rows of one file, continuing from next.
func scanFile(path string, found labels, next *int) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open fixture file: %w", err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		m := definition.FindStringSubmatch(scanner.Text())
		if m == nil {
			continue
		}
		if _, dup := found[m[1]]; dup {
			return fmt.Errorf("fixture label %q is defined twice (again in %s)", m[1], filepath.Base(path))
		}
		found[m[1]] = *next
		*next++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read fixture file %s: %w", filepath.Base(path), err)
	}
	return nil
}

func (l labels) id(label string) (int, error) {
	id, ok := l[label]
	if !ok {
		return 0, fmt.Errorf("fixture label %q is not defined by any row", label)
	}
	return id, nil
}

var (
	repoOnce   sync.Once
	repoLabels labels
	repoErr    error
)

// Label names one fixture row. It is a struct, not a string type, on purpose:
// an untyped string constant converts to a string type implicitly, so
// ID("typo") would compile, while nothing converts to a struct with an
// unexported field. The only ways to get a Label are the generated variables
// and LabelOf.
type Label struct{ name string }

// LabelOf is the escape hatch for a label that is not known when the test is
// written, such as one read from a request path. ID still panics for a name no
// row defines, so prefer the generated variables wherever the row is known.
func LabelOf(name string) Label {
	return Label{name: name}
}

// String is the label as it is written in the YAML.
func (l Label) String() string { return l.name }

// LabelNames lists every label the rows in dir define, in name order. It is
// what the generator of labels_gen.go reads.
func LabelNames(dir string) ([]string, error) {
	found, err := scan(dir)
	if err != nil {
		return nil, fmt.Errorf("number fixture rows: %w", err)
	}
	names := make([]string, 0, len(found))
	for name := range found {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// ID is the primary key of the fixture row labelled label in the repository's
// fixtures/ directory. It panics for a label no row defines.
func ID(label Label) int {
	repoOnce.Do(func() {
		var dir string
		if dir, repoErr = Dir(); repoErr == nil {
			repoLabels, repoErr = scan(dir)
		}
	})
	if repoErr != nil {
		panic(fmt.Sprintf("fixtures: %v", repoErr))
	}
	id, err := repoLabels.id(label.name)
	if err != nil {
		panic(fmt.Sprintf("fixtures: %v", err))
	}
	return id
}

// ID32 is ID as the int32 the HTTP contract uses for identifiers.
func ID32(label Label) int32 {
	return int32(ID(label))
}

// Dir finds the repository's fixtures/ directory from the working directory
// upwards, which is where a package's tests run from.
func Dir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "fixtures"), nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("inspect repository root: %w", err)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("find the repository root: go.mod not found")
		}
		dir = parent
	}
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

// Load replaces the contents of the tables in dir with the fixture rows. It
// refuses a database whose name does not contain "test" unless
// allowAnyDatabase is set, because loading wipes the fixture tables first.
func Load(db *sql.DB, dir string, allowAnyDatabase bool) error {
	numbered, err := scan(dir)
	if err != nil {
		return fmt.Errorf("number fixture rows: %w", err)
	}
	// The template options come before Directory: it reads and renders the files
	// while the options are applied, in order.
	options := []func(*testfixtures.Loader) error{
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateFuncs(template.FuncMap{
			"id":       numbered.id,
			"password": PasswordDigest,
		}),
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
