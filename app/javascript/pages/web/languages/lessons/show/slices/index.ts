import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector } from "react-redux";
import slice, { initialRootState } from "./RootSlice.ts";

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
