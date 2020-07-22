/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';

const slice = createSlice({
  name: 'editorSlice',
  initialState: {
    content: '',
  },
  reducers: {
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
