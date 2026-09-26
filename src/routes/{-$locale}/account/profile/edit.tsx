import { Anchor, Box, Button, Card, Center, Container, Stack, Title } from "@mantine/core";
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
import { zProfileInput } from "@/client/zod.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { requireAuth } from "@/lib/auth";
import { useAppForm } from "@/lib/form";

// The signed-in user's profile, ported from legacy account/profiles#edit
// without the passkey section (ADR-0015). The name rules (40 characters, no
// special characters) come from the generated `zProfileInput`, the same
// constraints the API enforces.
export const Route = createFileRoute("/{-$locale}/account/profile/edit")({
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },
  loader: ({ context }) => context.queryClient.ensureQueryData(getProfileOptions()),
  head: ({ match }) => {
    const { t } = match.context.i18n;
    return {
      meta: [
        { title: t(($) => $.account.profiles.edit.title) },
        { name: "description", content: t(($) => $.account.profiles.edit.meta.description) },
      ],
    };
  },
  component: Edit,
});

function Edit() {
  const { t } = useTranslation();
  const profile = Route.useLoaderData();

  const update = useResourceMutation({
    mutation: updateProfileMutation(),
    // The header shows the user's name from the current-user query.
    invalidate: [getProfileQueryKey(), getCurrentUserQueryKey()],
    successMessage: t(($) => $.flash.account.profiles.update.success),
    errorMessage: t(($) => $.flash.account.profiles.update.error),
  });

  const form = useAppForm({
    // The contract's nullable names as they are: TextField renders null as "".
    defaultValues: { firstName: profile.firstName, lastName: profile.lastName },
    validators: { onSubmit: zProfileInput },
    onSubmit: async ({ value }) => {
      await update.mutateAsync({ body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.account.profiles.edit.title)}
        </Title>

        <Center w="100%">
          <Card shadow="sm" withBorder p="xl" w={{ base: "100%", xs: "70%", sm: "50%" }}>
            <Stack
              component="form"
              onSubmit={(event) => {
                event.preventDefault();
                form.handleSubmit();
              }}
            >
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

              <Button type="submit" fullWidth mt="xl" loading={update.isPending}>
                {t(($) => $.helpers.submit.save)}
              </Button>
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
function DeleteAccount() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    ...deleteAccountMutation(),
    onSuccess: () => {
      queryClient.setQueryData(getCurrentUserQueryKey(), { user: null });
      queryClient.removeQueries({ queryKey: getProfileQueryKey() });
      notifications.show({
        color: "green",
        message: t(($) => $.flash.account.profiles.destroy.success),
      });
      navigate({ to: "/{-$locale}" });
    },
    onError: () =>
      notifications.show({
        color: "red",
        message: t(($) => $.flash.account.profiles.destroy.error),
      }),
  });

  const confirm = () =>
    modals.openConfirmModal({
      centered: true,
      title: t(($) => $.common.confirm),
      labels: {
        confirm: t(($) => $.common.boolean.yes),
        cancel: t(($) => $.common.boolean.no),
      },
      onConfirm: () => mutation.mutate({}),
    });

  return (
    <Box mt="xl">
      <Anchor component="button" type="button" c="red" onClick={confirm}>
        {t(($) => $.account.profiles.edit.delete)}
      </Anchor>
    </Box>
  );
}
