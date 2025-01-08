import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";


import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { User } from "@/types/serializers";

type Props = PropsWithChildren & {
};

export default function Index() {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <ApplicationLayout>
      <Container>
      </Container>
    </ApplicationLayout>
  );
}

