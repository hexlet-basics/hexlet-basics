import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetReviewOptions,
  adminGetReviewQueryKey,
  adminListReviewsQueryKey,
  adminUpdateReviewMutation,
} from "@/client/@tanstack/react-query.gen";
import { zReviewInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { reviewToForm, useReviewFields } from "@/components/admin/resources/review";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/reviews/$id")({
  component: EditReview,
});

function EditReview() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useReviewFields();
  const { id } = Route.useParams();
  const reviewId = Number(id);

  const { data, isLoading } = useQuery(adminGetReviewOptions({ path: { id: reviewId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/reviews" });

  const mutation = useResourceMutation({
    mutation: adminUpdateReviewMutation(),
    invalidate: [adminListReviewsQueryKey(), adminGetReviewQueryKey({ path: { id: reviewId } })],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <Center py="xl">
          <Loader />
        </Center>
      ) : (
        <CrudForm
          key={data.id}
          fields={fields}
          schema={zReviewInput}
          defaultValues={reviewToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: reviewId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
