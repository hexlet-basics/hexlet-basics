import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxReducers, { actions as tabsBoxActions, sliceName as tabsBoxSliceName } from './tabsBoxSlice.js';

export default combineReducers({
  [tabsBoxSliceName]: tabsBoxReducers,
});

export const actions = {
  ...tabsBoxActions,
};

export {
  tabsBoxSliceName,
};
