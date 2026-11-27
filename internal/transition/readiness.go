package transition

// Readiness is how finished a Course is (GLOSSARY, Course Readiness): a draft, in
// development, or completed. It is a property of the product, not a publication
// workflow, so it has a table of its own next to the Publication State's; it is a
// state all the same, and changes only through its own operations.
type Readiness string

// Vocabulary names this set of states for an adapter that words a refused move.
func (Readiness) Vocabulary() Vocabulary { return CourseReadiness }

// Values lists the readiness values, for ent's enum validation of the column.
func (Readiness) Values() []string { return values(CourseReadiness) }

const (
	// ReadinessDraft is a Course whose content is not started or not usable.
	ReadinessDraft Readiness = "draft"
	// ReadinessInDevelopment is a Course being written: usable, not finished.
	ReadinessInDevelopment Readiness = "in_development"
	// ReadinessCompleted is a finished Course.
	ReadinessCompleted Readiness = "completed"
)

// The legal moves of a Course's Readiness, one per operation. Readiness says how
// finished the Course is, not where it stands in a workflow, so any change
// between two distinct values is legal (a completed Course may go back to
// development, a draft may be declared completed). Each names the two values it
// may start from; repeating the current one is not among them, so it is refused
// like any illegal move.

// MarkDraft sends a Course back to draft from any other readiness.
func MarkDraft() Transition[Readiness] {
	return Transition[Readiness]{
		To:   ReadinessDraft,
		From: []Readiness{ReadinessInDevelopment, ReadinessCompleted},
	}
}

// MarkInDevelopment marks a Course as being written.
func MarkInDevelopment() Transition[Readiness] {
	return Transition[Readiness]{
		To:   ReadinessInDevelopment,
		From: []Readiness{ReadinessDraft, ReadinessCompleted},
	}
}

// MarkCompleted marks a Course as finished.
func MarkCompleted() Transition[Readiness] {
	return Transition[Readiness]{
		To:   ReadinessCompleted,
		From: []Readiness{ReadinessDraft, ReadinessInDevelopment},
	}
}
