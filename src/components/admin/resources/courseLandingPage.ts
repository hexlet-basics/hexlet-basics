import type { CourseLandingPage, CourseLandingPageInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { courseLandingPageFields } from "@/components/admin/generated/courseLandingPage";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

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

// Fields come from the contract (generated/courseLandingPage.ts); what is listed here is
// only what the schema cannot say.
export function useCourseLandingPageFields(): CrudFieldSpec<CourseLandingPageInput>[] {
  return useGeneratedFields<CourseLandingPageInput>(courseLandingPageFields, {
    // Searchable pickers for the redirect target are a later refinement.
    hide: ["landingPageToRedirectId"],
    order: [
      "courseId",
      "name",
      "slug",
      "main",
      "listed",
      "footer",
      "footerName",
      "order",
      "header",
      "description",
      "metaTitle",
      "metaDescription",
      "usedInHeader",
      "usedInDescription",
      "outcomesHeader",
      "outcomesDescription",
    ],
    // The schema allows null for these, the form insists.
    set: {
      name: { required: true },
      slug: { required: true },
      description: { type: "textarea" },
      metaDescription: { type: "textarea" },
      usedInDescription: { type: "textarea" },
      outcomesDescription: { type: "textarea" },
    },
  });
}
