import { Box, Container, Grid, TableOfContents, Text, Title, Typography } from "@mantine/core";
import type { TypographyProps } from "@mantine/core";
import { createFileRoute, notFound } from "@tanstack/react-router";
import type { TFunction } from "i18next";
import type { ComponentType, ReactElement } from "react";
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
import { seoHead } from "@/lib/seo-head";

// The static pages, at their legacy URLs (`resources :pages, only: :show`),
// ported from legacy pages/show. Their bodies are hardcoded markup with no data
// (ADR-0015), written in ru and en only.
const PAGES = {
  en: {
    about: AboutEn,
    authors: AuthorsEn,
    cookie_policy: CookiePolicyEn,
    privacy: PrivacyEn,
    tos: TosEn,
  },
  ru: {
    about: AboutRu,
    authors: AuthorsRu,
    cookie_policy: CookiePolicyRu,
    privacy: PrivacyRu,
    tos: TosRu,
  },
} satisfies Record<string, Record<string, ComponentType>>;

type PageId = keyof typeof PAGES.en;

function isPageId(id: string): id is PageId {
  return Object.hasOwn(PAGES.en, id);
}

// The legal pages are not meant to rank, so they carry no canonical link. Legacy
// meant the same (its DISALLOWED_PAGES) but set the canonical unconditionally
// first, so every page emitted one; this follows the intent.
const WITHOUT_CANONICAL: ReadonlySet<PageId> = new Set<PageId>(["cookie_policy", "privacy", "tos"]);

export const Route = createFileRoute("/{-$locale}/pages/$id")({
  // An unknown page is a 404, as legacy's RoutingError was.
  beforeLoad: ({ params }) => {
    if (!isPageId(params.id)) {
      throw notFound();
    }
  },
  component: Show,
  // Legacy meta: the page title and description, the social tags and, but for
  // the legal pages, the canonical link. es has no page strings of its own in
  // legacy, which fell back to en, and the es catalog carries those en strings.
  head: ({ match }) => {
    const { id } = match.params;
    if (!isPageId(id)) {
      return {};
    }
    const { i18n } = match.context;
    const { title, description } = pageMeta(i18n.t, id);
    return seoHead({
      canonicalPath: WITHOUT_CANONICAL.has(id) ? undefined : match.pathname,
      description,
      i18n,
      title,
    });
  },
});

// Typed selectors cannot take a runtime key, so each page names its own strings.
function pageMeta(t: TFunction, id: PageId): { title: string; description: string } {
  switch (id) {
    case "about": {
      return {
        description: t(($) => $.pages.parts.about.meta.description),
        title: t(($) => $.pages.parts.about.title),
      };
    }
    case "authors": {
      return {
        description: t(($) => $.pages.parts.authors.meta.description),
        title: t(($) => $.pages.parts.authors.title),
      };
    }
    case "cookie_policy": {
      return {
        description: t(($) => $.pages.parts.cookie_policy.meta.description),
        title: t(($) => $.pages.parts.cookie_policy.title),
      };
    }
    case "privacy": {
      return {
        description: t(($) => $.pages.parts.privacy.meta.description),
        title: t(($) => $.pages.parts.privacy.title),
      };
    }
    case "tos": {
      return {
        description: t(($) => $.pages.parts.tos.meta.description),
        title: t(($) => $.pages.parts.tos.title),
      };
    }
    default: {
      throw new Error(`No meta for the page ${String(id satisfies never)}`);
    }
  }
}

// Legacy typographyStyles: body text at the large line height, every heading
// below h1 brought down to the small sizes.
const typographyStyles: TypographyProps["styles"] = (theme) => ({
  root: {
    h2: theme.headings.sizes.h5,
    h3: theme.headings.sizes.h6,
    h4: theme.headings.sizes.h6,
    h5: theme.headings.sizes.h6,
    h6: theme.headings.sizes.h6,
    lineHeight: theme.lineHeights.lg,
    overflowWrap: "break-word",
  },
});

function Show(): ReactElement | null {
  const { t, i18n } = useTranslation();
  const { id } = Route.useParams();
  if (!isPageId(id)) {
    return null;
  }
  // Only ru and en bodies exist; es reads the ru one, as legacy did.
  const Body = PAGES[i18n.language === "en" ? "en" : "ru"][id];

  return (
    <Container size="lg" my="xl">
      <Title order={1} mb="xl">
        {pageMeta(t, id).title}
      </Title>
      {/* Keyed by page, so the contents are read again when the page changes. */}
      <Grid key={id} gap="xl">
        <Grid.Col span={{ base: 12, lg: 9 }}>
          <Typography id="page-content" styles={typographyStyles}>
            <Body />
          </Typography>
        </Grid.Col>
        <Grid.Col span={3} visibleFrom="lg">
          {/* The headings the page itself has: the legal texts use h3/h4 and the
              rest h2, so the first level found is the top of the list. */}
          <Box pos="sticky" top={80}>
            <Text fw="bold" mb="xs">
              {t(($) => $.common.tos)}
            </Text>
            <TableOfContents
              variant="light"
              size="sm"
              minDepthToOffset={2}
              scrollSpyOptions={{ selector: "#page-content :is(h2, h3)" }}
              getControlProps={({ data }) => ({
                children: data.value,
                onClick: () => {
                  data.getNode().scrollIntoView({ behavior: "smooth" });
                },
              })}
            />
          </Box>
        </Grid.Col>
      </Grid>
    </Container>
  );
}
