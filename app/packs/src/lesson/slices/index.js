// @ts-check

/* eslint-disable camelcase */
import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxSlice, { actions as tabsBoxActions } from './tabsBoxSlice.js';
import editorSlice, { actions as editorActions } from './editorSlice.js';
import solutionSlice, { actions as solutionActions } from './solutionSlice.js';
import lessonSlice, { actions as lessonActions } from './lessonSlice.js';
import checkInfoSlice, { actions as checkInfoActions } from './checkInfoSlice.js';
import { solutionStates, lessonMemberStates } from '../utils/maps.js';

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

// TODO move to preloadedState
export const setupState = (gon) => (dispatch) => {
  const { lesson_member } = gon;
  dispatch(solutionActions.setStartTime({ startTime: Date.now() }));
  const isFinished = lesson_member.state === lessonMemberStates.finished;
  if (isFinished) {
    dispatch(lessonActions.changeFinished({ finished: isFinished }));
    dispatch(solutionActions.changeSolutionProcessState({ processState: solutionStates.shown }));
  }
};
