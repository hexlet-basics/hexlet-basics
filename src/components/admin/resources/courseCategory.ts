import type { CourseCategory, CourseCategoryInput } from "@/client/types.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { useSchemaFields } from "@/components/admin/schemaFields";

// The course-category form model. Kept flat and all-string so it maps cleanly to
// the generated `zCourseCategoryInput` validator and the Mantine text fields; the
// nullable `description` is edited as an empty string and normalized on submit.
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

// Seed the form from an existing row; null columns become empty strings so the
// inputs stay controlled.
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

export function useCourseCategoryFields(): CrudFieldSpec<CourseCategoryFormValues>[] {
  return useSchemaFields<CourseCategoryFormValues>(zCourseCategoryInput, {
    namespace: "language_category",
    textarea: ["description"],
  });
}
