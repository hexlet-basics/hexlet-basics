import { Box, TypographyStylesProvider } from '@mantine/core';
// import rehypeShikiFromHighlighter from '@shikijs/rehype/core';
import type { Directives } from 'mdast-util-directive';
import { type ComponentPropsWithoutRef, useEffect, useState } from 'react';
import { MarkdownHooks } from 'react-markdown';
import rehypeExternalLinks from 'rehype-external-links';
import rehypeRaw from 'rehype-raw';
import remarkDirective from 'remark-directive';
import remarkGfm from 'remark-gfm';
import type { PluggableList } from 'unified';
import type { Node } from 'unist';
import { visit } from 'unist-util-visit';
import { LazyCodeHighlight } from './LazyCodeHighlight';

type CodeProps = ComponentPropsWithoutRef<'code'> & { node?: unknown };

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
    node.type === 'containerDirective' ||
    node.type === 'leafDirective' ||
    node.type === 'textDirective'
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

function MarkdownCodeHighlight({ className, children }: CodeProps) {
  const code = String(children).trim();
  const match = className?.match(/language-(\w+)/);
  const language = match ? match[1] : 'plaintext';
  const isInline = !match;

  if (isInline) {
    return <code className={className}>{code}</code>;
  }

  return (
    <LazyCodeHighlight code={code} language={language} className={className} />
  );
}

export default function MarkdownViewer({
  children,
  allowHtml = false,
  components = {},
}: MarkdownViewerProps) {
  const [rehypePlugins, setRehypePlugins] = useState<PluggableList>([]);

  const preparedComponents = {
    pre: (props: ComponentPropsWithoutRef<'pre'>) => <>{props.children}</>, // убираем обертку pre
    code: MarkdownCodeHighlight,
    ...components,
  };

  useEffect(() => {
    const loadPlugins = async () => {
      // const highlighter = await getBaseHighlighter();

      const rehypePlugins: PluggableList = [
        [
          rehypeExternalLinks,
          {
            target: '_blank',
            rel: ['noopener', 'noreferrer'],
          },
        ],
      ];

      if (allowHtml) {
        rehypePlugins.push(rehypeRaw);
      }

      setRehypePlugins(rehypePlugins);
    };

    loadPlugins();
  }, [allowHtml]);

  return (
    <TypographyStylesProvider>
      <Box className="markdown-viewer">
        <MarkdownHooks
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
        </MarkdownHooks>
      </Box>
    </TypographyStylesProvider>
  );
}
