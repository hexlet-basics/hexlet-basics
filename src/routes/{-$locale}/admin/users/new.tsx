import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateUserMutation,
  adminListUsersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zUserInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyUser, useUserFields } from "@/components/admin/resources/user";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/users/new")({
  component: NewUser,
});

export function NewUser() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useUserFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/users" });

  const mutation = useResourceMutation({
    mutation: adminCreateUserMutation(),
    invalidate: [adminListUsersQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zUserInput}
        defaultValues={emptyUser}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
