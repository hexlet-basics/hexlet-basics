import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import { XBreadcrumb } from "@/components/breadcrumbs";
import Application from "@/pages/layouts/Application";
import * as Routes from "@/routes.js";
import type { Language, LanguageCategory, User } from "@/types/serializers";
import type { BreadcrumbItem } from "@/types/types";

type Props = PropsWithChildren & {
  languageCategories: LanguageCategory[];
  languageCategory: LanguageCategory;
  language: Language;
  courses: Language[];
  user: User;
};

export default function New({
  languageCategories,
  courses,
  user,
  language,
  languageCategory,
}: Props) {
  const { t } = useTranslation();

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: languageCategory.name!,
      url: Routes.language_category_path(languageCategory.slug!),
    },
    {
      name: language.name!,
      url: Routes.language_path(language.slug!),
    },
  ];

  return (
    <Application languageCategories={languageCategories} courses={courses}>
      <Container>
        <XBreadcrumb items={breadcrumbItems} />
      </Container>
    </Application>
  );
}
