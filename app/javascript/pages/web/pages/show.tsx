import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout.tsx";
import type { User } from "@/types/serializers";
import i18next from "i18next";
import AboutEn from "./parts/about.en";
import AboutRu from "./parts/about.ru";
import AuthorsEn from "./parts/authors.en";
import AuthorsRu from "./parts/authors.ru";
import PrivaryRu from "./parts/privacy.ru";
import TosRu from "./parts/tos.ru";

type Props = PropsWithChildren & {
  page: "about" | "tos" | "privacy" | "cookie_policy" | "authors";
  title: string;
  user: User;
};

// https://termly.io/
function Empty() {
  return <></>;
}

const mapping = {
  ru: {
    about: AboutRu,
    tos: TosRu,
    privacy: PrivaryRu,
    cookie_policy: Empty,
    authors: AuthorsRu,
  },
  en: {
    about: AboutEn,
    tos: Empty,
    privacy: Empty,
    cookie_policy: Empty,
    authors: AuthorsEn,
  },
};

export default function New({ page, title }: Props) {
  const Component = mapping[i18next.language][page];
  return (
    <ApplicationLayout>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <h1 className="mb-5">{title}</h1>
            <Component />
          </Col>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
