import { CodeHighlight, InlineCodeHighlight } from "@mantine/code-highlight";
import { Table, Typography } from "@mantine/core";
import * as Sentry from "@sentry/react";
// import rehypeShikiFromHighlighter from '@shikijs/rehype/core';
import type { Directives } from "mdast-util-directive";
import type { ComponentPropsWithoutRef } from "react";
import Markdown, { type ExtraProps } from "react-markdown";
import rehypeExternalLinks from "rehype-external-links";
import rehypeRaw from "rehype-raw";
import remarkDirective from "remark-directive";
import remarkGfm from "remark-gfm";
import type { PluggableList } from "unified";
import type { Node } from "unist";
import { visit } from "unist-util-visit";
import { typographyStyles } from "@/lib/mantine";
import { plainTextLanguage, supportedLanguages } from "@/lib/shiki";

type HastElement = NonNullable<ExtraProps["node"]>;

type DirectiveComponents = Record<
  string,
  React.ComponentType<{
    attributes?: Record<string, string>;
    children?: React.ReactNode;
  }>
>;

type MarkdownViewerProps = {
  children: string;
  allowHtml?: boolean;
  components?: DirectiveComponents;
};

function isDirective(node: Node): node is Directives {
  return (
    node.type === "containerDirective" ||
    node.type === "leafDirective" ||
    node.type === "textDirective"
  );
}

function createDirectivePlugin(components: DirectiveComponents) {
  return function directivePlugin() {
    return (tree: Node) => {
      visit(tree, (node: Node) => {
        if (isDirective(node)) {
          const componentKey = `::${node.name}`;
          const component = components[componentKey];
          if (!component) return;

          if (!node.data) node.data = {};
          const data = node.data;
          data.hName = componentKey;
          data.hProperties = node.attributes || {};
        }
      });
    };
  };
}

// Блок кода (```lang ... ```) рендерим через хендлер `pre`, а не `code`.
// В react-markdown v10 у `code` больше нет флага `inline`, а блок без указания
// языка не получает класс `language-*`. Если ориентироваться только на этот класс,
// такой блок ошибочно считается inline и переносы строк схлопываются (issue #592).
// В `pre` же по hast-ноде надёжно виден и язык, и исходный текст.
function MarkdownCodeBlock({
  node,
}: ComponentPropsWithoutRef<"pre"> & ExtraProps) {
  const codeNode = node?.children.find(
    (child): child is HastElement =>
      child.type === "element" && child.tagName === "code",
  );

  const rawClassName = codeNode?.properties?.className;
  const className = Array.isArray(rawClassName)
    ? rawClassName.join(" ")
    : String(rawClassName ?? "");
  const requestedLanguage = className.match(/language-(\w+)/)?.[1];

  const code = (codeNode?.children ?? [])
    .map((child) => (child.type === "text" ? child.value : ""))
    .join("")
    .replace(/\n$/, "");

  // Неизвестный язык (нет соответствующей грамматики shiki) не рендерим как есть —
  // иначе highlighter падает с "Language `x` not found" (issue #597). Откатываемся
  // на обычный текст и сообщаем в Sentry, чтобы не терять такие случаи.
  let language = plainTextLanguage;
  if (requestedLanguage) {
    if (supportedLanguages.has(requestedLanguage)) {
      language = requestedLanguage;
    } else {
      Sentry.captureMessage(
        `MarkdownViewer: unsupported code language "${requestedLanguage}", falling back to plain text`,
        "warning",
      );
    }
  }

  return <CodeHighlight code={code} language={language} />;
}

// Только inline-код (`код` одинарными кавычками) доходит до хендлера `code`,
// блоки перехватываются в `pre`.
function MarkdownInlineCode({ children }: ComponentPropsWithoutRef<"code">) {
  return <InlineCodeHighlight code={String(children)} />;
}

export default function MarkdownViewer({
  children,
  allowHtml = false,
  components = {},
}: MarkdownViewerProps) {
  const preparedComponents = {
    pre: MarkdownCodeBlock,
    table: (props: ComponentPropsWithoutRef<"table">) => (
      <Table.ScrollContainer minWidth={800}>
        <table>{props.children}</table>
      </Table.ScrollContainer>
    ), // убираем обертку pre
    code: MarkdownInlineCode,
    ...components,
  };

  const rehypePlugins: PluggableList = [
    [
      rehypeExternalLinks,
      {
        target: "_blank",
        rel: ["noopener", "noreferrer"],
      },
    ],
  ];

  if (allowHtml) {
    rehypePlugins.push(rehypeRaw);
  }

  return (
    <Typography styles={typographyStyles}>
      <Markdown
        skipHtml={!allowHtml}
        remarkPlugins={[
          remarkGfm,
          remarkDirective,
          createDirectivePlugin(components),
        ]}
        rehypePlugins={rehypePlugins}
        components={preparedComponents}
      >
        {children}
      </Markdown>
    </Typography>
  );
}
