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

  return (
    <ApplicationLayout>
      <Container className="mb-lg-5 py-5">
        <h2 className="mb-5">{tViews("my.started")}</h2>
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
        <h2 className="mb-5">{tViews("my.finished")}</h2>
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
