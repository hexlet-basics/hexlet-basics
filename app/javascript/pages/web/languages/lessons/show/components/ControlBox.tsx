import { deleteFromStorage } from "@rehooks/local-storage";
import React from "react";
import { useTranslation } from "react-i18next";

import cn from "classnames";
import { Button, OverlayTrigger, Popover, Spinner } from "react-bootstrap";
// import Hotkeys from 'react-hot-keys';
import { useHotkeys } from "react-hotkeys-hook";

import slice, { runCheck } from "../slices/RootSlice.ts";
import * as Routes from "@/routes.js";
import { getKeyForStoringLessonCode } from "@/lib/utils.ts";
import { Link, usePage } from "@inertiajs/react";
import type { Props } from "../types.ts";
import { useAppDispatch, useAppSelector } from "../slices/index.ts";

export default function ControlBox() {
  const {
    course,
    prevLesson,
    nextLesson,
    auth: { user },
  } = usePage<Props>().props;

  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");

  const { lesson } = usePage<Props>().props;

  const processState = useAppSelector((state) => state.processState);
  const finished = useAppSelector((state) => state.finished);

  const dispatch = useAppDispatch();
  const handleRunCheck = () => {
    console.log('jopa')
    dispatch(runCheck(lesson));
  };

  const handleReset = () => {
    // NOTE: easier than state manipulating. dont touch, dont blame, be happy.

    if (window.confirm(tCommon("confirm"))) {
      deleteFromStorage(getKeyForStoringLessonCode(lesson));
      window.location.reload();
    }
  };

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

  const prevButtonClasses = cn(
    "btn btn-sm btn-outline-success me-3",
  );

  const nextButtonClasses = cn(
    "btn btn-sm btn-outline-success fw-normal",
    {
      disabled: !finished || !nextLesson,
    },
  );

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
            onClick={handleReset}
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
        {nextLesson && (
          <Link
            className={nextButtonClasses}
            href={Routes.language_lesson_path(course.slug!, nextLesson.slug!)}
          >
            {t("languages.lessons.show.next")}
          </Link>
        )}
      </div>
    </div>
  );
}
