// @ts-check

import { combineReducers } from "@reduxjs/toolkit";
import checkInfoSlice, {
  actions as checkInfoActions,
} from "./checkInfoSlice.js";
import editorSlice, { actions as editorActions } from "./editorSlice.js";
import lessonSlice, { actions as lessonActions } from "./lessonSlice.js";
import solutionSlice, { actions as solutionActions } from "./solutionSlice.js";
import tabsBoxSlice, { actions as tabsBoxActions } from "./tabsBoxSlice.js";

export default combineReducers({
  tabsBoxSlice,
  editorSlice,
  solutionSlice,
  lessonSlice,
  checkInfoSlice,
});

export const actions = {
  ...tabsBoxActions,
  ...editorActions,
  ...lessonActions,
  ...solutionActions,
  ...checkInfoActions,
};
