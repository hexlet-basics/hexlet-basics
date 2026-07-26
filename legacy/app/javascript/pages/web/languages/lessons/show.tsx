import { usePage } from "@inertiajs/react";
import { CodeHighlightAdapterProvider } from "@mantine/code-highlight";
import shikiAdapter from "@/lib/shiki.ts";
import type { LessonSharedProps } from "@/types";
import Index from "./show/index.tsx";
import { LessonProvider } from "./show/store.tsx";

export default function Show() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;
  return (
    <CodeHighlightAdapterProvider adapter={shikiAdapter}>
      <LessonProvider lesson={lesson} lessonMember={lessonMember}>
        <Index />
      </LessonProvider>
    </CodeHighlightAdapterProvider>
  );
}
