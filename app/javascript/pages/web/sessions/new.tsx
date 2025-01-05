import type { PropsWithChildren } from "react";
import { Button, Card, Container, Form, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { SignInForm } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  signInForm: SignInForm;
};

export default function New({ signInForm }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  useTranslation("activerecord");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("sessions.new.title")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <Form className="d-flex flex-column">
                  <XInput model={signInForm} attribute="email" />
                  <XInput model={signInForm} attribute="password" />
                  <div className="text-end text-muted small mb-4">
                    {t("sessions.new.forgot_password")}{" "}
                    <Link
                      href={Routes.new_remind_password_path()}
                      className="text-decoration-none"
                    >
                      {t("sessions.new.reset_password")}
                    </Link>
                  </div>
                  <Button size="lg" className="mb-3" type="submit">
                    {tHelpers("submit.user_sign_up_form.create")}
                  </Button>
                </Form>
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
