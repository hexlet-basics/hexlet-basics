// @ts-check

import { AnsiUp } from "ansi_up";
import cn from "classnames";
import { useTranslation } from "react-i18next";

import escape from "core-js/actual/escape.js";
import { useAppSelector } from "../slices";
import { usePage } from "@inertiajs/react";
import type { Props } from "../types";

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

  const message = t(`check.${result}.message`, {
    // url: getHelpByTutorUrl(language),
  });

  const messageForGuest = tCommon("signInSuggestion");
  const alertClassName = cn("mt-auto alert mb-0 small p-2", {
    "alert-success": passed,
    "alert-warning": !passed,
  });
  // NOTE: исправление неверной кодировки для кириллицы
  // https://developer.mozilla.org/en-US/docs/Glossary/Base64
  const outputAsHTML = ansi.ansi_to_html(decodeURIComponent(escape(output)));

  return (
    <div className="d-flex flex-column h-100">
      <pre>
        <code
          className="nohighlight"
          dangerouslySetInnerHTML={{ __html: outputAsHTML }}
        />
      </pre>
      <div
        className={alertClassName}
        dangerouslySetInnerHTML={{ __html: message }}
      />
      {!lessonMember && passed && (
        <div
          className="alert alert-warning mt-1 mb-0 small p-2"
          dangerouslySetInnerHTML={{ __html: messageForGuest }}
        />
      )}
    </div>
  );
}
