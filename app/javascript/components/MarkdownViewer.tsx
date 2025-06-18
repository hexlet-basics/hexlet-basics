import { useEffect, useState } from "react";
import { MarkdownHooks } from "react-markdown";
import rehypeShikiFromHighlighter from "@shikijs/rehype/core";
import remarkGfm from "remark-gfm";
import rehypeExternalLinks from "rehype-external-links";
import rehypeRaw from "rehype-raw";
import type { PluggableList } from "unified";
import { Box, TypographyStylesProvider, useMantineColorScheme } from "@mantine/core";
import { loadShiki } from "@/lib/shiki";

type MarkdownViewerProps = {
  children: string;
  allowHtml?: boolean;
};

export default function MarkdownViewer({
  children,
  allowHtml = false,
}: MarkdownViewerProps) {
  const [rehypePlugins, setRehypePlugins] = useState<PluggableList | null>(null);

  const { colorScheme } = useMantineColorScheme();

  useEffect(() => {
    loadShiki().then((highlighter) => {
      const plugins: PluggableList = [
        [rehypeShikiFromHighlighter, highlighter, {
          themes: { light: "github-light", dark: "github-dark" },
          theme: colorScheme,
        }],
        [rehypeExternalLinks, {
          target: "_blank",
          rel: ["noopener", "noreferrer"],
        }],
      ];

      if (allowHtml) {
        plugins.push(rehypeRaw)
      }

      setRehypePlugins(plugins);
    });
  }, [allowHtml, colorScheme]);

  if (!rehypePlugins) return null;

  // const styles = {
  //   root: {
  //     pre: {
  //       'border': '1px solid #dee2e6',
  //       'border-radius': '0.25rem',
  //     },
  //   },
  // }

  return (
    <TypographyStylesProvider>
      <Box className="markdown-viewer">
        <MarkdownHooks
          skipHtml={!allowHtml}
          remarkPlugins={[remarkGfm]}
          rehypePlugins={rehypePlugins}
        >
          {children}
        </MarkdownHooks>
      </Box>
    </TypographyStylesProvider>
  );
}
