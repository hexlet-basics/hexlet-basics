import type { PropsWithChildren } from "react";
import { Card, Container, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";
import SignUpFormBlock from "@/components/SignUpFormBlock";

type Props = PropsWithChildren & {
  user: User;
};

export default function New({ user }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center">
          <div className="col-sm-8 col-md-7 col-lg-5">
            <h1 className="text-center mb-3">{t("users.new.sign_up")}</h1>
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
