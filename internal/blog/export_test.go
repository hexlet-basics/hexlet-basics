package blog

// ReadingTime exposes the estimate to the boundary tests, which need bodies of
// an exact word count that no fixture post has.
func ReadingTime(richBodyHTML string) int { return readingTime(richBodyHTML) }
