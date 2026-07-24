import { usePage } from "@inertiajs/react";
import { CodeHighlight } from "@mantine/code-highlight";
import {
  Alert,
  Box,
  Button,
  Center,
  Image,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import dayjs from "dayjs";
import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import useCountdown from "@/hooks/useCountdown.ts";
import waitingClock from "@/images/waiting_clock.webp";
import { toSupportedLanguage } from "@/lib/shiki";
import { getEditorLanguage } from "@/lib/utils.ts";
import type { LessonSharedProps } from "@/types";
import { useLessonStore } from "../store.tsx";

const waitingTime = 20 * 60 * 1000; // 20 min
// const waitingTime = 3000;

type CountdownProps = {
  expiryTimestamp: Date;
  canBeShown: boolean;
  renderShowButton: () => ReactNode;
};

function Countdown({
  expiryTimestamp,
  canBeShown,
  renderShowButton,
}: CountdownProps) {
  const { t } = useTranslation();
  const countdown = useCountdown(expiryTimestamp);

  if (!countdown.isRunning || canBeShown) {
    return renderShowButton();
  }

  const remainingTime = dayjs
    .duration(countdown.remainingSeconds, "seconds")
    .format("mm:ss");

  return (
    <Stack align="center">
      <Text size="lg" fw={500}>
        {t(($) => $.common.solutionInstructions)}
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
}

export default function SolutionTab() {
  const { t } = useTranslation();
  const { course, lesson } = usePage<LessonSharedProps>().props;
  const content = useLessonStore((state) => state.content);
  const finished = useLessonStore((state) => state.finished);
  const solutionState = useLessonStore((state) => state.solutionState);
  const startTime = useLessonStore((state) => state.startTime);

  const changeSolutionState = useLessonStore(
    (state) => state.changeSolutionState,
  );

  const expiryTimestamp = new Date(startTime + waitingTime);

  const renderUserCode = () => {
    if (content === "") {
      return <Alert>{t(($) => $.common.userCodeInstructions)}</Alert>;
    }

    return (
      <CodeHighlight
        code={content}
        language={toSupportedLanguage(getEditorLanguage(course.slug!))}
        withCopyButton={false}
      />
    );
  };

  const renderSolution = () => {
    return (
      <>
        <Stack mb="xl">
          <Title order={2} mb="lg">
            {t(($) => $.common.teacherSolution)}
          </Title>
          <CodeHighlight
            code={lesson.original_code!}
            language={toSupportedLanguage(getEditorLanguage(course.slug!))}
            withCopyButton={false}
          />
        </Stack>
        <Stack>
          <Title order={2} mb="lg">
            {t(($) => $.common.userCode)}
          </Title>
          {renderUserCode()}
        </Stack>
      </>
    );
  };

  const handleShowSolution = () => {
    changeSolutionState("shown");
  };

  const renderShowButton = () => (
    <>
      <Text>{t(($) => $.common.solutionNotice)}</Text>
      <Center>
        <Button variant="light" onClick={handleShowSolution} px="xl">
          {t(($) => $.common.showSolution)}
        </Button>
      </Center>
    </>
  );

  return (
    <Box>
      {finished || solutionState === "shown" ? (
        renderSolution()
      ) : (
        <Countdown
          expiryTimestamp={expiryTimestamp}
          canBeShown={solutionState === "canBeShown"}
          renderShowButton={renderShowButton}
        />
      )}
    </Box>
  );
}
