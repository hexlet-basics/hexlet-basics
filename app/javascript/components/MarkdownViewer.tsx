import { CodeHighlight, InlineCodeHighlight } from "@mantine/code-highlight";
import { Table, Typography } from "@mantine/core";
// import rehypeShikiFromHighlighter from '@shikijs/rehype/core';
import type { Directives } from "mdast-util-directive";
import type { ComponentPropsWithoutRef } from "react";
import Markdown from "react-markdown";
import rehypeExternalLinks from "rehype-external-links";
import rehypeRaw from "rehype-raw";
import remarkDirective from "remark-directive";
import remarkGfm from "remark-gfm";
import type { PluggableList } from "unified";
import type { Node } from "unist";
import { visit } from "unist-util-visit";
import { typographyStyles } from "@/lib/mantine";

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

// Адаптер для CodeHighlight
function MarkdownCodeHighlight({
  className,
  children,
}: ComponentPropsWithoutRef<"code">) {
  const code = String(children).trim();
  const match = className?.match(/language-(\w+)/);
  const language = match ? match[1] : "plaintext";
  const isInline = !match;

  if (isInline) {
    return <InlineCodeHighlight code={code} />;
  }

  return <CodeHighlight code={code} language={language} />;
}

export default function MarkdownViewer({
  children,
  allowHtml = false,
  components = {},
}: MarkdownViewerProps) {
  const preparedComponents = {
    pre: (props: ComponentPropsWithoutRef<"pre">) => <>{props.children}</>, // убираем обертку pre
    table: (props: ComponentPropsWithoutRef<"table">) => (
      <Table.ScrollContainer minWidth={800}>
        <table>{props.children}</table>
      </Table.ScrollContainer>
    ), // убираем обертку pre
    ul: (props: ComponentPropsWithoutRef<"ul">) => (
      <ul className="list-disc ml-6">{props.children}</ul>
    ),
    ol: (props: ComponentPropsWithoutRef<"ol">) => (
      <ol className="list-decimal ml-6">{props.children}</ol>
    ),
    li: (props: ComponentPropsWithoutRef<"li">) => (
      <li className="my-1">{props.children}</li>
    ),
    code: MarkdownCodeHighlight,
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
    <Typography styles={typographyStyles} className="prose max-w-none">
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
