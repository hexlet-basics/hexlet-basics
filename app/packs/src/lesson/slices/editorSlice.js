// @ts-check

/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';

import { actions as tabActions } from './tabsBoxSlice.js';

export const EDITORS = {
  monaco: 'MONACO',
  codemirror: 'CODEMIRROR',
};

const slice = createSlice({
  name: 'editorSlice',
  initialState: {
    content: '',
    focusesCount: 1,
    editorType: EDITORS.monaco,
  },
  reducers: {
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
    changeEditor(state, { payload }) {
      state.editorType = payload.editor;
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
