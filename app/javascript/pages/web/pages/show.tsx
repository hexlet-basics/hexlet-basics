import type { PropsWithChildren } from "react";
import { Col, Container, Row } from "react-bootstrap";
import RuAbout from "./parts/about.ru.tsx";
import EnAbout from "./parts/about.en.tsx";

import RuAuthors from "./parts/authors.ru.tsx";
import EnAuthors from "./parts/authors.en.tsx";

import Application from "@/pages/layouts/ApplicationLayout.tsx";
import type { User } from "@/types/serializers";
import i18next from "i18next";

type Props = PropsWithChildren & {
  page: "about" | "tos";
  user: User;
};

// https://termly.io/
function Empty() {
  return <></>;
}

const mapping = {
  ru: {
    about: RuAbout,
    tos: RuAbout,
    privacy: RuAbout,
    cookie_policy: RuAbout,
    authors: RuAuthors,
  },
  en: {
    about: EnAbout,
    tos: Empty,
    privacy: Empty,
    cookie_policy: Empty,
    authors: EnAuthors,
  },
};

export default function New({ page }: Props) {
  const Component = mapping[i18next.language][page];
  return (
    <Application>
      <Container>
        <Row className="justify-content-center mb-5">
          <Col className="col-12 col-md-10 col-lg-8">
            <Component />
          </Col>
        </Row>
      </Container>
    </Application>
  );
}
