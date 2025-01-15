import type { FormEvent, PropsWithChildren } from "react";
import { Button, Card, Container, Form, Row } from "react-bootstrap";
import { useInertiaForm, Submit } from "use-inertia-form";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import XssContent from "@/components/XssContent";
import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  user: User;
};

export default function New({ user }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("users.new.sign_up")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <XForm
                  model="user_sign_up_form"
                  data={{ user_sign_up_form: user }}
                  to={Routes.users_path()}
                >
                  <XInput name="first_name" autoComplete="name" />
                  <XInput name="email" autoComplete="email" />
                  <XInput
                    name="password"
                    type="password"
                    autoComplete="current-password"
                  />
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
                    {tHelpers("submit.user_sign_up_form.create")}
                  </Submit>
                  <XssContent className="small text-muted">
                    {t("users.new.confirmation_html", {
                      url: Routes.page_path("tos"),
                    })}
                  </XssContent>
                </XForm>
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
