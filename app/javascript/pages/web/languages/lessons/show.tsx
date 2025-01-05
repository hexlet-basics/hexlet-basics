import type { PropsWithChildren } from "react";
import { Card, Col, Container, Nav, Row, Tab, Tabs } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import type { SharedProps } from "@/types/types";
import { usePage } from "@inertiajs/react";
import LessonLayout from "@/pages/layouts/LessonLayout";
import type { Language, LanguageLesson } from "@/types/serializers";
import Markdown from "react-markdown";

type Props = PropsWithChildren & {
  course: Language;
  lesson: LanguageLesson;
};

export default function Show({ course, lesson }: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t } = useTranslation();

  return (
    <LessonLayout>
      <Container fluid className="overflow-hidden mb-1 h-100">
        <Row>
          <Col className="h-100 col-12 col-md-6 col-lg-5 mb-3 mb-md-0 position-relative">
            <Tab.Container id="left-tabs-example" defaultActiveKey="lesson">
              <Nav variant="underline" fill justify className="mb-3 small">
                <Nav.Item>
                  <Nav.Link className="link-body-emphasis" eventKey="lesson">
                    {t('languages.lessons.show.lesson')}
                  </Nav.Link>
                </Nav.Item>
                <Nav.Item>
                  <Nav.Link className="link-body-emphasis" eventKey="discuss">
                    {t('languages.lessons.show.discuss')}
                  </Nav.Link>
                </Nav.Item>
                <Nav.Item>
                  <Nav.Link className="link-body-emphasis" eventKey="navigation">
                    {t('languages.lessons.show.navigation')}
                  </Nav.Link>
                </Nav.Item>
              </Nav>

              <Tab.Content className="p-2">
                <Tab.Pane eventKey="lesson">
                  <h1 className="h2">{lesson.name}</h1>
                  <Markdown>{lesson.instructions}</Markdown>
                  First tab content
                </Tab.Pane>
                <Tab.Pane eventKey="discuss">Second tab content</Tab.Pane>
                <Tab.Pane eventKey="navigation">Second tab content</Tab.Pane>
              </Tab.Content>
            </Tab.Container>
          </Col>
          <Col className="h-100 col-12 col-md-6 col-lg-7 pl-md-0">
            <Card className="h-100">
              <Card.Body>wow</Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </LessonLayout>
  );
}
