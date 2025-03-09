import { usePage } from "@inertiajs/react";
import { Nav } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type ReviewCrud from "@/types/serializers/ReviewCrud";

type Props = {
  data?: ReviewCrud;
};

export function Menu({ data }: Props) {
  const { url } = usePage();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Nav variant="tabs" className="mb-4" activeKey={url}>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.admin_reviews_path()}
        >
          {tHelpers("crud.list")}
        </Nav.Link>
      </Nav.Item>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.new_admin_review_path()}
        >
          {tHelpers("crud.add")}
        </Nav.Link>
      </Nav.Item>
      {data && (
        <Nav.Item>
          <Nav.Link
            className="link-body-emphasis"
            href={Routes.edit_admin_review_path(data.review.id)}
          >
            {tHelpers("crud.editing")}
          </Nav.Link>
        </Nav.Item>
      )}
    </Nav>
  );
}
