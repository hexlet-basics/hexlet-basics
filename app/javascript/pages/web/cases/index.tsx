import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Card, Col, Container, Row } from "react-bootstrap";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

export default function Index() {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container className="h-100 py-3">
        <Row className="row row-cols-1 row-cols-md-2 row-cols-lg-3 py-3">
          <Col>
            <Card className="bg-body-tertiary p-4 rounded-4 shadow-sm h-100 d-flex">
              <div className="h4 fw-bold mb-2">
                <span>{t("cases.for_teachers")}</span>
              </div>
              <a
                className="text-decoration-none stretched-link icon-link icon-link-hover mt-auto"
                href={Routes.for_teachers_cases_path()}
              >
                <span>{t("cases.link")}</span>
                <i className="bi bi-arrow-right lh-1" />
              </a>
            </Card>
          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
