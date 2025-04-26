import rehypeShiki from "@shikijs/rehype";
import type { Root } from "hast";
import rehypeExternalLinks from "rehype-external-links";
import rehypeRaw from "rehype-raw";
import rehypeStringify from "rehype-stringify";
import remarkGfm from "remark-gfm";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import { unified } from "unified";
import type { Plugin } from "unified";
import { visit } from "unist-util-visit";

export interface RenderMarkdownOptions {
  foldCode?: boolean;
  foldCodeMessage?: string;
  inlineCodeHighlight?: boolean;
}

/**
 * Plugin: Wrap <pre><code> blocks in <details>
 */
function rehypeFoldCode({
  enabled = false,
  message = "Show code",
}: {
  enabled: boolean;
  message: string;
}): Plugin<[], Root> {
  return () => (tree: Root) => {
    if (!enabled) return;

    visit(tree, "element", (node, index, parent) => {
      if (
        node.tagName === "pre" &&
        node.children?.[0]?.type === "element" &&
        node.children[0].tagName === "code" &&
        parent?.type === "element" &&
        typeof index === "number"
      ) {
        const innerHtml = unified()
          .use(rehypeStringify)
          .stringify({ type: "root", children: [node] });

        parent.children[index] = {
          type: "raw",
          value: `
<details><summary tabindex="-1">${message}</summary>
  <div class="m-2">${innerHtml}</div>
</details>`.trim(),
        };
      }
    });
  };
}

/**
 * Render Markdown to HTML using @shikijs/rehype for highlighting
 */
export default async function renderMarkdown(
  markdown: string,
  options: RenderMarkdownOptions = {},
): Promise<string> {
  const {
    foldCode = false,
    foldCodeMessage = "Show code",
    inlineCodeHighlight = true,
  } = options;

  const processor = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeRaw)
    .use(rehypeShiki, {
      themes: {
        light: "vitesse-light",
        dark: "vitesse-dark",
      },
      inline: inlineCodeHighlight ? "tailing-curly-colon" : false,
    })
    .use(
      rehypeFoldCode({
        enabled: foldCode,
        message: foldCodeMessage,
      }),
    )
    .use(rehypeExternalLinks, {
      target: "_blank",
      rel: ["noopener", "noreferrer"],
    })
    .use(rehypeStringify);

  const file = await processor.process(markdown);
  return String(file);
}
