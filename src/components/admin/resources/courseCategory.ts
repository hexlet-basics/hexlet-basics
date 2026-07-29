import { useTranslation } from "react-i18next";
import type { CourseCategory, CourseCategoryInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Form values ARE the generated CourseCategoryInput, so the form validates against
// zCourseCategoryInput directly. The nullable `description` renders as "" (the
// engine's text fields coerce null); a fresh form leaves it null.

export const emptyCourseCategory: CourseCategoryInput = {
  name: "",
  header: "",
  slug: "",
  description: null,
};

// Seed the form from an existing row: the required columns are nullable on the
// read model, so coalesce them to "" to keep the inputs controlled.
export function courseCategoryToForm(category: CourseCategory): CourseCategoryInput {
  return {
    name: category.name ?? "",
    header: category.header ?? "",
    slug: category.slug ?? "",
    description: category.description,
  };
}

// Field descriptors driving CrudForm. A hook so labels resolve through the typed
// i18n selector; shared by the create and edit screens.
export function useCourseCategoryFields(): CrudFieldSpec<CourseCategoryInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "name",
      label: t(($) => $.models.attributes.language_category.name),
      required: true,
    },
    {
      name: "header",
      label: t(($) => $.models.attributes.language_category.header),
      required: true,
    },
    {
      name: "slug",
      label: t(($) => $.models.attributes.language_category.slug),
      required: true,
    },
    {
      name: "description",
      label: t(($) => $.models.attributes.language_category.description),
      type: "textarea",
    },
  ];
}
