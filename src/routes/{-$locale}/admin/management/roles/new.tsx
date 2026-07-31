import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateRoleMutation,
  adminListRolesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zRoleInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyStaffRole, useStaffRoleFields } from "@/components/admin/resources/staffRole";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/new")({
  component: NewRole,
});

export function NewRole() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useStaffRoleFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/management/roles" });

  const mutation = useResourceMutation({
    mutation: adminCreateRoleMutation(),
    invalidate: [adminListRolesQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zRoleInput}
        defaultValues={emptyStaffRole}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
