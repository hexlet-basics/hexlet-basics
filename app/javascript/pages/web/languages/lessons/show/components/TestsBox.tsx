import { useTranslation } from "react-i18next";
import SyntaxHighlighter from "react-syntax-highlighter";
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs';

// import XssContent from "@/components/XssContent.js";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types.js";
import { getEditorLanguage } from "@/lib/utils.js";

export default function TestsBox() {
  const { t: tCommon } = useTranslation("common");
  const { lesson, course } = usePage<Props>().props;

  return (
    <div className="hexlet-basics-content">
      <p className="text-center lead">{tCommon("testInstructions")}</p>
      <SyntaxHighlighter
        style={github}
        showLineNumbers
        language={getEditorLanguage(course.slug!)}
      >
        {lesson.test_code}
      </SyntaxHighlighter>
    </div>
  );
}
