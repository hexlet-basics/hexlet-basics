import { usePage } from '@inertiajs/react';
import { Suspense } from 'react';
import { neededPreview } from '@/lib/utils.ts';
import HTMLPreview from '@/pages/web/languages/lessons/show/components/HTMLPreview.tsx';
import { useLessonStore } from '@/pages/web/languages/lessons/show/store.tsx';
import type { LessonSharedProps } from '@/pages/web/languages/lessons/show/types.ts';

export default function HTMLPreviewBlock() {
  const { course } = usePage<LessonSharedProps>().props;
  const currentTab = useLessonStore((state) => state.currentTab);
  const content = useLessonStore((state) => state.content);

  if (currentTab !== 'editor') {
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
}
