import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseLandingPageMutation,
  adminListCourseLandingPagesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { CourseLandingPageForm } from "@/components/admin/forms/CourseLandingPageForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { CourseLandingPageInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/new")({
  component: NewCourseLandingPage,
});

const emptyCourseLandingPage: CourseLandingPageInput = {
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

function NewCourseLandingPage(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

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
      <CourseLandingPageForm
        defaultValues={emptyCourseLandingPage}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
