// Package courseloader parses an exercises-<slug> course repository into a plain
// data tree and drives the build of a new course version from it. This file is
// the PURE parser: given a checked-out repo directory (and the course slug), it
// walks the on-disk layout and returns structs, touching neither the database
// nor the network. That purity is what makes it golden-testable against the
// committed fixture course.
//
// The layout it reads (a straight port of the legacy Rails ExerciseLoader, which
// is the source of truth for the format):
//
//	spec.yml                                      course/build metadata (`language:` map)
//	modules/<order>-<slug>/
//	  description.<locale>.yml                    module name/description per locale
//	  <order>-<slug>/                             a lesson
//	    <exercise_filename>                       student starter code   -> original/prepared
//	    <exercise_test_filename>                  the test               -> test_code
//	    <locale>/                                 one dir per translated locale
//	      data.yml                                name, tips[], definitions[]
//	      README.md                               lesson theory
//	      EXERCISE.md                             lesson instructions
package courseloader

import (
	"bytes"
	"cmp"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

// Course is the fully parsed repository: the spec plus modules in build order,
// each with its lessons in build order. natural_order is already assigned across
// the whole tree.
type Course struct {
	Spec    Spec
	Modules []Module
}

// Spec is the root spec.yml `language:` map. Every field is required (legacy
// fetches each, erroring on a missing key), so parsing fails loudly rather than
// silently building a half-configured version.
type Spec struct {
	Name                 string `yaml:"name"`
	Progress             string `yaml:"progress"`
	LearnAs              string `yaml:"learn_as"`
	DockerImage          string `yaml:"docker_image"`
	Extension            string `yaml:"extension"`
	ExerciseFilename     string `yaml:"exercise_filename"`
	ExerciseTestFilename string `yaml:"exercise_test_filename"`
}

// Module is one module directory: its numeric order, slug, per-locale infos, and
// ordered lessons.
type Module struct {
	Order   int
	Slug    string
	Infos   []ModuleInfo
	Lessons []Lesson
}

// ModuleInfo is one locale's name/description for a module.
type ModuleInfo struct {
	Locale      string
	Name        string
	Description string
}

// Lesson is one lesson directory: ordering, the three code artifacts, the
// container path used by the runtime runner, and per-locale infos.
type Lesson struct {
	Order        int
	Slug         string
	NaturalOrder int
	TestCode     string
	OriginalCode string
	PreparedCode string
	PathToCode   string
	Infos        []LessonInfo
}

// LessonInfo is one locale's content for a lesson. Dir is the absolute path to
// the locale directory, kept so the loader can resolve README image references
// (`./assets/foo.png`) relative to it.
type LessonInfo struct {
	Locale       string
	Dir          string
	Name         string
	Description  string
	Theory       string
	Instructions string
	Tips         []string
	Definitions  []Definition
}

// Definition is one term/explanation pair from a lesson's data.yml `definitions`.
type Definition struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Parse walks the repo rooted at dir and returns the parsed course. courseSlug is
// the DB course's slug (not derivable from the repo alone) and is only used to
// build each lesson's path_to_code, matching legacy exactly.
func Parse(dir, courseSlug string) (*Course, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, oops.Wrapf(err, "open course repo %s", dir)
	}
	defer func() { _ = root.Close() }()

	repo := repository{fsys: root.FS(), dir: dir}

	spec, err := parseSpec(repo)
	if err != nil {
		return nil, err
	}

	modules, err := parseModules(repo, "modules", spec, courseSlug)
	if err != nil {
		return nil, err
	}

	assignNaturalOrder(modules)

	return &Course{Spec: *spec, Modules: modules}, nil
}

// repository is a course repo opened for reading. Everything is read through
// fsys, which is rooted at the repo (os.Root): a symlink in a cloned repo that
// points outside of it — spec.yml -> /etc/passwd — fails to open instead of
// putting a host file into the database. Paths handed to it are slash-separated
// and relative to the repo; dir only turns them back into the absolute paths
// the loader needs for image uploads.
type repository struct {
	fsys fs.FS
	dir  string
}

// abs is the OS path of a repo-relative path.
func (r repository) abs(rel string) string {
	return filepath.Join(r.dir, filepath.FromSlash(rel))
}

func parseSpec(repo repository) (*Spec, error) {
	raw, err := fs.ReadFile(repo.fsys, "spec.yml")
	if err != nil {
		return nil, oops.Wrapf(err, "read %s", repo.abs("spec.yml"))
	}

	var doc struct {
		Course Spec `yaml:"language"`
	}
	if err := decodeYAML(raw, &doc); err != nil {
		return nil, oops.Wrapf(err, "parse %s", repo.abs("spec.yml"))
	}

	s := doc.Course
	// Mirror legacy's fetch("...") — a missing key is a hard error, not a blank
	// build. Every missing key is reported, in file order, so one run shows the
	// whole fix.
	var missing []error

	for _, key := range []struct{ name, value string }{
		{"name", s.Name},
		{"progress", s.Progress},
		{"learn_as", s.LearnAs},
		{"docker_image", s.DockerImage},
		{"extension", s.Extension},
		{"exercise_filename", s.ExerciseFilename},
		{"exercise_test_filename", s.ExerciseTestFilename},
	} {
		if key.value == "" {
			missing = append(missing, oops.Errorf("spec.yml: language.%s is required", key.name))
		}
	}

	if err := errors.Join(missing...); err != nil {
		return nil, err
	}

	return &s, nil
}

func parseModules(repo repository, modulesDir string, spec *Spec, courseSlug string) ([]Module, error) {
	dirs, err := subdirs(repo, modulesDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read modules dir %s", repo.abs(modulesDir))
	}

	modules := make([]Module, 0, len(dirs))
	for _, name := range dirs {
		order, slug, err := splitOrderSlug(name)
		if err != nil {
			return nil, oops.Wrapf(err, "module dir %q", name)
		}

		moduleDir := path.Join(modulesDir, name)

		infos, err := parseModuleInfos(repo, moduleDir)
		if err != nil {
			return nil, err
		}

		if len(infos) == 0 {
			return nil, oops.Errorf("module %q has no description.<locale>.yml", slug)
		}

		lessons, err := parseLessons(repo, moduleDir, name, spec, courseSlug)
		if err != nil {
			return nil, err
		}

		modules = append(modules, Module{Order: order, Slug: slug, Infos: infos, Lessons: lessons})
	}

	// Sort by the numeric prefix so build/display order follows intent, not the
	// filesystem's lexical directory order.
	slices.SortStableFunc(modules, func(a, b Module) int { return cmp.Compare(a.Order, b.Order) })

	return modules, nil
}

func parseModuleInfos(repo repository, moduleDir string) ([]ModuleInfo, error) {
	// Glob returns its matches sorted, which keeps the locale order deterministic.
	files, err := fs.Glob(repo.fsys, path.Join(moduleDir, "description.*.yml"))
	if err != nil {
		return nil, oops.Wrapf(err, "glob module infos in %s", repo.abs(moduleDir))
	}

	infos := make([]ModuleInfo, 0, len(files))
	for _, file := range files {
		locale := localeFromDescriptionFile(path.Base(file))
		if locale == "" {
			continue
		}

		raw, err := fs.ReadFile(repo.fsys, file)
		if err != nil {
			return nil, oops.Wrapf(err, "read %s", repo.abs(file))
		}

		var data struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
		}
		if err := decodeYAML(raw, &data); err != nil {
			return nil, oops.Wrapf(err, "parse %s", repo.abs(file))
		}

		infos = append(infos, ModuleInfo{Locale: locale, Name: data.Name, Description: data.Description})
	}

	return infos, nil
}

func parseLessons(repo repository, moduleDir, moduleDirName string, spec *Spec, courseSlug string) ([]Lesson, error) {
	dirs, err := subdirs(repo, moduleDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read module dir %s", repo.abs(moduleDir))
	}

	lessons := make([]Lesson, 0, len(dirs))
	for _, name := range dirs {
		order, slug, err := splitOrderSlug(name)
		if err != nil {
			return nil, oops.Wrapf(err, "lesson dir %q", name)
		}

		lessonDir := path.Join(moduleDir, name)

		testCode, err := readFile(repo, path.Join(lessonDir, spec.ExerciseTestFilename))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson %q test file", slug)
		}

		originalCode, err := readFile(repo, path.Join(lessonDir, spec.ExerciseFilename))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson %q exercise file", slug)
		}

		infos, err := parseLessonInfos(repo, lessonDir)
		if err != nil {
			return nil, err
		}

		if len(infos) == 0 {
			return nil, oops.Errorf("lesson %q has no <locale>/data.yml", slug)
		}

		lessons = append(lessons, Lesson{
			Order:        order,
			Slug:         slug,
			TestCode:     testCode,
			OriginalCode: originalCode,
			PreparedCode: prepareCode(originalCode),
			// Matches legacy path_to_code: an absolute path INSIDE the runtime
			// container's extracted /exercises-<slug> tree, used by `make -C <path>`.
			PathToCode:   "/exercises-" + courseSlug + "/modules/" + moduleDirName + "/" + name,
			Infos:        infos,
			NaturalOrder: 0,
		})
	}

	slices.SortStableFunc(lessons, func(a, b Lesson) int { return cmp.Compare(a.Order, b.Order) })

	return lessons, nil
}

func parseLessonInfos(repo repository, lessonDir string) ([]LessonInfo, error) {
	dirs, err := subdirs(repo, lessonDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read lesson dir %s", repo.abs(lessonDir))
	}

	infos := make([]LessonInfo, 0)

	for _, name := range dirs {
		localeDir := path.Join(lessonDir, name)
		// A locale dir is exactly one that carries a data.yml — this is how the
		// legacy loader tells `en/`/`ru/` apart from `assets/`.
		dataPath := path.Join(localeDir, "data.yml")
		if _, err := fs.Stat(repo.fsys, dataPath); err != nil {
			continue
		}

		raw, err := fs.ReadFile(repo.fsys, dataPath)
		if err != nil {
			return nil, oops.Wrapf(err, "read %s", repo.abs(dataPath))
		}

		var data struct {
			Name        string       `yaml:"name"`
			Description string       `yaml:"description"`
			Tips        []string     `yaml:"tips"`
			Definitions []Definition `yaml:"definitions"`
		}
		if err := decodeYAML(raw, &data); err != nil {
			return nil, oops.Wrapf(err, "parse %s", repo.abs(dataPath))
		}

		theory, err := readFile(repo, path.Join(localeDir, "README.md"))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson locale %q theory", name)
		}

		instructions, err := readFile(repo, path.Join(localeDir, "EXERCISE.md"))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson locale %q instructions", name)
		}

		infos = append(infos, LessonInfo{
			Locale:       name,
			Dir:          repo.abs(localeDir),
			Name:         data.Name,
			Description:  data.Description,
			Theory:       theory,
			Instructions: instructions,
			Tips:         data.Tips,
			Definitions:  data.Definitions,
		})
	}

	slices.SortFunc(infos, func(a, b LessonInfo) int { return strings.Compare(a.Locale, b.Locale) })

	return infos, nil
}

// assignNaturalOrder numbers every lesson 1..N in build order across ALL modules
// (module order, then lesson order within each), matching legacy's flattened
// each_with_index. natural_order is the global lesson sequence the player uses
// for prev/next navigation, so it must span modules, not reset per module.
func assignNaturalOrder(modules []Module) {
	n := 0

	for mi := range modules {
		for li := range modules[mi].Lessons {
			n++
			modules[mi].Lessons[li].NaturalOrder = n
		}
	}
}

// prepareCodeRe finds a solution region delimited by a line containing BEGIN and
// a line containing END. `(?s)` lets `.` span newlines (Ruby's /m); `(?m)` makes
// ^/$ match line boundaries. Lazy quantifiers keep the region minimal, so nested
// BEGIN/END blocks each collapse independently.
var prepareCodeRe = regexp.MustCompile(`(?sm)(?P<begin>^[^\n]*?BEGIN.*?$\s*)(?P<content>.+?)(?P<end>^[^\n]*?END.*?$)`)

// prepareCode derives the student template from the reference solution: it blanks
// out everything between BEGIN and END and relabels the BEGIN marker. If the code
// has no BEGIN/END markers the template is empty (legacy returns ""), signalling
// "no editable region" to the player. Straight port of legacy ExerciseLoader#prepare_code.
func prepareCode(code string) string {
	result := prepareCodeRe.ReplaceAllString(code, "${begin}\n${end}")
	if result == code {
		return ""
	}

	return strings.ReplaceAll(result, "BEGIN", "BEGIN (write your solution here)")
}

// splitOrderSlug parses a "<order>-<slug>" directory name. Split on the FIRST
// hyphen only, so slugs may contain hyphens (legacy split("-", 2)). The order
// prefix must be numeric.
func splitOrderSlug(name string) (int, string, error) {
	prefix, slug, found := strings.Cut(name, "-")
	if !found {
		return 0, "", oops.Errorf("%q is not <order>-<slug>", name)
	}

	order, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, "", oops.Errorf("%q has non-numeric order prefix %q", name, prefix)
	}

	return order, slug, nil
}

// localeFromDescriptionFile extracts `en` from `description.en.yml`. Returns ""
// for names that don't fit the pattern (skipped by the caller).
func localeFromDescriptionFile(base string) string {
	parts := strings.Split(base, ".")
	if len(parts) != 3 || parts[0] != "description" || parts[2] != "yml" {
		return ""
	}

	return parts[1]
}

// subdirs returns the names of immediate subdirectories of dir, sorted (fs.ReadDir
// already sorts by name). A missing directory is an error the caller wraps.
func subdirs(repo repository, dir string) ([]string, error) {
	entries, err := fs.ReadDir(repo.fsys, dir)
	if err != nil {
		return nil, errs.Wrapf(err, "read directory")
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}

	return names, nil
}

func readFile(repo repository, name string) (string, error) {
	raw, err := fs.ReadFile(repo.fsys, name)
	if err != nil {
		return "", errs.Wrapf(err, "read file")
	}

	return string(raw), nil
}

// decodeYAML rejects misspelled and unsupported keys in every course metadata
// file instead of silently dropping them and building incomplete course data.
func decodeYAML(raw []byte, target any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	decoder.KnownFields(true)

	err := decoder.Decode(target)
	if err != nil {
		return errs.Wrapf(err, "decode yaml")
	}

	return nil
}
