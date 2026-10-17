import { Anchor, Button, Card, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createSessionMutation, getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import { zSessionInput } from "@/client/zod.gen";
import { currentPasswordInputProps, loginEmailInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";
import { safeRedirectPath } from "@/lib/safe-redirect";

// Login page, ported from legacy sessions/new. Email + password, with links to
// the Password Reset and Magic Link request forms; the legacy phone / passkey / GitHub options are
// deferred until their backends and routes land. Submits through the generated `createSession`
// mutation, which starts the session server-side; validation reuses the
// generated `zSessionInput` schema.
//
// `redirect` is where to go once signed in: `requireAuth` sends a guest here with
// the page they asked for (the OAuth consent screen is one, and loses its request
// id if the visit is forgotten). Only a path on this site is honoured, so the
// page cannot be used to bounce a signed-in user off to another origin.
const signInSearchSchema = z.object({
  redirect: z.string().transform(safeRedirectPath).optional().catch(undefined),
});

export const Route = createFileRoute("/{-$locale}/session/new")({
  validateSearch: signInSearchSchema,
  component: New,
});

function New() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { redirect } = Route.useSearch();

  const mutation = useMutation({
    ...createSessionMutation(),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: getCurrentUserQueryKey() });
      notifications.show({
        color: "green",
        message: t(($) => $.flash.sessions.create.success),
      });
      if (redirect) void navigate({ href: redirect });
      else void navigate({ to: "/{-$locale}" });
    },
    onError: () =>
      notifications.show({
        color: "red",
        message: t(($) => $.flash.sessions.create.error),
      }),
  });

  const form = useAppForm({
    defaultValues: { email: "", password: "" },
    validators: { onSubmit: zSessionInput },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({ body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.sessions.new.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}>
          <Stack
            component="form"
            onSubmit={(event) => {
              event.preventDefault();
              void form.handleSubmit();
            }}
          >
            <form.AppField name="email">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.email)}
                  required
                  {...loginEmailInputProps}
                />
              )}
            </form.AppField>

            <form.AppField name="password">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.password)}
                  required
                  {...currentPasswordInputProps}
                />
              )}
            </form.AppField>

            <Button type="submit" fullWidth loading={mutation.isPending}>
              {t(($) => $.helpers.submit.user_sign_in_form.create)}
            </Button>

            <Anchor component={Link} to="/remind_password/new" ta="center">
              {t(($) => $.sessions.new.forgot_password)}
            </Anchor>

            <Anchor component={Link} to="/magic_links/new" ta="center">
              {t(($) => $.sessions.new.sign_in_with_magic_link)}
            </Anchor>
          </Stack>
        </Card>

        <Text mt="xs">
          {t(($) => $.sessions.new.dont_have_account)}{" "}
          <Anchor component={Link} to="/users/new" fw="bold">
            {t(($) => $.sessions.new.register)}
          </Anchor>
        </Text>
      </Stack>
    </Container>
  );
}
