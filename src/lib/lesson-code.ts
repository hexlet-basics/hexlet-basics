// Where a learner's buffer for a lesson is kept. The editor writes it and the
// Assistant reads it, both through Mantine's useLocalStorage, which keeps every
// Hook on the same key in step within the tab. Keyed per lesson and nothing
// Else — as in legacy, the key carries no lesson version, so a buffer survives
// The author changing the starter code.
export function lessonCodeKey(courseSlug: string, lessonSlug: string): string {
  return `lesson-code-${courseSlug}-${lessonSlug}`;
}
