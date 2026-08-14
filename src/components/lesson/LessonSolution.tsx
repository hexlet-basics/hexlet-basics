import { CodeHighlight } from "@mantine/code-highlight";
import { Alert, Button, Center, Stack, Text, Title } from "@mantine/core";
import { useInterval } from "@mantine/hooks";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { getEditorSettings } from "@/lib/editor-languages";
import { toSupportedLanguage } from "@/lib/shiki";

// How long a learner who has not passed waits before the reference solution can
// be opened. Long enough that reading it is a last resort, short enough that
// nobody is stuck on one lesson forever.
const waitingTime = 20 * 60 * 1000;

// The author's answer, beside the learner's own so the two can be compared on
// one screen.
//
// It opens as soon as the lesson is passed. Before that it is behind a wait,
// with the time left on screen so a stuck learner knows whether to keep trying.
// The wait restarts on every page load, exactly as legacy's does — reproduced
// rather than quietly fixed, so this port has one variable and not several.
export default function LessonSolution({
  courseSlug,
  referenceSolution,
  learnerCode,
  unlocked,
}: {
  courseSlug: string;
  referenceSolution: string;
  learnerCode: string;
  unlocked: boolean;
}) {
  const { t } = useTranslation();
  const { language } = getEditorSettings(courseSlug);
  const codeLanguage = toSupportedLanguage(language);

  const [expiresAt] = useState(() => Date.now() + waitingTime);
  const [remaining, setRemaining] = useState(waitingTime);
  const [revealed, setRevealed] = useState(false);
  useInterval(() => setRemaining(Math.max(0, expiresAt - Date.now())), 1000, { autoInvoke: true });

  if (unlocked || revealed) {
    return (
      <Stack p="xs">
        <Title order={2}>{t(($) => $.common.teacherSolution)}</Title>
        <CodeHighlight code={referenceSolution} language={codeLanguage} withCopyButton={false} />

        <Title order={2}>{t(($) => $.common.userCode)}</Title>
        {learnerCode === "" ? (
          <Alert>{t(($) => $.common.userCodeInstructions)}</Alert>
        ) : (
          <CodeHighlight code={learnerCode} language={codeLanguage} withCopyButton={false} />
        )}
      </Stack>
    );
  }

  if (remaining > 0) {
    const seconds = Math.ceil(remaining / 1000);
    const clock = `${String(Math.floor(seconds / 60)).padStart(2, "0")}:${String(seconds % 60).padStart(2, "0")}`;

    return (
      <Stack align="center" p="xs">
        <Text size="lg" fw={500}>
          {t(($) => $.common.solutionInstructions)}
        </Text>
        <Text fz={50}>{clock}</Text>
      </Stack>
    );
  }

  return (
    <Stack p="xs">
      <Text>{t(($) => $.common.solutionNotice)}</Text>
      <Center>
        <Button variant="light" px="xl" onClick={() => setRevealed(true)}>
          {t(($) => $.common.showSolution)}
        </Button>
      </Center>
    </Stack>
  );
}
