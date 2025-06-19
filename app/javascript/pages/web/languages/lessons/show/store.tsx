import { createContext, useContext, useMemo } from 'react';
import { createStore, useStore } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import type { } from '@redux-devtools/extension';
import analytics from "@/lib/analytics.ts";
import * as Routes from "@/routes.js";
import axios from "axios";
import i18next from "i18next";
import type { LanguageLesson, LessonCheckingResponse } from "@/types";
import { LessonProps, LessonState } from './types.ts';

type LessonStore = ReturnType<typeof createLessonStore>
export const LessonContext = createContext<LessonStore | null>(null)

export function createLessonStore(lesson: LanguageLesson, initProps?: Partial<LessonState>) {
  return createStore<LessonState>()(
    devtools(
      persist((set, get) => ({
        processState: "unchecked",
        currentTab: "editor",
        finished: false,
        defaultCode: "",
        result: null,
        resetsCount: 0,
        output: "",
        passed: false,
        content: "",
        focusesCount: 1,
        startTime: 0,
        solutionState: "notAllowedToBeShown",
        ...initProps,
        changeContent: (content) => set({ content }),
        resetContent: () =>
          set((state) => ({
            content: state.defaultCode,
            focusesCount: state.focusesCount + 1,
            resetsCount: state.resetsCount + 1,
            currentTab: "editor",
          })),
        changeTab: (currentTab) =>
          set((state) => ({
            currentTab,
            focusesCount: currentTab === "editor" ? state.focusesCount + 1 : state.focusesCount,
          })),
        setStartTime: (startTime) => set({ startTime }),
        changeSolutionState: (solutionState) => set({ solutionState }),
        runCheck: async ({ course, lesson }) => {
          set({ currentTab: "output", processState: "checking" });
          const { content } = get();
          const checkLessonPath = Routes.check_api_lesson_path(lesson.id!);
          const responsePromise = axios.post<LessonCheckingResponse>(checkLessonPath, {
            version_id: lesson.version!,
            data: { attributes: { code: content } },
          });
          let response: Awaited<typeof responsePromise> | null = null
          try {
            response = await responsePromise
          } catch (error) {
            if (axios.isAxiosError(error)) {
              set((state) => ({
                processState: "unchecked",
              }));
              return false
            } else {
              throw error
            }
          }
          const response_data = response.data;
          const {
            lesson_has_been_finished: lessonHasBeenFinished,
            language_has_been_finished: courseHasBeenFinished,
          } = response_data;

          analytics.track("solution_checked", {
            course_slug: course.slug,
            lesson_slug: lesson.slug,
            passed: response_data.passed,
            locale: i18next.language,
          });

          if (lessonHasBeenFinished) {
            analytics.track("lesson_finished", {
              course_slug: course.slug,
              lesson_slug: lesson.slug,
              locale: i18next.language,
            });
          }

          if (courseHasBeenFinished) {
            analytics.track("course_finished", {
              slug: course.slug,
              locale: i18next.language,
            });
          }

          const result: LessonCheckingResponse & { output: string } = {
            ...response_data,
            output: atob(response_data.output),
          };
          set((state) => ({
            solutionState: result.passed ? "shown" : state.solutionState,
            finished: state.finished || result.passed,
            result: result.result,
            output: result.output,
            passed: result.passed,
            processState: "checked",
          }));

          return true
        },
      }),
        {
          partialize: (state) => ({ content: state.content }),
          name: `lesson-${lesson.id}`,
        }
      )
    )
  )
}

type LessonProviderProps = React.PropsWithChildren<LessonProps>

export function LessonProvider({ children, ...props }: LessonProviderProps) {
  const { lesson, lessonMember } = props
  const store = useMemo(() => {
    const isFinished = lessonMember?.state === "finished";

    const startTime = Date.now()
    const defaultCode = lesson.prepared_code!
    const finished = isFinished
    const solutionState = isFinished ? "canBeShown" : "notAllowedToBeShown"

    return createLessonStore(lesson, {
      startTime, defaultCode, finished, solutionState, content: defaultCode,
    });
    // oxlint-disable-next-line exhaustive-deps
  }, []);

  return (
    <LessonContext.Provider value={store}>
      {children}
    </LessonContext.Provider>
  )
}

export function useLessonStore<T>(selector: (state: LessonState) => T): T {
  const store = useContext(LessonContext)
  if (!store) throw new Error('Missing LessonContext.Provider in the tree')
  return useStore(store, selector)
}
