import type { i18n as I18n } from "i18next";
import { siteUrl } from "@/lib/site-url";

// The site name every page title carries, as legacy's `set_meta_tags site:`
// set it; the root route's fallback title is the bare name.
export const SITE_NAME = "CodeBasics";

interface SeoHead {
  i18n: I18n;
  title: string;
  description?: string;
  // The page's own pathname (a route match's `pathname`), made absolute for the
  // canonical link. Left out, the page has none.
  canonicalPath?: string;
  image?: string | null;
  // Legacy gave most public pages Open Graph and Twitter card tags; the few it
  // did not (cases, the sitemap) say `social: false`.
  social?: boolean;
  // The fuller Open Graph block the course pages set in legacy: the object's
  // type and locale, og:url as the canonical, and the image repeated as an
  // `image_src` link. With `social: false` it still names the title and the
  // image, which is what the lesson page emitted without the rest of the
  // social block.
  openGraph?: { type: "website" | "article"; locale: string };
}

// What a route's `head` returns for the tags below: its meta entries and links.
interface SeoHeadTags {
  links: { href: string; rel: string }[];
  meta: { content?: string; name?: string; property?: string; title?: string }[];
}

// The head legacy's meta-tags gem rendered from a controller's `seo_tags`:
// `CodeBasics | <title>`, the description, the og/twitter pair and the
// canonical link.
export function seoHead({
  i18n,
  title,
  description,
  canonicalPath,
  image,
  social = true,
  openGraph,
}: SeoHead): SeoHeadTags {
  // `/reviews/` and `/reviews` are one page, so the canonical drops a trailing
  // slash, as legacy's `*_url` helpers never emitted one.
  const canonical =
    canonicalPath === undefined
      ? undefined
      : siteUrl(canonicalPath.replace(/(?<last>.)\/+$/u, "$<last>"));
  const meta = [
    { title: `${SITE_NAME} | ${title}` },
    ...(description === undefined ? [] : [{ content: description, name: "description" }]),
    ...(social
      ? [
          { content: title, property: "og:title" },
          ...(description === undefined
            ? []
            : [{ content: description, property: "og:description" }]),
          ...(image ? [{ content: image, property: "og:image" }] : []),
          { content: "summary", name: "twitter:card" },
          { content: i18n.t(($) => $.links.hexlet_twitter_handle), name: "twitter:site" },
        ]
      : []),
    ...(openGraph
      ? [
          ...(social
            ? []
            : [
                { content: title, property: "og:title" },
                ...(image ? [{ content: image, property: "og:image" }] : []),
              ]),
          { content: openGraph.type, property: "og:type" },
          { content: openGraph.locale, property: "og:locale" },
          // Og:url is the canonical, so the two cannot drift apart.
          ...(canonical ? [{ content: canonical, property: "og:url" }] : []),
        ]
      : []),
  ];
  const links = [
    ...(canonical ? [{ href: canonical, rel: "canonical" }] : []),
    ...(openGraph && image ? [{ href: image, rel: "image_src" }] : []),
  ];
  return { links, meta };
}
