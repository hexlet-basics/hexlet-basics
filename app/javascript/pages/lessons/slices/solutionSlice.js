// @ts-check

import { createSlice } from '@reduxjs/toolkit'
import { actions as checkInfoActions } from './checkInfoSlice.js'
import { solutionStates } from '../utils/maps.js'

const slice = createSlice({
  name: 'solutionSlice',
  initialState: {
    startTime: 0,
    processState: solutionStates.notAllowedToShown,
    waitingTime: 0,
  },
  reducers: {
    setStartTime(state, { payload }) {
      state.startTime = payload.startTime
    },
    changeSolutionProcessState(state, { payload }) {
      state.processState = payload.processState
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(checkInfoActions.runCheck.fulfilled, (state, { payload }) => {
        if (payload.passed) {
          state.processState = solutionStates.shown
        }
      })
  },
})

export const { actions } = slice

export default slice.reducer
