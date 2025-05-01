import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Col, Container, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

import successImg from "@/images/success-image.svg";

import * as Routes from "@/routes.js";

import i18next from "i18next";

export default function Success() {
  const { t: tViews } = useTranslation("web");

  return (
    <ApplicationLayout>
      <Container className="py-4 py-lg-5">
        <Row className="flex-column flex-lg-row align-items-center gy-5 pt-3">
          <Col className="col-lg-6 col-md-8 order-last order-lg-first d-none d-md-block text-center">
            <img src={successImg} alt="ковер" className="img-fluid" />
          </Col>
          <Col className="col-lg-6">
            <div className="text-center">
              <h1 className="h2 mb-4">{tViews("languages.success.header")}</h1>
              <div className="lh-sm mb-5">
                {tViews("languages.success.description")}
              </div>
              <div className="d-flex flex-column flex-md-row justify-content-center gap-3 gap-lg-4">
                <a
                  className="btn btn-primary btn-lg"
                  href={"https://hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral&utm_campaign=all-courses&utm_content=finished_course_page"}
                  target="_blank"
                  rel="noreferrer"
                >
                  <span>{tViews("languages.success.choose_profession")}</span>
                </a>
                <a className="btn btn-outline-secondary btn-lg" href={Routes.root_path()}>
                  <span className="me-2">
                    {tViews("languages.success.home")}
                  </span>
                  <i className="bi bi-arrow-right lh-1" />
                </a>
              </div>
            </div>
          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
