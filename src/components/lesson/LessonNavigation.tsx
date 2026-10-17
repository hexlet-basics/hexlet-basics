import { Box } from "@mantine/core";
import type { CourseLessonView } from "@/client/types.gen";
import LessonMark from "@/components/lesson/LessonMark";
import { NavLink } from "@/components/RouterLink";

// The course's lessons, in course order, with the current one marked and a check
// on the finished ones. Legacy's flat list: lessons may be taken in any order
// (ADR-0016), so every entry is a plain link.
export default function LessonNavigation({ view }: { view: CourseLessonView }) {
  const courseSlug = view.lesson.course.slug;

  // Names and order come from `lessons`, checks from `progress`, joined
  // by slug — the same pair the course page renders, so the two cannot disagree.
  const stateBySlug = new Map(view.progress?.lessons.map((item) => [item.slug, item]) ?? []);

  return (
    <Box p="lg">
      {view.lessons.map((item) => {
        const state = stateBySlug.get(item.slug);
        const current = item.slug === view.lesson.slug;

        return (
          <NavLink
            key={item.slug}
            to="/{-$locale}/languages/$slug/lessons/$lessonSlug"
            params={{ slug: courseSlug, lessonSlug: item.slug }}
            label={item.name}
            active={current}
            aria-current={current ? "page" : undefined}
            leftSection={<LessonMark state={state} />}
          />
        );
      })}
    </Box>
  );
}
