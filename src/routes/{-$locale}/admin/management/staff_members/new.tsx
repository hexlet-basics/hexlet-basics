import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateStaffMemberMutation,
  adminListStaffMembersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zStaffMemberInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyStaffMember, useStaffMemberFields } from "@/components/admin/resources/staffMember";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/new")({
  component: NewStaffMember,
});

export function NewStaffMember() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useStaffMemberFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/management/staff_members" });

  const mutation = useResourceMutation({
    mutation: adminCreateStaffMemberMutation(),
    invalidate: [adminListStaffMembersQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zStaffMemberInput}
        defaultValues={emptyStaffMember}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
