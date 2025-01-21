import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../types.ts";
import slice from "./RootSlice.ts";

export const initialRootState: RootState = {
  processState: "unchecked",
  currentTab: "editor",
  finished: false,
  result: null,
  output: "",
  passed: false,
  content: "",
  focusesCount: 1,
  startTime: 0,
  solutionState: "notAllowedToBeShown",
};

export default function getStore(incomingPreloadState?: Partial<AppState>) {
  const store = configureStore({
    reducer: slice.reducer,
    preloadedState: Object.assign({}, initialRootState, incomingPreloadState),
  });
  return store;
}
export type AppState = ReturnType<typeof slice.reducer>;
export type AppDispatch = ReturnType<typeof getStore>["dispatch"];

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<AppState>();
