package search

import "hexletbasics/internal/domain"

// SearchCatalogue is a search of the public catalogue: what the words are, in
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
	Page  domain.Page
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
}
