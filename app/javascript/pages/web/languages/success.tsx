import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Card, Col, Container, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

// import successImg from "@/images/success-image.svg";

import * as Routes from "@/routes.js";

import i18next from "i18next";
import { LanguageLandingPage, LeadCrud } from "@/types";
import XssContent from "@/components/XssContent";
import { XForm, XInput, XSelect } from "@/components/forms";
import { Submit } from "use-inertia-form";
import { enumToOptions } from "@/lib/utils";

type Props = {
  courseLandingPage: LanguageLandingPage;
  lead: LeadCrud
}

export default function Success(props: Props) {
  const { t: tViews } = useTranslation("web");
  const { t: tAr } = useTranslation("activerecord");
  const { t: tHelpers } = useTranslation("helpers");
  const { courseLandingPage, lead } = props

  const contactMethodEnum = tAr('attributes.user.contact_method/values', { returnObjects: true })
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  return (
    <ApplicationLayout>
      <Container className="py-4 py-lg-5">
        <Row>
          {/* <Col className="col-lg-4 col-md-8 order-last order-lg-first d-none d-md-block"> */}
          {/*   <img src={successImg} alt="ковер" className="img-fluid" /> */}
          {/* </Col> */}
          <Col className="col-12 col-lg-7 mb-5">
            <h1 className="h2 mb-4">{tViews("languages.success.header", { name: courseLandingPage.header })}</h1>
            <div className="mb-3">
              {tViews("languages.success.description")}
            </div>
            <div className="fw-bold mb-2">
              {tViews('languages.success.choose_your_path')}
            </div>
            <ul>
              <li>
                <XssContent>
                  {tViews('languages.success.changing_career_html')}
                </XssContent>
              </li>
              <li>
                <XssContent>
                  {tViews('languages.success.getting_new_skill_html')}
                </XssContent>
              </li>
            </ul>
            <div className="fw-bold mb-2">
              {tViews('languages.success.struggle_choosing')}
            </div>
            <div>
              {tViews('languages.success.leave_request')}
            </div>
          </Col>
          <Col className="col-lg-5">
            <Card className="h-100 p-4">
              <Card.Body>
                <XForm className="d-flex flex-column h-100" model="lead" data={lead} to={Routes.leads_path()}>
                  <XSelect
                    name="contact_method"
                    labelField="name"
                    valueField="id"
                    items={contactMethodOptions}
                  />
                  <XInput name="contact_value" />
                  <Submit className="btn d-block btn-lg btn-primary mt-auto">
                    {tHelpers("send")}
                  </Submit>
                </XForm>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
