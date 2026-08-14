import { CodeHighlight } from "@mantine/code-highlight";
import { Alert, Box, Text } from "@mantine/core";
import { IconAlertTriangle, IconCheck } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import type { LessonCheckingResponse } from "@/client/types.gen";
import { plainTextLanguage } from "@/lib/shiki";

// What the run said. A solution that never terminated is its own outcome, not an
// ordinary failure: it is a loop to find, not a wrong answer to rethink, and the
// copy for each of the three comes from `common.check` keyed by the result.
//
// The output is rendered exactly as it arrived. The runner scrubs it server-side
// (ADR-0013), so nothing here decodes or unescapes it — a departure from legacy,
// which base64-encoded it.
//
// Between pressing run and the answer there is nothing here, as in legacy: the
// run button's own spinner is the progress signal.
export default function LessonOutput({ result }: { result: LessonCheckingResponse | null }) {
  const { t } = useTranslation();

  if (!result) return null;

  return (
    <Box p="xs">
      <Alert
        color={result.passed ? "green" : "yellow"}
        icon={result.passed ? <IconCheck size={16} /> : <IconAlertTriangle size={16} />}
        title={t(($) => $.common.check[result.result].headline)}
        fz="sm"
        variant="light"
        radius="sm"
        mb="xs"
      >
        <Text fz="sm">{t(($) => $.common.check[result.result].message)}</Text>
      </Alert>

      <CodeHighlight code={result.output} language={plainTextLanguage} withCopyButton={false} />
    </Box>
  );
}
