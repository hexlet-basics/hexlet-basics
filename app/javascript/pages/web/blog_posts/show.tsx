import type { PropsWithChildren } from "react";
import { Col, Container, Image, Row } from "react-bootstrap";

import { dayjs } from "@/lib/utils.ts";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { BreadcrumbItem, SharedProps } from "@/types";
import type { BlogPost } from "@/types/serializers";
import i18next from "i18next";

import BlogPostBlock from "@/components/BlogPostBlock";
import MarkdownViewer from "@/components/MarkdownViewer.tsx";
import * as Routes from "@/routes.js";
import { Head, usePage } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import type { Article, WithContext } from "schema-dts";

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

export default function Show({ blogPost, recommendedBlogPosts }: Props) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");

  const postUrl = Routes.blog_post_url(blogPost.slug!);
  const items: BreadcrumbItem[] = [
    {
      name: t("blog_posts.index.header"),
      url: Routes.blog_posts_url(),
    },
    {
      name: blogPost.name!,
      url: postUrl,
    },
  ];

  const article: WithContext<Article> = {
    "@context": "https://schema.org",
    "@type": "Article",
    author: blogPost.creator.name!,
    name: blogPost.name!,
    datePublished: blogPost.created_at,
    headline: blogPost.description!,
    image: blogPost.cover_main_variant!,
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">{JSON.stringify(article)}</script>
      </Head>
      <ApplicationLayout items={items} center header={blogPost.name!}>
        <Container>
          <Row className="justify-content-center mb-3">
            <Col className="col-12 col-md-10 col-lg-8 d-flex flex-column">
              <div className="mb-4">
                <Image
                  className="img-fluid"
                  src={blogPost.cover_main_variant!}
                />
              </div>
              <div className="hexlet-basics-content">
                <MarkdownViewer>{blogPost.body || ""}</MarkdownViewer>
              </div>
              {i18next.language === "ru" && (
                <Row className="justify-content-center py-4">
                  <Col className="col-lg-11">
                    <div className="bg-primary rounded-4 shadow text-light py-4 px-4 px-md-5 position-relative">
                      <div className="d-flex flex-column flex-md-row gap-4">
                        <div className="fs-3 d-none d-md-block">
                          <i className="bi bi-chat" />
                        </div>
                        <div>
                          <div className="h2 fw-semibold mb-1">
                            {t("blog_posts.show.join_community")}
                          </div>
                          <p className="fs-5 lh-sm mb-4">
                            {t("blog_posts.show.discuss")}
                          </p>
                          <a
                            className="link-light text-decoration-none stretched-link icon-link icon-link-hover"
                            href="https://ttttt.me/HexletLearningBot"
                            target="_blank"
                            rel="noopener noreferrer"
                          >
                            <span className="me-2">
                              {t("blog_posts.show.link")}
                            </span>
                            <span className="bi bi-arrow-right lh-1" />
                          </a>
                        </div>
                      </div>
                    </div>
                  </Col>
                </Row>
              )}
              <div className="mt-5 d-flex text-muted">
                <div className="me-2">
                  <b>
                    <i className="bi bi-person-circle me-2" />
                    {blogPost.creator.name}
                  </b>
                </div>
                <div className="me-auto">{dayjs().to(blogPost.created_at)}</div>
                <div className="me-3">
                  <a href={postUrl} className="link-body-emphasis">
                    <i className="bi bi-hand-thumbs-up me-1" />
                  </a>
                  {blogPost.likes_count}
                </div>
                <div className="text-muted">
                  <i className="bi bi-clock me-2" />~
                  {tCommon("time.minutes", { count: 5 })}
                </div>
              </div>
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
      </ApplicationLayout>
    </>
  );
}
