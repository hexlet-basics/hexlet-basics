import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateCategoryQnaItemMutation,
  adminDeleteCategoryQnaItemMutation,
  adminGetCourseCategoryOptions,
  adminGetCourseCategoryQueryKey,
  adminListCategoryQnaItemsOptions,
  adminListCategoryQnaItemsQueryKey,
  adminListCourseCategoriesQueryKey,
  adminUpdateCategoryQnaItemMutation,
  adminUpdateCourseCategoryMutation,
} from "@/client/@tanstack/react-query.gen";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CourseCategoryForm } from "@/components/admin/forms/CourseCategoryForm";
import { QnaPanel } from "@/components/admin/QnaPanel";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { CourseCategory, CourseCategoryInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_categories/$id")({
  component: EditCourseCategory,
});

function courseCategoryToForm(category: CourseCategory): CourseCategoryInput {
  return {
    description: category.description,
    header: category.header ?? "",
    name: category.name ?? "",
    slug: category.slug ?? "",
  };
}

function EditCourseCategory(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = Route.useParams();
  const categoryId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseCategoryOptions({ path: { id: categoryId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/course_categories" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [
      adminListCourseCategoriesQueryKey(),
      adminGetCourseCategoryQueryKey({ path: { id: categoryId } }),
    ],
    mutation: adminUpdateCourseCategoryMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  // Nested FAQ (legacy `admin/course_categories/:id/qna_items`).
  const qnaPath = { path: { categoryId } };
  const qna = useQuery(adminListCategoryQnaItemsOptions(qnaPath));
  const qnaInvalidate = [adminListCategoryQnaItemsQueryKey(qnaPath)];
  const createQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: qnaInvalidate,
    mutation: adminCreateCategoryQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.created),
  });
  const updateQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: qnaInvalidate,
    mutation: adminUpdateCategoryQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.updated),
  });
  const deleteQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.deleteError),
    invalidate: qnaInvalidate,
    mutation: adminDeleteCategoryQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.deleted),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <>
          <CourseCategoryForm
            // Remount per row so the form re-seeds its default values from the
            // loaded category instead of keeping a stale initial snapshot.
            key={data.id}
            defaultValues={courseCategoryToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: categoryId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <QnaPanel
            items={qna.data ?? []}
            onCreate={(input) => {
              createQna.mutate({ ...qnaPath, body: input });
            }}
            onUpdate={(itemId, input) => {
              updateQna.mutate({ body: input, path: { categoryId, id: itemId } });
            }}
            onDelete={(itemId) => {
              deleteQna.mutate({ path: { categoryId, id: itemId } });
            }}
            isPending={createQna.isPending || updateQna.isPending || deleteQna.isPending}
          />
        </>
      )}
    </Stack>
  );
}
