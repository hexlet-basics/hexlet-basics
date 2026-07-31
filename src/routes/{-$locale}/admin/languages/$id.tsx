import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseOptions,
  adminGetCourseQueryKey,
  adminListCoursesQueryKey,
  adminUpdateCourseMutation,
} from "@/client/@tanstack/react-query.gen";
import { zCourseInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { courseToForm, useCourseFields } from "@/components/admin/resources/course";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/languages/$id")({
  component: EditCourse,
});

function EditCourse() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseFields();
  const { id } = Route.useParams();
  const courseId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseOptions({ path: { id: courseId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/languages" });

  const mutation = useResourceMutation({
    mutation: adminUpdateCourseMutation(),
    invalidate: [adminListCoursesQueryKey(), adminGetCourseQueryKey({ path: { id: courseId } })],
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
          schema={zCourseInput}
          defaultValues={courseToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: courseId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
