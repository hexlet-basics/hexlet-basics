import cn from "classnames";
import type { PropsWithChildren } from "react";
import { Col, Container, ListGroup, Row } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import { deviconClass } from "@/lib/utils";
import Application from "@/pages/layouts/Application";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageCategory,
  LanguageLesson,
  LanguageModule,
  User,
} from "@/types/serializers";
import type { BreadcrumbItem, SharedProps } from "@/types/types";
import { Link, usePage } from "@inertiajs/react";

type Props = PropsWithChildren & {
  courseCategory: LanguageCategory;
  course: Language;
  firstLesson: LanguageLesson;
  user: User;
  courseModules: LanguageModule[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
};

export default function New({
  firstLesson,
  course,
  courseCategory,
  courseModules,
  lessonsByModuleId,
}: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t } = useTranslation();

  console.log(courseModules, lessonsByModuleId)

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory.name!,
      url: Routes.language_category_path(courseCategory.slug!, { suffix }),
    },
    {
      name: course.name!,
      url: Routes.language_path(course.slug!, { suffix }),
    },
  ];

  return (
    <Application>
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
                { suffix },
              )}
            >
              <span className="me-2">{t("languages.show.start")}</span>
              <i className="bi bi-arrow-right" />
            </Link>
            <Link
              className="btn btn-outline-secondary"
              href={Routes.new_user_path({ suffix })}
            >
              {t("languages.show.registration")}
            </Link>
          </div>
        </div>
        <div>
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
                          href={Routes.language_lesson_path(course.slug!, l.slug!, { suffix })}
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
        </div>
      </Container>
    </Application>
  );
}
