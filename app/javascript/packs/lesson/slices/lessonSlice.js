/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import { actions as checkInfoActions } from './checkInfoSlice.js';

export const sliceName = 'lessonSlice';

const slice = createSlice({
  name: sliceName,
  initialState: {
    finished: false,
  },
  reducers: {
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
