import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseOptions,
  adminGetCourseQueryKey,
  adminListCoursesQueryKey,
  adminUpdateCourseMutation,
} from "@/client/@tanstack/react-query.gen";
import { CourseReadinessPanel } from "@/components/admin/CourseReadinessPanel";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CourseForm } from "@/components/admin/forms/CourseForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { Course, CourseInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/courses/$id")({
  component: EditCourse,
});

function courseToForm(course: Course): CourseInput {
  return {
    coverAttachmentId: null,
    hexletProgramLandingPage: course.hexletProgramLandingPage ?? "",
    learnAs: course.learnAs,
    repositoryUrl: null,
    slug: course.slug,
  };
}

function EditCourse(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = Route.useParams();
  const courseId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseOptions({ path: { id: courseId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/courses" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListCoursesQueryKey(), adminGetCourseQueryKey({ path: { id: courseId } })],
    mutation: adminUpdateCourseMutation(),
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
          <CourseForm
            key={data.id}
            defaultValues={courseToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: courseId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <CourseReadinessPanel course={data} />
        </>
      )}
    </Stack>
  );
}
