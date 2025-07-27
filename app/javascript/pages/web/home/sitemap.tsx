import { Accordion, Container, List, Stack, Text } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { Locale } from '@/types';
import type {
  LanguageCategory,
  LanguageSitemapLandingPage,
  LanguageSitemapLesson,
  SitemapBlogPost,
} from '@/types/serializers';

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

const getSuffix = (locale: Locale) => (locale === 'en' ? null : locale);

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
              <AppAnchor href={Routes.root_path({ suffix: getSuffix(locale) })}>
                {t('home.sitemap.home', { lng: locale })}
              </AppAnchor>
            </Text>

            <Accordion multiple variant="separated">
              {/* Курсы */}
              <Accordion.Item value="courses">
                <Accordion.Control>
                  {t('home.languages.courses', { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
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
                            <AppAnchor
                              href={Routes.language_path(landingPage.slug, {
                                suffix: getSuffix(locale),
                              })}
                            >
                              {landingPage.header}
                            </AppAnchor>
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
                  {t('blog_posts.index.header', { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {blogPostsByLocale[locale].map((post) => (
                      <List.Item key={post.id}>
                        <AppAnchor
                          href={Routes.blog_post_path(post.slug, {
                            suffix: getSuffix(post.locale),
                          })}
                        >
                          {post.name}
                        </AppAnchor>
                      </List.Item>
                    ))}
                  </List>
                </Accordion.Panel>
              </Accordion.Item>

              {/* Категории */}
              <Accordion.Item value="categories">
                <Accordion.Control>
                  {t('language_categories.index.header', { lng: locale })}
                </Accordion.Control>
                <Accordion.Panel>
                  <List listStyleType="none" spacing="xs" pl="md">
                    {categoriesByLocale[locale].map((category) => (
                      <List.Item key={category.id}>
                        <AppAnchor
                          href={Routes.language_category_path(category.slug!, {
                            suffix: getSuffix(category.locale),
                          })}
                        >
                          {category.name}
                        </AppAnchor>
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
