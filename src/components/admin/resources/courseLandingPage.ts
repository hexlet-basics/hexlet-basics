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
  description: "",
  footer: false,
  footerName: "",
  header: "",
  landingPageToRedirectId: null,
  listed: false,
  main: false,
  metaDescription: "",
  metaTitle: "",
  name: "",
  order: "",
  outcomesDescription: "",
  outcomesHeader: "",
  outcomesImageAttachmentId: null,
  slug: "",
  usedInDescription: "",
  usedInHeader: "",
};

// Seed the form from an existing row: pick only the writable fields and
// coalesce nullable strings to "" so the inputs stay controlled.
export function courseLandingPageToForm(page: CourseLandingPage): CourseLandingPageInput {
  return {
    courseId: page.courseId,
    description: page.description,
    footer: page.footer ?? false,
    footerName: page.footerName ?? "",
    header: page.header,
    landingPageToRedirectId: page.landingPageToRedirectId,
    listed: page.listed ?? false,
    main: page.main ?? false,
    metaDescription: page.metaDescription,
    metaTitle: page.metaTitle,
    name: page.name,
    order: page.order ?? "",
    outcomesDescription: page.outcomesDescription ?? "",
    outcomesHeader: page.outcomesHeader ?? "",
    outcomesImageAttachmentId: null,
    slug: page.slug,
    usedInDescription: page.usedInDescription ?? "",
    usedInHeader: page.usedInHeader ?? "",
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
      description: { type: "textarea" },
      metaDescription: { type: "textarea" },
      name: { required: true },
      outcomesDescription: { type: "textarea" },
      slug: { required: true },
      usedInDescription: { type: "textarea" },
    },
  });
}
