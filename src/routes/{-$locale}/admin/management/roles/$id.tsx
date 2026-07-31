import { Center, Loader, Stack, Title } from "@mantine/core";
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
import { CrudForm } from "@/components/admin/CrudForm";
import { staffRoleToForm, useStaffRoleFields } from "@/components/admin/resources/staffRole";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/roles/$id")({
  component: EditRole,
});

function EditRole() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useStaffRoleFields();
  const { id } = Route.useParams();
  const roleId = Number(id);

  const { data, isLoading } = useQuery(adminGetRoleOptions({ path: { id: roleId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/management/roles" });

  const mutation = useResourceMutation({
    mutation: adminUpdateRoleMutation(),
    invalidate: [adminListRolesQueryKey(), adminGetRoleQueryKey({ path: { id: roleId } })],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
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
          schema={zRoleInput}
          defaultValues={staffRoleToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: roleId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
