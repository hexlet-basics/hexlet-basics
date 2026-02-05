import { usePage } from '@inertiajs/react';
import { Alert, Box, Code, ScrollArea } from '@mantine/core';
import { AnsiUp } from 'ansi_up';
import { Check, TriangleAlert } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import XssContent from '@/components/XssContent';
import * as Routes from '@/routes.js';
import { useLessonStore } from '../store.tsx';
import type { LessonSharedProps } from '../types.ts';

const ansi = new AnsiUp();

export default function OutputTab() {
  // const { lessonMember } = usePage<LessonSharedProps>().props;
  const result = useLessonStore((state) => state.result);
  const processState = useLessonStore((state) => state.processState);
  const output = useLessonStore((state) => state.output);
  const passed = useLessonStore((state) => state.passed);
  const { t: tCommon } = useTranslation('common');

  if (processState !== 'checked') {
    return null;
  }

  const message = tCommon(($) => $.check[result!].message);
  // const messageForGuest = tCommon("signInSuggestion", {
  //   url: Routes.new_user_path(),
  // });

  const outputAsHTML = ansi.ansi_to_html(decodeURIComponent(output));

  return (
    <>
      <Alert
        color={passed ? 'green' : 'yellow'}
        icon={passed ? <Check size={16} /> : <TriangleAlert size={16} />}
        withCloseButton={false}
        fz="sm"
        variant="light"
        radius="sm"
        my="xs"
      >
        <XssContent fz="sm">{message}</XssContent>
      </Alert>
      <Box component="pre" p="sm">
        <XssContent fz="sm">{outputAsHTML}</XssContent>
      </Box>
    </>
  );
}
