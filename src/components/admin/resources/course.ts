import type { Course, CourseInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { courseFields } from "@/components/admin/generated/course";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Course form values ARE the generated request body (CourseInput).
// repositoryUrl is derived from the slug on read and ignored on write, and
// coverAttachmentId is deferred with the attachment wave — both stay null and
// get no field.

export const emptyCourse: CourseInput = {
  coverAttachmentId: null,
  hexletProgramLandingPage: "",
  learnAs: null,
  repositoryUrl: null,
  slug: "",
};

export function courseToForm(course: Course): CourseInput {
  return {
    coverAttachmentId: null,
    hexletProgramLandingPage: course.hexletProgramLandingPage ?? "",
    learnAs: course.learnAs,
    repositoryUrl: null,
    slug: course.slug,
  };
}

// Fields come from the contract (generated/course.ts); what is listed here is
// only what the schema cannot say.
export function useCourseFields(): CrudFieldSpec<CourseInput>[] {
  return useGeneratedFields<CourseInput>(courseFields, {
    set: { slug: { required: true } },
  });
}
