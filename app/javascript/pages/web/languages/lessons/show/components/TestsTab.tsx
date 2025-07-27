import { usePage } from '@inertiajs/react';
import { CodeHighlight } from '@mantine/code-highlight';
import { Center, Stack, Text } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { getEditorLanguage } from '@/lib/utils.ts';
import type { LessonSharedProps } from '../types';

export default function TestsTab() {
  const { t: tCommon } = useTranslation('common');
  const { lesson, course } = usePage<LessonSharedProps>().props;

  return (
    <Stack>
      <Center>
        <Text>{tCommon('testInstructions')}</Text>
      </Center>
      <CodeHighlight
        language={getEditorLanguage(course.slug!)}
        code={lesson.test_code ?? ''}
      />
    </Stack>
  );
}
