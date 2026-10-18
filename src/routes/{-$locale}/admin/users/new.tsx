import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateUserMutation,
  adminListUsersQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { UserForm } from "@/components/admin/forms/UserForm";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { UserInput } from "@/client/types.gen";
import type { ReactElement } from "react";

export const Route = createFileRoute("/{-$locale}/admin/users/new")({
  component: NewUser,
});

const emptyUser: UserInput = {
  admin: false,
  email: "",
  firstName: "",
  lastName: "",
};

function NewUser(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const backToList = (): void => {
    void navigate({ to: "/{-$locale}/admin/users" });
  };

  const mutation = useResourceMutation({
    errorMessage: t(($) => $.admin.crud.saveError),
    fieldErrors: true,
    invalidate: [adminListUsersQueryKey()],
    mutation: adminCreateUserMutation(),
    onDone: backToList,
    successMessage: t(($) => $.admin.crud.created),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <UserForm
        defaultValues={emptyUser}
        onSubmit={(values) => mutation.mutateAsync({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
      />
    </Stack>
  );
}
