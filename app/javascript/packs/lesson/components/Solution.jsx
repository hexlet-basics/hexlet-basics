import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import { Highlight } from 'react-fast-highlight';
import Countdown from 'react-countdown';
import { format } from 'date-fns';

import { actions } from '../slices/index.js';
import { getLanguage } from '../utils/editorUtils.js';
import { solutionStates } from '../utils/maps.js';
import EntityContext from '../EntityContext.js';

const Solution = () => {
  const { language, lessonVersion } = useContext(EntityContext);
  const { editor, solution } = useSelector((state) => ({
    editor: state.editorSlice,
    solution: state.solutionSlice,
  }));
  const { t } = useTranslation();
  const dispatch = useDispatch();

  const renderUserCode = () => {
    if (!editor.content) {
      return <p className="mt-3">{t('userCodeInstructions')}</p>;
    }

    return (
      <div>
        <p className="mt-3 mb-0">{t('userCode')}</p>
        <Highlight languages={[getLanguage(language)]}>
          {editor.content}
        </Highlight>
      </div>
    );
  };

  const renderSolution = () => (
    <div className="p-lg-3" id="basics-solution">
      <p className="mb-0">{t('teacherSolution')}</p>
      <Highlight languages={[getLanguage(language)]}>
        {lessonVersion.original_code}
      </Highlight>
      {renderUserCode()}
    </div>
  );

  const handleShowSolution = () => {
    dispatch(actions.changeSolutionProcessState({ processState: solutionStates.shown }));
  };

  const renderShowButton = () => (
    <>
      <p>{t('solutionNotice')}</p>
      <div className="text-center">
        <button
          type="button"
          className="btn btn-secondary px-4 mr-3"
          onClick={handleShowSolution}
        >
          {t('showSolution')}
        </button>
      </div>
    </>
  );

  const renderContent = (countdownData) => {
    const { completed } = countdownData;

    if (solutionStates.shown === solution.processState) {
      return renderSolution();
    }
    if (completed || solutionStates.canBeShown === solution.processState) {
      return renderShowButton();
    }

    const remainingTime = format(new Date(countdownData.total), 'mm:ss');

    return <p>{t('solutionInstructions', { remainingTime })}</p>;
  };

  return (
    <div>
      <Countdown date={solution.startTime + solution.waitingTime} renderer={renderContent} />
    </div>
  );
};

export default Solution;
