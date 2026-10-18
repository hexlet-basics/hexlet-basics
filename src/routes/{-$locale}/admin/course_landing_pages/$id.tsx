import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateLandingPageQnaItemMutation,
  adminDeleteLandingPageQnaItemMutation,
  adminGetCourseLandingPageOptions,
  adminGetCourseLandingPageQueryKey,
  adminListCourseLandingPagesQueryKey,
  adminListLandingPageQnaItemsOptions,
  adminListLandingPageQnaItemsQueryKey,
  adminUpdateCourseLandingPageMutation,
  adminUpdateLandingPageQnaItemMutation,
} from "@/client/@tanstack/react-query.gen";
import { zCourseLandingPageInput } from "@/client/zod.gen";
import { CourseLandingPageStatePanel } from "@/components/admin/CourseLandingPageStatePanel";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CrudForm } from "@/components/admin/CrudForm";
import { QnaPanel } from "@/components/admin/QnaPanel";
import {
  courseLandingPageToForm,
  useCourseLandingPageFields,
} from "@/components/admin/resources/courseLandingPage";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/course_landing_pages/$id")({
  component: EditCourseLandingPage,
});

function EditCourseLandingPage(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseLandingPageFields();
  const { id } = Route.useParams();
  const pageId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseLandingPageOptions({ path: { id: pageId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/course_landing_pages" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [
      adminListCourseLandingPagesQueryKey(),
      adminGetCourseLandingPageQueryKey({ path: { id: pageId } }),
    ],
    mutation: adminUpdateCourseLandingPageMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  // Nested FAQ (legacy `admin/course_landing_pages/:id/qna_items`).
  const qnaPath = { path: { landingPageId: pageId } };
  const qna = useQuery(adminListLandingPageQnaItemsOptions(qnaPath));
  const qnaInvalidate = [adminListLandingPageQnaItemsQueryKey(qnaPath)];
  const createQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: qnaInvalidate,
    mutation: adminCreateLandingPageQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.created),
  });
  const updateQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    invalidate: qnaInvalidate,
    mutation: adminUpdateLandingPageQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.updated),
  });
  const deleteQna = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.deleteError),
    invalidate: qnaInvalidate,
    mutation: adminDeleteLandingPageQnaItemMutation(),
    successMessage: t(($) => $.admin.crud.deleted),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <>
          <CrudForm
            key={data.id}
            fields={fields}
            schema={zCourseLandingPageInput}
            defaultValues={courseLandingPageToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: pageId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <CourseLandingPageStatePanel page={data} />
          <QnaPanel
            items={qna.data ?? []}
            onCreate={(input) => {
              createQna.mutate({ ...qnaPath, body: input });
            }}
            onUpdate={(itemId, input) => {
              updateQna.mutate({ body: input, path: { id: itemId, landingPageId: pageId } });
            }}
            onDelete={(itemId) => {
              deleteQna.mutate({ path: { id: itemId, landingPageId: pageId } });
            }}
            isPending={createQna.isPending || updateQna.isPending || deleteQna.isPending}
          />
        </>
      )}
    </Stack>
  );
}
