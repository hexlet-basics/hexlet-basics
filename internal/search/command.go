package search

import "hexletbasics/internal/domain"

// SearchCatalogue is a search of the catalogue, the public one by default: what the words are, in
// which language, for which kinds of record, and which window of the answer.
type SearchCatalogue struct {
	// Text is what the visitor typed. It is read as a web search query (quoted
	// phrases, `or`, `-word`), never as tsquery syntax, so no input is an error.
	Text string
	// Locale is the request locale: only records of it are searched, with the
	// text search configuration of its language.
	Locale string
	// Types narrows the search to some kinds of record; empty is all of them.
	Types []Type
	// Unrestricted is an editor's search: every record of the locale is found
	// whatever its state (a draft, archived or unlisted landing page, a course
	// not completed yet, a lesson of any build, an unpublished post), and each
	// hit says which state it is in. The public search never sets it; whoever
	// does must have decided that the caller may read drafts.
	Unrestricted bool
	Page         domain.Page
}

// Hit is one record that matches, in the shape of a result list.
type Hit struct {
	Type Type
	// ID is the typed id fetch-style reads take (CourseID, LessonID, PostID).
	ID    string
	Title string
	// Snippet is a few words of the record's text around the match, plain text.
	Snippet string
	// Path is the record's page on the website, with the locale prefix and
	// without an origin.
	Path string
	// State and Readiness say how published the record is: the publication state
	// of a landing page, lesson or post, and the readiness of the course. They
	// are filled by an unrestricted search; a public hit is published and its
	// course completed by construction.
	State     string
	Readiness string
}
