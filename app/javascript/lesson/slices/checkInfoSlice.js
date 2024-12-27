// @ts-check

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import hexletAxios from '../../lib/hexlet-axios.js'
import { checkInfoStates } from '../utils/maps.js'
import Routes from '../utils/configured_routes.js'

const runCheck = createAsyncThunk('runCheck', async ({ lessonVersion }, { getState }) => {
  const { editorSlice: { content } } = getState()
  const checkLessonPath = Routes.checkApiLessonPath(lessonVersion.lesson_id)
  const response = await hexletAxios.post(checkLessonPath, {
    version_id: lessonVersion.id,
    data: {
      attributes: {
        code: content,
      },
    },
  })

  const result = { ...response.data.attributes, output: atob(response.data.attributes.output) }
  return result
})

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
  extraReducers: (builder) => {
    builder
      .addCase(runCheck.pending, (state) => {
        state.processState = checkInfoStates.checking
      })
      .addCase(runCheck.fulfilled, (state, { payload }) => {
        state.result = payload.result
        state.output = payload.output
        state.passed = payload.passed
        state.processState = checkInfoStates.checked
      })
      .addCase(runCheck.rejected, (state) => {
        state.passed = false
        state.result = 'error'
        state.processState = checkInfoStates.checked
      })
  },
})

export const actions = {
  ...slice.actions,
  runCheck,
}

export default slice.reducer
