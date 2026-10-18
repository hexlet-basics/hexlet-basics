import { Box, Center, rem } from "@mantine/core";
import { IconCircleCheck } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";

// The mark beside a lesson's name: a check when it is finished, and an empty slot
// of the same width otherwise so the names line up. Lessons may be taken in any
// order (ADR-0016), so nothing marks one as out of reach.
//
// Shared by the player's lesson list and the course page, so the two screens
// mark the same lesson the same way from the same payload.
export default function LessonMark({ finished }: { finished: boolean }) {
  const { t } = useTranslation();

  if (finished) {
    return (
      <Center c="green" aria-label={t(($) => $.courses.lessons.show.finished)}>
        <IconCircleCheck size={16} />
      </Center>
    );
  }

  return <Box w={rem(16)} aria-hidden="true" />;
}
