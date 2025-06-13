import { useTranslation } from "react-i18next";

import { getEditorLanguage } from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import type { LessonSharedProps } from "../types.ts";
import { CodeHighlight } from "@mantine/code-highlight";

export default function TestsTab() {
  const { t: tCommon } = useTranslation("common");
  const { lesson, course, } = usePage<LessonSharedProps>().props;

  return (
    <>
      <p className="text-center lead">{tCommon("testInstructions")}</p>
      <CodeHighlight
        language={getEditorLanguage(course.slug!)}
        code={lesson.test_code ?? ''}
      />
    </>
  );
}
