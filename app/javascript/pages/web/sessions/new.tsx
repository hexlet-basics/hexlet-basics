import type { PropsWithChildren } from "react";
import { Button, Card, Container, Form, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { SignInForm } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { Submit } from "use-inertia-form";

type Props = PropsWithChildren & {
  signInForm: SignInForm;
};

export default function New({ signInForm }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("sessions.new.title")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <XForm
                  to={Routes.session_path()}
                  data={{ user_sign_in_form: signInForm }}
                  model="user_sign_in_form"
                  className="d-flex flex-column"
                >
                  <XInput name="email" autoComplete="email" />
                  <XInput
                    name="password"
                    type="password"
                    autoComplete="current-password"
                  />
                  <div className="text-end text-muted small mb-4">
                    {t("sessions.new.forgot_password")}{" "}
                    <Link
                      href={Routes.new_remind_password_path()}
                      className="text-decoration-none"
                    >
                      {t("sessions.new.reset_password")}
                    </Link>
                  </div>
                  <Submit className="btn w-100 btn-lg btn-primary mb-3">
                    {tHelpers("submit.user_sign_in_form.create")}
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
