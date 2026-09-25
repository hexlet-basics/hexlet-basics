import { Alert, Anchor, Button, Card, Container, Stack, Text, Title } from "@mantine/core";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createPasswordReminderMutation } from "@/client/@tanstack/react-query.gen";
import { zEmailInput } from "@/client/zod.gen";
import { emailInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";

// Password Reset request form, ported from legacy remind_passwords/new. Unlike
// legacy, an unknown email is no longer reported: the API answers the same for
// any email, so every submission ends on one confirmation. `invalid` is set by
// the password/$token/edit route when it refuses a followed link.
export const Route = createFileRoute("/{-$locale}/remind_password/new")({
  validateSearch: z.object({ invalid: z.boolean().optional() }),
  component: New,
});

function New() {
  const { t } = useTranslation();
  const { invalid } = Route.useSearch();

  const mutation = useMutation(createPasswordReminderMutation());

  const form = useAppForm({
    defaultValues: { email: "" },
    validators: { onSubmit: zEmailInput },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({ body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.remind_passwords.new.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}>
          {mutation.isSuccess ? (
            <Alert color="green">{t(($) => $.flash.remind_passwords.create.success)}</Alert>
          ) : (
            <Stack
              component="form"
              onSubmit={(event) => {
                event.preventDefault();
                form.handleSubmit();
              }}
            >
              {invalid && <Alert color="red">{t(($) => $.flash.passwords.edit.error)}</Alert>}
              {mutation.isError && <Alert color="red">{t(($) => $.common.errors.network)}</Alert>}

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
