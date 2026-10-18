import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetManagementUserOptions,
  adminGetManagementUserQueryKey,
  adminListManagementUsersQueryKey,
  adminUpdateManagementUserMutation,
} from "@/client/@tanstack/react-query.gen";
import { zUserInput } from "@/client/zod.gen";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CrudForm } from "@/components/admin/CrudForm";
import { useUserFields, userToForm } from "@/components/admin/resources/user";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/users/$id")({
  component: EditManagementUser,
});

function EditManagementUser(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useUserFields();
  const { id } = Route.useParams();
  const userId = Number(id);

  const { data, isLoading } = useQuery(adminGetManagementUserOptions({ path: { id: userId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/management/users" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [
      adminListManagementUsersQueryKey(),
      adminGetManagementUserQueryKey({ path: { id: userId } }),
    ],
    mutation: adminUpdateManagementUserMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <CrudForm
          key={data.id}
          fields={fields}
          schema={zUserInput}
          defaultValues={userToForm(data)}
          onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: userId } })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
        />
      )}
    </Stack>
  );
}
