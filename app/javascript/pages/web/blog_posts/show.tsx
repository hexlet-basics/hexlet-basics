import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { BreadcrumbItem, SharedProps } from "@/types";
import type { BlogPost } from "@/types/serializers";
import Markdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import rehypeSlug from "rehype-slug";
import remarkToc from "remark-toc";

import BlogPostBlock from "@/components/BlogPostBlock";
import * as Routes from "@/routes.js";
import { usePage } from "@inertiajs/react";
import type { PluggableList } from "node_modules/react-markdown/lib";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

const rehypePlugins = [rehypeHighlight, rehypeSlug];

export default function Show({ blogPost, recommendedBlogPosts }: Props) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");
  const { suffix } = usePage<SharedProps>().props;

  const heading = tCommon("tos");

  const remarkPlugins: PluggableList = [[remarkToc, { heading }]];

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
    <ApplicationLayout items={items} center header={blogPost.name!}>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <Markdown
              className="hexlet-basics-content"
              remarkPlugins={remarkPlugins}
              rehypePlugins={rehypePlugins}
            >
              {/* {blogPost.body} */}
              {`\n\n## ${heading}\n\n ${blogPost.body}`}
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
