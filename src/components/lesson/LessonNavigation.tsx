import { Box, Center, rem } from "@mantine/core";
import { IconCircleCheck, IconLock } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import type { CourseLessonView } from "@/client/types.gen";
import { NavLink } from "@/components/RouterLink";

// The course's lessons, in course order, with the current one marked. Legacy's
// flat list plus the two marks sequential progression needs: a check on the
// finished ones and a lock on the ones the gate has not opened yet (ADR-0012).
//
// A locked entry is still a link. Theory is public and is what a search engine
// indexes, so the lock says "you cannot take this yet", never "you cannot read
// this" — it is the editor's actions that are gated, and the server decides
// that, not this list.
export default function LessonNavigation({ view }: { view: CourseLessonView }) {
  const { t } = useTranslation();
  const courseSlug = view.lesson.course.slug;

  // Names and order come from `lessons`, checks and locks from `progress`, joined
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
            leftSection={
              state?.finished ? (
                <Center c="green" aria-label={t(($) => $.courses.lessons.show.finished)}>
                  <IconCircleCheck size={16} />
                </Center>
              ) : state && !state.available ? (
                <Center c="dimmed" aria-label={t(($) => $.courses.lessons.show.locked)}>
                  <IconLock size={16} />
                </Center>
              ) : (
                <Box w={rem(16)} aria-hidden="true" />
              )
            }
          />
        );
      })}
    </Box>
  );
}
