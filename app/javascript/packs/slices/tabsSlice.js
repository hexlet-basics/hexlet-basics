/* eslint-disable no-param-reassign */
import { createSlice } from '@reduxjs/toolkit';
import currentTabStates from 'packs/utils/currentTabStates.js';
import { useDispatch, useSelector } from 'react-redux';

export const sliceName = 'tabs';

const slice = createSlice({
  name: sliceName,
  initialState: {
    currentTab: currentTabStates.editor,
  },
  reducers: {
    changeTab(state, { payload }) {
      state.currentTab = payload;
    },
  },
});

const { actions } = slice;

const useAction = () => {
  const dispatch = useDispatch();

  return {
    changeTab: (newTabState) => {
      dispatch(actions.changeTab(newTabState));
    },
  };
};

const showSelectors = () => {
  const tabsSlice = useSelector((state) => state[sliceName]);

  return {
    getData: () => tabsSlice,
  };
};

export const hooks = {
  useAction,
  showSelectors,
};

export default slice.reducer;
