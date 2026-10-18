import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateReviewMutation,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { ReviewForm } from "@/components/admin/forms/ReviewForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReviewInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/reviews/new")({
  component: NewReview,
});

const emptyReview: ReviewInput = {
  body: "",
  courseId: null,
  firstName: "",
  lastName: "",
  pinned: false,
  userId: null,
};

function NewReview(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/reviews" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListReviewsQueryKey()],
    mutation: adminCreateReviewMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <ReviewForm
        defaultValues={emptyReview}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
