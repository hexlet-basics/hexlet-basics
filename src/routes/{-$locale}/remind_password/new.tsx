import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createPasswordReminderMutation } from "@/client/@tanstack/react-query.gen";
import { EmailRequestCard, useInvalidLinkNotice } from "@/components/EmailRequestCard";

// Password Reset request form, ported from legacy remind_passwords/new. Unlike
// Legacy, an unknown email is no longer reported: the API answers the same for
// Any email, so every submission ends on one confirmation. `invalid` is set by
// The password/$token/edit route when it refuses a followed link.
export const Route = createFileRoute("/{-$locale}/remind_password/new")({
  validateSearch: z.object({ invalid: z.boolean().optional() }),
  component: New,
});

function New() {
  const { t } = useTranslation();
  const { invalid } = Route.useSearch();

  const mutation = useMutation({
    ...createPasswordReminderMutation(),
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  // The refused-link notice is a flash, so it is shown once on arrival like the
  // Other flashes; the id keeps StrictMode's double effect from stacking two.
  useEffect(() => {
    if (invalid) {
      notifications.show({
        id: "invalid-link",
        color: "red",
        message: t(($) => $.flash.passwords.edit.error),
      });
    }
  }, [invalid, t]);

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

        <EmailRequestCard
          mutation={mutation}
          success={t(($) => $.flash.remind_passwords.create.success)}
          submit={t(($) => $.helpers.submit.remind_password_form.create)}
        />

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
