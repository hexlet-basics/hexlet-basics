import axios from "axios";
import {
  createAsyncThunk,
  createSlice,
  type PayloadAction,
} from "@reduxjs/toolkit";
import * as Routes from "@/routes.js";
import type { LanguageLesson } from "@/types/serializers";
import type { RootState } from ".";

interface CheckingResponse {
  result: boolean;
  output: string;
  passed: boolean;
}

export interface GeneralState {
  processState: "checked" | "unchecked" | "checking";
  currentTab: "editor" | "output" | "tests" | "solution";
  finished: boolean;
  result: unknown;
  output: string;
  passed: boolean;
  content: string;
  focusesCount: number;
  startTime: number;
  solutionState: "shown" | "canBeShown" | "notAllowedToBeShown";
  waitingTime: 0;
}

const initialState: GeneralState = {
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
  waitingTime: 0,
};

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
  initialState,
  reducers: {
    changeContent(state, { payload }) {
      state.content = payload.content;
    },
    changeTab(state, action: PayloadAction<GeneralState["currentTab"]>) {
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
      action: PayloadAction<GeneralState["solutionState"]>,
    ) {
      state.solutionState = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(runCheck.pending, (state) => {
        state.currentTab = "output";
        state.processState = "checking"
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
        state.processState = "checked"
      });
  },
});

export default slice;
