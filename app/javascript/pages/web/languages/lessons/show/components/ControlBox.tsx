import { Link, usePage } from "@inertiajs/react";
import { Box, Button, Divider, Group, Stack, Text } from "@mantine/core";
import { modals } from "@mantine/modals";
import { notifications } from "@mantine/notifications";
import {
  IconPlayerPlay,
  IconRepeat,
  IconThumbDown,
  IconThumbUp,
} from "@tabler/icons-react";
import { noop } from "es-toolkit";
import { useHotkeys } from "react-hotkeys-hook";
import { useTranslation } from "react-i18next";
import * as Routes from "@/routes.js";
import type { LessonSharedProps } from "@/types";
import { useLessonStore } from "../store.tsx";

export default function ControlBox() {
  const { props: sharedProps, url } = usePage<LessonSharedProps>();

  const {
    lesson,
    course,
    prevLesson,
    landingPage,
    nextLesson,
    auth: { user },
  } = sharedProps;

  const { t } = useTranslation();

  const processState = useLessonStore((state) => state.processState);
  const finished = useLessonStore((state) => state.finished);
  const runCheck = useLessonStore((state) => state.runCheck);

  const handleRunCheck = async () => {
    const result = await runCheck({ course, lesson });
    if (!result) {
      notifications.show({
        // title: 'Default notification',
        message: t(($) => $.common.errors.network),
      });
    }
  };

  const resetContent = useLessonStore((state) => state.resetContent);
  // const confirmResetting = useConfirmation({ callback: resetContent });
  const openModal = () =>
    modals.openConfirmModal({
      title: t(($) => $.common.confirm),
      labels: { confirm: <IconThumbUp />, cancel: <IconThumbDown /> },
      onCancel: noop,
      onConfirm: resetContent,
    });

  const isCodeChecking = processState === "checking";

  useHotkeys("ctrl+enter", handleRunCheck);

  const nextButtonDisabled = !finished;

  return (
    <Stack gap={0}>
      <Divider />
      <Box py="sm">
        <Group justify="center">
          <Button size="xs" me="xs" variant="light" onClick={openModal}>
            <IconRepeat size={18} />
          </Button>
          {/* <Popover width={200} position="top" withArrow shadow="md"> */}
          {/*   <Popover.Target> */}
          {/*     <Button size="xs" onClick={openModal}> */}
          {/*       <IconRepeat size={18} /> */}
          {/*     </Button> */}
          {/*   </Popover.Target> */}
          {/*   <Popover.Dropdown> */}
          {/*     <Text size="sm" fw={500}> */}
          {/*       {t('languages.lessons.show.controls.header')} */}
          {/*     </Text> */}
          {/*     <Text size="xs">{t('languages.lessons.show.controls.body')}</Text> */}
          {/*   </Popover.Dropdown> */}
          {/* </Popover> */}

          <Button
            component={Link}
            disabled={!prevLesson}
            size="xs"
            variant="outline"
            color="green"
            me="xs"
            href={
              prevLesson
                ? Routes.language_lesson_path(course.slug!, prevLesson.slug!)
                : "#"
            }
          >
            {t(($) => $.languages.lessons.show.prev)}
          </Button>

          <Button
            size="xs"
            // variant="outline"
            // color="blue"
            me="xs"
            loading={isCodeChecking}
            disabled={isCodeChecking}
            onClick={handleRunCheck}
          >
            <Group gap={4} align="center">
              <IconPlayerPlay size={18} />
              <Text size="xs">
                {t(($) => $.languages.lessons.show.controls.run)}
              </Text>
            </Group>
          </Button>

          {!user && (
            <Button
              component={Link}
              size="xs"
              variant="outline"
              color="green"
              me="xs"
              disabled={nextButtonDisabled}
              href={Routes.new_user_path({ demo: true, from: url })}
            >
              {t(($) => $.languages.lessons.show.next)}
            </Button>
          )}

          {user && nextLesson && (
            <Button
              component={Link}
              size="xs"
              variant="outline"
              color="green"
              disabled={nextButtonDisabled}
              href={Routes.language_lesson_path(course.slug!, nextLesson.slug!)}
            >
              {t(($) => $.languages.lessons.show.next)}
            </Button>
          )}

          {user && !nextLesson && (
            <Button
              component={Link}
              size="xs"
              variant="outline"
              color="green"
              disabled={nextButtonDisabled}
              href={Routes.success_language_path(landingPage.slug!)}
            >
              {t(($) => $.languages.lessons.show.finish)}
            </Button>
          )}
        </Group>
      </Box>
    </Stack>
  );
}
