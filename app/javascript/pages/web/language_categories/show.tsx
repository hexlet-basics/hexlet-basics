import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { Language, LanguageCategory, User } from "@/types/serializers";

import CourseBlock from "@/components/CourseBlock";
import * as Routes from "@/routes.js";

type Props = PropsWithChildren & {
  categoryCourses: Language[];
  courses: Language[];
  courseCategory: LanguageCategory;
  user: User;
};

export default function New({ courseCategory, categoryCourses }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");
  const { t: tAr } = useTranslation("activerecord");

  const items = [
    {
      name: courseCategory.name!,
      url: Routes.language_category_path(courseCategory.slug!),
    },
  ];

  return (
    <ApplicationLayout>
      <Container>
        <XBreadcrumb items={items} />
        <h1 className="mb-5">{courseCategory.name}</h1>
        <div className="row row-cols-2 row-cols-md-3 row-cols-lg-4">
          {categoryCourses.map((course) => (
            <div className="col mb-3" key={course.id}>
              <CourseBlock course={course} />
            </div>
          ))}
        </div>
      </Container>
    </ApplicationLayout>
  );
}
