import { useTranslation } from "react-i18next";
import type { CourseCategory, CourseCategoryInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// The course-category form model. Kept flat and all-string so it maps cleanly to
// the generated `zCourseCategoryInput` validator and the Mantine text fields; the
// nullable `description` is edited as an empty string and normalized on submit.
// A `type` (not `interface`) so it satisfies the engine's
// `Record<string, unknown>` constraint — interfaces have no implicit index
// signature and would be rejected there.
export type CourseCategoryFormValues = {
  name: string;
  header: string;
  slug: string;
  description: string;
};

export const emptyCourseCategory: CourseCategoryFormValues = {
  name: "",
  header: "",
  slug: "",
  description: "",
};

// Seed the form from an existing row (edit screen). Null columns become empty
// strings so the inputs stay controlled.
export function courseCategoryToForm(
  category: CourseCategory,
): CourseCategoryFormValues {
  return {
    name: category.name ?? "",
    header: category.header ?? "",
    slug: category.slug ?? "",
    description: category.description ?? "",
  };
}

// Map form values to the request body. An empty description is sent as null so
// the column is cleared rather than storing a blank string.
export function courseCategoryToInput(
  values: CourseCategoryFormValues,
): CourseCategoryInput {
  return { ...values, description: values.description || null };
}

// Field descriptors driving CrudForm. A hook so labels resolve through i18n;
// shared by the create and edit screens.
export function useCourseCategoryFields(): CrudFieldSpec<CourseCategoryFormValues>[] {
  const { t } = useTranslation();
  return [
    {
      name: "name",
      label: t(($) => $.models.attributes.language_category.name),
      required: true,
      autoFocus: true,
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
