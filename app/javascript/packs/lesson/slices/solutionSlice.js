/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import { actions as checkInfoActions } from './checkInfoSlice.js';
import { solutionStates } from '../utils/stateMachines.js';

export const sliceName = 'solutionSlice';

const waitingTime = 20 * 60 * 1000;

const slice = createSlice({
  name: sliceName,
  initialState: {
    startTime: null,
    processState: solutionStates.notAllowedToShown,
    waitingTime,
  },
  reducers: {
    setStartTime(state, { payload }) {
      state.startTime = payload.startTime;
    },
    changeSolutionProcessState(state, { payload }) {
      state.processState = payload.processState;
    },
  },
  extraReducers: {
    [checkInfoActions.runCheck.fulfilled](state, { payload }) {
      if (payload.passed) {
        state.processState = solutionStates.shown;
      }
    },
  },
});

export const { actions } = slice;

export default slice.reducer;
