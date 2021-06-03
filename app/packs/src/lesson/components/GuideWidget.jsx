import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import ReactJoyride, { STATUS } from 'react-joyride';
import { actions } from '../slices/index.js';

const steps = [
  {
    disableBeacon: true,
    disableOverlayClose: true,
    title: 'Guide page',
    content: 'This is a guide.',
    locale: {
      skip: 'Skip guide',
    },
    placement: 'center',
    target: 'body',
  },
  {
    target: '[data-guide-id="Lesson"]',
    title: 'Урок',
    content: 'Текст для урока',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="Discuss"]',
    title: 'Обсуждение',
    content: 'Текст для обсуждения',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="editor"]',
    title: 'Редактор',
    content: 'Текст для редактора',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="testForExercise"]',
    title: 'Тесты',
    content: 'Текст для тестов',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="output"]',
    title: 'Вывод',
    content: 'Текст для вывода',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="solution"]',
    title: 'Решение',
    content: 'Текст для решения',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="resetCode"]',
    title: 'Сброс',
    content: 'Текст для сброса',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="prevButton"]',
    title: 'Предыдущий',
    content: 'Текст для предыдущего',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    target: '[data-guide-id="checkCode"]',
    title: 'Проверить',
    content: 'Текст для проверки',
    locale: {
      skip: 'Skip guide',
    },
  },
  {
    disableBeacon: true,
    disableOverlayClose: true,
    target: '[data-guide-id="nextButton"]',
    title: 'Следующий',
    content: 'Текст для следующего',
    locale: {
      skip: 'Skip guide',
    },
  },
];

const GuideWidget = () => {
  const dispatch = useDispatch();
  const [isFirstTime, setIsFirstTime] = useState(window.localStorage.getItem('guidePassed') === null);
  const isShowGuide = useSelector((state) => state.editorSlice.isShowGuide);

  return (
    (isShowGuide || isFirstTime) && (
    <ReactJoyride
      continuous
      run
      scrollToFirstStep
      showProgress
      showSkipButton
      steps={steps}
      spotlightPadding={6}
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
          window.localStorage.setItem('guidePassed', 'true');
          setIsFirstTime(false);
          dispatch(actions.showGuide({ isShowGuide: false }));
        }
      }}
    />
    )
  );
};

export default GuideWidget;
