import type { PropsWithChildren } from "react";
import { Button, Card, Container, Form, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  user: User;
};

export default function New({ user }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  useTranslation("activerecord");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("users.new.sign_up")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <Form className="d-flex flex-column">
                  <XInput model={user} attribute="first_name" />
                  <XInput model={user} attribute="email" />
                  <div className="text-end text-muted small mb-4">
                    {t("users.new.have_account")}{" "}
                    <Link
                      href={Routes.new_session_path()}
                      className="text-decoration-none"
                    >
                      {t("users.new.sign_in")}
                    </Link>
                  </div>
                  <Button size="lg" className="mb-3" type="submit">
                    {tHelpers("submit.user_sign_up_form.create")}
                  </Button>
                  <div
                    className="small text-muted"
                    // biome-ignore lint/security/noDangerouslySetInnerHtml: <explanation>
                    dangerouslySetInnerHTML={{
                      __html: t("users.new.confirmation_html"),
                    }}
                  />
                </Form>
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
