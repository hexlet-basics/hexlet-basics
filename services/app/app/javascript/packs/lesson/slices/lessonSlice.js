/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import { actions as checkInfoActions } from './checkInfoSlice.js';

const slice = createSlice({
  name: 'lessonSlice',
  initialState: {
    finished: false,
  },
  reducers: {
    changeFinished(state, { payload }) {
      state.finished = payload.finished;
    },
  },
  extraReducers: {
    [checkInfoActions.runCheck.fulfilled](state, { payload }) {
      if (state.finished) {
        return;
      }
      state.finished = payload.passed;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
