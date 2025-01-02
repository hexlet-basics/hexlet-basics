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
  return (
    <Application languageCategories={languageCategories} languages={languages}>
      <Container>
        <Row className="row row-cols-1 row-cols-sm-2 row-cols-md-3 g-3">
          {blogPosts.map((blogPost) => (
            <Col key={blogPost.id}>
              <Card key={blogPost.id} className="card shadow-sm mb-3">
                <Card.Img variant="top" src={blogPost.cover_list_variant!} />
                <Card.Body>
                  <Card.Title>
                    <Link href={Routes.blog_post_path(blogPost.slug!)}>
                      {blogPost.name}
                    </Link>
                  </Card.Title>
                  <Card.Text>{blogPost.description}</Card.Text>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
        <Paging pagy={pagy} />
      </Container>
    </Application>
  );
}
