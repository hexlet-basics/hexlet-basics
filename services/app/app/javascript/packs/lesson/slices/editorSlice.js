/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';

const slice = createSlice({
  name: 'editorSlice',
  initialState: {
    content: '',
    cursorPosition: null,
  },
  reducers: {
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
    saveCursorPosition(state, { payload }) {
      state.cursorPosition = payload.cursorPosition;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
