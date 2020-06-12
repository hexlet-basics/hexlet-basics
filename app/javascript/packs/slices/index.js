import { combineReducers } from '@reduxjs/toolkit';
import tabsReducers, { hooks as tabsHooks, sliceName as tabsSliceName } from './tabsSlice.js';

export default combineReducers({
  [tabsSliceName]: tabsReducers,
});

export {
  tabsHooks,
};
