import * as Routes from "@/routes.js";
import { AnsiUp } from "ansi_up";
import { useTranslation } from "react-i18next";

import XssContent from "@/components/XssContent";
import { usePage } from "@inertiajs/react";
import { Alert } from "react-bootstrap";
import type { LessonSharedProps } from "../types.ts";
import { useLessonStore } from "../store.tsx";

const ansi = new AnsiUp();

export default function OutputTab() {
  const { lessonMember } = usePage<LessonSharedProps>().props;
  const result = useLessonStore((state) => state.result);
  const processState = useLessonStore((state) => state.processState);
  const output = useLessonStore((state) => state.output);
  const passed = useLessonStore((state) => state.passed);
  const { t: tCommon } = useTranslation("common");

  if (processState !== "checked") {
    return null;
  }

  const message = tCommon(`check.${result!}.message`);

  const messageForGuest = tCommon("signInSuggestion", {
    url: Routes.new_user_path(),
  });
  // NOTE: исправление неверной кодировки для кириллицы
  // https://developer.mozilla.org/en-US/docs/Glossary/Base64
  const outputAsHTML = ansi.ansi_to_html(decodeURIComponent(output));

  return (
    <div className="d-flex flex-column h-100">
      <pre>
        <code>
          <XssContent>{outputAsHTML}</XssContent>
        </code>
      </pre>
      <Alert
        variant={passed ? "success" : "warning"}
        className="mt-auto small border-0 p-2"
      >
        <XssContent>{message}</XssContent>
      </Alert>
      {/* {!lessonMember && passed && ( */}
      {/*   <Alert variant="primary" className="small border-0 p-2"> */}
      {/*     <XssContent>{messageForGuest}</XssContent> */}
      {/*   </Alert> */}
      {/* )} */}
    </div>
  );
}
