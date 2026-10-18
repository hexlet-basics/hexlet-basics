import type { Review, ReviewInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { reviewFields } from "@/components/admin/generated/review";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Review form values ARE the generated request body (ReviewInput). courseId and
// userId stay raw ids for now (a searchable picker is a later refinement); the
// input has no locale — the backend derives it, mirroring legacy.

export const emptyReview: ReviewInput = {
  body: "",
  courseId: null,
  firstName: "",
  lastName: "",
  pinned: false,
  userId: null,
};

// Seed the form from an existing row: pick only the writable fields and
// coalesce nullable strings to "" so the inputs stay controlled.
export function reviewToForm(review: Review): ReviewInput {
  return {
    body: review.body ?? "",
    courseId: review.courseId,
    firstName: review.firstName ?? "",
    lastName: review.lastName ?? "",
    pinned: review.pinned ?? false,
    userId: review.userId,
  };
}

// Fields come from the contract (generated/review.ts); what is listed here is
// only what the schema cannot say.
export function useReviewFields(): CrudFieldSpec<ReviewInput>[] {
  return useGeneratedFields<ReviewInput>(reviewFields, {
    order: ["pinned", "courseId", "userId", "firstName", "lastName", "body"],
    // The schema allows null for these, the form insists.
    set: {
      body: { required: true, type: "textarea" },
    },
  });
}
