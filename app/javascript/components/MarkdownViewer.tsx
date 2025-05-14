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
}
