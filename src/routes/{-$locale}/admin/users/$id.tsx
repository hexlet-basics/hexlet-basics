import { Center, Loader, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminGetUserOptions,
  adminGetUserQueryKey,
  adminListUsersQueryKey,
  adminUpdateUserMutation,
} from "@/client/@tanstack/react-query.gen";
import { zUserInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { userToForm, useUserFields } from "@/components/admin/resources/user";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/users/$id")({
  component: EditUser,
});

function EditUser() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useUserFields();
  const { id } = Route.useParams();
  const userId = Number(id);

  const { data, isLoading } = useQuery(adminGetUserOptions({ path: { id: userId } }));

  const backToList = () => navigate({ to: "/{-$locale}/admin/users" });

  const mutation = useResourceMutation({
    mutation: adminUpdateUserMutation(),
    invalidate: [adminListUsersQueryKey(), adminGetUserQueryKey({ path: { id: userId } })],
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
          schema={zUserInput}
          defaultValues={userToForm(data)}
          onSubmit={(values) => mutation.mutate({ path: { id: userId }, body: values })}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
        />
      )}
    </Stack>
  );
}
