import { Box, Center, rem } from "@mantine/core";
import { IconCircleCheck, IconLock } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import type { LessonProgressItem } from "@/client/types.gen";

// The mark beside a lesson's name: a check when it is finished, a lock when the
// gate has not opened it yet (ADR-0012), and an empty slot of the same width
// otherwise so the names line up.
//
// Shared by the player's lesson list and the course page, so the two screens
// mark the same lesson the same way from the same payload.
export default function LessonMark({ state }: { state: LessonProgressItem | undefined }) {
  const { t } = useTranslation();

  if (state?.finished) {
    return (
      <Center c="green" aria-label={t(($) => $.courses.lessons.show.finished)}>
        <IconCircleCheck size={16} />
      </Center>
    );
  }

  if (state && !state.available) {
    return (
      <Center c="dimmed" aria-label={t(($) => $.courses.lessons.show.locked)}>
        <IconLock size={16} />
      </Center>
    );
  }

  return <Box w={rem(16)} aria-hidden="true" />;
}
