import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetBannerOptions,
  adminGetBannerQueryKey,
  adminListBannersQueryKey,
  adminUpdateBannerMutation,
} from "@/client/@tanstack/react-query.gen";
import { zBannerInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { bannerToForm, useBannerFields } from "@/components/admin/resources/banner";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/banners/$id")({
  component: EditBanner,
});

function EditBanner() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useBannerFields();
  const { id } = Route.useParams();
  const bannerId = Number(id);

  const { data, isLoading } = useQuery(adminGetBannerOptions({ path: { id: bannerId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/banners" });

  const mutation = useResourceMutation({
    mutation: adminUpdateBannerMutation(),
    invalidate: [adminListBannersQueryKey(), adminGetBannerQueryKey({ path: { id: bannerId } })],
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
          schema={zBannerInput}
          defaultValues={bannerToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: bannerId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
