import { useChat } from "@ai-sdk/react";
import {
  Anchor,
  Button,
  Group,
  List,
  Paper,
  Skeleton,
  Stack,
  Text,
  Textarea,
  ThemeIcon,
} from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { APICallError, TextStreamChatTransport } from "ai";
import type { UIMessage } from "ai";
import { IconCheck, IconLock } from "@tabler/icons-react";
import { useLocation } from "@tanstack/react-router";
import { useState } from "react";
import type { SubmitEvent } from "react";
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
import { ButtonLink } from "@/components/RouterLink";
import { isApiError } from "@/lib/api-error";
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
    return <GuestAssistant />;
  }

  if (history.isPending) {
    return (
      <Stack p="lg" gap="sm">
        <Skeleton h={14} w="70%" />
        <Skeleton h={14} />
        <Skeleton h={14} w="40%" />
      </Stack>
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

// What a visitor sees instead of the chat: who Tota is, what she does for a
// learner, a sample exchange in the real chat's bubbles, and the way in. The
// sign-up and sign-in links carry the way back to this lesson.
function GuestAssistant() {
  const { t } = useTranslation();
  const redirect = useLocation({ select: (location) => location.href });
  const benefits = [
    t(($) => $.courses.lessons.show.chat.guest_benefit_hint),
    t(($) => $.courses.lessons.show.chat.guest_benefit_errors),
    t(($) => $.courses.lessons.show.chat.guest_benefit_free),
  ];

  return (
    <Stack p="lg" gap="lg">
      <Paper p="sm" radius="lg" bg="var(--mantine-color-default-hover)">
        <Stack gap="sm">
          <Text size="sm">{t(($) => $.courses.lessons.show.chat.guest)}</Text>
          <List
            spacing="xs"
            size="sm"
            icon={
              <ThemeIcon variant="transparent" size={20}>
                <IconCheck size={16} />
              </ThemeIcon>
            }
          >
            {benefits.map((text) => (
              <List.Item key={text}>{text}</List.Item>
            ))}
          </List>
        </Stack>
      </Paper>

      <Stack gap="xs" aria-hidden>
        <Group justify="flex-end">
          <Paper p="sm" radius="lg" maw="85%" bg="var(--mantine-primary-color-light)">
            <Text size="sm">{t(($) => $.courses.lessons.show.chat.guest_example_question)}</Text>
          </Paper>
        </Group>
        <Group justify="flex-start">
          <Paper p="sm" radius="lg" maw="85%" bg="var(--mantine-color-default-hover)">
            <Text size="sm">{t(($) => $.courses.lessons.show.chat.guest_example_answer)}</Text>
          </Paper>
        </Group>
        <Textarea
          disabled
          rows={2}
          leftSection={<IconLock size={16} />}
          placeholder={t(($) => $.courses.lessons.show.chat.question)}
          tabIndex={-1}
        />
      </Stack>

      <Group>
        <ButtonLink to="/{-$locale}/users/new" search={{ redirect }}>
          {t(($) => $.courses.lessons.show.chat.guest_sign_up)}
        </ButtonLink>
        <ButtonLink to="/{-$locale}/session/new" search={{ redirect }} variant="default">
          {t(($) => $.courses.lessons.show.sign_in)}
        </ButtonLink>
      </Group>
    </Stack>
  );
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
      if (isQuotaExceeded(failure)) {
        return;
      }
      notifications.show({ message: t(($) => $.common.errors.network) });
    },
  });

  const busy = status === "submitted" || status === "streaming";
  const quotaExceeded = chat.quotaExceeded || isQuotaExceeded(error);
  const communityUrl = t(($) => $.common.community_url).trim();

  const submit = (event: SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    const text = input.trim();
    if (!text || busy) {
      return;
    }

    // The editor state travels with each question, read at the moment it is
    // asked rather than when the chat was created.
    void sendMessage({ text }, { body: { userCode: code, output } });
    setInput("");
  };

  return (
    <Stack p="lg">
      <MarkdownViewer>{t(($) => $.courses.lessons.show.chat.hi)}</MarkdownViewer>

      {messages.map((message) => {
        const own = message.role === "user";
        // Messenger layout: the learner's replies are filled and sit on the
        // right, the assistant's are neutral and sit on the left.
        return (
          <Group key={message.id} justify={own ? "flex-end" : "flex-start"}>
            <Paper
              p="sm"
              radius="lg"
              maw="85%"
              bg={own ? "var(--mantine-primary-color-light)" : "var(--mantine-color-default-hover)"}
            >
              <MarkdownViewer>{textOf(message)}</MarkdownViewer>
            </Paper>
          </Group>
        );
      })}

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
// generated client — so the path, the session cookie and URL
// locale come from the same place as every other call. `parseAs: "stream"`
// hands back the body as a stream, which the transport reads as plain text.
// The request body the chat hook forwards is untyped; only strings are sent on.
function stringOrNull(value: unknown): string | null {
  return typeof value === "string" ? value : null;
}

function assistantTransport(lessonId: number) {
  return new TextStreamChatTransport<UIMessage>({
    prepareSendMessagesRequest: ({ messages, body }) => {
      const input: AssistantMessageInput = {
        message: textOf(messages.at(-1)),
        userCode: stringOrNull(body?.userCode),
        output: stringOrNull(body?.output),
      };
      return { body: input };
    },
    fetch: async (_url, init) => {
      // The AI SDK serializes the prepared body to a JSON string.
      if (typeof init?.body !== "string") {
        throw new TypeError("assistant transport expects a JSON string body");
      }
      try {
        const { data, response } = await createAssistantMessage({
          path: { lessonId },
          body: JSON.parse(init.body) as AssistantMessageInput,
          parseAs: "stream",
          signal: init?.signal ?? undefined,
          throwOnError: true,
        });
        return new Response(data, { status: response.status });
      } catch (error) {
        // Every status is handed to the transport, which reports a refusal
        // (the quota's 429 among them) as an APICallError with its status.
        if (isApiError(error)) {
          return new Response(null, { status: error.status });
        }
        throw error;
      }
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
