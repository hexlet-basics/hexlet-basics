import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxReducers, { actions as tabsBoxActions, sliceName as tabsBoxSliceName } from './tabsBoxSlice.js';
import editorReducers, { actions as editorActions, sliceName as editorSliceName } from './editorSlice.js';
import soltionReducers, { actions as solutionActions, sliceName as solutionSliceName } from './solutionSlice.js';
import lessonReducers, { actions as lessonActions, sliceName as lessonSliceName } from './lessonSlice.js';
import checkInfoReducers, { actions as checkInfoActions, sliceName as checkInfoSliceName } from './checkInfoSlice.js';

export default combineReducers({
  [tabsBoxSliceName]: tabsBoxReducers,
  [editorSliceName]: editorReducers,
  [lessonSliceName]: lessonReducers,
  [solutionSliceName]: soltionReducers,
  [checkInfoSliceName]: checkInfoReducers,
});

export const actions = {
  ...tabsBoxActions,
  ...editorActions,
  ...lessonActions,
  ...solutionActions,
  ...checkInfoActions,
};

export {
  tabsBoxSliceName,
  editorSliceName,
  lessonSliceName,
  solutionSliceName,
  checkInfoSliceName,
};

export const setupState = (gon) => (dispatch) => {
  const { lesson } = gon;
  dispatch(editorActions.changeContent({ content: lesson.preparedCode }));
  dispatch(solutionActions.setStartTime({ startTime: Date.now() }));
};
