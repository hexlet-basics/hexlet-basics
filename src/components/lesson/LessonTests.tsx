import { CodeHighlight } from "@mantine/code-highlight";
import { Center, Stack, Text } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { getEditorSettings } from "@/lib/editor-languages";
import { toSupportedLanguage } from "@/lib/shiki";

// The lesson tests, readable at any time by anyone — deliberately, as in legacy:
// a learner who opens this tab is reading the specification of the exercise,
// which is the thing they are being asked to satisfy.
export default function LessonTests({
  courseSlug,
  testCode,
}: {
  courseSlug: string;
  testCode: string;
}) {
  const { t } = useTranslation();
  const { language } = getEditorSettings(courseSlug);

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
