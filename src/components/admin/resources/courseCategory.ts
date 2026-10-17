import type { CourseCategory, CourseCategoryInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { courseCategoryFields } from "@/components/admin/generated/courseCategory";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

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

// Fields come from the contract (generated/courseCategory.ts); what is listed here is
// only what the schema cannot say.
export function useCourseCategoryFields(): CrudFieldSpec<CourseCategoryInput>[] {
  return useGeneratedFields<CourseCategoryInput>(courseCategoryFields, {
    set: { description: { type: "textarea" } },
  });
}
