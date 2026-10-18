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
import FormSkeleton from "@/components/admin/FormSkeleton";
import { StaffRoleForm } from "@/components/admin/forms/StaffRoleForm";
import { RolePermissions } from "@/components/admin/RolePermissions";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { RoleInput, StaffRoleDetail } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/$id")({
  component: EditRole,
});

function staffRoleToForm(role: StaffRoleDetail): RoleInput {
  return {
    description: role.description ?? "",
    name: role.name,
  };
}

function EditRole(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
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
          <StaffRoleForm
            key={data.id}
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
