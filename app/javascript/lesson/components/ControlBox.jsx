// @ts-check

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import { deleteFromStorage } from '@rehooks/local-storage';

import cn from 'classnames';
// import Hotkeys from 'react-hot-keys';
import { useHotkeys } from 'react-hotkeys-hook';
import {
  Button, Spinner, Popover, OverlayTrigger,
} from 'react-bootstrap';

import Routes from '../utils/configured_routes.js';
import { actions } from '../slices/index.js';
import { checkInfoStates } from '../utils/maps.js';
import EntityContext from '../EntityContext.js';

function ControlBox() {
  const { t } = useTranslation();
  const { checkInfo, lessonInfo } = useSelector((state) => ({
    checkInfo: state.checkInfoSlice,
    lessonInfo: state.lessonSlice,
  }));

  const dispatch = useDispatch();
  const {
    lessonVersion, language, lesson,
  } = useContext(EntityContext);

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lessonVersion }));
  };

  const handleReset = () => {
    // NOTE: easier than state manipulating. dont touch, dont blame, be happy.
    // eslint-disable-next-line no-alert
    if (window.confirm(t('confirm'))) {
      const localStorageKey = `lesson-version-${lessonVersion.id}`;
      deleteFromStorage(localStorageKey);
      window.location.reload();
    }
  };

  const isCodeChecking = checkInfo.processState === checkInfoStates.checking;

  const renderRunButtonContent = () => {
    const text = t('run');
    if (isCodeChecking) {
      return (
        <>
          <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true" />
          <span className="visually-hidden">{t('loading')}</span>
          <span className="d-none d-sm-block d-md-none d-lg-block ms-1">{text}</span>
        </>
      );
    }

    return (
      <>
        <span className="bi bi-play-circle" />
        <span className="d-none d-sm-block d-md-none d-lg-block ms-1">{text}</span>
      </>
    );
  };

  const prevButtonClasses = cn(`btn btn-outline-secondary
    fw-normal me-3 order-first order-sm-0 order-md-first order-lg-0`);

  const nextButtonClasses = cn('btn btn-outline-primary fw-normal', {
    disabled: !lessonInfo.finished,
  });

  const nextLessonPath = lessonInfo.finished ? Routes.nextLessonLanguageLessonPath(language, lesson.slug) : null;
  const prevLessonPath = Routes.prevLessonLanguageLessonPath(language, lesson.slug);

  useHotkeys('ctrl+enter', handleRunCheck);

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">{t('help.controls.header')}</Popover.Header>
      <Popover.Body>{t('help.controls.body')}</Popover.Body>
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
        <OverlayTrigger trigger={['hover', 'focus']} placement="top" overlay={popover}>
          <Button
            variant="outline-secondary"
            className="me-3"
            onClick={handleReset}
            // title={t('resetCode')}
          >
            <span className="bi bi-arrow-repeat" />
          </Button>
        </OverlayTrigger>
        <a className={prevButtonClasses} href={prevLessonPath}>
          <span className="bi bi-arrow-left-short d-sm-none d-md-block d-lg-none" />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('prevLesson')}</span>
        </a>
        <Button
          variant="primary"
          className="me-3 d-inline-flex align-items-center"
          onClick={handleRunCheck}
          disabled={isCodeChecking}
        >
          {renderRunButtonContent()}
        </Button>
        <a className={nextButtonClasses} href={nextLessonPath}>
          <span className="bi bi-arrow-right-short d-sm-none d-md-block d-lg-none" />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('nextLesson')}</span>
        </a>
      </div>
    </div>
  );
}

export default ControlBox;
