import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';

import cn from 'classnames';
import Hotkeys from 'react-hot-keys';
import { Button, Spinner } from 'react-bootstrap';

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

  const isCheckInfoChecking = checkInfo.processState === checkInfoStates.checking;

  const renderRunButtonContent = () => {
    const text = t('run');
    if (isCheckInfoChecking) {
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
        <span className="fas fa-play-circle" />
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
          <span className="fas fa-sync-alt" />
          <span className="d-none d-sm-block d-md-none d-lg-block">&nbsp;</span>
        </a>
        <a className={prevButtonClasses} href={prevLessonPath}>
          <span className="d-sm-none d-md-block d-lg-none fas fa-arrow-left" />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('prevLesson')}</span>
        </a>
        <Button
          variant="primary"
          className="mr-3 d-inline-flex align-items-center"
          onClick={handleRunCheck}
          disabled={isCheckInfoChecking}
        >
          {renderRunButtonContent()}
        </Button>
        <a className={nextButtonClasses} href={nextLessonPath}>
          <span className="d-sm-none d-md-block d-lg-none fas fa-arrow-right" />
          <span className="d-none d-sm-block d-md-none d-lg-block">{t('nextLesson')}</span>
        </a>
      </div>
    </Hotkeys>
  );
};

export default ControlBox;
