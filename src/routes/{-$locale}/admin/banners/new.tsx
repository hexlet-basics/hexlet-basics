import { Stack, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminCreateBannerMutation,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  bannerFormSchema,
  bannerToInput,
  emptyBanner,
  useBannerFields,
} from "@/components/admin/resources/banner";

export const Route = createFileRoute("/{-$locale}/admin/banners/new")({
  component: NewBanner,
});

function NewBanner() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const fields = useBannerFields();
  const [serverError, setServerError] = useState<string | null>(null);

  const backToList = () => navigate({ to: "/{-$locale}/admin/banners" });

  const mutation = useMutation({
    ...adminCreateBannerMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminListBannersQueryKey() });
      notifications.show({
        color: "green",
        message: t(($) => $.admin.crud.created),
      });
      backToList();
    },
    onError: () => setServerError(t(($) => $.admin.crud.saveError)),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={bannerFormSchema}
        defaultValues={emptyBanner}
        onSubmit={async (values) => {
          setServerError(null);
          await mutation.mutateAsync({ body: bannerToInput(values) });
        }}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
        serverError={serverError}
      />
    </Stack>
  );
}
