import { Anchor, Box } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import type { ComponentPropsWithoutRef } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { getActiveBannerOptions } from "@/client/@tanstack/react-query.gen";
import type { BannerBackground } from "@/client/types.gen";

// Background presets; the keys are the contract's BannerBackground. The colours
// are Mantine's own variables, which resolve per colour scheme, so the strip
// needs no scheme check in JS.
const backgrounds: Record<BannerBackground, string> = {
  blue: "var(--mantine-color-blue-6)",
  cta_gradient:
    "linear-gradient(90deg, var(--mantine-color-blue-filled), var(--mantine-color-cyan-5))",
  dark: "var(--mantine-color-dark-6)",
};

// Inline markdown: the block wrapper and links are dropped, because the link is
// the banner's own `url`, set on the whole strip.
const inlineComponents = {
  a: ({ children }: ComponentPropsWithoutRef<"a">) => children,
  p: ({ children }: ComponentPropsWithoutRef<"p">) => children,
};

// The site banner, ported from legacy Banner: the live banner of the request
// locale above the page, or nothing when none is live. It is the pages that
// carry one (the home page and the course page) that render it, as legacy's
// `showBanner` did.
export function Banner() {
  const { data } = useQuery(getActiveBannerOptions());
  const banner = data?.banner;
  if (!banner) {
    return null;
  }

  const content = (
    <Markdown
      skipHtml
      remarkPlugins={[remarkGfm]}
      disallowedElements={["a"]}
      unwrapDisallowed
      components={inlineComponents}
    >
      {banner.body}
    </Markdown>
  );

  return (
    <Box bg={backgrounds[banner.background]} c="white" p="xs" ta="center">
      {banner.url ? (
        <Anchor href={banner.url} c="inherit" target="_blank" rel="noopener noreferrer">
          {content}
        </Anchor>
      ) : (
        content
      )}
    </Box>
  );
}
