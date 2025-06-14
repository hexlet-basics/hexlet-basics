import { useTranslation } from "react-i18next";
import SyntaxHighlighter from "react-syntax-highlighter";
import { github } from "react-syntax-highlighter/dist/esm/styles/hljs";

import { getEditorLanguage } from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import type { LessonSharedProps } from "../types.ts";

export default function TestsTab() {
  const { t: tCommon } = useTranslation("common");
  const { lesson, course } = usePage<LessonSharedProps>().props;

  return (
    <div className="hexlet-basics-content">
      <p className="text-center lead">{tCommon("testInstructions")}</p>
      <SyntaxHighlighter
        style={github}
        showLineNumbers
        language={getEditorLanguage(course.slug!)}
      >
        {lesson.test_code!}
      </SyntaxHighlighter>
    </div>
  );
}
