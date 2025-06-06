import CourseBlock from "@/components/CourseBlock";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { Link } from "@inertiajs/react";
import {
  Card,
  Col,
  Container,
  ProgressBar,
  Row
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
    <Card className="overflow-hidden h-100">
      <Row className="g-0 h-100">
        <Col className="col-lg-4">
          <img
            alt={lp.header}
            className="img-fluid w-100 h-100 object-fit-fill"
            src={lp.language.cover_list_variant}
          />
        </Col>
        <Col className="col-6 col-lg-8">
          <Card.Body className="d-flex flex-column h-100">
            <Card.Title className="m-0">
              <Link
                href={Routes.language_url(
                  lp.slug!,
                )}
                className="stretched-link text-decoration-none link-body-emphasis h2"
              >
                <h3>{lp.header}</h3>
              </Link>
            </Card.Title>
            <Card.Text className="pe-3">{cm.next_lesson.name} →</Card.Text>
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
        <Row className="row-cols-1 row-cols-md-2 gy-4 gy-lg-5 pb-5">
          {startedCourseMembers.map((cm) => (
            <Col key={cm.id}>
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
