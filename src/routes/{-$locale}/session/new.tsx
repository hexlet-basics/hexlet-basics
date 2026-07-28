import {
  Alert,
  Anchor,
  Button,
  Card,
  Container,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  createSessionMutation,
  getCurrentUserQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zSessionInput } from "@/client/zod.gen";
import {
  currentPasswordInputProps,
  loginEmailInputProps,
} from "@/lib/authFieldProps";
import { useAppForm } from "@/lib/form";

// Login page, ported from legacy sessions/new. Email + password only; the
// legacy phone / magic-link / passkey / GitHub options are deferred until their
// backends and routes land. Submits through the generated `createSession`
// mutation, which sets the JWT cookie server-side; validation reuses the
// generated `zSessionInput` schema.
export const Route = createFileRoute("/{-$locale}/session/new")({
  component: New,
});

function New() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [serverError, setServerError] = useState<string | null>(null);

  const mutation = useMutation({
    ...createSessionMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getCurrentUserQueryKey() });
      navigate({ to: "/{-$locale}" });
    },
    onError: () => setServerError(t(($) => $.flash.sessions.create.error)),
  });

  const form = useAppForm({
    defaultValues: { email: "", password: "" },
    validators: { onSubmit: zSessionInput },
    onSubmit: async ({ value }) => {
      setServerError(null);
      await mutation.mutateAsync({ body: value });
    },
  });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.sessions.new.title)}
        </Title>

        <Card
          withBorder
          p="xl"
          w={{ base: "100%", sm: "80%", md: "70%", lg: "50%" }}
        >
          <Stack
            component="form"
            onSubmit={(event) => {
              event.preventDefault();
              form.handleSubmit();
            }}
          >
            {serverError && <Alert color="red">{serverError}</Alert>}

            <form.AppField name="email">
              {(field) => (
                <field.TextField
                  label={t(($) => $.models.attributes.user.email)}
                  required
                  autoFocus
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
