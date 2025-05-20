import type { Locale, SharedProps } from "@/types";
import type { PropsWithChildren } from "react";
import React from "react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type {
  LanguageCategory,
  LanguageSitemapLandingPage,
  LanguageSitemapLesson,
  SitemapBlogPost,
} from "@/types/serializers";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Collapse, Container } from "react-bootstrap";
import { Link } from "@inertiajs/react";

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
  lessonsByLocaleAndLanguageId: Props["lessonsByLocaleAndLanguageId"];
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

const getSuffix = (locale: Locale) => (locale === "en" ? null : locale);

function LessonsBlock({ lessons, landingPage }: LessonsBlockProps) {
  return (
    <ul className="m-0 ms-3 list-unstyled pb-2">
      {lessons.map((lesson) => (
        <li key={`${landingPage.id}-${lesson.id}`}>
          <Link
            className="text-decoration-none"
            href={Routes.language_lesson_path(landingPage.slug, lesson.slug, {
              suffix: getSuffix(lesson.locale),
            })}
          >
            <span className="me-1">{lesson.natural_order}.</span>
            {lesson.name}
          </Link>
        </li>
      ))}
    </ul>
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
    <div className="ms-3">
      {filteredLandingPages.map((landingPage, index) => {
        const landingPageOpened = index === activeLandingPageIdx;

        return (
          <div key={landingPage.id}>
            <Link
              className="text-decoration-none text-body"
              href={Routes.language_path(landingPage.slug, {
                suffix: getSuffix(landingPage.locale),
              })}
            >
              {landingPage.header}
            </Link>
            <button
              type="button"
              onClick={() =>
                setActiveLandingPageIdx(landingPageOpened ? null : index)
              }
              aria-controls={String(index)}
              className="btn btn-link py-0 shadow-none"
            >
              <span
                className={`bi ${landingPageOpened ? "bi-chevron-up" : "bi-chevron-down"}`}
              />
            </button>
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
    </div>
  );
}

function BlogPostsBlock({ posts, opened }: BlogPostsBlockProps) {
  return (
    <Collapse in={opened}>
      <ul className="m-0 ms-3 list-unstyled" id="blog-collapse">
        {posts.map((post) => (
          <li key={post.id}>
            <Link
              className="text-decoration-none"
              href={Routes.blog_post_path(post.slug, {
                suffix: getSuffix(post.locale),
              })}
            >
              {post.name}
            </Link>
          </li>
        ))}
      </ul>
    </Collapse>
  );
}

function LanguageCategoriesBlock({
  categories,
  opened,
}: LanguageCategoriesBlockProps) {
  return (
    <Collapse in={opened}>
      <ul className="m-0 ms-3 list-unstyled" id="blog-collapse">
        {categories.map((category) => (
          <li key={category.id}>
            <Link
              className="text-decoration-none"
              href={Routes.language_category_path(category.slug!, {
                suffix: getSuffix(category.locale),
              })}
            >
              {category.name}
            </Link>
          </li>
        ))}
      </ul>
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
          <Container key={locale} className="mb-5">
            <h2 className="h5 mb-3">
              <Link
                id={`main-${locale}`}
                className="text-decoration-none link-body-emphasis"
                href={Routes.root_path({ suffix: getSuffix(locale) })}
              >
                {t("home.sitemap.home", { lng: locale })}
              </Link>
            </h2>
            <h2 className="h5 mb-3">
              <Link
                id={`courses-${locale}`}
                className="text-decoration-none link-body-emphasis"
                href={`#courses-${locale}`}
              >
                {t("home.languages.courses", { lng: locale })}
              </Link>
            </h2>
            <LadingPagesBlock
              locale={locale}
              landingPages={landingPagesByLocale[locale]}
              lessonsByLocaleAndLanguageId={lessonsByLocaleAndLanguageId}
            />

            <h2 className="h5 my-2">
              <Link
                id={`blog-${locale}`}
                className="text-decoration-none link-body-emphasis"
                href={`#blog-${locale}`}
              >
                {t("blog_posts.index.header", { lng: locale })}
              </Link>
              <button
                type="button"
                onClick={() => setBlogOpened(!blogOpened)}
                aria-controls="blog-collapse"
                className="btn btn-link py-0 shadow-none"
              >
                <span
                  className={`bi ${blogOpened ? "bi-chevron-up" : "bi-chevron-down"}`}
                />
              </button>
            </h2>
            <BlogPostsBlock
              posts={blogPostsByLocale[locale]}
              opened={blogOpened}
            />
            <h2 className="h5 py-2">
              <Link
                id={`categories-${locale}`}
                className="text-decoration-none link-body-emphasis"
                href={`#categories-${locale}`}
              >
                {t("language_categories.index.header", { lng: locale })}
              </Link>
              <button
                type="button"
                onClick={() => setCategoriesOpened(!categoriesOpened)}
                aria-controls="categories-collapse"
                className="btn btn-link py-0 shadow-none"
              >
                <span
                  className={`bi ${categoriesOpened ? "bi-chevron-up" : "bi-chevron-down"}`}
                />
              </button>
            </h2>
            <LanguageCategoriesBlock
              categories={categoriesByLocale[locale]}
              opened={categoriesOpened}
            />
          </Container>
        );
      })}
    </ApplicationLayout>
  );
}
