import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxReducers, { actions as tabsBoxActions, sliceName as tabsBoxSliceName } from './tabsBoxSlice.js';
import editorReducers, { actions as editorActions, sliceName as editorSliceName } from './editorSlice.js';

export default combineReducers({
  [tabsBoxSliceName]: tabsBoxReducers,
  [editorSliceName]: editorReducers,
});

export const actions = {
  ...tabsBoxActions,
  ...editorActions,
};

export {
  tabsBoxSliceName,
  editorSliceName,
};

export const setupState = (gon) => (dispatch) => {
  const { language } = gon;
  dispatch(editorActions.changeLanguage({ language }));
};
