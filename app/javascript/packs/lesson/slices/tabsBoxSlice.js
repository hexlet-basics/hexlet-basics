/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import { actions as checkInfoActions } from './checkInfoSlice.js';
import { currentTabValues } from '../utils/stateMachines.js';

export const sliceName = 'tabsBoxSlice';

const slice = createSlice({
  name: sliceName,
  initialState: {
    currentTab: currentTabValues.editor,
  },
  reducers: {
    changeTab(state, { payload }) {
      state.currentTab = payload.newTabState;
    },
  },
  extraReducers: {
    [checkInfoActions.runCheck.pending](state) {
      state.currentTab = currentTabValues.console;
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
