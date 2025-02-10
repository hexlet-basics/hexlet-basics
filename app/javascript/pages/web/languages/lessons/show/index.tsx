import { Provider } from "react-redux";

import App from "./components/App.tsx";

import { usePage } from "@inertiajs/react";
import getStore, { type AppState } from "./slices/index.ts";
import type { LessonSharedProps } from "./types.ts";
import { getKeyForStoringLessonCode } from "@/lib/utils.ts";

export default function Index() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;

  const isFinished = lessonMember?.state === "finished";

  const localStorageKey = getKeyForStoringLessonCode(lesson);
  const localStorageValue = localStorage.getItem(localStorageKey);
  const defaultCode = lesson.prepared_code || "";
  const code = localStorageValue ? JSON.parse(localStorageValue) : defaultCode;

  const preloadedState: Partial<AppState> = {
    startTime: Date.now(),
    defaultCode,
    content: code,
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
