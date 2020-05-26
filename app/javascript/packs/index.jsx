import React from 'react';
import ReactDOM from 'react-dom';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';

import { rootReducer } from './slices';

export default () => {
  const store = configureStore({
    reducer: rootReducer,
  });

  ReactDOM.render(
    <Provider store={store}>
      <div>Hello</div>
    </Provider>,
    document.getElementById('main'),
  );
};
