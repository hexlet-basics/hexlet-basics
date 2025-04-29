import * as Routes from "@/routes.js";
import { createConsumer } from "@rails/actioncable";
import axios from "axios";
import { debounce } from "es-toolkit";
import { useSnackbar } from "notistack";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

const url = `wss://${import.meta.env.VITE_APP_HOST}/cable`;
const cableInstance = createConsumer(url);

export type AssistantMessage = {
  id?: string;
  role: "user" | "assistant" | "system";
  content: string;
};

type StreamMessage = {
  delta: string[];
  message_id: string;
  index: number;
};

export function useAssistantStream(
  lessonMemberId: number,
  lessonId: number,
  userCode: string,
  output: string,
) {
  const { enqueueSnackbar } = useSnackbar();

  const [messages, setMessages] = useState<AssistantMessage[]>([]);
  const [input, setInput] = useState("");
  const [status, setStatus] = useState<"awaiting_message" | "in_progress">(
    "awaiting_message",
  );

  const { t: tViews } = useTranslation("web");

  // Buffer per message_id
  const buffers = useRef<Record<string, string[][]>>({});

  // Debounced update to reduce flickering + broken words
  const scheduleUpdate = useRef(
    debounce((messageId: string) => {
      const parts = buffers.current[messageId];
      const fullText = parts.join("");

      setMessages((prev) => {
        const existing = prev.find((m) => m.id === messageId);

        if (existing) {
          return prev.map((m) =>
            m.id === messageId ? { ...m, content: fullText } : m,
          );
        }

        return [
          ...prev,
          { id: messageId, role: "assistant", content: fullText },
        ];
      });
    }, 50),
  ).current;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    const subscription = cableInstance.subscriptions.create(
      { channel: "AssistantChannel", id: lessonMemberId },
      {
        connected() {
          console.log("connected");
        },
        disconnected() {
          setStatus("awaiting_message");
          console.log("disconnected");
        },
        rejected() {
          console.log("rejected");
        },

        received(data: StreamMessage) {
          const { message_id, delta, index } = data;
          console.log(delta);

          if (!buffers.current[message_id]) {
            buffers.current[message_id] = [];
          }

          if (delta[0] === "DONE") {
            console.log("📡 Stream finished for message:", message_id);
            setStatus("awaiting_message");
          } else {
            buffers.current[message_id][index] = delta;
          }

          scheduleUpdate(message_id);
        },
      },
    );

    return () => {
      subscription.unsubscribe();
      scheduleUpdate.cancel();
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
        user_code: userCode,
        output,
      });
    } catch (error) {
      setStatus("awaiting_message");
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 429) {
          const message = tViews("languages.lessons.show.chat.disabled");

          setMessages((prev) => [
            ...prev,
            {
              id: crypto.randomUUID(),
              role: "assistant",
              content: message,
            },
          ]);
        }
        enqueueSnackbar(error.message);
        console.error(error);
      } else {
        throw error;
      }
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
