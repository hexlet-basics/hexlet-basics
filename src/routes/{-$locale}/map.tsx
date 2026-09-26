import { Accordion, Container, List, Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { getSitemapOptions } from "@/client/@tanstack/react-query.gen";
import { TextLink } from "@/components/RouterLink";
import type { Locale } from "@/lib/i18n";
import { seoHead } from "@/lib/seo-head";

// The HTML sitemap at its legacy URL, ported from legacy home/sitemap. Legacy
// served it on the ru site only (`/ru/map`); `/map` and `/es/map` are 404. It
// lists the ru site first, then the en one, each in its own language; es is
// left out because its links did not work on the legacy site.
const SECTION_LOCALES: Locale[] = ["ru", "en"];

// A link's `{-$locale}` param: en is the unprefixed default.
const localeParam = (locale: Locale) => (locale === "en" ? undefined : locale);

export const Route = createFileRoute("/{-$locale}/map")({
  beforeLoad: ({ params }) => {
    if (params.locale !== "ru") {
      throw notFound();
    }
  },
  loader: ({ context }) => context.queryClient.ensureQueryData(getSitemapOptions()),
  // Legacy set the title alone: no description, canonical or social tags.
  head: ({ match }) => {
    const { i18n } = match.context;
    return seoHead({ i18n, title: i18n.t(($) => $.home.sitemap.title), social: false });
  },
  component: SitemapPage,
});

function SitemapPage() {
  const { t } = useTranslation();
  const { data } = useQuery(getSitemapOptions());
  const landingPages = data?.landingPages ?? [];
  const blogPosts = data?.blogPosts ?? [];
  const categories = data?.categories ?? [];

  return (
    <>
      <Container size="lg" my="xl">
        <Title order={1}>{t(($) => $.home.sitemap.title)}</Title>
      </Container>

      {SECTION_LOCALES.map((locale) => (
        <Container key={locale} mb="xl">
          <Stack gap="xs">
            <Text fw={500} size="lg">
              <TextLink to="/{-$locale}" params={{ locale: localeParam(locale) }} span>
                {t(($) => $.home.sitemap.home, { lng: locale })}
              </TextLink>
            </Text>

            <Accordion multiple variant="separated">
              <Accordion.Item value="courses">
                <Accordion.Control>
                  {t(($) => $.home.courses.courses, { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
                  <Text fw={500} size="sm" mb="sm">
                    <TextLink
                      to="/{-$locale}/languages"
                      params={{ locale: localeParam(locale) }}
                      span
                    >
                      {t(($) => $.pages.courses.index.header, { lng: locale })}
                    </TextLink>
                  </Text>
                  {landingPages
                    .filter((landingPage) => landingPage.locale === locale)
                    .map((landingPage) => (
                      <Text key={landingPage.id} fw={500} size="sm" mt="sm">
                        <TextLink
                          to="/{-$locale}/languages/$slug"
                          params={{ locale: localeParam(locale), slug: landingPage.slug }}
                          span
                        >
                          {landingPage.header}
                        </TextLink>
                      </Text>
                    ))}
                </Accordion.Panel>
              </Accordion.Item>

              <Accordion.Item value="blog">
                <Accordion.Control>
                  {t(($) => $.blog_posts.index.header, { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {blogPosts
                      .filter((post) => post.locale === locale)
                      .map((post) => (
                        <List.Item key={post.id}>
                          <TextLink
                            to="/{-$locale}/blog_posts/$slug"
                            params={{ locale: localeParam(locale), slug: post.slug }}
                            span
                          >
                            {post.name}
                          </TextLink>
                        </List.Item>
                      ))}
                  </List>
                </Accordion.Panel>
              </Accordion.Item>

              <Accordion.Item value="categories">
                <Accordion.Control>
                  {t(($) => $.course_categories.index.header, { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {categories
                      .filter((category) => category.locale === locale)
                      .map((category) => (
                        <List.Item key={category.id}>
                          <TextLink
                            to="/{-$locale}/language_categories/$slug"
                            params={{ locale: localeParam(locale), slug: category.slug ?? "" }}
                            span
                          >
                            {category.name}
                          </TextLink>
                        </List.Item>
                      ))}
                  </List>
                </Accordion.Panel>
              </Accordion.Item>
            </Accordion>
          </Stack>
        </Container>
      ))}
    </>
  );
}
