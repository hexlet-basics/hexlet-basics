package relatedcourse

// OfferedCandidates exposes the candidate query to the external tests, which
// cannot be internal because testsupport imports this package.
var OfferedCandidates = offeredCandidates

// Candidate names the unexported candidate type for the external tests.
type Candidate = candidate
