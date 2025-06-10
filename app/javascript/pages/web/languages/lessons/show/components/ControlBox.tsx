import { useTranslation } from "react-i18next";

import cn from "classnames";
import { Button, OverlayTrigger, Popover, Spinner } from "react-bootstrap";
import { useHotkeys } from "react-hotkeys-hook";

import useConfirmation from "@/hooks/useConfirmation.ts";
import * as Routes from "@/routes.js";
import { Link, usePage } from "@inertiajs/react";
import type { LessonSharedProps } from "../types.ts";
import { useLessonStore } from "../store.tsx";
import { enqueueSnackbar } from "notistack";

export default function ControlBox() {
  const {
    props: sharedProps,
    url,
  } = usePage<LessonSharedProps>();

  const {
    lesson,
    course,
    prevLesson,
    landingPage,
    nextLesson,
    auth: { user },
  } = sharedProps

  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");

  const processState = useLessonStore((state) => state.processState);
  const finished = useLessonStore((state) => state.finished);
  const runCheck = useLessonStore((state) => state.runCheck);

  const handleRunCheck = async () => {
    const result = await runCheck({ course, lesson });
    if (!result) {
      enqueueSnackbar(tCommon('errors.network'));
    }
  };

  const resetContent = useLessonStore((state) => state.resetContent);
  const confirmResetting = useConfirmation({
    callback: resetContent
  });

  const isCodeChecking = processState === "checking";

  const renderRunButtonContent = () => {
    const text = t("languages.lessons.show.controls.run");
    if (isCodeChecking) {
      return (
        <>
          <Spinner
            as="span"
            animation="border"
            size="sm"
            role="status"
            aria-hidden="true"
          />
          <span className="visually-hidden">{tCommon("loading")}</span>
          <span className="d-none d-sm-block d-md-none d-lg-block ms-1">
            {text}
          </span>
        </>
      );
    }

    return (
      <>
        <span className="bi bi-play-circle" />
        <span className="d-none d-sm-block d-md-none d-lg-block ms-1">
          {text}
        </span>
      </>
    );
  };

  const prevButtonClasses = cn("btn btn-sm btn-outline-success me-3");

  const nextButtonClasses = cn("btn btn-sm btn-outline-success fw-normal", {
    disabled: !finished,
  });

  useHotkeys("ctrl+enter", handleRunCheck);

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">
        {t("languages.lessons.show.controls.header")}
      </Popover.Header>
      <Popover.Body>{t("languages.lessons.show.controls.body")}</Popover.Body>
    </Popover>
  );

  return (
    <div className="d-flex p-3 border-top">
      {/* <OverlayTrigger trigger="click" placement="top" overlay={popover}>
        <Button variant="link" className="text-muted my-auto me-3 p-1">
          <span aria-label="help" className="bi bi-question-circle fs-5" />
        </Button>
      </OverlayTrigger> */}
      <div className="m-auto">
        <OverlayTrigger
          trigger={["hover", "focus"]}
          placement="top"
          overlay={popover}
        >
          <Button
            variant="outline-secondary"
            size="sm"
            className="me-3"
            onClick={confirmResetting}
          // title={t('resetCode')}
          >
            <span className="bi bi-arrow-repeat" />
          </Button>
        </OverlayTrigger>

        {prevLesson && (
          <Link
            href={Routes.language_lesson_path(course.slug!, prevLesson.slug!)}
            className={prevButtonClasses}
          >
            {t("languages.lessons.show.prev")}
          </Link>
        )}
        <Button
          variant="outline-primary"
          size="sm"
          className="me-3 d-inline-flex align-items-center"
          onClick={handleRunCheck}
          disabled={isCodeChecking}
        >
          {renderRunButtonContent()}
        </Button>
        {user.guest && (
          <Link
            className={nextButtonClasses}
            href={Routes.new_user_path({ demo: true, from: url })}
          >
            {t("languages.lessons.show.next")}
          </Link>
        )}
        {!user.guest && nextLesson && (
          <Link
            className={nextButtonClasses}
            href={Routes.language_lesson_path(course.slug!, nextLesson.slug!)}
          >
            {t("languages.lessons.show.next")}
          </Link>
        )}
        {!user.guest && !nextLesson && (
          <Link
            className={nextButtonClasses}
            href={Routes.success_language_url(landingPage.slug!)}
          >
            {t("languages.lessons.show.finish")}
          </Link>
        )}
      </div>
    </div>
  );
}
