import { usePage } from "@inertiajs/react";
import { CodeHighlight } from "@mantine/code-highlight";
import { Center, Stack, Text } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { toSupportedLanguage } from "@/lib/shiki";
import { getEditorLanguage } from "@/lib/utils.ts";
import type { LessonSharedProps } from "@/types";

export default function TestsTab() {
  const { t } = useTranslation();

  const { lesson, course } = usePage<LessonSharedProps>().props;

  return (
    <Stack>
      <Center>
        <Text>{t(($) => $.common.testInstructions)}</Text>
      </Center>
      <CodeHighlight
        language={toSupportedLanguage(getEditorLanguage(course.slug!))}
        code={lesson.test_code ?? ""}
        withCopyButton={false}
      />
    </Stack>
  );
}
