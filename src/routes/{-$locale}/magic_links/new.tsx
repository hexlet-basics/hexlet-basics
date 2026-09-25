import { Alert, Anchor, Button, Card, Container, Stack, Text, Title } from "@mantine/core";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { createMagicLinkMutation } from "@/client/@tanstack/react-query.gen";
import { zEmailInput } from "@/client/zod.gen";
import { emailInputProps } from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";

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

  const mutation = useMutation(createMagicLinkMutation());

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

        <Card withBorder p="xl" w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}>
          {mutation.isSuccess ? (
            <Alert color="green">{t(($) => $.flash.magic_links.create.success)}</Alert>
          ) : (
            <Stack
              component="form"
              onSubmit={(event) => {
                event.preventDefault();
                form.handleSubmit();
              }}
            >
              {invalid && <Alert color="red">{t(($) => $.flash.magic_links.show.error)}</Alert>}
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
