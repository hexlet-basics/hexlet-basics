import cn from "classnames";
import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import { deviconClass } from "@/lib/utils";
import Application from "@/pages/layouts/Application";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageCategory,
  LanguageLesson,
  User,
} from "@/types/serializers";
import type { BreadcrumbItem } from "@/types/types";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  languageCategory: LanguageCategory;
  course: Language;
  firstLesson: LanguageLesson;
  user: User;
};

export default function New({ firstLesson, course, languageCategory }: Props) {
  const { t } = useTranslation();

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: languageCategory.name!,
      url: Routes.language_category_path(languageCategory.slug!),
    },
    {
      name: course.name!,
      url: Routes.language_path(course.slug!),
    },
  ];

  return (
    <Application>
      <Container>
        <XBreadcrumb items={breadcrumbItems} />
        <div className="p-5 text-center bg-body-tertiary rounded-3">
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
            <Link
              className="btn btn-outline-secondary"
              href={Routes.new_user_path()}
            >
              {t("languages.show.registration")}
            </Link>
          </div>
        </div>
      </Container>
    </Application>
  );
}
