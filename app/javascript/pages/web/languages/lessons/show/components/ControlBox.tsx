import { Link, usePage } from '@inertiajs/react';
import { Box, Button, Group, Loader, Text } from '@mantine/core';
import { modals } from '@mantine/modals';
import { notifications } from '@mantine/notifications';
import { noop } from 'es-toolkit';
import { Play, Repeat, ThumbsDown, ThumbsUp } from 'lucide-react';
import { useHotkeys } from 'react-hotkeys-hook';
import { useTranslation } from 'react-i18next';
import * as Routes from '@/routes.js';
import { useLessonStore } from '../store.tsx';
import type { LessonSharedProps } from '../types.ts';

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
  const { t: tCommon } = useTranslation('common');

  const processState = useLessonStore((state) => state.processState);
  const finished = useLessonStore((state) => state.finished);
  const runCheck = useLessonStore((state) => state.runCheck);

  const handleRunCheck = async () => {
    const result = await runCheck({ course, lesson });
    if (!result) {
      notifications.show({
        // title: 'Default notification',
        message: tCommon(($) => $.errors.network),
      });
    }
  };

  const resetContent = useLessonStore((state) => state.resetContent);
  // const confirmResetting = useConfirmation({ callback: resetContent });
  const openModal = () =>
    modals.openConfirmModal({
      title: tCommon(($) => $.confirm),
      labels: { confirm: <ThumbsUp />, cancel: <ThumbsDown /> },
      onCancel: noop,
      onConfirm: resetContent,
    });

  const isCodeChecking = processState === 'checking';

  useHotkeys('ctrl+enter', handleRunCheck);

  const renderRunButtonContent = () => {
    const text = t(($) => $.languages.lessons.show.controls.run);
    if (isCodeChecking) {
      return (
        <Group gap={4} align="center">
          <Loader size="xs" />
          <Text size="xs">{text}</Text>
        </Group>
      );
    }

    return (
      <Group gap={4} align="center">
        <Play size={18} />
        <Text size="xs">{text}</Text>
      </Group>
    );
  };

  const nextButtonDisabled = !finished;

  return (
    <Box py="sm" className="border-t border-gray-400">
      <Group justify="center">
        <Button size="xs" me="xs" onClick={openModal}>
          <Repeat size={18} />
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

        {prevLesson && (
          <Button
            component={Link}
            size="xs"
            variant="outline"
            color="green"
            radius="sm"
            me="xs"
            href={Routes.language_lesson_path(course.slug!, prevLesson.slug!)}
          >
            {t(($) => $.languages.lessons.show.prev)}
          </Button>
        )}

        <Button
          size="xs"
          variant="outline"
          color="blue"
          radius="sm"
          me="xs"
          disabled={isCodeChecking}
          onClick={handleRunCheck}
        >
          {renderRunButtonContent()}
        </Button>

        {user.guest && (
          <Button
            component={Link}
            size="xs"
            variant="outline"
            color="green"
            radius="sm"
            me="xs"
            disabled={nextButtonDisabled}
            href={Routes.new_user_path({ demo: true, from: url })}
          >
            {t(($) => $.languages.lessons.show.next)}
          </Button>
        )}

        {!user.guest && nextLesson && (
          <Button
            component={Link}
            size="xs"
            variant="outline"
            color="green"
            radius="sm"
            disabled={nextButtonDisabled}
            href={Routes.language_lesson_path(course.slug!, nextLesson.slug!)}
          >
            {t(($) => $.languages.lessons.show.next)}
          </Button>
        )}

        {!user.guest && !nextLesson && (
          <Button
            component={Link}
            size="xs"
            variant="outline"
            color="green"
            radius="sm"
            disabled={nextButtonDisabled}
            href={Routes.success_language_url(landingPage.slug!)}
          >
            {t(($) => $.languages.lessons.show.finish)}
          </Button>
        )}
      </Group>
    </Box>
  );
}
