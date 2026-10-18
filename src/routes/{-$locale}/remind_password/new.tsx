import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createPasswordReminderMutation } from "@/client/@tanstack/react-query.gen";
import { zEmailInput } from "@/client/zod.gen";
import { emailInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";
import type { ReactElement } from "react";

// Password Reset request form, ported from legacy remind_passwords/new. Unlike
// legacy, an unknown email is no longer reported: the API answers the same for
// any email, so every submission ends on one confirmation. `invalid` is set by
// the password/$token/edit route when it refuses a followed link.
export const Route = createFileRoute("/{-$locale}/remind_password/new")({
  component: New,
  validateSearch: z.object({ invalid: z.boolean().optional() }),
});

function New(): ReactElement {
  const { t } = useTranslation();
  const { invalid } = Route.useSearch();

  const mutation = useMutation({
    ...createPasswordReminderMutation(),
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  // The refused-link notice is a flash, so it is shown once on arrival like the
  // other flashes; the id keeps StrictMode's double effect from stacking two.
  useEffect(() => {
    if (invalid) {
      notifications.show({
        color: "red",
        id: "invalid-link",
        message: t(($) => $.flash.passwords.edit.error),
      });
    }
  }, [invalid, t]);

  const form = useAppForm({
    defaultValues: { email: "" },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({ body: value });
    },
    validators: { onSubmit: zEmailInput },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.remind_passwords.new.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", lg: "50%", md: "70%", sm: "80%" }}>
          {mutation.isSuccess ? (
            <Alert color="green">{t(($) => $.flash.remind_passwords.create.success)}</Alert>
          ) : (
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
                    {...emailInputProps}
                  />
                )}
              </form.AppField>

              <Button type="submit" fullWidth loading={mutation.isPending}>
                {t(($) => $.helpers.submit.remind_password_form.create)}
              </Button>
            </Stack>
          )}
        </Card>

        <Text mt="xs">
          {t(($) => $.remind_passwords.new.trying_to_login)}{" "}
          <Anchor component={Link} to="/session/new" fw="bold">
            {t(($) => $.remind_passwords.new.login)}
          </Anchor>
        </Text>
      </Stack>
    </Container>
  );
}
