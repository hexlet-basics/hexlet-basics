import { Anchor, Box, Card, Center, Container, Stack, Title } from "@mantine/core";
import { modals } from "@mantine/modals";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  deleteAccountMutation,
  getCurrentUserQueryKey,
  getProfileOptions,
  getProfileQueryKey,
  updateProfileMutation,
} from "@/client/@tanstack/react-query.gen";
import type { ProfileInput } from "@/client/types.gen";
import { zProfileInput } from "@/client/zod.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { requireAuth } from "@/lib/auth";
import { submitForm, useAppForm } from "@/lib/form";
import { submitWithServerErrors } from "@/lib/server-validation";
import type { ReactElement } from "react";

// The signed-in user's profile, ported from legacy account/profiles#edit
// without the passkey section (ADR-0015). The name rules (40 characters, no
// special characters) come from the generated `zProfileInput`, the same
// constraints the API enforces.
export const Route = createFileRoute("/{-$locale}/account/profile/edit")({
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },

  loader: ({ context }) => context.queryClient.query(getProfileOptions()),

  head: ({ match }) => {
    const { t } = match.context.i18n;
    return {
      meta: [
        { title: t(($) => $.account.profiles.edit.title) },
        { content: t(($) => $.account.profiles.edit.meta.description), name: "description" },
      ],
    };
  },

  component: Edit,
});

function Edit(): ReactElement {
  const { t } = useTranslation();
  const profile = Route.useLoaderData();

  const update = useResourceMutation({
    errorMessage: t(($) => $.flash.account.profiles.update.error),
    fieldErrors: true,
    // The header shows the user's name from the current-user query.
    invalidate: [getProfileQueryKey(), getCurrentUserQueryKey()],
    mutation: updateProfileMutation(),
    successMessage: t(($) => $.flash.account.profiles.update.success),
  });

  const form = useAppForm({
    // The contract's nullable names as they are: TextField renders null as "".
    defaultValues: { firstName: profile.firstName, lastName: profile.lastName },
    validators: {
      onSubmit: zProfileInput,
      onSubmitAsync: submitWithServerErrors((value: ProfileInput) =>
        update.mutateAsync({ body: value }),
      ),
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.account.profiles.edit.title)}
        </Title>

        <Center w="100%">
          <Card shadow="sm" withBorder p="xl" w={{ base: "100%", sm: "50%", xs: "70%" }}>
            <Stack component="form" onSubmit={submitForm(form)}>
              <form.AppField name="firstName">
                {(field) => (
                  <field.TextField
                    label={t(($) => $.models.attributes.user.first_name)}
                    autoComplete="given-name"
                  />
                )}
              </form.AppField>

              <form.AppField name="lastName">
                {(field) => (
                  <field.TextField
                    label={t(($) => $.models.attributes.user.last_name)}
                    autoComplete="family-name"
                  />
                )}
              </form.AppField>

              <form.AppForm>
                <form.SubmitButton fullWidth mt="xl">
                  {t(($) => $.helpers.submit.save)}
                </form.SubmitButton>
              </form.AppForm>
            </Stack>

            <DeleteAccount />
          </Card>
        </Center>
      </Stack>
    </Container>
  );
}

// Account deletion behind the same yes/no confirmation legacy's useConfirmation
// showed. The API signs the user out, so the current-user cache is set to the
// guest before leaving: the root guard would otherwise keep the removed user.
function DeleteAccount(): ReactElement {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    ...deleteAccountMutation(),
    onError: () =>
      notifications.show({
        color: "red",
        message: t(($) => $.flash.account.profiles.destroy.error),
      }),
    onSuccess: () => {
      queryClient.setQueryData(getCurrentUserQueryKey(), { user: null });
      queryClient.removeQueries({ queryKey: getProfileQueryKey() });
      notifications.show({
        color: "green",
        message: t(($) => $.flash.account.profiles.destroy.success),
      });
      void navigate({ to: "/{-$locale}" });
    },
  });

  const confirm = (): string =>
    modals.openConfirmModal({
      centered: true,
      labels: {
        cancel: t(($) => $.common.boolean.no),
        confirm: t(($) => $.common.boolean.yes),
      },
      onConfirm: () => {
        mutation.mutate({});
      },
      title: t(($) => $.common.confirm),
    });

  return (
    <Box mt="xl">
      <Anchor
        component="button"
        type="button"
        c="red"
        onClick={() => {
          confirm();
        }}
      >
        {t(($) => $.account.profiles.edit.delete)}
      </Anchor>
    </Box>
  );
}
