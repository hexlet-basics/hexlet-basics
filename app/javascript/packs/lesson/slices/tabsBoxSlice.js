/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import currentTabStates from '../utils/currentTabStates.js';

export const sliceName = 'tabsBoxSlice';

const slice = createSlice({
  name: sliceName,
  initialState: {
    currentTab: currentTabStates.editor,
  },
  reducers: {
    changeTab(state, { payload }) {
      state.currentTab = payload.newTabState;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
