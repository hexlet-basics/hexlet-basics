import type { Course, CourseInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { courseFields } from "@/components/admin/generated/course";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Course form values ARE the generated request body (CourseInput).
// RepositoryUrl is derived from the slug on read and ignored on write, and
// CoverAttachmentId is deferred with the attachment wave — both stay null and
// Get no field.

export const emptyCourse: CourseInput = {
  slug: "",
  learnAs: null,
  hexletProgramLandingPage: "",
  repositoryUrl: null,
  coverAttachmentId: null,
};

export function courseToForm(course: Course): CourseInput {
  return {
    slug: course.slug,
    learnAs: course.learnAs,
    hexletProgramLandingPage: course.hexletProgramLandingPage ?? "",
    repositoryUrl: null,
    coverAttachmentId: null,
  };
}

// Fields come from the contract (generated/course.ts); what is listed here is
// Only what the schema cannot say.
export function useCourseFields(): CrudFieldSpec<CourseInput>[] {
  return useGeneratedFields<CourseInput>(courseFields, {
    set: { slug: { required: true } },
  });
}
