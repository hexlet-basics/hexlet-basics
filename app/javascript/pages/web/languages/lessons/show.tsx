import { usePage } from "@inertiajs/react";
import { useMemo } from "react";
import { Provider } from "react-redux";
import Index from "./show/index.tsx";
import getStore, { type AppState } from "./show/slices/index.ts";
import type { LessonSharedProps } from "./show/types.ts";

// Replace rehype plugins with remark ones whenever possible
export default function Show() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;
  const isFinished = lessonMember?.state === "finished";

  const store = useMemo(() => {
    const preloadedState: Partial<AppState> = {
      startTime: Date.now(),
      defaultCode: lesson.prepared_code!,
      finished: isFinished,
      solutionState: isFinished ? "canBeShown" : "notAllowedToBeShown",
    };
    return getStore(preloadedState);
  }, [lesson.prepared_code, isFinished]);

  return (
    <Provider store={store}>
      <Index />
    </Provider>
  );
}
