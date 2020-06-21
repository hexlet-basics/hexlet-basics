/* eslint-disable no-param-reassign */
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { checkInfoStates } from '../utils/stateMachines.js';

export const sliceName = 'checkInfoSlice';

const runCheck = createAsyncThunk('runCheck', async ({ lesson, editor }) => {
  const response = await new Promise((resolve) => {
    setTimeout(() => {
      resolve({ passed: true, result: 'passed' });
    }, 1000, lesson, editor);
  });
  return response;
});

const slice = createSlice({
  name: sliceName,
  initialState: {
    processState: checkInfoStates.unchecked,
    result: null,
    output: '',
    passed: false,
  },
  reducers: {
  },
  extraReducers: {
    [runCheck.pending](state) {
      state.processState = checkInfoStates.checking;
    },
    [runCheck.fulfilled](state, { payload }) {
      state.result = payload.result;
      state.output = payload.output;
      state.passed = payload.passed;
      state.processState = checkInfoStates.checked;
    },
    [runCheck.rejected](state) {
      state.passed = false;
      state.processState = checkInfoStates.checked;
    },
  },
});

export const actions = {
  ...slice.actions,
  runCheck,
};


export default slice.reducer;
