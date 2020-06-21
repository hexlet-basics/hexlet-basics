/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import { solutionStates } from '../utils/stateMachines.js';

export const sliceName = 'solutionSlice';

const waitingTime = 20 * 60 * 1000;

const slice = createSlice({
  name: sliceName,
  initialState: {
    startTime: null,
    processState: solutionStates.pending,
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
});

export const { actions } = slice;

export default slice.reducer;
