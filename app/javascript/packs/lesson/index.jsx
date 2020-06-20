import React from 'react';
import ReactDOM from 'react-dom';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import gon from 'gon';

import App from './components/App.jsx';
import reducer, { setupState } from './slices/index.js';
import resources from '../locales/index.js';


export default async () => {
  await i18n
    .use(initReactI18next)
    .init({
      resources,
      load: 'languageOnly',
      fallbackLng: false,
      lng: 'ru',
      debug: process.env.NODE_ENV !== 'production',
      react: {
        wait: true,
      },
    });

  const store = configureStore({
    reducer,
  });
  store.dispatch(setupState(gon));

  ReactDOM.render(
    <Provider store={store}>
      <App />
    </Provider>,
    document.querySelector('#basics-lesson-container'),
  );
};
