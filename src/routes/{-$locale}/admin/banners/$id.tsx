import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetBannerOptions,
  adminGetBannerQueryKey,
  adminListBannersQueryKey,
  adminUpdateBannerMutation,
} from "@/client/@tanstack/react-query.gen";
import { BannerStatePanel } from "@/components/admin/BannerStatePanel";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { BannerForm } from "@/components/admin/forms/BannerForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/banners/$id")({
  component: EditBanner,
});

function EditBanner(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = Route.useParams();
  const bannerId = Number(id);

  const { data, isLoading } = useQuery(adminGetBannerOptions({ path: { id: bannerId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/banners" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListBannersQueryKey(), adminGetBannerQueryKey({ path: { id: bannerId } })],
    mutation: adminUpdateBannerMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <>
          <BannerForm
            key={data.id}
            defaultValues={{
              background: data.background,
              body: data.body,
              finishesAt: data.finishesAt,
              locale: data.locale,
              startsAt: data.startsAt,
              url: data.url ?? "",
            }}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: bannerId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <BannerStatePanel banner={data} />
        </>
      )}
    </Stack>
  );
}
