import * as Routes from "@/routes.js";
import { createConsumer } from "@rails/actioncable";
import axios from "axios";
import { useEffect, useRef, useState } from "react";

const url = `wss://${import.meta.env.VITE_APP_HOST}/cable`;
const cableInstance = createConsumer(url);

export type AssistantMessage = {
  id?: string;
  role: "user" | "assistant";
  content: string;
};

type StreamMessage = {
  delta: string;
  message_id: string;
};

export function useAssistantStream(lessonMemberId: number, lessonId: number) {
  const [messages, setMessages] = useState<AssistantMessage[]>([]);
  const [input, setInput] = useState("");
  const [status, setStatus] = useState<"awaiting_message" | "in_progress">(
    "awaiting_message",
  );

  // Buffer per message_id
  const buffers = useRef<Record<string, string>>({});

  useEffect(() => {
    const subscription = cableInstance.subscriptions.create(
      { channel: "AssistantChannel", id: lessonMemberId },
      {
        connected() {
          console.log("connected");
        },
        disconnected() {
          console.log("disconnected");
        },
        rejected() {
          console.log("rejected");
        },

        received(data: StreamMessage) {
          const { message_id, delta } = data;
          if (!buffers.current[message_id]) {
            buffers.current[message_id] = "";
          }

          buffers.current[message_id] += delta;

          setMessages((prev) => {
            const existing = prev.find((m) => m.id === message_id);

            if (existing) {
              return prev.map((m) =>
                m.id === message_id
                  ? { ...m, content: buffers.current[message_id] }
                  : m,
              );
            }

            return [
              ...prev,
              {
                id: message_id,
                role: "assistant",
                content: buffers.current[message_id],
              },
            ];
          });
        },
      },
    );

    return () => {
      subscription.unsubscribe();
    };
  }, [lessonMemberId]);

  const submitMessage = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!input.trim()) return;

    setStatus("in_progress");

    const userMessage: AssistantMessage = {
      id: crypto.randomUUID(),
      role: "user",
      content: input,
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput("");

    try {
      await axios.post(Routes.ai_lesson_messages_path(lessonId), {
        message: input,
      });
    } catch (error) {
      console.error("Failed to send message", error);
    } finally {
      setStatus("awaiting_message");
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
  };

  return {
    messages,
    input,
    handleInputChange,
    submitMessage,
    status,
  };
}
