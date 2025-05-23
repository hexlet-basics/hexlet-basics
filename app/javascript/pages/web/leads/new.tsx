import type { PropsWithChildren } from "react";
import { Card, Container, Row } from "react-bootstrap";
import * as Routes from "@/routes.js";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { XForm, XInput, XSelect } from "@/components/forms";
import { Submit } from "use-inertia-form";
import { enumToOptions } from "@/lib/utils";
import { LeadCrud } from "@/types";

type Props = PropsWithChildren & {
  lead: LeadCrud
};

export default function New({ lead }: Props) {
  const { t } = useTranslation();
  const { t: tAr } = useTranslation("activerecord");
  const { t: tHelpers } = useTranslation("helpers");

  const contactMethodEnum = tAr('attributes.user.contact_method/values', { returnObjects: true })
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center my-5 mb-5">{t("leads.new.header")}</h1>
            <Card className="p-4 border-0">
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

          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
