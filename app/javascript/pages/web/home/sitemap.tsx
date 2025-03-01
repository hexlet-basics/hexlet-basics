import type { Locale, SharedProps } from "@/types";
import type { PropsWithChildren } from "react";
import React from "react";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type {
  LanguageSitemapLesson,
  SitemapBlogPost,
  SitemapLanguage,
} from "@/types/serializers";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Accordion, Container } from "react-bootstrap";

type Props = PropsWithChildren & {
  title: string;
  orderedLocales: Locale[];
  coursesByLocale: Record<Locale, SitemapLanguage[]>;
  lessonsByLocaleAndLanguageId: Record<
    Locale,
    Record<number, LanguageSitemapLesson[]>
  >;
  blogPostsByLocale: Record<Locale, SitemapBlogPost[]>;
};

type CoursesBlockProps = {
  locale: Locale;
  courses: SitemapLanguage[];
  lessonsByLocaleAndLanguageId: Props["lessonsByLocaleAndLanguageId"];
};

type LessonsBlockProps = {
  lessons: LanguageSitemapLesson[];
  course: SitemapLanguage;
};

type BlogPostsBlockProps = {
  posts: SitemapBlogPost[];
};

const getSuffix = (locale: Locale) => (locale === "en" ? null : locale);

function LessonsBlock({ lessons, course }: LessonsBlockProps) {
  return (
    <ul className="m-0 ms-3 list-unstyled">
      {lessons.map((lesson) => (
        <li key={lesson.id}>
          <a
            className="text-decoration-none"
            href={Routes.language_lesson_path(course.slug, lesson.slug, {
              suffix: getSuffix(lesson.locale),
            })}
          >
            <span className="me-1">{lesson.natural_order}.</span>
            {lesson.name}
          </a>
        </li>
      ))}
    </ul>
  );
}

function CoursesBlock({
  courses,
  lessonsByLocaleAndLanguageId,
}: CoursesBlockProps) {
  const filteredCourses = courses.filter(
    (course) => lessonsByLocaleAndLanguageId[course.locale][course.id],
  );

  return (
    <Accordion defaultActiveKey="0">
      {filteredCourses.map((course, index) => (
        <Accordion.Item eventKey={String(index)} key={course.id}>
          <Accordion.Header as="h3">
            <a
              className="link-body-emphasis"
              href={Routes.language_path(course.slug, {
                suffix: getSuffix(course.locale),
              })}
            >
              {course.name}
            </a>
          </Accordion.Header>
          <Accordion.Body>
            <LessonsBlock
              course={course}
              lessons={lessonsByLocaleAndLanguageId[course.locale][course.id]}
            />
          </Accordion.Body>
        </Accordion.Item>
      ))}
    </Accordion>
  );
}

function BlogPostsBlock({ posts }: BlogPostsBlockProps) {
  return (
    <ul className="m-0 ms-3 list-unstyled">
      {posts.map((post) => (
        <li key={post.id}>
          <a
            className="text-decoration-none"
            href={Routes.blog_post_path(post.slug, {
              suffix: getSuffix(post.locale),
            })}
          >
            {post.name}
          </a>
        </li>
      ))}
    </ul>
  );
}

export default function SiteMap({
  title,
  orderedLocales,
  coursesByLocale,
  lessonsByLocaleAndLanguageId,
  blogPostsByLocale,
}: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout header={title}>
      {orderedLocales.map((locale) => (
        <Container key={locale} className="mb-5">
          <h2 className="mb-3">
            <a
              id={`courses-${locale}`}
              className="text-decoration-none link-body-emphasis"
              href={`#courses-${locale}`}
            >
              {t("home.languages.courses", { lng: locale })}
            </a>
          </h2>
          <CoursesBlock
            locale={locale}
            courses={coursesByLocale[locale]}
            lessonsByLocaleAndLanguageId={lessonsByLocaleAndLanguageId}
          />

          <h2 className="my-3">
            <a
              id={`blog-${locale}`}
              className="text-decoration-none link-body-emphasis"
              href={`#blog-${locale}`}
            >
              {t("blog_posts.index.header", { lng: locale })}
            </a>
          </h2>
          <BlogPostsBlock posts={blogPostsByLocale[locale]} />
        </Container>
      ))}
    </ApplicationLayout>
  );
}
