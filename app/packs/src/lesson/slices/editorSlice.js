// @ts-check

/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';

import { actions as tabActions } from './tabsBoxSlice.js';

const slice = createSlice({
  name: 'editorSlice',
  initialState: {
    content: '',
    focusesCount: 1,
    isShowGuide: false,
  },
  reducers: {
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
    showGuide: (state, { payload }) => {
      Object.assign(state, payload);
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(tabActions.changeTab, (state, { payload }) => {
        const { newTabState } = payload;
        if (newTabState === 'editor') {
          state.focusesCount += 1;
        }
      });
  },
});

export const { actions } = slice;

export default slice.reducer;
