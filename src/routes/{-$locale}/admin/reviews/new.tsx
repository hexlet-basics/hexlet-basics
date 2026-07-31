import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateReviewMutation,
  adminListReviewsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zReviewInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyReview, useReviewFields } from "@/components/admin/resources/review";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/reviews/new")({
  component: NewReview,
});

export function NewReview() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useReviewFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/reviews" });

  const mutation = useResourceMutation({
    mutation: adminCreateReviewMutation(),
    invalidate: [adminListReviewsQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zReviewInput}
        defaultValues={emptyReview}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
