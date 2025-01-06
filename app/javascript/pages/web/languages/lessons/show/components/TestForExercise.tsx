import { useTranslation } from "react-i18next";
import SyntaxHighlighter from "react-syntax-highlighter";

import XssContent from "@/components/XssContent.js";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types.js";
import { getLanguage } from "../utils/editorUtils.js";

export default function TestForExercise() {
  const { t } = useTranslation();
  const { lesson, course } = usePage<Props>().props;

  return (
    <div>
      <p className="text-center lead">
        {t("languages.lessons.show.instructions")}
      </p>
      <pre>
        <code>
          <SyntaxHighlighter language={course.slug}>
            {lesson.test_code}
          </SyntaxHighlighter>
        </code>
      </pre>
    </div>
  );
}
