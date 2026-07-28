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
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	spec, err := parseSpec(dir)
	if err != nil {
		return nil, err
	}

	modules, err := parseModules(filepath.Join(dir, "modules"), spec, courseSlug)
	if err != nil {
		return nil, err
	}

	assignNaturalOrder(modules)

	return &Course{Spec: *spec, Modules: modules}, nil
}

func parseSpec(dir string) (*Spec, error) {
	path := filepath.Join(dir, "spec.yml")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, oops.Wrapf(err, "read %s", path)
	}
	var doc struct {
		Language Spec `yaml:"language"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, oops.Wrapf(err, "parse %s", path)
	}
	s := doc.Language
	// Mirror legacy's fetch("...") — a missing key is a hard error, not a blank
	// build. Report the first offender so the fix is obvious.
	for k, v := range map[string]string{
		"name":                   s.Name,
		"progress":               s.Progress,
		"learn_as":               s.LearnAs,
		"docker_image":           s.DockerImage,
		"extension":              s.Extension,
		"exercise_filename":      s.ExerciseFilename,
		"exercise_test_filename": s.ExerciseTestFilename,
	} {
		if v == "" {
			return nil, oops.Errorf("spec.yml: language.%s is required", k)
		}
	}
	return &s, nil
}

func parseModules(modulesDir string, spec *Spec, courseSlug string) ([]Module, error) {
	dirs, err := subdirs(modulesDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read modules dir %s", modulesDir)
	}

	modules := make([]Module, 0, len(dirs))
	for _, name := range dirs {
		order, slug, err := splitOrderSlug(name)
		if err != nil {
			return nil, oops.Wrapf(err, "module dir %q", name)
		}
		moduleDir := filepath.Join(modulesDir, name)

		infos, err := parseModuleInfos(moduleDir)
		if err != nil {
			return nil, err
		}
		if len(infos) == 0 {
			return nil, oops.Errorf("module %q has no description.<locale>.yml", slug)
		}

		lessons, err := parseLessons(moduleDir, name, spec, courseSlug)
		if err != nil {
			return nil, err
		}

		modules = append(modules, Module{Order: order, Slug: slug, Infos: infos, Lessons: lessons})
	}

	// Sort by the numeric prefix so build/display order follows intent, not the
	// filesystem's lexical directory order.
	sort.Slice(modules, func(i, j int) bool { return modules[i].Order < modules[j].Order })
	return modules, nil
}

func parseModuleInfos(moduleDir string) ([]ModuleInfo, error) {
	files, err := filepath.Glob(filepath.Join(moduleDir, "description.*.yml"))
	if err != nil {
		return nil, oops.Wrapf(err, "glob module infos in %s", moduleDir)
	}
	sort.Strings(files) // deterministic locale order

	infos := make([]ModuleInfo, 0, len(files))
	for _, file := range files {
		locale := localeFromDescriptionFile(filepath.Base(file))
		if locale == "" {
			continue
		}
		raw, err := os.ReadFile(file)
		if err != nil {
			return nil, oops.Wrapf(err, "read %s", file)
		}
		var data struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
		}
		if err := yaml.Unmarshal(raw, &data); err != nil {
			return nil, oops.Wrapf(err, "parse %s", file)
		}
		infos = append(infos, ModuleInfo{Locale: locale, Name: data.Name, Description: data.Description})
	}
	return infos, nil
}

func parseLessons(moduleDir, moduleDirName string, spec *Spec, courseSlug string) ([]Lesson, error) {
	dirs, err := subdirs(moduleDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read module dir %s", moduleDir)
	}

	lessons := make([]Lesson, 0, len(dirs))
	for _, name := range dirs {
		order, slug, err := splitOrderSlug(name)
		if err != nil {
			return nil, oops.Wrapf(err, "lesson dir %q", name)
		}
		lessonDir := filepath.Join(moduleDir, name)

		testCode, err := readFile(filepath.Join(lessonDir, spec.ExerciseTestFilename))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson %q test file", slug)
		}
		originalCode, err := readFile(filepath.Join(lessonDir, spec.ExerciseFilename))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson %q exercise file", slug)
		}

		infos, err := parseLessonInfos(lessonDir)
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
			PathToCode: "/exercises-" + courseSlug + "/modules/" + moduleDirName + "/" + name,
			Infos:      infos,
		})
	}

	sort.Slice(lessons, func(i, j int) bool { return lessons[i].Order < lessons[j].Order })
	return lessons, nil
}

func parseLessonInfos(lessonDir string) ([]LessonInfo, error) {
	dirs, err := subdirs(lessonDir)
	if err != nil {
		return nil, oops.Wrapf(err, "read lesson dir %s", lessonDir)
	}

	infos := make([]LessonInfo, 0)
	for _, name := range dirs {
		localeDir := filepath.Join(lessonDir, name)
		// A locale dir is exactly one that carries a data.yml — this is how the
		// legacy loader tells `en/`/`ru/` apart from `assets/`.
		dataPath := filepath.Join(localeDir, "data.yml")
		if _, err := os.Stat(dataPath); err != nil {
			continue
		}

		raw, err := os.ReadFile(dataPath)
		if err != nil {
			return nil, oops.Wrapf(err, "read %s", dataPath)
		}
		var data struct {
			Name        string       `yaml:"name"`
			Description string       `yaml:"description"`
			Tips        []string     `yaml:"tips"`
			Definitions []Definition `yaml:"definitions"`
		}
		if err := yaml.Unmarshal(raw, &data); err != nil {
			return nil, oops.Wrapf(err, "parse %s", dataPath)
		}

		theory, err := readFile(filepath.Join(localeDir, "README.md"))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson locale %q theory", name)
		}
		instructions, err := readFile(filepath.Join(localeDir, "EXERCISE.md"))
		if err != nil {
			return nil, oops.Wrapf(err, "lesson locale %q instructions", name)
		}

		infos = append(infos, LessonInfo{
			Locale:       name,
			Dir:          localeDir,
			Name:         data.Name,
			Description:  data.Description,
			Theory:       theory,
			Instructions: instructions,
			Tips:         data.Tips,
			Definitions:  data.Definitions,
		})
	}

	sort.Slice(infos, func(i, j int) bool { return infos[i].Locale < infos[j].Locale })
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

// subdirs returns the names of immediate subdirectories of dir, sorted (os.ReadDir
// already sorts by name). A missing directory is an error the caller wraps.
func subdirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func readFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
