import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import Paging from "@/components/Paging";
import Application from "@/pages/layouts/Application";
import type { BlogPost, Pagy } from "@/types/serializers";

import BlogPostBlock from "@/components/BlogPostBlock";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  pagy: Pagy;
};

export default function New({ blogPosts, pagy }: Props) {
  const { t } = useTranslation();

  return (
    <Application>
      <Container>
        <h1>{t("blog_posts.index.header")}</h1>
        <Row className="row row-cols-1 row-cols-sm-2 row-cols-md-3 g-3">
          {blogPosts.map((post) => (
            <Col key={post.id}>
              <BlogPostBlock post={post} />
            </Col>
          ))}
        </Row>
        <Paging pagy={pagy} />
      </Container>
    </Application>
  );
}