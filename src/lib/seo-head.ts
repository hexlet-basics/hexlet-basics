import type { i18n as I18n } from "i18next";
import { siteUrl } from "@/lib/site-url";

// The site name every page title carries, as legacy's `set_meta_tags site:`
// set it; the root route's fallback title is the bare name.
export const SITE_NAME = "CodeBasics";

type SeoHead = {
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
};

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
}: SeoHead) {
  // `/reviews/` and `/reviews` are one page, so the canonical drops a trailing
  // slash, as legacy's `*_url` helpers never emitted one.
  const canonical =
    canonicalPath === undefined ? undefined : siteUrl(canonicalPath.replace(/(.)\/+$/, "$1"));
  const meta = [
    { title: `${SITE_NAME} | ${title}` },
    ...(description === undefined ? [] : [{ name: "description", content: description }]),
    ...(social
      ? [
          { property: "og:title", content: title },
          ...(description === undefined
            ? []
            : [{ property: "og:description", content: description }]),
          ...(image ? [{ property: "og:image", content: image }] : []),
          { name: "twitter:card", content: "summary" },
          { name: "twitter:site", content: i18n.t(($) => $.links.hexlet_twitter_handle) },
        ]
      : []),
    ...(openGraph
      ? [
          ...(social
            ? []
            : [
                { property: "og:title", content: title },
                ...(image ? [{ property: "og:image", content: image }] : []),
              ]),
          { property: "og:type", content: openGraph.type },
          { property: "og:locale", content: openGraph.locale },
          // og:url is the canonical, so the two cannot drift apart.
          ...(canonical ? [{ property: "og:url", content: canonical }] : []),
        ]
      : []),
  ];
  const links = [
    ...(canonical ? [{ rel: "canonical", href: canonical }] : []),
    ...(openGraph && image ? [{ rel: "image_src", href: image }] : []),
  ];
  return { meta, links };
}
