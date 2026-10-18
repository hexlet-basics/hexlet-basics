package transition

// PublicationState is where a piece of content stands in its life (GLOSSARY,
// Publication State): a draft nobody but staff sees, published for everyone, or
// archived and withdrawn. Blog posts, banners, reviews and landing pages all
// carry it as a plain string column, so the legal moves are written once, here,
// over that vocabulary instead of once per resource over each generated enum.
type PublicationState string

// Vocabulary names this set of states for an adapter that words a refused move.
func (PublicationState) Vocabulary() string { return "publication_state" }

const (
	Draft     PublicationState = "draft"
	Published PublicationState = "published"
	Archived  PublicationState = "archived"
)

// The legal moves of a Publication State, one per operation (publish,
// unpublish, archive). Each names the states it may start from; repeating the
// current state is not among them, so it is refused like any illegal move.
var (
	Publish = Transition[PublicationState]{
		To:   Published,
		From: []PublicationState{Draft, Archived},
	}
	Unpublish = Transition[PublicationState]{
		To:   Draft,
		From: []PublicationState{Published},
	}
	Archive = Transition[PublicationState]{
		To:   Archived,
		From: []PublicationState{Draft, Published},
	}
)

// CurrentOf reads the state column of a row. A row with no state (the column is
// nullable in the legacy schema) counts as a draft, which is what the public
// queries treat it as: not published.
func CurrentOf(column *string) PublicationState {
	if column == nil {
		return Draft
	}
	return PublicationState(*column)
}
