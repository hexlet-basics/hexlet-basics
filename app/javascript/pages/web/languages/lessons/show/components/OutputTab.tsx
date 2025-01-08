// @ts-check

import { AnsiUp } from "ansi_up";
import cn from "classnames";
import { useTranslation } from "react-i18next";

import escape from "core-js/actual/escape.js";
import { useAppSelector } from "../slices";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types";
import XssContent from "@/components/XssContent";
import { Alert } from "react-bootstrap";

const ansi = new AnsiUp();

export default function OutputTab() {
  const { lessonMember } = usePage<Props>().props;
  const result = useAppSelector((state) => state.result);
  const processState = useAppSelector((state) => state.processState);
  const output = useAppSelector((state) => state.output);
  const passed = useAppSelector((state) => state.passed);
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");

  if (processState !== "checked") {
    return null;
  }

  const message = tCommon(`check.${result}.message`);

  const messageForGuest = tCommon("signInSuggestion");
  // NOTE: исправление неверной кодировки для кириллицы
  // https://developer.mozilla.org/en-US/docs/Glossary/Base64
  const outputAsHTML = ansi.ansi_to_html(decodeURIComponent(escape(output)));

  return (
    <div className="d-flex flex-column h-100">
      <pre>
        <code>
          <XssContent>{outputAsHTML}</XssContent>
        </code>
      </pre>
      <Alert
        variant={passed ? "success" : "warning"}
        className="mt-auto small p-2"
      >
        <XssContent>{message}</XssContent>
      </Alert>
      {!lessonMember && passed && (
        <Alert variant="warning" className="mt-auto small p-2">
          <XssContent>{messageForGuest}</XssContent>
        </Alert>
      )}
    </div>
  );
}
