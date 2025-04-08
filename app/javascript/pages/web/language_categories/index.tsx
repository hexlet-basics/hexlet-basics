import type { PropsWithChildren } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { LanguageCategory } from "@/types/serializers";

import * as Routes from "@/routes.js";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  categories: LanguageCategory[];
};

export default function Index({ categories }: Props) {
  const { t } = useTranslation();
  const header = t("language_categories.index.header");

  const items = [
    {
      name: header,
      url: Routes.language_categories_path(),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        <Row className="row row-cols-1 row-cols-md-2 row-cols-lg-3 py-3">
          {categories.map((category) => (
            <Col className="col mb-3" key={category.id}>
              <Card className="bg-body-tertiary p-4 rounded-4 shadow-sm h-100 d-flex">
                <div className="h4 fw-bold mb-2">{category.name}</div>
                <Link
                  className="text-decoration-none stretched-link icon-link icon-link-hover mt-auto"
                  href={Routes.language_category_path(category.slug!)}
                >
                  <span className="me-1">
                    {t("language_categories.index.link")}
                  </span>
                  <i className="bi bi-arrow-right lh-1" />
                </Link>
              </Card>
            </Col>
          ))}
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
