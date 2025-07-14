import { Button, Group, Loader, Stack, Textarea } from '@mantine/core';
import i18next from 'i18next';
import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  type AssistantMessage,
  useAssistantStream,
} from '@/hooks/useAssistantStream';
import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
} from '@/types/serializers';
// import { useInView } from "react-intersection-observer";
import MarkdownViewer from './MarkdownViewer';

type Props = {
  enabled: boolean;
  lesson: LanguageLesson;
  focusesCount: number;
  course: Language;
  userCode: string;
  output: string;
  previousMessages: AssistantMessage[];
  lessonMember?: LanguageLessonMember;
};

function MessagePresenter({ message }: { message: AssistantMessage }) {
  // const classesLine = cn("hexlet-basics-content mt-2", {
  //   "text-end bg-light ms-5 p-3 rounded": message.role === "user",
  // });
  return <MarkdownViewer>{message.content}</MarkdownViewer>;
}

export default function Chat({
  userCode,
  enabled,
  focusesCount,
  output,
  course,
  lesson,
  lessonMember,
  previousMessages,
}: Props) {
  const { t: tViews } = useTranslation('web');
  const { t: tHelpers } = useTranslation('helpers');
  const { t: tCommon } = useTranslation('common');
  const inputRef = useRef<HTMLTextAreaElement>(null);

  // biome-ignore lint/correctness/useExhaustiveDependencies: -
  useEffect(() => {
    inputRef.current?.focus();
  }, [focusesCount]);

  const { input, status, messages, submitMessage, handleInputChange } =
    useAssistantStream(lessonMember?.id, lesson.id, userCode, output);

  if (!lessonMember || !course.openai_assistant_id || !enabled) {
    let content = '';

    if (!lessonMember) {
      content = tViews('languages.lessons.show.chat.guest');
    } else if (!enabled) {
      content = tViews('languages.lessons.show.chat.disabled_html');
    } else if (!course.openai_assistant_id) {
      content = tViews('languages.lessons.show.chat.not_available');
    }

    const disabledMessage: AssistantMessage = {
      role: 'assistant',
      content: content,
    };

    return <MessagePresenter message={disabledMessage} />;
  }

  const initMessage: AssistantMessage = {
    role: 'assistant',
    content: tViews('languages.lessons.show.chat.hi'),
  };

  return (
    <Stack>
      <MessagePresenter message={initMessage} />
      {previousMessages.map((m: AssistantMessage) => (
        <MessagePresenter key={m.id} message={m} />
      ))}
      {messages.map((m: AssistantMessage) => (
        <MessagePresenter key={m.id} message={m} />
      ))}

      <form onSubmit={submitMessage}>
        <Textarea
          ref={inputRef}
          disabled={status === 'in_progress'}
          value={input}
          onChange={handleInputChange}
          rows={5}
        />
        <Group justify="flex-end" pt="md">
          {i18next.language === 'ru' && (
            <Button
              component="a"
              variant="light"
              href={tCommon('community_url')}
              target="_blank"
            >
              {tViews('languages.lessons.show.chat.community')}
            </Button>
          )}
          <Button disabled={status !== 'awaiting_message'} type="submit">
            {status === 'in_progress' && <Loader size="sm" mr="xs" />}
            {tHelpers('send')}
          </Button>
        </Group>
      </form>
    </Stack>
  );
}
