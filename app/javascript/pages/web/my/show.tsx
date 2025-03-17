import CourseBlock from "@/components/CourseBlock";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type {
  LanguageLandingPage,
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { Button, Col, Container, Row } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = {
  startedCourseMembers: LanguageMember[];
  finishedCourseMembers: LanguageMember[];
  landingPageResourcesByCourseId: Record<number, LanguageLandingPageForLists>;
};

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
        <Row className="row-cols-2 row-cols-md-3 row-cols-lg-4">
          {startedCourseMembers.map((cm) => (
            <Col key={cm.id} className="mb-5">
              <CourseBlock
                courseMember={cm}
                landingPage={landingPageResourcesByCourseId[cm.language_id]}
                // continueButton
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
