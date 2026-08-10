import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseMutation,
  adminListCoursesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zCourseInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyCourse, useCourseFields } from "@/components/admin/resources/course";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/courses/new")({
  component: NewCourse,
});

export function NewCourse() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/courses" });

  const mutation = useResourceMutation({
    mutation: adminCreateCourseMutation(),
    invalidate: [adminListCoursesQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zCourseInput}
        defaultValues={emptyCourse}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
