import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetRoleOptions,
  adminGetRoleQueryKey,
  adminListRolesQueryKey,
  adminUpdateRoleMutation,
} from "@/client/@tanstack/react-query.gen";
import { zRoleInput } from "@/client/zod.gen";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CrudForm } from "@/components/admin/CrudForm";
import { staffRoleToForm, useStaffRoleFields } from "@/components/admin/resources/staffRole";
import { RolePermissions } from "@/components/admin/RolePermissions";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/$id")({
  component: EditRole,
});

function EditRole(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useStaffRoleFields();
  const { id } = Route.useParams();
  const roleId = Number(id);

  const { data, isLoading } = useQuery(adminGetRoleOptions({ path: { id: roleId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/management/roles" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListRolesQueryKey(), adminGetRoleQueryKey({ path: { id: roleId } })],
    mutation: adminUpdateRoleMutation(),
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
          <CrudForm
            key={data.id}
            fields={fields}
            schema={zRoleInput}
            defaultValues={staffRoleToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: roleId } })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <RolePermissions roleId={roleId} />
        </>
      )}
    </Stack>
  );
}
