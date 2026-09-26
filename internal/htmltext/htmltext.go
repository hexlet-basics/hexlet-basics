// Package htmltext reduces stored editor HTML to its plain text — the Go stand-in
// for ActionText's `to_plain_text`, which the legacy blog used both for the
// reading-time estimate and for the text it sent to the LLM.
package htmltext

import (
	"strings"

	"golang.org/x/net/html"
)

// PlainText returns the text nodes of an HTML fragment joined by single
// spaces, so adjacent blocks (`<p>a</p><p>b</p>`) never glue into one word.
// Markup is tokenized, not parsed: the fragment needs no well-formedness.
func PlainText(fragment string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(fragment))
	var plain strings.Builder
	for {
		switch tokenizer.Next() {
		case html.TextToken:
			plain.Write(tokenizer.Text())
			plain.WriteByte(' ')
		case html.ErrorToken:
			return plain.String()
		}
	}
}
