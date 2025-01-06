// @ts-check

import { createSlice } from "@reduxjs/toolkit";
import { currentTabValues } from "../utils/maps.js";
import { actions as checkInfoActions } from "./checkInfoSlice.js";

const slice = createSlice({
  name: "tabsBoxSlice",
  initialState: {
    currentTab: currentTabValues.editor,
  },
  reducers: {
    changeTab(state, { payload }) {
      state.currentTab = payload.newTabState;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(checkInfoActions.runCheck.pending, (state) => {
      state.currentTab = currentTabValues.output;
    });
  },
});

export const { actions } = slice;

export default slice.reducer;
