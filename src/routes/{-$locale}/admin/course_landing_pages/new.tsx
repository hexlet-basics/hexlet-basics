import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseLandingPageMutation,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zCourseLandingPageInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  emptyCourseLandingPage,
  useCourseLandingPageFields,
} from "@/components/admin/resources/courseLandingPage";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/new")({
  component: NewCourseLandingPage,
});

export function NewCourseLandingPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseLandingPageFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/course_landing_pages" });

  const mutation = useResourceMutation({
    mutation: adminCreateCourseLandingPageMutation(),
    invalidate: [adminListCourseLandingPagesQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zCourseLandingPageInput}
        defaultValues={emptyCourseLandingPage}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
