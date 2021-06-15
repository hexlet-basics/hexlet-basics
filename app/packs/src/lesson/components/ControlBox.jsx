// @ts-check

import React, { useContext, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import { useLocalStorage, deleteFromStorage } from '@rehooks/local-storage';

import cn from 'classnames';
// import Hotkeys from 'react-hot-keys';
import { useHotkeys } from 'react-hotkeys-hook';
import { Button, Spinner } from 'react-bootstrap';

// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
// import {
//   faArrowRight, faSyncAlt, faArrowLeft, faPlayCircle,
// } from '@fortawesome/free-solid-svg-icons';

import routes from 'vendor/appRoutes.js';
import { actions } from '../slices/index.js';
import { checkInfoStates } from '../utils/maps.js';
import EntityContext from '../EntityContext.js';

const ControlBox = () => {
  const { t } = useTranslation();
  const { checkInfo, lessonInfo, editor } = useSelector((state) => ({
    editor: state.editorSlice,
    checkInfo: state.checkInfoSlice,
    lessonInfo: state.lessonSlice,
  }));

  const dispatch = useDispatch();
  const {
    lessonVersion, language, lesson,
  } = useContext(EntityContext);

  const localStorageKey = `lesson-version-${lessonVersion.id}`;
  const [localStorageContent, setContent] = useLocalStorage(localStorageKey);

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lessonVersion, editor }));
  };

  const handleReset = () => {
    // NOTE: easier than state manipulating. dont touch, dont blame, be happy.
    // eslint-disable-next-line no-alert
    if (window.confirm(t('confirm'))) {
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

    // <FontAwesomeIcon icon={faPlayCircle} />
    return (
      <>
        <span className="bi bi-play-circle" />
        <span className="d-none d-sm-block d-md-none d-lg-block ms-1">{text}</span>
      </>
    );
  };

  const prevButtonClasses = cn(`btn btn-outline-secondary
    font-weight-normal me-3 order-first order-sm-0 order-md-first order-lg-0`);

  const nextButtonClasses = cn('btn btn-outline-primary font-weight-normal', {
    disabled: !lessonInfo.finished,
  });

  const nextLessonPath = routes.nextLessonLanguageLessonPath(language, lesson.slug);
  const prevLessonPath = routes.prevLessonLanguageLessonPath(language, lesson.slug);

  useHotkeys('ctrl+enter', handleRunCheck);

  const isFirstSoluton = !localStorageContent && editor.content;
  const isCheckDone = checkInfo.processState === 'checked';

  useEffect(() => {
    if (isCheckDone && checkInfo.passed || isFirstSoluton) {
      setContent(editor.content);
    }
  }, [checkInfo.processState]);

  // <FontAwesomeIcon icon={faSyncAlt} />
  return (
    <div className="d-flex justify-content-center p-3 border-top flex-shrink-0">
      <Button
        variant="secondary"
        className="me-3 d-inline-flex align-items-center"
        onClick={handleReset}
        title={t('resetCode')}
      >
        <span className="bi bi-arrow-repeat" />
      </Button>
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
  );
};

export default ControlBox;
