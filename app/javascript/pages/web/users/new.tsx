import type { PropsWithChildren } from "react";
import { Alert, Card, Container, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import SignUpFormBlock from "@/components/SignUpFormBlock";
import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import { useForm } from "@inertiajs/react";

type Props = PropsWithChildren & {
  user: User;
  demo: boolean;
};

export default function New({ user, demo }: Props) {
  const { t } = useTranslation();
  const data = useForm()
  console.log(data)

  return (
    <ApplicationLayout>
      <Container>
        {demo && (
          <Alert variant="info">
            <XssContent>{t("users.new.demo_html")}</XssContent>
          </Alert>
        )}
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center my-5 mb-3">{t("users.new.sign_up")}</h1>
            <Card className="p-4 border-0">
              <Card.Body>
                <SignUpFormBlock user={user} />
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
