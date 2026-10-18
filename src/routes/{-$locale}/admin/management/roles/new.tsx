import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateRoleMutation,
  adminListRolesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { StaffRoleForm } from "@/components/admin/forms/StaffRoleForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { RoleInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/new")({
  component: NewRole,
});

const emptyStaffRole: RoleInput = {
  description: "",
  name: "",
};

function NewRole(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/management/roles" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListRolesQueryKey()],
    mutation: adminCreateRoleMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <StaffRoleForm
        defaultValues={emptyStaffRole}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
