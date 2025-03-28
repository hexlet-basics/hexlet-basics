import { usePage } from "@inertiajs/react";
import { Nav } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type { BlogPostCrud } from "@/types/serializers";

type Props = {
  data?: BlogPostCrud;
};

export function Menu({ data }: Props) {
  const { url } = usePage();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Nav variant="tabs" className="mb-4" activeKey={url}>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.admin_blog_posts_path()}
        >
          {tHelpers("crud.list")}
        </Nav.Link>
      </Nav.Item>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.new_admin_blog_post_path()}
        >
          {tHelpers("crud.add")}
        </Nav.Link>
      </Nav.Item>
      {data && (
        <>
          <Nav.Item>
            <Nav.Link
              className="link-body-emphasis"
              href={Routes.edit_admin_blog_post_path(data.blog_post.id)}
            >
              {tHelpers("crud.editing")}
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link
              target="_blank"
              className="link-body-emphasis"
              href={Routes.blog_post_path(data.blog_post.slug!)}
            >
              <i className="bi bi-arrow-up-right-square" />
            </Nav.Link>
          </Nav.Item>
        </>
      )}
    </Nav>
  );
}
