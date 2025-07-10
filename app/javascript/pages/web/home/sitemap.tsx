import { Link } from '@inertiajs/react';
import {
  Button,
  Collapse,
  Container,
  Group,
  List,
  Stack,
  Text,
} from '@mantine/core';
import type { PropsWithChildren } from 'react';
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { Locale, SharedProps } from '@/types';
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

type LandingPagesBlockProps = {
  locale: Locale;
  landingPages: LanguageSitemapLandingPage[];
  lessonsByLocaleAndLanguageId: Props['lessonsByLocaleAndLanguageId'];
};

type LessonsBlockProps = {
  lessons: LanguageSitemapLesson[];
  landingPage: LanguageSitemapLandingPage;
};

type BlogPostsBlockProps = {
  posts: SitemapBlogPost[];
  opened: boolean;
};

type LanguageCategoriesBlockProps = {
  categories: LanguageCategory[];
  opened: boolean;
};

const getSuffix = (locale: Locale) => (locale === 'en' ? null : locale);

function LessonsBlock({ lessons, landingPage }: LessonsBlockProps) {
  return (
    <List listStyleType="none" pl="md" pb="sm">
      {lessons.map((lesson) => (
        <List.Item key={`${landingPage.id}-${lesson.id}`}>
          <Link
            style={{ textDecoration: 'none' }}
            href={Routes.language_lesson_path(landingPage.slug, lesson.slug, {
              suffix: getSuffix(lesson.locale),
            })}
          >
            <Group gap="xs">
              <Text span>{lesson.natural_order}.</Text>
              <Text span>{lesson.name}</Text>
            </Group>
          </Link>
        </List.Item>
      ))}
    </List>
  );
}

function LadingPagesBlock({
  landingPages,
  lessonsByLocaleAndLanguageId,
}: LandingPagesBlockProps) {
  const filteredLandingPages = landingPages.filter(
    (landingPage) =>
      lessonsByLocaleAndLanguageId[landingPage.locale][landingPage.language_id],
  );

  const [activeLandingPageIdx, setActiveLandingPageIdx] = useState<
    number | null
  >(null);

  return (
    <Stack gap="xs" pl="md">
      {filteredLandingPages.map((landingPage, index) => {
        const landingPageOpened = index === activeLandingPageIdx;

        return (
          <div key={landingPage.id}>
            <Group gap="xs">
              <Link
                style={{ textDecoration: 'none' }}
                href={Routes.language_path(landingPage.slug, {
                  suffix: getSuffix(landingPage.locale),
                })}
              >
                {landingPage.header}
              </Link>
              <Button
                variant="subtle"
                size="xs"
                onClick={() =>
                  setActiveLandingPageIdx(landingPageOpened ? null : index)
                }
                aria-controls={String(index)}
              >
                <span
                  className={`bi ${landingPageOpened ? 'bi-chevron-up' : 'bi-chevron-down'}`}
                />
              </Button>
            </Group>
            <Collapse in={landingPageOpened}>
              <div id={String(index)}>
                <LessonsBlock
                  landingPage={landingPage}
                  lessons={
                    lessonsByLocaleAndLanguageId[landingPage.locale][
                      landingPage.language_id
                    ]
                  }
                />
              </div>
            </Collapse>
          </div>
        );
      })}
    </Stack>
  );
}

function BlogPostsBlock({ posts, opened }: BlogPostsBlockProps) {
  return (
    <Collapse in={opened}>
      <List listStyleType="none" pl="md" id="blog-collapse">
        {posts.map((post) => (
          <List.Item key={post.id}>
            <Link
              style={{ textDecoration: 'none' }}
              href={Routes.blog_post_path(post.slug, {
                suffix: getSuffix(post.locale),
              })}
            >
              {post.name}
            </Link>
          </List.Item>
        ))}
      </List>
    </Collapse>
  );
}

function LanguageCategoriesBlock({
  categories,
  opened,
}: LanguageCategoriesBlockProps) {
  return (
    <Collapse in={opened}>
      <List listStyleType="none" pl="md" id="blog-collapse">
        {categories.map((category) => (
          <List.Item key={category.id}>
            <Link
              style={{ textDecoration: 'none' }}
              href={Routes.language_category_path(category.slug!, {
                suffix: getSuffix(category.locale),
              })}
            >
              {category.name}
            </Link>
          </List.Item>
        ))}
      </List>
    </Collapse>
  );
}

// NOTE для карты сайта выводим ссылки для всех локалей
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
      {orderedLocales.map((locale) => {
        const [blogOpened, setBlogOpened] = useState(false);
        const [categoriesOpened, setCategoriesOpened] = useState(false);

        return (
          <Container key={locale} mb="xl">
            <Stack gap="xs">
              <Text fw={500} size="lg">
                <Link
                  id={`main-${locale}`}
                  style={{ textDecoration: 'none' }}
                  href={Routes.root_path({ suffix: getSuffix(locale) })}
                >
                  {t('home.sitemap.home', { lng: locale })}
                </Link>
              </Text>
              <Text fw={500} size="lg">
                <Link
                  id={`courses-${locale}`}
                  style={{ textDecoration: 'none' }}
                  href={`#courses-${locale}`}
                >
                  {t('home.languages.courses', { lng: locale })}
                </Link>
              </Text>
              <LadingPagesBlock
                locale={locale}
                landingPages={landingPagesByLocale[locale]}
                lessonsByLocaleAndLanguageId={lessonsByLocaleAndLanguageId}
              />

              <Group gap="xs">
                <Text fw={500} size="lg">
                  <Link
                    id={`blog-${locale}`}
                    style={{ textDecoration: 'none' }}
                    href={`#blog-${locale}`}
                  >
                    {t('blog_posts.index.header', { lng: locale })}
                  </Link>
                </Text>
                <Button
                  variant="subtle"
                  size="xs"
                  onClick={() => setBlogOpened(!blogOpened)}
                  aria-controls="blog-collapse"
                >
                  <span
                    className={`bi ${blogOpened ? 'bi-chevron-up' : 'bi-chevron-down'}`}
                  />
                </Button>
              </Group>
              <BlogPostsBlock
                posts={blogPostsByLocale[locale]}
                opened={blogOpened}
              />
              <Group gap="xs">
                <Text fw={500} size="lg">
                  <Link
                    id={`categories-${locale}`}
                    style={{ textDecoration: 'none' }}
                    href={`#categories-${locale}`}
                  >
                    {t('language_categories.index.header', { lng: locale })}
                  </Link>
                </Text>
                <Button
                  variant="subtle"
                  size="xs"
                  onClick={() => setCategoriesOpened(!categoriesOpened)}
                  aria-controls="categories-collapse"
                >
                  <span
                    className={`bi ${categoriesOpened ? 'bi-chevron-up' : 'bi-chevron-down'}`}
                  />
                </Button>
              </Group>
              <LanguageCategoriesBlock
                categories={categoriesByLocale[locale]}
                opened={categoriesOpened}
              />
            </Stack>
          </Container>
        );
      })}
    </ApplicationLayout>
  );
}
