import { Anchor, Container, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createMagicLinkMutation } from "@/client/@tanstack/react-query.gen";
import { zEmailInput } from "@/client/zod.gen";
import { emailInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";
import type { ReactElement } from "react";

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

  // The refused-link notice is a flash, so it is shown once on arrival like the
  // other flashes; the id keeps StrictMode's double effect from stacking two.
  useEffect(() => {
    if (invalid) {
      notifications.show({
        color: "red",
        id: "invalid-link",
        message: t(($) => $.flash.magic_links.show.error),
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
          {t(($) => $.magic_links.new.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", lg: "50%", md: "70%", sm: "80%" }}>
          {mutation.isSuccess ? (
            <Alert color="green">{t(($) => $.flash.magic_links.create.success)}</Alert>
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
                {t(($) => $.magic_links.new.submit)}
              </Button>
            </Stack>
          )}
        </Card>

        <Text mt="xs">
          <Anchor component={Link} to="/session/new" fw="bold">
            {t(($) => $.magic_links.new.other_methods)}
          </Anchor>
        </Text>
      </Stack>
    </Container>
  );
}
