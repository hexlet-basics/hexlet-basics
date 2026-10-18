import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import type { ReactElement } from "react";
import { z } from "zod";
import { createPasswordReminderMutation } from "@/client/@tanstack/react-query.gen";
import { EmailRequestCard, useInvalidLinkNotice } from "@/components/EmailRequestCard";

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

  useInvalidLinkNotice(
    invalid,
    t(($) => $.flash.passwords.edit.error),
  );

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
