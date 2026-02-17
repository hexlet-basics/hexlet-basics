import { Alert, Anchor, Group, Stack } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { ChatKit, useChatKit } from "@openai/chatkit-react";
import { IconInfoCircle } from "@tabler/icons-react";
import i18next from "i18next";
import { Trans, useTranslation } from "react-i18next";
import type { AssistantMessage } from "@/types/assistantMessage";
import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
} from "@/types/serializers";
import AppAnchor from "./Elements/AppAnchor";
import MarkdownViewer from "./MarkdownViewer";

type Props = {
  enabled: boolean;
  lesson: LanguageLesson;
  course: Language;
  userCode: string;
  output: string;
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
  course,
  output,
  lessonMember,
  lesson,
}: Props) {
  const { t } = useTranslation();

  const apiUrl = "/api/chatkit";
  const domainKey = import.meta.env.VITE_CHATKIT_DOMAIN_KEY ?? "local-dev";
  const lessonMemberId = lessonMember?.id;

  const { control } = useChatKit({
    api: {
      url: apiUrl,
      domainKey,
      fetch: async (input, init) => {
        const body = init?.body;
        const payload =
          typeof body === "string" ? JSON.parse(body) : (body ?? {});
        const metadata = {
          ...(payload.metadata ?? {}),
          lesson_id: lesson.id,
          lesson_member_id: lessonMemberId,
          user_code: userCode,
          output,
        };

        return fetch(input, {
          ...init,
          credentials: "include",
          headers: {
            ...init?.headers,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            ...payload,
            metadata,
          }),
        });
      },
    },
    locale: i18next.language as "en" | "ru-RU" | "es-ES" | "es-419" | undefined,
    initialThread: lessonMember?.openai_thread_id ?? null,
    composer: {
      attachments: {
        enabled: false,
      },
    },
    onError: ({ error }) => {
      notifications.show({
        message: error.message,
      });
    },
  });

  if (!lessonMember || !course.openai_assistant_id || !enabled) {
    if (!enabled) {
      return (
        <Trans
          t={t}
          i18nKey={($) => $.languages.lessons.show.chat.disabled_html}
          components={{
            a: <AppAnchor external href="https://t.me/HexletLearningBot" />,
          }}
        />
      );
    }

    let content = "";

    if (!lessonMember) {
      content = t(($) => $.languages.lessons.show.chat.guest);
    } else if (!course.openai_assistant_id) {
      content = t(($) => $.languages.lessons.show.chat.not_available);
    }

    const disabledMessage: AssistantMessage = {
      role: "assistant",
      content: content,
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
      {i18next.language === "ru" && (
        <Alert icon={<IconInfoCircle size={16} />} variant="light">
          <Group gap={6}>
            <span>{t(($) => $.languages.lessons.show.chat.community)}</span>
            <Anchor href={t(($) => $.common.community_url)} target="_blank">
              {t(($) => $.common.community_url)}
            </Anchor>
          </Group>
        </Alert>
      )}
      <ChatKit control={control} className="h-[560px] w-full" />
      {/*
        ChatKit web component is loaded from CDN as required by @openai/chatkit-react.
        Script tag is injected in application layout.
      */}
    </Stack>
  );
}
