import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateBannerMutation,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zBannerInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyBanner, useBannerFields } from "@/components/admin/resources/banner";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/banners/new")({
  component: NewBanner,
});

function NewBanner() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useBannerFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/banners" });

  const mutation = useResourceMutation({
    mutation: adminCreateBannerMutation(),
    invalidate: [adminListBannersQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zBannerInput}
        defaultValues={emptyBanner}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
