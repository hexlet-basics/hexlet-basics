import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";

import Application from "@/pages/layouts/ApplicationLayout.tsx";
import type { User } from "@/types/serializers";
import i18next from "i18next";
import TosRu from "./parts/tos.ru";
import AboutRu from "./parts/about.ru";
import PrivaryRu from "./parts/privacy.ru";
import AuthorsRu from "./parts/authors.ru";
import AboutEn from "./parts/about.en";
import AuthorsEn from "./parts/authors.en";

type Props = PropsWithChildren & {
  page: "about" | "tos" | "privacy" | "cookie_policy" | "authors";
  title: string
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
    <Application>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <h1 className="mb-5">{title}</h1>
            <Component />
          </Col>
        </Row>
      </Container>
    </Application>
  );
}
