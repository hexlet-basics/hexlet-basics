package courseloader

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/oops"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"

	"hexletbasics/ent"
	"hexletbasics/ent/courselesson"
	"hexletbasics/ent/coursemodule"
	"hexletbasics/ent/courseversion"
	"hexletbasics/internal/assetstore"
	"hexletbasics/internal/store"
)

// Course version states, mirroring the legacy AASM machine on
// `Language::Version`. A build is only attempted from `created`; it moves to
// `building` for the duration and finishes at `built` (promoted live) or `failed`
// (previous live version untouched).
const (
	stateCreated  = "created"
	stateBuilding = "building"
	stateBuilt    = "built"
	stateFailed   = "failed"
)

// lessonStateCreated is the initial state a freshly-upserted lesson gets, mirroring
// the legacy Language::Lesson AASM initial state.
const lessonStateCreated = "created"

// Loader builds a new course version from its exercises repository: fetch →
// parse → write. It is the Go port of the legacy ExerciseLoader#run, including
// the load-bearing invariant that a version becomes live (Course.currentVersionId
// flips) ONLY on success, in the same transaction that marks it built — so a
// broken rebuild never takes a live course offline.
type Loader struct {
	db      *ent.Client
	txStore store.Transactor
	assets  *assetstore.Store
	fetcher Fetcher
}

// NewLoader wires the loader to its dependencies. The fetcher is an interface so
// tests can supply a fixture directory instead of cloning.
func NewLoader(db *ent.Client, txStore store.Transactor, assets *assetstore.Store, fetcher Fetcher) *Loader {
	return &Loader{db: db, txStore: txStore, assets: assets, fetcher: fetcher}
}

// Run builds the given course version. It is idempotent against dead runs: a
// version not in `created` is skipped (not re-run), matching legacy may_build?.
// Any failure marks the version `failed` with the error and returns it; the
// previously-live version is never disturbed.
func (l *Loader) Run(ctx context.Context, versionID int) error {
	version, err := l.db.CourseVersion.Get(ctx, versionID)
	if err != nil {
		return oops.Wrapf(err, "load course version %d", versionID)
	}

	claimed, err := l.claim(ctx, versionID)
	if err != nil {
		return err
	}
	if !claimed {
		return nil
	}

	course, err := l.db.Course.Get(ctx, version.LanguageID)
	if err != nil {
		return oops.Wrapf(err, "load course %d for version %d", version.LanguageID, versionID)
	}
	if course.Slug == nil || *course.Slug == "" {
		return l.fail(ctx, version, oops.Errorf("course %d has no slug", course.ID))
	}
	slug := *course.Slug

	dir, cleanup, err := l.fetcher.Fetch(ctx, slug)
	if err != nil {
		return l.fail(ctx, version, oops.Wrapf(err, "fetch course %q", slug))
	}
	defer cleanup()

	parsed, err := Parse(dir, slug)
	if err != nil {
		return l.fail(ctx, version, oops.Wrapf(err, "parse course %q", slug))
	}

	// Upload theory images BEFORE the write transaction, rewriting each info's
	// theory in place. Blob writes are network I/O and are NOT transactional, so
	// keeping them inside the DB transaction would neither buy atomicity (bytes
	// already landed on a rollback) nor be safe at scale (a course has hundreds of
	// lessons — an S3-writing transaction held open that long risks statement
	// timeouts and pins a connection). The transaction that follows only writes
	// rows, and still flips the course live only on success.
	if err := l.uploadImages(ctx, parsed); err != nil {
		return l.fail(ctx, version, oops.Wrapf(err, "upload images for %q", slug))
	}

	if err := l.build(ctx, version, course, parsed); err != nil {
		return l.fail(ctx, version, err)
	}
	return nil
}

// claim atomically transitions a freshly-created version to building. Checking
// the affected-row count is the concurrency guard: only one worker can satisfy
// the state predicate, regardless of how many duplicate jobs reach Run.
func (l *Loader) claim(ctx context.Context, versionID int) (bool, error) {
	updated, err := l.db.CourseVersion.Update().
		Where(
			courseversion.ID(versionID),
			courseversion.State(stateCreated),
		).
		SetState(stateBuilding).
		Save(ctx)
	if err != nil {
		return false, oops.Wrapf(err, "claim course version %d", versionID)
	}
	return updated == 1, nil
}

// fail records the build error on the version and transitions it to `failed`,
// then returns the original cause (so the worker surfaces it). The state write
// uses the passed context; a cancelled job leaves the version `building` for the
// reaper, which is the correct fallback.
func (l *Loader) fail(ctx context.Context, version *ent.CourseVersion, cause error) error {
	_, _ = l.db.CourseVersion.UpdateOne(version).
		SetState(stateFailed).
		SetResult(fmt.Sprintf("Error: %s", cause.Error())).
		Save(ctx)
	return cause
}

// build writes the whole parsed course in ONE transaction: the version's spec
// metadata, every module/lesson and their per-locale infos, and finally the
// built-state + live-promotion. Atomicity is the point — either the new version
// is fully present and live, or nothing changed and the old version stays live.
func (l *Loader) build(ctx context.Context, version *ent.CourseVersion, course *ent.Course, parsed *Course) error {
	return l.txStore.WithinTx(ctx, func(_ *sql.Tx, txClient *ent.Client) error {
		return l.buildTx(ctx, txClient, version.ID, course.ID, parsed)
	})
}

func (l *Loader) buildTx(ctx context.Context, tx *ent.Client, versionID, languageID int, parsed *Course) error {
	spec := parsed.Spec
	if _, err := tx.CourseVersion.UpdateOneID(versionID).
		SetName(spec.Name).
		SetProgress(spec.Progress).
		SetLearnAs(spec.LearnAs).
		SetExtension(spec.Extension).
		SetDockerImage(spec.DockerImage).
		SetExerciseFilename(spec.ExerciseFilename).
		SetExerciseTestFilename(spec.ExerciseTestFilename).
		Save(ctx); err != nil {
		return oops.Wrapf(err, "update version spec")
	}

	totalLessons := 0
	for _, m := range parsed.Modules {
		moduleID, err := upsertModule(ctx, tx, languageID, m.Slug)
		if err != nil {
			return err
		}

		mv, err := tx.CourseModuleVersion.Create().
			SetLanguageID(languageID).
			SetLanguageVersionID(versionID).
			SetModuleID(moduleID).
			SetOrder(m.Order).
			Save(ctx)
		if err != nil {
			return oops.Wrapf(err, "create module version %q", m.Slug)
		}

		for _, info := range m.Infos {
			if _, err := tx.CourseModuleTranslation.Create().
				SetLanguageID(languageID).
				SetLanguageVersionID(versionID).
				SetVersionID(mv.ID).
				SetLocale(info.Locale).
				SetName(info.Name).
				SetDescription(info.Description).
				Save(ctx); err != nil {
				return oops.Wrapf(err, "create module info %q/%s", m.Slug, info.Locale)
			}
		}

		for _, lesson := range m.Lessons {
			lessonID, err := upsertLesson(ctx, tx, languageID, moduleID, lesson.Slug)
			if err != nil {
				return err
			}

			lv, err := tx.CourseLessonVersion.Create().
				SetLanguageID(languageID).
				SetLanguageVersionID(versionID).
				SetLessonID(lessonID).
				SetModuleVersionID(mv.ID).
				SetOrder(lesson.Order).
				SetNaturalOrder(lesson.NaturalOrder).
				SetTestCode(lesson.TestCode).
				SetOriginalCode(lesson.OriginalCode).
				SetPreparedCode(lesson.PreparedCode).
				SetPathToCode(lesson.PathToCode).
				Save(ctx)
			if err != nil {
				return oops.Wrapf(err, "create lesson version %q", lesson.Slug)
			}

			for _, info := range lesson.Infos {
				if err := l.createLessonInfo(ctx, tx, languageID, versionID, lessonID, lv.ID, info); err != nil {
					return err
				}
			}
			totalLessons++
		}
	}

	// Mark built + promote to live, atomically with all the content above. This is
	// the moment the new version becomes the one learners see.
	if _, err := tx.CourseVersion.UpdateOneID(versionID).
		SetState(stateBuilt).
		SetResult("Success").
		SetLessonsCount(totalLessons).
		Save(ctx); err != nil {
		return oops.Wrapf(err, "mark version built")
	}
	if _, err := tx.Course.UpdateOneID(languageID).
		SetCurrentVersionID(versionID).
		SetLessonsCount(totalLessons).
		Save(ctx); err != nil {
		return oops.Wrapf(err, "promote version to live")
	}
	return nil
}

// createLessonInfo writes one locale's lesson info. Theory images were already
// uploaded and the markdown rewritten by uploadImages before this transaction, so
// info.Theory is used as-is here. tips/definitions are serialized as YAML arrays
// for Rails `serialize type: Array` compatibility.
func (l *Loader) createLessonInfo(ctx context.Context, tx *ent.Client, languageID, versionID, lessonID, lessonVersionID int, info LessonInfo) error {
	create := tx.CourseLessonTranslation.Create().
		SetLanguageID(languageID).
		SetLanguageVersionID(versionID).
		SetLanguageLessonID(lessonID).
		SetVersionID(lessonVersionID).
		SetLocale(info.Locale).
		SetName(info.Name).
		SetTheory(info.Theory).
		SetInstructions(info.Instructions)

	if info.Description != "" {
		create.SetDescription(info.Description)
	}
	if len(info.Tips) > 0 {
		s, err := marshalYAML(info.Tips)
		if err != nil {
			return oops.Wrapf(err, "serialize tips")
		}
		create.SetTips(s)
	}
	if len(info.Definitions) > 0 {
		s, err := marshalYAML(info.Definitions)
		if err != nil {
			return oops.Wrapf(err, "serialize definitions")
		}
		create.SetDefinitions(s)
	}

	if _, err := create.Save(ctx); err != nil {
		return oops.Wrapf(err, "create lesson info %s", info.Locale)
	}
	return nil
}

// upsertModule atomically creates or finds the stable module row for (course,
// slug). Identity is kept across rebuilds so downstream references survive; the
// per-build ordering lives on the module version, not here.
func upsertModule(ctx context.Context, tx *ent.Client, languageID int, slug string) (int, error) {
	id, err := tx.CourseModule.Create().
		SetLanguageID(languageID).
		SetSlug(slug).
		OnConflictColumns(coursemodule.FieldLanguageID, coursemodule.FieldSlug).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return 0, oops.Wrapf(err, "upsert module %q", slug)
	}
	return id, nil
}

// upsertLesson atomically creates or finds the stable lesson row for (course,
// slug), and (re)points it at its current module — a lesson can move between
// modules across rebuilds. Learner progress FKs this stable id, so it must NOT be
// recreated.
func upsertLesson(ctx context.Context, tx *ent.Client, languageID, moduleID int, slug string) (int, error) {
	id, err := tx.CourseLesson.Create().
		SetLanguageID(languageID).
		SetModuleID(moduleID).
		SetSlug(slug).
		SetState(lessonStateCreated).
		OnConflictColumns(courselesson.FieldLanguageID, courselesson.FieldSlug).
		Update(func(update *ent.CourseLessonUpsert) {
			update.SetModuleID(moduleID)
		}).
		ID(ctx)
	if err != nil {
		return 0, oops.Wrapf(err, "upsert lesson %q", slug)
	}
	return id, nil
}

// uploadImages walks every lesson info and rewrites its theory in place: each
// locally-referenced Markdown image is uploaded to the blob bucket and its URL
// swapped for the served /storage path. This runs BEFORE the write transaction so
// no network I/O is held inside it. Attachment rows/blobs it creates for a build
// that later fails are harmless orphans (the /storage read path keys off the blob
// directly, not the Attachment row).
func (l *Loader) uploadImages(ctx context.Context, parsed *Course) error {
	for mi := range parsed.Modules {
		for li := range parsed.Modules[mi].Lessons {
			infos := parsed.Modules[mi].Lessons[li].Infos
			for ii := range infos {
				theory, err := l.processImages(ctx, infos[ii].Theory, infos[ii].Dir)
				if err != nil {
					return oops.Wrapf(err, "lesson info %s", infos[ii].Locale)
				}
				infos[ii].Theory = theory
			}
		}
	}
	return nil
}

// processImages uploads each locally-referenced theory image to the blob bucket
// and rewrites its Markdown URL to the served /storage path. Remote (http) images
// are left untouched. A missing referenced file fails the build (the repo is
// broken), matching legacy strictness.
func (l *Loader) processImages(ctx context.Context, theory, localeDir string) (string, error) {
	source := []byte(theory)
	document := goldmark.DefaultParser().Parse(text.NewReader(source))
	var replacements []markdownImageReplacement

	err := ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		image, ok := node.(*ast.Image)
		if !ok || image.Reference != nil || len(image.Destination) == 0 {
			return ast.WalkContinue, nil
		}

		ref := string(image.Destination)
		if strings.HasPrefix(ref, "http") {
			return ast.WalkSkipChildren, nil
		}

		url, err := l.uploadImage(ctx, localeDir, ref)
		if err != nil {
			return ast.WalkStop, err
		}

		start, ok := sourceSliceOffset(source, image.Destination)
		if !ok {
			return ast.WalkStop, oops.Errorf("locate Markdown image destination %q", ref)
		}
		replacements = append(replacements, markdownImageReplacement{
			start: start,
			end:   start + len(image.Destination),
			url:   url,
		})
		return ast.WalkSkipChildren, nil
	})
	if err != nil {
		return "", err
	}

	sort.Slice(replacements, func(i, j int) bool {
		return replacements[i].start < replacements[j].start
	})

	var result bytes.Buffer
	cursor := 0
	for _, replacement := range replacements {
		result.Write(source[cursor:replacement.start])
		result.WriteString(replacement.url)
		cursor = replacement.end
	}
	result.Write(source[cursor:])
	return result.String(), nil
}

type markdownImageReplacement struct {
	start int
	end   int
	url   string
}

// sourceSliceOffset finds a parser-owned subslice in the original Markdown.
// Goldmark keeps inline destinations backed by the source buffer, which lets us
// replace only the URL while preserving every other byte of the author's text.
func sourceSliceOffset(source, part []byte) (int, bool) {
	if len(part) == 0 || len(part) > len(source) {
		return 0, false
	}
	for i := 0; i+len(part) <= len(source); i++ {
		if &source[i] == &part[0] {
			return i, true
		}
	}
	return 0, false
}

// uploadImage opens a theory image referenced relative to the locale dir and
// delegates its complete storage lifecycle to the shared asset store. OpenInRoot
// confines symlink resolution as well as lexical traversal to localeDir.
func (l *Loader) uploadImage(ctx context.Context, localeDir, ref string) (string, error) {
	file, err := os.OpenInRoot(localeDir, ref)
	if err != nil {
		return "", oops.Wrapf(err, "open image %q", ref)
	}
	defer func() { _ = file.Close() }()

	attachment, err := l.assets.Put(ctx, assetstore.Upload{
		Filename: filepath.Base(ref),
		Body:     file,
	})
	if err != nil {
		return "", oops.Wrapf(err, "store image %q", ref)
	}

	return attachment.URL, nil
}

// marshalYAML serializes tips/definitions to the YAML-array form Rails'
// `serialize type: Array` produced, keeping loaded rows readable the same way
// across the cutover.
func marshalYAML(v any) (string, error) {
	out, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
