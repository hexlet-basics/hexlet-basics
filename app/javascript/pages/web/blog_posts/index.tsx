import type { PropsWithChildren } from "react";
import { Card, CardImg, CardText, Col, Container, Row } from "react-bootstrap";

import Paging from "@/components/Paging";
import Application from "@/pages/layouts/Application";
import type {
  BlogPost,
  Language,
  LanguageCategory,
  Pagy,
} from "@/types/serializers";

import * as Routes from "@/routes.js";
import { Link } from "@inertiajs/react";
import BlogPostBlock from "@/components/BlogPostBlock";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  languages: Language[];
  blogPosts: BlogPost[];
  pagy: Pagy;
};

export default function New({
  languageCategories,
  languages,
  blogPosts,
  pagy,
}: Props) {
  const { t } = useTranslation()

  return (
    <Application languageCategories={languageCategories} languages={languages}>
      <Container>
        <h1>{t('blog_posts.index.header')}</h1>
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
