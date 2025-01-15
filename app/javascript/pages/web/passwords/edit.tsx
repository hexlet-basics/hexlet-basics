import type { PropsWithChildren } from "react";
import { Card, Container, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { UserPassword } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  userPassword: UserPassword;
};

export default function New({ userPassword }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("passwords.edit.title")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <XForm
                  method="patch"
                  model="user_password_form"
                  data={{ user_password_form: userPassword }}
                  to={Routes.password_path()}
                >
                  <XInput name="password" type="password" autoComplete="new-password" />
                  <Submit className="btn w-100 btn-lg btn-primary mb-3">
                    {tHelpers("submit.replace")}
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
