import React, { StrictMode } from "react";
import { Provider } from "react-redux";

import App from "./components/App.tsx";

import { usePage } from "@inertiajs/react";
import getStore, { type AppState } from "./slices/index.ts";
import type { LessonSharedProps } from "./types.ts";

export default function Index() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;

  const isFinished = lessonMember?.state === "finished";

  const preloadedState: Partial<AppState> = {
    startTime: Date.now(),
    defaultCode: lesson.prepared_code!,
    finished: isFinished,
    solutionState: isFinished ? "canBeShown" : "notAllowedToBeShown",
  };
  const store = getStore(preloadedState);

  return (

    <Provider store={store}>
      <App />
    </Provider>

  );
}
