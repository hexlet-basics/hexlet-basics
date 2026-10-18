import { CodeHighlight } from "@mantine/code-highlight";
import { Center, Stack, Text } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { getTestLanguage } from "@/lib/editor-languages";
import { toSupportedLanguage } from "@/lib/shiki";
import type { ReactElement } from "react";

// The lesson tests, readable at any time by anyone — deliberately, as in legacy:
// a learner who opens this tab is reading the specification of the exercise,
// which is the thing they are being asked to satisfy.
export default function LessonTests({
  courseSlug,
  testExtension,
  testCode,
}: {
  courseSlug: string;
  // The test file's own extension: the tests need not be in the course's language.
  testExtension: string | null;
  testCode: string;
}): ReactElement {
  const { t } = useTranslation();
  const language = getTestLanguage(courseSlug, testExtension);

  return (
    <Stack p="xs">
      <Center>
        <Text fz="sm">{t(($) => $.common.testInstructions)}</Text>
      </Center>
      <CodeHighlight
        code={testCode}
        language={toSupportedLanguage(language)}
        withCopyButton={false}
      />
    </Stack>
  );
}
