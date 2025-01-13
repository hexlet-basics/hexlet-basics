import { Nav } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type { BlogPost } from "@/types/serializers";
import { url } from "@/lib/utils";

type Props = {
  data?: BlogPost;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation("helpers");
  return (
    <Nav variant="tabs" className="mb-4" activeKey={url({ onlyPath: true })}>
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
        <Nav.Item>
          <Nav.Link
            className="link-body-emphasis"
            href={Routes.edit_admin_blog_post_path(data.id)}
          >
            {tHelpers("crud.editing")}
          </Nav.Link>
        </Nav.Item>
      )}
    </Nav>
  );
}

