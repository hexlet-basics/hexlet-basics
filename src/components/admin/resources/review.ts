import { useTranslation } from "react-i18next";
import type { Review, ReviewInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Review form values ARE the generated request body (ReviewInput). courseId and
// userId stay raw ids for now (a searchable picker is a later refinement); the
// input has no locale — the backend derives it, mirroring legacy.

export const emptyReview: ReviewInput = {
  state: "draft",
  pinned: false,
  courseId: null,
  userId: null,
  body: "",
  firstName: "",
  lastName: "",
};

// Seed the form from an existing row: pick only the writable fields and
// coalesce nullable strings to "" so the inputs stay controlled.
export function reviewToForm(review: Review): ReviewInput {
  return {
    state: review.state,
    pinned: review.pinned ?? false,
    courseId: review.courseId,
    userId: review.userId,
    body: review.body ?? "",
    firstName: review.firstName ?? "",
    lastName: review.lastName ?? "",
  };
}

export function useReviewFields(): CrudFieldSpec<ReviewInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "state",
      label: t(($) => $.models.attributes.review.state),
      type: "select",
      required: true,
      options: [
        {
          value: "draft",
          label: t(($) => $.models.attributes.review["state/values"].draft),
        },
        {
          value: "published",
          label: t(($) => $.models.attributes.review["state/values"].published),
        },
        {
          value: "archived",
          label: t(($) => $.models.attributes.review["state/values"].archived),
        },
      ],
    },
    {
      name: "pinned",
      label: t(($) => $.models.attributes.review.pinned),
      type: "checkbox",
    },
    {
      name: "courseId",
      label: t(($) => $.models.attributes.review.course_id),
      type: "number",
    },
    {
      name: "userId",
      label: t(($) => $.models.attributes.review.user_id),
      type: "number",
    },
    {
      name: "firstName",
      label: t(($) => $.models.attributes.review.first_name),
    },
    {
      name: "lastName",
      label: t(($) => $.models.attributes.review.last_name),
    },
    {
      name: "body",
      label: t(($) => $.models.attributes.review.body),
      type: "textarea",
      required: true,
    },
  ];
}
