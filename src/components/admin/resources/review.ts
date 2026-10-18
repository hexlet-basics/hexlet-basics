import type { Review, ReviewInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { reviewFields } from "@/components/admin/generated/review";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Review form values ARE the generated request body (ReviewInput). courseId and
// userId stay raw ids for now (a searchable picker is a later refinement); the
// input has no locale — the backend derives it, mirroring legacy.

export const emptyReview: ReviewInput = {
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
    pinned: review.pinned ?? false,
    courseId: review.courseId,
    userId: review.userId,
    body: review.body ?? "",
    firstName: review.firstName ?? "",
    lastName: review.lastName ?? "",
  };
}

// Fields come from the contract (generated/review.ts); what is listed here is
// only what the schema cannot say.
export function useReviewFields(): CrudFieldSpec<ReviewInput>[] {
  return useGeneratedFields<ReviewInput>(reviewFields, {
    order: ["pinned", "courseId", "userId", "firstName", "lastName", "body"],
    // The schema allows null for these, the form insists.
    set: {
      body: { type: "textarea", required: true },
    },
  });
}
