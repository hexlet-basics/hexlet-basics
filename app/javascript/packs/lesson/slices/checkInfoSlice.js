/* eslint-disable no-param-reassign */
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { checkInfoStates } from '../utils/stateMachines.js';
import hexletAxios from '../../hexlet-axios';
import routes from '../../routes';

export const sliceName = 'checkInfoSlice';

const runCheck = createAsyncThunk('runCheck', async ({ lesson, editor }) => {
  const response = await hexletAxios.post(
    routes.checkLessonVersionPath(lesson.language_id, lesson.lesson_id, lesson.id),
    {
      data: {
        attributes: {
          code: editor.content,
        },
      },
    },
  );

  const result = { ...response.data.attributes, output: atob(response.data.attributes.output) };
  return result;
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
