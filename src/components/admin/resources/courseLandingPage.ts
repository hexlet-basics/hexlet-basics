import { useTranslation } from "react-i18next";
import type { CourseLandingPage, CourseLandingPageInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Landing-page form values ARE the generated request body. courseId /
// landingPageToRedirectId stay raw ids for now (searchable pickers are a later
// refinement); outcomesImageAttachmentId is deferred with the whole
// cover/attachment wave, so the form always submits null for it.

export const emptyCourseLandingPage: CourseLandingPageInput = {
  courseId: 0,
  slug: "",
  name: "",
  main: false,
  listed: false,
  footer: false,
  footerName: "",
  state: "draft",
  order: "",
  landingPageToRedirectId: null,
  metaTitle: "",
  metaDescription: "",
  header: "",
  description: "",
  usedInHeader: "",
  usedInDescription: "",
  outcomesHeader: "",
  outcomesDescription: "",
  outcomesImageAttachmentId: null,
};

// Seed the form from an existing row: pick only the writable fields and
// coalesce nullable strings to "" so the inputs stay controlled.
export function courseLandingPageToForm(page: CourseLandingPage): CourseLandingPageInput {
  return {
    courseId: page.courseId,
    slug: page.slug,
    name: page.name,
    main: page.main ?? false,
    listed: page.listed ?? false,
    footer: page.footer ?? false,
    footerName: page.footerName ?? "",
    state: page.state,
    order: page.order ?? "",
    landingPageToRedirectId: page.landingPageToRedirectId,
    metaTitle: page.metaTitle,
    metaDescription: page.metaDescription,
    header: page.header,
    description: page.description,
    usedInHeader: page.usedInHeader ?? "",
    usedInDescription: page.usedInDescription ?? "",
    outcomesHeader: page.outcomesHeader ?? "",
    outcomesDescription: page.outcomesDescription ?? "",
    outcomesImageAttachmentId: null,
  };
}

export function useCourseLandingPageFields(): CrudFieldSpec<CourseLandingPageInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "courseId",
      label: t(($) => $.models.attributes.course_landing_page.course_id),
      type: "number",
      required: true,
    },
    {
      name: "name",
      label: t(($) => $.models.attributes.course_landing_page.name),
      required: true,
    },
    {
      name: "slug",
      label: t(($) => $.models.attributes.course_landing_page.slug),
      required: true,
    },
    {
      name: "state",
      label: t(($) => $.models.attributes.course_landing_page.state),
      type: "select",
      required: true,
      options: [
        {
          value: "draft",
          label: t(($) => $.models.attributes.course_landing_page["state/values"].draft),
        },
        {
          value: "published",
          label: t(($) => $.models.attributes.course_landing_page["state/values"].published),
        },
        {
          value: "archived",
          label: t(($) => $.models.attributes.course_landing_page["state/values"].archived),
        },
      ],
    },
    {
      name: "main",
      label: t(($) => $.models.attributes.course_landing_page.main),
      type: "checkbox",
    },
    {
      name: "listed",
      label: t(($) => $.models.attributes.course_landing_page.listed),
      type: "checkbox",
    },
    {
      name: "footer",
      label: t(($) => $.models.attributes.course_landing_page.footer),
      type: "checkbox",
    },
    {
      name: "footerName",
      label: t(($) => $.models.attributes.course_landing_page.footer_name),
    },
    {
      name: "order",
      label: t(($) => $.models.attributes.course_landing_page.order),
    },
    {
      name: "header",
      label: t(($) => $.models.attributes.course_landing_page.header),
    },
    {
      name: "description",
      label: t(($) => $.models.attributes.course_landing_page.description),
      type: "textarea",
    },
    {
      name: "metaTitle",
      label: t(($) => $.models.attributes.course_landing_page.meta_title),
    },
    {
      name: "metaDescription",
      label: t(($) => $.models.attributes.course_landing_page.meta_description),
      type: "textarea",
    },
    {
      name: "usedInHeader",
      label: t(($) => $.models.attributes.course_landing_page.used_in_header),
    },
    {
      name: "usedInDescription",
      label: t(($) => $.models.attributes.course_landing_page.used_in_description),
      type: "textarea",
    },
    {
      name: "outcomesHeader",
      label: t(($) => $.models.attributes.course_landing_page.outcomes_header),
    },
    {
      name: "outcomesDescription",
      label: t(($) => $.models.attributes.course_landing_page.outcomes_description),
      type: "textarea",
    },
  ];
}
