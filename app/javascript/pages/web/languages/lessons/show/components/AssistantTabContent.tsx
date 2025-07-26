import { usePage } from '@inertiajs/react';
import { Alert, Box } from '@mantine/core';
import i18next from 'i18next';
import { Info } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import Chat from '@/components/Chat.tsx';
import XssContent from '@/components/XssContent.tsx';
import { useLessonStore } from '@/pages/web/languages/lessons/show/store.tsx';
import type { LessonSharedProps } from '@/pages/web/languages/lessons/show/types.ts';

export default function AssistantTabContent({
  focusesCount,
}: {
  focusesCount: number;
}) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation('common');

  const {
    previousMessages,
    canCreateAssistantMessage,
    course,
    lesson,
    lessonMember,
  } = usePage<LessonSharedProps>().props;

  const userCode = useLessonStore((state) => state.content);
  const output = useLessonStore((state) => state.output);

  return (
    <Box p="lg">
      {i18next.language === 'ru' && (
        <Alert icon={<Info />} mb="lg">
          <XssContent>
            {t('languages.lessons.show.if_stuck_html', {
              url: tCommon('community_url'),
            })}
          </XssContent>
        </Alert>
      )}
      <Chat
        focusesCount={focusesCount}
        previousMessages={previousMessages}
        enabled={canCreateAssistantMessage}
        userCode={userCode}
        output={output}
        course={course}
        lesson={lesson}
        lessonMember={lessonMember}
      />
    </Box>
  );
}
