import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";

type Props = {
  code: string;
  message: string;
};

export default function Show({ code, message }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container>
        <h1 className="mb-5">{`${code} ${message}`}</h1>
      </Container>
    </ApplicationLayout>
  );
}
