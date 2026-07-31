import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseLandingPageOptions,
  adminGetCourseLandingPageQueryKey,
  adminListCourseLandingPagesQueryKey,
  adminUpdateCourseLandingPageMutation,
} from "@/client/@tanstack/react-query.gen";
import { zCourseLandingPageInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
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
          schema={zCourseLandingPageInput}
          defaultValues={courseLandingPageToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: pageId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
