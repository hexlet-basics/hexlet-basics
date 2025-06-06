import { neededPreview } from "@/lib/utils.ts";
import { usePage } from "@inertiajs/react";
import { Suspense } from "react";
import type { LessonSharedProps } from "../types.ts";
import ControlBox from "./ControlBox.tsx";
import HTMLPreview from "./HTMLPreview.tsx";
import TabsBox from "./TabsBox.tsx";
import { useLessonStore } from "../store.tsx";

function App() {
  const { course } = usePage<LessonSharedProps>().props;

  const content = useLessonStore((state) => state.content);
  const currentTab = useLessonStore((state) => state.currentTab);

  const renderHtmlPreview = () => {
    if (currentTab !== "editor") {
      return null;
    }
    if (!neededPreview(course.slug!)) {
      return null;
    }

    return (
      <Suspense fallback={<div>Loading...</div>}>
        <HTMLPreview html={content} />
      </Suspense>
    );
  };

  return (
    <>
      <TabsBox />
      {renderHtmlPreview()}
      <ControlBox />
    </>
  );
}

export default App;
