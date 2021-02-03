import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import { deleteFromStorage } from '@rehooks/local-storage';

import cn from 'classnames';
import Hotkeys from 'react-hot-keys';
import { Button, Spinner } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faArrowRight, faSyncAlt, faArrowLeft, faPlayCircle,
} from '@fortawesome/free-solid-svg-icons';

import routes from '@/appRoutes.js';
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

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lessonVersion, editor }));
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
          <span className="sr-only">{t('loading')}</span>
          <span className="d-none d-sm-block d-md-none d-lg-block ml-1">{text}</span>
        </>
      );
    }

    return (
      <>
        <FontAwesomeIcon icon={faPlayCircle} />
        <span className="d-none d-sm-block d-md-none d-lg-block ml-1">{text}</span>
      </>
    );
  };

  const prevButtonClasses = cn(`btn btn-outline-secondary
    font-weight-normal mr-3 order-first order-sm-0 order-md-first order-lg-0`);

  const nextButtonClasses = cn('btn btn-outline-primary font-weight-normal', {
    disabled: !lessonInfo.finished,
  });

  const nextLessonPath = routes.nextLessonLanguageLessonPath(language, lesson.slug);
  const prevLessonPath = routes.prevLessonLanguageLessonPath(language, lesson.slug);

  return (
    <Hotkeys keyName="ctrl+Enter" onKeyUp={handleRunCheck}>
      <div className="d-flex justify-content-center p-3 border-top flex-shrink-0">
        <Button
          variant="secondary"
          className="mr-3 d-inline-flex align-items-center"
          onClick={handleReset}
          title={t('resetCode')}
        >
          <FontAwesomeIcon icon={faSyncAlt} />
        </Button>
        <a className={prevButtonClasses} href={prevLessonPath}>
          <FontAwesomeIcon className="d-sm-none d-md-block d-lg-none" icon={faArrowLeft} />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('prevLesson')}</span>
        </a>
        <Button
          variant="primary"
          className="mr-3 d-inline-flex align-items-center"
          onClick={handleRunCheck}
          disabled={isCodeChecking}
        >
          {renderRunButtonContent()}
        </Button>
        <a className={nextButtonClasses} href={nextLessonPath}>
          <FontAwesomeIcon className="d-sm-none d-md-block d-lg-none" icon={faArrowRight} />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('nextLesson')}</span>
        </a>
      </div>
    </Hotkeys>
  );
};

export default ControlBox;
