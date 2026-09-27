// Where a learner's buffer for a lesson is kept. The editor writes it and the
// assistant reads it, both through Mantine's useLocalStorage, which keeps every
// hook on the same key in step within the tab. Keyed per lesson and nothing
// else — as in legacy, the key carries no lesson version, so a buffer survives
// the author changing the starter code.
export function lessonCodeKey(courseSlug: string, lessonSlug: string): string {
  return `lesson-code-${courseSlug}-${lessonSlug}`;
}
