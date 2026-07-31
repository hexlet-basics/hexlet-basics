import { Center, Loader, Stack, Title } from "@mantine/core";
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
import { CrudForm } from "@/components/admin/CrudForm";
import { QnaPanel } from "@/components/admin/QnaPanel";
import {
  courseLandingPageToForm,
  useCourseLandingPageFields,
} from "@/components/admin/resources/courseLandingPage";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/language_landing_pages/$id")({
  component: EditCourseLandingPage,
});

function EditCourseLandingPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useCourseLandingPageFields();
  const { id } = Route.useParams();
  const pageId = Number(id);

  const { data, isLoading } = useQuery(adminGetCourseLandingPageOptions({ path: { id: pageId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/language_landing_pages" });

  const mutation = useResourceMutation({
    mutation: adminUpdateCourseLandingPageMutation(),
    invalidate: [
      adminListCourseLandingPagesQueryKey(),
      adminGetCourseLandingPageQueryKey({ path: { id: pageId } }),
    ],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  // Nested FAQ (legacy `admin/language_landing_pages/:id/qna_items`).
  const qnaPath = { path: { landingPageId: pageId } };
  const qna = useQuery(adminListLandingPageQnaItemsOptions(qnaPath));
  const qnaInvalidate = [adminListLandingPageQnaItemsQueryKey(qnaPath)];
  const createQna = useResourceMutation({
    mutation: adminCreateLandingPageQnaItemMutation(),
    invalidate: qnaInvalidate,
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
  });
  const updateQna = useResourceMutation({
    mutation: adminUpdateLandingPageQnaItemMutation(),
    invalidate: qnaInvalidate,
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
  });
  const deleteQna = useResourceMutation({
    mutation: adminDeleteLandingPageQnaItemMutation(),
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
            key={data.id}
            fields={fields}
            schema={zCourseLandingPageInput}
            defaultValues={courseLandingPageToForm(data)}
            onSubmit={(values) => mutation.mutate({ path: { id: pageId }, body: values })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
            isPending={mutation.isPending}
          />
          <QnaPanel
            items={qna.data ?? []}
            onCreate={(input) => createQna.mutate({ ...qnaPath, body: input })}
            onUpdate={(itemId, input) =>
              updateQna.mutate({ path: { landingPageId: pageId, id: itemId }, body: input })
            }
            onDelete={(itemId) => deleteQna.mutate({ path: { landingPageId: pageId, id: itemId } })}
            isPending={createQna.isPending || updateQna.isPending || deleteQna.isPending}
          />
        </>
      )}
    </Stack>
  );
}
