import { combineReducers } from '@reduxjs/toolkit';
import tabsBoxReducers, { actions as tabsBoxActons, sliceName as tabsBoxSliceName } from './tabsBoxSlice.js';

export default combineReducers({
  [tabsBoxSliceName]: tabsBoxReducers,
});

export const actions = {
  ...tabsBoxActons,
};

export {
  tabsBoxSliceName,
};
