import { Container, Grid, Title, Typography, type TypographyProps } from "@mantine/core";
import { createFileRoute, notFound } from "@tanstack/react-router";
import type { TFunction } from "i18next";
import type { ComponentType } from "react";
import { useTranslation } from "react-i18next";
import AboutEn from "@/components/pages/AboutEn";
import AboutRu from "@/components/pages/AboutRu";
import AuthorsEn from "@/components/pages/AuthorsEn";
import AuthorsRu from "@/components/pages/AuthorsRu";
import CookiePolicyEn from "@/components/pages/CookiePolicyEn";
import CookiePolicyRu from "@/components/pages/CookiePolicyRu";
import PrivacyEn from "@/components/pages/PrivacyEn";
import PrivacyRu from "@/components/pages/PrivacyRu";
import TosEn from "@/components/pages/TosEn";
import TosRu from "@/components/pages/TosRu";
import { siteUrl } from "@/lib/site-url";

// The static pages, at their legacy URLs (`resources :pages, only: :show`),
// ported from legacy pages/show. Their bodies are hardcoded markup with no data
// (ADR-0015), written in ru and en only.
const PAGES = {
  ru: {
    about: AboutRu,
    authors: AuthorsRu,
    cookie_policy: CookiePolicyRu,
    privacy: PrivacyRu,
    tos: TosRu,
  },
  en: {
    about: AboutEn,
    authors: AuthorsEn,
    cookie_policy: CookiePolicyEn,
    privacy: PrivacyEn,
    tos: TosEn,
  },
} satisfies Record<string, Record<string, ComponentType>>;

type PageId = keyof typeof PAGES.en;

function isPageId(id: string): id is PageId {
  return Object.hasOwn(PAGES.en, id);
}

// The legal pages are not meant to rank, so they carry no canonical link. Legacy
// meant the same (its DISALLOWED_PAGES) but set the canonical unconditionally
// first, so every page emitted one; this follows the intent.
const WITHOUT_CANONICAL: readonly PageId[] = ["cookie_policy", "privacy", "tos"];

export const Route = createFileRoute("/{-$locale}/pages/$id")({
  // An unknown page is a 404, as legacy's RoutingError was.
  beforeLoad: ({ params }) => {
    if (!isPageId(params.id)) throw notFound();
  },
  // Legacy meta: the page title and description, Open Graph and a summary
  // Twitter card. es has no page strings of its own in legacy, which fell back
  // to en, and the es catalog carries those en strings.
  head: ({ match }) => {
    const { id } = match.params;
    if (!isPageId(id)) return {};
    const { t } = match.context.i18n;
    const { title, description } = pageMeta(t, id);
    return {
      meta: [
        { title: `CodeBasics | ${title}` },
        { name: "description", content: description },
        { property: "og:title", content: title },
        { property: "og:description", content: description },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: t(($) => $.links.hexlet_twitter_handle) },
      ],
      links: WITHOUT_CANONICAL.includes(id)
        ? []
        : [{ rel: "canonical", href: siteUrl(match.pathname) }],
    };
  },
  component: Show,
});

// Typed selectors cannot take a runtime key, so each page names its own strings.
function pageMeta(t: TFunction, id: PageId): { title: string; description: string } {
  switch (id) {
    case "about":
      return {
        title: t(($) => $.pages.parts.about.title),
        description: t(($) => $.pages.parts.about.meta.description),
      };
    case "authors":
      return {
        title: t(($) => $.pages.parts.authors.title),
        description: t(($) => $.pages.parts.authors.meta.description),
      };
    case "cookie_policy":
      return {
        title: t(($) => $.pages.parts.cookie_policy.title),
        description: t(($) => $.pages.parts.cookie_policy.meta.description),
      };
    case "privacy":
      return {
        title: t(($) => $.pages.parts.privacy.title),
        description: t(($) => $.pages.parts.privacy.meta.description),
      };
    case "tos":
      return {
        title: t(($) => $.pages.parts.tos.title),
        description: t(($) => $.pages.parts.tos.meta.description),
      };
  }
}

// Legacy typographyStyles: body text at the large line height, every heading
// below h1 brought down to the small sizes.
const typographyStyles: TypographyProps["styles"] = (theme) => ({
  root: {
    overflowWrap: "break-word",
    lineHeight: theme.lineHeights.lg,
    h2: theme.headings.sizes.h5,
    h3: theme.headings.sizes.h6,
    h4: theme.headings.sizes.h6,
    h5: theme.headings.sizes.h6,
    h6: theme.headings.sizes.h6,
  },
});

function Show() {
  const { t, i18n } = useTranslation();
  const { id } = Route.useParams();
  if (!isPageId(id)) return null;
  // Only ru and en bodies exist; es reads the ru one, as legacy did.
  const Body = PAGES[i18n.language === "en" ? "en" : "ru"][id];

  return (
    <>
      {/* Legacy ApplicationLayout's `header` block, centered. */}
      <Container size="lg" my="xl">
        <Title order={1} ta="center">
          {pageMeta(t, id).title}
        </Title>
      </Container>
      <Container size="xl">
        <Grid justify="center" mb="xl">
          <Grid.Col span={{ base: 12, md: 10, lg: 8 }}>
            <Typography styles={typographyStyles}>
              <Body />
            </Typography>
          </Grid.Col>
        </Grid>
      </Container>
    </>
  );
}
