import { usePage } from '@inertiajs/react';
import { CodeHighlightAdapterProvider } from '@mantine/code-highlight';
import shikiAdapter from '@/lib/shiki.ts';
import Index from './show/index.tsx';
import { LessonProvider } from './show/store.tsx';
import type { LessonSharedProps } from './show/types.ts';

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
