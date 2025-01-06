import { configureStore } from "@reduxjs/toolkit";
import React from "react";
import { Provider } from "react-redux";

import App from "./components/App.tsx";
import reducer from "./slices/index.ts";

import type { LanguageLesson, LanguageLessonMember } from "@/types/serializers";
import type { PropsWithChildren } from "react";
import { usePage } from "@inertiajs/react";
import type { Props } from "./types";

const waitingTime = 20 * 60 * 1000; // 20 min

// type Props = PropsWithChildren & {
//   lesson: LanguageLesson;
//   lessonMember: LanguageLessonMember;
// };

export default function Index() {
  const {
    // courseCategory,
    // course,
    // lessons,
    lessonMember,
    lesson,
    // prevLesson,
    // nextLesson,
    // auth: { user },
  } = usePage<Props>().props;

  const entities = {
    lessonVersion: lesson.version,
    // lesson: gon.lesson,
    // language: gon.language,
    // lessonMember: gon.lesson_member,
  };
  //
  const isFinished = lessonMember && lessonMember.state === "finished";
  //
  const localStorageKey = `lesson-version-${lesson.id}`;
  const locallySavedContent = localStorage.getItem(localStorageKey);
  //
  const preloadedState = {
    editorSlice: {
      content: locallySavedContent || lesson.prepared_code || "",
      focusesCount: 1,
    },
    solutionSlice: {
      // TODO move counter to server
      startTime: Date.now(),
      processState: isFinished ? "shown" : "notAllowedToBeShown",
      waitingTime,
    },
    lessonSlice: {
      finished: isFinished,
    },
  };

  const store = configureStore({
    preloadedState,
    reducer,
  });

  return (
    <Provider store={store}>
      <App />
    </Provider>
  );
}
