import { usePage } from "@inertiajs/react";
import { Anchor, Box } from "@mantine/core";
import type { ComponentPropsWithoutRef } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { propsForExternalLink } from "@/lib/utils";
import type { BannerBackground } from "@/types/serializers";

// Пресеты фона. Ключи синхронизированы с Banner::Background (бэкенд).
const backgroundStyles: Record<BannerBackground, string> = {
  cta_gradient: "var(--app-cta-gradient)",
  dark: "var(--mantine-color-dark-6)",
  blue: "var(--mantine-color-blue-6)",
};

// Инлайн-рендер markdown: убираем блочную обёртку <p> и ссылки
// (ссылка задаётся целиком полем url на уровне всего баннера).
const inlineComponents = {
  p: ({ children }: ComponentPropsWithoutRef<"p">) => <>{children}</>,
  a: ({ children }: ComponentPropsWithoutRef<"a">) => <>{children}</>,
};

export function Banner() {
  const { activeBanner } = usePage().props;

  if (!activeBanner) {
    return null;
  }

  const background =
    backgroundStyles[activeBanner.background] ?? backgroundStyles.cta_gradient;

  const content = (
    <Markdown
      skipHtml
      remarkPlugins={[remarkGfm]}
      disallowedElements={["a"]}
      unwrapDisallowed
      components={inlineComponents}
    >
      {activeBanner.body}
    </Markdown>
  );

  return (
    <Box bg={background} p="xs" my="xl" ta="center">
      {activeBanner.url ? (
        <Anchor href={activeBanner.url} {...propsForExternalLink()}>
          {content}
        </Anchor>
      ) : (
        content
      )}
    </Box>
  );
}
