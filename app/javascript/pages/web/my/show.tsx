import CourseBlock from "@/components/CourseBlock";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { Link } from "@inertiajs/react";
import {
  Button,
  Card,
  Col,
  Container,
  ProgressBar,
  Row,
} from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = {
  startedCourseMembers: LanguageMember[];
  finishedCourseMembers: LanguageMember[];
  landingPageResourcesByCourseId: Record<number, LanguageLandingPageForLists>;
};

function StartedCourse({
  lp,
  cm,
}: { lp: LanguageLandingPageForLists; cm: LanguageMember }) {
  return (
    <Card border="0">
      <Row className="g-0">
        <Col className="col-md-4 d-flex align-items-stretch">
          <img
            alt={lp.header}
            className="img-fluid w-100 h-100 object-fit-cover rounded-start"
            src={lp.language.cover_list_variant}
          />
        </Col>
        <Col className="col-md-8 border rounded-end">
          <Card.Body className="d-flex flex-column h-100">
            <Card.Title className="m-0">
              <a
                href={Routes.language_lesson_url(
                  lp.language.slug!,
                  cm.next_lesson.slug,
                )}
                className="stretched-link text-decoration-none link-body-emphasis h2"
              >
                <h3>{lp.header}</h3>
              </a>
            </Card.Title>
            <Card.Text>{cm.next_lesson.name} →</Card.Text>
            <ProgressBar
              className="mt-auto"
              now={cm.progress}
              label={`${cm.progress}%`}
            />
          </Card.Body>
        </Col>
      </Row>
    </Card>
  );
}

export default function My(props: Props) {
  const {
    startedCourseMembers,
    finishedCourseMembers,
    landingPageResourcesByCourseId,
  } = props;
  const { t: tViews } = useTranslation("web");
  const { t: tCommon } = useTranslation("common");

  return (
    <ApplicationLayout>
      <Container className="mb-lg-5 py-5">
        <div className="mb-5">
          <h2>{tViews("my.started")}</h2>
          {startedCourseMembers.length === 0 && (
            <p className="mt-5">{tCommon("empty")}</p>
          )}
        </div>
        <Row className="row-cols-1 row-cols-sm-2">
          {startedCourseMembers.map((cm) => (
            <Col key={cm.id} className="mb-5">
              <StartedCourse
                cm={cm}
                lp={landingPageResourcesByCourseId[cm.language_id]}
              />
            </Col>
          ))}
        </Row>
        <div className="mb-5">
          <h2>{tViews("my.finished")}</h2>
          {finishedCourseMembers.length === 0 && (
            <p className="mt-5">{tCommon("empty")}</p>
          )}
        </div>
        <Row className="mb-5 row-cols-2 row-cols-md-3 row-cols-lg-4">
          {finishedCourseMembers.map((cm) => (
            <Col key={cm.id}>
              <CourseBlock
                courseMember={cm}
                landingPage={landingPageResourcesByCourseId[cm.language_id]}
              />
            </Col>
          ))}
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
