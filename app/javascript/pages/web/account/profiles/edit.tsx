import type { FormEvent, PropsWithChildren } from "react";
import { Button, Card, Container, Form, Row } from "react-bootstrap";
import { useInertiaForm, Submit } from "use-inertia-form";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import XssContent from "@/components/XssContent";
import { XForm, XInput } from "@/components/forms";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User, UserProfileForm } from "@/types/serializers";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function New({ form }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">
              {t("account.profiles.edit.title")}
            </h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <XForm
                  model="user"
                  method="patch"
                  data={{ user: form }}
                  to={Routes.account_profile_path()}
                >
                  <XInput name="first_name" autoComplete="name" />
                  <Submit className="btn w-100 btn-lg btn-primary mb-3">
                    {tHelpers("submit.save")}
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
