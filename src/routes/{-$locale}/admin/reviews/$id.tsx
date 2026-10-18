import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetReviewOptions,
  adminGetReviewQueryKey,
  adminListReviewsQueryKey,
  adminUpdateReviewMutation,
} from "@/client/@tanstack/react-query.gen";
import { ReviewStatePanel } from "@/components/admin/ReviewStatePanel";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { ReviewForm } from "@/components/admin/forms/ReviewForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { Review, ReviewInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/reviews/$id")({
  component: EditReview,
});

function reviewToForm(review: Review): ReviewInput {
  return {
    body: review.body ?? "",
    courseId: review.courseId,
    firstName: review.firstName ?? "",
    lastName: review.lastName ?? "",
    pinned: review.pinned ?? false,
    userId: review.userId,
  };
}

function EditReview(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = Route.useParams();
  const reviewId = Number(id);

  const { data, isLoading } = useQuery(adminGetReviewOptions({ path: { id: reviewId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/reviews" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListReviewsQueryKey(), adminGetReviewQueryKey({ path: { id: reviewId } })],
    mutation: adminUpdateReviewMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <>
          <ReviewForm
            key={data.id}
            defaultValues={reviewToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: reviewId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <ReviewStatePanel review={data} />
        </>
      )}
    </Stack>
  );
}
