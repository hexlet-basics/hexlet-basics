import { Center, Loader, Stack, Title } from "@mantine/core";
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
import { zCourseCategoryInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { QnaPanel } from "@/components/admin/QnaPanel";
import {
  courseCategoryToForm,
  useCourseCategoryFields,
} from "@/components/admin/resources/courseCategory";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/course_categories/$id")({
  component: EditCourseCategory,
});

function EditCourseCategory() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseCategoryFields();
  const { id } = Route.useParams();
  const categoryId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseCategoryOptions({ path: { id: categoryId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/course_categories" });

  const mutation = useResourceMutation({
    mutation: adminUpdateCourseCategoryMutation(),
    invalidate: [
      adminListCourseCategoriesQueryKey(),
      adminGetCourseCategoryQueryKey({ path: { id: categoryId } }),
    ],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  // Nested FAQ (legacy `admin/course_categories/:id/qna_items`).
  const qnaPath = { path: { categoryId } };
  const qna = useQuery(adminListCategoryQnaItemsOptions(qnaPath));
  const qnaInvalidate = [adminListCategoryQnaItemsQueryKey(qnaPath)];
  const createQna = useResourceMutation({
    mutation: adminCreateCategoryQnaItemMutation(),
    invalidate: qnaInvalidate,
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
  });
  const updateQna = useResourceMutation({
    mutation: adminUpdateCategoryQnaItemMutation(),
    invalidate: qnaInvalidate,
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
  });
  const deleteQna = useResourceMutation({
    mutation: adminDeleteCategoryQnaItemMutation(),
    invalidate: qnaInvalidate,
    successMessage: t(($) => $.admin.crud.deleted),
    errorMessage: t(($) => $.admin.crud.deleteError),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <Center py="xl">
          <Loader />
        </Center>
      ) : (
        <>
          <CrudForm
            // Remount per row so the form re-seeds its default values from the
            // loaded category instead of keeping a stale initial snapshot.
            key={data.id}
            fields={fields}
            schema={zCourseCategoryInput}
            defaultValues={courseCategoryToForm(data)}
            onSubmit={(values) => mutation.mutate({ path: { id: categoryId }, body: values })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
            isPending={mutation.isPending}
          />
          <QnaPanel
            items={qna.data ?? []}
            onCreate={(input) => createQna.mutate({ ...qnaPath, body: input })}
            onUpdate={(itemId, input) =>
              updateQna.mutate({ path: { categoryId, id: itemId }, body: input })
            }
            onDelete={(itemId) => deleteQna.mutate({ path: { categoryId, id: itemId } })}
            isPending={createQna.isPending || updateQna.isPending || deleteQna.isPending}
          />
        </>
      )}
    </Stack>
  );
}
