/* eslint-disable no-param-reassign */
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { checkInfoStates } from '../utils/stateMachines.js';
import hexletAxios from '../../../lib/hexlet-axios.js';
import routes from '../../routes.js';

const runCheck = createAsyncThunk('runCheck', async ({ lesson, editor }) => {
  const checkLessonPath = routes.checkLessonPath(lesson.lesson_id);
  const response = await hexletAxios.post(checkLessonPath, {
    version_id: lesson.id,
    data: {
      attributes: {
        code: editor.content,
      },
    },
  });

  const result = { ...response.data.attributes, output: atob(response.data.attributes.output) };
  return result;
});

const slice = createSlice({
  name: 'checkInfoSlice',
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
