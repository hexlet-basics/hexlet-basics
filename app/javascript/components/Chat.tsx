import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
} from "@/types/serializers";
import { type Message, useAssistant } from "@ai-sdk/react";
import { usePage } from "@inertiajs/react";
import axios from "axios";
import cn from "classnames";
import { useEffect, useRef } from "react";
import { Button, Form, Spinner } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { useInView } from "react-intersection-observer";
import Markdown from "react-markdown";
import useSWRImmutable from "swr/immutable";

type Props = {
  lesson: LanguageLesson;
  course: Language;
  lessonMember?: LanguageLessonMember;
};

const fetcher = (url: string) => axios.get(url).then((res) => res.data);

function MessagePresenter({ message }: { message: Message }) {
  const classesLine = cn("hexlet-basics-content mt-2", {
    "text-end bg-light ms-5 p-3 rounded": message.role === "user",
  });
  return (
    <div className={classesLine}>
      <Markdown>{message.content}</Markdown>
    </div>
  );
}

// https://sdk.vercel.ai/cookbook/next/stream-assistant-response
export default function Chat({ course, lesson, lessonMember }: Props) {
  const { t: tViews } = useTranslation("web");
  const {
    auth: { user },
  } = usePage<SharedProps>().props;
  const formRef = useRef<HTMLTextAreaElement>(null);
  const { ref } = useInView({
    threshold: 0,
    // onChange: () => formRef.current?.focus(),
  });

  const { status, messages, input, submitMessage, error, handleInputChange } =
    useAssistant({
      threadId: lessonMember
        ? lessonMember.openai_thread_id || undefined
        : undefined,
      api: Routes.ai_lesson_messages_path(lesson.id),
    });

  useEffect(() => {
    if (status !== "in_progress") {
      formRef.current?.focus();
    }
  }, [status]);

  const result = useSWRImmutable<Message[]>(
    user.guest ? null : Routes.ai_lesson_messages_path(lesson.id),
    fetcher,
  );
  const previousMessages = result.data ?? [];

  const { t: tHelpers } = useTranslation("helpers");

  useEffect(() => {
    if (error) {
      // toast.error(error.message);
      console.log(error);
    }
  }, [error]);

  let content: string;

  if (user.guest) {
    content = tViews("languages.lessons.show.chat.guest");
  } else if (course.openai_assistant_id) {
    content = tViews("languages.lessons.show.chat.hi");
  } else {
    content = tViews("languages.lessons.show.chat.disabled");
  }

  const initMessage: Message = {
    id: "0",
    role: "assistant",
    content,
  };

  return (
    <div ref={ref} className="h-100">
      <div className="mb-3">
        <MessagePresenter message={initMessage} />
        {previousMessages.map((m: Message) => (
          <MessagePresenter key={m.id} message={m} />
        ))}
        {messages.map((m: Message) => (
          <MessagePresenter key={m.id} message={m} />
        ))}
      </div>

      {course.openai_assistant_id && !user.guest && (
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
      )}
    </div>
  );
}
