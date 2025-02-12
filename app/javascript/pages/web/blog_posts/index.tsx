import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BlogPost, Pagy } from "@/types/serializers";

import BlogPostBlock from "@/components/BlogPostBlock";
import type { BreadcrumbItem, SharedProps } from "@/types";
import { useTranslation } from "react-i18next";
import { usePage } from "@inertiajs/react";

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  pagy: Pagy;
};

export default function Index({ blogPosts }: Props) {
  const { t } = useTranslation();
  const header = t("blog_posts.index.header");
  const { suffix } = usePage<SharedProps>().props;

  const items: BreadcrumbItem[] = [
    {
      name: header,
      url: Routes.blog_posts_path({ suffix }),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        <Row className="row row-cols-1 row-cols-sm-2 row-cols-md-3 g-3">
          {blogPosts.map((post) => (
            <Col key={post.id}>
              <BlogPostBlock post={post} />
            </Col>
          ))}
        </Row>
        {/* <XPaging pagy={pagy} /> */}
      </Container>
    </ApplicationLayout>
  );
}
