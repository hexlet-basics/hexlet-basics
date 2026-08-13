import { CodeHighlight, InlineCodeHighlight } from "@mantine/code-highlight";
import { Table, type TypographyProps, Typography } from "@mantine/core";
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
import { classifyLanguage, plainTextLanguage } from "@/lib/shiki";

// Markdown renderer for authored content (lesson theory, instructions, tips),
// ported from legacy. Content is GFM with raw HTML allowed, so tables, external
// links and inline markup all render the way course authors wrote them.

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

// A fenced block (```lang ... ```) is handled through `pre`, not `code`.
// react-markdown v10 dropped `code`'s `inline` flag, and a fence with no
// language gets no `language-*` class — keying off that class alone would treat
// such a block as inline code and collapse its newlines (legacy issue #592).
// The `pre` hast node carries both the language and the original text.
function MarkdownCodeBlock({ node }: ComponentPropsWithoutRef<"pre"> & ExtraProps) {
  const codeNode = node?.children.find(
    (child): child is HastElement => child.type === "element" && child.tagName === "code",
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

  // A language with no shiki grammar is never passed through: the highlighter
  // throws "Language `x` not found" on it (legacy issue #597). Unknown tags fall
  // back to plain text.
  const language = requestedLanguage
    ? classifyLanguage(requestedLanguage).language
    : plainTextLanguage;

  return <CodeHighlight code={code} language={language} />;
}

// Only inline code (single backticks) reaches the `code` handler; fenced blocks
// are intercepted by `pre` above.
function MarkdownInlineCode({ children }: ComponentPropsWithoutRef<"code">) {
  return <InlineCodeHighlight code={String(children)} />;
}

// Headings inside authored prose are one step down from the page's own, so a
// theory `##` never competes with the lesson title.
const typographyStyles: TypographyProps["styles"] = (theme) => ({
  root: {
    overflowWrap: "break-word",
    lineHeight: theme.lineHeights.lg,
    h2: theme.headings.sizes.h5,
    h3: theme.headings.sizes.h6,
    h4: theme.headings.sizes.h6,
    h5: theme.headings.sizes.h6,
    h6: theme.headings.sizes.h6,
  },
});

export default function MarkdownViewer({
  children,
  allowHtml = false,
  components = {},
}: MarkdownViewerProps) {
  const preparedComponents = {
    pre: MarkdownCodeBlock,
    // A wide table scrolls inside its own container instead of stretching the
    // pane it sits in.
    table: (props: ComponentPropsWithoutRef<"table">) => (
      <Table.ScrollContainer minWidth={800}>
        <table>{props.children}</table>
      </Table.ScrollContainer>
    ),
    code: MarkdownInlineCode,
    ...components,
  };

  const rehypePlugins: PluggableList = [
    [rehypeExternalLinks, { target: "_blank", rel: ["noopener", "noreferrer"] }],
  ];

  if (allowHtml) {
    rehypePlugins.push(rehypeRaw);
  }

  return (
    <Typography styles={typographyStyles}>
      <Markdown
        skipHtml={!allowHtml}
        remarkPlugins={[remarkGfm, remarkDirective, createDirectivePlugin(components)]}
        rehypePlugins={rehypePlugins}
        components={preparedComponents}
      >
        {children}
      </Markdown>
    </Typography>
  );
}
