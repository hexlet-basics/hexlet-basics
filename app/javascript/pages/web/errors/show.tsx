import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";

type Props = {
  code: string;
  header: string;
  description: string;
};

export default function Show({ code, header, description }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container className="text-center my-5 py-lg-5">
        <div className="fw-bold">{code}</div>
        <h1 className="display-1">{header}</h1>
        <p className="lead">{description}</p>
      </Container>
    </ApplicationLayout>
  );
}
