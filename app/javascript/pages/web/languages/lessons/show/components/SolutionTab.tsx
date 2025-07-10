import { usePage } from '@inertiajs/react';
import { CodeHighlight } from '@mantine/code-highlight';
import {
  Alert,
  Box,
  Button,
  Center,
  Group,
  Image,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import dayjs from 'dayjs';
import { useTranslation } from 'react-i18next';
import { useTimer } from 'react-timer-hook';
import waitingClock from '@/images/waiting_clock.png';
import { getEditorLanguage } from '@/lib/utils.ts';
import { useLessonStore } from '../store.tsx';
import type { LessonSharedProps } from '../types.ts';

const waitingTime = 20 * 60 * 1000; // 20 min
// const waitingTime = 3000;

export default function SolutionTab() {
  const { course, lesson } = usePage<LessonSharedProps>().props;
  const content = useLessonStore((state) => state.content);
  const finished = useLessonStore((state) => state.finished);
  const solutionState = useLessonStore((state) => state.solutionState);
  const startTime = useLessonStore((state) => state.startTime);
  const { t: tCommon } = useTranslation('common');
  const changeSolutionState = useLessonStore(
    (state) => state.changeSolutionState,
  );

  const expiryTimestamp = new Date(startTime + waitingTime);
  const timerData = useTimer({ expiryTimestamp });
  // console.log(timerData.isRunning, timerData.totalSeconds, expiryTimestamp)

  const renderUserCode = () => {
    if (content === '') {
      return <Alert>{tCommon('userCodeInstructions')}</Alert>;
    }

    return (
      <CodeHighlight
        code={content}
        language={getEditorLanguage(course.slug!)}
      />
    );
  };

  const renderSolution = () => {
    return (
      <>
        <Stack mb="xl">
          <Title order={2} mb="lg">
            {tCommon('teacherSolution')}
          </Title>
          <CodeHighlight
            code={lesson.original_code!}
            language={getEditorLanguage(course.slug!)}
          />
        </Stack>
        <Stack>
          <Title order={2} mb="lg">
            {tCommon('userCode')}
          </Title>
          {renderUserCode()}
        </Stack>
      </>
    );
  };

  const handleShowSolution = () => {
    changeSolutionState('shown');
  };

  const renderShowButton = () => (
    <>
      <Text>{tCommon('solutionNotice')}</Text>
      <Center>
        <Button variant="light" onClick={handleShowSolution} px="xl">
          {tCommon('showSolution')}
        </Button>
      </Center>
    </>
  );

  const Countdown = () => {
    const { isRunning, totalSeconds } = timerData;

    if (!isRunning || solutionState === 'canBeShown') {
      return renderShowButton();
    }
    const remainingTime = dayjs
      .duration(totalSeconds, 'seconds')
      .format('mm:ss');

    return (
      <Stack align="center">
        <Text size="lg" fw={500}>
          {tCommon('solutionInstructions')}
        </Text>
        <Text fz={50}>{remainingTime}</Text>
        <Image
          src={waitingClock}
          fit="contain"
          loading="lazy"
          alt="waiting_clock"
        />
      </Stack>
    );
  };

  return (
    <Box>
      {finished || solutionState === 'shown' ? renderSolution() : <Countdown />}
    </Box>
  );
}
