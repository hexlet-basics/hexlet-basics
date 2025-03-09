import cn from "classnames";
import { Alert, Col, Container, ListGroup, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import CourseBlock from "@/components/CourseBlock";
import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem, SharedProps } from "@/types";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageForLists,
  LanguageLesson,
  LanguageMember,
  LanguageModule,
} from "@/types/serializers";
import { Link, usePage } from "@inertiajs/react";

type Props = {
  courseMember?: LanguageMember;
  finishedLessonIds: number[];
  courseCategory: LanguageCategory;
  courseLandingPage: LanguageLandingPage;
  firstLesson: LanguageLesson;
  nextLesson?: LanguageLesson;
  // user: User;
  courseModules: LanguageModule[];
  recommendedCourseLandingPages: LanguageLandingPageForLists[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
};

export default function Show({
  firstLesson,
  nextLesson,
  courseLandingPage,
  courseMember,
  courseCategory,
  courseModules,
  finishedLessonIds,
  recommendedCourseLandingPages,
  lessonsByModuleId,
}: Props) {
  const { t } = useTranslation();
  const { auth } = usePage<SharedProps>().props;

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory.name,
      url: Routes.language_category_path(courseCategory.slug!),
    },
    {
      name: courseLandingPage.header,
      url: Routes.language_path(courseLandingPage.slug),
    },
  ];

  return (
    <ApplicationLayout items={breadcrumbItems}>
      <Container>
        {courseMember?.state === "finished" && (
          <Alert variant="success">
            <XssContent>{t("languages.show.completed_html")}</XssContent>
          </Alert>
        )}
        <div className="p-5 text-center bg-body-tertiary rounded-3 mb-5 border">
          <div className="d-flex justify-content-center align-items-center">
            <h1>{courseLandingPage.header}</h1>
          </div>
          <p className="col-lg-8 mx-auto fs-5 text-muted">
            {courseLandingPage.description}
          </p>
          <Row className="row-cols-1 row-cols-sm-2 justify-content-center">
            {!courseMember && (
              <Col className="mb-3 col-lg-3">
                <Link
                  className="btn d-block btn-primary"
                  href={Routes.language_lesson_path(
                    courseLandingPage.language.slug!,
                    firstLesson.slug,
                  )}
                >
                  <span className="me-2">{t("languages.show.start")}</span>
                  <i className="bi bi-arrow-right" />
                </Link>
              </Col>
            )}
            {courseMember && nextLesson && (
              <Col className="mb-3 col-lg-3">
                <Link
                  className="btn d-block btn-outline-primary"
                  href={Routes.language_lesson_path(
                    courseLandingPage.language.slug!,
                    nextLesson.slug,
                  )}
                >
                  <span className="me-2">{t("languages.show.continue")}</span>
                  <i className="bi bi-arrow-right" />
                </Link>
              </Col>
            )}
            {auth.user.guest && (
              <Col className="mb-3 col-lg-3">
                <Link
                  className="btn d-block btn-outline-secondary"
                  href={Routes.new_user_path()}
                >
                  {t("languages.show.registration")}
                </Link>
              </Col>
            )}
          </Row>
        </div>
        <div className="mb-5">
          {courseModules.map((m) => (
            <div key={m.id} className="mb-5">
              <h2 className="mb-4">{m.name!}</h2>
              <Row className="row-cols-1 row-cols-md-2">
                <Col className="mb-3">
                  <ListGroup>
                    {(lessonsByModuleId[m.id] ?? []).map((l) => (
                      <ListGroup.Item key={l.id}>
                        <Link
                          className="text-decoration-none stretched-link d-flex"
                          href={Routes.language_lesson_path(
                            courseLandingPage.language.slug!,
                            l.slug!,
                          )}
                        >
                          <div className="me-auto">
                            <span className="me-1">{l.natural_order!}.</span>
                            {l.name}
                          </div>
                          {finishedLessonIds.includes(l.id) && (
                            <i className="bi bi-check-lg" />
                          )}
                        </Link>
                      </ListGroup.Item>
                    ))}
                  </ListGroup>
                </Col>
                <Col className="d-none d-sm-block">{m.description}</Col>
              </Row>
            </div>
          ))}

          {courseMember?.state !== "finished" && (
            <div className="p-3 p-sm-4 bg-body-tertiary rounded-3 mb-5">
              <div className="d-flex justify-content-between">
                <div>
                  <div className="display-6 fw-bold">
                    {t("languages.show.ready")}
                  </div>
                  <div className="fs-4 fw-bold text-primary">
                    {t("languages.show.no_registration")}
                  </div>
                </div>
                <div className="align-content-around">
                  <Link
                    className="btn btn-lg btn-outline-primary"
                    href={Routes.language_lesson_path(
                      courseLandingPage.language.slug!,
                      firstLesson.slug!,
                    )}
                  >
                    {t("languages.show.start_demo_lesson")}
                  </Link>
                </div>
              </div>
            </div>
          )}

          <div className="mb-5">
            <div className="d-flex justify-content-between border-bottom mb-3">
              <h2 className="mb-2">{t("languages.show.similar_courses")}</h2>
              <div className="align-content-around">
                <Link
                  className="link-body-emphasis text-muted text-decoration-none"
                  href={Routes.language_category_path(courseCategory.slug!, {})}
                >
                  {t("languages.show.see_all_courses_in_category", {
                    name: courseCategory.name,
                  })}
                </Link>
              </div>
            </div>

            <Row className="row row-cols-2 row-cols-md-4 g-3">
              {recommendedCourseLandingPages.map((lp) => (
                <Col key={lp.id}>
                  <CourseBlock landingPage={lp} />
                </Col>
              ))}
            </Row>
          </div>
        </div>
      </Container>
    </ApplicationLayout>
  );
}
