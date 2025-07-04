import { MarkdownHooks } from "react-markdown";
import rehypeShikiFromHighlighter from "@shikijs/rehype/core";
import remarkGfm from "remark-gfm";
import rehypeExternalLinks from "rehype-external-links";
import rehypeRaw from "rehype-raw";
import type { PluggableList } from "unified";
import { Box, TypographyStylesProvider, useMantineColorScheme } from "@mantine/core";
import getHighlighter from "@/lib/shiki";
import { useEffect, useState } from "react";

type MarkdownViewerProps = {
  children: string;
  allowHtml?: boolean;
};

export default function MarkdownViewer({
  children,
  allowHtml = false,
}: MarkdownViewerProps) {
  // const { colorScheme } = useMantineColorScheme();
  const [rehypePlugins, setRehypePlugins] = useState<PluggableList>([]);

  useEffect(() => {
    const loadPlugins = async () => {
      const highlighter = await getHighlighter();

      const rehypePlugins: PluggableList = [
        [rehypeShikiFromHighlighter, highlighter, {
          themes: { light: "github-light", dark: "github-dark" },
          theme: 'github-light',
        }],
        [rehypeExternalLinks, {
          target: "_blank",
          rel: ["noopener", "noreferrer"],
        }],
      ];

      if (allowHtml) {
        rehypePlugins.push(rehypeRaw)
      }

      setRehypePlugins(rehypePlugins);
    }

    loadPlugins();
  }, [allowHtml]);

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
