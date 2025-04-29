import {
  type AssistantMessage,
  useAssistantStream,
} from "@/hooks/useAssistantStream";

import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
} from "@/types/serializers";
import cn from "classnames";
import { useRef } from "react";
import { Button, Form, Spinner } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { useInView } from "react-intersection-observer";
import MarkdownViewer from "./MarkdownViewer";

type Props = {
  enabled: boolean
  lesson: LanguageLesson;
  course: Language;
  userCode: string;
  output: string;
  previousMessages: AssistantMessage[];
  lessonMember?: LanguageLessonMember;
};

function MessagePresenter({ message }: { message: AssistantMessage }) {
  const classesLine = cn("hexlet-basics-content mt-2", {
    "text-end bg-light ms-5 p-3 rounded": message.role === "user",
  });
  return (
    <div className={classesLine}>
      <MarkdownViewer>{message.content}</MarkdownViewer>
    </div>
  );
}

// https://sdk.vercel.ai/cookbook/next/stream-assistant-response
export default function Chat({
  userCode,
  enabled,
  output,
  course,
  lesson,
  lessonMember,
  previousMessages,
}: Props) {
  const { t: tViews } = useTranslation("web");

  if (!lessonMember || !course.openai_assistant_id || !enabled) {
    let content = "";

    if (!enabled) {
      content = tViews("languages.lessons.show.chat.disabled");
    } else if (!lessonMember) {
      content = tViews("languages.lessons.show.chat.guest");
    } else if (!course.openai_assistant_id) {
      content = tViews("languages.lessons.show.chat.not_available");
    }

    const disabledMessage: AssistantMessage = {
      role: "assistant",
      content: content,
    };

    return (
      <div className="h-100">
        <div className="mb-3">
          <MessagePresenter message={disabledMessage} />
        </div>
      </div>
    );
  }

  const formRef = useRef<HTMLTextAreaElement>(null);
  const { ref } = useInView({
    threshold: 0,
    // onChange: () => formRef.current?.focus(),
  });

  const { input, status, messages, submitMessage, handleInputChange } =
    useAssistantStream(lessonMember.id, lesson.id, userCode, output);

  const { t: tHelpers } = useTranslation("helpers");

  const initMessage: AssistantMessage = {
    role: "assistant",
    content: tViews("languages.lessons.show.chat.hi"),
  };

  return (
    <div ref={ref} className="h-100">
      <div className="mb-3">
        <MessagePresenter message={initMessage} />
        {previousMessages.map((m: AssistantMessage) => (
          <MessagePresenter key={m.id} message={m} />
        ))}
        {messages.map((m: AssistantMessage) => (
          <MessagePresenter key={m.id} message={m} />
        ))}
      </div>

      <form onSubmit={submitMessage}>
        <Form.Control
          ref={formRef}
          disabled={status === "in_progress"}
          as="textarea"
          value={input}
          onChange={handleInputChange}
          aria-label="With textarea"
        />
        <div className="d-flex justify-content-end">
          <Button
            disabled={status !== "awaiting_message"}
            className="mt-3"
            type="submit"
          >
            {status === "in_progress" && (
              <Spinner
                className="me-2"
                as="span"
                animation="grow"
                size="sm"
                role="status"
                aria-hidden="true"
              />
            )}
            {tHelpers("send")}
          </Button>
        </div>
      </form>
    </div>
  );
}
