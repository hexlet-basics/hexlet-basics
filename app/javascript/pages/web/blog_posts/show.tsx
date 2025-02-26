import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import { highlightingLanguages } from "@/lib/utils.ts";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { BreadcrumbItem, SharedProps } from "@/types";
import type { BlogPost } from "@/types/serializers";
import Markdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import rehypeSlug from "rehype-slug";
import remarkGfm from "remark-gfm";
import remarkToc from "remark-toc";

import BlogPostBlock from "@/components/BlogPostBlock";
import * as Routes from "@/routes.js";
import { usePage } from "@inertiajs/react";
import { useTranslation } from "react-i18next";
import type { Pluggable } from "unified";

type Props = PropsWithChildren & {
  blogPost: BlogPost;
  recommendedBlogPosts: BlogPost[];
};

const rehypePlugins: Pluggable[] = [
  [rehypeHighlight, { languages: highlightingLanguages }],
  rehypeSlug,
];

export default function Show({ blogPost, recommendedBlogPosts }: Props) {
  const { t } = useTranslation();
  const { t: tCommon } = useTranslation("common");
  const { suffix } = usePage<SharedProps>().props;

  const heading = tCommon("tos");

  const remarkPlugins: Pluggable[] = [[remarkToc, { heading }], [remarkGfm]];

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
            <div className="hexlet-basics-content">
              <Markdown
                remarkPlugins={remarkPlugins}
                rehypePlugins={rehypePlugins}
              >
                {/* {blogPost.body} */}
                {`\n\n## ${heading}\n\n ${blogPost.body}`}
              </Markdown>
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
  );
}
