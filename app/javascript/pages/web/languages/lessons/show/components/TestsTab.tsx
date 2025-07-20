import { usePage } from '@inertiajs/react';
import { useTranslation } from 'react-i18next';
import { LazyCodeHighlight } from '@/components/LazyCodeHighlight.tsx';
import { getEditorLanguage } from '@/lib/utils.ts';
import type { LessonSharedProps } from '../types.ts';

export default function TestsTab() {
  const { t: tCommon } = useTranslation('common');
  const { lesson, course } = usePage<LessonSharedProps>().props;

  return (
    <>
      <p className="text-center lead">{tCommon('testInstructions')}</p>
      <LazyCodeHighlight
        language={getEditorLanguage(course.slug!)}
        code={lesson.test_code ?? ''}
      />
    </>
  );
}
