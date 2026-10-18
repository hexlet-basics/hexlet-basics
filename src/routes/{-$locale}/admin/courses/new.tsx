import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseMutation,
  adminListCoursesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { CourseForm } from "@/components/admin/forms/CourseForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { CourseInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/courses/new")({
  component: NewCourse,
});

const emptyCourse: CourseInput = {
  coverAttachmentId: null,
  hexletProgramLandingPage: "",
  learnAs: null,
  repositoryUrl: null,
  slug: "",
};

function NewCourse(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/courses" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListCoursesQueryKey()],
    mutation: adminCreateCourseMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CourseForm
        defaultValues={emptyCourse}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
