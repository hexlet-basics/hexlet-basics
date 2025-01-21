import cn from "classnames";
import type { PropsWithChildren } from "react";
import { Col, Container, ListGroup, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import { deviconClass } from "@/lib/utils";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageCategory,
  LanguageLesson,
  LanguageModule,
  User,
} from "@/types/serializers";
import type { BreadcrumbItem, SharedProps } from "@/types";
import { Link, usePage } from "@inertiajs/react";
import CourseBlock from "@/components/CourseBlock";

type Props = PropsWithChildren & {
  courseCategory: LanguageCategory;
  course: Language;
  firstLesson: LanguageLesson;
  user: User;
  courseModules: LanguageModule[];
  recommendedCourses: Language[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
};

export default function Show({
  firstLesson,
  course,
  courseCategory,
  courseModules,
  recommendedCourses,
  lessonsByModuleId,
}: Props) {
  const { t } = useTranslation();
  const { auth } = usePage<SharedProps>().props;

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory.name!,
      url: Routes.language_category_path(courseCategory.slug!),
    },
    {
      name: course.name!,
      url: Routes.language_path(course.slug!),
    },
  ];

  return (
    <ApplicationLayout>
      <Container>
        <XBreadcrumb items={breadcrumbItems} />
        <div className="p-5 text-center bg-body-tertiary rounded-3 mb-5 border">
          <div className="d-flex justify-content-center align-items-center">
            <i
              className={cn(
                deviconClass(course.slug!),
                "colored",
                "fs-3",
                "me-2",
              )}
            />
            <h1>{course.name}</h1>
          </div>
          <p className="col-lg-8 mx-auto fs-5 text-muted">
            {course.description!}
          </p>
          <div className="d-inline-flex gap-2 mb-5">
            <Link
              className="btn btn-primary"
              href={Routes.language_lesson_path(
                course.slug!,
                firstLesson.slug!,
              )}
            >
              <span className="me-2">{t("languages.show.start")}</span>
              <i className="bi bi-arrow-right" />
            </Link>
            {auth.user.guest && <Link
              className="btn btn-outline-secondary"
              href={Routes.new_user_path()}
            >
              {t("languages.show.registration")}
            </Link>}
          </div>
        </div>
        <div className="mb-5">
          {courseModules.map((m) => (
            <div key={m.id} className="mb-5">
              <h2 className="mb-4">{m.name!}</h2>
              <Row className="col-12">
                <Col>
                  <ListGroup>
                    {lessonsByModuleId[m.id].map((l) => (
                      <ListGroup.Item key={l.id}>
                        <Link
                          className="text-decoration-none stretched-link"
                          href={Routes.language_lesson_path(
                            course.slug!,
                            l.slug!,
                          )}
                        >
                          <span className="me-1">{l.natural_order!}.</span>
                          {l.name}
                        </Link>
                      </ListGroup.Item>
                    ))}
                  </ListGroup>
                </Col>
                <Col>{m.description}</Col>
              </Row>
            </div>
          ))}

          <div className="p-5 bg-body-tertiary rounded-3 mb-5">
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
                    course.slug!,
                    firstLesson.slug!,
                  )}
                >
                  {t("languages.show.start_demo_lesson")}
                </Link>
              </div>
            </div>
          </div>

          <div className="mb-5">
            <div className="d-flex justify-content-between border-bottom mb-3">
              <h2 className="mb-2">{t("languages.show.similar_courses")}</h2>
              <div className="align-content-around">
                <Link
                  className="link-body-emphasis text-muted text-decoration-none"
                  href={Routes.language_category_path(courseCategory.slug!, {
                  })}
                >
                  {t("languages.show.see_all_courses_in_category", {
                    name: courseCategory.name,
                  })}
                </Link>
              </div>
            </div>

            <Row className="row row-cols-2 row-cols-md-4 g-3">
              {recommendedCourses.map((l) => (
                <Col key={l.id}>
                  <CourseBlock course={l} />
                </Col>
              ))}
            </Row>
          </div>
        </div>
      </Container>
    </ApplicationLayout>
  );
}
