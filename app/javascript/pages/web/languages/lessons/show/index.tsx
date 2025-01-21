import React, { StrictMode } from "react";
import { Provider } from "react-redux";

import App from "./components/App.tsx";

import { getKeyForStoringLessonCode } from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import useLocalStorage from "@rehooks/local-storage";
import getStore, { type AppState } from "./slices/index.ts";
import type { Props } from "./types";

export default function Index() {
  const { lessonMember, lesson } = usePage<Props>().props;

  const isFinished = lessonMember?.state === "finished";
  //
  const [content] = useLocalStorage<string>(
    getKeyForStoringLessonCode(lesson),
    lesson.prepared_code || "",
  );

  const preloadedState: Partial<AppState> = {
    startTime: Date.now(),
    content,
    finished: isFinished,
  };
  const store = getStore(preloadedState);

  return (
    <StrictMode>
      <Provider store={store}>
        <App />
      </Provider>
    </StrictMode>
  );
}
