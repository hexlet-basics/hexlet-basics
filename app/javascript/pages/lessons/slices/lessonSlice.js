// @ts-check

import { createSlice } from "@reduxjs/toolkit";
import { actions as checkInfoActions } from "./checkInfoSlice.js";

const slice = createSlice({
  name: "lessonSlice",
  initialState: {
    finished: false,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      checkInfoActions.runCheck.fulfilled,
      (state, { payload }) => {
        if (state.finished) {
          return;
        }
        state.finished = payload.passed;
      },
    );
  },
});

export const { actions } = slice;

export default slice.reducer;
