import { Button, Group, Loader, Stack, Textarea } from "@mantine/core";
import i18next from "i18next";
import { useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useAssistantStream } from "@/hooks/useAssistantStream";
import type { AssistantMessage } from "@/types/assistantMessage";
import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
} from "@/types/serializers";
import MarkdownViewer from "./MarkdownViewer";

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
  const { t } = useTranslation();
  const inputRef = useRef<HTMLTextAreaElement>(null);

  // biome-ignore lint/correctness/useExhaustiveDependencies: focus is intentionally retriggered by tab selection count
  useEffect(() => {
    inputRef.current?.focus();
  }, [focusesCount]);

  const { input, status, messages, submitMessage, handleInputChange } =
    useAssistantStream(lessonMember?.id, lesson.id, userCode, output);

  if (!lessonMember || !course.openai_assistant_id || !enabled) {
    let content = "";

    if (!lessonMember) {
      content = t(($) => $.languages.lessons.show.chat.guest);
    } else if (!enabled) {
      content = t(($) => $.languages.lessons.show.chat.disabled_html);
    } else if (!course.openai_assistant_id) {
      content = t(($) => $.languages.lessons.show.chat.not_available);
    }

    const disabledMessage: AssistantMessage = {
      role: "assistant",
      content,
    };

    return <MessagePresenter message={disabledMessage} />;
  }

  const initMessage: AssistantMessage = {
    role: "assistant",
    content: t(($) => $.languages.lessons.show.chat.hi),
  };

  return (
    <Stack>
      <MessagePresenter message={initMessage} />
      {previousMessages.map((message) => (
        <MessagePresenter key={message.id} message={message} />
      ))}
      {messages.map((message) => (
        <MessagePresenter key={message.id} message={message} />
      ))}
      <form onSubmit={submitMessage}>
        <Textarea
          ref={inputRef}
          disabled={status === "in_progress"}
          value={input}
          onChange={handleInputChange}
          rows={5}
        />
        <Group justify="flex-end" pt="md">
          {i18next.language === "ru" && (
            <Button
              component="a"
              variant="light"
              href={t(($) => $.common.community_url)}
              target="_blank"
            >
              {t(($) => $.languages.lessons.show.chat.community)}
            </Button>
          )}
          <Button disabled={status !== "awaiting_message"} type="submit">
            {status === "in_progress" && <Loader size="sm" mr="xs" />}
            {t(($) => $.helpers.send)}
          </Button>
        </Group>
      </form>
    </Stack>
  );
}
