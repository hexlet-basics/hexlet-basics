import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateStaffMemberMutation,
  adminListStaffMembersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { StaffMemberForm } from "@/components/admin/forms/StaffMemberForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { StaffMemberInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/new")({
  component: NewStaffMember,
});

const emptyStaffMember: StaffMemberInput = {
  allowedLocales: [],
  roleId: 0,
  userId: 0,
};

function NewStaffMember(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/management/staff_members" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListStaffMembersQueryKey()],
    mutation: adminCreateStaffMemberMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <StaffMemberForm
        defaultValues={emptyStaffMember}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
