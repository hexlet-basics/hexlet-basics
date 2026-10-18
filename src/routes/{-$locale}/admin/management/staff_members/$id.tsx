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
import { zStaffMemberInput } from "@/client/zod.gen";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CrudForm } from "@/components/admin/CrudForm";
import { staffMemberToForm, useStaffMemberFields } from "@/components/admin/resources/staffMember";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/management/staff_members/$id")({
  component: EditStaffMember,
});

function EditStaffMember() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useStaffMemberFields();
  const { id } = Route.useParams();
  const memberId = Number(id);

  const { data, isLoading } = useQuery(adminGetStaffMemberOptions({ path: { id: memberId } }));

  const backToList = () => {
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
        <CrudForm
          key={data.id}
          fields={fields}
          schema={zStaffMemberInput}
          defaultValues={staffMemberToForm(data)}
          onSubmit={(values) => mutation.mutateAsync({ body: values, path: { id: memberId } })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
        />
      )}
    </Stack>
  );
}
