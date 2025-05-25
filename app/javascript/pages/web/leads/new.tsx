import type { PropsWithChildren } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import * as Routes from "@/routes.js";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { XForm, XInput, XSelect } from "@/components/forms";
import { Submit } from "use-inertia-form";
import { enumToOptions } from "@/lib/utils";
import { LeadCrud } from "@/types";
import LeadFormBlock from "@/components/LeadFormBlock";
import XssContent from "@/components/XssContent";

type Props = PropsWithChildren & {
  lead: LeadCrud
};

export default function New({ lead }: Props) {
  const { t } = useTranslation();
  // const { t: tAr } = useTranslation("activerecord");
  // const { t: tHelpers } = useTranslation("helpers");
  const { t: tViews } = useTranslation("web");
  const helpItems = tViews('leads.new.help_items', { returnObjects: true })

  return (
    <ApplicationLayout>
      <Container className="py-4 py-lg-5">
        <Row className="justify-content-center">
          <Col className="col-12 col-lg-7 mb-5">
            <h1 className="h2 mb-4">{tViews("leads.new.header")}</h1>
            <div className="mb-3">
              {tViews("leads.new.description")}
            </div>
            <div className="fw-bold mb-2">
              {tViews('leads.new.how_can_we_help')}
            </div>
            <ul>
              {helpItems.map((item, index) =>
                <li key={index}>
                  <XssContent>
                    {item}
                  </XssContent>
                </li>
              )}
            </ul>
            <div>
              <XssContent>{tViews('leads.new.do_it')}</XssContent>
            </div>
          </Col>
          <Col className="col-lg-5">
            <Card className="bg-body-tertiary h-100 border p-4">
              <Card.Body>
                <LeadFormBlock lead={lead} />
              </Card.Body>
            </Card>

          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
