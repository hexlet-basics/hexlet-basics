import { configureStore } from "@reduxjs/toolkit";
import React from "react";
import { Provider } from "react-redux";

import EntityContext from "./EntityContext.js";
import App from "./components/App.jsx";
import reducer from "./slices/index.js";
// import { lessonMemberStates, solutionStates } from "./utils/maps.js";

import type { LanguageLesson, LanguageLessonMember } from "@/types/serializers";
import type { PropsWithChildren } from "react";

const waitingTime = 20 * 60 * 1000; // 20 min

type Props = PropsWithChildren & {
  lesson: LanguageLesson;
  lessonMember: LanguageLessonMember;
};

export default function Editor({ lessonMember, lesson }: Props) {
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

  // const store = configureStore({
  //   preloadedState,
  //   reducer,
  // });

  return "JOPA";

  // return (
  //   <Provider store={store}>
  //     <EntityContext.Provider value={entities}>
  //       <App />
  //     </EntityContext.Provider>
  //   </Provider>
  // );
}
