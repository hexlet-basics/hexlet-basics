import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createMagicLinkMutation } from "@/client/@tanstack/react-query.gen";
import { EmailRequestCard, useInvalidLinkNotice } from "@/components/EmailRequestCard";

// Magic Link request form, ported from legacy magic_links/new. The API answers
// the same whether or not the email is registered, so the page does too: every
// submission ends on one "check your inbox" confirmation. `invalid` is set by
// the $token route when it refuses a followed link.
export const Route = createFileRoute("/{-$locale}/magic_links/new")({
  validateSearch: z.object({ invalid: z.boolean().optional() }),
  component: New,
});

function New() {
  const { t } = useTranslation();
  const { invalid } = Route.useSearch();

  const mutation = useMutation({
    ...createMagicLinkMutation(),
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  // The refused-link notice is a flash, so it is shown once on arrival like the
  // other flashes; the id keeps StrictMode's double effect from stacking two.
  useEffect(() => {
    if (invalid) {
      notifications.show({
        id: "invalid-link",
        color: "red",
        message: t(($) => $.flash.magic_links.show.error),
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
          {t(($) => $.magic_links.new.title)}
        </Title>

        <EmailRequestCard
          mutation={mutation}
          success={t(($) => $.flash.magic_links.create.success)}
          submit={t(($) => $.magic_links.new.submit)}
        />

        <Text mt="xs">
          <Anchor component={Link} to="/session/new" fw="bold">
            {t(($) => $.magic_links.new.other_methods)}
          </Anchor>
        </Text>
      </Stack>
    </Container>
  );
}
