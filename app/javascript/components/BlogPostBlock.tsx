import type { BlogPost } from "@/types/serializers";
import type { PropsWithChildren } from "react";

import { Card } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";
import dayjs from "dayjs";

type Props = PropsWithChildren & {
  post: BlogPost;
};

export default function BlogPostBlock({ post }: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t: tCommon } = useTranslation("common");

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
            href={Routes.blog_post_path(post.slug!, { suffix })}
            className="link-body-emphasis text-decoration-none stretched-link"
          >
            {post.name}
          </a>
        </Card.Title>
        <Card.Text>{post.description}</Card.Text>
        <div className="mt-2 d-flex small text-muted">
          <div className="me-auto">{dayjs().to(post.created_at)}</div>
          <div className="me-3">
            <i className="bi bi-hand-thumbs-up me-1" />
            {post.likes_count}
          </div>
          <div className="me-1">
            <i className="bi bi-clock me-1" />
          </div>
          <div className="me-2">~{tCommon("time.minutes", { count: 5 })}</div>
          {/* <div> */}
          {/*   {tHelpers("read")} */}
          {/*   <i className="bi bi-arrow-right ms-1 lh-1" /> */}
          {/* </div> */}
        </div>
      </Card.Body>
    </Card>
  );
}
