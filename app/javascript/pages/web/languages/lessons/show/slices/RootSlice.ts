import * as Routes from "@/routes.js";
import type { LanguageLesson } from "@/types/serializers";
import {
  type PayloadAction,
  createAsyncThunk,
  createSlice,
} from "@reduxjs/toolkit";
import axios from "axios";
import type { CheckingResponse, RootState } from "../types";
import { initialRootState } from ".";

export const runCheck = createAsyncThunk(
  "runCheck",
  async (lesson: LanguageLesson, thunkAPI) => {
    const { content } = thunkAPI.getState() as RootState;
    const checkLessonPath = Routes.check_api_lesson_path(lesson.id!);
    const response = await axios.post(checkLessonPath, {
      version_id: lesson.version!,
      data: {
        attributes: {
          code: content,
        },
      },
    });

    const result = {
      ...response.data.attributes,
      output: atob(response.data.attributes.output),
    };
    return result;
  },
);

const slice = createSlice({
  name: "GeneralSlice",
  initialState: initialRootState,
  reducers: {
    changeContent(state, action: PayloadAction<RootState["content"]>) {
      state.content = action.payload;
    },
    changeTab(state, action: PayloadAction<RootState["currentTab"]>) {
      state.currentTab = action.payload;
      if (action.payload === "editor") {
        state.focusesCount += 1;
      }
    },
    setStartTime(state, { payload }) {
      state.startTime = payload.startTime;
    },
    changeSolutionState(
      state,
      action: PayloadAction<RootState["solutionState"]>,
    ) {
      state.solutionState = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(runCheck.pending, (state) => {
        state.currentTab = "output";
        state.processState = "checking";
      })
      .addCase(
        runCheck.fulfilled,
        (state, action: PayloadAction<CheckingResponse>) => {
          if (state.finished) {
            return;
          }

          if (action.payload.passed) {
            state.solutionState = "canBeShown";
          }
          state.finished = action.payload.passed;
          state.result = action.payload.result;
          state.output = action.payload.output;
          state.passed = action.payload.passed;
          state.processState = "checked";
        },
      )
      .addCase(runCheck.rejected, (state) => {
        state.passed = false;
        state.result = "error";
        state.processState = "checked";
      });
  },
});

export default slice;
