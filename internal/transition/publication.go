package transition

// PublicationState is where a piece of content stands in its life (GLOSSARY,
// Publication State): a draft nobody but staff sees, published for everyone, or
// archived and withdrawn. Blog posts, banners, reviews and landing pages carry
// it as the Go type of their state column (an ent enum with GoType), so the
// legal moves are written once, here, and a query or a write that names a state
// can only name one of these.
type PublicationState string

// Vocabulary names this set of states for an adapter that words a refused move.
func (PublicationState) Vocabulary() Vocabulary { return PublicationStates }

// Values lists the states, for ent's enum validation of the column.
func (PublicationState) Values() []string { return values(PublicationStates) }

const (
	// Draft is a record not shown to the public yet.
	Draft PublicationState = "draft"
	// Published is a record shown to the public.
	Published PublicationState = "published"
	// Archived is a record withdrawn from the public but kept.
	Archived PublicationState = "archived"
)

// The legal moves of a Publication State, one per operation (publish,
// unpublish, archive). Each names the states it may start from; repeating the
// current state is not among them, so it is refused like any illegal move.

// Publish moves a draft or an archived record to published.
func Publish() Transition[PublicationState] {
	return Transition[PublicationState]{
		To:   Published,
		From: []PublicationState{Draft, Archived},
	}
}

// Unpublish takes a published record back to draft.
func Unpublish() Transition[PublicationState] {
	return Transition[PublicationState]{
		To:   Draft,
		From: []PublicationState{Published},
	}
}

// Archive withdraws a draft or published record.
func Archive() Transition[PublicationState] {
	return Transition[PublicationState]{
		To:   Archived,
		From: []PublicationState{Draft, Published},
	}
}
