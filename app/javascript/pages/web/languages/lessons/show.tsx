import { Head, usePage } from '@inertiajs/react';
import Index from './show/index.tsx';
import { LessonProvider } from './show/store.tsx';
import type { LessonSharedProps } from './show/types.ts';

export default function Show() {
  const { lessonMember, lesson } = usePage<LessonSharedProps>().props;
  return (
    <>
      <LessonProvider lesson={lesson} lessonMember={lessonMember}>
        <Index />
      </LessonProvider>
    </>
  );
}
