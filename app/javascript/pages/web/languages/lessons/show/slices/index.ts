import { configureStore } from "@reduxjs/toolkit";
import slice from "./GeneralSlice.ts";
import { useDispatch, useSelector } from "react-redux";

const store = configureStore({
  reducer: slice.reducer,
});
export default store;

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
