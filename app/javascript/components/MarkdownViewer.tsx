// import { type RenderMarkdownOptions } from "@/lib/markdown";
// import { useEffect, useState } from "react";
import Markdown from 'react-markdown'
import rehypeHighlight from 'rehype-highlight'
import remarkGfm from 'remark-gfm'
import rehypeRaw from 'rehype-raw'
import rehypeExternalLinks from "rehype-external-links";
import { PluggableList } from "unified";

type MarkdownViewerProps = {
  children: string;
  allowHtml?: boolean
};

export default function MarkdownViewer({
  children,
  allowHtml = false,
  // foldCode,
  // foldCodeMessage,
  // inlineCodeHighlight,
}: MarkdownViewerProps) {

  const rehypePlugins: PluggableList = [
    rehypeHighlight,
    [rehypeExternalLinks, {
      target: "_blank",
      rel: ["noopener", "noreferrer"],
    }]
  ]
  if (allowHtml) {
    rehypePlugins.push(rehypeRaw)
  }

  return <Markdown
    remarkPlugins={[remarkGfm]}
    rehypePlugins={rehypePlugins}
  >
    {children}
  </Markdown>
  // const [html, setHtml] = useState<string>("");
  //
  // useEffect(() => {
  //   let cancelled = false;
  //
  //   const render = async () => {
  //     const result = await renderMarkdown(children, {
  //       foldCode,
  //       foldCodeMessage,
  //       inlineCodeHighlight,
  //     });
  //     if (!cancelled) setHtml(result);
  //   };
  //
  //   render();
  //
  //   return () => {
  //     cancelled = true;
  //   };
  // }, [children, foldCode, foldCodeMessage, inlineCodeHighlight]);
  //
  // return (
  //   <div
  //     className="markdown-body"
  //     // biome-ignore lint/security/noDangerouslySetInnerHtml: markdown is sanitized by pipeline
  //     dangerouslySetInnerHTML={{ __html: html }}
  //   />
  // );
}
