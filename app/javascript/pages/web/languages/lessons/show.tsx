import { usePage } from "@inertiajs/react";
import Index from "./show/index.tsx";
import type { LessonSharedProps } from "./show/types.ts";
import { LessonProvider } from "./show/store.tsx";

export default function Show() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;
  return (
    <LessonProvider
      lesson={lesson}
      lessonMember={lessonMember}
    >
      <Index />
    </LessonProvider>
  );
}
