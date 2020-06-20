/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';

export const sliceName = 'editorSlice';

const slice = createSlice({
  name: sliceName,
  initialState: {
    language: '',
    content: '',
  },
  reducers: {
    changeLanguage(state, { payload }) {
      state.language = payload.language;
    },
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
