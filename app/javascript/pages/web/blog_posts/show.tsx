import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import { XBreadcrumb } from "@/components/breadcrumbs";
import Application from "@/pages/layouts/Application";
import type {
  BlogPost,
  Language,
  LanguageCategory,
  User,
} from "@/types/serializers";
import type { BreadcrumbItem } from "@/types/types";
import i18next from "i18next";
import Markdown from "react-markdown";
import remarkSlug from "remark-slug";
import remarkToc from "remark-toc";

import BlogPostBlock from "@/components/BlogPostBlock";
import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  languages: Language[];
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

const plugins = [
  () => remarkToc({ heading: "Содержание" }),
  remarkSlug,
];

export default function New({
  languageCategories,
  languages,
  blogPost,
  recommendedBlogPosts,
}: Props) {
  const { t } = useTranslation();

  const items: BreadcrumbItem[] = [
    {
      name: t("blog_posts.index.header"),
      url: Routes.blog_posts_path(),
    },
    {
      name: blogPost.name!,
      url: Routes.blog_post_path(blogPost.slug!),
    },
  ];

  return (
    <Application languageCategories={languageCategories} languages={languages}>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <XBreadcrumb items={items} />
            <h1 className="mb-5">{blogPost.name}</h1>
            <Markdown remarkPlugins={plugins}>
              {`\n\n## Содержание\n\n ${blogPost.body}`}
            </Markdown>
          </Col>
        </Row>
        <Row className="justify-content-center row-cols-sm-2 row-cols-md-3 row-cols-1">
          {recommendedBlogPosts.map((post) => (
            <Col key={post.id} className="border-top pt-5">
              <BlogPostBlock post={post} />
            </Col>
          ))}
        </Row>
      </Container>
    </Application>
  );
}
