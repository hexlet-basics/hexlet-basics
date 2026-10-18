// Package search finds a course, a lesson or a blog post of the public catalogue
// by what it says, with PostgreSQL full-text search (#877).
//
// It searches only text the site's authors wrote (names, descriptions, theory,
// instructions, post bodies) and only records a visitor can read on the website
// today: the same scopes the public reads use (internal/landingpage,
// internal/blog), so a draft, an unpublished post, a lesson of a retired
// build or a record of another locale is not found by search any more than it is
// by URL. Nothing a learner wrote is searched (no review, no message to the
// assistant, no lead): the tables of learner content are not touched at all, and
// keeping it that way is what lets an assistant treat a result as the site's own
// text.
package search

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	entgo "entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"hexletbasics/ent"
	"hexletbasics/internal/domain"
	"hexletbasics/internal/localization"
)

// maxWindow is how deep into the answer a page may reach: the first thousand
// hits. A search is for finding a page, not for copying the catalogue (the lists
// do that), and each type's query must produce page*perPage rows to merge, so an
// unbounded page number would be an unbounded query.
const maxWindow = 1000

// maxTextBytes bounds the query text the database is asked to parse.
const maxTextBytes = 512

// snippetRunes is the longest snippet, in characters.
const snippetRunes = 300

// Catalogue is the full-text search of the public catalogue.
type Catalogue struct {
	db *ent.Client
}

// NewCatalogue builds the search over an ent client.
func NewCatalogue(db *ent.Client) *Catalogue {
	return &Catalogue{db: db}
}

// match is one search, as the three sources need it: the text search
// configuration of the locale and the words.
type match struct {
	config string
	locale string
	text   string
	// unrestricted drops the public storefront's filters (see Query).
	unrestricted bool
}

// configs are the text search configurations of the languages the site is in,
// which the migration indexes (catalogue_search_vector, one GIN index each). The
// name is written into the SQL as a literal, never taken from a request: only
// these values can be there, and a literal is what lets the planner match the
// index expression. A locale without one is searched with the language-neutral
// `simple`, with no index, which is right while it has no content.
func configs() map[string]string {
	return map[string]string{"ru": "russian", "en": "english"}
}

// candidate is a hit before its page is known: how well it matches, and
// what is needed to tell it from another.
type candidate struct {
	rank float64
	typ  Type
	id   int
	hit  Hit
}

// source is one kind of record the search covers.
type source interface {
	// count is how many records of the kind match.
	count(ctx context.Context, m match) (int, error)
	// top is the n best matches, best first, without snippets.
	top(ctx context.Context, m match, n int) ([]candidate, error)
	// snippets are the snippets of the records with ids, by id.
	snippets(ctx context.Context, m match, ids []int) (map[int]string, error)
}

// Search answers the command: the hits of the window, and how many records match
// in all, best match first. A tie is broken by the kind and then by the record's
// age, so a page never changes under the reader between two calls.
//
// Each kind is queried for its best page*perPage matches, which is enough to
// merge: the best N of the union are among the best N of each part. Only the hits
// of the window get a snippet, in a second query per kind, because producing
// one reads the whole text of the record.
func (c *Catalogue) Search(ctx context.Context, cmd Query) ([]Hit, int, error) {
	// PostgreSQL text holds neither a NUL nor an invalid sequence, and refuses
	// the statement that carries one: a stray byte in what was typed is dropped,
	// not an error.
	text := strings.TrimSpace(strings.ToValidUTF8(strings.ReplaceAll(cmd.Text, "\x00", ""), ""))
	if len(text) > maxTextBytes {
		text = strings.ToValidUTF8(text[:maxTextBytes], "")
	}

	if text == "" || cmd.Page.Limit <= 0 {
		return []Hit{}, 0, nil
	}

	config, ok := configs()[cmd.Locale]
	if !ok {
		config = "simple"
	}

	matcher := match{config: config, locale: cmd.Locale, text: text, unrestricted: cmd.Unrestricted}

	sources := c.sources(cmd.Types)
	total := 0

	for _, src := range sources {
		n, err := src.count(ctx, matcher)
		if err != nil {
			return nil, 0, domain.FromStore(err)
		}

		total += n
	}

	window := min(cmd.Page.Offset+cmd.Page.Limit, maxWindow)
	if cmd.Page.Offset >= window {
		return []Hit{}, total, nil
	}

	var candidates []candidate

	for _, src := range sources {
		found, err := src.top(ctx, matcher, window)
		if err != nil {
			return nil, 0, domain.FromStore(err)
		}

		candidates = append(candidates, found...)
	}

	slices.SortStableFunc(candidates, func(a, b candidate) int {
		return cmp.Or(
			-cmp.Compare(a.rank, b.rank),
			cmp.Compare(typeOrder(a.typ), typeOrder(b.typ)),
			cmp.Compare(a.id, b.id),
		)
	})
	candidates = candidates[min(cmd.Page.Offset, len(candidates)):min(window, len(candidates))]

	hits, err := withSnippets(ctx, matcher, sources, candidates)
	if err != nil {
		return nil, 0, err
	}

	return hits, total, nil
}

// kind is a source with the type of record it holds.
type kind struct {
	source

	typ Type
}

// sources are the kinds the command asks for, all of them when it names none.
func (c *Catalogue) sources(types []Type) []kind {
	all := map[Type]source{
		Course: landingPages{c.db},
		Lesson: lessons{c.db},
		Post:   posts{c.db},
	}

	var out []kind

	for _, t := range Types() {
		if len(types) == 0 || slices.Contains(types, t) {
			out = append(out, kind{source: all[t], typ: t})
		}
	}

	return out
}

func typeOrder(t Type) int { return slices.Index(Types(), t) }

// withSnippets fills the snippet of each candidate's hit, asking each kind once
// for the ids it contributed.
func withSnippets(ctx context.Context, matcher match, kinds []kind, candidates []candidate) ([]Hit, error) {
	byType := map[Type][]int{}
	for _, cand := range candidates {
		byType[cand.typ] = append(byType[cand.typ], cand.id)
	}

	snippets := map[Type]map[int]string{}

	for _, entry := range kinds {
		if len(byType[entry.typ]) == 0 {
			continue
		}

		found, err := entry.snippets(ctx, matcher, byType[entry.typ])
		if err != nil {
			return nil, domain.FromStore(err)
		}

		snippets[entry.typ] = found
	}

	hits := make([]Hit, len(candidates))
	for i, cand := range candidates {
		hit := cand.hit
		hit.Snippet = tidy(snippets[cand.typ][cand.id])
		hits[i] = hit
	}

	return hits, nil
}

// pagePath is the page of a record on the website: the locale prefix and the path.
func pagePath(locale, path string) string {
	return localization.PathPrefix(locale) + path
}

// titleOf is the name of a record, or its slug when it has none, so a hit is
// never untitled.
func titleOf(slug string, names ...*string) string {
	for _, name := range names {
		if name != nil && strings.TrimSpace(*name) != "" {
			return strings.TrimSpace(*name)
		}
	}

	return slug
}

// tidy makes a snippet plain text on one line of a bounded length: the text of a
// post is stored as HTML, and a headline can cut through a tag.
func tidy(snippet string) string {
	plain := strings.Join(strings.Fields(htmlText(snippet)), " ")

	runes := []rune(plain)
	if len(runes) > snippetRunes {
		return string(runes[:snippetRunes]) + "…"
	}

	return plain
}

// ranked is the part of every source query that scores a match: the document's
// rank selected as `search_rank`, best first, then oldest first.
func ranked(rank func(*sql.Selector) sql.Querier) func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.AppendSelectExprAs(rank(s), "search_rank")
		s.OrderBy(sql.Desc("search_rank"), sql.Asc("id"))
	}
}

// rankOf reads the score the ranked option selected.
func rankOf(value func(string) (entgo.Value, error)) (float64, error) {
	raw, err := value("search_rank")
	if err != nil {
		return 0, errs.Wrapf(err, "read search rank")
	}

	switch rank := raw.(type) {
	case float64:
		return rank, nil
	case float32:
		return float64(rank), nil
	default:
		return 0, errs.With("type", fmt.Sprintf("%T", raw)).Errorf("search rank has an unexpected type")
	}
}
