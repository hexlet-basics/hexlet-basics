// @ts-check

import { AnsiUp } from "ansi_up";
import cn from "classnames";
import { useTranslation } from "react-i18next";
import { useSelector } from "react-redux";

import escape from "core-js/actual/escape.js";
import { getHelpByTutorUrl } from "../utils/languagesUtils.js";
import { checkInfoStates } from "../utils/maps.js";

const ansi = new AnsiUp();

function Output() {
  const checkInfo = useSelector((state) => state.checkInfoSlice);
  const { t } = useTranslation();

  if (checkInfoStates.checked !== checkInfo.processState) {
    return null;
  }

  const message = t(`check.${checkInfo.result}.message`, {
    url: getHelpByTutorUrl(language),
  });
  const messageForGuest = t("signInSuggestion");
  const alertClassName = cn("mt-auto alert mb-0 small p-2", {
    "alert-success": checkInfo.passed,
    "alert-warning": !checkInfo.passed,
  });
  // NOTE: исправление неверной кодировки для кириллицы
  // https://developer.mozilla.org/en-US/docs/Glossary/Base64
  const output = ansi.ansi_to_html(
    decodeURIComponent(escape(checkInfo.output)),
  );

  return (
    <div className="d-flex flex-column h-100">
      <pre>
        <code
          className="nohighlight"
          dangerouslySetInnerHTML={{ __html: output }}
        />
      </pre>
      <div
        className={alertClassName}
        dangerouslySetInnerHTML={{ __html: message }}
      />
      {!lessonMember.id && checkInfo.passed && (
        <div
          className="alert alert-warning mt-1 mb-0 small p-2"
          dangerouslySetInnerHTML={{ __html: messageForGuest }}
        />
      )}
    </div>
  );
}

export default Output;
