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
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/new")({
  component: NewCourseLandingPage,
});

function NewCourseLandingPage(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseLandingPageFields();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/course_landing_pages" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListCourseLandingPagesQueryKey()],
    mutation: adminCreateCourseLandingPageMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zCourseLandingPageInput}
        defaultValues={emptyCourseLandingPage}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
