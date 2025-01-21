import type { PropsWithChildren } from "react";
import { Card, Container, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { PasswordReminder } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  passwordReminder: PasswordReminder;
};

export default function New({ passwordReminder }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">
              {t("remind_passwords.new.title")}
            </h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <XForm
                  model="remind_password_form"
                  data={{ remind_password_form: passwordReminder }}
                  to={Routes.remind_password_path()}
                >
                  <XInput name="email" autoComplete="email" />
                  <div className="text-end text-muted small mb-4">
                    {t("users.new.have_account")}{" "}
                    <Link
                      href={Routes.new_session_path()}
                      className="text-decoration-none"
                    >
                      {t("users.new.sign_in")}
                    </Link>
                  </div>
                  <Submit className="btn w-100 btn-lg btn-primary mb-3">
                    {tHelpers("submit.remind_password_form.create")}
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
