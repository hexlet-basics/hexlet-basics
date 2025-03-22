import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Col, Container, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

export default function Success() {
  const { t: tViews } = useTranslation("web");

  return (
    <ApplicationLayout>
      <Container className="pt-3 pt-lg-5">
        <Row className="text-center py-5 my-5">
          <Col>
            <h1 className="h3">{tViews("languages.success.header")}</h1>
            <p>{tViews("languages.success.description")}</p>
          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
