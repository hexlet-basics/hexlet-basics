import type { BlogPost } from "@/types/serializers";
import type { PropsWithChildren } from "react";

import { Breadcrumb, Card } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

type Props = PropsWithChildren & {
  post: BlogPost;
};

export default function BlogPostBlock({ post }: Props) {
  const { t } = useTranslation();
  return (
    <Card className="border-0">
      {post.cover_thumb_variant && (
        <Card.Img
          src={post.cover_list_variant!}
          // className="img-fluid"
          alt={`Cover for ${post.name}`}
        />
      )}
      <Card.Body className="px-1">
        <Card.Title>
          <a
            href={Routes.blog_post_path(post.slug!)}
            className="link-body-emphasis text-decoration-none stretched-link"
          >
            {post.name}
          </a>
        </Card.Title>
        <Card.Text>{post.description}</Card.Text>
      </Card.Body>
    </Card>
  );
}
