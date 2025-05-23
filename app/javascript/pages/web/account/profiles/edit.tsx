import type { PropsWithChildren } from "react";
import { Card, Container, Row } from "react-bootstrap";
import { Submit } from "use-inertia-form";

import { useTranslation } from "react-i18next";

import * as Routes from "@/routes.js";

import { XForm, XInput, XSelect } from "@/components/forms";
import useConfirmation from "@/hooks/useConfirmation";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { UserProfileForm } from "@/types/serializers";
import { Link } from "@inertiajs/react";
import { enumToOptions } from "@/lib/utils";

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function Edit(props: Props) {
  const { form } = props
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tAr } = useTranslation("activerecord");

  const confirmDeleting = useConfirmation();

  return (
    <ApplicationLayout>
      <Container className="mt-5">
        <Row className="justify-content-center">
          <div className="col-sm-10 col-lg-7">
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
                  <XInput name="last_name" autoComplete="name" />
                  <Submit className="btn w-100 btn-lg btn-primary mb-3">
                    {tHelpers("submit.save")}
                  </Submit>
                </XForm>
                <Link
                  onClick={confirmDeleting}
                  as="button"
                  className="mt-5 btn btn-sm btn-outline-danger"
                  method="delete"
                  href={Routes.account_profile_path()}
                >
                  {t("account.profiles.edit.delete")}
                </Link>
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
