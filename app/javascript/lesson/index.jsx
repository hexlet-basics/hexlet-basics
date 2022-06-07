// @ts-check

import React from 'react';
import { createRoot } from 'react-dom/client';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
// eslint-disable-next-line import/no-unresolved
import gon from 'gon';

import App from './components/App.jsx';
import reducer from './slices/index.js';
import resources from '../../assets/locales/index.js';
import EntityContext from './EntityContext.js';
import { lessonMemberStates, solutionStates } from './utils/maps.js';

const waitingTime = 20 * 60 * 1000; // 20 min

export default async () => {
  await i18n
    .use(initReactI18next)
    .init({
      resources,
      load: 'languageOnly',
      fallbackLng: false,
      lng: gon.locale,
      debug: process.env.NODE_ENV !== 'production',
      react: {
        useSuspense: true,
      },
    });

  // eslint-disable-next-line react/jsx-no-constructed-context-values
  const entities = {
    lessonVersion: gon.lesson_version,
    lesson: gon.lesson,
    language: gon.language,
    lessonMember: gon.lesson_member,
  };

  const isFinished = gon.lesson_member.state === lessonMemberStates.finished;

  const localStorageKey = `lesson-version-${gon.lesson_version.id}`;
  const locallySavedContent = localStorage.getItem(localStorageKey);

  const preloadedState = {
    editorSlice: {
      content: locallySavedContent || gon.lesson_version.prepared_code || '',
      focusesCount: 1,
    },
    solutionSlice: {
      // TODO move counter to server
      startTime: Date.now(),
      processState: isFinished
        ? solutionStates.shown
        : solutionStates.notAllowedToShown,
      waitingTime,
    },
    lessonSlice: {
      finished: isFinished,
    },
  };

  const store = configureStore({
    preloadedState,
    reducer,
  });

  const container = document.querySelector('#basics-lesson-container');

  const root = createRoot(container);

  root.render(
    <Provider store={store}>
      <EntityContext.Provider value={entities}>
        <App />
      </EntityContext.Provider>
    </Provider>,
  );
};
