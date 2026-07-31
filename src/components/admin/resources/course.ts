import { useTranslation } from "react-i18next";
import type { Course, CourseInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Course form values ARE the generated request body (CourseInput).
// repositoryUrl is derived from the slug on read and ignored on write, and
// coverAttachmentId is deferred with the attachment wave — both stay null and
// get no field.

export const emptyCourse: CourseInput = {
  slug: "",
  learnAs: null,
  progress: "draft",
  hexletProgramLandingPage: "",
  repositoryUrl: null,
  coverAttachmentId: null,
};

export function courseToForm(course: Course): CourseInput {
  return {
    slug: course.slug,
    learnAs: course.learnAs,
    progress: course.progress,
    hexletProgramLandingPage: course.hexletProgramLandingPage ?? "",
    repositoryUrl: null,
    coverAttachmentId: null,
  };
}

export function useCourseFields(): CrudFieldSpec<CourseInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "slug",
      label: t(($) => $.models.attributes.language.slug),
      required: true,
    },
    {
      name: "learnAs",
      label: t(($) => $.models.attributes.language.learn_as),
      type: "select",
      options: [
        {
          value: "first_language",
          label: t(($) => $.models.attributes.language["learn_as/values"].first_language),
        },
        {
          value: "second_language",
          label: t(($) => $.models.attributes.language["learn_as/values"].second_language),
        },
      ],
    },
    {
      name: "progress",
      label: t(($) => $.models.attributes.language.progress),
      type: "select",
      options: [
        {
          value: "draft",
          label: t(($) => $.models.attributes.language["progress/values"].draft),
        },
        {
          value: "in_development",
          label: t(($) => $.models.attributes.language["progress/values"].in_development),
        },
        {
          value: "completed",
          label: t(($) => $.models.attributes.language["progress/values"].completed),
        },
      ],
    },
    {
      name: "hexletProgramLandingPage",
      label: t(($) => $.models.attributes.language.hexlet_program_landing_page),
    },
  ];
}
