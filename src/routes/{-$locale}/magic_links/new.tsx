import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import type { ReactElement } from "react";
import { z } from "zod";
import { createMagicLinkMutation } from "@/client/@tanstack/react-query.gen";
import { EmailRequestCard } from "@/components/EmailRequestCard";
import { useInvalidLinkNotice } from "@/hooks/useInvalidLinkNotice";

// Magic Link request form, ported from legacy magic_links/new. The API answers
// the same whether or not the email is registered, so the page does too: every
// submission ends on one "check your inbox" confirmation. `invalid` is set by
// the $token route when it refuses a followed link.
export const Route = createFileRoute("/{-$locale}/magic_links/new")({
  component: New,
  validateSearch: z.object({ invalid: z.boolean().optional() }),
});

function New(): ReactElement {
  const { t } = useTranslation();
  const { invalid } = Route.useSearch();

  const mutation = useMutation({
    ...createMagicLinkMutation(),
    onError: () => notifications.show({ color: "red", message: t(($) => $.common.errors.network) }),
  });

  useInvalidLinkNotice(
    invalid,
    t(($) => $.flash.magic_links.show.error),
  );

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
