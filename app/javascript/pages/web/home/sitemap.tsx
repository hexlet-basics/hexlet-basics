import { Link } from "@inertiajs/react";
import { Accordion, Container, List, Stack, Text } from "@mantine/core";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { Locale } from "@/types";
import type {
  LanguageCategory,
  LanguageSitemapLandingPage,
  LanguageSitemapLesson,
  SitemapBlogPost,
} from "@/types/serializers";

type Props = PropsWithChildren & {
  title: string;
  orderedLocales: Locale[];
  landingPagesByLocale: Record<Locale, LanguageSitemapLandingPage[]>;
  lessonsByLocaleAndLanguageId: Record<
    Locale,
    Record<number, LanguageSitemapLesson[]>
  >;
  blogPostsByLocale: Record<Locale, SitemapBlogPost[]>;
  categoriesByLocale: Record<Locale, LanguageCategory[]>;
};

const getSuffix = (locale: Locale) => (locale === "en" ? null : locale);

export default function SiteMap({
  title,
  orderedLocales,
  landingPagesByLocale,
  lessonsByLocaleAndLanguageId,
  blogPostsByLocale,
  categoriesByLocale,
}: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout header={title}>
      {orderedLocales.map((locale) => (
        <Container key={locale} mb="xl">
          <Stack gap="xs">
            <Text fw={500} size="lg">
              <Text
                component={Link}
                href={Routes.root_path({ suffix: getSuffix(locale) })}
                span
              >
                {t(($) => $.home.sitemap.home, {
                  lng: locale,
                })}
              </Text>
            </Text>

            <Accordion multiple variant="separated">
              {/* Курсы */}
              <Accordion.Item value="courses">
                <Accordion.Control>
                  {t(($) => $.home.languages.courses, {
                    lng: locale,
                  })}
                </Accordion.Control>
                <Accordion.Panel>
                  <Text fw={500} size="sm" mb="sm">
                    <Text
                      component={Link}
                      href={Routes.languages_path({
                        suffix: getSuffix(locale),
                      })}
                      span
                    >
                      {t(($) => $.pages.languages.index.header, {
                        lng: locale,
                      })}
                    </Text>
                  </Text>
                  {landingPagesByLocale[locale]
                    .filter(
                      (landingPage) =>
                        lessonsByLocaleAndLanguageId[landingPage.locale]?.[
                          landingPage.language_id
                        ]?.length,
                    )
                    .map((landingPage) => {
                      const lessons =
                        lessonsByLocaleAndLanguageId[landingPage.locale]?.[
                          landingPage.language_id
                        ] ?? [];

                      return (
                        <div key={landingPage.id}>
                          <Text fw={500} size="sm" mt="sm">
                            <Text
                              component={Link}
                              href={Routes.language_path(landingPage.slug, {
                                suffix: getSuffix(locale),
                              })}
                              span
                            >
                              {landingPage.header}
                            </Text>
                          </Text>
                          {/* <List listStyleType="none" spacing="xs" pl="md"> */}
                          {/*   {lessons.map((lesson) => ( */}
                          {/*     <List.Item key={lesson.id}> */}
                          {/*       <AppAnchor */}
                          {/*         href={Routes.language_lesson_path( */}
                          {/*           landingPage.slug, */}
                          {/*           lesson.slug, */}
                          {/*           { suffix: getSuffix(locale) }, */}
                          {/*         )} */}
                          {/*       > */}
                          {/*         {lesson.natural_order}. {lesson.name} */}
                          {/*       </AppAnchor> */}
                          {/*     </List.Item> */}
                          {/*   ))} */}
                          {/* </List> */}
                        </div>
                      );
                    })}
                </Accordion.Panel>
              </Accordion.Item>

              {/* Блог */}
              <Accordion.Item value="blog">
                <Accordion.Control>
                  {t(($) => $.blog_posts.index.header, {
                    lng: locale,
                  })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {blogPostsByLocale[locale].map((post) => (
                      <List.Item key={post.id}>
                        <Text
                          component={Link}
                          href={Routes.blog_post_path(post.slug, {
                            suffix: getSuffix(post.locale),
                          })}
                          span
                        >
                          {post.name}
                        </Text>
                      </List.Item>
                    ))}
                  </List>
                </Accordion.Panel>
              </Accordion.Item>

              {/* Категории */}
              <Accordion.Item value="categories">
                <Accordion.Control>
                  {t(($) => $.language_categories.index.header, {
                    lng: locale,
                  })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {categoriesByLocale[locale].map((category) => (
                      <List.Item key={category.id}>
                        <Text
                          component={Link}
                          href={Routes.language_category_path(category.slug!, {
                            suffix: getSuffix(category.locale),
                          })}
                          span
                        >
                          {category.name}
                        </Text>
                      </List.Item>
                    ))}
                  </List>
                </Accordion.Panel>
              </Accordion.Item>
            </Accordion>
          </Stack>
        </Container>
      ))}
    </ApplicationLayout>
  );
}
