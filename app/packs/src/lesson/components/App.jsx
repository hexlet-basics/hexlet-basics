// @ts-check

import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import ReactJoyride, { STATUS } from 'react-joyride';
import { useTranslation } from 'react-i18next';

import TabsBox from './TabsBox.jsx';
import ControlBox from './ControlBox.jsx';
import HTMLPreview from './HTMLPreview.jsx';
import { currentTabValues } from '../utils/maps.js';
import { neededPreview } from '../utils/languagesUtils.js';
import EntityContext from '../EntityContext.js';

const Guidewidget = ({ runGuide, setRunGuide }) => {
  const { t } = useTranslation();

  const [isFirstTime, setIsFirstTime] = useState(window.localStorage.getItem('guidePassed') === null);

  const steps = [
    {
      disableBeacon: true,
      disableOverlayClose: true,
      title: t('.guidePage'),
      content: 'This is a guide.',
      locale: {
        skip: 'Skip guide',
      },
      placement: 'center',
      target: 'body',
    },
    {
      disableBeacon: true,
      disableOverlayClose: true,
      target: '[data-guide-id="Lesson"]',
      title: 'Урок',
      content: 'Текст для урока',
      locale: {
        skip: 'Skip guide',
      },
    },
    {
      disableBeacon: true,
      disableOverlayClose: true,
      target: '[data-guide-id="Discuss"]',
      title: 'Обсуждение',
      content: 'Текст для обсуждения',
      locale: {
        skip: 'Skip guide',
      },
    },
  ];

  return (
    (runGuide || isFirstTime) && (
    <ReactJoyride
      continuous
      scrollToFirstStep
      showProgress
      showSkipButton
      runGuide={runGuide}
      steps={steps}
      styles={{
        options: {
          primaryColor: '#0275d8',
          zIndex: 1000,
        },
        buttonNext: {
          borderRadius: 'unset',
        },
      }}
      callback={({ status }) => {
        if ([STATUS.FINISHED, STATUS.SKIPPED].includes(status)) {
          console.log(status);
          window.localStorage.setItem('guidePassed', 'true');
          setIsFirstTime(false);
          setRunGuide(false);
        }
      }}
    />
    )
  );
};

const App = () => {
  const { currentTab, content } = useSelector((state) => ({ ...state.tabsBoxSlice, ...state.editorSlice }));
  const [runGuide, setRunGuide] = useState(false);
  const { language } = React.useContext(EntityContext);

  const handleClickStart = (e) => {
    e.preventDefault();
    setRunGuide(true);
  };

  const renderHtmlPreview = () => {
    if (currentTab !== currentTabValues.editor) {
      return null;
    }
    if (!neededPreview(language)) {
      return null;
    }

    return <HTMLPreview html={content} />;
  };

  return (
    <div className="card vh-100 x-h-md-100">
      <TabsBox handleClickStart={handleClickStart} />
      {renderHtmlPreview()}
      <ControlBox />
      <Guidewidget runGuide={runGuide} setRunGuide={setRunGuide} />
    </div>
  );
};

export default App;
