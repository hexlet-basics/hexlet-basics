import { Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetStaffMemberOptions,
  adminGetStaffMemberQueryKey,
  adminListStaffMembersQueryKey,
  adminUpdateStaffMemberMutation,
} from "@/client/@tanstack/react-query.gen";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { StaffMemberForm } from "@/components/admin/forms/StaffMemberForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { StaffMember, StaffMemberInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/$id")({
  component: EditStaffMember,
});

function staffMemberToForm(member: StaffMember): StaffMemberInput {
  return {
    allowedLocales: member.allowedLocales,
    roleId: member.roleId,
    userId: member.userId,
  };
}

function EditStaffMember(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = Route.useParams();
  const memberId = Number(id);

  const { data, isLoading } = useQuery(adminGetStaffMemberOptions({ path: { id: memberId } }));

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/management/staff_members" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [
      adminListStaffMembersQueryKey(),
      adminGetStaffMemberQueryKey({ path: { id: memberId } }),
    ],
    mutation: adminUpdateStaffMemberMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.updated),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <StaffMemberForm
          key={data.id}
          defaultValues={staffMemberToForm(data)}
          onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: memberId } })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
        />
      )}
    </Stack>
  );
}
