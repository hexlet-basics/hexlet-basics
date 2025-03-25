import * as Routes from "@/routes.js";
import type { LanguageLesson, LanguageLessonMember } from "@/types/serializers";
import { type Message, useAssistant } from "@ai-sdk/react";
import { useEffect } from "react";
import { Button, Form, Spinner } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = {
  lesson: LanguageLesson;
  lessonMember?: LanguageLessonMember;
};

export default function Chat({ lesson, lessonMember }: Props) {
  const { status, messages, input, submitMessage, error, handleInputChange } =
    useAssistant({
      threadId: lessonMember
        ? lessonMember.openai_thread_id || undefined
        : undefined,
      api: Routes.ai_lesson_messages_path(lesson.id),
    });

  const { t: tHelpers } = useTranslation("helpers");

  useEffect(() => {
    if (error) {
      // toast.error(error.message);
    }
  }, [error]);

  return (
    <div>
      {messages.map((m: Message) => (
        <div key={m.id}>
          <strong>{`${m.role}: `}</strong>
          {m.role !== "data" && m.content}
          {m.role === "data" && (
            <>
              {/* {m.data!.description} */}
              <br />
              <pre className={"bg-gray-200"}>
                {JSON.stringify(m.data, null, 2)}
              </pre>
            </>
          )}
        </div>
      ))}

      {status === "in_progress" && <Spinner />}

      <form onSubmit={submitMessage}>
        <Form.Control
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
            {tHelpers("send")}
          </Button>
        </div>
        {/* <input */}
        {/*   disabled={status !== "awaiting_message"} */}
        {/*   value={input} */}
        {/*   placeholder="What is the temperature in the living room?" */}
        {/*   onChange={handleInputChange} */}
        {/* /> */}
      </form>
    </div>
  );
}
