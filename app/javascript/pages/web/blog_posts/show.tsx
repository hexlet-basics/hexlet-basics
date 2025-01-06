import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import { XBreadcrumb } from "@/components/breadcrumbs";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { BlogPost } from "@/types/serializers";
import type { BreadcrumbItem, SharedProps } from "@/types/types";
import Markdown from "react-markdown";
import rehypeSlug from "rehype-slug";
import remarkToc from "remark-toc";
import rehypeHighlight from "rehype-highlight";

import BlogPostBlock from "@/components/BlogPostBlock";
import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";
import { usePage } from "@inertiajs/react";
import type { PluggableList } from "node_modules/react-markdown/lib";

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

const remarkPlugins: PluggableList = [[remarkToc, { heading: "Содержание" }]];
const rehypePlugins = [rehypeHighlight, rehypeSlug];

export default function New({ blogPost, recommendedBlogPosts }: Props) {
  const { t } = useTranslation();
  const { suffix } = usePage<SharedProps>().props;

  const items: BreadcrumbItem[] = [
    {
      name: t("blog_posts.index.header"),
      url: Routes.blog_posts_path({ suffix }),
    },
    {
      name: blogPost.name!,
      url: Routes.blog_post_path(blogPost.slug!, { suffix }),
    },
  ];

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <XBreadcrumb items={items} />
            <h1 className="mb-5">{blogPost.name}</h1>
            <Markdown
              remarkPlugins={remarkPlugins}
              rehypePlugins={rehypePlugins}
            >
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
    </ApplicationLayout>
  );
}
