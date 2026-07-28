import { Center, Loader, Stack, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetBannerOptions,
  adminGetBannerQueryKey,
  adminListBannersQueryKey,
  adminUpdateBannerMutation,
} from "@/client/@tanstack/react-query.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  bannerFormSchema,
  bannerToForm,
  bannerToInput,
  useBannerFields,
} from "@/components/admin/resources/banner";

export const Route = createFileRoute("/{-$locale}/admin/banners/$id")({
  component: EditBanner,
});

function EditBanner() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const fields = useBannerFields();
  const { id } = Route.useParams();
  const bannerId = Number(id);
  const [serverError, setServerError] = useState<string | null>(null);

  const { data, isLoading } = useQuery(
    adminGetBannerOptions({ path: { id: bannerId } }),
  );

  const backToList = () => navigate({ to: "/{-$locale}/admin/banners" });

  const mutation = useMutation({
    ...adminUpdateBannerMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminListBannersQueryKey() });
      queryClient.invalidateQueries({
        queryKey: adminGetBannerQueryKey({ path: { id: bannerId } }),
      });
      notifications.show({
        color: "green",
        message: t(($) => $.admin.crud.updated),
      });
      backToList();
    },
    onError: () => setServerError(t(($) => $.admin.crud.saveError)),
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
          schema={bannerFormSchema}
          defaultValues={bannerToForm(data)}
          onSubmit={async (values) => {
            setServerError(null);
            await mutation.mutateAsync({
              path: { id: bannerId },
              body: bannerToInput(values),
            });
          }}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
          serverError={serverError}
        />
      )}
    </Stack>
  );
}
