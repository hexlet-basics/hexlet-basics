import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';

import cn from 'classnames';
import Hotkeys from 'react-hot-keys';
import { Button, Spinner } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faArrowRight, faSyncAlt, faArrowLeft, faPlayCircle,
} from '@fortawesome/free-solid-svg-icons';

import {
  actions, lessonSliceName,
  checkInfoSliceName, editorSliceName,
} from '../slices/index.js';
import { checkInfoStates } from '../utils/stateMachines.js';
import EntityContext from '../EntityContext';

const ControlBox = () => {
  const { t } = useTranslation();
  const { checkInfo, lessonInfo, editor } = useSelector((state) => ({
    editor: state[editorSliceName],
    checkInfo: state[checkInfoSliceName],
    lessonInfo: state[lessonSliceName],
  }));

  const dispatch = useDispatch();
  const { lesson } = useContext(EntityContext);

  const handleRunCheck = () => {
    dispatch(actions.runCheck({ lesson, editor }));
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
    font-weight-normal mr-3 order-first order-sm-0 order-md-first order-lg-0`, {
      disabled: false,
    });

  const nextButtonClasses = cn('btn btn-outline-primary font-weight-normal', {
    disabled: !lessonInfo.finished,
  });

  const nextLessonPath = '#';
  const prevLessonPath = '#';

  return (
    <Hotkeys keyName="ctrl+Enter" onKeyUp={handleRunCheck}>
      <div className="mx-auto d-flex align-items-center text-center my-3">
        <a
          className="btn btn-outline-secondary mr-3 d-inline-flex align-items-center"
          href={window.location.href}
          title={t('resetCode')}
          // TODO: Add modal window instead browser confirmation
          data-confirm={t('confirm')}
        >
          <FontAwesomeIcon icon={faSyncAlt} />
          <span className="d-none d-sm-block d-md-none d-lg-block">&nbsp;</span>
        </a>
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
