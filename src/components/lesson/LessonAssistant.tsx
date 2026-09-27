import { useChat } from "@ai-sdk/react";
import {
  Anchor,
  Box,
  Button,
  Center,
  Group,
  Loader,
  Paper,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { APICallError, TextStreamChatTransport, type UIMessage } from "ai";
import { type SubmitEvent, useState } from "react";
import { Trans, useTranslation } from "react-i18next";
import {
  listAssistantMessagesOptions,
  listAssistantMessagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { createAssistantMessage } from "@/client/sdk.gen";
import type {
  AssistantMessageInput,
  CourseLessonView,
  LessonAssistantChat,
  LessonAssistantMessage,
} from "@/client/types.gen";
import MarkdownViewer from "@/components/MarkdownViewer";
import { lessonCodeKey } from "@/lib/lesson-code";

// The in-lesson assistant (Tota). A signed-in feature: a visitor is told to
// sign in, a learner gets their chat for this lesson and can ask about the code
// in the editor and the output of their last check. Ported from legacy's Chat;
// the stream is read by the AI SDK, not by hand.
export default function LessonAssistant({
  view,
  output,
  signedIn,
}: {
  view: CourseLessonView;
  output: string;
  signedIn: boolean;
}) {
  const { t } = useTranslation();
  const lessonId = view.lesson.id;
  const history = useQuery({
    ...listAssistantMessagesOptions({ path: { lessonId } }),
    enabled: signedIn,
  });

  if (!signedIn) {
    return (
      <Box p="lg">
        <MarkdownViewer>{t(($) => $.courses.lessons.show.chat.guest)}</MarkdownViewer>
      </Box>
    );
  }

  if (history.isPending) {
    return (
      <Center p="lg">
        <Loader />
      </Center>
    );
  }

  if (history.isError) {
    return (
      <Text c="red" p="lg">
        {t(($) => $.common.errors.network)}
      </Text>
    );
  }

  // Keyed by lesson: useChat takes its messages once, so another lesson needs
  // another chat rather than new props.
  return <AssistantChat key={lessonId} view={view} output={output} chat={history.data} />;
}

function AssistantChat({
  view,
  output,
  chat,
}: {
  view: CourseLessonView;
  output: string;
  chat: LessonAssistantChat;
}) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { lesson } = view;
  const historyKey = listAssistantMessagesQueryKey({ path: { lessonId: lesson.id } });
  const [code] = useLocalStorage({
    key: lessonCodeKey(lesson.course.slug, lesson.slug),
    defaultValue: lesson.preparedCode ?? "",
    getInitialValueInEffect: false,
  });
  const [input, setInput] = useState("");

  const { messages, sendMessage, status, error } = useChat({
    id: `lesson-${lesson.id}`,
    messages: chat.messages.map(toUIMessage),
    transport: assistantTransport(lesson.id),
    // Re-read the history's quota flag after every exchange: the answer that
    // used up today's last question is what switches the panel off.
    onFinish: () => {
      void queryClient.invalidateQueries({ queryKey: historyKey });
    },
    onError: (failure) => {
      if (isQuotaExceeded(failure)) return;
      notifications.show({ message: t(($) => $.common.errors.network) });
    },
  });

  const busy = status === "submitted" || status === "streaming";
  const quotaExceeded = chat.quotaExceeded || isQuotaExceeded(error);
  const communityUrl = t(($) => $.common.community_url).trim();

  const submit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    const text = input.trim();
    if (!text || busy) return;

    // The editor state travels with each question, read at the moment it is
    // asked rather than when the chat was created.
    void sendMessage({ text }, { body: { userCode: code, output } });
    setInput("");
  };

  return (
    <Stack p="lg">
      <MarkdownViewer>{t(($) => $.courses.lessons.show.chat.hi)}</MarkdownViewer>

      {messages.map((message) => (
        <Paper
          key={message.id}
          p="sm"
          withBorder={message.role === "user"}
          bg={message.role === "user" ? undefined : "var(--mantine-color-default-hover)"}
        >
          <MarkdownViewer>{textOf(message)}</MarkdownViewer>
        </Paper>
      ))}

      {quotaExceeded ? (
        <Text>
          <Trans
            t={t}
            i18nKey={($) => $.courses.lessons.show.chat.disabled_html}
            components={{
              a: <Anchor href={communityUrl} target="_blank" rel="noopener noreferrer" />,
            }}
          />
        </Text>
      ) : (
        <form onSubmit={submit}>
          <Textarea
            aria-label={t(($) => $.courses.lessons.show.chat.question)}
            placeholder={t(($) => $.courses.lessons.show.chat.question)}
            value={input}
            onChange={(event) => setInput(event.currentTarget.value)}
            disabled={busy}
            rows={5}
          />
          <Group justify="flex-end" pt="md">
            <Button
              component="a"
              variant="light"
              href={communityUrl}
              target="_blank"
              rel="noopener noreferrer"
            >
              {t(($) => $.courses.lessons.show.chat.community)}
            </Button>
            <Button type="submit" loading={busy} disabled={!input.trim()}>
              {t(($) => $.helpers.send)}
            </Button>
          </Group>
        </form>
      )}
    </Stack>
  );
}

// The transport speaks the contract rather than the AI SDK's own request shape:
// it sends the question with the editor state, and it sends it through the
// generated client — so the path, the session cookie, XSRF header and URL
// locale come from the same place as every other call. Axios' fetch adapter
// hands back the body as a stream, which the transport reads as plain text.
function assistantTransport(lessonId: number) {
  return new TextStreamChatTransport<UIMessage>({
    prepareSendMessagesRequest: ({ messages, body }) => {
      const input: AssistantMessageInput = {
        message: textOf(messages.at(-1)),
        userCode: body?.userCode ?? null,
        output: body?.output ?? null,
      };
      return { body: input };
    },
    fetch: async (_url, init) => {
      const response = await createAssistantMessage({
        path: { lessonId },
        body: JSON.parse(String(init?.body)),
        adapter: "fetch",
        responseType: "stream",
        // Every status is handed to the transport, which reports a refusal
        // (the quota's 429 among them) as an APICallError with its status.
        validateStatus: () => true,
        signal: init?.signal ?? undefined,
        throwOnError: true,
      });
      return new Response(response.data, { status: response.status });
    },
  });
}

function isQuotaExceeded(error: unknown): boolean {
  return APICallError.isInstance(error) && error.statusCode === 429;
}

function toUIMessage(message: LessonAssistantMessage): UIMessage {
  return {
    id: String(message.id),
    role: message.role === "user" ? "user" : "assistant",
    parts: [{ type: "text", text: message.content }],
  };
}

function textOf(message: UIMessage | undefined): string {
  return (message?.parts ?? [])
    .flatMap((part) => (part.type === "text" ? [part.text] : []))
    .join("");
}
