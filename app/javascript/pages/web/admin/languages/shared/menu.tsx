import { usePage } from "@inertiajs/react";
import { Nav } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";
import type { OriginalLanguage } from "@/types/serializers";

type Props = {
  data?: OriginalLanguage;
};

export function Menu({ data }: Props) {
  const { url } = usePage();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <Nav variant="tabs" className="mb-4" activeKey={url}>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.admin_languages_path()}
        >
          {tHelpers("crud.list")}
        </Nav.Link>
      </Nav.Item>
      <Nav.Item>
        <Nav.Link
          className="link-body-emphasis"
          href={Routes.new_admin_language_path()}
        >
          {tHelpers("crud.add")}
        </Nav.Link>
      </Nav.Item>
      {data && (
        <>
          <Nav.Item>
            <Nav.Link
              className="link-body-emphasis"
              href={Routes.edit_admin_language_path(data.id)}
            >
              {tHelpers("crud.editing")}
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link
              target="_blank"
              className="link-body-emphasis"
              href={Routes.language_path(data.slug!)}
            >
              <i className="bi bi-arrow-up-right-square" />
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link
              target="_blank"
              className="link-body-emphasis"
              href={data.repository_url!}
            >
              <i className="bi bi-github" />
            </Nav.Link>
          </Nav.Item>
        </>
      )}
    </Nav>
  );
}
