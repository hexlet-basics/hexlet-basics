import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageForLists,
} from "@/types/serializers";

import CourseBlock from "@/components/CourseBlock";
import * as Routes from "@/routes.js";

type Props = PropsWithChildren & {
  categoryLandingPages: LanguageLandingPageForLists[];
  landingPages: LanguageLandingPage[];
  courseCategory: LanguageCategory;
};

export default function Show({ courseCategory, categoryLandingPages }: Props) {
  const { t } = useTranslation();

  const header = t("language_categories.show.header", {
    name: courseCategory.name!,
  });

  const items = [
    {
      name: t("language_categories.index.header"),
      url: Routes.language_categories_path(),
    },
    {
      name: header,
      url: Routes.language_category_path(courseCategory.slug!),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        <div className="row row-cols-2 row-cols-md-3 row-cols-lg-4">
          {categoryLandingPages.map((lp) => (
            <div className="col mb-3" key={lp.id}>
              <CourseBlock landingPage={lp} />
            </div>
          ))}
        </div>
      </Container>
    </ApplicationLayout>
  );
}
