/* eslint-disable camelcase */
import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxSlice, { actions as tabsBoxActions } from './tabsBoxSlice.js';
import editorSlice, { actions as editorActions } from './editorSlice.js';
import solutionSlice, { actions as solutionActions } from './solutionSlice.js';
import lessonSlice, { actions as lessonActions } from './lessonSlice.js';
import checkInfoSlice, { actions as checkInfoActions } from './checkInfoSlice.js';
import { solutionStates, lessonMemberStates } from '../utils/stateMachines.js';

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

export const setupState = (gon) => (dispatch) => {
  const { lesson_version, lesson_member } = gon;
  dispatch(editorActions.changeContent({ content: lesson_version.prepared_code }));
  dispatch(solutionActions.setStartTime({ startTime: Date.now() }));
  const isFinished = lesson_member.state === lessonMemberStates.finished;
  if (isFinished) {
    dispatch(lessonActions.changeFinished({ finished: isFinished }));
    dispatch(solutionActions.changeSolutionProcessState({ processState: solutionStates.shown }));
  }
};
