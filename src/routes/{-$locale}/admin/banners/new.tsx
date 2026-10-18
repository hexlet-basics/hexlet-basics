import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import type { BannerInput } from "@/client/types.gen";
import {
  adminCreateBannerMutation,
  adminListBannersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { BannerForm } from "@/components/admin/forms/BannerForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/banners/new")({
  component: NewBanner,
});

const emptyBanner: BannerInput = {
  background: "cta_gradient",
  body: "",
  finishesAt: null,
  locale: "en",
  startsAt: null,
  url: "",
};

function NewBanner(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/banners" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListBannersQueryKey()],
    mutation: adminCreateBannerMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <BannerForm
        defaultValues={emptyBanner}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
