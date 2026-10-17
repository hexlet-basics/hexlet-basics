package transitions

// Readiness is how finished a Course is (GLOSSARY, Course Readiness): a draft, in
// development, or completed. It is a property of the product, not a publication
// workflow, so it has a table of its own next to the Publication State's; it is a
// state all the same, and changes only through its own operations.
type Readiness string

const (
	ReadinessDraft         Readiness = "draft"
	ReadinessInDevelopment Readiness = "in_development"
	ReadinessCompleted     Readiness = "completed"
)

// The legal moves of a Course's Readiness, one per operation. Readiness says how
// finished the Course is, not where it stands in a workflow, so any change
// between two distinct values is legal (a completed Course may go back to
// development, a draft may be declared completed). Each names the two values it
// may start from; repeating the current one is not among them, so it is refused
// like any illegal move.
var (
	MarkDraft = Transition[Readiness]{
		To:   ReadinessDraft,
		From: []Readiness{ReadinessInDevelopment, ReadinessCompleted},
	}
	MarkInDevelopment = Transition[Readiness]{
		To:   ReadinessInDevelopment,
		From: []Readiness{ReadinessDraft, ReadinessCompleted},
	}
	MarkCompleted = Transition[Readiness]{
		To:   ReadinessCompleted,
		From: []Readiness{ReadinessDraft, ReadinessInDevelopment},
	}
)

// ReadinessOf reads the readiness column of a Course. A Course with none (the
// column is nullable in the legacy schema) counts as a draft, which is how the
// public queries treat it: not completed, so off the public lists.
func ReadinessOf(column *string) Readiness {
	if column == nil {
		return ReadinessDraft
	}
	return Readiness(*column)
}
