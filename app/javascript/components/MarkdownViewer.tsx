import renderMarkdown, { type RenderMarkdownOptions } from "@/lib/markdown";
import { useEffect, useState } from "react";

type MarkdownViewerProps = RenderMarkdownOptions & {
  children: string;
};

export default function MarkdownViewer({
  children,
  foldCode,
  foldCodeMessage,
  inlineCodeHighlight,
}: MarkdownViewerProps) {
  const [html, setHtml] = useState<string>("");

  useEffect(() => {
    let cancelled = false;

    const render = async () => {
      const result = await renderMarkdown(children, {
        foldCode,
        foldCodeMessage,
        inlineCodeHighlight,
      });
      if (!cancelled) setHtml(result);
    };

    render();

    return () => {
      cancelled = true;
    };
  }, [children, foldCode, foldCodeMessage, inlineCodeHighlight]);

  return (
    <div
      className="markdown-body"
      // biome-ignore lint/security/noDangerouslySetInnerHtml: markdown is sanitized by pipeline
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
