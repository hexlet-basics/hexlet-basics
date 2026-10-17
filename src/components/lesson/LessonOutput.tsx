import { CodeHighlight } from "@mantine/code-highlight";
import { Alert, Anchor, Box, Center, Loader, Stack, Text } from "@mantine/core";
import { IconAlertTriangle, IconCheck, IconInfoCircle } from "@tabler/icons-react";
import { useLocation, useRouteContext } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import type { LessonCheckingResponse } from "@/client/types.gen";
import { TextLink } from "@/components/RouterLink";
import { ansiLanguage } from "@/lib/shiki";

// What the run said. A solution that never terminated is its own outcome, not an
// ordinary failure: it is a loop to find, not a wrong answer to rethink, and the
// copy for each of the three comes from `common.check` keyed by the result.
//
// The output is rendered exactly as it arrived, highlighted as ANSI so the test
// runner's colours show instead of their escape codes. The runner scrubs it server-side
// (ADR-0013), so nothing here decodes or unescapes it — a departure from legacy,
// which base64-encoded it.
//
// The run can take a while (a container starts per check), and the learner is
// switched to this tab the moment they press run. The tab therefore says what is
// going on — a spinner while the check runs, a hint before the first one — rather
// than showing an empty pane the run button's spinner cannot explain from here.
//
// The result copy carries markup (`<a>` in the failed message), so it goes through
// Trans; the help article it links to is the one legacy linked.
// The line a guest reads after passing: signing up keeps what they have done.
function GuestPrompt() {
  const { t } = useTranslation();
  const redirect = useLocation({ select: (location) => location.href });

  return (
    <Alert
      color="blue"
      icon={<IconInfoCircle size={16} />}
      fz="sm"
      variant="light"
      radius="sm"
      mb="xs"
    >
      <Trans
        t={t}
        i18nKey={($) => $.courses.lessons.show.sign_up_for_tracking_progress_html}
        components={{ a: <TextLink to="/{-$locale}/users/new" search={{ redirect }} inherit /> }}
      />
    </Alert>
  );
}

const FIND_ERRORS_URL = "https://help.hexlet.io/article/20623";

export default function LessonOutput({
  result,
  running = false,
}: {
  result: LessonCheckingResponse | null;
  running?: boolean;
}) {
  const { t } = useTranslation();
  const user = useRouteContext({ from: "__root__", select: (context) => context.user });

  if (running) {
    return (
      <Center h="100%" p="xl">
        <Stack align="center" gap="xs">
          <Loader size="sm" />
          <Text fz="sm" c="dimmed">
            {t(($) => $.courses.lessons.show.output_running)}
          </Text>
        </Stack>
      </Center>
    );
  }

  if (!result) {
    return (
      <Text p="md" fz="sm" c="dimmed">
        {t(($) => $.courses.lessons.show.output_empty)}
      </Text>
    );
  }

  return (
    <Box
      p="xs"
      mih="100%"
      display="flex"
      style={{ flexDirection: "column", boxSizing: "border-box" }}
    >
      <Alert
        color={result.passed ? "green" : "yellow"}
        icon={result.passed ? <IconCheck size={16} /> : <IconAlertTriangle size={16} />}
        title={t(($) => $.common.check[result.result].headline)}
        fz="sm"
        variant="light"
        radius="sm"
        mb="xs"
      >
        <Text fz="sm">
          <Trans
            t={t}
            i18nKey={($) => $.common.check[result.result].message}
            components={{
              a: <Anchor href={FIND_ERRORS_URL} target="_blank" rel="noreferrer" inherit />,
            }}
          />
        </Text>
      </Alert>

      {/* The one thing this port adds rather than copies: legacy carried the
          string and rendered it nowhere. It sits with the result it is about, not
          under the controls, where its appearing moved the buttons — and in an
          alert of its own, because inside the success one it read as more of the
          same message. */}
      {result.passed && !user && <GuestPrompt />}

      <CodeHighlight code={result.output} language={ansiLanguage} withCopyButton={false} flex={1} />
    </Box>
  );
}
